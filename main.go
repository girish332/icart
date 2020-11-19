package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Product struct {
	Id   string
	Name string
}

type icart interface {
	getCart()
	get()
	post()
	update()
	remove()
}

type cartHandlers struct {
	sync.Mutex
	store map[string]Product
}

func newCartHandlers() *cartHandlers {

	return &cartHandlers{
		store: map[string]Product{
			"Id1": Product{
				Id:   "1",
				Name: "girish",
			},
		},
	}
}

func (h *cartHandlers) get(w http.ResponseWriter, r *http.Request) {
	carts := make([]Product, len(h.store))
	h.Lock()
	i := 0

	for _, prod := range h.store {
		carts[i] = prod
		i++
	}

	h.Unlock()

	jsonBytes, err := json.Marshal(carts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *cartHandlers) getCart(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()

	prod, ok := h.store[parts[2]]

	h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(prod)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (h *cartHandlers) post(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")

	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content type application/json but got '%s'", ct)))
	}

	var prod Product
	err = json.Unmarshal(bodyBytes, &prod)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	prod.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()

	h.store[prod.Id] = prod
	defer h.Unlock()

}

func (h *cartHandlers) remove(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()

	_, ok := h.store[parts[3]]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(h.store, parts[3])

	h.Unlock()

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Delete Successful")))

}

func (h *cartHandlers) update(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// id := parts[3]

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")

	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content type application/json but got '%s'", ct)))
	}

	var prod Product
	err = json.Unmarshal(bodyBytes, &prod)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	h.Lock()

	obj, ok := h.store[parts[3]]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	obj.Id = prod.Id
	obj.Name = prod.Name

	defer h.Unlock()

}

func main() {
	cartHandler := newCartHandlers()
	http.HandleFunc("/cart", cartHandler.get)
	http.HandleFunc("/cart/post", cartHandler.post)
	http.HandleFunc("/cart/", cartHandler.getCart)
	http.HandleFunc("/cart/remove/Id", cartHandler.remove)
	http.HandleFunc("/carts/update/Id", cartHandler.update)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}
