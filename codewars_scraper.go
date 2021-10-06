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
	Code string `json:"code"`
}

type Kata struct {
	Kyu          string     `json:"kyu"`
	KataLink     string     `json:"kataLink"`
	KataTitle    string     `json:"kata"`
	KataSolution []Solution `json:"solution"`
}

func main() {
	allKatas := make([]Kata, 0)

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
		kataTitle := e.ChildText(".item-title a")
		kataLink := e.ChildAttr("a", "href")
		kyu := e.ChildText(".inner-small-hex")

		allSolutions := make([]Solution, 0)

		solution := Solution{
			Code: code,
		}

		allSolutions = append(allSolutions, solution)

		kata := Kata{
			Kyu:          kyu,
			KataLink:     kataLink,
			KataTitle:    kataTitle,
			KataSolution: allSolutions,
		}

		allKatas = append(allKatas, kata)
	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	collector.Visit("https://www.codewars.com/users/yzAlvin/completed_solutions")

	writeJSON(allKatas)
}

func writeJSON(data []Kata) {
	file, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	ioutil.WriteFile("codewars_solutions.json", file, 0644)
}
