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
)

// Set default values
var (
    MaxSize = flag.Int("max-size", 1024, "maximum key size")
    srvPort = flag.Int("srv-port", 1123, "server listening port")
)

func Run() {
    flag.Parse()

    // handle requests to /key
    http.HandleFunc("/key/", KeyHandler)

    // Start server
    addr := fmt.Sprintf(":%d", *srvPort)
    log.Printf("Starting server on %s (max-size=%d)\n", addr, *MaxSize)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

func KeyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Extract the length component
    lengthStr := path.Base(r.URL.Path)
    length, err := strconv.Atoi(lengthStr)
    if err != nil || length < 1 {
        http.Error(w, "Invalid key length", http.StatusBadRequest)
        return
    }

    if length > *MaxSize {
        http.Error(w, fmt.Sprintf("Requested length exceeds max-size (%d)", *MaxSize), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", strconv.Itoa(length))

    // Generate and write random bytes
    buf := make([]byte, length)
    if _, err := io.ReadFull(rand.Reader, buf); err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if _, err := w.Write(buf); err != nil {
        log.Printf("Failed writing response: %v", err)
    }
}
