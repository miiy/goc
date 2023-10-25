package markdown

import (
	"fmt"
	"testing"
)

func TestParseBlock(t *testing.T) {
	var parser Parser
	blocks, err := parser.ParseBlock([]byte(testMarkdown))
	if err != nil {
		t.Error(err)
	}
	for _, block := range blocks {
		fmt.Println("-----------------------------------------------------", block.Kind)
		fmt.Print(string(block.Content))
	}
	t.Log(blocks)
}
