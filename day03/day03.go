package main

import (
	"bufio"
  "errors"
	"fmt"
	"os"
  "regexp"
	"strconv"
  "time"

  "adventOfCode2024/util"
)

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

func printData(d []string) {
  for _, s := range d {
    fmt.Println(s)
  }
}

func task1(data []string, debug bool) {
  result := 0
  re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
  for _, str := range data {
    matches := re.FindAllStringSubmatch(str, -1)
    for _, match := range matches {
      if debug {
        fmt.Println(match[0])
      }
      m1, _ := strconv.Atoi(match[1])
      m2, _ := strconv.Atoi(match[2])
      result += m1*m2
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []string, debug bool) {
  result := 0
  mulEnabled := true
  re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
  for _, str := range data{
    matches := re.FindAllStringSubmatch(str, -1)
    for _, match := range matches {
      switch match[0] {
      case "do()":
        mulEnabled = true
      case "don't()":
        mulEnabled = false
      default:
        if debug {
          fmt.Printf("Do: %v <> %v\n", mulEnabled, match[0])
        }
        if !mulEnabled {
          continue
        }
        m1, _ := strconv.Atoi(match[1])
        m2, _ := strconv.Atoi(match[2])
        result += m1*m2
      }
    }
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
