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

type equation struct {
  result int
  numbers []int
}

func newEquation(result int, numbers []int) *equation {
  return &equation{result: result, numbers: numbers}
}

func (e *equation) String() string {
  return fmt.Sprintf("%d: %v", e.result, e.numbers)
}

func printData(data []*equation) {
  for _, d := range data {
    fmt.Println(d)
  }
}

func operatorSymbol(o string) string {
  switch o {
  case "ADD":
    return "+"
  case "MUL":
    return "*"
  case "CON":
    return "||"
  default:
    return "UNKNOWN"
  }
}

func printSolution(eq *equation, operators []string) {
  output := fmt.Sprintf("%d = %d", eq.result, eq.numbers[0])
  for i := 1; i < len(eq.numbers); i++ {
    output += fmt.Sprintf(" %v %v", operatorSymbol(operators[i-1]), eq.numbers[i])
  }
  fmt.Println(output)
}

func readInput(path string) ([]*equation, error) {
  var data []*equation 
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    splitLine := strings.Split(line, ":")
    result, _ := strconv.Atoi(splitLine[0])
    numbersStr := strings.Split(strings.Trim(splitLine[1], " "), " ")
    numbers := make([]int, len(numbersStr))
    for i, numStr := range numbersStr {
      num, _ := strconv.Atoi(numStr)
      numbers[i] = num
    }
    data = append(data, newEquation(result, numbers))
  }
  return data, nil
}

func getResult(numbers []int, operators []string) int {
  result := numbers[0]
  for i := 1; i < len(numbers); i++ {
    switch operators[i-1] {
    case "ADD":
      result += numbers[i]
    case "MUL":
      result *= numbers[i]
    case "CON":
      result, _ = strconv.Atoi(strconv.Itoa(result)+strconv.Itoa(numbers[i]))
    default:
      continue
    }
  }
  return result
}

func getOperatorPermutations(operators []string, l int, current []string, permutations *[][]string) {
  if len(current) == l {
    permCopy := make([]string, len(current))
    copy(permCopy, current)
    *permutations = append(*permutations, permCopy)
    return
  }
  for _, o := range operators {
    getOperatorPermutations(operators, l, append(current, o), permutations)
  }
}

func getCalibration(data []*equation, operators []string, debug bool) int {
  result := 0
  if debug {
    fmt.Println("\nValid equations:")
  }
  for _, eq := range data {
    var permutations [][]string
    getOperatorPermutations(operators, len(eq.numbers)-1, []string{}, &permutations)
    for _, p := range permutations {
      tempRes := getResult(eq.numbers, p)
      if tempRes == eq.result {
        result += eq.result
        if debug {
          printSolution(eq, p)
        }
        break
      }
    }
  }
  return result
}

func task1(data []*equation, debug bool) {
  operators := []string{"ADD", "MUL"}
  result := getCalibration(data, operators, debug)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []*equation, debug bool) {
  operators := []string{"ADD", "MUL", "CON"}
  result := getCalibration(data, operators, debug)
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
