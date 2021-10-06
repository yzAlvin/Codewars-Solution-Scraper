package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gocolly/colly/v2"
)

type Solution struct {
	Kyu      string `json:"kyu"`
	KataLink string `json:"kataLink"`
	Kata     string `json:"kata"`
	Code     string `json:"code"`
}

func main() {
	allSolutions := make([]Solution, 0)

	collector := colly.NewCollector()

	cookie := &http.Cookie{
		Name:  "_session_id",
		Value: os.Getenv("session_secret"),
	}

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, cookie)

	collector.SetCookies("https://www.codewars.com/users/sign_in", cookies)

	collector.OnHTML(".list-item-solutions", func(e *colly.HTMLElement) {
		code := e.ChildText(".mb-5px")
		kata := e.ChildText(".item-title a")
		kataLink := e.ChildAttr("a", "href")
		kyu := e.ChildText(".inner-small-hex")

		solution := Solution{
			Code:     code,
			Kyu:      kyu,
			KataLink: kataLink,
			Kata:     kata,
		}

		allSolutions = append(allSolutions, solution)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://www.codewars.com/users/yzAlvin/completed_solutions")

	writeJSON(allSolutions)
}

func writeJSON(data []Solution) {
	file, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	ioutil.WriteFile("codewars_solutions.json", file, 0644)
}
