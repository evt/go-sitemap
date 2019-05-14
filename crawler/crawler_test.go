package crawler

import (
	"testing"
)

// ClientTest is a file and it's test parser
type ClientTest struct {
	urls       []string
	resultURLs []string
	maxWorkers int
}

// TestTable is a list of tests to run
var TestTable = []ClientTest{
	ClientTest{
		urls: []string{"https://www.sitemaps.org", "https://www.sitemaps.org/faq.php"},
		resultURLs: []string{
			"https://www.sitemaps.org/faq.php",
			"https://www.sitemaps.org/protocol.php",
		},
	},
}

func TestCrawl(t *testing.T) {
	for i, tt := range TestTable {
		resultURLs, err := NewCrawler(len(tt.urls)).Crawl(tt.urls)
		if err != nil {
			t.Errorf("#%d: %v", i, err)
		}
		// Prepare index of found URLs
		resultURLIndex := map[string]bool{}
		for _, resultURL := range resultURLs {
			resultURLIndex[resultURL] = true
		}
		// Check all golden URLs have been found
		for _, goldenResultURL := range tt.resultURLs {
			if resultURLIndex[goldenResultURL] {
				t.Logf("#%d: Golden URL %s has been found by crawler", i, goldenResultURL)
			} else {
				t.Errorf("#%d: Golden URL %s wasn't found by crawler", i, goldenResultURL)
			}
		}
	}
}
