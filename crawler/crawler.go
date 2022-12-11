package crawler

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/91.0.4472.114 Safari/537.36"
)

var linkArr []string

func Crawler() {
	c := colly.NewCollector(
		colly.UserAgent(agent),

		colly.AllowedDomains("eworldtrade.com", "www.eworldtrade.com"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		Delay:       2 * time.Second,
		RandomDelay: 3 * time.Second,
	})

	companyProfileCollector := c.Clone()

	c.OnRequest(func(h *colly.Request) {
		fmt.Println("Visiting: ", h.URL)
	})

	c.OnError(func(h *colly.Response, err error) {
		fmt.Printf("Something went wrong: %s: %v ", h.Request.URL, err)
	})

	c.OnHTML("div.buyer-listing-result-row  div.com-flex", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")
		companyProfileCollector.Visit(link)
	})

	companyProfileCollector.OnHTML(`div.row > div.col-lg-8`, func(h *colly.HTMLElement) {
		companyUrl := h.ChildAttr("a", "href")

		linkArr = append(linkArr, companyUrl)

		fmt.Println("This are company Url: ", linkArr)
	})

	companyProfileCollector.OnResponse(func(h *colly.Response) {
		fmt.Println("Crawling: ", h.Request.URL.String())
	})

	companyProfileCollector.OnScraped(func(h *colly.Response) {
		fmt.Println("FINISHED: Company Profile Collector", h.Request.URL.String())
	})

	for page := 1; page <= 1108; page++ {
		crawlUrl := fmt.Sprintf("https://www.eworldtrade.com/c/?page=%d", page)
		fmt.Printf("\r# Scraping page number %.2d  ---> DONE #", page)
		c.Visit(crawlUrl)
	}

	c.Wait()

}

func scrapeEmails(urls []string) []string {
	var emails []string

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		webpage, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		emailRegexp := regexp.MustCompile("[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+")
		emailsFromPage := emailRegexp.FindAllString(string(webpage), -1)

		for _, email := range emailsFromPage {
			if !strings.Contains(email, "address") && !strings.Contains(email, "example") {
				emails = append(emails, email)
			}
		}
	}
	// var emails := scrapeEmails(linkArr)
	return emails
}

// Error check
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//Write to csv

func writeToCSV(data []string) {
	f, err := os.Create("emails.csv")
	checkError(err)
	
	defer f.Close()

	w := csv.NewWriter(f)
	w.Write([]string{"Emails"})

	for _, email := range data {
		err := w.Write([]string{email})
		checkError(err)
		
	}

	w.Flush()
}

func init() {
	Crawler()
	emails := scrapeEmails(linkArr)
	writeToCSV(emails)
}

/*

// removeDuplicates removes duplicate emails
func removeDuplicates(emails []string) []string {
	var seen = make(map[string]bool)
	var newEmails []string

	for _, v := range emails {
		if _, ok := seen[v]; !ok {
			newEmails = append(newEmails, v)
			seen[v] = true
		}
	}
	return newEmails
}


//Write to csv

func writeToCSV(data []string) {
    f, err := os.Create("emails.csv")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    w := csv.NewWriter(f)
    w.Write([]string{"Emails"})

    for _, email := range data {
        err := w.Write([]string{email})
        if err != nil {
            panic(err)
        }
    }

    w.Flush()
}

func main() {
    Crawler()
    emails := scrapeEmails(linkArr)
    writeToCSV(emails)
}

*/

/*


func scrapeEmails(urls []string) []string {
		var emails []string
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}

			webpage, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			r := regexp.MustCompile("[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+")
			e := r.FindAllString(string(webpage), -1)
			for _, email := range e {
				if !strings.Contains(email, "address") && !strings.Contains(email, "example") {
					emails = append(emails, email)
				}
			}
		}
		emails := scrapeEmails(linkArr)
	}

    func checkError(err error) {
        if err != nil {
            panic(err)
        }
    }
*/
