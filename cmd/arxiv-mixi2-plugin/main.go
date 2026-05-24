package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/arxiv"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/config"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/mixi2"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/posttext"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/state"
)

const defaultCategories = "math.CT,math.AG,math.AT,math.RT,math.NT,math.AC,math.KT,math.OA,math.FA,math.RA"

type posterFactory func(config.Credentials) (mixi2.Poster, error)

func main() {
	var (
		categoriesFlag    = flag.String("categories", defaultCategories, "comma-separated arXiv categories")
		statePath         = flag.String("state", "data/posted.json", "path to posted state JSON")
		dryRun            = flag.Bool("dry-run", false, "print planned posts without posting or saving state")
		initializeOnly    = flag.Bool("initialize-only", false, "record fetched papers without posting")
		initializeOnEmpty = flag.Bool("initialize-on-empty", true, "record fetched papers without posting when state is empty")
		postInterval      = flag.Duration("post-interval", 4*time.Second, "delay between posts in one category")
		requestTimeout    = flag.Duration("request-timeout", 30*time.Second, "timeout for arXiv and mixi2 requests")
	)
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	categories := parseCategories(*categoriesFlag)
	if len(categories) == 0 {
		log.Fatal("no categories configured")
	}

	store, err := state.Load(*statePath)
	if err != nil {
		log.Fatalf("load state: %v", err)
	}
	initialRun := store.IsEmpty()
	if initialRun && *initializeOnEmpty {
		log.Printf("state is empty; running in initialize-only mode")
	}

	ctx := context.Background()
	httpClient := &http.Client{Timeout: *requestTimeout}
	shouldSave := !*dryRun
	hadError := false

	for _, category := range categories {
		initializeCategory := *initializeOnly || (initialRun && *initializeOnEmpty)
		if err := runCategory(ctx, category, store, httpClient, *dryRun, initializeCategory, *postInterval, *requestTimeout, newMixi2Poster); err != nil {
			hadError = true
			log.Printf("category %s failed: %v", category, err)
		}
	}

	if shouldSave {
		if err := state.Save(*statePath, store); err != nil {
			log.Fatalf("save state: %v", err)
		}
		log.Printf("saved state to %s", *statePath)
	}

	if hadError {
		log.Printf("completed with category errors")
		os.Exit(1)
	}
}

func newMixi2Poster(creds config.Credentials) (mixi2.Poster, error) {
	return mixi2.New(creds)
}

func runCategory(ctx context.Context, category string, store state.Store, httpClient arxiv.HTTPClient, dryRun, initializeOnly bool, postInterval, requestTimeout time.Duration, newPoster posterFactory) error {
	papers, err := arxiv.FetchNew(ctx, httpClient, category)
	if err != nil {
		return err
	}
	log.Printf("%s: fetched %d new submissions", category, len(papers))

	var pending []arxiv.Paper
	for _, paper := range papers {
		if !store.Has(category, paper.ID) {
			pending = append(pending, paper)
		}
	}
	log.Printf("%s: %d papers pending", category, len(pending))
	if len(pending) == 0 {
		return nil
	}

	if initializeOnly {
		for _, paper := range pending {
			log.Printf("%s: initialize %s", category, paper.ID)
			if !dryRun {
				store.Mark(category, paper.ID)
			}
		}
		return nil
	}

	if dryRun {
		for _, paper := range pending {
			log.Printf("%s: dry-run post %s: %s", category, paper.ID, posttext.Format(paper))
		}
		return nil
	}

	creds, err := config.LoadCredentials(category)
	if err != nil {
		return err
	}
	poster, err := newPoster(creds)
	if err != nil {
		return err
	}
	defer poster.Close()

	for i, paper := range pending {
		text := posttext.Format(paper)
		postCtx, cancel := context.WithTimeout(ctx, requestTimeout)
		err := poster.Post(postCtx, text)
		cancel()
		if err != nil {
			return fmt.Errorf("post %s: %w", paper.ID, err)
		}
		store.Mark(category, paper.ID)
		log.Printf("%s: posted %s", category, paper.ID)
		if i < len(pending)-1 {
			time.Sleep(postInterval)
		}
	}
	return nil
}

func parseCategories(s string) []string {
	parts := strings.Split(s, ",")
	var categories []string
	seen := map[string]bool{}
	for _, part := range parts {
		category := strings.TrimSpace(part)
		if category != "" && !seen[category] {
			categories = append(categories, category)
			seen[category] = true
		}
	}
	return categories
}
