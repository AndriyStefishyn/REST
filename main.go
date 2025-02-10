package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", getShopHandler).Methods(http.MethodGet)
	r.HandleFunc("/getshop/{id}", getShopDyIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/createshop", createShopHandler).Methods(http.MethodPost)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic("server error")
	}
}

type Shop struct {
	Id          string `json:"id"`
	Version     int    `json:"version"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func (s *Shop) isValid() bool {
	return s.Id != "" && s.Version != 0 && s.Name != "" && s.Location != "" && s.Description != ""
}

// localhost:80/getshop?id=123
func getShopHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("shop.json")
	if err != nil {
		panic(err)
	}

	w.Write(file)
}

func getShopDyIdHandler(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)

	file, err := os.ReadFile("shop.json")
	if err != nil {
		panic(err)
	}

	var shops []Shop
	err = json.Unmarshal(file, &shops)
	if err != nil {
		panic(err)
	}

	shop, err := findShopById(m["id"], shops)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("shop doesnt exist"))
		return
	}

	byteShop, err := json.Marshal(shop)
	if err != nil {
		e := fmt.Errorf("could not serialize")
		panic(e)
	}
	w.Write(byteShop)
}

func findShopById(id string, shops []Shop) (Shop, error) {
	for _, shop := range shops {
		if shop.Id == id {
			return shop, nil
		}
	}
	return Shop{}, fmt.Errorf("shop not found")
}

func createShopHandler(w http.ResponseWriter, r *http.Request) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var newShop Shop
	json.Unmarshal(rBody, &newShop)

	if !newShop.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid data"))
		return
	}

	err = os.WriteFile("newStore.json", rBody, os.ModePerm)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("shop was successfully created"))
}
