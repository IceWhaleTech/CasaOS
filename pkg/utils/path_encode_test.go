package utils

import "testing"

// TestEncodePathEscapesPercentOnce pins the non-all escape contract: a literal
// percent sign must be escaped exactly once. Applying the "%" -> "%25" rule
// twice turns "a%b" into "a%2525b" and corrupts paths used in URLs.
func TestEncodePathEscapesPercentOnce(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "plain", "plain"},
		{"single percent", "a%b", "a%25b"},
		{"trailing percent", "100%", "100%25"},
		{"percent with reserved chars", "a%?b#c", "a%25%3Fb%23c"},
		{"reserved chars", "a?b#c", "a%3Fb%23c"},
		{"percent in each segment", "x%y/a%b", "x%25y/a%25b"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := EncodePath(tc.in); got != tc.want {
				t.Fatalf("EncodePath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestEncodePathAllModeStillEscapesPercent(t *testing.T) {
	if got, want := EncodePath("a b%c", true), "a%20b%25c"; got != want {
		t.Fatalf("EncodePath(%q, true) = %q, want %q", "a b%c", got, want)
	}
}
