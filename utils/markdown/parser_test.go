package markdown

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	testMarkdown1, err := os.ReadFile("./test-empty.md")
	if err != nil {
		t.Log(err)
	}
	var parser Parser
	blocks, err := parser.ParseBlock(testMarkdown1)
	if err != nil {
		t.Error(err)
	}
	t.Log(blocks)
}

func TestParseBlock(t *testing.T) {
	testMarkdown1, err := os.ReadFile("./test1.md")
	if err != nil {
		t.Log(err)
	}

	var parser Parser
	blocks, err := parser.ParseBlock(testMarkdown1)
	if err != nil {
		t.Error(err)
	}
	for _, block := range blocks {
		fmt.Println("-----------------------------------------------------", block.Kind)
		fmt.Print(string(block.Content))
	}
	t.Log(blocks)
}

func TestOutput(t *testing.T) {
	testMarkdown1, err := os.ReadFile("./test0.md")
	if err != nil {
		t.Log(err)
	}

	var parser Parser
	blocks, err := parser.ParseBlock(testMarkdown1)
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	for _, block := range blocks {
		_, err := buf.Write(block.Content)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("-----------------------------------------------------", block.Kind)
		fmt.Print(string(block.Content))
	}

	// 写入目标文件
	err = os.WriteFile("./test0_dst.md", buf.Bytes(), 0644)
	if err != nil {
		t.Error(err)
	}
	t.Log(blocks)
}

func TestIsMetaData(t *testing.T) {
	c1 := "---"
	t.Log(isMetadata([]byte(c1)))
	c2 := "----"
	t.Log(isMetadata([]byte(c2)))

	c3 := "--- "
	t.Log(isMetadata([]byte(c3)))
	c4 := " ---"
	t.Log(isMetadata([]byte(c4)))
	c5 := "--"
	t.Log(isMetadata([]byte(c5)))
	c6 := "--"
	t.Log(isMetadata([]byte(c6)))
}

// 0-3个空白符+3个```或~~~
func TestFenceLine(t *testing.T) {
	c1 := "```"
	t.Log(isFenceLine([]byte(c1)))
	c2 := "~~~"
	t.Log(isFenceLine([]byte(c2)))
	c3 := "   ```"
	t.Log(isFenceLine([]byte(c3)))

	c4 := ""
	t.Log(isFenceLine([]byte(c4)))
	c5 := " "
	t.Log(isFenceLine([]byte(c5)))
	c6 := "    ```"
	t.Log(isFenceLine([]byte(c6)))
	c7 := "	```"
	t.Log(isFenceLine([]byte(c7)))
}

// 4个空白符或1个tab开始
func TestCodePrefix(t *testing.T) {
	c1 := "    ls"
	t.Log(isCodePrefix([]byte(c1)))
	c2 := "	ls"
	t.Log(isCodePrefix([]byte(c2)))

	c3 := "   ls"
	t.Log(isCodePrefix([]byte(c3)))
	c4 := "ls"
	t.Log(isCodePrefix([]byte(c4)))
	c5 := ""
	t.Log(isCodePrefix([]byte(c5)))
}
