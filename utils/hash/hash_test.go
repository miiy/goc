package hash

import "testing"

func TestSHA256Base64URL(t *testing.T) {
	got := SHA256Base64URL("session-id")
	const want = "S98eFd9xbyf_brzBGapLiGOiIc1U6Hdy2CSIj0rqxcA"
	if got != want {
		t.Fatalf("SHA256Base64URL() = %q, want %q", got, want)
	}
}
