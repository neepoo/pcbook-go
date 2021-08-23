package service_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/neepoo/pcbook/pb"
	"github.com/neepoo/pcbook/sample"
	"github.com/neepoo/pcbook/service"
)

func requireSameLaptop(t *testing.T, l1, l2 *pb.Laptop)  {
	d1, err := protojson.Marshal(l1)
	require.NoError(t, err)
	d2, err := protojson.Marshal(l2)
	require.NoError(t, err)
	require.Equal(t, d1, d2)

}

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddr := startTestLaptopServer(t, service.NewInMemoryLaptopStore())
	laptopClient := newLaptopClient(t, serverAddr)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	// 需要确保新创建的笔记本确实存储了
	other, err := laptopServer.LaptopStore.Find(laptop.Id)
	require.NoError(t, err)
	require.NotNil(t, other)
	requireSameLaptop(t, laptop, other)


}


func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()
	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		},
	}
	store := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)
	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{
				Value: 3,
				Unit:  pb.Memory_GIGABYTE,
			}
		case 5, 6:
			t.Log("------------------")
			laptop.PriceUsd = 1900
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{
				Value: 16,
				Unit:  pb.Memory_GIGABYTE,
			}
			expectedIDs[laptop.Id] = true

		}
		err := store.Save(laptop)
		require.NoError(t, err)
	}

	_, serverAddr := startTestLaptopServer(t, store)
	laptopClient := newLaptopClient(t, serverAddr)
	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)
	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF{
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())
		found+=1
	}
	require.Equal(t, found, len(expectedIDs))
}

func newLaptopClient(t *testing.T, addr string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func startTestLaptopServer(t *testing.T, store service.LaptopStore) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(store)

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	// 开始监听grpc请求
	go grpcServer.Serve(listener) // block call


	return laptopServer, listener.Addr().String()
}
