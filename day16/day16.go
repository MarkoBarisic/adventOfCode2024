package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"time"

	"adventOfCode2024/util"
)

type node struct {
  x         int
  y         int
  direction rune
  cost      int
}

func newNode(x, y int, direction rune, cost int) *node {
  return &node{x: x, y: y, direction: direction, cost: cost}
}

func (n *node) String() string {
  return fmt.Sprintf("(%d,%d,%s)", n.x, n.y, string(n.direction))
}

type priorityQueue struct {
  data []*node
}

func newPriorityQueue() *priorityQueue{
  return &priorityQueue{data: []*node{}}
}

func (pq *priorityQueue) Append(n *node) {
  pq.data = append(pq.data, n)
  pq.Sort()
}

func (pq *priorityQueue) Pop() *node {
  outNode := pq.data[0]
  pq.data = pq.data[1:]
  return outNode
}

func (pq *priorityQueue) Sort() {
  slices.SortFunc(pq.data, func(n1, n2 *node) int {
    if n1.cost == n2.cost {
      return 0
    }
    if n1.cost < n2.cost {
      return -1
    }
    return 1
  })
}

func printData(d [][]rune) {
  for i := range d {
    for j := range d[i] {
      fmt.Print(string(d[i][j]))
    }
    fmt.Println()
  }
}

func readInput(path string) ([][]rune, error) {
  var data [][]rune
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return nil, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    data = append(data, []rune(line))
  }
  return data, nil
}

func findStart(data [][]rune) (int, int) {
  for i := range data {
    for j := range data[i] {
      if data[i][j] == rune('S') {
        return i, j
      }
    }
  }
  return -1, -1
}

func findEnd(data [][]rune) (int, int) {
  for i := range data {
    for j := range data[i] {
      if data[i][j] == rune('E') {
        return i, j
      }
    }
  }
  return -1, -1
}

func findShortest(data [][]rune, startX, startY, endX, endY int, startDir rune) (int, map[rune][][]int) {
  pq := newPriorityQueue()
  pq.Append(newNode(startX, startY, startDir, 0))
  costMatrix := make(map[rune][][]int)
  costMatrix[rune('^')] = make([][]int, len(data))
  costMatrix[rune('v')] = make([][]int, len(data))
  costMatrix[rune('>')] = make([][]int, len(data))
  costMatrix[rune('<')] = make([][]int, len(data))
  for i := range data {
    costMatrix[rune('^')][i] = make([]int, len(data[i]))
    costMatrix[rune('v')][i] = make([]int, len(data[i]))
    costMatrix[rune('>')][i] = make([]int, len(data[i]))
    costMatrix[rune('<')][i] = make([]int, len(data[i]))
    for j := range data[i] {
      costMatrix[rune('^')][i][j] = math.MaxInt32
      costMatrix[rune('v')][i][j] = math.MaxInt32
      costMatrix[rune('>')][i][j] = math.MaxInt32
      costMatrix[rune('<')][i][j] = math.MaxInt32
    }
  }
  costMatrix[rune('>')][startX][startY] = 0
  for len(pq.data) > 0 {
    currNode := pq.Pop()
    if currNode.x == endX && currNode.y == endY {
      return currNode.cost, costMatrix
    }
    // Check straight
    dx, dy := util.TranslateDirection(currNode.direction)
    newX := currNode.x + dx
    newY := currNode.y + dy
    if data[newX][newY] != rune('#') {
      newCost := currNode.cost + 1
      if newCost < costMatrix[currNode.direction][newX][newY] {
        costMatrix[currNode.direction][newX][newY] = newCost
        pq.Append(newNode(newX, newY, currNode.direction, newCost))
      }
    }
    // Check right
    newDirection := util.TurnRight(currNode.direction)
    dx, dy = util.TranslateDirection(newDirection)
    newX = currNode.x + dx
    newY = currNode.y + dy
    if data[newX][newY] != rune('#') {
      newCost := currNode.cost + 1001
      if newCost < costMatrix[newDirection][newX][newY] {
        costMatrix[newDirection][currNode.x][currNode.y] = newCost-1
        costMatrix[newDirection][newX][newY] = newCost
        pq.Append(newNode(newX, newY, newDirection, newCost))
      }
    }
    // Check left
    newDirection = util.TurnLeft(currNode.direction)
    dx, dy = util.TranslateDirection(newDirection)
    newX = currNode.x + dx
    newY = currNode.y + dy
    if data[newX][newY] != rune('#') {
      newCost := currNode.cost + 1001
      if newCost < costMatrix[newDirection][newX][newY] {
        costMatrix[newDirection][currNode.x][currNode.y] = newCost-1
        costMatrix[newDirection][newX][newY] = newCost
        pq.Append(newNode(newX, newY, newDirection, newCost))
      }
    }
  }
  return -1, costMatrix
}

func task1(data [][]rune, debug bool) {
  result := 0
  endX, endY := findStart(data)
  startX, startY := findEnd(data)
  result, _ = findShortest(data, startX, startY, endX, endY, rune('>'))
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]rune, debug bool) {
  result := 0
  endX, endY := findStart(data)
  startX, startY := findEnd(data)
  shortestPath, costMatrix1 := findShortest(data, startX, startY, endX, endY, rune('>'))
  _, costMatrix2 := findShortest(data, endX, endY, startX, startY, rune('^'))
  _, costMatrix3 := findShortest(data, endX, endY, startX, startY, rune('>'))
  _, costMatrix4 := findShortest(data, endX, endY, startX, startY, rune('v'))
  _, costMatrix5 := findShortest(data, endX, endY, startX, startY, rune('<'))
  for d := range costMatrix1 {
    flippedDir := util.TurnRight(util.TurnRight(d))
    for i := range costMatrix1[d] {
      for j := range costMatrix1[d][i] {
        for _, v := range []int{costMatrix2[flippedDir][i][j], 
        costMatrix3[flippedDir][i][j], 
        costMatrix4[flippedDir][i][j], 
        costMatrix5[flippedDir][i][j]} {
          flipSum := costMatrix1[d][i][j] + v
          if flipSum == shortestPath {
            data[i][j] = rune('O')
          }
        }
      }
    }
  }
  for i := range data {
    for j := range data[i] {
      if data[i][j] == rune('O') || data[i][j] == rune('S') || data[i][j] == rune('E') {
        result++
      }
    }
  }
  if debug {
    printData(data)
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
