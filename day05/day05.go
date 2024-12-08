package main

import (
	"bufio"
  "errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"strconv"
  "time"

  "adventOfCode2024/util"
)

func printRules(r map[string][]string) {
  for k, v := range r {
    fmt.Printf("Page %v comes before pages %v\n", k, strings.Join(v, ","))
  }
}

func readInput(path string) (map[string][]string, [][]string, error) {
  rules := make(map[string][]string)
  var pages [][]string
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return rules, pages, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      break
    }
    elements := strings.Split(line, "|")
    _, ok := rules[elements[0]]
    if !ok {
      rules[elements[0]] = []string{}
    }
    rules[elements[0]] = append(rules[elements[0]], elements[1])
  }
  for scanner.Scan() {
    line := scanner.Text()
    pages = append(pages, strings.Split(line, ","))
  }
  return rules, pages, nil
}

func isPageValid(rules map[string][]string, page []string) bool {
  for i, p := range page {
    for _, pt := range page[i:] {
      if slices.Contains(rules[pt], p){
        return false
      }
    }
  }
  return true
}

func fixPage(rules map[string][]string, page []string) []string {
  fixedPage := make([]string, len(page))
  for i := 0; i < len(fixedPage); i++ {
    for j := 0; j < len(page); j++ {
      broken := false
      for k := 0; k < len(page); k++ {
        if j == k {
          continue
        }
        if slices.Contains(rules[page[k]], page[j]) {
          broken = true
          break
        }
      }
      if !broken {
        fixedPage[i] = page[j]
        page = slices.Delete(page, j, j+1)
        break
      }
    }
  }
  return fixedPage
}

func task1(rules map[string][]string, pages [][]string, debug bool) {
  var validPages [][]string
  result := 0
  for _, p := range pages {
    if isPageValid(rules, p) {
      validPages = append(validPages, p)
    }
  }
  if debug {
    fmt.Printf("Valid pages: %v\n", validPages)
  }
  for _, p := range validPages {
    middle, _ := strconv.Atoi(p[(len(p)-1)/2])
    result += middle
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(rules map[string][]string, pages [][]string, debug bool) {
  var invalidPages [][]string
  result := 0
  for _, p := range pages {
    if !isPageValid(rules, p) {
      invalidPages = append(invalidPages, p)
    }
  }
  if debug {
    fmt.Printf("Invalid pages: %v\n", invalidPages)
  }

  fixedPages := make([][]string, len(invalidPages))
  for i, ip := range invalidPages {
    fixedPages[i] = fixPage(rules, ip)
  }
  if debug {
    fmt.Printf("Fixed pages: %v\n", fixedPages)
  }
  for _, p := range fixedPages {
    middle, _ := strconv.Atoi(p[(len(p)-1)/2])
    result += middle
  }
  fmt.Printf("Task 2: %d\n", result)
}

func Run(path string, taskId int, debug bool) error {
  rules, pages, err := readInput(path)
  if err != nil {
    return err
  }
  if debug {
    fmt.Printf("\nRunning task %d\n", taskId)
    fmt.Println("Starting rules:")
    printRules(rules)
    fmt.Println("Starting pages:")
    fmt.Println(pages)
  }
  switch taskId {
  case 1:
    task1(rules, pages, debug)
  case 2:
    task2(rules, pages, debug)
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
