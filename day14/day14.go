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

type point struct {
  x int
  y int
}

func newPoint(x, y int) *point {
  return &point{x: x, y: y}
}

func (p *point) String() string {
  return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type robot struct {
  p0       *point
  p        *point
  v        *point
  arena    *arena
  quadrant int
}

func newRobot(p0, v *point) *robot {
  return &robot{p0: p0, v: v, p: newPoint(p0.x, p0.y), arena: nil, quadrant: 0}
}

func (r *robot) String() string {
  return fmt.Sprintf("p0=%v v=%v p=%v", r.p0, r.v, r.p)
}

func (r *robot) setQuadrant() {
  if r.arena == nil {
    fmt.Println("Can't set quadrant when arena is undefined")
  }
  if r.p.x == r.arena.midX || r.p.y == r.arena.midY {
    r.quadrant = 0
    return
  }
  if r.p.x < r.arena.midX {
    if r.p.y < r.arena.midY {
      r.quadrant = 1
      return
    }
    r.quadrant = 3
    return
  }
  if r.p.y < r.arena.midY {
    r.quadrant = 2
    return
  }
  r.quadrant = 4
}

func (r *robot) move() {
  if r.arena == nil {
    fmt.Println("Can't move when arena is undefined")
    return
  }
  r.p.x += r.v.x
  if r.p.x < 0 {
    r.p.x += r.arena.width
  } else if r.p.x >= r.arena.width {
    r.p.x -= r.arena.width
  }
  r.p.y += r.v.y
  if r.p.y < 0 {
    r.p.y += r.arena.height
  } else if r.p.y >= r.arena.height {
    r.p.y -= r.arena.height
  }
  r.setQuadrant()
}

type arena struct {
  width  int
  height int
  midX   int
  midY   int
}

func newArena(width, height int) *arena {
  return &arena{width: width, height: height, midX: findMiddle(width), midY: findMiddle(height)}
}

func findMiddle(x int) int {
  if x%2 == 1 {
    return x/2
  }
  return x/2 - 1
}

func printData(d []*robot) {
  for _, r := range d {
    fmt.Printf("%v --> quadrant %d\n", r, r.quadrant)
  }
}

func printImage(d []*robot) {
  image := make([][]string, d[0].arena.height)
  for i := range image {
    row := make([]string, d[0].arena.width)
    for j := range row {
      row[j] = "."
    }
    image[i] = row
  }
  for _, r := range d {
    image[r.p.y][r.p.x] = "#"
  }
  for i := range image {
    for j := range image[i] {
      fmt.Print(image[i][j])
    }
    fmt.Println()
  }
}

func readInput(path string) ([]*robot, error) {
  var data []*robot
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    splitLine := strings.Split(line, " ")
    pLine := strings.Split(strings.Split(splitLine[0], "=")[1], ",")
    px, _ := strconv.Atoi(pLine[0])
    py, _ := strconv.Atoi(pLine[1])
    vLine := strings.Split(strings.Split(splitLine[1], "=")[1], ",")
    vx, _ := strconv.Atoi(vLine[0])
    vy, _ := strconv.Atoi(vLine[1])
    data = append(data, newRobot(newPoint(px, py), newPoint(vx, vy)))
  }
  return data, nil
}

func containsPoint(s []*point, p *point) bool {
  for _, ps := range s {
    if ps.x == p.x && ps.y == p.y {
      return true
    }
  }
  return false
}

func task1(data []*robot, debug bool) {
  result := 0
  nSeconds := 100
  counter := map [int]int{
    0: 0,
    1: 0,
    2: 0,
    3: 0,
    4: 0,
  }
  arena := newArena(101, 103)
  for _, r := range data {
    r.arena = arena
    for range nSeconds {
      r.move()
    }
    counter[r.quadrant]++
  }
  if debug {
    fmt.Printf("\nAfter %d seconds:\n", nSeconds)
    printData(data)
  }
  result = counter[1] * counter[2] * counter[3] * counter[4]
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []*robot, debug bool) {
  result := 0
  arena := newArena(101, 103)
  for {
    result++
    var points []*point
    isEasterEgg := true
    for _, r := range data {
      r.arena = arena
      r.move()
      if !isEasterEgg {
        continue
      }
      if containsPoint(points, r.p) {
        isEasterEgg = false
        continue
      }
      points = append(points, r.p)
    }
    if !isEasterEgg {
      continue
    }
    if debug {
      fmt.Println("Easter egg picture:")
      printImage(data)
    }
    break
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
