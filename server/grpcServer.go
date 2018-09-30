package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "ocr-engine/grpc"

	"ocr-engine/ocr"
	"ocr-engine/parser"
	"ocr-engine/rasterize"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "GRPC server port")
)

type ocrServer struct {
}

func (server *ocrServer) OCR(context context.Context, request *pb.OCRRequest) (*pb.OCRResponse, error) {

	log.Printf("Received file %s\n", request.Filename)

	doc := rasterize.NewSourceDocument(request.Filename)
	rasterOpts := &rasterize.DocumentRasterizeOptions{rasterize.EngineGhostScriptLib}

	image, err := doc.ToImage(rasterOpts)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Rasterized %vkb\n", (len(image) / 1204))

	hocr, err := ocr.GetHOCR(image)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	blocks := parser.Parse(hocr)

	response := &pb.OCRResponse{Blocks: blocks}

	return response, nil
}

func newServer() *ocrServer {
	server := &ocrServer{}
	return server
}

func main() {

	address := fmt.Sprintf("localhost:%v", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen %s. %v", address, err)
	}

	log.Printf("Listening on %s\n", address)

	grpcServer := grpc.NewServer()

	server := newServer()
	pb.RegisterOCRServiceServer(grpcServer, server)

	grpcServer.Serve(listener)

}
