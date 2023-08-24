package server

import (
	"context"
	postv1 "github.com/miiy/goc/examples/apiserver/gen/goc/post/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type postServer struct {
	postv1.UnimplementedPostServiceServer
}

func NewPostServiceServer() postv1.PostServiceServer {
	return &postServer{}
}

func (s *postServer) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.GetPostResponse, error) {
	return &pb.Message{
		Id:       1,
		UserName: "test",
		Message:  "hello",
	}, nil
}

func (s *echoServer) EchoQueryId(ctx context.Context, req *pb.EchoQueryIdRequest) (*pb.Message, error) {
	return &pb.Message{
		Id:       req.Id,
		UserName: "test",
		Message:  "hello",
	}, nil
}

func (s *echoServer) EchoPost(ctx context.Context, req *pb.EchoPostRequest) (*pb.Message, error) {
	return &pb.Message{
		Id:       req.Id,
		UserName: req.UserName,
		Message:  "",
	}, nil
}

func (s *echoServer) EchoError(ctx context.Context, req *emptypb.Empty) (*pb.Message, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid parameters")
}
