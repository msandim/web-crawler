package main

import (
	"flag"
	"fmt"
)

func parseArguments() (domain string, maxDepth int) {
	flag.StringVar(&domain, "domain", "", "the domain to crawl")
	flag.IntVar(&maxDepth, "depth", 5, "the maximum depth of domain crawling")
	flag.Parse()

	// TODO: check if the domain is valid

	return
}

func main() {
	domain, maxDepth := parseArguments()

	fmt.Println("domain: ", domain, " maxDepth: ", maxDepth)

	lol()

	// Agora tenho de analisar
	// get page do domain
	// get urls
	// dec depth

}
