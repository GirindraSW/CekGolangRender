package main

import (
 "encoding/json"
 "log"
 "net/http"
 "os"
 "time"
)

func main() {
 port := os.Getenv("PORT")
 if port == "" {
  port = "10000"
 }
 mux := http.NewServeMux()
 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" && r.URL.Path != "/health" {
   http.NotFound(w, r)
   return
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.Header().Set("Cache-Control", "no-store")
  if r.Method != http.MethodGet && r.Method != http.MethodHead {
   w.Header().Set("Allow", "GET, HEAD")
   w.WriteHeader(http.StatusMethodNotAllowed)
   return
  }
  log.Printf("%s %s", r.Method, r.URL.Path)
  if r.Method == http.MethodHead {
   return
  }
  if err := json.NewEncoder(w).Encode(map[string]string{
   "status": "ok",
   "message": "Koneksi ke server Go berhasil!",
   "time": time.Now().UTC().Format(time.RFC3339),
  }); err != nil {
   log.Printf("response error: %v", err)
  }
 })
 server := &http.Server{
  Addr: "0.0.0.0:" + port,
  Handler: mux,
  ReadHeaderTimeout: 5 * time.Second,
  WriteTimeout: 10 * time.Second,
  IdleTimeout: 60 * time.Second,
 }
 log.Printf("Server Go berjalan di %s", server.Addr)
 log.Fatal(server.ListenAndServe())
}
