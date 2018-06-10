package fetcher

// RateLimiter is a struct that controlls how many concurrent requests can
// be executed in a given context, by calling the function Limit() and Free()
// when the request starts and ends.
type RateLimiter struct {
	semaphore chan bool
}

// NewRateLimiter generates a RateLimiter with a given capacity.
func NewRateLimiter(capacity int) *RateLimiter {
	semaphore := make(chan bool, capacity)

	// Fill channel
	for i := 0; i < capacity; i++ {
		semaphore <- true
	}

	return &RateLimiter{semaphore}
}

// Limit limits the number of concurrent requests by 1 and blocks
// if the number of concurrent requests reached a maximum.
func (rater *RateLimiter) Limit() {
	<-rater.semaphore
}

// Free increases the number of concurrent requests by 1
// This function must be called after a Limit call or it will block.
func (rater *RateLimiter) Free() {
	rater.semaphore <- true
}
