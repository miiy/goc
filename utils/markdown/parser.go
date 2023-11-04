package markdown

import (
	"bytes"
)

type BlockType int

type Block struct {
	Kind    BlockType
	Content []byte
}

type Parser struct {
	Blocks []Block
}

const (
	BlockKindDefault BlockType = iota
	BlockKindBlankLine
	BlockKindMetadata
	BlockKindHeading
	BlockKindCode
	BlockKindParagraph
)

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseBlock(data []byte) ([]Block, error) {
	for len(data) > 0 {
		// metadata
		if len(p.Blocks) == 0 && isMetadata(data) {
			end := p.metadata(data)
			data = data[end:]
			continue
		}

		// empty line
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

		// fenced code block
		if isFenceLine(data) {
			end := p.fencedCodeBlock(data)
			data = data[end:]
			continue
		}

		// indent code block
		if isCodePrefix(data) {
			end := p.code(data)
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
	// the whole line must be the char
	n := 0
	for i := 0; i < len(data) && data[i] != '\n'; i++ {
		switch {
		case data[i] == '-':
			n++
		case data[i] == '\r': // ---/r/n
		case data[i] != '-':
			return false
		}
	}

	return n >= 3
}

// metadata
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

// isFenceLine
func isFenceLine(data []byte) bool {
	i := 0
	// skip up to three spaces
	for i < len(data) && i < 3 && data[i] == ' ' {
		i++
	}

	if i >= len(data) {
		return false
	}

	if data[i] != '`' && data[i] != '~' {
		return false
	}

	c := data[i]

	// the whole line must be the char or whitespace
	n := 0
	for i < len(data) && data[i] == c {
		n++
		i++
	}

	return n >= 3
}

func (p *Parser) fencedCodeBlock(data []byte) int {
	// first line ```
	end := bytes.IndexByte(data, '\n') + 1
	beg := end
	for end < len(data) {
		end = bytes.IndexByte(data[beg:], '\n') + 1
		if isFenceLine(data[beg:]) {
			break
		}
		beg += end
	}
	p.addBlock(BlockKindCode, data[:beg+end])
	return beg + end
}

func isCodePrefix(data []byte) bool {
	if len(data) >= 1 && data[0] == '\t' {
		return true
	}
	if len(data) >= 4 && data[0] == ' ' && data[1] == ' ' && data[2] == ' ' && data[3] == ' ' {
		return true
	}
	return false
}

func (p *Parser) code(data []byte) int {
	i := 0
	for i < len(data) {
		end := bytes.IndexByte(data[i:], '\n') + 1
		if isCodePrefix(data[i:]) {
			i += end
			continue
		}
		if isEmpty(data[i:]) {
			i += end
			continue
		}
		break
	}
	p.addBlock(BlockKindCode, data[:i])
	return i
}

func (p *Parser) paragraph(data []byte) int {
	var i, end int
	for i <= len(data) {
		end = bytes.IndexByte(data[i:], '\n') + 1
		if end == 0 {
			end = len(data[i:]) // end without \n
		}
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
		if data[i] != ' ' && data[i] != '\r' && data[i] != '\t' {
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
