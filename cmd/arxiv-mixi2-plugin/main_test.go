package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/config"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/mixi2"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/state"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type testPoster struct {
	err error
}

func (p *testPoster) Post(context.Context, string) error { return p.err }
func (p *testPoster) Close() error                       { return nil }

func TestRunCategoryInitializeOnlyMarksWithoutCredentials(t *testing.T) {
	store := state.Store{}
	client := htmlClient(`<h3>New submissions</h3><dl>
<dt><span class="list-identifier"><a href="/abs/2501.00001">arXiv:2501.00001</a></span></dt>
<dd><div class="list-title">Title: Paper</div><div class="list-authors"><a>Jane Doe</a></div></dd>
</dl>`)

	err := runCategory(context.Background(), "math.CT", store, client, false, true, 0, time.Second, func(config.Credentials) (mixi2.Poster, error) {
		t.Fatal("poster factory should not be called in initialize-only mode")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("runCategory() error = %v", err)
	}
	if !store.Has("math.CT", "2501.00001") {
		t.Fatalf("state was not marked: %+v", store)
	}
}

func TestRunCategoryPostFailureDoesNotMark(t *testing.T) {
	t.Setenv("MIXI2_MATH_CT_CLIENT_ID", "client-id")
	t.Setenv("MIXI2_MATH_CT_CLIENT_SECRET", "client-secret")
	t.Setenv(config.TokenURLEnv, "https://token.example")
	t.Setenv(config.APIAddressEnv, "api.example:443")
	t.Setenv(config.CommunityIDEnv, "community-id")
	store := state.Store{}
	client := htmlClient(`<h3>New submissions</h3><dl>
<dt><span class="list-identifier"><a href="/abs/2501.00001">arXiv:2501.00001</a></span></dt>
<dd><div class="list-title">Title: Paper</div><div class="list-authors"><a>Jane Doe</a></div></dd>
</dl>`)

	err := runCategory(context.Background(), "math.CT", store, client, false, false, 0, time.Second, func(config.Credentials) (mixi2.Poster, error) {
		return &testPoster{err: errors.New("post failed")}, nil
	})
	if err == nil {
		t.Fatal("runCategory() error = nil, want post error")
	}
	if store.Has("math.CT", "2501.00001") {
		t.Fatalf("failed post was marked: %+v", store)
	}
}

func htmlClient(html string) roundTripFunc {
	return func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewBufferString(html)),
		}, nil
	}
}
