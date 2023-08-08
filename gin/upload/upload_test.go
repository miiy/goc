package upload

import "testing"

func TestGenerateFileName(t *testing.T) {
	t.Log(generateFileName())
}

func TestUploadDir(t *testing.T) {
	t.Log(uploadDir())
}
