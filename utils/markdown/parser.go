package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Parser struct {
	Blocks []Block
}

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

func (p *Parser) addBlock(kind BlockType, data []byte) {
	p.Blocks = append(p.Blocks, Block{
		Kind:    kind,
		Content: data,
	})
}

func (p *Parser) ParseBlock(data []byte) error {
	var i int
	for i < len(data) {
		//current := data[i:]
		////line := i
		//fmt.Print(string(current))

		// heading
		if isPrefixHeading(data) {
			end := p.prefixHeading(data)
			data = data[end:]
			continue
		}

		if isMetadata(data) {
			end := p.metadata(data)
			data = data[end:]
			continue
		}

		// 扫描行
		nl := bytes.IndexByte(data, '\n')
		if nl >= 0 {
			i = i + nl + 1
		}
		fmt.Print(string(data[:i]))
		data = data[i:]
	}

	return nil
}

func isPrefixHeading(data []byte) bool {
	i := skipChar(data, 0, ' ')
	if data[i] != '#' {
		return false
	}
	return true
}

func (p *Parser) prefixHeading(data []byte) int {
	end := bytes.IndexByte(data, '\n')
	p.addBlock(BlockKindMetadata, data[:end])
	return end + 1
}

func isMetadata(data []byte) bool {
	// look at the metadata char
	i := 0
	if data[i] != '-' {
		return false
	}

	// the whole line must be the char or whitespace
	n := 0
	for i < len(data) && data[i] != '\n' {
		switch {
		case data[i] == '-':
			n++
		case data[i] != ' ':
			return false
		}
		i++
	}

	return n >= 3
}

func (p *Parser) metadata(data []byte) int {
	// first line ---
	end := bytes.IndexByte(data, '\n')
	beg := end + 1
	for beg < len(data) {
		end = bytes.IndexByte(data, '\n')
		end = end + end + 1
		line := data[beg:end]
		beg = end + 1
		if isMetadata(line) {
			break
		}
	}
	p.addBlock(BlockKindMetadata, data[:end])
	return end
}

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

func skipChar(data []byte, start int, char byte) int {
	i := start
	for i < len(data) {
		if data[i] != char {
			break
		}
		i++
	}
	return i
}
