package server

import (
	"context"
	"errors"
	pb "github.com/miiy/goc/component/file/api/v1"
	"github.com/miiy/goc/component/file/entity"
	"github.com/miiy/goc/component/file/repository"
	"github.com/miiy/goc/db/gorm/paginate"
	"github.com/miiy/goc/db/gorm/scope"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type fileServer struct {
	pb.UnimplementedFileServiceServer

	logger *zap.Logger
	repo   repository.FileRepository
}

func NewFileServer(logger *zap.Logger, repo repository.FileRepository) pb.FileServiceServer {
	return &fileServer{
		logger: logger,
		repo:   repo,
	}
}

func (s *fileServer) CreateFile(ctx context.Context, request *pb.CreateFileRequest) (*pb.File, error) {
	reqFile := request.GetFile()

	file := &entity.File{
		SysId:    reqFile.SysId,
		CatId:    reqFile.CatId,
		ItemId:   reqFile.ItemId,
		UserId:   reqFile.UserId,
		FileType: int(reqFile.FileType),
		Name:     reqFile.Name,
		Ext:      reqFile.Ext,
		Path:     reqFile.Path,
		Hash:     reqFile.Hash,
		Status:   int(reqFile.Status),
	}

	_, err := s.repo.CreateFile(ctx, file)
	if err != nil {
		s.logger.Error("repo.CreateFile", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	filePb := fileToProto(file)
	return filePb, nil
}

func (s *fileServer) GetFile(ctx context.Context, request *pb.GetFileRequest) (*pb.File, error) {
	file, err := s.repo.GetFileById(ctx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "file not found")
		}
		return nil, err
	}

	return fileToProto(file), nil
}

func (s *fileServer) UpdateFile(ctx context.Context, request *pb.UpdateFileRequest) (*pb.RowsAffected, error) {
	file := request.GetFile()
	if !request.UpdateMask.IsValid(file) {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	selects := []string{"*"}
	if request.UpdateMask.GetPaths() != nil {
		selects = request.UpdateMask.GetPaths()
	}

	b := &entity.File{
		SysId:    file.SysId,
		CatId:    file.CatId,
		ItemId:   file.ItemId,
		UserId:   file.UserId,
		FileType: int(file.FileType),
		Name:     file.Name,
		Ext:      file.Ext,
		Path:     file.Path,
		Hash:     file.Hash,
		Status:   int(file.Status),
	}
	ra, err := s.repo.UpdateFile(ctx, request.Id, b, selects)
	if err != nil {
		s.logger.Error("repo.UpdateFile", zap.Error(err))
		return nil, status.Error(codes.Internal, "update fail")
	}
	return &pb.RowsAffected{RowsAffected: ra}, err
}

func (s *fileServer) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*emptypb.Empty, error) {
	err := s.repo.DeleteFileById(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "delete fail")
	}
	return nil, nil
}

func (s *fileServer) ListFiles(ctx context.Context, request *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// count
	total, err := s.repo.FindCount(ctx,
		scope.ScopeOfStatus([]int64{entity.FileStatusActive}),
		repository.ScopeOfFileSys(request.SysId),
		repository.ScopeOfFileCat(request.CatId),
	)
	if err != nil {
		return nil, err
	}
	scopePaginate, totalPages := paginate.Paginate(int(request.Page), int(request.PageSize), 50, int(total))

	// list
	files, err := s.repo.Find(ctx,
		scope.ScopeOfStatus([]int64{entity.FileStatusActive}),
		repository.ScopeOfFileSys(request.SysId),
		repository.ScopeOfFileCat(request.CatId),
		scopePaginate,
	)
	if err != nil {
		return nil, err
	}

	var items []*pb.File
	for _, v := range files {
		items = append(items, fileToProto(v))
	}

	return &pb.ListFilesResponse{
		Total:       total,
		TotalPages:  int64(totalPages),
		PageSize:    request.PageSize,
		CurrentPage: request.Page,
		Files:       items,
	}, nil
}

func fileToProto(v *entity.File) *pb.File {
	return &pb.File{
		Id:         v.ID,
		SysId:      v.SysId,
		CatId:      v.CatId,
		ItemId:     v.ItemId,
		UserId:     v.UserId,
		FileType:   pb.File_FileType(v.FileType),
		Name:       v.Name,
		Ext:        v.Ext,
		Path:       v.Path,
		Hash:       v.Hash,
		Status:     pb.File_FileStatus(v.Status),
		CreateTime: timestamppb.New(v.CreateTime),
		UpdateTime: timestamppb.New(v.UpdateTime),
		DeleteTime: timestamppb.New(v.DeleteTime.Time),
	}
}
