package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	reporter "reporter/pkg/reporter/service"
	timer "reporter/pkg/timer/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/xuri/excelize/v2"
)

const (
	// Random Image API
	riApiKey = "IJDHKmA5qlYGvBjD/CyPRw==3CIIQF9Nmym5mVwB"
	riUrl    = "https://api.api-ninjas.com/v1/randomimage"
	riAccept = "image/jpg"

	// FishText API
	ftNumber = 10
	ftUrl    = "https://fish-text.ru/get"
	ftType   = "sentence"

	// Timer gRPC
	grpcPort        = 50052
	address         = "timer-service:50051"
	defaultLocation = "Asia/Novosibirsk"

	url = "http://localhost:3000/download"
)

type ftResponse struct {
	Status    string `json:"status"`
	Text      string `json:"text"`
	ErrorCode int    `json:"errorCode"`
}

type server struct {
	reporter.UnimplementedReporterServer
}

func (s *server) GetReport(_ *reporter.GetRequest, stream reporter.Reporter_GetReportServer) error {
	log.Println("1 new connection")

	timech := make(chan string, 1)
	textch := make(chan string, 1)
	imagech := make(chan []byte, 1)
	errch := make(chan error, 3)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func(ch chan<- string, ech chan<- error) {
		defer wg.Done()
		t, err := getLocalTime()
		if err != nil {
			log.Println("failed to get local time")
			ech <- err
		}
		ch <- t
		close(ch)
	}(timech, errch)

	go func(ch chan<- []byte, ech chan<- error) {
		defer wg.Done()
		img, err := getRandomImage()
		if err != nil {
			log.Println("failed to get random image")
			ech <- err
		}
		ch <- img
		close(ch)
	}(imagech, errch)

	go func(ch chan<- string, ech chan<- error) {
		defer wg.Done()
		text, err := getFishText()
		if err != nil {
			log.Println("failed to get fish text")
			ech <- err
		}
		ch <- text
		close(ch)
	}(textch, errch)

	wg.Wait()

	if len(errch) != 0 {
		log.Println("error: ", <-errch)
		close(errch)
		return fmt.Errorf("failed to get report")
	}

	t := <-timech
	buf, err := getReport(t, <-imagech, <-textch, stream)
	if err != nil {
		log.Println("failed to generate report")
		return fmt.Errorf("failed to get report")

	}
	fileName := fmt.Sprintf("%s.xlsx", t)
	contentDisposition := fmt.Sprintf("attachment; filename=%s", fileName)
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", contentDisposition)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")

		_, err = io.Copy(w, buf)
		if err != nil {
			log.Println("failed to copy file")
		}
	})

	go func() {
		log.Fatal(http.ListenAndServe(":3000", nil))
	}()

	err = stream.Send(&reporter.GetResponse{
		Res: &reporter.GetResponse_Url{
			Url: url,
		},
	})
	return nil
}

func getReport(time string, image []byte, text string, stream reporter.Reporter_GetReportServer) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Print("failed to close report")
			panic(err)
		}
	}()

	enable, disable := true, false
	pic := &excelize.Picture{
		Extension: ".jpg",
		File:      image,
		Format: &excelize.GraphicOptions{
			PrintObject:     &enable,
			LockAspectRatio: false,
			OffsetX:         15,
			OffsetY:         10,
			Locked:          &disable,
		},
	}
	if err := f.AddPictureFromBytes("Sheet1", "A1", pic); err != nil {
		log.Print("failed to add picture")
		return nil, err
	}

	const (
		wordsPerPage int = 180
		wordsPerLine int = 6
	)

	// text = strings.Repeat(text, 100)

	var (
		pageIndex int
		pageName  string
		progress  uint32
		words     []string = strings.Fields(text)
	)

	wordsCount := len(words)
	fivePercent := (wordsCount * 5) / 100

	for i := 0; i < len(words); i++ {
		if i%wordsPerPage == 0 {
			pageIndex++
			pageName = fmt.Sprint("Page ", pageIndex)
			_, err := f.NewSheet(pageName)
			if err != nil {
				log.Print("failed to add page")
				return nil, err
			}
			err = f.SetColWidth(pageName, "A", "F", 30)
			if err != nil {
				log.Print("failed to set column width")
				return nil, err
			}
		}

		if i%fivePercent == 0 {
			err := stream.Send(&reporter.GetResponse{
				Res: &reporter.GetResponse_Status{
					Status: progress,
				},
			})
			if err != nil {
				return nil, err
			}
			progress += 5
		}

		col := i%wordsPerLine + 1
		row := (i/wordsPerLine)%(wordsPerPage/wordsPerLine) + 1
		cell, err := excelize.CoordinatesToCellName(col, row)
		if err != nil {
			return nil, err
		}
		f.SetCellValue(pageName, cell, words[i])
	}

	name := filepath.Join("..", "storage", "reports", fmt.Sprintf("%s.xlsx", time))
	if err := f.SaveAs(name); err != nil {
		log.Print("failed to save report")
		return nil, err
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		log.Print("failed to write report to buffer")
		return nil, err
	}
	return buf, nil
}

func getLocalTime() (string, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print("failed to connect to time service")
		return "", err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			panic("failed to close file")
		}
	}()
	c := timer.NewTimerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := c.GetTime(ctx, &emptypb.Empty{})
	if err != nil {
		log.Print("failed to get current time")
		return "", err
	}

	t, err := getTimeInTimeZone(defaultLocation, res.Time.AsTime())
	if err != nil {
		log.Print("failed to get local time")
		return "", err
	}
	ts := fmt.Sprintf("%02d:%02d:%02d:%v", t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	return ts, nil
}

func getTimeInTimeZone(timeZone string, utcTime time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Print("failed to get time zone")
		return time.Time{}, err
	}
	return utcTime.In(loc), nil
}

func getRandomImage() ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", riUrl, nil)
	if err != nil {
		log.Print("failed to create http request")
		return nil, err
	}

	req.Header.Add("X-Api-Key", riApiKey)
	req.Header.Add("Accept", riAccept)
	resp, err := c.Do(req)
	if err != nil {
		log.Print("failed to do http request")
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("failed to read http response body")
		return nil, err
	}
	return data, nil
}

func getFishText() (string, error) {
	c := &http.Client{}
	url := fmt.Sprintf("%s?&type=%s&number=%d", ftUrl, ftType, ftNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("failed to create http request")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		log.Print("failed to do http request")
		return "", err
	}
	defer func() {
		resp.Body.Close()
		if err != nil {
			panic("failed to close response body")
		}
	}()

	decoder := json.NewDecoder(resp.Body)
	var ftResp ftResponse
	err = decoder.Decode(&ftResp)
	if err != nil {
		log.Print("failed to decode response")
		return "", err
	}
	if ftResp.ErrorCode != 0 {
		log.Print("failed to get fish text")
		return "", err
	}

	return ftResp.Text, nil
}

func main() {
	runtime.SetBlockProfileRate(1)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen %d: %s", grpcPort, err)
	}
	go func() {
		http.ListenAndServe(":6666", nil)
	}()
	gs := grpc.NewServer()
	reflection.Register(gs)
	reporter.RegisterReporterServer(gs, &server{})

	log.Printf("reporter is listening at %v", lis.Addr())
	log.Fatal(gs.Serve(lis))
}
