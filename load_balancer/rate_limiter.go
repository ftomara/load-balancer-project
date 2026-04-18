package lb

import (
	"sync"
	"time"
)

type Token struct {
	token_count float64
	last_refill time.Time
}

var ip_map = make(map[string]Token)
var max_token float64 = 10
var ip_mutex sync.Mutex

func checkLimit(ip string) bool {
	ip_mutex.Lock()
	defer ip_mutex.Unlock()
	token := ip_map[ip]
	time_elapsed := time.Since(token.last_refill)
	token.token_count = min(max_token, token.token_count+(time_elapsed.Seconds()*10))
	token.last_refill = time.Now()
	if token.token_count >= 1 {
		token.token_count--
		ip_map[ip] = token
		return true
	}
	ip_map[ip] = token
	return false
}
