package main

import "fmt"
import "os"
import "time"
import "net/http"
import "exp/html"

func traverse(node *html.Node, ch chan string) string {
  if node.Type == html.ElementNode && node.Data == "td" {
    for _, a := range node.Attr {
      if a.Key == "class" && a.Val == "title" && len(node.Child[0].Child) > 0 {
        title := node.Child[0].Child[0].Data

        if title != "More" {
          ch <- title
        } else {
          for _, a := range node.Child[0].Attr {
            if a.Key == "href" {
              return a.Val
            }
          }
        }

        return ""
      }
    }
  }

  var x string
  for _, child := range node.Child {
    r := traverse(child, ch)
    if r != "" {
      x = r
    }
  }

  return x
}

func scrapeurl(url string, channel chan string) string {
  r, err := http.Get(url)
  if err != nil {
    fmt.Println("response error: ", err)
    return ""
  }
  defer r.Body.Close()

  doc, err := html.Parse(r.Body)
  if err != nil {
    fmt.Println("response error: ", err)
    return ""
  }

  return traverse(doc, channel)
}


func write(ch chan string) {
  t := time.Now()
  filename := fmt.Sprintf("titles/%d%02d%02d.%02d%02d%02d.txt", t.Year(),
    t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

  file, _ := os.Create(filename)
  defer file.Close()

  for {
    title, open := <-ch
    if !open { return }

    file.WriteString(title)
    file.WriteString("\n")
  }
}

var DOMAIN = "http://news.ycombinator.com/"

func scrapepages(next string, channel chan string, n int) {
  for i := 0; i < n; i++ {
    page := scrapeurl(next, channel)
    next = DOMAIN + page
  }
}

func main() {
  channel := make(chan string)

  go write(channel)
  scrapepages(DOMAIN, channel, 5)
  scrapepages(DOMAIN + "newest", channel, 2)
  close(channel)
}

