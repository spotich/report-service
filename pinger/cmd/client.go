package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	pinger "pinger/pkg/reporter/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	pingerAddress   = ":50050"
	reporterAddress = "reporter-service:50052"
	place           = "Asia/Novosibirsk"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("pinger is listening at", pingerAddress)
	log.Fatal(http.ListenAndServe(pingerAddress, nil))
}

func handler(w http.ResponseWriter, _ *http.Request) {
	log.Println("1 new request")
	conn, err := grpc.Dial(reporterAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		w.Write([]byte("failed to connect to report server"))
		log.Println("failed to connect to report server")
		return
	}
	defer conn.Close()

	c := pinger.NewReporterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := c.GetReport(ctx, &pinger.GetRequest{Location: place})
	if err != nil {
		w.Write([]byte("failed to get report"))
		log.Println("failed to get report")
		return
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			w.Write([]byte("failed to get recieve data"))
			log.Println("failed to get recieve data")
			return
		}
		st := res.GetStatus()
		fmt.Printf("%d%%\n", st)
		if st == 100 {
			w.Write([]byte("report generated successfully"))
			log.Println("report generated successfully")
			break
		}
	}

	res, err := stream.Recv()
	if err != nil {
		w.Write([]byte("failed to get download link"))
		log.Println("failed to get download link")
		return
	}
	resp := fmt.Sprintf("Report is available at: %s", res.GetUrl())
	w.Write([]byte(resp))
	log.Println("report is available at ", resp)
}
