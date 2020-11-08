package generator

// typedCacheTemplate holds the base template for a typed cache
const typedCacheTemplate = `
// Code generated by github.com/Hartimer/loadingcache/cmd/typedcache, DO NOT EDIT.
package {{.Package}}

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
)

type {{.Name}} interface {
	Get(key {{.KeyType}}) ({{.ValueType}}, error)
	Put(key {{.KeyType}}, value {{.ValueType}})
	Invalidate(key {{.KeyType}}, keys ...{{.KeyType}})
	InvalidateAll()
}

type {{.Name}}Option func({{.Name}})

type LoadFunc func({{.KeyType}}) ({{.ValueType}}, error)

type RemovalNotification struct {
	Key    {{.KeyType}}
	Value  {{.ValueType}}
	Reason loadingcache.RemovalReason
}

type RemovalListenerFunc func(RemovalNotification)

type internalImplementation struct {
	genericCache loadingcache.Cache
	cacheOptions []loadingcache.CacheOption
}

func Clock(clk clock.Clock) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.Clock(clk))
		}
	}
}

func ExpireAfterWrite(duration time.Duration) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.ExpireAfterWrite(duration))
		}
	}
}

func ExpireAfterRead(duration time.Duration) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.ExpireAfterRead(duration))
		}
	}
}

func Load(f LoadFunc) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.Load(func(key interface{}) (interface{}, error) {
				typedKey, ok := key.({{.KeyType}})
				if !ok {
					return 0, fmt.Errorf("Key expeceted to be a {{.KeyType}} but got %T", key)
				}
				return f(typedKey)
			}))
		}
	}
}

func MaxSize(maxSize int32) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.MaxSize(maxSize))
		}
	}
}

func RemovalListener(listener RemovalListenerFunc) {{.Name}}Option {
	return func(cache {{.Name}}) {
		if g, ok := cache.(*internalImplementation); ok {
			g.cacheOptions = append(g.cacheOptions, loadingcache.RemovalListener(func(notification loadingcache.RemovalNotification) {
				typedNofication := RemovalNotification{Reason: notification.Reason}
				var ok bool
				typedNofication.Key, ok = notification.Key.({{.KeyType}})
				if !ok {
					panic(fmt.Sprintf("Somehow the key is a %T instead of a {{.KeyType}}", notification.Key))
				}
				typedNofication.Value, ok = notification.Value.({{.ValueType}})
				if !ok {
					panic(fmt.Sprintf("Somehow the value is a %T instead of an {{.ValueType}}", notification.Value))
				}
				listener(typedNofication)
			}))
		}
	}
}

func NewCache(options ...{{.Name}}Option) {{.Name}} {
	internal := &internalImplementation{}
	for _, option := range options {
		option(internal)
	}

	internal.genericCache = loadingcache.NewGenericCache(internal.cacheOptions...)
	return internal
}

func (i *internalImplementation) Get(key {{.KeyType}}) ({{.ValueType}}, error) {
	val, err := i.genericCache.Get(key)
	if err != nil {
		return 0, err
	}
	typedVal, ok := val.({{.ValueType}})
	if !ok {
		// TODO type mismatch error
	}
	return typedVal, nil
}

func (i *internalImplementation) Put(key {{.KeyType}}, value {{.ValueType}}) {
	i.genericCache.Put(key, value)
}

func (i *internalImplementation) Invalidate(key {{.KeyType}}, keys ...{{.KeyType}}) {
	genericKeys := make([]interface{}, len(keys))
	for i, k := range keys {
		genericKeys[i] = k
	}
	i.genericCache.Invalidate(key, genericKeys...)
}

func (i *internalImplementation) InvalidateAll() {
	i.genericCache.InvalidateAll()
}
`