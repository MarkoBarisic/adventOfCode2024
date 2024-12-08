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

func printData(data [][]int) {
  for _, x := range data {
    for _, y := range x {
      fmt.Printf("%d  ", y)
    }
    fmt.Println()
  }
}

func readInput(path string) ([][]int, error) {
  var matrix [][]int
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return matrix, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    var newRow []int
    line := scanner.Text()
    elements := strings.Split(line, " ")
    for _, value := range elements {
      intVal, err := strconv.Atoi(value)
      if err != nil {
        return matrix, err
      }
      newRow = append(newRow, intVal)
    }
    matrix = append(matrix, newRow)
  }
  return matrix, nil
}

func checkSafeCmp(x, y int, trend string, debug bool) bool {
  diff := x-y
  if debug {
    fmt.Printf("Comparing %d <> %d with trend '%s'\n", x, y, trend)
  }
  switch trend {
  case "asc":
    if diff > -1 || diff < -3 {
      return false
    }
    return true
  case "desc":
    if diff < 1 || diff > 3 {
      return false
    }
    return true
  default:
    fmt.Println("Only acceptable trends are 'asc' or 'desc'")
    return false
  }
}

func checkSafeRow(slice []int, debug bool, useDampener bool) (bool ) {
  if debug {
    fmt.Printf("Is %v safe?\n", slice)
  }
  if slice[0] == slice[len(slice)-1] {
    return false 
  }
  trend := "asc"
  if slice[0] > slice[len(slice)-1] {
    trend = "desc"
  }
  for i := range len(slice)-1 {
    if checkSafeCmp(slice[i], slice[i+1], trend, debug) {
      continue
    }
    if useDampener {
      if checkSafeRow(slices.Concat(slice[:i], slice[i+1:]), debug, false) {
        return true
      }
      if checkSafeRow(slices.Concat(slice[:i+1], slice[i+2:]), debug, false) {
        return true
      }
    }
    return false
  }
  return true
}

func task1(data [][]int, debug bool) {
  result := 0
  for _, row := range data {
    isSafe := checkSafeRow(row, debug, false)
    if debug {
      fmt.Printf("--> %v\n", isSafe)
    }
    if isSafe {
      result += 1
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]int, debug bool) {
  result := 0
  for _, row := range data {
    isSafe := checkSafeRow(row, debug, true)
    if debug {
      fmt.Printf("--> %v\n", isSafe)
    }
    if isSafe {
      result += 1
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
