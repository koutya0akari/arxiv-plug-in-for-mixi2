package arxiv

import (
	"strings"
	"testing"
)

func TestParseNewOnlyNewSubmissions(t *testing.T) {
	html := `
<html><body>
<h3>New submissions</h3>
<dl>
  <dt><span class="list-identifier"><a href="/abs/2501.00001">arXiv:2501.00001</a></span></dt>
  <dd>
    <div class="list-title mathjax">Title: First Paper</div>
    <div class="list-authors">Authors: <a>Jane Doe</a>, <a>John Roe</a></div>
  </dd>
  <dt><span class="list-identifier"><a href="/abs/2501.00002">arXiv:2501.00002</a></span></dt>
  <dd>
    <div class="list-title mathjax">Title: Second   Paper</div>
    <div class="list-authors">Authors: <a>Ada Lovelace</a></div>
  </dd>
</dl>
<h3>Cross submissions</h3>
<dl>
  <dt><span class="list-identifier"><a href="/abs/2501.99999">arXiv:2501.99999</a></span></dt>
  <dd><div class="list-title mathjax">Title: Cross Listed</div></dd>
</dl>
<h3>Replacements</h3>
<dl>
  <dt><span class="list-identifier"><a href="/abs/2501.88888">arXiv:2501.88888</a></span></dt>
  <dd><div class="list-title mathjax">Title: Replacement</div></dd>
</dl>
</body></html>`

	papers, err := ParseNew(strings.NewReader(html), "math.CT")
	if err != nil {
		t.Fatalf("ParseNew() error = %v", err)
	}
	if got, want := len(papers), 2; got != want {
		t.Fatalf("len(papers) = %d, want %d", got, want)
	}
	if papers[0].ID != "2501.00001" || papers[0].Title != "First Paper" || papers[0].Category != "math.CT" {
		t.Fatalf("unexpected first paper: %+v", papers[0])
	}
	if got, want := len(papers[0].Authors), 2; got != want {
		t.Fatalf("len(authors) = %d, want %d", got, want)
	}
	if papers[1].Title != "Second Paper" {
		t.Fatalf("normalized title = %q", papers[1].Title)
	}
}

func TestParseNewWhenHeadingIsInsideDL(t *testing.T) {
	html := `
<html><body>
<dl id="articles">
  <h3>New submissions (showing 1 of 1 entries)</h3>
  <dt><span class="list-identifier"><a href="/abs/2605.00773">arXiv:2605.00773</a></span></dt>
  <dd>
    <div class="list-title mathjax"><span class="descriptor">Title:</span> The Synthetic Sierpinski Cone</div>
    <div class="list-authors">Authors: <a>Fredrik Bakke</a>, <a>Jonathan Sterling</a></div>
  </dd>
  <h3>Cross submissions (showing 1 of 1 entries)</h3>
  <dt><span class="list-identifier"><a href="/abs/2605.00097">arXiv:2605.00097</a></span></dt>
  <dd><div class="list-title mathjax"><span class="descriptor">Title:</span> Cross Listed</div></dd>
</dl>
</body></html>`

	papers, err := ParseNew(strings.NewReader(html), "math.CT")
	if err != nil {
		t.Fatalf("ParseNew() error = %v", err)
	}
	if got, want := len(papers), 1; got != want {
		t.Fatalf("len(papers) = %d, want %d", got, want)
	}
	if papers[0].ID != "2605.00773" {
		t.Fatalf("paper ID = %q", papers[0].ID)
	}
}
