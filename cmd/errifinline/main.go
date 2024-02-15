package main

import (
	"log"

	"github.com/bastianccm/errifinline"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	analyzer, err := errifinline.NewAnalyzer()
	if err != nil {
		log.Fatal(err)
	}

	singlechecker.Main(analyzer)
}
