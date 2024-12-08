package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"time"
	"unicode"

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
  return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

func containsPoint(s []*point, p *point) bool {
  for _, ps := range s {
    if ps.x == p.x && ps.y == p.y {
      return true
    }
  }
  return false
}

func printData(d [][]rune) {
  for _, s := range d {
    fmt.Println(string(s))
  }
}

func printAntennas(a map[rune][]*point) {
  for k, v := range a {
    fmt.Printf("%v: %v\n", string(k), v)
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

func getAntennas(data [][]rune) map[rune][]*point {
  antennas := make(map[rune][]*point)
  for i, line := range data {
    for j, sym := range line {
      if !unicode.IsUpper(sym) && !unicode.IsLower(sym) && !unicode.IsDigit(sym) {
        continue
      }
      _, ok := antennas[sym]
      if !ok {
        antennas[sym] = []*point{}
      }
      antennas[sym] = append(antennas[sym], newPoint(i, j))
    }
  }
  antennas[rune('#')] = []*point{}
  return antennas
}

func getNewPoints(p *point, dx, dy, cnt, rowLen, colLen int) []*point {
  // newPoint(p1.x - dx, p1.y + dy), newPoint(p2.x + dx, p2.y - dy)
  var newPoints []*point
  np := newPoint(p.x, p.y)
  for range cnt {
    np = newPoint(np.x - dx, np.y + dy)
    if np.x < 0 || np.x >= rowLen || np.y < 0 || np.y >= colLen {
      break
    }
    newPoints = append(newPoints, np)
  }
  return newPoints
}

func getAntinodes(p1, p2 *point, cnt, rowLen, colLen int) []*point {
  dx := p2.x - p1.x
  dy := p1.y - p2.y
  newPoints1 := getNewPoints(p1, dx, dy, cnt, rowLen, colLen)
  newPoints2 := getNewPoints(p2, dx*-1, dy*-1, cnt, rowLen, colLen)
  return slices.Concat(newPoints1, newPoints2)
}

func task1(data [][]rune, debug bool) {
  antennas := getAntennas(data)
  rowLen := len(data[0])
  colLen := len(data)
  for sym, pts := range antennas {
    if sym == rune('#') {
      continue
    }
    for i := 0; i < len(pts); i++ {
      for j := i+1; j <len(pts); j++ {
        antinodes := getAntinodes(pts[i], pts[j], 1, rowLen, colLen)
        for _, a := range antinodes {
          if containsPoint(antennas[rune('#')], a) {
            continue
          }
          antennas[rune('#')] = append(antennas[rune('#')], a)
          data[a.x][a.y] = rune('#')
        }
      }
    }
  }
  if debug {
    printAntennas(antennas)
    printData(data)
  }
  result := len(antennas[rune('#')])
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]rune, debug bool) {
  antennas := getAntennas(data)
  rowLen := len(data[0])
  colLen := len(data)
  for  sym := range maps.Keys(antennas) {
    if sym == rune('#') {
      continue
    }
    antennas[rune('#')] = slices.Concat(antennas[rune('#')], antennas[sym])
  }
  for sym, pts := range antennas {
    if sym == rune('#') {
      continue
    }
    for i := 0; i < len(pts); i++ {
      for j := i+1; j <len(pts); j++ {
        antinodes := getAntinodes(pts[i], pts[j], max(rowLen, colLen), rowLen, colLen)
        for _, a := range antinodes {
          if containsPoint(antennas[rune('#')], a) {
            continue
          }
          antennas[rune('#')] = append(antennas[rune('#')], a)
          data[a.x][a.y] = rune('#')
        }
      }
    }
  }
  if debug {
    printAntennas(antennas)
    printData(data)
  }
  result := len(antennas[rune('#')])
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
