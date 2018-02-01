package main

// Visit https://norasandler.com/2017/11/29/Write-a-Compiler.html

import (
	"log"

	"github.com/jimmyjames85/gohandlebars/parser"
)

// An artificial input source.
const input = "1234 5678 1234567901234567890"

func main() {

	_ = []byte(`{ int foo = 23234;

// 3this is a comment
}
int a = 3; // here is another comment

/* here is multi-line
   comment

*/
int foo = 34;
/* yet another line comment */ int a = 34;

puts a;

int jim = "hello \"jim\" how are you : )";
this is an identifier but this is unknown

int main(char *args){
  fmt.Printf("hello world");
}


=
==
===x



`)

	_ = []byte(`
/* this is a comment
multi line */

if (a <= b) {
    c = 2;
    return c;
} else if (if b == 4) {
    c = 3;
}

float f = 34.23443254;

`)

	return2 := []byte(`

	return 334;

	/*
	  This is a multiline comment
	*/

	int main() {
	    // line comment
	    return 2;
	}

`)

	err := parser.ParseReturn2(return2)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

}
