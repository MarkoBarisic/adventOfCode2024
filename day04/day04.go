package main

import (
	"bufio"
  "errors"
	"fmt"
	"os"
  "strings"
  "slices"
  "regexp"
  "time"

  "adventOfCode2024/util"
)

func printData(d []string) {
  for _, s := range d {
    fmt.Println(s)
  }
}

func readInput(path string) ([]string, error) {
  var data []string
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    data = append(data, line)
  }
  return data, nil
}

func reverseData(d []string) []string {
  rd := make([]string, len(d))
  for i, s := range d {
    var b strings.Builder
    for i := len(s) - 1; i >= 0; i-- {
        b.WriteByte(s[i])
    }
    rd[i] = b.String()
  }
    return rd
}

func transposeData(d []string) []string {
  lineLen := len(d[0])
  td := make([]string, lineLen)
  for i:= 0; i < lineLen; i++ {
    var b strings.Builder
    for _, s := range d {
      b.WriteByte(s[i])
    }
    td[i] = b.String()
  }
  return td
}

func upperTriangle(d []string) []string {
  lineLen := len(d[0])
  ut := make([]string, lineLen)
  // Includes the main diagonal
  for j := 0; j < lineLen; j++ {
    var b strings.Builder
    for i:= 0; i < len(d) && j+i < lineLen; i++ {
      b.WriteByte(d[i][j+i])
    }
    ut[j] = b.String()
  }
  return ut
}

func lowerTriangle(d []string) []string {
  lineLen := len(d[0])
  lt := make([]string, len(d)-1)
  // Doesn't include the main diagonal
  for i := 1; i < len(d); i++ {
    var b strings.Builder
    for j := 0; j < lineLen && i+j < len(d); j++ {
      b.WriteByte(d[i+j][j])
    }
    lt[i-1] = b.String()
  }
  return lt
}

func diagonalData(d []string) []string {
  return slices.Concat(upperTriangle(d), lowerTriangle(d))
}

func checkXmas(data []string) int {
  cnt := 0
  re := regexp.MustCompile(`XMAS`)
  for _, str := range data {
    matches := re.FindAllStringSubmatch(str, -1)
    cnt += len(matches)
  }
  return cnt
}

func buildString(chars...byte) string {
  var b strings.Builder
  for _, c := range chars {
    b.WriteByte(c)
  }
  return b.String()
}

func checkMasX(data []string) int {
  cnt := 0
  for i := 1; i < len(data)-1; i++ {
    for j := 1; j < len(data[i])-1; j++ {
      if data[i][j] != 'A' {
        continue
      }
      d1 := buildString(data[i-1][j-1], data[i][j], data[i+1][j+1])
      d2 := buildString(data[i+1][j-1], data[i][j], data[i-1][j+1])
      if d1 != "MAS" && d1 != "SAM" {
        continue
      }
      if d2 != "MAS" && d2 != "SAM" {
        continue
      }
      cnt += 1
    }
  }
  return cnt
}

func task1(data []string, debug bool) {
  result := 0
  rd := reverseData(data)
  drd := diagonalData(rd)
  rdrd := reverseData(drd)
  td := transposeData(data)
  rtd := reverseData(td)
  dd := diagonalData(data)
  rdd := reverseData(dd)
  result += checkXmas(data)
  result += checkXmas(rd)
  result += checkXmas(drd)
  result += checkXmas(rdrd)
  result += checkXmas(td)
  result += checkXmas(rtd)
  result += checkXmas(dd)
  result += checkXmas(rdd)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []string, debug bool) {
  result := checkMasX(data)
  fmt.Printf("Task 2: %d\n", result)
}

func Run(path string, taskId int, debug bool) error {
  data, err := readInput(path)
  if err != nil {
    return err
  }
  if debug {
    fmt.Printf("\nRunning task %d\n", taskId)
    fmt.Println("Starting data:")
    printData(data)
  }
  switch taskId {
  case 1:
    task1(data, debug)
  case 2:
    task2(data, debug)
  default:
    return errors.New("Invalid value for taskId, please use 1 or 2")
  }
  return nil
}

func main() {
  path, debug, err := util.ProcessArgs(os.Args[1:])
  if err != nil {
    fmt.Printf("%v\nExiting!\n", err)
    os.Exit(1)
  }
  for taskId := 1; taskId <= 2; taskId++ {
    tStart := time.Now()
    err := Run(path, taskId, debug)
    if err != nil {
      fmt.Println(err)
    }
    tEnd := time.Now()
    duration := tEnd.Sub(tStart)
    fmt.Printf("Task %d execution time: %v\n", taskId, duration)
  }
}
