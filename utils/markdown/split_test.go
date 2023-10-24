package markdown

import (
	"fmt"
	"strings"
	"testing"
)

var testMarkdown = `---
title: Test
linkTitle: Test
weight: 9
description: >
    Test
aliases:
  - /test/test
---

test

# Heading1

content1

## Subheading1.1

content1.1

## Subheading1.2

content1.2

# Heading2

content2
`

var testMarkdown2 = `# Heading1

content1

## Subheading1.1

content1.1

## Subheading1.2

content1.2

# Heading2

# Heading3

content3

# Heading3

content3 repeat

`

func TestSplit(t *testing.T) {
	contents, err := SplitByHeading(strings.NewReader(testMarkdown))
	if err != nil {
		t.Log(err)
	}
	for _, content := range contents {
		fmt.Print(string(content))
	}
	fmt.Println("---")
	contents, err = SplitByHeading(strings.NewReader(testMarkdown2))
	if err != nil {
		t.Log(err)
	}
	for _, content := range contents {
		fmt.Print(string(content))
	}
}
