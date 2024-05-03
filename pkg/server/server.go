package server

import (
	"WB0/pkg/memcache"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

var tpl = template.Must(template.ParseFiles("page/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)

}

func showOrder(w http.ResponseWriter, r *http.Request, c *memcache.Cache) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}
	order, found := c.Get(id)
	fmt.Println(order)
	if !found {
		http.Error(w, "Order not found in cache", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func RunServer(c *memcache.Cache) error {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/orderID/", func(w http.ResponseWriter, r *http.Request) {
		showOrder(w, r, c)
	}) // GET
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return err
	}
	select {}
}
