package server

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/miiy/goc/auth/jwt"
	pb "github.com/miiy/goc/service/auth/api/v1"
	"github.com/miiy/goc/service/auth/entity"
	"github.com/miiy/goc/service/auth/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type AuthServiceServer struct {
	repo      repository.AuthRepository
	tokenRepo repository.AuthTokenRepository
	jwtAuth   *jwt.JWTAuth
	pb.UnimplementedAuthServiceServer
}

const (
	AuthTokenKey = "user_token:%s" // user_token:md5({user_id})
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrPasswordsDiffer = errors.New("passwords differ")
	ErrUnauthenticated = errors.New("unauthenticated")

	ErrUsernameOrEmailExist = errors.New("username or email already exist")

	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

func NewAuthServiceServer(repo repository.AuthRepository, tokenRepo repository.AuthTokenRepository, jwtAuth *jwt.JWTAuth) pb.AuthServiceServer {
	return &AuthServiceServer{
		repo:      repo,
		tokenRepo: tokenRepo,
		jwtAuth:   jwtAuth,
	}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := registerValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	exist, err := s.repo.UserExist(ctx, entity.UserColumnUsername, req.Username)
	if err != nil {
		grpclog.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exist {
		return nil, status.Error(codes.AlreadyExists, ErrUsernameOrEmailExist.Error())
	}

	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		grpclog.Errorln(err)
		return nil, err
	}

	user := entity.User{
		Username:          req.Username,
		Password:          string(hashPasswd),
		Email:             req.Email,
		EmailVerifiedTime: nil,
		Phone:             "",
		Status:            0,
	}

	// register
	err = s.repo.Create(ctx, &user)
	if err != nil {
		grpclog.Errorln(err)
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.AuthenticatedUser{
			Username: user.Username,
		},
	}, nil
}

func (s *AuthServiceServer) UsernameCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnUsername, req.Value)
}

func (s *AuthServiceServer) EmailCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnEmail, req.Value)
}

func (s *AuthServiceServer) PhoneCheck(ctx context.Context, req *pb.FieldCheckRequest) (*pb.FieldCheckResponse, error) {
	return s.userExist(ctx, entity.UserColumnPhone, req.Value)
}

func (s *AuthServiceServer) userExist(ctx context.Context, field, value string) (*pb.FieldCheckResponse, error) {
	exist, err := s.repo.UserExist(ctx, field, value)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.FieldCheckResponse{
		Exist: exist,
	}, nil
}

// Login
func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := loginValidate(req); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	user, err := s.repo.FirstByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		grpclog.Error(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrWrongPassword
	}

	claims := s.jwtAuth.CreateClaims(user.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		return nil, err
	}

	// store token
	err = s.tokenRepo.Set(ctx, formatTokenKey(token), token, time.Second*time.Duration(s.jwtAuth.Options.ExpiresIn))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(claims.ExpiresAt.Time),
		User: &pb.AuthenticatedUser{
			Username: user.Username,
		},
	}, nil
}

// VerifyToken
func (s *AuthServiceServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {

	if err := verifyTokenValidate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	claims, err := s.jwtAuth.ParseToken(req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.VerifyTokenResponse{
		User: &pb.AuthenticatedUser{
			Username: claims.Username,
		},
	}, nil
}

// RefreshToken
// 1. validate old token
// 2. delete old token
// 3. create new token
func (s *AuthServiceServer) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// validate old token and create new token
	oldClaims, err := s.jwtAuth.ParseToken(request.AccessToken)
	if err != nil {
		return nil, err
	}
	// delete old token
	err = s.tokenRepo.Del(ctx, formatTokenKey(request.AccessToken))
	if err != nil {
		return nil, err
	}
	// create new token
	claims := s.jwtAuth.CreateClaims(oldClaims.Username)
	token, err := s.jwtAuth.CreateTokenByClaims(claims)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenResponse{
		TokenType:   "Bearer",
		AccessToken: token,
		ExpiresAt:   timestamppb.New(claims.ExpiresAt.Time),
		User:        &pb.AuthenticatedUser{Username: claims.Username},
	}, nil
}

// Logout
// 1. delete token
func (s *AuthServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.tokenRepo.Del(ctx, formatTokenKey(req.AccessToken))
	if err != nil {
		return nil, err
	}
	return &pb.LogoutResponse{}, nil
}

func formatTokenKey(token string) string {
	return fmt.Sprintf(AuthTokenKey, fmt.Sprintf("%x", md5.Sum([]byte(token))))
}
