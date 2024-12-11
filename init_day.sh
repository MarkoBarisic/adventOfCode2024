#!/usr/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
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
DAY_DIR="$SCRIPT_DIR/day${day}"
mkdir -p $DAY_DIR/inputs
touch $DAY_DIR/inputs/real.txt
touch $DAY_DIR/inputs/test.txt
cat $SCRIPT_DIR/templates/dayX.go > $DAY_DIR/day${day}.go
cat $SCRIPT_DIR/templates/dayX_README.md > $DAY_DIR/README.md
cd $DAY_DIR
go mod init adventOfCode2024/day${day}
go mod edit -replace adventOfCode2024/util=../util
go mod tidy
cd $SCRIPT_DIR
