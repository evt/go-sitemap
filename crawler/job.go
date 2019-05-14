package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// JobTimeout is a job timeout
const JobTimeout = 60

// Job is a job to run
type Job struct {
	// Context        context.Context
	URL            string
	ResultsChannel chan JobResults
	TimeoutChannel chan struct{}
	OnCompleted    JobCallback
}

// JobCallback is a callback called once job completed
type JobCallback func()

// JobResults
type JobResults struct {
	URLs  []string
	Error error
}

// Run runs a job
func (job Job) Run(workerSeq int) {
	log.Printf("[Worker %d] Crawling %s", workerSeq, job.URL)

	done := make(chan struct{})

	go func() {
		urls, err := Crawl(job.URL)
		job.ResultsChannel <- JobResults{URLs: urls, Error: err}
		close(done)
	}()

	select {
	case <-done:
		log.Printf("[Worker %d] Completed", workerSeq)
		if job.OnCompleted != nil {
			job.OnCompleted()
		}
	case <-time.After(time.Second * time.Duration(JobTimeout)):
		log.Printf("[Worker %d] Timeout %d sec", workerSeq, JobTimeout)
		if job.OnCompleted != nil {
			job.OnCompleted()
		}
		close(job.TimeoutChannel)
	}
}

// HTTPTimeout is HTTP client timeout in seconds
const HTTPTimeout = 30

// HTTPClient is a client with timeout set (default one doesn't have it and may hang forever)
var HTTPClient = &http.Client{Timeout: time.Duration(HTTPTimeout) * time.Second}

// Crawl fetches provided URL and returns <a> URLs found on the page
func Crawl(startURL string) ([]string, error) {
	// Fetch it
	response, err := HTTPClient.Get(startURL)
	if err != nil {
		return nil, errors.Wrap(err, "Crawl->HttpClient.Get")
	}
	defer response.Body.Close()
	// Read body
	body, err := ioutil.ReadAll(response.Body)
	// Parse links
	links := findHTMLLinks(body)
	// Parse <base> tag
	base := findBase(body)
	// Convert relative links to absolute
	resultLinks, err := convertRelativeLinks(startURL, base, links)
	if err != nil {
		return nil, errors.Wrap(err, "Crawl->convertRelativeLinks")
	}
	// spew.Dump(resultLinks)
	return resultLinks, nil
}
