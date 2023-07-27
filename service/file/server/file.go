package server

import (
	"context"
	"docxlib.com/pkg/database"
	"docxlib.com/pkg/database/gorm/paginate"
	"docxlib.com/pkg/database/gorm/scope"
	"docxlib.com/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "zsaix.com/api/file/v1"
	"zsaix.com/apiserver/internal/file/repository"
)

type FileServer struct {
	Repository repository.FileRepository
	Logger     *log.Logger
	pb.UnimplementedFileServiceServer
}

func NewFileServer(db *database.GormDB, logger *log.Logger) pb.FileServiceServer {
	r := repository.NewFileRepository(db, logger)
	return &FileServer{
		Repository: r,
		Logger:     logger,
	}
}

func (s *FileServer) CreateFile(ctx context.Context, request *pb.CreateFileRequest) (*pb.File, error) {
	reqFile := request.GetFile()

	file := &repository.File{
		UserId: reqFile.UserId,
		Type:   int(reqFile.Type),
		Name:   reqFile.Name,
		Ext:    reqFile.Ext,
		Path:   reqFile.Path,
		Hash:   reqFile.Hash,
		Status: int(reqFile.Status),
	}

	_, err := s.Repository.CreateFile(ctx, file)
	if err != nil {
		return nil, err
	}
	filePb := fileToProto(file)
	return filePb, nil
}

func (s *FileServer) GetFile(ctx context.Context, request *pb.GetFileRequest) (*pb.File, error) {
	file, err := s.Repository.GetFileById(ctx, request.Id)
	if file != nil {
		if err == repository.ErrRecordNotFound {
			st := status.New(codes.NotFound, err.Error())
			return nil, st.Err()
		}
		return nil, err
	}

	return fileToProto(file), nil
}

func (s *FileServer) UpdateFile(ctx context.Context, request *pb.UpdateFileRequest) (*pb.RowsAffected, error) {
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
		UserId: file.UserId,
		Type:   int(file.Type),
		Name:   file.Name,
		Ext:    file.Ext,
		Path:   file.Path,
		Hash:   file.Hash,
		Status: int(file.Status),
	}
	ra, err := s.Repository.UpdateFile(ctx, request.Id, b, selects)
	return &pb.RowsAffected{RowsAffected: ra}, err
}

func (s *FileServer) DeleteFile(ctx context.Context, request *pb.DeleteFileRequest) (*emptypb.Empty, error) {
	err := s.Repository.DeleteFileById(ctx, request.Id)
	return nil, err
}

func (s *FileServer) ListFiles(ctx context.Context, request *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// count
	total, err := s.Repository.FindCount(ctx,
		scope.ScopeOfStatus([]int64{repository.StatusActive}),
		repository.ScopeOfFileCategory(request.CategoryId),
	)
	if err != nil {
		return nil, err
	}
	scopePaginate, totalPages := paginate.Paginate(int(request.Page), int(request.PageSize), int(total))

	// list
	files, err := s.Repository.Find(ctx,
		scope.ScopeOfStatus([]int64{repository.StatusActive}),
		repository.ScopeOfFileCategory(request.CategoryId),
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
		UserId:     v.UserId,
		Type:       pb.File_Type(v.Type),
		Name:       v.Name,
		Ext:        v.Ext,
		Path:       v.Path,
		Hash:       v.Hash,
		Status:     pb.File_Status(v.Status),
		CreateTime: timestamppb.New(v.CreateTime),
		UpdateTime: timestamppb.New(v.UpdateTime),
		DeleteTime: timestamppb.New(v.DeleteTime.Time),
	}
}
