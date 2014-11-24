package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strconv"
	"strings"
)

const (
	CCYFORMAT_GOFILE = "ccyfmts.go"
)

var formatFile string
var targetDir string

func main() {

	// read input params
	flag.StringVar(&formatFile, "filename", "ccyformats.json", "The filename of the currencies in json.")
	flag.StringVar(&formatFile, "f", "ccyformats.json", "The filename of the currencies in json. (shorthand)")
	flag.StringVar(&targetDir, "targetdir", ".", "The directoryname of the generated gofile. Put it into the ccyfmt package src and recompile.")
	flag.StringVar(&targetDir, "td", "..", "The directoryname of the generated gofile. (shorthand)")
	flag.Parse()

	fmt.Println("Source file")
	fmt.Println(formatFile)

	// open formats file
	ccyformatsFile, err := os.Open(formatFile)
	fi, err := ccyformatsFile.Stat()
	ccyformatsFileSize := fi.Size()
	if err != nil {
		fmt.Println("opening formats file: " + err.Error())
	}

	// close file after running
	defer func() {
		if err := ccyformatsFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// read content
	contentBuf := make([]byte, ccyformatsFileSize)
	ccyformatsFile.Read(contentBuf)

	// is in json format
	if !isJsonFormat(contentBuf) {
		fmt.Println("input data format malformed")
		return
	}

	targetDir = strings.TrimRight(targetDir, string(os.PathSeparator))
	fileinfo, err := os.Stat(targetDir)
	if err != nil || !fileinfo.IsDir() {
		fmt.Println("target parameter is wrong: " + err.Error())
		return
	}

	targetFile := targetDir + string(os.PathSeparator) + CCYFORMAT_GOFILE
	// write the generated file into a golang file
	//err = ioutil.WriteFile(targetFile, []byte(src), 0655)

	// create go source code in
	f := createASTFile("CCYFORMATS", strconv.Quote(string(contentBuf)), "ccyfmt")
	// Create empty ast fileset
	fset := token.NewFileSet()
	// instantiate a writer
	file, err := os.Create(targetFile)
	if err != nil {
		fmt.Println("create file failed " + err.Error())
		return
	}
	// write out file
	err = format.Node(file, fset, f)
	if err != nil {
		fmt.Println("writing file failed " + err.Error())
		return
	}

	// show status
	fmt.Print("generated file placed in: ")
	fmt.Println(targetFile)
	fmt.Println("if not already in there, put it into the ccyfmt src package - recompile and use.")

}

// create the ast file to generate the go file from
func createASTFile(constantName string, constantValue string, packageName string) *ast.File {

	// constant name
	constantIdent := ast.NewIdent(constantName)
	var identList = []*ast.Ident{constantIdent}

	// value spec containing the constant
	val := ast.ValueSpec{
		Doc:   nil,
		Names: identList,
		Type:  nil,
		Values: []ast.Expr{
			&ast.BasicLit{
				ValuePos: 0,
				Kind:     token.STRING,
				Value:    constantValue,
			},
		},
		Comment: nil,
	}
	// slice of valueSpecs
	vall := []ast.Spec{&val}

	// add a genDec for the CONST
	genDecl := ast.GenDecl{
		nil,
		0,
		token.CONST,
		1,
		vall,
		1,
	}
	var decls = []ast.Decl{&genDecl}

	// go source file
	fIdent := ast.NewIdent(packageName)
	f := &ast.File{
		nil,
		1,
		fIdent,
		decls,
		nil,
		nil,
		nil,
		nil,
	}

	return f
}

// check input format
func isJsonFormat(s []byte) bool {
	var v interface{}
	err := json.Unmarshal(s, &v)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
