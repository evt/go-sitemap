package sitemap

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sitemap/crawler"
	"time"

	"github.com/pkg/errors"
)

const (
	DefaultSitemapFile   = "sitemap.xml"
	SiteMapLastModFormat = "2006-01-02"
	DefaultWorkers       = 3
	DefaultMaxDepth      = 2
)

// GenerateSitemap generates site map by site URL
func GenerateSitemap(startURL, file string, workers, maxDepth, currentLevel int) error {
	if startURL == "" {
		return errors.New("No start URL provided")
	}
	if file == "" {
		file = GetDefaultSitemapFile(startURL)
	}
	if workers == 0 {
		workers = DefaultWorkers
	}
	if maxDepth == 0 {
		maxDepth = DefaultMaxDepth
	}

	// List of result links
	sitemapLinks := map[string]bool{}
	// List of crawled links (to skip crawling the same link twice)
	crawledLinks := map[string]bool{}

	generateSitemap(crawler.NewCrawler(workers), []string{startURL}, file, workers, maxDepth, currentLevel, sitemapLinks, crawledLinks)

	if len(sitemapLinks) == 0 {
		return errors.New("No links found")
	}

	sitemapURLs := make([]SiteMapURL, 0, len(sitemapLinks))

	for link := range sitemapLinks {
		sitemapURLs = append(sitemapURLs, SiteMapURL{
			Loc:     link,
			LastMod: time.Now().UTC().Format(SiteMapLastModFormat),
		})
	}

	sitemap := &SiteMap{URLs: sitemapURLs}
	sitemapBytes, err := xml.MarshalIndent(sitemap, "", "\t")
	if err != nil {
		return err
	}

	// Create path
	os.MkdirAll(filepath.Dir(file), 0755)

	log.Printf("Writing sitemap to %s", file)

	return ioutil.WriteFile(file, sitemapBytes, 0755)
}

// generateSitemap generates site map by site URL
func generateSitemap(crawler *crawler.Crawler, urls []string, file string, workers, maxDepth, currentLevel int, sitemapLinks, crawledLinks map[string]bool) error {
	// Do nothing if already crawled
	filteredURLs := urls[:0]
	for _, url := range urls {
		if crawledLinks[url] {
			log.Printf("URL %s has been already crawled", url)
			continue
		}
		filteredURLs = append(filteredURLs, url)
	}
	// Fetch it
	links, err := crawler.Crawl(filteredURLs)
	if err != nil {
		return errors.Wrap(err, "generateSitemap->client.Crawl")
	}
	// Track crawled URL
	for _, url := range urls {
		crawledLinks[url] = true
	}
	log.Printf("[Level %d] %d links found", currentLevel, len(links))
	// Save found links
	for _, link := range links {
		sitemapLinks[link] = true
	}
	// Crawl found URLs if recursion level hasn't been reached yet
	if maxDepth > currentLevel {
		generateSitemap(crawler, links, file, workers, maxDepth, currentLevel+1, sitemapLinks, crawledLinks)
	}
	return nil
}

func GetDefaultSitemapFile(startURL string) string {
	if startURL == "" {
		return DefaultSitemapFile
	}
	u, err := url.Parse(startURL)
	if err != nil {
		log.Println(err)
		return DefaultSitemapFile
	}
	return u.Host + "_" + "sitemap.xml"
}
