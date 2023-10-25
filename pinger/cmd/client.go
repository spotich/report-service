package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pinger "pinger/pkg/reporter/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50052"
	place   = "Asia/Novosibirsk"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to report server: %e", err)
	}
	defer conn.Close()

	c := pinger.NewReporterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := c.GetReport(ctx, &pinger.GetRequest{Location: place})
	if err != nil {
		log.Fatalf("failed to get report: %s", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		st := res.GetStatus()
		fmt.Printf("%d%%\n", st)
		if st == 100 {
			fmt.Println("report generated successfully.")
			break
		}
	}

	res, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Report is available at: ", res.GetUrl())
}
