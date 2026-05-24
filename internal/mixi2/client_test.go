package mixi2

import (
	"context"
	"errors"
	"testing"
)

type fakePoster struct {
	posted []string
	err    error
}

func (f *fakePoster) Post(_ context.Context, text string) error {
	if f.err != nil {
		return f.err
	}
	f.posted = append(f.posted, text)
	return nil
}

func (f *fakePoster) Close() error { return nil }

func TestPosterInterface(t *testing.T) {
	var poster Poster = &fakePoster{}
	if err := poster.Post(context.Background(), "hello"); err != nil {
		t.Fatalf("Post() error = %v", err)
	}
}

func TestFakePosterFailure(t *testing.T) {
	poster := &fakePoster{err: errors.New("boom")}
	if err := poster.Post(context.Background(), "hello"); err == nil {
		t.Fatal("Post() error = nil, want error")
	}
}
