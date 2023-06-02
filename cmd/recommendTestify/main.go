package main

import (
	"github.com/ry023/recommendTestify"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(recommendTestify.Analyzer) }
