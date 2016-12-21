package main

type TokenType int

const (
	IfType TokenType = iota
	ElseType
	ForType
	FuncType
	LParenType
	RParenType
	IntType
	DoubleType
	PlusType
	MinusType
	MultiType
	DivideType
	AssignType
	LessType
	GreatType
	LessEqType
	GreateEqType
	EqualType
	NotEqType
	SpaceType
	IdentType
	NumberType
	EOFType
)

type Token struct {
	val  string
	kind TokenType
}

type CharType int

const (
	LETTER CharType = iota
	DIGIT
	DOUBLE_QUOTE
	OTHER
)
