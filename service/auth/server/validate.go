package server

import (
	pb "github.com/miiy/goc/service/auth/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func registerValidate(req *pb.RegisterRequest) error {
	// trim space
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.PasswordConfirmation = strings.TrimSpace(req.PasswordConfirmation)

	// validate
	if req.Username == "" || req.Email == "" || req.Password == "" || req.PasswordConfirmation == "" {
		return ErrInvalidArgument
	}
	if req.Password != req.PasswordConfirmation {
		return ErrPasswordsDiffer
	}
	return nil
}

func loginValidate(req *pb.LoginRequest) error {
	if req.Username == "" || req.Password == "" {
		return ErrInvalidArgument
	}
	return nil
}

func mpLoginValidate(req *pb.MpLoginRequest) error {
	if req.Code == "" {
		return status.New(codes.InvalidArgument, "code can not empty").Err()
	}
	return nil
}

func verifyTokenValidate(request *pb.VerifyTokenRequest) error {
	if request.AccessToken == "" {
		return ErrInvalidArgument
	}
	return nil
}
