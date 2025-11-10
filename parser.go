package main

import (
	"fmt"
	"os"
)

type BinaryKind int
const (
    BinaryPlus BinaryKind = iota
    BinaryMinus
    BinaryMultiply
    BinaryDivide
)

type ExBinary struct {
    Kind  BinaryKind
    Left  Expr
    Right Expr
}

func (expr ExBinary) Print() {
    switch (expr.Kind) {
    case BinaryPlus:
        expr.Left.Print()
        fmt.Printf(" + ")
        expr.Right.Print()
    case BinaryMinus:
        expr.Left.Print()
        fmt.Printf(" - ")
        expr.Right.Print()
    case BinaryMultiply:
        expr.Left.Print()
        fmt.Printf(" * ")
        expr.Right.Print()
    case BinaryDivide:
        expr.Left.Print()
        fmt.Printf(" / ")
        expr.Right.Print()
    }
}

type UnaryKind int
const (
    UnaryNegate UnaryKind = iota
)

type ExUnary struct {
    Kind UnaryKind
    Val  Expr
}

func (expr ExUnary) Print() {
    switch (expr.Kind) {
    case UnaryNegate:
        fmt.Print("-")
        expr.Val.Print()
    }
}

type ExprKind int
const (
    ExprNumlit ExprKind = iota
    ExprBinary
    ExprUnary
    ExprIdent
    ExprGroup
)

type Expr struct {
    Kind ExprKind
    Val  any
}

func (expr Expr) Print() {
    switch (expr.Kind) {
    case ExprNumlit:
        n := expr.Val.(float64)
        fmt.Printf("%v", n)
    case ExprBinary:
        e := expr.Val.(ExBinary)
        e.Print()
    case ExprUnary:
        e := expr.Val.(ExUnary)
        e.Print()
    }
}

func expr_numlit(numlit Token) Expr {
    return Expr{
        Kind: ExprNumlit,
        Val: numlit.Val,
    }
}
func expr_ident(ident Token) Expr {
    return Expr{
        Kind: ExprIdent,
        Val: ident.Val.(string),
    }
}
func expr_binary(bin ExBinary) Expr {
    return Expr{
        Kind: ExprBinary,
        Val: bin,
    }
}
func expr_unary(un ExUnary) Expr {
    return Expr{
        Kind: ExprUnary,
        Val: un,
    }
}
func expr_group(expr Expr) Expr {
    return Expr{
        Kind: ExprGroup,
        Val: expr,
    }
}

type StmntKind int
const (
    StmntExprs StmntKind = iota
    StmntClear
    StmntExit
    StmntAssign
)

type Assign struct {
    Ident string
    Val   Expr
}

type Stmnt struct {
    Kind StmntKind
    Val  any
}

func (s Stmnt) Print() {
    switch (s.Kind) {
    case StmntExprs:
        expr := s.Val.(Expr)
        expr.Print()
        fmt.Println()
    case StmntClear:
        fmt.Println("Clear")
    case StmntExit:
        fmt.Println("Exit")
    }
}
func stmnt_clear() Stmnt {
    return Stmnt{
        Kind: StmntClear,
    }
}
func stmnt_exit() Stmnt {
    return Stmnt{
        Kind: StmntExit,
    }
}
func (p *Parser) assign() Stmnt {
    // remove `let`
    p.next()

    ident := p.expect(TokenIdent)
    p.expect(TokenEqual)
    expr := p.expr()

    return Stmnt{
        Kind: StmntAssign,
        Val: Assign{
            Ident: ident.Val.(string),
            Val: expr,
        },
    }
}

type Parser struct {
    toks []Token
}

func (p *Parser) peek() Token {
    if len(p.toks) == 0 {
        return token_none()
    }

    tok := p.toks[0]
    return tok
}

func (p *Parser) next() Token {
    if len(p.toks) == 0 {
        return token_none()
    }

    tok := p.peek()
    p.toks = p.toks[1:]
    return tok
}

func (p *Parser) expect(te TokenKind) Token {
    tok := p.next()
    if te != tok.Kind {
        fmt.Printf("error: expected %v, got %v\n", te, tok.Kind)
        os.Exit(1)
    }
    return tok
}

func (p *Parser) primary() Expr {
    tok := p.peek()

    switch (tok.Kind) {
    case TokenNumlit:
        p.next()
        return expr_numlit(tok)
    case TokenIdent:
        p.next()
        return expr_ident(tok)
    case TokenLeftBracket:
        p.next()
        expr := p.expr()
        p.expect(TokenRightBracket)
        return expr_group(expr)
    default:
        fmt.Printf("error: unexpected token %v\n", tok.Kind)
        os.Exit(1)
    }

    return Expr{}
}

func (p *Parser) fn_call() Expr {
    expr := p.primary()

    return expr
}

func (p *Parser) unary() Expr {
    op := p.peek()
    if op.Kind != TokenMinus {
        return p.fn_call()
    }
    p.next()

    val := p.unary()
    if (op.Kind == TokenMinus) {
        return expr_unary(ExUnary{
            Kind: UnaryNegate,
            Val: val,
        })
    }

    return Expr{}
}

func (p *Parser) factor() Expr {
    expr := p.unary()

    for op := p.peek(); op.Kind == TokenStar || op.Kind == TokenSlash; op = p.peek() {
        p.next()

        left := expr
        right := p.unary()

        if op.Kind == TokenStar {
            return expr_binary(ExBinary{
                Kind: BinaryMultiply,
                Left: left,
                Right: right,
            })
        } else {
            return expr_binary(ExBinary{
                Kind: BinaryDivide,
                Left: left,
                Right: right,
            })
        }
    }

    return expr
}

func (p *Parser) term() Expr {
    expr := p.factor()

    for op := p.peek(); op.Kind == TokenPlus || op.Kind == TokenMinus; op = p.peek() {
        p.next()

        left := expr
        right := p.factor()

        if op.Kind == TokenPlus {
            return expr_binary(ExBinary{
                Kind: BinaryPlus,
                Left: left,
                Right: right,
            })
        } else {
            return expr_binary(ExBinary{
                Kind: BinaryMinus,
                Left: left,
                Right: right,
            })
        }
    }

    return expr
}

func (p *Parser) expr() Expr {
    return p.term()
}

func (p *Parser) exprs() Stmnt {
    return Stmnt{
        Kind: StmntExprs,
        Val: p.expr(),
    }
}

func Parse(toks []Token) Stmnt {
    parser := Parser{
        toks: toks,
    }

    tok := parser.peek()
    switch (tok.Kind) {
    case TokenClear:
        return stmnt_clear()
    case TokenExit:
        return stmnt_exit()
    case TokenLet:
        return parser.assign()
    case TokenIdent:
        fallthrough
    case TokenNumlit:
        return parser.exprs()
    case TokenLeftBracket:
        return parser.exprs()
    default:
        fmt.Printf("error: unexpected token %v\n", tok.Kind)
        os.Exit(1)
    }

    return Stmnt{}
}
