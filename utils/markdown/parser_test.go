package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func TestTranslateParse(t *testing.T) {
	blocks, err := TranslateParse(strings.NewReader(testMarkdown))
	if err != nil {
		t.Error(err)
	}
	for _, block := range blocks {
		fmt.Print(string(block.Content))
	}
	t.Log(blocks)
}
