# ratelimit-go

```go
// Initialize

// Eg. Consume API maximum 10 times in 15 seconds, reset usage after every 15 seconds
var rt = ratelimit.NewWindow(10, time.Second*15)

// use it in middleware
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ip := getIP(req)

        // use unique identifier, Consume API.
		_, canUseAPI := rt.Use(ip)
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

> Note: Implementing redis integration soon.
