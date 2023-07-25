package main

import (
	"context"
	"flag"
	pb "github.com/miiy/goc/service/auth/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
	//login(c)
	//mpLogin(c)
	verifyToken(c)
}

func register(c pb.AuthServiceClient) {
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

func mpLogin(c pb.AuthServiceClient) {
	code := ""
	req := pb.MpLoginRequest{
		Code: code,
	}
	resp, err := callMpLogin(c, &req)
	if err != nil {
		log.Fatalf("client.callMpLogin(_) = _, %v", err)
	}
	log.Println("MpLogin:", resp)
}

func verifyToken(c pb.AuthServiceClient) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6IiIsImV4cCI6MTY5MDI4MTI3NSwiaWF0IjoxNjkwMjc0MDc1fQ.BH7oNkJmZF20eUQs36d_to7txg8zEaYbWxXfZLNNwKk"
	resp, err := callVerifyToken(c, token)
	if err != nil {
		log.Fatalf("client.callMpLogin(_) = _, %v", err)
	}
	log.Println("VerifyToken:", resp)
}

func callMpLogin(client pb.AuthServiceClient, req *pb.MpLoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	return client.MpLogin(ctx, req)
}

func callRegister(client pb.AuthServiceClient, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.Register(ctx, req)
}

func callVerifyToken(client pb.AuthServiceClient, accessToken string) (*pb.VerifyTokenResponse, error) {
	md := metadata.Pairs("Authorization", "Bearer "+accessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
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
