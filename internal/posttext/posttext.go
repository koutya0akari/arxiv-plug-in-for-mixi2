package posttext

import (
	"fmt"
	"strings"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/arxiv"
)

const MaxRunes = 149

func Format(p arxiv.Paper) string {
	url := "https://arxiv.org/abs/" + p.ID
	prefix := fmt.Sprintf("[%s] ", p.Category)
	author := firstAuthor(p.Authors)
	suffix := " / " + author + " " + url
	if author == "" {
		suffix = " " + url
	}

	text := prefix + p.Title + suffix
	if runeLen(text) <= MaxRunes {
		return text
	}

	budget := MaxRunes - runeLen(prefix) - runeLen(suffix)
	if budget < 8 {
		suffix = " " + url
		budget = MaxRunes - runeLen(prefix) - runeLen(suffix)
	}
	if budget <= 0 {
		return trimRunes(prefix+url, MaxRunes)
	}
	return prefix + ellipsize(p.Title, budget) + suffix
}

func firstAuthor(authors []string) string {
	if len(authors) == 0 {
		return ""
	}
	name := strings.TrimSpace(authors[0])
	if len(authors) > 1 && name != "" {
		return name + " et al."
	}
	return name
}

func ellipsize(s string, limit int) string {
	if runeLen(s) <= limit {
		return s
	}
	if limit <= 3 {
		return trimRunes(s, limit)
	}
	return strings.TrimSpace(trimRunes(s, limit-3)) + "..."
}

func trimRunes(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit])
}

func runeLen(s string) int {
	return len([]rune(s))
}
