package main

import (
	"fmt"

	"github.com/lumix103/tiny-basic/internal/parser"
)

func main() {
	input := `10 REM This is my program
20 LET X = 5 * (3 + 1)
30 IF X - 2 <= 3 THEN GOTO 50
40 PRINT "FALSE"
50 END`
	p := parser.ParserInfo{
		Input:          []byte(input),
		Line:           1,
		Col:            1,
		Is_Interactive: true,
	}
	for {
		got, next, err := parser.Parse_Statement(p)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println("tb > ", got)
		}

		if next.Is_Done {
			break
		}
		p = next
	}
}
