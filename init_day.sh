#!/usr/bin/bash
day=-1
while [[ $# -gt 0 ]]; do
  case "$1" in
    --day)
      day=$2
      shift 2
    ;;
    *)
      echo "Unknown argument ${1}"
      exit 1
    ;;
  esac
done
if [ "$day" == -1 ]; then
  echo "Argument --day is not passed"
  exit 1
fi

mkdir day${day}
cd day${day}
mkdir inputs
touch inputs/real.txt
touch inputs/test.txt
touch day${day}.go
go mod init adventOfCode2024/day${day}
cat ../dayX.go.template > day${day}.go
go mod edit -replace adventOfCode2024/util=../util
go mod tidy
cd ..
