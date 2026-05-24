package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/config"
	"github.com/koutya0akari/arxiv-plug-in-for-mixi2/internal/mixi2"
)

const defaultCategories = "math.CT,math.AG,math.AT,math.RT,math.NT,math.AC,math.KT,math.OA,math.FA,math.RA"

func main() {
	var (
		categoriesFlag = flag.String("categories", defaultCategories, "comma-separated arXiv categories")
		requestTimeout = flag.Duration("request-timeout", 30*time.Second, "timeout for mixi2 requests")
	)
	flag.Parse()

	log.SetFlags(0)
	categories := parseCategories(*categoriesFlag)
	if len(categories) == 0 {
		log.Fatal("no categories configured")
	}

	fmt.Println("category\tcommunity_id\tname\tapplication_version_id")
	hadError := false
	for _, category := range categories {
		if err := listCategory(category, *requestTimeout); err != nil {
			hadError = true
			log.Printf("%s: %v", category, err)
		}
	}
	if hadError {
		os.Exit(1)
	}
}

func listCategory(category string, requestTimeout time.Duration) error {
	creds, err := config.LoadApplicationCredentials(category)
	if err != nil {
		return err
	}
	client, err := mixi2.New(creds)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	communities, err := client.Communities(ctx)
	if err != nil {
		return err
	}
	for _, community := range communities {
		fmt.Printf("%s\t%s\t%s\t%s\n", category, community.ID, community.Name, community.ApplicationVersionID)
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
