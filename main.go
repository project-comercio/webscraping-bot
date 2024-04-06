package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"

	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type topic struct {
	Page        string `json:"page"`
	Image       string `json:"image"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func handleGetAllPosts(w http.ResponseWriter, r *http.Request) {
	data := getAllPosts()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func getAllPosts() []topic {
	var items []topic

	callEndeazor := colly.NewCollector(
		colly.AllowedDomains("endeavor.org.br"),
	)

	callEndeazor.OnHTML("li.Card", func(h *colly.HTMLElement) {
		item := topic{
			Page:        "endeavor",
			Image:       h.ChildAttr("img.card-img.lazyloaded", "src"),
			Title:       h.ChildText("h3.card-title"),
			Description: h.ChildText("a.card-categoria"),
			Url:         h.ChildAttr("div.desc a", "href"),
		}

		items = append(items, item)
	})

	callEndeazor.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	err := callEndeazor.Visit("https://endeavor.org.br/")
	if err != nil {
		log.Fatal(err)
	}

	topicsContent, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("data.json", topicsContent, 0644)
	return items
}

func main() {
	router := mux.NewRouter()

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router.Use(corsMiddleware)

	router.HandleFunc("/data", handleGetAllPosts).Methods("GET")

	http.ListenAndServe(":3030", router)

}
