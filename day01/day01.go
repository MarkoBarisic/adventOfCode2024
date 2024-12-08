package main

import (
	"bufio"
  "errors"
	"fmt"
	"os"
  "strings"
  "strconv"
  "sort"
  "math"
  "time"

  "adventOfCode2024/util"
)

func readInput(path string) ([]int, []int, error) {
  var s1, s2 []int
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return s1, s2, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    elements := strings.Split(line, " ")
    e1, err := strconv.Atoi(elements[0])
    if err != nil {
      return s1, s2, err
    }
    e2, err := strconv.Atoi(elements[len(elements)-1])
    if err != nil {
      return s1, s2, err
    }
    s1 = append(s1, e1)
    s2 = append(s2, e2)
  }
  return s1, s2, nil
}

func absDiff(s1, s2 []int) int {
  sort.SliceStable(s1, func(i, j int) bool {return s1[i] < s1[j]})
  sort.SliceStable(s2, func(i, j int) bool {return s2[i] < s2[j]})
  absDiff := 0
  for i := range len(s1) {
    absDiff += int(math.Abs(float64(s1[i])-float64(s2[i])))
  }
  return absDiff
}

func similarityScore(s1, s2 []int) int {
  similarityScore := 0
  elementSimilarityScore := make(map[int]int)
  for _, e1 := range s1 {
    value, ok := elementSimilarityScore[e1]
    if ok {
      similarityScore += e1*value
      continue
    }
    elementSimilarityScore[e1] = 0
    for _, e2 := range s2 {
      if e2 == e1 {
        elementSimilarityScore[e1] += 1
      }
    }
    similarityScore += e1*elementSimilarityScore[e1]
  }
  return similarityScore
}

func task1(first, second []int) {
  result := absDiff(first, second)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(first, second []int) {
  result := similarityScore(first, second)
  fmt.Printf("Task 2: %d\n", result)
}

func Run(path string, taskId int, debug bool) error {
  first, second, err := readInput(path)
  if err != nil {
    return err
  }
  if debug {
    fmt.Printf("\nRunning task %d\n", taskId)
    fmt.Println("Starting data:")
    fmt.Println("First:")
    fmt.Println(first)
    fmt.Println("Second:")
    fmt.Println(second)
  }
  switch taskId {
  case 1:
    task1(first, second)
  case 2:
    task2(first, second)
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
