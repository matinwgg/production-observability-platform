package main

import (
  "encoding/json"
  "log"
  "net/http"
  "sync/atomic"
  "time"
)

var requests uint64

func metrics(w http.ResponseWriter, _ *http.Request) {
  w.Header().Set("Content-Type", "text/plain; version=0.0.4")
  w.Write([]byte("http_requests_total " + formatUint(atomic.LoadUint64(&requests)) + "\n"))
}
func formatUint(v uint64) string { if v == 0 { return "0" }; b:=make([]byte,0,20); for v>0 { b=append([]byte{byte('0'+v%10)},b...); v/=10 }; return string(b) }
func handler(w http.ResponseWriter, r *http.Request) { start:=time.Now(); atomic.AddUint64(&requests,1); w.Header().Set("X-Request-Id", time.Now().UTC().Format("20060102T150405.000000000Z")); json.NewEncoder(w).Encode(map[string]any{"path":r.URL.Path,"method":r.Method,"latency_ms":time.Since(start).Microseconds()/1000}) }
func main(){ http.HandleFunc("/metrics",metrics); http.HandleFunc("/",handler); log.Println("observability collector listening on :8080"); log.Fatal(http.ListenAndServe(":8080",nil)) }
