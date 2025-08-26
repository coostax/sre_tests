package keyserver

import (
    "crypto/rand"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "path"
    "strconv"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Set default values
var (
    MaxSize = flag.Int("max-size", 1024, "maximum key size (default 1024)")
    srvPort = flag.Int("srv-port", 1123, "server listening port (default 1123)")
    //---- Prometheus metrics
    // **key_length_histogram**: 20 linear buckets from 0 to maxSize.
    keyLengthHistogram = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "keyserver_key_length_bytes",
            Help:    "Histogram of key lengths generated.",
            Buckets: prometheus.LinearBuckets(0, float64(*MaxSize/20), 20),
        },
    )
    // **http_status_counter**: count HTTP responses by status code.
    httpStatusCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "keyserver_http_status_total",
            Help: "Total number of HTTP responses, labeled by status code.",
        },
        []string{"status"},
    )
)

func Run() {
    flag.Parse()

    // handle requests to /key
    http.HandleFunc("/key/", KeyHandler)

    // handle requests to /metrics (prometheus)
    prometheus.MustRegister(keyLengthHistogram, httpStatusCounter)
    http.Handle("/metrics", promhttp.Handler())

    // Start server
    addr := fmt.Sprintf(":%d", *srvPort)
    log.Printf("Starting server on %s (max-size=%d)\n", addr, *MaxSize)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

// Send a non 200 code and increments counter
func badRequest(w http.ResponseWriter, status string, code int) {
    httpStatusCounter.WithLabelValues(strconv.Itoa(code)).Inc()
    http.Error(w, status, code)
}

func KeyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        badRequest(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract the length component
    lengthStr := path.Base(r.URL.Path)
    length, err := strconv.Atoi(lengthStr)
    if err != nil || length < 1 {
        badRequest(w, "Invalid key length", http.StatusBadRequest)
        return
    }

    if length > *MaxSize {
        badRequest(w, fmt.Sprintf("Requested length exceeds max-size (%d)", *MaxSize), http.StatusBadRequest)
        return
    }
    // Observe the requested length (only if valid)
    keyLengthHistogram.Observe(float64(length))
    //set headers
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", strconv.Itoa(length))

    // Generate and write random bytes
    buf := make([]byte, length)
    if _, err := io.ReadFull(rand.Reader, buf); err != nil {
        badRequest(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if _, err := w.Write(buf); err != nil {
        log.Printf("Failed writing response: %v", err)
    }
    //record 200 OK
    httpStatusCounter.WithLabelValues("200").Inc()
}
