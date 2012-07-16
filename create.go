package main

import "fmt"
import "os"
import "io/ioutil"
import "bufio"
import "strings"
import "time"
import "math/rand"
import "strconv"

type Word struct {
  str string
  next []*Word
}

var Wordmap map[string]*Word
var Titles map[string]bool = make(map[string]bool)

func getword(str string) *Word {
  x, exists := Wordmap[str]

  if !exists {
    x = &Word{str, make([]*Word, 0, 256)}
    Wordmap[str] = x
  }

  return x
}

var prepositions = []string {
  "a", "an", "as", "at", "but", "by", "for", "from", "in", "into", "of",
  "on", "onto", "over", "per", "since", "than", "the", "till", "times", "to",
  "unto", "up", "upon", "with", "within", "without" }

func cleanword(s string) string {
  x := strings.ToLower(s)

  for _,prep := range prepositions {
    if x == prep {
      return x
    }
  }
  return s
}

func addwords(one, two string) {
  word1 := getword(cleanword(one))
  word2 := getword(cleanword(two))

  word1.next = append(word1.next, word2)
}

func readtitlesfile(filename string) {
  file, err := os.Open(filename)
  if err != nil {
    fmt.Println("could not open titles file")
    return
  }

  reader := bufio.NewReader(file)

  for {
    bytes,_,err := reader.ReadLine()
    if err != nil {
      break
    }

    str := string(bytes)
    Titles[strings.ToLower(str)] = true
    strs := strings.Split(string(bytes), " ")

    addwords("", strs[0])
    addwords(strs[len(strs)-1], "")

    for i := 0; i < len(strs) - 1; i++ {
      addwords(strs[i], strs[i+1])
    }
  }
}

func build() {
  Wordmap = make(map[string]*Word)

  fileinfos, err := ioutil.ReadDir("titles")

  if err != nil {
    fmt.Println("Could not open titles directory")
    return
  }

//  for _,fi := range fileinfos {
//    readtitlesfile("titles/" + fi.Name())
//  }

  readtitlesfile("titles/" + fileinfos[len(fileinfos) - 1].Name())
}

func createword(random *rand.Rand) string {
  x := getword("")
  s := ""

  for {
    x = x.next[random.Int31n(int32(len(x.next)))]

    if x.str == "" {
      return s
    }

    s += x.str + " "
  }

  return ""
}

func filter(s string) bool {
  if len(s) < 25 || len(s) > 130 {
    return false
  }

  if _,ok := Titles[strings.Trim(strings.ToLower(s), " ")]; ok {
    return false
  }

  return true
}

func create(n int) {
  var random *rand.Rand = rand.New(rand.NewSource(time.Now().Unix()))
  var s string

  for i := 0; i < n; i++ {
    s = createword(random)

    if filter(s) {
      fmt.Println(s, "\n")
    } else {
      i--
    }
  }
}

func main() {
  build()
  n := 30

  if len(os.Args) > 1 {
    n,_ = strconv.Atoi(os.Args[1])
  }

  fmt.Println()
  create(n)
}

