package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"

	"github.com/neepoo/pcbook/pb"
)

type LaptopServer struct {
	LaptopStore LaptopStore
	pb.UnimplementedLaptopServiceServer
}

func (server *LaptopServer) CreateLaptop(
	ctx context.Context,
	request *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := request.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.GetId())

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop id is not a valid uuid: %s", err)
		}
	}else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate uuid: %s", err)
		}
		laptop.Id = id.String()
	}
	// some heavy processing
	//time.Sleep(6 *time.Second)

	// check if cancel by client
	if ctx.Err() == context.Canceled{
		log.Print("request is canceled")
		return nil, status.Error(codes.DeadlineExceeded, "request is canceled")

	}

	// check if exceed deadline
	if ctx.Err() == context.DeadlineExceeded{
		log.Print("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	// save the laptop to in-mem
	err := server.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop: %s", err)
	}
	log.Printf("saved laptop with id: %s", laptop.Id)
	return &pb.CreateLaptopResponse{Id: laptop.Id}, nil
}



func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{LaptopStore: store}
}
