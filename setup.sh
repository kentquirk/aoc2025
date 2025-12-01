#! /bin/bash

if [ -z $1 ]; then
  echo "Day number required"
  exit 1
fi

if [ ! -d _template_$2 ]; then
  echo "Language required (go, py)"
  exit 1
fi

cp -r _template_$2 day$1_$2
cd day$1_$2
if [ -f go.mod ]; then
  sed s/XXX/day$1/ <../_template_go/go.mod >go.mod
  go mod tidy
fi
code .
