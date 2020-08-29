# ratelimit-go

### Prerequisite

- redis

```
go get github.com/go-redis/redis/v8
go get github.com/Marvin9/ratelimit
```

```go
// Initialize

var redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

// Eg. Consume API maximum 10 times in 15 seconds, reset usage after every 15 seconds
var rt = ratelimit.NewWindow(10, time.Second*15, redisClient)

// use it in middleware
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := getIP(req)

		// use unique identifier, Consume API.
		profiler, canUseAPI := rt.Use(ip)

		log.Prtinf("\nFor %v: \n%v APIs left for usage.\nNext reset: %v", ip, profiler.APIUsageLeft, profiler.NextWindow)

		if !canUseAPI {
		    fmt.Fprintf(w, "API call Limit exceeded. ")
		    return
		}

		next(w, req)
	})
}

func main() {
    http.HandleFunc("/hello", rateLimitMiddleware(hello))

}

```
