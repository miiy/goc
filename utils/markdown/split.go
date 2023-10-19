package markdown

import (
	"regexp"
	"strings"
)

func Split(markdown string) []string {
	// Define the regular expression for matching headings
	regex := regexp.MustCompile(`(?m)^#+\s+(.*)$`)

	// Find all the headings using the regular expression
	matches := regex.FindAllStringSubmatch(markdown, -1)

	var contents []string
	startIndex := 0
	// Iterate over the matches
	for _, match := range matches {
		titleIndex := strings.Index(markdown, match[0])
		// if the first line is Heading
		if titleIndex == 0 {
			continue
		}
		content := markdown[startIndex:titleIndex]
		startIndex = titleIndex
		contents = append(contents, content)
	}
	if startIndex < len(markdown) {
		content := markdown[startIndex:]
		contents = append(contents, content)
	}

	//fmt.Println(strings.Join(contents, ""))
	return contents
}
