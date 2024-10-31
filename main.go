package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"sync"
	cli "testTask/cmd/client"
	database "testTask/cmd/database"
	s "testTask/cmd/server"
	pb "testTask/pkg/api"
)

var (
	asyncFlag  = flag.Bool("async", false, "Download thumbnails asynchronously")
	serverAddr = flag.String("server_addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	db := database.InitDB()
	defer database.CloseDB()

	lis, err := net.Listen("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	defer grpcServer.Stop()

	service := s.NewThumbnailServer(db)
	pb.RegisterThumbnailServiceServer(grpcServer, service)

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)

	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s <url1> <url2> ...", os.Args[0])
	}
	urls := flag.Args()

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		if *asyncFlag {
			go func(url string) {
				defer wg.Done()
				cli.DownloadAndSaveThumbnail(service, url)
			}(url)
		} else {
			cli.DownloadAndSaveThumbnail(service, url)
			wg.Done()
		}
	}
	wg.Wait()
}
