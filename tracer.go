package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

var appContext string
var tracingCall string

// go build && ./tracer ./server/serve.go,./server/sample.txt
func main() {
	changedFiles := getFilepathsFromArgs()
	for _, f := range changedFiles {
		checkFileForTracing(f)
	}
}

func getFilepathsFromArgs() []string {
	allChangedFiles := ""

	args := os.Args
	if len(args) > 1 {
		allChangedFiles = args[1]
	}

	filepaths := strings.Split(allChangedFiles, ",")
	return filepaths
}

func checkFileForTracing(filepath string) {
	if filepath == "" {
		return
	}

	fset, f := astForFile(filepath) // positions are relative to fset

	// TODO: also handle methods

	// TODO: let the tracing pattern be configurable for nested contextts
	appContext = "echo.Context"
	tracingCall = "trace"

	exportedFuncs := exportedFunctions(f)
	for _, ef := range exportedFuncs {
		if shouldHaveTracing(ef, appContext) && !hasTracing(ef, tracingCall) {
			fmt.Printf("🌈 %s:%d %s() needs tracing\n", filepath, fset.Position(ef.Pos()).Line, ef.Name)
		}
	}
}

func astForFile(filepath string) (*token.FileSet, *ast.File) {
	if filepath == "" {
		return nil, nil
	}

	src := fileToText(filepath)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	return fset, f
}

func print(fset *token.FileSet, f *ast.File) {
	ast.Print(fset, f)
}

func fileToText(filepath string) string {
	// TODO: refactor to read in chunks, maybe
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("⛈ error reading file")
	}
	return string(content)
}

func addTracing(fd *ast.FuncDecl) {
	// TODO
}

func hasTracing(fd *ast.FuncDecl, tracingCall string) bool {
	hasTracingCall := false

	// TODO strings.Contains(tracingCall, ".")
	if stmt, ok := fd.Body.List[0].(*ast.ExprStmt); ok {
		if call, ok := stmt.X.(*ast.CallExpr); ok {
			if fmt.Sprintf("%s", call.Fun) == tracingCall {
				hasTracingCall = true
			}
		}
	}

	return hasTracingCall
}

func shouldHaveTracing(ef *ast.FuncDecl, appContext string) bool {
	hasAppContext := false

	params := ef.Type.Params.List
	if len(params) != 0 {
		if strings.Contains(appContext, ".") {
			appContextPackageCall := strings.Split(appContext, ".")
			t, ok := params[0].Type.(*ast.SelectorExpr)
			if ok {
				if fmt.Sprintf("%s", t.X) == appContextPackageCall[0] && fmt.Sprintf("%s", t.Sel) == appContextPackageCall[1] {
					hasAppContext = true
				}
			}
		} else if fmt.Sprintf("%s", params[0].Type) == appContext {
			fmt.Println(ef.Name, "should have tracing")
			hasAppContext = true
		}
	}

	return hasAppContext
}

func exportedFunctions(f *ast.File) []*ast.FuncDecl {
	exportedFunctions := []*ast.FuncDecl{}

	ast.Inspect(f, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if ok && fd.Name.IsExported() {
			exportedFunctions = append(exportedFunctions, fd)
		}
		return true
	})

	return exportedFunctions
}
