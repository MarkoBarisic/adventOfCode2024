package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
  "slices"
	"strings"
	"time"

	"adventOfCode2024/util"
)

type neibghours struct {
  left *gardenPlot
  right *gardenPlot
  up *gardenPlot
  down *gardenPlot
}

func newNeibghours() *neibghours {
  return &neibghours{left: nil, right: nil, up: nil, down: nil}
}

type gardenPlot struct {
  x int
  y int
  symbol string
  neibghours *neibghours
}

type region struct {
  area int
  perimiter int
  sides int
  symbol string
  gardenPlots []*gardenPlot
}

func newGardenPlot(x, y int, symbol string) *gardenPlot {
  return &gardenPlot{x: x, y: y, symbol: symbol, neibghours: newNeibghours()}
}

func (gp *gardenPlot) String() string {
  return fmt.Sprintf("%s(%d, %d)", gp.symbol, gp.x, gp.y)
}

func (gp1 *gardenPlot) addNeibghour(gp2 *gardenPlot) {
  dx := gp1.x - gp2.x
  dy := gp1.y - gp2.y
  if dy == 0 {
    if dx == -1 {
      gp1.neibghours.up = gp2
      return
    }
    if dx == 1 {
      gp1.neibghours.down = gp2
      return
    }
    return
  }
  if dx == 0 {
    if dy == -1 {
      gp1.neibghours.left = gp2
      return
    }
    if dy == 1 {
      gp1.neibghours.right = gp2
      return
    }
    return
  }
}

func newRegion(gardenPlots []*gardenPlot) *region {
  return &region{area: -1, perimiter: -1, sides: -1, gardenPlots: gardenPlots, symbol: gardenPlots[0].symbol}
}

func (r *region) String() string {
  return fmt.Sprintf("Plots: %v\nArea: %d\nPerimiter: %d\nSides: %d\n", r.gardenPlots, r.area, r.perimiter, r.sides)
}

func (r *region) appendGardenPlot(gp *gardenPlot) {
  if ! containsPlot(r.gardenPlots, gp) {
    r.gardenPlots = append(r.gardenPlots, gp)
  }
}

func areNeibghours(gp1, gp2 *gardenPlot) bool {
  dx := gp1.x - gp2.x
  dy := gp1.y - gp2.y
  if (dx == 1 || dx == -1) && dy == 0{
    return true
  }
  if (dy == 1 || dy == -1) && dx == 0{
    return true
  }
  return false
}

func (r *region) calculateArea() {
  r.area = len(r.gardenPlots)
}

func (r *region) calculatePerimiter() {
  r.sortPlots()
  r.findNeibghours()
  r.perimiter = 0
  for _, gp1 := range r.gardenPlots {
    r.perimiter += 4
    if gp1.neibghours.up != nil {
      r.perimiter--
    }
    if gp1.neibghours.down != nil {
      r.perimiter--
    }
    if gp1.neibghours.left != nil {
      r.perimiter--
    }
    if gp1.neibghours.right != nil {
      r.perimiter--
    }
  }
}

func (r *region) findNeibghours() {
  for i := range r.gardenPlots {
    for j := range r.gardenPlots {
      r.gardenPlots[i].addNeibghour(r.gardenPlots[j])
    }
  }
}

func (r *region) sortPlots() {
  slices.SortFunc(r.gardenPlots, func(gp1, gp2 *gardenPlot) int {
    if gp1.x == gp2.x {
      if gp1.y == gp2.y {
        return 0
      }
      if gp1.y > gp2.y {
        return 1
      }
      return -1
    }
    if gp1.x > gp2.x {
      return 1
    }
    return -1
  })
}

func (r *region) calculateSides() {
  r.sortPlots()
  r.findNeibghours()
  r.sides = 0
  for _, gp := range r.gardenPlots {
    // Check up
    if gp.neibghours.up == nil {
      if gp.neibghours.left == nil {
        r.sides++
      } else if gp.neibghours.left.neibghours.up != nil {
        r.sides++
      }
    }
    // Check down
    if gp.neibghours.down == nil {
      if gp.neibghours.left == nil {
        r.sides++
      } else if gp.neibghours.left.neibghours.down != nil {
        r.sides++
      }
    }
    // Check left
    if gp.neibghours.left == nil {
      if gp.neibghours.up == nil {
        r.sides++
      } else if gp.neibghours.up.neibghours.left != nil {
        r.sides++
      }
    }
    // Check right
    if gp.neibghours.right == nil {
      if gp.neibghours.up == nil {
        r.sides++
      } else if gp.neibghours.up.neibghours.right != nil {
        r.sides++
      }
    }
  }
}

func printData(d [][]*gardenPlot) {
  for _, s := range d {
    for _, si := range s {
      fmt.Print(si.symbol)
    }
    fmt.Println()
  }
}

func readInput(path string) ([][]*gardenPlot, error) {
  var data [][]*gardenPlot
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  i := -1
  for scanner.Scan() {
    i++
    line := scanner.Text()
    splitLine := strings.Split(line, "")
    row := make([]*gardenPlot, len(splitLine))
    for j, s := range splitLine {
      row[j] = newGardenPlot(i, j, s)
    }
    data = append(data, row)
  }
  return data, nil
}

func containsPlot(gps []*gardenPlot, gp *gardenPlot) bool {
  for _, g := range gps {
    if g.x == gp.x && g.y == gp.y && g.symbol == gp.symbol{
      return true
    }
  }
  return false
}

func mapRegion(data [][]*gardenPlot, x, y int, usedPlots *[]*gardenPlot, region *region) {
  if x < 0 || y < 0 || x >= len(data) || y >= len(data[0]) {
    return
  }
  if data[x][y].symbol != region.symbol {
    return
  }
  if containsPlot(*usedPlots, data[x][y]) {
    return
  }
  region.appendGardenPlot(data[x][y])
  *usedPlots = append(*usedPlots, data[x][y])
  mapRegion(data, x+1, y, usedPlots, region)
  mapRegion(data, x-1, y, usedPlots, region)
  mapRegion(data, x, y+1, usedPlots, region)
  mapRegion(data, x, y-1, usedPlots, region)
  return
}

func task1(data [][]*gardenPlot, debug bool) {
  result := 0
  var usedPlots []*gardenPlot
  for i := range data {
    for j := range data[i] {
      if containsPlot(usedPlots, data[i][j]) {
        continue
      }
      currRegion := newRegion([]*gardenPlot{data[i][j]})
      mapRegion(data, i, j, &usedPlots, currRegion)
      currRegion.calculateArea()
      currRegion.calculatePerimiter()
      result += currRegion.area*currRegion.perimiter
      if debug {
        fmt.Println(currRegion)
      }
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]*gardenPlot, debug bool) {
  result := 0
  var usedPlots []*gardenPlot
  for i := range data {
    for j := range data[i] {
      if containsPlot(usedPlots, data[i][j]) {
        continue
      }
      currRegion := newRegion([]*gardenPlot{data[i][j]})
      mapRegion(data, i, j, &usedPlots, currRegion)
      currRegion.calculateArea()
      currRegion.calculatePerimiter()
      currRegion.calculateSides()
      result += currRegion.area*currRegion.sides
      if debug {
        fmt.Println(currRegion)
      }
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
