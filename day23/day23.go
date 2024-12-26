package main

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"adventOfCode2024/util"
)

type set struct {
  data []string
}

func newSet(data  []string) *set {
  return &set{data: data}
}

func (s *set) Append(el string) bool {
  for _, d := range s.data {
    if d == el {
      return false
    }
  }
  s.data = append(s.data, el)
  slices.Sort(s.data)
  return true
}

func (s *set) Contains(el string) bool {
  for _, d := range s.data {
    if d == el {
      return true
    }
  }
  return false
}

func (s *set) Equals(s2 *set) bool {
  if len(s.data) != len(s2.data) {
    return false
  }
  for i := range s.data {
    if s.data[i] != s2.data[i] {
      return false
    }
  }
  return true
}

func (s *set) Len() int {
  return len(s.data)
}

func (s *set) String() string {
  return strings.Join(s.data, ",")
}

func printData(d [][]string) {
  for _, s := range d {
    fmt.Printf("%s-%s\n", s[0], s[1])
  }
}

func readInput(path string) ([][]string, error) {
  var data [][]string
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    data = append(data, strings.Split(line, "-"))
  }
  return data, nil
}

func printConnectionMap(cm map[string][]string) {
  for pc, conns := range cm {
    fmt.Printf("%s is connected to %v\n", pc, conns)
  }
}

func getConnectionMap(data [][]string) map[string][]string {
  result := make(map[string][]string)
  for _, conn := range data {
    if _, ok := result[conn[0]]; !ok {
      result[conn[0]] = []string{}
    }
    result[conn[0]] = append(result[conn[0]], conn[1])
    if _, ok := result[conn[1]]; !ok {
      result[conn[1]] = []string{}
    }
    result[conn[1]] = append(result[conn[1]], conn[0])
  }
  return result
}

func containsSet(slice []*set, s *set) bool {
  for _, el := range slice {
    if el.Equals(s) {
      return true
    }
  }
  return false
}

func getLanParty(connectionMap map[string][]string, oldSet *set, currPC string, limit int) *set {
  if limit != -1 && oldSet.Len() == limit {
    return oldSet
  }
  for _, oldPc := range oldSet.data {
    if !slices.Contains(connectionMap[currPC], oldPc) {
      return oldSet
    }
  }
  newSet := oldSet
  newSet.Append(currPC)
  bestSet := newSet
  for _, potentialPC := range connectionMap[currPC] {
    if newSet.Contains(potentialPC) {
      continue
    }
    potentialSet := getLanParty(connectionMap, newSet, potentialPC, limit)
    if potentialSet.Len() > bestSet.Len() {
      bestSet = potentialSet
    }
  }
  return bestSet
}

func getConnectionSets(data [][]string, debug bool) []*set {
  connectionMap := getConnectionMap(data)
  if debug {
    fmt.Println("Connection map:")
    printConnectionMap(connectionMap)
  }
  result := []*set{}
  for pc1 := range maps.Keys(connectionMap) {
    for _, pc2 := range connectionMap[pc1] {
      for _, pc3 := range connectionMap[pc2] {
        if !slices.Contains(connectionMap[pc3], pc1) {
          continue
        }
        s := newSet([]string{pc1, pc2, pc3})
        if containsSet(result, s) {
          continue
        }
        result = append(result, s)
      }
    }
  }
  return result
}

func task1(data [][]string, debug bool) {
  result := 0
  connectionSets := getConnectionSets(data, debug)
  if debug {
    fmt.Println("Connection sets:")
  }
  for _, s := range connectionSets {
    if debug {
      fmt.Println(s)
    }
    for _, si := range s.data {
      if strings.HasPrefix(si, "t") {
        result++
        break
      }
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data [][]string, debug bool) {
  result := ""
  connectionMap := getConnectionMap(data)
  if debug {
    fmt.Println("Connection map:")
    printConnectionMap(connectionMap)
  }
  besSet := newSet([]string{})
  for k := range connectionMap {
    potentialSet := getLanParty(connectionMap, newSet([]string{}), k, -1)
    if potentialSet.Len() > besSet.Len() {
      besSet = potentialSet
    }
  }
  result = strings.Join(besSet.data, ",")
  fmt.Printf("Task 2: %s\n", result)
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
