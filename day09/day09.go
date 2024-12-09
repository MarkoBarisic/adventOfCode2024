package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"adventOfCode2024/util"
)

func printData(d []string) {
  fmt.Println(strings.Join(d, ""))
}

func readInput(path string) ([]string, error) {
  var data []string
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    data = slices.Concat(data, strings.Split(line, ""))
  }
  return data, nil
}

func transformData(data []string) []string {
  var td []string
  currId := 0
  for i, dS := range data {
    dI, _ := strconv.Atoi(dS)
    if i%2 == 0 {
      s := slices.Repeat([]string{strconv.Itoa(currId)}, dI)
      td = slices.Concat(td, s)
      currId += 1
    } else {
      s := slices.Repeat([]string{"."}, dI)
      td = slices.Concat(td, s)
    }
  }
  return td
}

func lastDigit(data []string) int {
  for i := len(data)-1; i >= 0; i-- {
    if strings.ContainsAny(data[i], "0123456789") {
      return i
    }
  }
  return -1
}

func digitBlock(data []string, digit int) (int, int) {
  var ind, size int
  digitS := strconv.Itoa(digit)
  for i := len(data)-1; i >= 0; i-- {
    if data[i] == digitS {
      size = 0
      for j := i; j >= 0; j-- {
        if data[j] != digitS {
          break
        }
        ind = j
        size += 1
      }

      return ind, size
    }
  }
  return -1, -1
}

func firstSpace(data []string) int {
  for i := 0; i < len(data); i++ {
    if data[i] == "." {
      return i
    }
  }
  return -1
}

func findSpaceBlock(data []string, size int) int {
  for i := 0; i < len(data); {
    if data[i] != "." {
      i++
      continue
    }
    iSize := 0
    for _, d := range data[i:] {
      if d != "." {
        break
      }
      iSize++
    }
    if iSize < size {
      i += iSize
      continue
    }
    return i
  }
  return -1
}

func checksum(data []string) int {
  chksm := 0
  for i, dS := range data {
    if dS == "0" || dS == "." {
      continue
    }
    dI, _ := strconv.Atoi(dS)
    chksm += i*dI
  }
  return chksm
}

func task1(data []string, debug bool) {
  result := 0
  data = transformData(data)
  if debug {
    fmt.Println("Transformed data:")
    printData(data)
    fmt.Println("-----------------")
  }
  for {
    lastDigitIndex := lastDigit(data)
    firstSpaceIndex := firstSpace(data)
    if debug {
      printData(data)
    }
    if lastDigitIndex < firstSpaceIndex {
      break
    }
    data[firstSpaceIndex] = data[lastDigitIndex]
    data[lastDigitIndex] = "."
  }
  result = checksum(data)
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []string, debug bool) {
  result := 0
  data = transformData(data)
  if debug {
    fmt.Println("Transformed data:")
    printData(data)
    fmt.Println("-----------------")
  }
  currIdS := data[lastDigit(data)]
  currId, _ := strconv.Atoi(currIdS)
  if debug {
    fmt.Printf("Highest index is %d\n", currId)
  }
  for currId > 0 {
    digitIndex, digitBlockSize := digitBlock(data, currId)
    if debug {
      printData(data)
      fmt.Printf("Looking for space block of size %d for digit %d\n", digitBlockSize, currId)
    }
    spaceBlockIndex := findSpaceBlock(data, digitBlockSize)
    if spaceBlockIndex == -1 {
      currId--
      continue
    }
    if digitIndex < spaceBlockIndex {
      currId--
      continue
    }
    for i := range digitBlockSize {
      data[spaceBlockIndex+i] = data[digitIndex+i]
      data[digitIndex+i] = "."
    }
    currId--
  }
  result = checksum(data)
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
