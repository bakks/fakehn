package main

import "fmt"
import "os"
import "bufio"
import "strings"
import "math/rand"
import "time"

type Word struct {
  str string
  next []*Word
}

var Wordmap map[string]*Word

func getword(str string) *Word {
  x, exists := Wordmap[str]

  if !exists {
    x = &Word{str, make([]*Word, 0, 128)}
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

func build() {
  Wordmap = make(map[string]*Word)

  file, err := os.Open("titles.txt")
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

    strs := strings.Split(string(bytes), " ")
    addwords("", strs[0])
    addwords(strs[len(strs)-1], "")

    for i := 0; i < len(strs) - 1; i++ {
      addwords(strs[i], strs[i+1])
    }
  }
}

func create() {
  x := getword("")
  random := rand.New(rand.NewSource(time.Now().Unix()))

  for {
    x = x.next[random.Int31n(int32(len(x.next)))]

    if x.str == "" || len(x.next) == 0 {//&& random.Float32() > .8 {
      fmt.Println()
      return
    }

    fmt.Print(x.str, " ")
  }
}

func main() {
  build()
  create()
}

