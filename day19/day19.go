package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"adventOfCode2024/util"
)

type input struct {
  patterns []string
  designs []string
  validDesigns map[string]bool
  arrangementCount map[string]int
}

func newInput(patterns, designs []string) *input {
  validDesigns := make(map[string]bool)
  arrangementCount := make(map[string]int)
  for _, p := range patterns {
    validDesigns[p] = true
  }
  return &input{patterns: patterns, designs: designs, validDesigns: validDesigns, arrangementCount: arrangementCount}
}

func printData(d *input) {
  fmt.Println("Available patterns:")
  for i := range d.patterns {
    fmt.Printf("%s ", d.patterns[i])
  }
  fmt.Println()
  fmt.Println("Designs:")
  for i := range d.designs {
    fmt.Printf("%s\n", d.designs[i])
  }
}

func readInput(path string) (*input, error) {
  var patterns []string
  var designs []string
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return nil, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      break
    }
    splitLine := strings.Split(line, ", ")
    for _, p := range splitLine {
      patterns = append(patterns, p)
    }
  }
  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      break
    }
    designs = append(designs, line)
  }
  return newInput(patterns, designs), nil
}

func checkDesign(data *input, design string) bool {
  // Check if design is stored as valid
  if valid, contains := data.validDesigns[design]; contains {
    return valid
  }
  // Break the design into smaller chunks
  if len(design) == 1 {
    data.validDesigns[design] = false
    return false
  }
  for _, pattern := range data.patterns {
      if strings.HasPrefix(design, pattern) {
          remainingDesign := design[len(pattern):]
          if checkDesign(data, remainingDesign) {
            data.validDesigns[remainingDesign] = true
            return true
          }
          data.validDesigns[remainingDesign] = false
      }
  }
  return false
}

func countArrangements(data *input, design string) int {
  // Check if arrangements count for design is memorized
  if cnt, contains := data.arrangementCount[design]; contains {
      return cnt
  }
  // Valid
  if design == "" {
      return 1
  }
  cnt := 0
  for _, pattern := range data.patterns {
      if strings.HasPrefix(design, pattern) {
          remainingDesign := design[len(pattern):]
          cnt += countArrangements(data, remainingDesign)
      }
  }
  data.arrangementCount[design] = cnt
  return cnt
}

func task1(data *input, debug bool) {
  result := 0
  for _, d := range data.designs {
    if debug {
      fmt.Printf("Testing design %s\n", string(d))
    }
    if checkDesign(data, d) {
      result++
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data *input, debug bool) {
  result := 0
  for _, d := range data.designs {
    if debug {
      fmt.Printf("Testing design %s\n", string(d))
    }
    result += countArrangements(data, d)
  }
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
