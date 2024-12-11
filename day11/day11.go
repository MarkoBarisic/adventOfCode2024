package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"adventOfCode2024/util"
)

func printData(d []string) {
  for _, s := range d {
    fmt.Printf("%s ", s)
  }
  fmt.Println()
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
    splitLine := strings.Split(line, " ")
    data = slices.Concat(data, splitLine)
  }
  return data, nil
}

func trimNumber(number string) string {
  for i, n := range number {
    if n != rune('0') {
      return strings.Join(strings.Split(number, "")[i:], "")
    }
  }
  return "0"
}

func blink(digit string, blinksLeft int, memo map[string][]int, totalBlinks int) int{
  if blinksLeft == -1 {
    return 1
  }
  _, ok := memo[digit]
  if !ok {
    memo[digit] = make([]int, totalBlinks)
  }
  if memo[digit][blinksLeft] != 0 {
    return memo[digit][blinksLeft]
  }
  result := 0
  if digit == "0" {
    result += blink("1", blinksLeft-1, memo, totalBlinks)
  } else if len(digit) % 2 == 0 {
    splitDigit := strings.Split(digit, "")
    d1 := trimNumber(strings.Join(splitDigit[:len(digit)/2], ""))
    result += blink(d1, blinksLeft-1, memo, totalBlinks)
    d2 := trimNumber(strings.Join(splitDigit[len(digit)/2:], ""))
    result += blink(d2, blinksLeft-1, memo, totalBlinks)
  } else {
  dI, _ := strconv.Atoi(digit)
  result += blink(strconv.Itoa(dI * 2024), blinksLeft-1, memo, totalBlinks)
  }
  memo[digit][blinksLeft] = result
  return result
}

func task1(data []string, debug bool) {
  result := 0
  nBlinks := 25
  memo := make(map[string][]int)
  for _, d := range data {
    dResult := blink(d, nBlinks-1, memo, nBlinks)
    result += dResult
    if debug {
      fmt.Printf("Stone %s produces %d stones in %d blinks\n", d, dResult, nBlinks)
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []string, debug bool) {
  result := 0
  nBlinks := 75
  memo := make(map[string][]int)
  for _, d := range data {
    dResult := blink(d, nBlinks-1, memo, nBlinks)
    result += dResult
    if debug {
      fmt.Printf("Stone %s produces %d stones in %d blinks\n", d, dResult, nBlinks)
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
