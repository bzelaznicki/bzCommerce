package main

import "net/http"

type CategoryPageData struct {
	CategoryName string
	Products     []Product
	CurrentPage  int
	PrevPage     int
	NextPage     int
}

func (cfg *apiConfig) handleCategoryPage(w http.ResponseWriter, r *http.Request) {

}
