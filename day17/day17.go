package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"adventOfCode2024/util"
)

type computer struct {
  registerA          int
  registerB          int
  registerC          int
  instructionPointer int
  program            []int
  output             []string
}

func newComputer(a, b, c int, program []int) *computer {
  return &computer{registerA: a, registerB: b, registerC: c, program: program, instructionPointer: 0}
}

func (c *computer) String() string {
  outStr := fmt.Sprintf("Register A: %d\n", c.registerA)
  outStr += fmt.Sprintf("Register B: %d\n", c.registerB)
  outStr += fmt.Sprintf("Register C: %d\n", c.registerC)
  outStr += "Program: "
  for i := range c.program {
    if i == c.instructionPointer {
      outStr += fmt.Sprintf("[%d],", c.program[i])
      continue
    }
    outStr += fmt.Sprintf("%d,", c.program[i])
  }
  outStr = fmt.Sprintf("%s\n", outStr[:len(outStr)-1])
  outStr += fmt.Sprintf("Output: %s\n", strings.Join(c.output, ","))
  return outStr
}

func (c *computer) comboOperand(operand int) int {
  switch operand {
  case 0,1,2,3:
    return operand
  case 4:
    return c.registerA
  case 5:
    return c.registerB
  case 6:
    return c.registerC
  default:
    fmt.Printf("ERROR: operand %d is not suppoerted!\n", operand)
    return -1
  }
}

func (c *computer) runOpcode() {
  opcode := c.program[c.instructionPointer]
  literalOperand := c.program[c.instructionPointer+1]
  comboOperand := c.comboOperand(literalOperand)
  switch opcode {
  case 0:
    // adv
    c.registerA = int(float64(c.registerA) / math.Pow(float64(2), float64(comboOperand)))
  case 1:
    // bxl
    c.registerB = c.registerB ^ literalOperand
  case 2:
    // bst
    c.registerB = comboOperand % 8
  case 3:
    // jnz
    if c.registerA != 0 {
      c.instructionPointer = literalOperand
      return
    }
  case 4:
    // bxc
    c.registerB = c.registerB ^ c.registerC
  case 5:
    // out
    c.output = append(c.output, strconv.Itoa(comboOperand % 8))
  case 6:
    // bdv
    c.registerB = int(float64(c.registerA) / math.Pow(float64(2), float64(comboOperand)))
  case 7:
    // bdv
    c.registerC = int(float64(c.registerA) / math.Pow(float64(2), float64(comboOperand)))
  }
  c.instructionPointer += 2
}

func printData(d *computer) {
  fmt.Println(d)
}

func readInput(path string) (*computer, error) {
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return nil, err
  }
  scanner := bufio.NewScanner(file)
  var regA, regB, regC int
  var program []int
  for scanner.Scan() {
    line := scanner.Text()
    splitLine := strings.Split(line, ": ")
    switch splitLine[0] {
    case "Register A":
      regA, _ = strconv.Atoi(splitLine[1])
    case "Register B":
      regB, _ = strconv.Atoi(splitLine[1])
    case "Register C":
      regC, _ = strconv.Atoi(splitLine[1])
    case "Program":
      programLine := strings.Split(splitLine[1], ",")
      program = make([]int, len(programLine))
      for i := range programLine {
        pInt, _ := strconv.Atoi(programLine[i])
        program[i] = pInt
      }
    default:
      continue
    }
  }
  return newComputer(regA, regB, regC, program), nil
}

func (c *computer) getProgramStr() string {
  outStr := ""
  for i := range c.program {
    outStr += fmt.Sprintf("%d,", c.program[i])
  }
  return fmt.Sprintf("%s", outStr[:len(outStr)-1])
}

func (c *computer) runProgram(debug bool) {
  for c.instructionPointer < len(c.program) {
    c.runOpcode()
    if debug {
      printData(c)
    }
  }
}

func task1(data *computer, debug bool) {
  result := ""
  data.runProgram(debug)
  result = strings.Join(data.output, ",")
  fmt.Printf("Task 1: %s\n", result)
}

func subFunction(data *computer, regA int, targetOutput string, wg *sync.WaitGroup, ch chan<-int, debug bool) {
  defer wg.Done()
  for regA < regA+int(math.Pow10(14)) {
    if debug {
      fmt.Printf("Checking %d\n", regA)
    }
    dataCopy := newComputer(regA, data.registerB, data.registerC, data.program)
    dataCopy.runProgram(false)
    if targetOutput == strings.Join(dataCopy.output, ","){
      ch<-regA
      return
    }
    regA++
  }
}

func task2(data *computer, debug bool) {
  result := 0
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
