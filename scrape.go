package main

import "fmt"
import "os"
import "time"
import "net/http"
import "net/url"
import "exp/html"
import "io/ioutil"
import "strings"

type CJ struct {
  cookies map[*url.URL][]*http.Cookie
}

func NewCJ() CJ {
  var x CJ
  x.cookies = make(map[*url.URL][]*http.Cookie)
  return x
}

func (x CJ) SetCookies(u *url.URL, cookies []*http.Cookie) {
  fmt.Println("set cookies")
  x.cookies[u] = cookies
}

func (x CJ) Cookies(u *url.URL) []*http.Cookie {
  fmt.Println("get cookies")
  return x.cookies[u]
}


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
              if a.Val[0] == '/' {
                return a.Val[1:len(a.Val)-1]
              } else {
                return a.Val
              }
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

func scrapeurl(url string, lastUrl string, cookies []*http.Cookie, channel chan string) (string, []*http.Cookie) {
  client := &http.Client{}
  client.Jar = NewCJ()

  fmt.Println(url)
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
  req.Header.Add("Referer", lastUrl)
  req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:14.0) Gecko/20100101 Firefox/14.0.1")
  req.Header.Add("Accept-Language", "en-us,en;q=0.5")
  req.Header.Add("Accept-Encoding", "gzip, deflate")
  req.Header.Add("HTTP-Connection", "keep-alive")

  for i := 0; i < len(cookies); i++ {
    fmt.Println(cookies[i].String())
    req.AddCookie(cookies[i])
  }

  r, err := client.Do(req)
  cookies = r.Cookies()

  if err != nil {
    fmt.Println("response error: ", err)
    return "", nil
  }
  defer r.Body.Close()

  fmt.Println(len(req.Cookies()))

  x,_ := ioutil.ReadAll(r.Body)

  doc, err := html.Parse(strings.NewReader(string(x)))
  if err != nil {
    fmt.Println("response error: ", err)
    return "", nil
  }

  next := traverse(doc,channel)

  if(next == "") {
    fmt.Println(string(x))
  }

  return next, cookies
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
  last := ""
  var cookies []*http.Cookie = nil
  var page string

  for i := 0; i < n; i++ {
    page, cookies = scrapeurl(next, last, cookies, channel)

    last = next
    next = DOMAIN + page
  }
}

func main() {
  channel := make(chan string)

  go write(channel)
  scrapepages(DOMAIN, channel, 4)
  //scrapepages(DOMAIN + "newest", channel, 2)
  close(channel)
}

