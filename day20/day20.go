package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

  "adventOfCode2024/util"
)

type node struct {
  x      int
  y      int
  symbol rune
  cost   int
}

func (n *node) String() string {
  return fmt.Sprintf("%s(%d,%d)", string(n.symbol), n.x, n.y)
}

type cheat struct {
  start *node
  end   *node
  cost  int
}

func (c *cheat) String() string {
  return fmt.Sprintf("Cheat: %v --> %v saves %d picoseconds", c.start, c.end, c.cost)
}

type queue struct {
  data []*node
}

func (q *queue) pop() *node {
  el := q.data[0]
  q.data = q.data[1:]
  return el
}

func (q *queue) push(el *node) {
  q.data = append(q.data, el)
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

func findStart(data [][]rune) (int, int) {
  for i := range data {
    for j := range data[i] {
      if data[i][j] == rune('S') {
        return i,j
      }
    }
  }
  return -1, -1
}

func getCheatArea(size int) [][]int {
  cheatArea := [][]int{}
  for i := 0; i <= size; i++ {
    for j := 0; j <= size - i; j++ {
      for _, combo := range [][]int{
        {i,j},
        {i,-j},
        {-i,j},
        {-i,-j},
      } {
        contains := false
        for _, ca := range cheatArea {
          if ca[0] == combo[0] && ca[1] == combo[1] {
            contains = true
            break
          }
        }
        if contains {
          continue
        }
        cheatArea = append(cheatArea, combo)
      }
    }
  }
  return cheatArea
}

func containsCheat(s []*cheat, cheat *cheat) bool {
  for _, c := range s {
    if c.start.x != cheat.start.x || c.start.y != cheat.start.y {
      continue
    }
    if c.end.x != cheat.end.x || c.end.y != cheat.end.y {
      continue
    }
    return true
  }
  return false
}

func abs(x int) int {
  if x < 0 {
    return -x
  }
  return x
}

func findCheats(data [][]rune, costMatrix [][]int, path *queue, cheatSize int, limit int, debug bool) int {
  result := 0
  cheatArea := getCheatArea(cheatSize)
  cheats := []*cheat{}
  originalLen := len(path.data)
  cnt := 0
  for len(path.data) > 0 {
    if debug {
      cnt++
      fmt.Printf("%d/%d\n", cnt, originalLen)
    }
    currNode := path.pop()
    for _, direction := range cheatArea {
      newX := currNode.x + direction[0]
      newY := currNode.y + direction[1]
      if newX < 0 || newY < 0 || newX >= len(costMatrix) || newY >= len(costMatrix[0]) {
        continue
      }
      newCost := costMatrix[newX][newY]
      if newCost == -1 {
        continue
      }
      if newCost < currNode.cost {
        continue
      }
      stepsTaken := abs(direction[0]) + abs(direction[1])
      saving := newCost - currNode.cost - stepsTaken
      if saving <= 0 {
        continue
      }
      newCheat := &cheat{
        start: currNode, 
        end: &node{x: newX, y: newY, cost: newCost, symbol: data[newX][newY]},
        cost: newCost - currNode.cost - stepsTaken,
      }
      if containsCheat(cheats, newCheat) {
        continue
      }
      if newCheat.cost >= limit {
        result++
      }
      cheats = append(cheats, newCheat)
    }
  }
  if debug {
    cheatCnt := make(map[int]int)
    for _, c := range cheats {
      if c.cost < limit {
        continue
      }
      if _, ok := cheatCnt[c.cost]; !ok {
        cheatCnt[c.cost] = 0
      }
      cheatCnt[c.cost]++
      fmt.Println(c)
    }
    for k,v := range cheatCnt {
      fmt.Printf("- There are %d cheats that save %d picoseconds\n", v, k)
    }
  }
  return result
}

func runTrack(data [][]rune, debug bool) ([][]int, *queue) {
  q := &queue{data: []*node{}}
  path := &queue{data: []*node{}}
  startX, startY := findStart(data)
  q.push(&node{x: startX, y: startY, symbol:data[startX][startY], cost: 0})
  costMatrix := make([][]int, len(data))
  for i := range len(data) {
    row := make([]int, len(data[0]))
    for j := range len(data[0]) {
      row[j] = -1
    }
    costMatrix[i] = row
  }
  costMatrix[startX][startY] = 0
  for len(q.data) > 0 {
    currNode := q.pop()
    path.push(currNode)
    if currNode.symbol == rune('E') {
      if debug {
        for i := range costMatrix {
          for j := range costMatrix[i] {
            if costMatrix[i][j] == -1 {
              fmt.Print("  #")
              continue
            }
            fmt.Printf("%3d", costMatrix[i][j])
          }
          fmt.Println()
        }
      }
      return costMatrix, path
    }
    for _, direction := range [][]int{
      {1,0},
      {-1,0},
      {0,1},
      {0,-1},
    }{
      newX := currNode.x + direction[0]
      newY := currNode.y + direction[1]
      if data[newX][newY] == rune('#') {
        continue
      }
      if costMatrix[newX][newY] != -1 {
        continue
      }
      newCost := currNode.cost + 1
      costMatrix[newX][newY] = newCost
      q.push(&node{x: newX, y: newY, cost: newCost, symbol: data[newX][newY]})
    }
  }
  return costMatrix, path
}

func task1(data [][]rune, debug bool) {
  result := 0
  costMatrix, path := runTrack(data, debug)
  result = findCheats(data, costMatrix, path, 2, 100, debug)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]rune, debug bool) {
  result := 0
  costMatrix, path := runTrack(data, debug)
  result = findCheats(data, costMatrix, path, 20, 100, debug)
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
