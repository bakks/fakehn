package main

import "fmt"
import "os"
import "net/http"
import "exp/html"

func traverse(node *html.Node, ch chan string) {
  if node.Type == html.ElementNode && node.Data == "td" {
    for _, a := range node.Attr {
      if a.Key == "class" && a.Val == "title" && len(node.Child[0].Child) > 0 {
        title := node.Child[0].Child[0].Data

        if title != "More" {
          ch <- title
        }

        return
      }
    }
  }

  for _, child := range node.Child {
    traverse(child, ch)
  }
}

func scrapeurl(url string, channel chan string) {
  r, err := http.Get(url)
  if err != nil {
    fmt.Println("response error: ", err)
    return
  }
  defer r.Body.Close()

  doc, err := html.Parse(r.Body)
  if err != nil {
    fmt.Println("response error: ", err)
    return
  }

  traverse(doc, channel)
}


func write(ch chan string) {
  file, _ := os.Create("titles.txt")
  defer file.Close()

  for {
    title, open := <-ch
    if !open { return }

    file.WriteString(title)
    file.WriteString("\n")
  }
}

func main() {
  channel := make(chan string)

  go write(channel)
  scrapeurl("http://news.ycombinator.com", channel)
  scrapeurl("http://news.ycombinator.com/newest", channel)
  close(channel)
}

