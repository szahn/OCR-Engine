package main

import (
	"fmt"
	pb "ocr-engine/grpc"
	"ocr-engine/parser"

	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error dialing %s: %v", address, err)
	}

	defer conn.Close()
	client := pb.NewOCRServiceClient(conn)

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filename := "samples/pdf-test.pdf"
	response, err := client.OCR(context, &pb.OCRRequest{Filename: filename})
	if err != nil {
		log.Fatalf("Error ocring %s: %v", filename, err)
	}

	for _, b := range response.Blocks {
		fmt.Printf("%s\n", parser.DescribeBlock(b))
	}

}
