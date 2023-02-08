package middleware

import (
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func Cache() gin.HandlerFunc {
	memoryStore := persist.NewMemoryStore(1 * time.Minute)
	return cache.CacheByRequestURI(memoryStore, 2*time.Second)
}
