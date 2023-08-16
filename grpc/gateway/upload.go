package gateway

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func handleBinaryFileUpload(w http.ResponseWriter, r *http.Request, params map[string]string) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
		return
	}

	f, header, err := r.FormFile("attachment")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get file 'attachment': %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer f.Close()

	//
	// Now do something with the io.Reader in `f`, i.e. read it into a buffer or stream it to a gRPC client side stream.
	// Also `header` will contain the filename, size etc of the original file.
	//

	fileSize := header.Size
	ext := filepath.Ext(header.Filename)
	fmt.Println(fileSize, ext, header)
}
