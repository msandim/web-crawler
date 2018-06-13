package crawler

import (
	"testing"
	"webcrawler/fetcher"
)

func TestCrawler1(t *testing.T) {
	nWorkers := 5
	rateLimit := 4
	domain := "http://www.monzo.com/"

	crawler := New(nWorkers, rateLimit, domain)

	if crawler.pool == nil {
		t.Errorf("Pool was not initiliazed")
	}

	if crawler.results == nil || crawler.results != crawler.pool.GetResultsChannel() {
		t.Errorf("Results channel is not initialized")
	}

	if crawler.domain != domain {
		t.Errorf("Domain attribute was not initialized")
	}

	if crawler.checkedUrls == nil {
		t.Errorf("checkedURLs attribute was not initialized")
	}

	if crawler.finishedFlag == nil {
		t.Errorf("Finished flag channel was not initialized")
	}

	if pageFetcher == nil {
		t.Errorf("PageFetcher global variable was not initialized")
	}
	if _, ok := pageFetcher.(*fetcher.HTTPFetcher); !ok {
		t.Errorf("PageFetcher global variable is not set for production")
	}
}

func TestCrawler2(t *testing.T) {
	setUpTest()

	crawler := newTesting(10, "A")
	crawler.Run()

	testLog := log.(*testPrinter)

	if len(testLog.errorMsgs) != 0 {
		t.Errorf("Number of error messages in crawling should be 0.")
	}

	var nA, nB, nC, nD, nE int

	for _, page := range testLog.domainMap {
		switch page.parentURL {
		case "A":
			nA++
			checkMatchingChildren(t, "A", []string{"B", "C"}, page.childrenURLs)
		case "B":
			nB++
			checkMatchingChildren(t, "B", []string{"C", "D"}, page.childrenURLs)
		case "C":
			nC++
			checkMatchingChildren(t, "C", []string{"A", "B", "E", "D"}, page.childrenURLs)
		case "D":
			nD++
			checkMatchingChildren(t, "D", []string{}, page.childrenURLs)
		case "E":
			nE++
			checkMatchingChildren(t, "E", []string{}, page.childrenURLs)
		}
	}

	if nA != 1 {
		t.Errorf("Multiple occurances of A in the domain map")
	}

	if nB != 1 {
		t.Errorf("Multiple occurances of B in the domain map")
	}

	if nC != 1 {
		t.Errorf("Multiple occurances of C in the domain map")
	}

	if nD != 1 {
		t.Errorf("Multiple occurances of D in the domain map")
	}

	if nE != 1 {
		t.Errorf("Multiple occurances of E in the domain map")
	}
}

func checkMatchingChildren(t *testing.T, page string, expectedChildren []string, obtainedChildren []string) {
	if !checkEqualSlices(expectedChildren, obtainedChildren) {
		t.Errorf("Children URLs for %s are not correct. Expected: %v, Obtained: %v",
			page, expectedChildren, obtainedChildren)
	}
}

func checkEqualSlices(s1 []string, s2 []string) bool {
	if s1 == nil && s2 == nil {
		return true
	}

	if s1 == nil || s2 == nil {
		return false
	}

	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func setUpTest() {
	pageFetcher = &TestFetcher{}
	log = &testPrinter{}
}

type TestFetcher struct {
}

func (fetcher *TestFetcher) Fetch(urlArg string) ([]string, []error) {
	switch urlArg {
	case "A":
		return []string{"B", "C"}, nil
	case "B":
		return []string{"C", "D"}, nil
	case "C":
		return []string{"A", "B", "E", "D"}, nil
	default:
		return []string{}, nil
	}
}

type testPrinter struct {
	domainMap []parentPage
	errorMsgs []string
}

type parentPage struct {
	parentURL    string
	childrenURLs []string
}

func (log *testPrinter) logPage(parentURL string, childrenURLs []string) {
	log.domainMap = append(log.domainMap, parentPage{
		parentURL:    parentURL,
		childrenURLs: childrenURLs,
	})
}

func (log *testPrinter) logError(msg string) {
	log.errorMsgs = append(log.errorMsgs, msg)
}
