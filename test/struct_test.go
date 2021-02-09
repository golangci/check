package main_test

import (
	"fmt"
	"go/types"
	"log"
	"sync"
	"testing"

	"github.com/golangci/check/cmd/structcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/loader"
)

func TestStruct(t *testing.T) {

	log.SetFlags(log.Lshortfile)
	testdata := analysistest.TestData()
	var mu sync.Mutex
	var res []goanalysis.Issue
	a := &analysis.Analyzer{
		Name: "structchexck",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			prog := goanalysis.MakeFakeLoaderProgram(pass)

			structcheckIssues := structcheck.Run(prog, true)
			if len(structcheckIssues) == 0 {
				return nil, nil
			}

			issues := make([]goanalysis.Issue, 0, len(structcheckIssues))
			for _, i := range structcheckIssues {
				pass.Reportf(i.Posx, fmt.Sprintf("%s is unused", i.FieldName))
			}

			mu.Lock()
			res = append(res, issues...)
			mu.Unlock()
			return nil, nil
		},

		Requires: []*analysis.Analyzer{
			buildssa.Analyzer,
		},
	}
	analysistest.Run(t, testdata, a, "a")
}

func MakeFakeLoaderProgram(pass *analysis.Pass) *loader.Program {
	prog := &loader.Program{
		Fset: pass.Fset,
		Created: []*loader.PackageInfo{
			{
				Pkg:                   pass.Pkg,
				Importable:            true, // not used
				TransitivelyErrorFree: true, // TODO

				Files:  pass.Files,
				Errors: nil,
				Info:   *pass.TypesInfo,
			},
		},
		AllPackages: map[*types.Package]*loader.PackageInfo{
			pass.Pkg: {
				Pkg:                   pass.Pkg,
				Importable:            true,
				TransitivelyErrorFree: true,
				Files:                 pass.Files,
				Errors:                nil,
				Info:                  *pass.TypesInfo,
			},
		},
	}
	return prog
}
