package main

import (
	"errors"
	"fmt"
	"os"
)

type State struct {
    Ans     float64
    Show    bool
    Err     error
    Assigns map[string]float64
    FnDecls map[string]FnDecl
}

func (state *State) binary(expr ExBinary) float64 {
    switch (expr.Kind) {
    case BinaryPlus:
        x := state.expr(expr.Left)
        y := state.expr(expr.Right)
        return x + y
    case BinaryMinus:
        x := state.expr(expr.Left)
        y := state.expr(expr.Right)
        return x - y
    case BinaryMultiply:
        x := state.expr(expr.Left)
        y := state.expr(expr.Right)
        return x * y
    case BinaryDivide:
        x := state.expr(expr.Left)
        y := state.expr(expr.Right)
        return x / y
    }

    return 0
}

func (state *State) unary(expr ExUnary) float64 {
    switch expr.Kind {
    case UnaryNegate:
        x := state.expr(expr.Val)
        return -x
    }

    return 0
}

func (state *State) fncall(fncall ExFnCall) float64 {
    fndecl, ok := state.FnDecls[fncall.Ident]
    if !ok {
        state.Err = errors.New("function declaration not found")
        return 0
    }
    if len(fndecl.Args) != len(fncall.Args) {
        state.Err = errors.New("function call and declaration have different arguments")
        return 0
    }

    for i, darg := range fndecl.Args {
        carg := fncall.Args[i]
        state.Assigns[darg.Val.(string)] = state.expr(carg)
    }

    result := state.expr(fndecl.Body)

    for _, darg := range fndecl.Args {
        delete(state.Assigns, darg.Val.(string))
    }

    return result
}

func (state *State) expr(expr Expr) float64 {
    switch expr.Kind {
    case ExprNumlit:
        return expr.Val.(float64)
    case ExprIdent:
        val, ok := state.Assigns[expr.Val.(string)]
        if !ok {
            state.Err = errors.New("identifier not found")
            return 0
        }
        return val
    case ExprAns:
        return state.Ans
    case ExprFnCall:
        return state.fncall(expr.Val.(ExFnCall))
    case ExprBinary:
        return state.binary(expr.Val.(ExBinary))
    case ExprUnary:
        return state.unary(expr.Val.(ExUnary))
    case ExprGroup:
        return state.expr(expr.Val.(Expr))
    }

    return 0
}

func (state *State) assign(assign Assign) {
    ident := assign.Ident
    expr := state.expr(assign.Val)
    state.Assigns[ident] = expr
}

func (state *State) fndecl(fndecl FnDecl) {
    ident := fndecl.Ident
    state.FnDecls[ident] = fndecl
}

func (state *State) Eval(stmnt Stmnt) {
    switch stmnt.Kind {
    case StmntExprs:
        state.Ans = state.expr(stmnt.Val.(Expr))
        state.Show = true
    case StmntClear:
        fmt.Print("\033[H\033[2J")
        state.Show = false
    case StmntExit:
        os.Exit(0)
    case StmntAssign:
        state.assign(stmnt.Val.(Assign))
        state.Show = false
    case StmntFnDecl:
        state.fndecl(stmnt.Val.(FnDecl))
        state.Show = false
    }
}
