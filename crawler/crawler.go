package crawler

import (
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Crawler is a crawler
type Crawler struct {
	Dispatcher *Dispatcher
}

func NewCrawler(maxWorkers int) *Crawler {
	return &Crawler{
		Dispatcher: NewDispatcher(maxWorkers),
	}
}

// Crawl fetches provided URL and returns <a> URLs found on the page
func (crawler *Crawler) Crawl(urls []string) ([]string, error) {
	if len(urls) == 0 {
		return nil, errors.New("No URLs provided")
	}
	resultsChannel := make(chan JobResults)
	var wg sync.WaitGroup
	for _, url := range urls {
		job := Job{
			URL:            url,
			ResultsChannel: resultsChannel,
			OnCompleted:    func() { wg.Done() },
		}
		crawler.Dispatcher.JobQueue <- job
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		close(resultsChannel)
	}()

	resultURLs := []string{}
	// Wait for results (timeout error can be returned)
	for jobResults := range resultsChannel {
		if jobResults.Error != nil {
			return nil, errors.Wrap(jobResults.Error, "Crawl->job.ResultsChannel")
		}
		resultURLs = append(resultURLs, jobResults.URLs...)
	}
	return resultURLs, nil
}

// RegexAnchorLinks is a regexp for <a href="...">
var RegexAnchorLinks = regexp.MustCompile(`<a[^>]+\bhref=["']([^"']+)["']`)

// findHTMLLinks parses <a href="..."> links from provided HTML
func findHTMLLinks(html []byte) []string {
	parsedLinks := RegexAnchorLinks.FindAllStringSubmatch(string(html), -1)
	// spew.Dump(parsedLinks)
	resultLinks := make([]string, len(parsedLinks))
	seen := map[string]bool{}
	for i, parsedLinkData := range parsedLinks {
		if len(parsedLinkData) < 2 {
			continue
		}
		if seen[parsedLinkData[1]] {
			continue
		}
		resultLinks[i] = parsedLinkData[1]
		seen[parsedLinkData[1]] = true
	}
	// spew.Dump(resultLinks)
	return resultLinks
}

// RegexBase is a regexp for <base href="...">
var RegexBase = regexp.MustCompile(`<base[^>]+\bhref=["']([^"']+)["']`)

// findBase returns <base href="..."> link if any
func findBase(html []byte) string {
	parsedBase := RegexBase.FindAllStringSubmatch(string(html), -1)
	// spew.Dump(parsedBase)
	if len(parsedBase) == 0 {
		return ""
	}
	return parsedBase[0][1]
}

// convertRelativeLinks converts relative links to absolute by adding base href or domain
func convertRelativeLinks(startURL, base string, links []string) ([]string, error) {
	if len(links) == 0 {
		return links, nil
	}
	// Parse start URL to prepare schemeAndHost
	u, err := url.Parse(startURL)
	if err != nil {
		return links, errors.Wrap(err, "convertRelativeLinks->url.Parse")
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	schemeAndHost := u.String()
	// log.Printf("schemeAndHost - %s", schemeAndHost)

	// Filtering without allocating
	resultLinks := links[:0]

	for _, link := range links {
		// Skip empty, XSD and hash links
		if link == "" || strings.HasPrefix(link, "#") || strings.HasSuffix(link, ".xsd") {
			continue
		}
		// Add full links as is (same host only)
		if strings.HasPrefix(link, "http") {
			if !strings.Contains(link, u.Host) {
				continue
			}
			resultLinks = append(resultLinks, link)
			continue
		}
		// Add schemeAndHost if links starts with slash
		if strings.HasPrefix(link, "/") {
			resultLinks = append(resultLinks, schemeAndHost+link)
		} else {
			// Relative link with no slash - use base if defined
			if base != "" {
				resultLinks = append(resultLinks, base+link)
			} else {
				resultLinks = append(resultLinks, schemeAndHost+"/"+link)
			}
		}
	}

	return resultLinks, nil
}
