package main

import (
	"context"
	"flag"
	pb "github.com/miiy/goc/service/auth/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	var addr = flag.String("addr", "localhost:50051", "the address to connect to")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dit not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	//register(c)
	login(c)
}

func register(c pb.AuthServiceClient) {
	log.Println("--- calling up-api.Auth/SignUp ---")
	regReq := pb.RegisterRequest{
		Email:                "test@test.com",
		Username:             "username",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	rResp, err := callRegister(c, &regReq)
	if err != nil {
		log.Fatalf("client.callRegister(_) = _, %v", err)
	}
	log.Println("SignUp:", rResp)
}

func login(c pb.AuthServiceClient) {
	log.Println("--- calling up-api.Auth/login ---")
	regReq := pb.LoginRequest{
		Username: "username",
		Password: "123456",
	}
	rResp, err := callLogin(c, &regReq)
	if err != nil {
		log.Fatalf("client.callRegister(_) = _, %v", err)
	}
	log.Println("Login:", rResp)
}

func callRegister(client pb.AuthServiceClient, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.Register(ctx, req)
}

func callVerifyToken(client pb.AuthServiceClient, accessToken string) (*pb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := pb.VerifyTokenRequest{
		AccessToken: accessToken,
	}
	return client.VerifyToken(ctx, &req)
}

func callLogin(client pb.AuthServiceClient, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.Login(ctx, req)
}
