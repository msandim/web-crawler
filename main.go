package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/msandim/web-crawler/crawler"
)

func parseArguments() (nWorkers int, rateLimit int, timeoutSeconds int, domain string) {
	flag.IntVar(&nWorkers, "nworkers", 4, "the number of workers to crawl the domain")
	flag.IntVar(&rateLimit, "ratelimit", 4, "the number of HTTP requests that can be done at the same time")
	flag.IntVar(&timeoutSeconds, "timeoutseconds", 10, "The number of seconds to wait for a HTTP GET request")
	flag.StringVar(&domain, "domain", "https://www.monzo.com", "the domain to crawl")
	flag.Parse()

	if !isnWorkersValid(nWorkers) {
		fmt.Fprintln(os.Stderr, "main::parseArguments() - Error: Number of workers is invalid: ", nWorkers)
		os.Exit(-1)
	}

	if !isRateLimitValid(rateLimit) {
		fmt.Fprintln(os.Stderr, "main::parseArguments() - Error: Rate limit is invalid: ", rateLimit)
		os.Exit(-1)
	}

	if !isTimeoutSecondsValid(timeoutSeconds) {
		fmt.Fprintln(os.Stderr, "main::parseArguments() - Error: Timeout (seconds) is invalid: ", rateLimit)
		os.Exit(-1)
	}

	if !isDomainValid(domain) {
		fmt.Fprintln(os.Stderr, "main::parseArguments() - Error: Domain is invalid: ", domain)
		os.Exit(-1)
	}
	return
}

func main() {
	nWorkers, rateLimit, timeoutSeconds, domain := parseArguments()

	fmt.Println("nworkers: ", nWorkers, " ratelimit: ", rateLimit, " timeoutseconds: ", timeoutSeconds, " domain: ", domain)
	crawler := crawler.New(nWorkers, rateLimit, timeoutSeconds, domain)
	crawler.Run()
}

func isnWorkersValid(nWorkers int) bool {
	return nWorkers > 0
}

func isRateLimitValid(rateLimit int) bool {
	return rateLimit > 0
}

func isTimeoutSecondsValid(timeoutSeconds int) bool {
	return timeoutSeconds > 0
}

func isDomainValid(domain string) bool {
	domainParsed, err := url.Parse(domain)
	if err != nil {
		return false
	}

	if domainParsed.Scheme == "" || domainParsed.Hostname() == "" {
		return false
	}

	return true
}
