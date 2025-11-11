package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var keywords = []string{"clear", "exit", "let"}

type TokenKind int
const (
    TokenNone TokenKind = iota
    TokenPlus
    TokenMinus
    TokenStar
    TokenSlash
    TokenNumlit
    TokenClear
    TokenExit
    TokenLet
    TokenIdent
    TokenEqual
    TokenLeftBracket
    TokenRightBracket
    TokenComma
)

type Token struct {
    Kind TokenKind
    Val  any
}

func token_none() Token {
    return Token{
        Kind: TokenNone,
    }
}
func token_plus() Token {
    return Token{
        Kind: TokenPlus,
        Val: nil,
    }
}
func token_minus() Token {
    return Token{
        Kind: TokenMinus,
        Val: nil,
    }
}
func token_star() Token {
    return Token{
        Kind: TokenStar,
        Val: nil,
    }
}
func token_slash() Token {
    return Token{
        Kind: TokenSlash,
        Val: nil,
    }
}
func token_numlit(numlit float64) Token {
    return Token{
        Kind: TokenNumlit,
        Val: numlit,
    }
}
func token_keyword(keyword string) Token {
    switch keyword {
    case "clear":
        return Token{Kind: TokenClear}
    case "exit":
        return Token{Kind: TokenExit}
    case "let":
        return Token{Kind: TokenLet}
    }

    return token_none()
}
func token_ident(ident string) Token {
    return Token{
        Kind: TokenIdent,
        Val: ident,
    }
}
func token_equal() Token {
    return Token{
        Kind: TokenEqual,
    }
}
func token_left_bracket() Token {
    return Token{
        Kind: TokenLeftBracket,
    }
}
func token_right_bracket() Token {
    return Token{
        Kind: TokenRightBracket,
    }
}
func token_comma() Token {
    return Token{
        Kind: TokenComma,
    }
}

type Lexer struct {
    buf strings.Builder
    Tokens []Token
}

func (l *Lexer) push_rune(r rune) {
    _, err := l.buf.WriteRune(r)
    if err != nil {
        fmt.Printf("error: %v", err.Error())
        os.Exit(1)
    }
}

func (l *Lexer) push_tok(tok Token) {
    l.Tokens = append(l.Tokens, tok)
}

func (l *Lexer) resolve_buf() {
    if l.buf.Len() == 0 {
        return
    }

    buf := l.buf.String()
    float, ferr := strconv.ParseFloat(buf, 64)
    integer, ierr := strconv.ParseInt(buf, 0, 64)

    if ferr == nil {
        l.push_tok(token_numlit(float))
    } else if ierr == nil {
        l.push_tok(token_numlit(float64(integer)))
    } else if slices.Contains(keywords, buf) {
        l.push_tok(token_keyword(buf))
    } else {
        l.push_tok(token_ident(buf))
    }

    l.buf.Reset()
}

func (l *Lexer) resolve_and_push(tok Token) {
    l.resolve_buf()
    l.push_tok(tok)
}

func Lex(line string) Lexer {
    lexer := Lexer{}

    for _, ch := range line {
        switch ch {
        case '\r':
            continue
        case '\t':
            fallthrough
        case '\n':
            fallthrough
        case ' ':
            lexer.resolve_buf()

        case '+':
            lexer.resolve_and_push(token_plus())
        case '-':
            lexer.resolve_and_push(token_minus())
        case '*':
            lexer.resolve_and_push(token_star())
        case '/':
            lexer.resolve_and_push(token_slash())
        case '=':
            lexer.resolve_and_push(token_equal())
        case '(':
            lexer.resolve_and_push(token_left_bracket())
        case ')':
            lexer.resolve_and_push(token_right_bracket())
        case ',':
            lexer.resolve_and_push(token_comma())
        default:
            lexer.push_rune(ch)
        }
    }

    return lexer
}
