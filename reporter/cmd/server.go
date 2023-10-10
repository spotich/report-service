package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	reporter "reporter/pkg/reporter/service"
	timer "reporter/pkg/timer/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	apiKey = "IJDHKmA5qlYGvBjD/CyPRw==3CIIQF9Nmym5mVwB"
	url    = "https://api.api-ninjas.com/v1/randomimage"
	accept = "image/jpg"

	grpcPort  = 50052
	message   = "hmm, what time is it?"
	address   = "localhost:50051"
	chunkSize = 1024

	defaultLocation = "Asia/Novosibirsk"
)

type server struct {
	reporter.UnimplementedReporterServer
}

func (s *server) GetReport(req *reporter.GetRequest, stream reporter.Reporter_GetReportServer) error {
	time, err := getLocalTime()
	if err != nil {
		log.Print("failed to get local time")
		return err
	}

	err = sendLocalTime(time, stream)
	if err != nil {
		log.Print("failed to send local time")
		return err
	}

	data, err := getRandomImage()
	if err != nil {
		log.Print("failed to get random image")
		return err
	}

	path := filepath.Join("storage", fmt.Sprint(time, ".jpg"))
	err = saveDataToFile(data, path)
	if err != nil {
		log.Print("failed to save image")
		return err
	}

	err = sendDataToStream(data, stream)
	if err != nil {
		log.Print("failed to send image")
		return err
	}

	return nil
}

func getLocalTime() (string, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print("failed to connect to time service")
		return "", err
	}
	defer conn.Close()

	c := timer.NewTimerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	res, err := c.GetTime(ctx, &timer.GetRequest{Message: message})
	if err != nil {
		log.Print("failed to get current time")
		return "", err
	}

	t, err := getTimeInTimeZone(defaultLocation, res.Time.AsTime())
	if err != nil {
		log.Print("failed to get local time")
		return "", err
	}
	time := fmt.Sprintf("%02d:%02d:%02d:%v", t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	return time, nil
}

func getTimeInTimeZone(timeZone string, utcTime time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Print("failed to get time zone")
		return time.Time{}, err
	}
	return utcTime.In(loc), nil
}

func sendLocalTime(time string, stream reporter.Reporter_GetReportServer) error {
	err := stream.Send(&reporter.GetResponse{
		Value: &reporter.GetResponse_Time{
			Time: time,
		},
	})
	if err != nil {
		log.Print("failed to send time")
		return err
	}
	return nil
}

func getRandomImage() ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("failed to create http request")
		return nil, err
	}

	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("Accept", accept)
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

func saveDataToFile(data []byte, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Print("failed to create file")
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		log.Print("failed to save file")
		return err
	}
	log.Printf("saved file is %d bytes", n)
	return nil
}

func sendDataToStream(data []byte, stream reporter.Reporter_GetReportServer) error {
	buf := make([]byte, chunkSize)
	reader := bytes.NewReader(data)
	chunksSent := 0
	bytesSent := 0

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print("failed to read chunk")
			return err
		}
		err = stream.Send(&reporter.GetResponse{
			Value: &reporter.GetResponse_Chunk{
				Chunk: buf[:n],
			},
		})
		if err != nil {
			log.Print("failed to send bytes")
			return err
		}
		chunksSent++
		bytesSent += n
	}

	log.Printf("Total amount of chunks sent: %d", chunksSent)
	log.Printf("Total amount of bytes sent: %d", bytesSent)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen %d: %s", grpcPort, err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	reporter.RegisterReporterServer(s, &server{})

	log.Printf("reporter is listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
