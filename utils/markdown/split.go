package markdown

import (
	"bufio"
	"bytes"
	"io"
)

func SplitByHeading(reader io.Reader) ([][]byte, error) {
	// Create a scanner
	scanner := bufio.NewScanner(reader)

	var contents [][]byte
	var curr = new(bytes.Buffer)
	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasPrefix(line, []byte("#")) {
			contents = append(contents, curr.Bytes())
			curr = new(bytes.Buffer)
		}

		curr.Write(line)
		curr.WriteByte('\n')
	}

	return contents, nil
}
