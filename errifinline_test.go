package errifinline

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	a, err := NewAnalyzer()
	if err != nil {
		t.Fatal(err)
	}

	analysistest.Run(t, analysistest.TestData(), a, "basic")
}
