package commands

import "google.golang.org/grpc"

func dialGrpc() *grpc.ClientConn, error {
	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	} else {
		return &conn, nil
	}
}
