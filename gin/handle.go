package gin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// UploadTmpFile upload file to tmp dir
func UploadTmpFile(c *gin.Context) (string, func(), error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", nil, err
	}
	fileSuffix := path.Ext(file.Filename)
	if fileSuffix != ".xlsx" {
		return "", nil, errors.New("file type error")
	}

	// 上传文件至指定的完整文件路径
	dst := fmt.Sprintf("/tmp/%s%d%s", strings.TrimSuffix(file.Filename, fileSuffix), time.Now().Unix(), path.Ext(file.Filename))
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		return "", nil, err
	}
	return dst, func() {
		os.Remove(dst)
	}, nil
}

// ExportFile 导出文件
func ExportFile(c *gin.Context, fileName string, data []byte) {
	Header := c.Writer.Header()
	Header.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(fileName)))
	Header.Set("Content-Description", "File Transfer")
	Header.Set("Content-Transfer-Encoding", "binary")
	Header.Set("Expires", "0")
	Header.Set("Cache-Control", "must-revalidate")
	Header.Set("Pragma", "public")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
