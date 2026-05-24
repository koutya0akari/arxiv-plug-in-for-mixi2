package posttext

import (
	"strings"
	"testing"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/arxiv"
)

func TestFormatShortPaper(t *testing.T) {
	text := Format(arxiv.Paper{
		ID:       "2501.00001",
		Category: "math.CT",
		Title:    "A Short Title",
		Authors:  []string{"Jane Doe", "John Roe"},
	})
	want := "[math.CT] A Short Title / Jane Doe et al. https://arxiv.org/abs/2501.00001"
	if text != want {
		t.Fatalf("Format() = %q, want %q", text, want)
	}
}

func TestFormatLongPaperKeepsURLAndLimit(t *testing.T) {
	text := Format(arxiv.Paper{
		ID:       "2501.12345",
		Category: "math.NT",
		Title:    strings.Repeat("Very Long Title ", 20),
		Authors:  []string{"A Very Long Author Name", "Another Author"},
	})
	if got := len([]rune(text)); got > MaxRunes {
		t.Fatalf("len(text) = %d, want <= %d: %q", got, MaxRunes, text)
	}
	if !strings.Contains(text, "https://arxiv.org/abs/2501.12345") {
		t.Fatalf("text does not keep URL: %q", text)
	}
}
