package markdown

import (
	"bufio"
	"bytes"
	"io"
)

func SplitByHeading(reader io.Reader) ([][]byte, error) {
	// Create a reader
	r := bufio.NewReader(reader)

	var contents [][]byte
	var curr = new(bytes.Buffer)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			contents = append(contents, curr.Bytes())
			break
		}
		if err != nil {
			return contents, err
		}

		if bytes.HasPrefix(line, []byte("#")) {
			contents = append(contents, curr.Bytes())
			curr = new(bytes.Buffer)
		}
		curr.Write(line)
		curr.WriteByte('\n')
	}

	return contents, nil
}
