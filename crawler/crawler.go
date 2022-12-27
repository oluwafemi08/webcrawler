package crawler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/lawzava/emailscraper"
)

const (
	agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59"
)

var companyUrls []string

// Error check
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func makeRequest(url string) (*http.Response, error) {
	/* TIMEOUT STARTS HERE  */
	// Set a timeout for the web scraper to prevent it from hanging indefinitely
	timeout := time.Duration(15 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	// Make a request to the website using the HTTP client
	resp, err := client.Get(url)
	if err != nil {
		// Handle the error and return if the request fails or times out
		log.Println("Error making request: ", err)
		return nil, err
	}

	// Handle the response from the website
	if resp.StatusCode != 200 {
		// Handle non-200 status codes and return if the request is unsuccessful
		log.Println("Non-200 status code: ", resp.StatusCode)
		return nil, fmt.Errorf("non-200 status code %d", resp.StatusCode)
	}
	return resp, nil

}

// make Request to the websit to check its vaialibaility and status
func CrawlPage(page int) error {
	// Generate the URL of the page to be crawled
	url, err := generatePageURL(page)
	if err != nil {
		return err
	}

	// Make a request to the website using the HTTP client
	resp, err := makeRequest(url)
	if err != nil {
		log.Println("Error making request to the website: ", err)
		return err
	}
	defer resp.Body.Close()

	// Scrape the page
	if err := Crawler(); err != nil {
		return err
	}

	return nil
}

func generatePageURL(page int) (string, error) {
	return fmt.Sprintf("https://www.eworldtrade.com/c/?page=%d", page), nil
}

/*
progress bar function starts here
*/
func showProgress(total int) {
	// Start a goroutine to update the progress bar continuously
	go func() {
		for i := 0; i <= 100; i++ {
			// Calculate the percentage of items processed
			percent := float64(i) / float64(total) * 100
			// Print the progress bar
			fmt.Printf("\r[%-100s] %.2f%%", strings.Repeat("=", i), percent)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println()
	}()
}

/*
progress bar function ends here
*/

// CRAWLER FUNCTION
func Crawler() error {

	// Start the progress bar
	showProgress(100)

	var companyUrls []string

	c := colly.NewCollector(
		colly.UserAgent(agent),

		colly.AllowedDomains("eworldtrade.com", "www.eworldtrade.com"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       9 * time.Second,
		RandomDelay: 4 * time.Second,
	})

	companyProfileCollector := c.Clone()

	c.OnRequest(func(h *colly.Request) {
		fmt.Println("Visiting: ", h.URL)
	})

	c.OnError(func(h *colly.Response, err error) {
		fmt.Println("Request URL:", h.Request.URL, "failed with response:", h, "\nError:", err)

	})

	c.OnHTML("div.buyer-listing-result-row div.com-flex", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")
		// Add a delay of 3 seconds before making the next request

		time.Sleep(3 * time.Second)
		companyProfileCollector.Visit(link)
	})

	companyProfileCollector.OnHTML("div.row > div.col-lg-8", func(h *colly.HTMLElement) {
		// var companyUrl *string
		companyUrl := h.ChildAttr("a", "href")
		fmt.Println("this the company url: ", companyUrl)

		companyUrls = append(companyUrls, companyUrl)
		fmt.Println("these are the company urls: ", companyUrls)

	})

	companyProfileCollector.OnResponse(func(h *colly.Response) {
		fmt.Println("Crawling: ", h.Request.URL.String())
	})

	for page := 1; page <= 1108; page++ {
		crawlUrl, err := generatePageURL(page)
		if err != nil {
			log.Println("Error generating page URL: ", err)
			continue
		}
		fmt.Printf("\r# Scraping page number %.2d  ---> DONE #", page)
		c.Visit(crawlUrl)
	}

	c.Wait()
	companyProfileCollector.Wait()

	return nil

}

func writeUrlToJSON(urls []string) error {
	jsonData, err := json.Marshal(urls)
	if err != nil {
		return fmt.Errorf("Error encoding URLs to JSON: %v", err)
	}

	err = ioutil.WriteFile("links.json", jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON to file: %v", err)
	}

	fmt.Println("URLs saved to links.json")
	return nil
}

func readUrlFromJSON() ([]string, error) {
	jsonData, err := ioutil.ReadFile("links.json")
	if err != nil {
		return nil, fmt.Errorf("Error reading JSON file: %v", err)
	}

	var urls []string
	err = json.Unmarshal(jsonData, &urls)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}
	fmt.Println(urls)
	return urls, nil

}

func scrapeEmails(urls []string) []string {
	var emails []string
	mailer := emailscraper.New(emailscraper.DefaultConfig())

	for _, url := range urls {
		defer func() {
			if reco := recover(); reco != nil {
				log.Println("Error scraping email from URL: ", url)
			}
		}()
		// Use the emailscraper package to scrape email addresses from the URL
		scrapedEmails, err := mailer.Scrape(url)
		checkError(err)

		// Append the scraped email addresses to the emails slice
		emails = append(emails, scrapedEmails...)

	}

	return emails
}

// Write to csv
func writeEmailsToCSV(data []string) error {

	file, err := os.Create("emails.csv")
	checkError(err)

	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"Emails"})

	// Write the email addresses to the CSV file
	for _, email := range data {
		err := writer.Write([]string{email})
		checkError(err)

	}
	//save the changes
	writer.Flush()
	return nil
}

func init() {
	Crawler()
	readUrlFromJSON()
	writeUrlToJSON(companyUrls)
	emails := scrapeEmails(companyUrls)
	err := writeEmailsToCSV(emails)
	if err != nil {
		log.Println("Error writing to CSV file: ", err)
		return
	}

}
