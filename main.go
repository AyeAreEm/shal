package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    state := State{
        Ans: 0,
        Show: false,
        Err: nil,
        Assigns: make(map[string]float64),
        FnDecls: make(map[string]FnDecl),
    }

    for true {
        fmt.Print("> ")
        input, _ := reader.ReadString('\n')

        lexer := Lex(input)
        stmnt := Parse(lexer.Tokens)
        state.Eval(stmnt)

        if state.Err != nil {
            fmt.Println("error:", state.Err.Error())
        } else if state.Show {
            fmt.Println(state.Ans)
        }
    }
}
