package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"adventOfCode2024/util"
)

type data struct {
  warehouse    [][]rune
  instructions []rune
  robotX       int
  robotY       int
}

func newData() *data {
  return &data{warehouse: [][]rune{}, instructions: []rune{}, robotX: -1, robotY: -1}
}

func (d *data) printWarehouse() {
  for _, di := range d.warehouse {
    fmt.Println(string(di))
  }
}

func (d *data) printInstructions() {
  fmt.Println(string(d.instructions))
}

func (d *data) widenWarehouse() {
  newWarehouse := make([][]rune, len(d.warehouse))
  for i := range d.warehouse {
    row := make([]rune, len(d.warehouse[i])*2)
    for j := range d.warehouse[i] {
      switch d.warehouse[i][j] {
      case rune('#'):
        row[j*2] = rune('#')
        row[j*2+1] = rune('#')
      case rune('O'):
        row[j*2] = rune('[')
        row[j*2+1] = rune(']')
      case rune('.'):
        row[j*2] = rune('.')
        row[j*2+1] = rune('.')
      case rune('@'):
        row[j*2] = rune('@')
        row[j*2+1] = rune('.')
      }
    }
    newWarehouse[i] = row
  }
  d.warehouse = newWarehouse
}

func printData(d *data) {
  fmt.Println("Warehouse map:")
  d.printWarehouse()
  fmt.Println("Instructions:")
  d.printInstructions()
}

func readInput(path string) (*data, error) {
  d := newData()
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return d, err
  }
  scanner := bufio.NewScanner(file)
  // Warehouse
  for scanner.Scan() {
    line := scanner.Text()
    if line == "" {
      break
    }
    d.warehouse = append(d.warehouse, []rune(line))
  }
  // Instructions
  for scanner.Scan() {
    line := scanner.Text()
    d.instructions = slices.Concat(d.instructions, []rune(line))
  }
  return d, nil
}

func findRobot(d *data) (int, int) {
  for i := range d.warehouse {
    for j := range d.warehouse[i] {
      if d.warehouse[i][j] == rune('@') {
        return i, j
      }
    }
  }
  return -1, -1
}

func translateInstruction(instruction rune) (int, int) {
  switch instruction {
  case rune('^'):
    return -1, 0
  case rune('v'):
    return 1, 0
  case rune('>'):
    return 0, 1
  case rune('<'):
    return 0, -1
  default:
    fmt.Printf("Instruction %v is invalid\n", string(instruction))
    return 0, 0
  }
}

func move(d *data, x, y, dx, dy int) bool {
  newX := x + dx
  newY := y + dy
  if d.warehouse[newX][newY] == rune('#') {
    return false
  }
  if d.warehouse[newX][newY] == rune('.') {
    d.warehouse[newX][newY] = d.warehouse[x][y]
    d.warehouse[x][y] = rune('.')
    return true
  }
  if move(d, newX, newY, dx, dy) {
    d.warehouse[newX][newY] = d.warehouse[x][y]
    d.warehouse[x][y] = rune('.')
    return true
  }
  return false
}

func addToBlock(block *[][]int, x, y int) {
  for _, b := range *block {
    if x == b[0] && y == b[0] {
      return
    }
  }
  *block = append(*block, []int{x, y})
}

func checkMoveWideUD(d *data, x, y, dx, dy int, block *[][]int, paired bool) bool {
  addToBlock(block, x, y)
  newX := x + dx
  newY := y + dy
  if d.warehouse[newX][newY] == rune('#') {
    return false
  }
  if d.warehouse[newX][newY] == rune('.') {
    if !paired && d.warehouse[x][y] == rune(']') {
      if !checkMoveWideUD(d, x, y-1, dx, dy, block, true) {
        return false
      }
    } else if !paired && d.warehouse[x][y] == rune('[') {
      if !checkMoveWideUD(d, x, y+1, dx, dy, block, true) {
        return false
      }
    }
    return true
  }
  if checkMoveWideUD(d, newX, newY, dx, dy, block, false) {
    if !paired && d.warehouse[x][y] == rune(']') {
      if !checkMoveWideUD(d, x, y-1, dx, dy, block, true) {
        return false
      }
    } else if !paired && d.warehouse[x][y] == rune('[') {
      if !checkMoveWideUD(d, x, y+1, dx, dy, block, true) {
        return false
      }
    }
    return true
  }
  return false
}

func moveWideUD(d *data, x, y, dx, dy int, paired bool) {
  newX := x + dx
  newY := y + dy
  if d.warehouse[newX][newY] == rune('#') {
    return
  }
  if d.warehouse[newX][newY] == rune('.') {
    if !paired && d.warehouse[x][y] == rune(']') {
      moveWideUD(d, x, y-1, dx, dy, true)
    } else if !paired && d.warehouse[x][y] == rune('[') {
      moveWideUD(d, x, y+1, dx, dy, true)
    }
    d.warehouse[newX][newY] = d.warehouse[x][y]
    d.warehouse[x][y] = rune('.')
    return
  }
  moveWideUD(d, newX, newY, dx, dy, false)
  if !paired && d.warehouse[x][y] == rune(']') {
    moveWideUD(d, x, y-1, dx, dy, true)
  } else if !paired && d.warehouse[x][y] == rune('[') {
    moveWideUD(d, x, y+1, dx, dy, true)
  }
  d.warehouse[newX][newY] = d.warehouse[x][y]
  d.warehouse[x][y] = rune('.')
}

func task1(d *data, debug bool) {
  result := 0
  for _, instruction := range d.instructions {
    robotX, robotY := findRobot(d)
    dx, dy := translateInstruction(instruction)
    move(d, robotX, robotY, dx, dy)
  }
  if debug {
    fmt.Println("Warehouse after all instructions:")
    d.printWarehouse()
  }
  for i := range d.warehouse {
    for j := range d.warehouse[0] {
      if d.warehouse[i][j] != rune('O') {
        continue
      }
      result += i*100 + j
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(d *data, debug bool) {
  result := 0
  d.widenWarehouse()
  if debug {
    fmt.Println("Warehouse after widening")
    d.printWarehouse()
  }
  for _, instruction := range d.instructions {
    robotX, robotY := findRobot(d)
    dx, dy := translateInstruction(instruction)
    switch instruction {
    case rune('<'), rune('>'):
      move(d, robotX, robotY, dx, dy)
    case rune('^'), rune('v'):
      if checkMoveWideUD(d, robotX, robotY, dx, dy, &[][]int{}, false) {
        moveWideUD(d, robotX, robotY, dx, dy, false)
      }
    }
  }
  if debug {
    fmt.Println("Warehouse after all instructions:")
    d.printWarehouse()
  }
  for i := range d.warehouse {
    for j := range d.warehouse[0] {
      if d.warehouse[i][j] != rune('[') {
        continue
      }
      result += i*100 + j
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
