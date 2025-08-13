package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"sisupass.com/sisupass/cmd/api/app"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      *sync.RWMutex
	r       rate.Limit
	b       int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		mu:      &sync.RWMutex{},
		r:       r,
		b:       b,
	}
}

func (rl *RateLimiter) cleanupOldClients() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, client := range rl.clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getClient(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, found := rl.clients[ip]; !found {
		rl.clients[ip] = &client{
			limiter: rate.NewLimiter(rl.r, rl.b),
		}

		// start cleanup goroutine for first client
		if len(rl.clients) == 1 {
			go rl.cleanupOldClients()
		}
	}

	rl.clients[ip].lastSeen = time.Now()
	return rl.clients[ip].limiter
}

func RateLimit(app *app.Application, rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !app.Config.Limiter.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.ServerErrorResponse(w, r, err)
			}

			limiter := rateLimiter.getClient(ip)
			if !limiter.Allow() {
				app.RateLimitExceededResponse(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
