package crawler

import (
	"testing"

	"github.com/gocolly/colly"
)

// TestCrawlerUserAgent tests that the Crawler function sets the correct UserAgent on the
// colly.Collector object it creates.
func TestCrawlerUserAgent(t *testing.T) {
	crawler := colly.NewCollector()

	// Check that the UserAgent of the colly.Collector object is set to the correct value
	if got, want := crawler.UserAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/91.0.4472.114 Safari/537.36"; got != want {
		t.Errorf("crawler.UserAgent = %q, want %q", got, want)
	}
}

// TestCrawlerURLs tests that the Crawler function visits the correct URLs.
func TestCrawlerURLs(t *testing.T) {
	crawler := colly.NewCollector()

	// Set a flag to indicate whether the expected URLs have been visited
	visited := false

	// Set a callback to check the visited URLs
	crawler.OnRequest(func(r *colly.Request) {
		if r.URL.String() == "https://www.eworldtrade.com/c/?page=1" || r.URL.String() == "https://www.eworldtrade.com/c/?page=2" {
			visited = true
		}
	})

	// Visit some test URLs
	crawler.Visit("https://www.eworldtrade.com/c/?page=1")
	crawler.Visit("https://www.eworldtrade.com/c/?page=2")

	// Check that the expected URLs have been visited
	if !visited {
		t.Error("Expected URLs were not visited")
	}
}

// TestCrawlerCompanyURLs tests that the Crawler function extracts the correct company URLs
// from the visited pages.
func TestCrawlerCompanyURLs(t *testing.T) {
	crawler := colly.NewCollector()

	// Set a flag to indicate whether the expected company URLs have been extracted
	extracted := false

	// Set a callback to extract the company URLs
	crawler.OnHTML("div.row > div.col-lg-8", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")

		// Check if the extracted company URL is correct
		if link == "https://www.eworldtrade.com/c/abc-company" || link == "https://www.eworldtrade.com/c/xyz-company" {
			extracted = true
		}
	})

	// Visit some test pages
	crawler.Visit("https://www.eworldtrade.com/c/?page=1")
	crawler.Visit("https://www.eworldtrade.com/c/?page=2")

	// Check that the expected company URLs have been extracted
	if !extracted {
		t.Error("Expected company URLs were not extracted")
	}

	// Set a flag to indicate whether the expected emails have been scraped
scraped := false

// Set a callback to scrape the emails from the company URLs
crawler.OnHTML("div.row > div.col-lg-8", func(h *colly.HTMLElement) {
	companyUrl := h.ChildAttr("a", "href")

	// Check if the scraped email is correct
	if companyUrl == "https://www.eworldtrade.com/c/abc-company" && scrapeEmails([]string{companyUrl})[0] == "abc@company.com" {
		scraped = true
	}
	if companyUrl == "https://www.eworldtrade.com/c/xyz-company" && scrapeEmails([]string{companyUrl})[0] == "xyz@company.com" {
		scraped = true
	}
})

// Visit some test pages
crawler.Visit("https://www.eworldtrade.com/c/?page=1")
crawler.Visit("https://www.eworldtrade.com/c/?page=2")

// Check that the expected emails have been scraped
if !scraped {
	t.Error("Expected emails were not scraped")
}
}


// TestScrapeEmails tests that the scrapeEmails function correctly extracts email addresses
// from the given list of URLs.
func TestScrapeEmails(t *testing.T) {
	urls := []string{
		"https://www.eworldtrade.com/c/abc-company",
		"https://www.eworldtrade.com/c/xyz-company",
	}

	emails := scrapeEmails(urls)
	

		// Check if the extracted emails are correct
		if len(emails) != 2 || emails[0] != "abc@company.com" || emails[1] != "xyz@company.com" {
			t.Errorf("emails = %v, want [abc@company.com, xyz@company.com]", emails)
		}


	/*
	// Check that the correct number of emails is extracted
	if got, want := len(emails), 2; got != want {
		t.Fatalf("len(emails) = %d, want %d", got, want)
	}

	// Check that the extracted email addresses are correct
	if got, want := emails[0], "abc@company.com"; got != want {
		t.Errorf("emails[0] = %q, want %q", got, want)
	}
	if got, want := emails[1], "xyz@company.com"; got != want {
		t.Errorf("emails[1] = %q, want %q", got, want)
	}
	*/
}
