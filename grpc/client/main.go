package main

//
//import (
//	"context"
//	authpb "docxlib.com/api/auth/v1"
//	"flag"
//	"google.golang.org/grpc"
//	"log"
//	"time"
//)
//
//type user struct {
//	email    string
//	username string
//	password string
//}
//
//func main() {
//	var addr = flag.String("addr", "localhost:50051", "the address to connect to")
//	flag.Parse()
//
//	// Set up a connection to the server.
//	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		log.Fatalf("dit not connect: %v", err)
//	}
//	defer conn.Close()
//
//	tc := authpb.NewAuthServiceClient(conn)
//	u1 := &user{
//		email:    "a@a.com",
//		username: "a3",
//		password: "a3",
//	}
//
//	log.Println("--- calling up-api.Auth/SignUp ---")
//	rResp, err := callSignUp(tc, u1.email, u1.username, u1.password, u1.password)
//	if err != nil {
//		log.Fatalf("client.SignUp(_) = _, %v", err)
//	}
//	log.Println("SignUp:", rResp)
//
//	log.Println("--- calling up-api.Auth/SignIn ---")
//	lResp, err := callSignIn(tc, u1.username, u1.password)
//	if err != nil {
//		log.Fatalf("client.SignIn(_) = _, %v", err)
//	}
//	log.Println("SignIn:", lResp)
//
//	log.Println("--- calling up-api.Auth/VerifyToken ---")
//	vResp, err := callVerifyToken(tc, lResp.AccessToken)
//	if err != nil {
//		log.Fatalf("client.VerifyToken(_) = _, %v", err)
//	}
//	log.Println("VerifyToken:", vResp)
//}
//
//func callSignUp(client authpb.AuthServiceClient, email, username, password, passwordConfirmation string) (*authpb.SignUpResponse, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	req := authpb.SignUpRequest{
//		Email:                email,
//		Username:             username,
//		Password:             password,
//		PasswordConfirmation: passwordConfirmation,
//	}
//	return client.SignUp(ctx, &req)
//}
//
//func callVerifyToken(client authpb.AuthServiceClient, accessToken string) (*authpb.VerifyTokenResponse, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	req := authpb.VerifyTokenRequest{
//		AccessToken: accessToken,
//	}
//	return client.VerifyToken(ctx, &req)
//}
//
//func callSignIn(client authpb.AuthServiceClient, username, password string) (*authpb.SignInResponse, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	req := &authpb.SignInRequest{
//		Username: username,
//		Password: password,
//	}
//	return client.SignIn(ctx, req)
//}
