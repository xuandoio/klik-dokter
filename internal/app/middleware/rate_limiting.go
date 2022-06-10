package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimitingMiddleware /**
func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		visitor := limiter.GetVisitor(GetIP(c.Request))
		if !visitor.limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}

var limiter = NewIPRateLimiter(1, 25) // rate is 1, token bucket size is 25

func init() {
	go limiter.CleanUpVisitors()
}

// Visitor struct
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter /**
type IPRateLimiter struct {
	visitors map[string]*Visitor
	mu       *sync.RWMutex
	r        rate.Limit
	b        int
}

// NewIPRateLimiter /**
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		visitors: make(map[string]*Visitor),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}
}

// AddVisitor creates a new rate limiter and adds it to the ips map,
// using the IP address as the key
func (i *IPRateLimiter) AddVisitor(ip string) *Visitor {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.visitors[ip] = &Visitor{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return i.visitors[ip]
}

// GetVisitor /**
// GetLimiter returns the rate limiter for the provided IP address if it exists
// Otherwise call AddVisitor to add IP address to the map
func (i *IPRateLimiter) GetVisitor(ip string) *Visitor {
	i.mu.Lock()
	visitor, ok := i.visitors[ip]

	if !ok {
		i.mu.Unlock()
		return i.AddVisitor(ip)
	}

	i.mu.Unlock()
	return visitor
}

// CleanUpVisitors /**
// Every minute check the map for IPs that haven't been seen for
// more than 3 minutes and delete the entries
func (i *IPRateLimiter) CleanUpVisitors() {
	for {
		time.Sleep(time.Minute)

		i.mu.Lock()

		for ip, v := range i.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(i.visitors, ip)
			}
		}

		i.mu.Unlock()
	}
}

// GetIP /**
func GetIP(r *http.Request) string {
	// get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	// get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	// get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	return ""
}
