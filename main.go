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
	r.HandleFunc("/createshop", createShopHandler).Methods(http.MethodPut)
	r.HandleFunc("/deleteshop/{id}", deleteShopByIdHandler).Methods(http.MethodDelete)

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

func getShops() []Shop {
	file, err := os.ReadFile("shop.json")
	if err != nil {
		panic(err)
	}

	var shops []Shop
	err = json.Unmarshal(file, &shops)
	if err != nil {
		panic(err)
	}

	return shops
}

func deleteShopByIdHandlerV2(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)

	shops := getShops()

	_, err := findShopById(m["id"], shops)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("shop doesnt exist"))
		return
	}

	var index int

	for i := range shops {
		if shops[i].Id == m["id"] {
			index = i

			break
		}
	}

	newShops := append(shops[:index], shops[index+1:]...)

	// save updated shops to  file

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

	shops := getShops()

	for _, shop := range shops {
		if newShop.Id == shop.Id {
			w.Write([]byte("shop already exists"))

			return
		}
	}

	shops = append(shops, newShop)
	w.Write([]byte("shop was successfully created"))

	updatedShops, err := json.MarshalIndent(shops, " ", "")
	if err != nil {
		fmt.Println("Failed to marshal")
	}

	err = os.WriteFile("shop.json", updatedShops, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to write")
	}

}

func deleteShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)

	shops := getShops()
	shop, err := findShopById(m["id"], shops)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("shop doesnt exist"))
		return
	}

	shops = append(shops[:shop])
}
