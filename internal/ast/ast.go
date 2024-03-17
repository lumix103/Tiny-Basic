package ast

import "fmt"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Line struct {
	LN   int
	Stmt Statement
}

type Program []Line

type Identifier struct {
	Name string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Name
}

type Number struct {
	Value int
}

func (n *Number) expressionNode() {}
func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string {
	return "\"" + sl.Value + "\""
}

type Unary struct {
	Op   string
	Expr Expression
}

func (u *Unary) expressionNode() {}
func (u *Unary) String() string {
	return u.Op + u.Expr.String()
}

type Infix struct {
	Left  Expression
	Op    string
	Right Expression
}

func (i *Infix) expressionNode() {}
func (i *Infix) String() string {
	return i.Left.String() + i.Op + i.Right.String()
}

type REM struct {
	Remark string
}

func (r *REM) expressionNode() {}
func (r *REM) String() string {
	return "REM " + r.Remark
}

type Let struct {
	Var   *Identifier
	Value Expression
}

func (l *Let) statementNode() {}
func (l *Let) String() string {
	return "LET " + l.Var.String() + " = " + l.Value.String()
}

type If struct {
	Condition  *Infix
	Consequent Statement
}

func (i *If) statementNode() {}
func (i *If) String() string {
	return "IF " + i.Condition.String() + " THEN " + i.Consequent.String()
}

type Return struct {
}

func (r *Return) statementNode() {}
func (r *Return) String() string {
	return "RETURN"
}

type End struct {
}

func (e *End) statementNode() {}
func (e *End) String() string {
	return "END"
}

type Gosub struct {
	Expr Expression
}

func (g *Gosub) statementNode() {}
func (g *Gosub) String() string {
	return "GOSUB " + g.Expr.String()
}

type Goto struct {
	Expr Expression
}

func (g *Goto) statementNode() {}
func (g *Goto) String() string {
	return "GOTO " + g.Expr.String()
}

type Print struct {
	Values []Expression
}

func (p *Print) statementNode() {}
func (p *Print) String() string {
	str := "PRINT"
	for _, val := range p.Values {
		str += " " + val.String()
	}
	return str
}
