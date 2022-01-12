//+build ignore
package main

import (
	"encoding/base64"
	"html"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

func main() {
	doc, err := goquery.NewDocumentFromReader(os.Stdin)
	if err != nil {
		panic(err)
	}
	m := minify.New()
	m.AddFunc("text/html", mhtml.Minify)
	m.AddFunc("application/javascript", js.Minify)
	m.AddFunc("text/css", css.Minify)
	doc.Find("script").Each(func(i int, s *goquery.Selection) {

		shtml := ""
		if src, ok := s.Attr("src"); ok {
			b, err := os.ReadFile(src)
			if err != nil {
				panic(err)
			}
			shtml = string(b)
			s.RemoveAttr("src")
		} else {
			shtml, _ = s.Html()
			shtml = html.UnescapeString(shtml)
		}
		htmlReader := strings.NewReader(shtml)
		var htmlWriter strings.Builder
		if err := m.Minify("application/javascript", &htmlWriter, htmlReader); err != nil {
			panic(err)
		}
		s.SetHtml(htmlWriter.String())
	})
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		if rel, _ := s.Attr("rel"); rel != "stylesheet" {
			return
		}
		shtml := ""
		src, _ := s.Attr("href")
		b, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}
		shtml = string(b)
		htmlReader := strings.NewReader(shtml)
		var htmlWriter strings.Builder
		if err := m.Minify("text/css", &htmlWriter, htmlReader); err != nil {
			panic(err)
		}
		//io.Copy(&htmlWriter, htmlReader)
		s.AfterHtml("<style>" + htmlWriter.String() + "</style>")
		s.Remove()
	})
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		b, err := os.ReadFile(src)
		if err != nil {
			panic(err)
		}
		str := base64.StdEncoding.EncodeToString(b)
		ext := "XXX"
		if strings.HasSuffix(src, ".gif") {
			ext = "gif"
		}
		s.SetAttr("src", "data:image/"+ext+";base64,"+str)
	})

	newDoc, _ := doc.Html()
	htmlReader := strings.NewReader(newDoc)
	if err := m.Minify("text/html", os.Stdout, htmlReader); err != nil {
		panic(err)
	}
}
