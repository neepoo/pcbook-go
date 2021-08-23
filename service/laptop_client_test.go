package service_test

import (
	"context"
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

	laptopServer, serverAddr := startTestLaptopServer(t)
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

func newLaptopClient(t *testing.T, addr string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	// 开始监听grpc请求
	go grpcServer.Serve(listener) // block call


	return laptopServer, listener.Addr().String()
}
