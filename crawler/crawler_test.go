package crawler

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/gocolly/colly"
)

func TestCrawlerUserAgent(t *testing.T) {
	var agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59"

	// crawler := colly.NewCollector()
	c := Crawler()

	c.OnResponse(func(h *colly.Response) {
		// Check that the UserAgent of the colly.Collector object is set to the correct value
		// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59"

		if got, want := c.UserAgent, agent; got != want {
			t.Fatalf("c.UserAgent = %q, want %q", got, want)
		}
	})

	c.Visit("https://www.eworldtrade.com/c/page=1108")

}

func TestCrawlerURLs(t *testing.T) {
	c := crawler.Crawler()

	// Set a flag to indicate whether the expected URLs have been visited
	visited := false

	// Set a callback to check the visited URLs
	c.OnRequest(func(h *colly.Request) {
		if h.URL.String() == "https://www.eworldtrade.com/c/?page=1" || h.URL.String() == "https://www.eworldtrade.com/c/?page=2" {
			visited = true
		}
	})

	// Visit some test URLs
	c.Visit("https://www.eworldtrade.com/c/?page=1")
	c.Visit("https://www.eworldtrade.com/c/?page=2")

	// Check that the expected URLs have been visited
	if !visited {
		t.Error("Expected URLs were not visited")
	}
}

func TestCrawlerCompanyURLs(t *testing.T) {
	c := Crawler()

	// Set a flag to indicate whether the expected company URLs have been extracted
	extracted := false

	// Set a callback to extract the company URLs
	c.OnHTML("div.row > div.col-lg-8", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")

		// Check if the extracted company URL is correct
		if link == "https://www.eworldtrade.com/c/abc-company" || link == "https://www.eworldtrade.com/c/xyz-company" {
			extracted = true
		}
	})

	// Visit some test pages
	c.Visit("https://www.eworldtrade.com/c/?page=1")
	c.Visit("https://www.eworldtrade.com/c/?page=2")

	// Check that the expected company URLs have been extracted
	if !extracted {
		t.Error("Expected company URLs were not extracted")
	}
}

func TestScrapeEmails(t *testing.T) {
	c := Crawler()

	// Set a flag to indicate whether the expected emails have been scraped
	scraped := false

	// Set a callback to scrape the emails from the company URLs
	c.OnHTML("div.row > div.col-lg-8", func(h *colly.HTMLElement) {
		companyURL := h.ChildAttr("a", "href")

		// Check if the scraped email is correct
		if companyURL == "https://www.eworldtrade.com/c/abc-company" && scrapeEmails([]string{companyURL})[0] == "abc@company.com" {
			scraped = true
		}
		if companyURL == "https://www.eworldtrade.com/c/xyz-company" && scrapeEmails([]string{companyURL})[0] == "xyz@company.com" {
			scraped = true
		}
	})

	// Visit some test pages
	c.Visit("https://www.eworldtrade.com/c/?page=1")
	c.Visit("https://www.eworldtrade.com/c/?page=2")

	// Check that the expected emails have been scraped
	if !scraped {
		t.Error("Expected emails were not scraped")
	}

}

// TestShowProgress tests that the showProgress function prints the correct progress bar.
func TestShowProgress(t *testing.T) {
	// Capture the output of the showProgress function
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	// Test the showProgress function with different total values
	showProgress(50)
	if got, want := buf.String(), "\r[==================================================] 100%"; got != want {
		t.Errorf("showProgress(50) = %q, want %q", got, want)
	}

	buf.Reset()
	showProgress(75)
	if got, want := buf.String(), "\r[==============================================   ] 75%"; got != want {
		t.Errorf("showProgress(75) = %q, want %q", got, want)
	}

	buf.Reset()
	showProgress(0)
	if got, want := buf.String(), "\r[                                                  ] 0%"; got != want {
		t.Errorf("showProgress(0) = %q, want %q", got, want)
	}

}

// TestRemoveDuplicates tests that the removeDuplicates function returns a slice without duplicate values.
func TestRemoveDuplicates(t *testing.T) {
	// Test the removeDuplicates function with a slice of strings
	slice := []string{"abc@example.com", "xyz@example.com", "abc@example.com"}
	if got, want := removeDuplicates(slice), []string{"abc@example.com", "xyz@example.com"}; !equal(got, want) {
		t.Errorf("removeDuplicates(%q) = %q, want %q", slice, got, want)
	}

	// Test the removeDuplicates function with a slice of integers
	slice = []string{"1", "2", "2"}
	if got, want := removeDuplicates(slice), []string{"1", "2"}; !equal(got, want) {
		t.Errorf("removeDuplicates(%q) = %q, want %q", slice, got, want)
	}

	// Test the removeDuplicates function with an empty slice
	if got, want := removeDuplicates([]string{}), []string{}; !equal(got, want) {
		t.Errorf("removeDuplicates(%q) = %q, want %q", []string{}, got, want)
	}
}

// equal checks if two slices are equal
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
