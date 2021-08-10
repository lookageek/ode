package parser

import (
	"lookageek.com/ode/ast"
	"lookageek.com/ode/lexer"
	"lookageek.com/ode/token"
)

// Parser has a reference to the lexer
// curToken and peekToken holds the current token and the
// next token
// sometimes we will have to look two tokens at once to determine
// the proper parsing, for example with statement `5;` we need to see 5 and
// the semicolon to determine if `5` does not start off an Expression
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
