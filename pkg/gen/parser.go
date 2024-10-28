package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

type TaskMetadata struct {
	Name        string
	PayloadType string
}

const asyncMethodName = "Async"

// parseFile parses a go file and returns ???
func parseFile(fpath string) TaskMetadata {
	file, err := os.Open(fpath)
	if err != nil {
		fmt.Println(os.Getwd())
		log.Fatal(fmt.Sprintf("cannot open file %s with error %v ", fpath, err))
	}
	defer file.Close()

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fpath, file, parser.AllErrors)
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot parse file %s with error %v ", fpath, err))
	}

	var structName string
	var payloadType string
	for _, decl := range node.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						structName = typeSpec.Name.Name
					}
				}
			}
		case *ast.FuncDecl:
			if decl.Name.Name == asyncMethodName {
				params := decl.Type.Params.List
				if len(params) <= 1 {
					log.Fatal("Async method must have at least one parameter")
				}
				payloadParam := params[1]
				payloadType = payloadParam.Type.(*ast.Ident).Name
			}
		}
	}

	return TaskMetadata{
		Name:        structName,
		PayloadType: payloadType,
	}
}

func findFiles(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}
	return matches
}

func ParseFiles(pattern string) []TaskMetadata {
	fnames := findFiles(pattern)
	fmt.Println("fnames: ", fnames)
	var taskMetas []TaskMetadata
	for _, fname := range fnames {
		taskMeta := parseFile(fname)
		if taskMeta.Name == "" {
			log.Fatal("invalid task file:", fname)
		}
		taskMetas = append(taskMetas, taskMeta)
	}
	return taskMetas
}
