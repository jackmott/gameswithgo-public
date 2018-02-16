package apt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type tokenType int

const (
	openParen tokenType = iota
	closeParen
	op
	constant
)

type token struct {
	typ   tokenType
	value string
}

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

func parse(tokens chan token) Node {
	for {
		token, ok := <-tokens
		if !ok {
			panic("no more tokens")
		}
		fmt.Print(token.value, ",")
	}
	return nil
}

const eof rune = -1

type stateFunc func(*lexer) stateFunc

func BeginLexing(s string) Node {
	l := &lexer{input: s, tokens: make(chan token, 100)}
	go l.run()
	return parse(l.tokens)
}

func (l *lexer) run() {
	for state := determineToken; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func determineToken(l *lexer) stateFunc {
	for {
		switch r := l.next(); {
		case isWhiteSpace(r):
			l.ignore()
		case r == '(':
			l.emit(openParen)
		case r == ')':
			l.emit(closeParen)
		case isStartOfNumber(r):
			return lexNumber
		case r == eof:
			return nil
		default:
			return lexOp
		}
	}
}

func lexOp(l *lexer) stateFunc {
	l.acceptRun("+-/*abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	l.emit(op)
	return determineToken
}

func lexNumber(l *lexer) stateFunc {
	l.accept("+-.")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	/// this would accept ".234.234" as number -> we will catch this in the parser

	if l.input[l.start:l.pos] == "-" {
		l.emit(op)
	} else {
		l.emit(constant)
	}

	return determineToken
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t' || r == '\r'
}

func isStartOfNumber(r rune) bool {
	return (r >= '0' && r <= '9') || r == '-' || r == '+' || r == '.'
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) peek() (r rune) {
	r, _ = utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}
