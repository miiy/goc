package markdown

import (
	"bufio"
	"io"
	"strings"
)

func SplitByHeading(reader io.Reader) ([]string, error) {
	// Create a reader
	r := bufio.NewReader(reader)

	var contents []string
	var currContent strings.Builder
	for {
		b, _, err := r.ReadLine()
		if err == io.EOF {
			contents = append(contents, currContent.String())
			break
		}
		line := string(b) + "\n"
		if strings.HasPrefix(line, "#") {
			contents = append(contents, currContent.String())
			currContent.Reset()
		}
		currContent.WriteString(line)
	}

	return contents, nil
}
