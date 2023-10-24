package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type BlockType int

type Block struct {
	Kind    BlockType
	Content []byte
}

const (
	BlockKindDefault BlockType = iota
	BlockKindBlankLine
	BlockKindMetadata
	BlockKindHeading
	BlockKindCode
	BlockKindParagraph
)

// TranslateParse
// metadata: 第一行至少是三个连续的短横线，结束至少是三个连续的短横线
func TranslateParse(reader io.Reader) ([]Block, error) {
	// lines
	var lines [][]byte

	// Create a scanner
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		lines = append(lines, line)
	}

	var blocks []Block
	var block Block
	//var flag bool
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		fmt.Print(string(line))

		// metadata begin
		if i == 0 && bytes.HasPrefix(line, []byte("---")) {
			block.Kind = BlockKindMetadata
			block.Content = append(line, '\n')
			continue
		}
		// metadata
		if block.Kind == BlockKindMetadata {
			block.Content = append(block.Content, line...)
			block.Content = append(block.Content, '\n')
			if bytes.HasPrefix(line, []byte("---")) {
				blocks = append(blocks, block)
				block = Block{}
			}
			continue
		}

		// paragraph begin
		if block.Kind = BlockKindDefault; len(line) > 0 {
			block.Kind = BlockKindParagraph
			block.Content = append(block.Content, line...)
			block.Content = append(block.Content, '\n')
		}
		if block.Kind == BlockKindParagraph && len(line) == 0 {
			blocks = append(blocks, block)
			block = Block{}
			continue
		}

		// blank line
		if len(line) == 0 {
			block.Kind = BlockKindBlankLine
			block.Content = append(block.Content, '\n')
			blocks = append(blocks, block)
			block = Block{}
			continue
		}

		// heading
		if bytes.HasPrefix(line, []byte("#")) {
			block.Kind = BlockKindHeading
			block.Content = append(block.Content, line...)
			block.Content = append(block.Content, '\n')
			blocks = append(blocks, block)
			block = Block{}
			continue
		}

	}
	return blocks, nil

	//
	//
	//var prevLine, currLine []byte
	//
	//for {
	//	line, err := r.ReadSlice('\n')
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		return nil, err
	//	}
	//	prevLine = currLine
	//	currLine = line
	//	nextLinePeek, err := r.Peek(1)
	//
	//	if err != nil {
	//		if err == io.EOF {
	//
	//		} else {
	//			return nil, err
	//		}
	//	}
	//	fmt.Printf("prevLine: %s, currLine: %s, nextLinePeek: %s\n", prevLine, currLine, nextLinePeek)
	//
	//
	//
	//}

}
