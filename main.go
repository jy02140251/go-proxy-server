package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

type Handler struct {
    store sync.Map
}

func NewHandler() *Handler {
    return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.handleGet(w, r)
    case http.MethodPost:
        h.handlePost(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    if val, ok := h.store.Load(key); ok {
        json.NewEncoder(w).Encode(map[string]interface{}{"value": val})
        return
    }
    http.Error(w, "Not found", http.StatusNotFound)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    for k, v := range data {
        h.store.Store(k, v)
    }
    w.WriteHeader(http.StatusCreated)
}

func main() {
    handler := NewHandler()
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }
    server.ListenAndServe()
}