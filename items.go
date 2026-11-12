package main

import (
 "crypto/rand"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "log"
 "mime"
 "net/http"
 "strings"
 "sync"
 "time"
)

type item struct {
 ID string `json:"id"`
 Name string `json:"name"`
 CreatedAt time.Time `json:"created_at"`
}

type itemStore struct {
 mu sync.Mutex
 items map[string]item
}

func newItemHandler() http.Handler {
 store := &itemStore{items: make(map[string]item)}
 return http.HandlerFunc(store.create)
}

func (s *itemStore) create(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json; charset=utf-8")
 w.Header().Set("Cache-Control", "no-store")
 if r.Method != http.MethodPost {
  w.Header().Set("Allow", "POST")
  writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Gunakan POST /items"})
  return
 }
 mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
 if err != nil || mediaType != "application/json" {
  writeJSON(w, http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type harus application/json"})
  return
 }
 r.Body = http.MaxBytesReader(w, r.Body, 4096)
 defer r.Body.Close()
 var input struct {
  Name string `json:"name"`
 }
 decoder := json.NewDecoder(r.Body)
 decoder.DisallowUnknownFields()
 if err := decoder.Decode(&input); err != nil {
  invalidBody(w, err)
  return
 }
 var extra any
 if err := decoder.Decode(&extra); err != io.EOF {
  invalidBody(w, err)
  return
 }
 input.Name = strings.TrimSpace(input.Name)
 if input.Name == "" || len([]rune(input.Name)) > 100 {
  writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name wajib diisi, maksimal 100 karakter"})
  return
 }
 var randomID [16]byte
 if _, err := rand.Read(randomID[:]); err != nil {
  writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat ID"})
  return
 }
 created := item{ID: hex.EncodeToString(randomID[:]), Name: input.Name, CreatedAt: time.Now().UTC()}
 s.mu.Lock()
 s.items[created.ID] = created
 s.mu.Unlock()
 writeJSON(w, http.StatusCreated, map[string]any{
  "message": "Item berhasil dibuat (sementara di memori)",
  "data": created,
  "storage": "memory",
 })
}

func invalidBody(w http.ResponseWriter, err error) {
 var tooLarge *http.MaxBytesError
 if errors.As(err, &tooLarge) {
  writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "Body maksimal 4096 byte"})
  return
 }
 writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Body harus satu objek JSON dengan field name"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
 w.WriteHeader(status)
 if err := json.NewEncoder(w).Encode(value); err != nil {
  log.Printf("response error: %v", err)
 }
}
