package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/neepoo/pcbook/pb"
	"github.com/neepoo/pcbook/sample"
)

func main() {
	serverAddr := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s\n", *serverAddr)
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect grpc server error", err)
	}
	laptopClient := pb.NewLaptopServiceClient(conn)
	laptop := sample.NewLaptop()
	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 5 *time.Second)
	defer cancel()
	res, err :=laptopClient.CreateLaptop(ctx, &pb.CreateLaptopRequest{Laptop: laptop})
	if err != nil{
		log.Fatal("client create laptop error", err)
	}
	log.Printf("laptop create success id: %s", res.GetId())
}
