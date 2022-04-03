# NS MapCache

mapCache is in-memory LRU key-value store/cache with concurrency support. 

The key type restriction is `comparable`. The value type restriction is `any`.

Write asymptotic is O(1), Read asymptotic is O(1) amortized. 

No additional goroutines to clean up cache by TTL used

### Usage

#### Using TTL
```go
// create new instance
mc := NewMapCache[int64, uint](10_000, time.Hour)
// adding new element to cache
mc.Set(12, 123456)
// reading from cache
value, found := mc.Get(12)
// use result
fmt.Println(value, found)
```

#### Cache with no expiration

```go
// create new instance
mc := NewMapCache[int64, uint](10_000, 0)
// adding new element to cache
mc.Set(12, 123456)
// reading from cache
value, found := mc.Get(12)
// use result
fmt.Println(value, found)
```