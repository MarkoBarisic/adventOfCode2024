package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"adventOfCode2024/util"
)

func printData(d [][]int) {
  for _, x := range d {
    for _, y := range x {
      fmt.Print(y)
    }
    fmt.Println()
  }
}

func readInput(path string) ([][]int, error) {
  var data [][]int
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    lineSplit := strings.Split(line, "")
    lineInt := make([]int, len(lineSplit)) 
    for i, dS := range lineSplit {
      dI,_ := strconv.Atoi(dS)
      lineInt[i] = dI
    }
    data = append(data, lineInt)
  }
  return data, nil
}

func isVisited(vP [][]int, ei, ej int) bool {
  for _, el := range vP {
    if el[0] == ei && el[1] == ej {
      return true
    }
  }
  return false
}

func scoreTrailhead(data [][]int, visitedPeaks *[][]int, i, j, expected int, debug bool) int {
  if i < 0 || j < 0 || i >= len(data) || j >= len(data[0]) {
    return 0
  }
  if data[i][j] != expected {
    return 0
  }
  if data[i][j] == 9 {
    if isVisited(*visitedPeaks, i, j) {
      return 0
    }
    *visitedPeaks = append(*visitedPeaks, []int{i, j})
    return 1
  }
  score := 0
  score += scoreTrailhead(data, visitedPeaks, i+1, j, data[i][j]+1, debug)
  score += scoreTrailhead(data, visitedPeaks, i-1, j, data[i][j]+1, debug)
  score += scoreTrailhead(data, visitedPeaks, i, j+1, data[i][j]+1, debug)
  score += scoreTrailhead(data, visitedPeaks, i, j-1, data[i][j]+1, debug)
  return score
}

func rateTrailhead(data [][]int, i, j, expected int, debug bool) int {
  if i < 0 || j < 0 || i >= len(data) || j >= len(data[0]) {
    return 0
  }
  if data[i][j] != expected {
    return 0
  }
  if data[i][j] == 9 {
    return 1
  }
  score := 0
  score += rateTrailhead(data, i+1, j, data[i][j]+1, debug)
  score += rateTrailhead(data, i-1, j, data[i][j]+1, debug)
  score += rateTrailhead(data, i, j+1, data[i][j]+1, debug)
  score += rateTrailhead(data, i, j-1, data[i][j]+1, debug)
  return score
}

func task1(data [][]int, debug bool) {
  result := 0
  for i := range data {
    for j := range data[i] {
      if data[i][j] != 0 {
        continue
      }
      result += scoreTrailhead(data, &[][]int{}, i, j, 0, debug)
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]int, debug bool) {
  result := 0
  for i := range data {
    for j := range data[i] {
      if data[i][j] != 0 {
        continue
      }
      result += rateTrailhead(data, i, j, 0, debug)
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
