package main

import (
	"flag"
	"fmt"
	"webcrawler/crawler"
)

func parseArguments() (nWorkers int, domain string, maxDepth int) {
	flag.IntVar(&nWorkers, "nworkers", 4, "the number of workers to crawl the domain")
	flag.StringVar(&domain, "domain", "", "the domain to crawl")
	flag.IntVar(&maxDepth, "depth", 5, "the maximum depth of domain crawling")
	flag.Parse()

	// TODO: check if the domain is valid, and the ints > 0
	return
}

func main() {
	nWorkers, domain, maxDepth := parseArguments()

	fmt.Println("nworkers: ", nWorkers, " domain: ", domain, " maxDepth: ", maxDepth)
	crawler.New(nWorkers, domain, maxDepth)
}
