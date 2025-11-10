package main

import (
	"errors"
	"fmt"
	"os"
)

type State struct {
    Val     float64
    Show   bool
    Err     error
    Assigns map[string]float64
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
    case ExprBinary:
        return state.binary(expr.Val.(ExBinary))
    case ExprUnary:
        return state.unary(expr.Val.(ExUnary))
    case ExprGroup:
        return state.expr(expr.Val.(Expr))
    }

    return 0
}

func (state *State) assign(assign Assign) (string, float64) {
    ident := assign.Ident
    expr := state.expr(assign.Val)
    return ident, expr
}

func (state *State) Eval(stmnt Stmnt) {
    switch stmnt.Kind {
    case StmntExprs:
        state.Val = state.expr(stmnt.Val.(Expr))
        state.Show = true
    case StmntClear:
        fmt.Print("\033[H\033[2J")
        state.Show = false
    case StmntExit:
        os.Exit(0)
    case StmntAssign:
        ident, expr := state.assign(stmnt.Val.(Assign))
        state.Assigns[ident] = expr
        state.Show = false
    }
}
