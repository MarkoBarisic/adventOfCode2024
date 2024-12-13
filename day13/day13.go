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

type button struct {
  symbol string
  dx     int
  dy     int
  cost   int
}

func newButton(symbol string, dx, dy, cost int) *button {
  return &button{symbol: symbol, dx: dx, dy: dy, cost: cost}
}

func (b *button) String() string {
  return fmt.Sprintf("Button %s(%d token): X+%d, Y+%d", b.symbol, b.cost, b.dx, b.dy)
}

type machine struct {
  prizeX  int
  prizeY  int
  buttonA *button
  buttonB *button
}

func newMachine(pX, pY int, bA, bB *button) *machine {
  return &machine{prizeX: pX, prizeY: pY, buttonA: bA, buttonB: bB}
}

func (m *machine) String() string {
  return fmt.Sprintf("%v\n%v\nPrize: X=%d, Y=%d", m.buttonA, m.buttonB, m.prizeX, m.prizeY)
}

func printData(d []*machine) {
  for _, s := range d {
    fmt.Printf("Machine:\n%v\n", s)
  }
}

func readInput(path string) ([]*machine, error) {
  var data []*machine
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  var buttonA *button
  var buttonB *button
  for scanner.Scan() {
    line := scanner.Text()
    splitLine := strings.Split(line, ": ")
    switch splitLine[0] {
    case "Button A":
      buttonLine := strings.Split(splitLine[1], ", ")
      dx, _ := strconv.Atoi(strings.Split(buttonLine[0], "+")[1])
      dy, _ := strconv.Atoi(strings.Split(buttonLine[1], "+")[1])
      buttonA = newButton("A", dx, dy, 3)
    case "Button B":
      buttonLine := strings.Split(splitLine[1], ", ")
      dx, _ := strconv.Atoi(strings.Split(buttonLine[0], "+")[1])
      dy, _ := strconv.Atoi(strings.Split(buttonLine[1], "+")[1])
      buttonB = newButton("B", dx, dy, 1)
    case "Prize":
      prizeLine := strings.Split(splitLine[1], ", ")
      px, _ := strconv.Atoi(strings.Split(prizeLine[0], "=")[1])
      py, _ := strconv.Atoi(strings.Split(prizeLine[1], "=")[1])
      data = append(data, newMachine(px, py, buttonA, buttonB))
    default:
      continue
    }
  }
  return data, nil
}

func howToWin(m *machine, maxIter int) int {
  cntB := ((m.prizeY*m.buttonA.dx)-(m.prizeX*m.buttonA.dy))/((m.buttonB.dy*m.buttonA.dx)-(m.buttonB.dx*m.buttonA.dy))
  cntA := (m.prizeX-(cntB*m.buttonB.dx))/m.buttonA.dx
  if maxIter != -1 && (cntA > maxIter || cntB > maxIter) {
    return 0
  }
  finalX := (m.buttonA.dx*cntA) + (m.buttonB.dx*cntB)
  if m.prizeX % finalX != 0 {
    return 0
  }
  finalY := (m.buttonA.dy*cntA) + (m.buttonB.dy*cntB)
  if m.prizeY % finalY != 0 {
    return 0
  }
  return (m.buttonA.cost*cntA) + (m.buttonB.cost*cntB)
}

func task1(data []*machine) {
  result := 0
  for _, d := range data {
    result += howToWin(d, 100)
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []*machine) {
  result := 0
  for _, d := range data {
    d.prizeX += 10000000000000
    d.prizeY += 10000000000000
    result += howToWin(d, -1)
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
    task1(data)
  case 2:
    task2(data)
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
