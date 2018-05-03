package main

// Visit https://norasandler.com/2017/11/29/Write-a-Compiler.html

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jimmyjames85/gohandlebars/pkg/lexer"
	"github.com/jimmyjames85/gohandlebars/pkg/parser"
)

// An artificial input source.
const ()

type parserFunc func(srcCode []byte) ([]byte, error)

func test() {

	myExp := []byte(`2*3+(~7-34)`)
	tokens, err := lexer.Scan(myExp)
	if err != nil {
		panic(err)
	}

	expAst, err := parser.ParseExp(tokens)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n\n", expAst.ExpString())
}

func parseFunc(srcCode []byte) ([]byte, error) {
	tokens, err := lexer.Scan(srcCode)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFunction(tokens)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s\n\n", f.FuncString())), nil
}

func main() {

	curParser := parser.ParseReturn2
	curParser = parseFunc

	//////////////////////////////////////////////////////////////////////
	// Compile
	//////////////////////////////////////////////////////////////////////
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "please")
	}
	fileloc := os.Args[1]

	_, err := os.Stat(fileloc)
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(fileloc)
	name := filepath.Base(fileloc)
	if !strings.HasSuffix(name, ".c") {
		log.Fatal("expect .c file")
	}
	// rm .c suffix for output file
	name = name[:len(name)-2]

	srcCode, err := ioutil.ReadFile(fileloc)
	if err != nil {
		log.Fatal(err)
	}

	asmCode, err := curParser(srcCode)
	if err != nil {
		log.Fatal(err)
	}

	asmFilepath := filepath.Join(dir, fmt.Sprintf("%s.s", name))
	err = ioutil.WriteFile(asmFilepath, asmCode, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", asmCode)

	binFilePath := filepath.Join(dir, name)
	cmd := exec.Command("gcc", asmFilepath, "-o", binFilePath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command("rm", asmFilepath)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// go run cmd/test/main.go > ./return_2.s && gcc -m32 return_2.s -o return_2

	return

	b := []byte(``)
	err = parser.Parse(b)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}
