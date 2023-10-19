package markdown

import "testing"

var testMarkdown = `
---
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
`

func TestSplit(t *testing.T) {
	Split(testMarkdown)
	Split(testMarkdown2)
}
