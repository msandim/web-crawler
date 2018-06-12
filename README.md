# Web Crawler

Implementation of a web crawler in Go. The program searches for links to internal page of a certain domain (e.g. reddit.com) and crawls them recursively in order to obtain a sitemap containing the links between the pages.

## Compiling 

```go build```

## Run tests

```go test -v ./...```

## Example usage

```./webcrawler -nworkers=40 -ratelimit=40 -domain=https://www.random.org/ > output.txt 2> error.txt```

- **nworkers:** number of workers (go routines) in the pool of workers implemented.
- **ratelimit:** number of workers that can perform an HTTP GET request at the same time.
- **domain:** domain to crawl and obtain the sitemap. URLs without a path will be considered invalid: use `https://www.random.org/` and not `https://www.random.org`.

The program outputs the sitemap to stdout with the following format:
```
. websiteA
  -> websiteB
  -> websiteC
. websiteB
  -> websiteA
. websiteC
```

Errors and warnings are outputed to stderr.