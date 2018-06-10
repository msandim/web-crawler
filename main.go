package main

import (
	"flag"
	"fmt"
	"webcrawler/crawler"
)

func parseArguments() (nWorkers int, rateLimit int, domain string) {
	flag.IntVar(&nWorkers, "nworkers", 4, "the number of workers to crawl the domain")
	flag.IntVar(&rateLimit, "ratelimit", 4, "the number of HTTP requests that can be done at the same time")
	flag.StringVar(&domain, "domain", "", "the domain to crawl")
	flag.Parse()

	// TODO: check if the domain is valid, and the ints > 0
	return
}

func main() {
	nWorkers, rateLimit, domain := parseArguments()

	fmt.Println("nworkers: ", nWorkers, " ratelimit: ", rateLimit, " domain: ", domain)
	crawler := crawler.New(nWorkers, rateLimit, domain)
	crawler.Run()
}
