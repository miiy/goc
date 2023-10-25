package markdown

import (
	"bytes"
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

func (p *Parser) ParseBlock(data []byte) ([]Block, error) {
	for len(data) > 0 {

		if isEmpty(data) {
			end := p.blankLine(data)
			data = data[end:]
			continue
		}

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

		end := p.paragraph(data)
		data = data[end:]
	}

	return p.Blocks, nil
}

func (p *Parser) addBlock(kind BlockType, data []byte) {
	p.Blocks = append(p.Blocks, Block{
		Kind:    kind,
		Content: data,
	})
}

func (p *Parser) blankLine(data []byte) int {
	end := bytes.IndexByte(data, '\n') + 1
	p.addBlock(BlockKindBlankLine, data[:end])
	return end
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

// metadata
// ---\n
// aaa\n
// ---\n
func (p *Parser) metadata(data []byte) int {
	// first line ---
	end := bytes.IndexByte(data, '\n') + 1
	beg := end
	for end < len(data) {
		end = bytes.IndexByte(data[beg:], '\n') + 1
		if isMetadata(data[beg:]) {
			break
		}
		beg += end
	}
	p.addBlock(BlockKindMetadata, data[:beg+end])
	return beg + end
}

func isPrefixHeading(data []byte) bool {
	i := skipChar(data, 0, ' ')
	if data[i] != '#' {
		return false
	}
	return true
}

func (p *Parser) prefixHeading(data []byte) int {
	end := bytes.IndexByte(data, '\n') + 1
	p.addBlock(BlockKindHeading, data[:end])
	return end
}

func (p *Parser) paragraph(data []byte) int {
	var i, end int
	for i <= len(data) {
		end = bytes.IndexByte(data[i:], '\n') + 1
		if isEmpty(data[i:]) {
			p.addBlock(BlockKindParagraph, data[:i])
			return i
		}
		i += end
	}
	return i
}

func isEmpty(data []byte) bool {
	for i := 0; i < len(data) && data[i] != '\n'; i++ {
		if data[i] != ' ' && data[i] != '\t' {
			return false
		}
	}
	return true
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
