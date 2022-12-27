package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/oluwafemi08/collyscraper/crawler"
)

func main() {

	startTime := time.Now()
	var wg sync.WaitGroup
	// Create a goroutine for each page
	for i := 1; i <= 1108; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			if err := crawler.CrawlPage(page); err != nil {
				fmt.Println("Error crawling page: ", err)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(startTime)
	fmt.Printf("Done! Code started %v\n Code execution time: %v\n", startTime, elapsed)

}
