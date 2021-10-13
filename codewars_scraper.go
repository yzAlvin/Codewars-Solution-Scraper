package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
)

type Kata struct {
	Kyu             string   `json:"kyu"`
	KataLink        string   `json:"kataLink"`
	KataTitle       string   `json:"kata"`
	LanguagesSolved []string `json:"languages"`
	Solutions       []string `json:"solutions"`
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found"))
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	username := viper.GetString("username")
	session_id := viper.GetString("session_id")

	allKatas := make([]Kata, 0)

	collector := colly.NewCollector()

	cookie := &http.Cookie{
		Name:  "_session_id",
		Value: session_id,
	}

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, cookie)

	collector.SetCookies("https://www.codewars.com", cookies)

	end := "<div class=\"p-10px js-infinite-marker\" data-page=\"10\"><h5>Loading more items...</h5></div>"

	collector.OnResponse(func(r *colly.Response) {
		if string(r.Body) == end {
			panic("Exit")
		}
	})

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Finished Scraping")
		}
	}()

	collector.OnHTML(".list-item-solutions", func(e *colly.HTMLElement) {
		kataTitle := e.ChildText(".item-title a")
		kataLink := e.ChildAttr("a", "href")
		kyu := e.ChildText(".inner-small-hex")

		allSolutions := make([]string, 0)
		e.ForEach("code", func(_ int, cd *colly.HTMLElement) {
			solution := cd.Text
			if solution == "" {
				return
			}
			allSolutions = append(allSolutions, solution)
		})

		allLanguages := make([]string, 0)
		e.ForEach("h6", func(_ int, l *colly.HTMLElement) {
			language := l.Text
			if language == "" {
				return
			}
			allLanguages = append(allLanguages, language)
		})

		kata := Kata{
			Kyu:             kyu,
			KataLink:        kataLink,
			KataTitle:       kataTitle,
			LanguagesSolved: allLanguages,
			Solutions:       allSolutions,
		}

		allKatas = append(allKatas, kata)
	})

	collector.OnRequest(func(request *colly.Request) {
		request.Headers.Set("x-requested-with", "XMLHttpRequest")
	})

	page := 2
	url := fmt.Sprintf("https://www.codewars.com/users/%s/completed_solutions?page=%d", username, page)
	fmt.Println("Visiting ", url)
	collector.Visit(url)
	writeJSON(allKatas, page)
}

func writeJSON(data []Kata, page int) {
	file, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	filepath := fmt.Sprintf("codewars_solutions_%d.json", page)
	ioutil.WriteFile(filepath, file, 0644)
}
