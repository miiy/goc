package upload

import (
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type File struct {
	Name string
	Dst  string
	Ext  string
	Hash string
	Size int64
}

func generateFileName() string {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return strconv.FormatInt(millis, 10)
}

func uploadDir() string {
	yearMonth := time.Now().Format("200601")
	day := time.Now().Format("02")
	return filepath.Join(yearMonth, day)
}

func UploadFile(basePath string, c *gin.Context) (*File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	fileSize := file.Size
	ext := filepath.Ext(file.Filename)
	// file type check

	fileName := generateFileName()
	path := uploadDir()
	dst := filepath.Join(basePath, path, fileName+ext)

	err = SaveUploadedFile(file, dst)
	if err != nil {
		return nil, err
	}
	return &File{
		Name: fileName,
		Dst:  path,
		Ext:  ext,
		Hash: "",
		Size: fileSize,
	}, nil
}

// SaveUploadedFile uploads the form file to specific dst.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
