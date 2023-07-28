package server

import (
	"context"
	"github.com/miiy/goc/db/gorm/paginate"
	"github.com/miiy/goc/db/gorm/scope"
	pb "github.com/miiy/goc/service/file/api/v1"
	"github.com/miiy/goc/service/file/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type fileServer struct {
	pb.UnimplementedFileServiceServer

	repo repository.FileRepository
}

func NewFileServer(repo repository.FileRepository) pb.FileServiceServer {
	return &fileServer{
		repo: repo,
	}
}

func (s *fileServer) CreateFile(ctx context.Context, request *pb.CreateFileRequest) (*pb.File, error) {
	reqFile := request.GetFile()

	file := &repository.File{
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
		return nil, err
	}
	filePb := fileToProto(file)
	return filePb, nil
}

func (s *fileServer) GetFile(ctx context.Context, request *pb.GetFileRequest) (*pb.File, error) {
	file, err := s.repo.GetFileById(ctx, request.Id)
	if file != nil {
		return nil, err
	}

	return fileToProto(file), nil
}

func (s *fileServer) UpdateFile(ctx context.Context, request *pb.UpdateFileRequest) (*pb.RowsAffected, error) {
	file := request.GetFile()
	if !request.UpdateMask.IsValid(file) {
		st := status.New(codes.InvalidArgument, "invalid field mask")
		return nil, st.Err()
	}

	selects := []string{"*"}
	if request.UpdateMask.GetPaths() != nil {
		selects = request.UpdateMask.GetPaths()
	}

	b := &repository.File{
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
	return &pb.RowsAffected{RowsAffected: ra}, err
}

func (s *fileServer) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*emptypb.Empty, error) {
	err := s.repo.DeleteFileById(ctx, request.Id)
	return nil, err
}

func (s *fileServer) ListFiles(ctx context.Context, request *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// count
	total, err := s.repo.FindCount(ctx,
		scope.ScopeOfStatus([]int64{repository.StatusActive}),
		repository.ScopeOfFileSys(request.SysId),
		repository.ScopeOfFileCat(request.CatId),
	)
	if err != nil {
		return nil, err
	}
	scopePaginate, totalPages := paginate.Paginate(int(request.Page), int(request.PageSize), 50, int(total))

	// list
	files, err := s.repo.Find(ctx,
		scope.ScopeOfStatus([]int64{repository.StatusActive}),
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

func fileToProto(v *repository.File) *pb.File {
	return &pb.File{
		Id:         v.Id,
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
