package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type New struct {
	Img 			string  `json:"img"`
	Description 	string  `json:"description"`
	Title 			string	`json:"title"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Wecome the my GO API!")
	})
	router.HandleFunc("/news/{id}", handlescraping).Methods("GET")

	log.Fatal(http.ListenAndServe(":1323", router))
}
func handlescraping(w http.ResponseWriter, r* http.Request){
	vars := mux.Vars(r)
	id:= vars["id"]
	log.Println(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":scraping(id),
	})

}
func scraping(id string) []New {
	news := make([]New,0)
	collector := colly.NewCollector(
		colly.AllowedDomains("pj.gob.pe","www.pj.gob.pe"),
	)
	collector.OnHTML("#Cuerpo",func(element *colly.HTMLElement){
		goquerySelection := element.DOM
		img := fmt.Sprintf(goquerySelection.Find("td").Find("img").Attr("src"))
		descripcion := fmt.Sprintf(goquerySelection.Find("td").Children().Text())
		title := fmt.Sprintf(goquerySelection.Find(".titulo").Children().Text())
		n := New{
			Img:img,
			Description:descripcion,
			Title:title,
		}
		news = append(news,n)
	})
	collector.OnRequest(func(request *colly.Request){
		fmt.Println("Visiting",request.URL.String())
	})
	collector.Visit("http://www.pj.gob.pe/wps/wcm/connect/CorteSuperiorTacnaPJ/s_csj_de_tacna/as_inicio/as_imagen_prensa/as_noticias/?WCM_PI=1&WCM_Page.41e5b68046e9346dbff5ff4c973a96be="+id)
	return news
}