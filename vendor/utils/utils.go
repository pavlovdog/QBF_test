package utils
import (
	// "encoding/json"
	// "net/http"
    "strings"
    "time"
)

// Format the API provider endpoint with appropriate parameters, such as API key and so on
func GetUrl(basicUrl string, apiKey string, ticker string) string {
	urlWithTicker := strings.Replace(basicUrl, "__SYMBOL__", ticker, 1)
	urlWithKey := strings.Replace(urlWithTicker, "__APIKEY__", apiKey, 1)

	return urlWithKey
}

func EscapeDots(r string) string {
	return strings.Replace(r, ". ", "-", -1)
}

func filterInt(values []int, f func(int) bool) []int {
    result := make([]int, 0)

    for _, v := range values {
        if f(v) {
            result = append(result, v)
        }
    }

    return result
}

type RateLimiter struct {
    RPM         int
    Allow       chan bool
    Pending		[]int
    Executed  	[]int
}

// By calling this function, you'll block the execution until the rate limit
// allows you to perform the call
func (r *RateLimiter) WaitForRateLimitApprove() {
    r.Pending = append(r.Pending, int(time.Now().Unix()))

    // This is the "block point"
    _ = <- r.Allow
}

// Main logic of the rate limiter, should be runned as a routine
func (r *RateLimiter) Initialize(RPM int) {
    r.RPM = RPM
    r.Allow = make(chan bool)

    // Perform the check each second
    // There is more to life than simply increasing its speed d: (c)
    repeater := time.NewTicker(1 * time.Second)

    // Check how much requests has been performed for this period
    // And pass / not pass pending request
    for currentTimestamp := range repeater.C {

        // - Filter out those requests, which were executed more than a minute ago
        minuteAgoTimestamp := int(currentTimestamp.Unix()) - 60

        r.Executed = filterInt(r.Executed, func(v int) bool {
            if v > minuteAgoTimestamp {
                return true
            }

            return false
        })

        // - Rate limit is okey - pass the request execution
        if len(r.Executed) < r.RPM && len(r.Pending) > 0 {
            r.Executed = append(r.Executed, r.Pending[0])

            // - Remove the request from the pending list
            copy(r.Pending, r.Pending[1:])
            r.Pending = r.Pending[:len(r.Pending) - 1]

            r.Allow <- true
        }
    }
}
