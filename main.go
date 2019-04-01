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

const baseURL = "https://aws.amazon.com/whitepapers/"

func main() {

	os.Mkdir("aws_whitepapers", 0700)

	type MyJSON struct {
		Filename string
		Link     string
	}

	c := colly.NewCollector()

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
		ftype := e.ChildText("a")
		link = e.ChildAttr("a", "href")
		if ftype == "PDF" {
			if !strings.Contains(link, "http") {
				link = "https:" + link
				jsondat = MyJSON{
					Filename: fileName,
					Link:     link,
				}
				encjson, _ := json.Marshal(jsondat)
				fmt.Println(string(encjson))
				d := colly.NewCollector()
				d.OnResponse(func(in *colly.Response) {
					in.Save("aws_whitepapers/" + fileName)
				})
				d.Visit(link)
			}
		}
		return
	})

	c.Visit(baseURL)
}
