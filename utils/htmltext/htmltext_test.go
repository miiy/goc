package htmltext

import (
	"strings"
	"testing"
)

func TestText(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "paragraphs",
			raw:  "<p>First</p><p>Second</p>",
			want: "First Second ",
		},
		{
			name: "skip scripts",
			raw:  "<p>First</p><script>alert(1)</script><p>Second</p>",
			want: "First Second ",
		},
		{
			name: "headings and links",
			raw:  "<h1>Title</h1><p>Hello <a href=\"/about\">world</a></p>",
			want: "Title Hello world ",
		},
		{
			name: "invalid html",
			raw:  "<p>Hello",
			want: "Hello ",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := strings.Join(strings.Fields(Text(tc.raw)), " ")
			want := strings.TrimSpace(tc.want)
			if got != want {
				t.Fatalf("Text() = %q, want %q", got, want)
			}
		})
	}
}
