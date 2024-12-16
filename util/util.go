package util

import (
  "errors"
  "fmt"
  "path/filepath"
  "runtime"
  "slices"
  "strings"
)

func ProcessArgs(args []string) (string, bool, error) {
  debug := slices.Contains(args, "--debug")
  inputsFile := "real.txt"
  if slices.Contains(args, "--test") {
    inputsFile = "test.txt"
  }
  _, srcPath, _, ok := runtime.Caller(1)
  if !ok {
    return "", false, errors.New("Error getting the file path")
  }
  srcDir := filepath.Dir(srcPath)
  inputsPath := filepath.Join(srcDir, "inputs", inputsFile)
  absInputsPath, err := filepath.Abs(inputsPath)
  if err != nil {
      return "", false, err
  }
  if debug {
    inputsMsg := fmt.Sprintf("Using inputs from: %s", absInputsPath)
    fmt.Println("Running in debug mode")
    fmt.Println(inputsMsg)
    fmt.Println(strings.Repeat("-", len(inputsMsg)))
  }
  return absInputsPath, debug, nil
}

func TurnRight(direction rune) rune {
  switch direction {
  case rune('^'):
    return rune('>')
  case rune('>'):
    return rune('v')
  case rune('v'):
    return rune('<')
  case rune('<'):
    return rune('^')
  default:
    fmt.Println("Wrong direction provided")
    return direction
  }
}

func TurnLeft(direction rune) rune {
  switch direction {
  case rune('^'):
    return rune('<')
  case rune('<'):
    return rune('v')
  case rune('v'):
    return rune('>')
  case rune('>'):
    return rune('^')
  default:
    fmt.Println("Wrong direction provided")
    return direction
  }
}

func TranslateDirection(direction rune) (int, int) {
  switch direction {
  case rune('^'):
    return -1, 0
  case rune('<'):
    return 0, -1
  case rune('v'):
    return 1, 0
  case rune('>'):
    return 0, 1
  default:
    fmt.Println("Wrong direction provided")
    return 0, 0
  }
}
