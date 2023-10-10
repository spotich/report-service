package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pinger "pinger/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type emptyStruct struct{}

const (
	address         = "localhost:50052"
	chunkSize       = 4
	defaultLocation = "Asia/Novosibirsk"
)

var (
	void      emptyStruct
	timeZones map[string]emptyStruct
)

func registerLocations() error {
	file, err := os.Open("asia.txt")
	if err != nil {
		log.Print("failed to open locations file")
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		timeZones[fileScanner.Text()] = void
	}
	return nil
}

func fetchResponse(stream pinger.Reporter_GetReportClient) (*string, []byte, error) {
	resp, err := stream.Recv()
	if err != nil {
		log.Printf("failed to fetch reponse")
		return nil, nil, err
	}

	switch v := resp.Value.(type) {
	case *pinger.GetResponse_Time:
		return &v.Time, nil, nil
	case *pinger.GetResponse_Chunk:
		return nil, v.Chunk, nil
	default:
		err := errors.New("failed to fetch response")
		return nil, nil, err
	}
}

func main() {
	place := os.Args[1]
	_, registered := timeZones[place]

	if !registered {
		place = defaultLocation
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to report server: %e", err)
	}
	defer conn.Close()

	c := pinger.NewReporterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.GetReport(ctx, &pinger.GetRequest{Location: place})
	if err != nil {
		log.Fatalf("failed to get report: %s", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		log.Fatalf("failed to recieve a response: %s", err)
	}

	var localTime string
	switch v := resp.Value.(type) {
	case *pinger.GetResponse_Time:
		localTime = v.Time
		break
	default:
		log.Fatalf("failed to recieve time: %s", err)
	}

	path := filepath.Join("storage", localTime)
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}
	defer file.Close()

	log.Print("fetching the file...")
	num := 0
downloadLoop:
	for {
		res, err := stream.Recv()

		switch err {
		case nil:
			num++
			file.Write(res.Chunk)
			continue
		case io.EOF:
			log.Print("file fetched successfully")
			break downloadLoop
		default:
			log.Print(err)
			break downloadLoop
		}
	}
	log.Printf("amount of chunks fetched: %d", num)
	log.Printf("amount of bytes fetched: %d", num*chunkSize)
}
