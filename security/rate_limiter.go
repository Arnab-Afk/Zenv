package security

import (
	"github.com/gin-gonic/gin"
	"github.com/throttled/throttled/store/memstore"
	"github.com/throttled/throttled/v2"
)

func SetupRateLimiter() gin.HandlerFunc {
	store, _ := memstore.New(65536)
	quota := throttled.RateQuota{MaxRate: throttled.PerMin(10), MaxBurst: 5} // 10 requests per minute
	rateLimiter, _ := throttled.NewGCRARateLimiter(store, quota)
	rateLimiterMiddleware := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	return func(c *gin.Context) {
		rateLimiterMiddleware.RateLimit(c.Writer, c.Request)
		c.Next()
	}
}
