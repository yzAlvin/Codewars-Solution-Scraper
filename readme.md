# Codewars Solution Scraper

This go program will scrape all your codewars solutions and output them in json file.

The schema of a single kata is 
```javascript
{
	"Kyu":			"string",
	"KataLink":		"string",
	"KataTitle":		"string",
	"LanguagesSolved":	["string"],
	"Solutions":		["string"]
}
```

## How to Run

1. Clone the repo

2. In config.yaml
	- Specify codewars username
	- Specify session_id cookie
	- Specify number of pages to scrape

3. Run `go run codewars_scraper.go`

## Work to be done

- Merge output of all codewars solution pages and save to a single json

- Try to link language with the solution, something like:
```javascript
{
	"Kyu":		"string",
	"KataLink":	"string",
	"KataTitle":	"string",
	"KataSolution":	[{
				"Language": "string",
				"Solutions": ["string"]
			}]
}
```
