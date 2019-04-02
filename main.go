package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// constants definitions
const baseURL = "https://aws.amazon.com/whitepapers/"

func main() {
	// create folder for downloads
	os.Mkdir("aws_whitepapers", 0700)

	type MyJSON struct {
		Filename string
		Link     string
	}

	// instantiate a new collector
	c := colly.NewCollector()

	// every time you hit this HTML element, that contains whitepaper
	c.OnHTML("div[class='aws-text-box section'] div[class='  '] ul li", func(e *colly.HTMLElement) {
		var fileName string
		var link string
		var jsondat MyJSON
		title := strings.Replace(e.ChildText("b"), " ", "_", -1)
		reg, err := regexp.Compile("[^a-zA-Z0-9._]+")
		if err != nil {
			log.Fatal(err)
		}
		fileName = reg.ReplaceAllString(title, "") + ".pdf"

		// some whitepapers have HTML and PDF links
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			// only care about PDF
			if el.Text == "PDF" {
				link = el.Attr("href")
				// some links already have the right prefix, but most
				if !strings.Contains(link, "http") {
					link = "https:" + link
				}
				jsondat = MyJSON{
					Filename: fileName,
					Link:     link,
				}
				encjson, _ := json.Marshal(jsondat)
				fmt.Println(string(encjson))

				// create another collector
				d := colly.NewCollector()
				d.OnResponse(func(in *colly.Response) {
					in.Save("aws_whitepapers/" + fileName)
				})
				// start saving files
				d.Visit(link)
			}
			return
		})

	})
	// visit the baseURL
	c.Visit(baseURL)
}
