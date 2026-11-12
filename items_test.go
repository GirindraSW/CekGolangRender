package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestCreateItem(t *testing.T) {
 store := &itemStore{items: make(map[string]item)}
 request := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"name":"  Kursi  "}`))
 request.Header.Set("Content-Type", "application/json")
 response := httptest.NewRecorder()
 store.create(response, request)
 if response.Code != http.StatusCreated {
  t.Fatalf("status %d: %s", response.Code, response.Body.String())
 }
 var body struct { Data item `json:"data"` }
 if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil { t.Fatal(err) }
 if body.Data.Name != "Kursi" || body.Data.ID == "" || body.Data.CreatedAt.IsZero() {
  t.Fatalf("invalid created item: %+v", body.Data)
 }
 if stored, ok := store.items[body.Data.ID]; !ok || stored != body.Data {
  t.Fatal("created item was not stored")
 }
}

func TestCreateRejectsInvalidRequests(t *testing.T) {
 cases := []struct { name, method, contentType, body string; status int }{
  {"method", "GET", "application/json", "", 405},
  {"content type", "POST", "text/plain", `{"name":"Kursi"}`, 415},
  {"empty name", "POST", "application/json", `{"name":" "}`, 400},
  {"invalid json", "POST", "application/json", `{`, 400},
  {"extra json", "POST", "application/json", `{"name":"Kursi"} {}`, 400},
  {"unknown field", "POST", "application/json", `{"name":"Kursi","price":10}`, 400},
  {"long name", "POST", "application/json", `{"name":"` + strings.Repeat("a", 101) + `"}`, 400},
  {"large body", "POST", "application/json", `{"name":"` + strings.Repeat("a", 5000) + `"}`, 413},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   store := &itemStore{items: make(map[string]item)}
   request := httptest.NewRequest(tc.method, "/items", strings.NewReader(tc.body))
   request.Header.Set("Content-Type", tc.contentType)
   response := httptest.NewRecorder()
   store.create(response, request)
   if response.Code != tc.status { t.Fatalf("got %d, want %d: %s", response.Code, tc.status, response.Body.String()) }
   if len(store.items) != 0 { t.Fatal("invalid request created an item") }
  })
 }
}
