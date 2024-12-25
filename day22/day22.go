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

type diffSequence struct {
  price    int
  sequence []int
}

type buyer struct {
  currSecret *secretNumber
  secrets    []*secretNumber
}

func newBuyer(init int) *buyer {
  s := newSecretNumber(init)
  return &buyer{secrets: []*secretNumber{s}, currSecret: s}
}

type secretNumber struct {
  value     int
  price     int
  priceDiff string
}

func newSecretNumber(value int) *secretNumber {
  return &secretNumber{value: value, price: value%10, priceDiff: ""}
}

func (sn *secretNumber) String() string {
  return fmt.Sprintf("%d: %d (%s)", sn.value, sn.price, sn.priceDiff)
}

func printData(d []*buyer) {
  for _, b := range d {
    for _, s := range b.secrets {
      fmt.Println(s)
    }
  }
}

func readInput(path string) ([]*buyer, error) {
  var data []*buyer
  file, err := os.Open(path)
  defer file.Close()
  if err != nil {
    return data, err
  }
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    init, _ := strconv.Atoi(line)
    data = append(data, newBuyer(init))
  }
  return data, nil
}

func (sn *secretNumber) mix(val int) {
  sn.value = sn.value ^ val
  sn.price = sn.value%10
}

func (sn *secretNumber) prune() {
  sn.value = sn.value % 16777216
  sn.price = sn.value%10
}

func (b *buyer) evolveSecret(iters int) {
  for range iters {
    newSecret := newSecretNumber(b.currSecret.value)
    newSecret.mix(newSecret.value * 64)
    newSecret.prune()
    newSecret.mix(int(math.Floor(float64(newSecret.value)/float64(32))))
    newSecret.prune()
    newSecret.mix(newSecret.value * 2048)
    newSecret.prune()
    newSecret.priceDiff = strconv.Itoa(newSecret.price - b.currSecret.price)
    b.currSecret = newSecret
    b.secrets = append(b.secrets, newSecret)
  }
}

func containsSeq(s []*diffSequence, seq *diffSequence) bool {
  for _, el := range s {
    if el.sequence[0] == seq.sequence[0] && 
    el.sequence[1] == seq.sequence[1] &&
    el.sequence[2] == seq.sequence[2] &&
    el.sequence[3] == seq.sequence[3] {
      return true
    }
  }
  return false
}

func rateSequence(data []*buyer, seqTarget string) int {
  result := 0
  for i := range data {
    for j := 3; j < len(data[i].secrets); j++ {
      seq := strings.Join([]string{
        data[i].secrets[j-3].priceDiff,
        data[i].secrets[j-2].priceDiff,
        data[i].secrets[j-1].priceDiff,
        data[i].secrets[j].priceDiff,
      }, ",")
      if seqTarget != seq {
        continue
      }
      result += data[i].secrets[j].price
      break
    }
  }
  return result
}

func checkExistance(data []*buyer, seqTarget string, seqLen int) bool {
  for i := range data {
    for j := seqLen-1; j < len(data[i].secrets); j++ {
      seq := make([]string, seqLen)
      ki := 0
      for k := seqLen-1; k > -1; k-- {
        seq[ki] = data[i].secrets[j-k].priceDiff
        ki++
      }
      seqStr := strings.Join(seq, ",")
      if seqTarget == seqStr {
        return true
      }
    }
  }
  return false

}

func subFucntion(data []*buyer, seq string, wg *sync.WaitGroup, ch chan<- string) {
  defer wg.Done()
  ch <- fmt.Sprintf("%s %d", seq, rateSequence(data, seq))
}

func rateDiffSequences(data []*buyer, debug bool) int {
  maxPrice := -1
  var wg sync.WaitGroup
  ch := make(chan string, 19*19*19*19)
  for a := 9; a > -10; a-- {
    if !checkExistance(data, fmt.Sprintf("%d", a), 1) {
      fmt.Printf("%12s doesn't exist\n", fmt.Sprintf("%d", a))
      continue
    }
    for b := 9; b > -10; b-- {
      if !checkExistance(data, fmt.Sprintf("%d,%d", b,a), 2) {
        fmt.Printf("%12s doesn't exist\n", fmt.Sprintf("%d,%d", b, a))
        continue
      }
      for c := 9; c > -10; c-- {
        if !checkExistance(data, fmt.Sprintf("%d,%d,%d", c,b,a), 3) {
          fmt.Printf("%12s doesn't exist\n", fmt.Sprintf("%d,%d,%d", c, b, a))
          continue
        }
        for d := 9; d > -10; d-- {
          wg.Add(1)
          seq := fmt.Sprintf("%d,%d,%d,%d", d,c,b,a)
          if debug {
            fmt.Printf("Checking %12s\n", seq)
          }
          go subFucntion(data, seq, &wg, ch)
        }
      }
    }
  }
  go func() {
    wg.Wait()
    close(ch)
  }()
  for msg := range ch {
    msgSplit := strings.Split(msg, " ")
    seqPrice, _ := strconv.Atoi(msgSplit[1])
    if seqPrice > maxPrice {
      maxPrice = seqPrice
    }
  }
  return maxPrice
}

func task1(data []*buyer, debug bool) {
  result := 0
  for _, d := range data {
    d.evolveSecret(2000)
    result += d.currSecret.value
    if debug {
      fmt.Printf("%v --(%d iters)--> %v\n", d.secrets[0], len(d.secrets)-1, d.currSecret)
    }
  }
  fmt.Printf("Task 1: %d\n", result)
}

func task2(data []*buyer, debug bool) {
  result := 0
  for _, d := range data {
    d.evolveSecret(2000)
  }
  result = rateDiffSequences(data, debug)
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
