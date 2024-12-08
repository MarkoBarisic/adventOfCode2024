package main

import (
	"bufio"
  "errors"
	"fmt"
	"os"
	"time"

  "adventOfCode2024/util"
)

type obstacle struct {
  i int
  j int
  d rune
}

func newObstacle(i, j int, d rune) *obstacle {
  return &obstacle{i: i, j: j, d: d}
}

func printData(d [][]rune) {
  for _, s := range d {
    fmt.Println(string(s))
  }
}

func readInput(path string) ([][]rune, error) {
  var data [][]rune
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    data = append(data, []rune(line))
  }
  return data, nil
}

func initGuard(data [][]rune) (int, int, rune){
  for i, s := range data {
    for j, c := range s {
      switch c {
      case rune('^'), rune('v'), rune('>'), rune('<'):
        return i, j, c
      default:
        continue
      }
    }
  }
  fmt.Println("Couldn't find the guard!")
  return 0, 0, 0
}

func nextDirection(direction rune) rune {
  switch direction {
  case rune('^'):
    return rune('>')
  case rune('>'):
    return rune('v')
  case rune('v'):
    return rune('<')
  case rune('<'):
    return rune('^')
  default:
    fmt.Printf("Wrong direction %v\n", string(direction))
    return 0
  }
}

func nextGuard(i, j int, direction rune) (int, int) {
  switch direction {
  case rune('^'):
    return i-1, j
  case rune('v'):
    return i+1, j
  case rune('>'):
    return i, j+1
  case rune('<'):
    return i, j-1
  default:
    fmt.Printf("Wrong direction %v\n", string(direction))
    return i, j
  }
}

func countX(data [][]rune) int {
  cnt := 0
  for _, s := range data {
    for _, c := range s {
      switch c {
      case rune('X'):
        cnt += 1
      default:
        continue
      }
    }
  }
  return cnt
}

func locateAllX(data [][]rune) [][]int {
  var xLocations [][]int
  for i, s := range data {
    for j, c := range s {
      switch c {
      case rune('X'):
        xLocations = append(xLocations, []int{i, j})
      default:
        continue
      }
    }
  }
  return xLocations
}

func task1(data [][]rune, debug bool) {
  iGurad, jGuard, dGuard := initGuard(data)
  i := 0
  j := 0
  for i >= 0 && i < len(data) && j >= 0 && j < len(data[0]) {
    iNext, jNext := nextGuard(iGurad, jGuard, dGuard)
    if iNext < 0 || iNext >= len(data) || jNext < 0 || jNext >= len(data[0]) {
      data[iGurad][jGuard] = rune('X')
      if debug {
        printData(data)
      }
      break
    }
    for x := 1; x < 4; x++ {
      if data[iNext][jNext] != '#' {
        break
      }
      dGuard = nextDirection(dGuard)
      iNext, jNext = nextGuard(iGurad, jGuard, dGuard)
    }
    data[iNext][jNext] = dGuard
    data[iGurad][jGuard] = rune('X')
    iGurad = iNext
    jGuard = jNext
    if debug {
      fmt.Printf("Guard is going to: (%d,%d) going %s\n", iNext, jNext, string(dGuard))
    }
  }
  result := countX(data)
  fmt.Printf("Task 1: %d\n", result)
}

func containsObstacle(obstacleList []*obstacle, target *obstacle) bool {
  for _, o := range obstacleList {
    if o.d == target.d && o.i == target.i && o.j == target.j {
      return true
    }
  }
  return false
}

func task2Traverse(data [][]rune, iGurad, jGuard int, dGuard rune, debug bool) bool {
  var obstacleList []*obstacle
  i := 0
  j := 0
  for i >= 0 && i < len(data) && j >= 0 && j < len(data[0]) {
    iNext, jNext := nextGuard(iGurad, jGuard, dGuard)
    if iNext < 0 || iNext >= len(data) || jNext < 0 || jNext >= len(data[0]) {
      if debug {
        printData(data)
      }
      break
    }
    for x := 1; x <= 2; x++ {
      if data[iNext][jNext] != '#' {
        break
      }
      newObs := newObstacle(iNext, jNext, dGuard)
      if containsObstacle(obstacleList, newObs) {
        return false
      }
      obstacleList = append(obstacleList, newObs)
      dGuard = nextDirection(dGuard)
      iNext, jNext = nextGuard(iGurad, jGuard, dGuard)
    }
    data[iNext][jNext] = dGuard
    data[iGurad][jGuard] = rune('X')
    iGurad = iNext
    jGuard = jNext
    if debug {
      fmt.Printf("Guard is going to: (%d,%d) going %s\n", iNext, jNext, string(dGuard))
    }
  }
  return true
}

func task2(data [][]rune, path string, debug bool) {
  iGurad, jGuard, dGuard := initGuard(data)
  if debug {
    fmt.Printf("Guard is located at: (%d,%d) going %s\n", iGurad, jGuard, string(dGuard))
  }
  task1(data, debug)
  xLocations := locateAllX(data)
  result := 0
  for _, xLoc := range xLocations {
    cleanData, _ := readInput(path)
    cleanData[xLoc[0]][xLoc[1]] = rune('#')
    if !task2Traverse(cleanData, iGurad, jGuard, dGuard, debug) {
      if debug {
        fmt.Println("Solution:")
        cleanData[xLoc[0]][xLoc[1]] = rune('O')
        printData(cleanData)
      }
      result += 1
    }
    cleanData[xLoc[0]][xLoc[1]] = rune('X')
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
    task2(data, path, debug)
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
