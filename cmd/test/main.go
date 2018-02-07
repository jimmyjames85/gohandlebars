package main

// Visit https://norasandler.com/2017/11/29/Write-a-Compiler.html

import (
	"log"

	"github.com/jimmyjames85/gohandlebars/parser"
)

// An artificial input source.
const ()

func parse2() {

	return2 := []byte(`
		/*
		  This is a multiline comment
		*/

//		return 4;

	// what is this
	/* 		return 2342; // <- should not get parsed
this is a multi line comment */


//		return 2; //should get parsed


		int main() {
		    // line comment
		    return 32;
		}
	`)

	err := parser.ParseReturn2(return2)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}

func main() {

	// go run cmd/test/main.go > ./return_2.s && gcc -m32 return_2.s -o return_2

	parse2()
	return

	b := []byte(``)
	err := parser.Parse(b)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}
