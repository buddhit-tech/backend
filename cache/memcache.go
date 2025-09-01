package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)
func New(addr string) *memcache.Client {
return memcache.New(addr)
}