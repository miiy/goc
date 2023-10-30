package markdown

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestSplitByHeading(t *testing.T) {
	testMarkdown1, err := os.ReadFile("./test1.md")
	if err != nil {
		t.Log(err)
	}

	contents, err := SplitByHeading(bytes.NewReader(testMarkdown1))
	if err != nil {
		t.Log(err)
	}
	for _, content := range contents {
		fmt.Print(string(content))
	}
	fmt.Println("---")

	testMarkdown2, err := os.ReadFile("./test2.md")
	if err != nil {
		t.Log(err)
	}
	contents, err = SplitByHeading(bytes.NewReader(testMarkdown2))
	if err != nil {
		t.Log(err)
	}
	for _, content := range contents {
		fmt.Print(string(content))
	}
}
