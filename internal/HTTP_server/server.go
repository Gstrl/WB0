package HTTP_server

import (
	"WB0/internal/config"
	"WB0/internal/memcache"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func RunServer(c *memcache.Cache, cfg *config.Config) error {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/orderID/", func(w http.ResponseWriter, r *http.Request) {
		showOrder(w, r, c)
	}) // GET
	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	select {}
}
