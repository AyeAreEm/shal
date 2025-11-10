package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    state := State{
        Val: 0,
        Show: false,
        Err: nil,
        Assigns: make(map[string]float64),
    }

    for true {
        fmt.Print("> ")
        input, _ := reader.ReadString('\n')

        lexer := Lex(input)
        stmnt := Parse(lexer.Tokens)
        state.Eval(stmnt)

        if state.Show {
            fmt.Println(state.Val)
        }
    }
}
