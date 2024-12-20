package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"adventOfCode2024/util"
)

type node struct {
  x    int
  y    int
  cost int
}

func newNode(x, y, cost int) *node {
  return &node{x: x, y: y, cost: cost}
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

func readInput(path string, size int) ([][]rune, [][]int, error) {
  data := make([][]rune, size+1)
  for i := range data {
    row := make([]rune, size+1)
    for j := range row {
      row[j] = rune('.')
    }
    data[i] = row
  }
  var corruptedData [][]int
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, corruptedData, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    splitLine := strings.Split(line, ",")
    x, _ := strconv.Atoi(splitLine[1])
    y, _ := strconv.Atoi(splitLine[0])
    corruptedData = append(corruptedData, []int{x, y})
  }
  return data, corruptedData, nil
}

func corruptData(data [][]rune, corruptedData [][]int, count int) {
  for _, cd := range corruptedData {
    if count == 0 {
      break
    }
    data[cd[0]][cd[1]] = rune('#')
    count --
  }
}

func shortestPath(data [][]rune, startX, startY, endX, endY int, debug bool) int {
  minCost := -1
  pq := newPriorityQueue()
  costMatrix := make([][]int, len(data))
  for i := range costMatrix {
    row := make([]int, len(data[0]))
    for j := range data[0] {
      row[j] = math.MaxInt32
    }
    costMatrix[i] = row
  }
  costMatrix[startX][startY] = 0
  pq.Append(newNode(startX, startY, 0))
  directions := [][]int{
    {1,0},
    {-1,0},
    {0,1},
    {0,-1},
  }
  for len(pq.data) > 0 {
    currentNode := pq.Pop()
    if currentNode.x == endX && currentNode.y == endY {
      minCost = currentNode.cost
      break
    }
    for _, d := range directions {
      newX := currentNode.x + d[0]
      newY := currentNode.y + d[1]
      if newX < 0 || newY < 0 || newX >= len(data) || newY >= len(data[0]){
        continue
      }
      if data[newX][newY] == rune('#') {
        continue
      }
      newCost := currentNode.cost + 1
      if newCost >= costMatrix[newX][newY] {
        continue
      }
      costMatrix[newX][newY] = newCost
      pq.Append(newNode(newX, newY, newCost))
    }
  }
  if debug {
    fmt.Println("Cost matrix: ")
    for i := range costMatrix {
      for j := range costMatrix[i] {
        if costMatrix[i][j] == math.MaxInt32 {
          fmt.Printf("%5d", -1)
          continue
        }
        fmt.Printf("%5d", costMatrix[i][j])
      }
      fmt.Println()
    }
  }
  return minCost
}

func task1(data [][]rune, corruptedData [][]int, debug bool) {
  result := 0
  corruptCount := 1024
  corruptData(data, corruptedData, corruptCount)
  if debug {
    fmt.Printf("Data after %d corrupted bytes:\n", corruptCount)
    printData(data)
  }
  result = shortestPath(data, 0, 0, len(data)-1, len(data[0])-1, debug)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]rune, corruptedData [][]int, debug bool) {
  resultx := 0
  resulty := 0
  for corruptCount := 1024; corruptCount < len(corruptedData); corruptCount++ {
    corruptData(data, corruptedData, corruptCount)
    if debug {
      fmt.Printf("Data after %d corrupted bytes:\n", corruptCount)
      printData(data)
    }
    if shortestPath(data, 0, 0, len(data)-1, len(data[0])-1, debug) == -1 {
      resultx = corruptedData[corruptCount-1][1]
      resulty = corruptedData[corruptCount-1][0]
      break
    }
  }
  fmt.Printf("Task 2: %d,%d\n", resultx, resulty)
}

func Run(path string, taskId int, debug bool) error {
  data, corruptedData, err := readInput(path, 70)
  if err != nil {
    return err
  }
  if debug {
    fmt.Printf("\nRunning task %d\n", taskId)
  }
  switch taskId {
  case 1:
    task1(data, corruptedData, debug)
  case 2:
    task2(data, corruptedData, debug)
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
