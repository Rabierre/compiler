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

func KeywordType(token string) (kind TokenType) {
	switch token {
	case If:
		kind = IfType
	case Else:
		kind = ElseType
	case For:
		kind = ForType
	case Func:
		kind = FuncType
	case LParen:
		kind = LParenType
	case RParen:
		kind = RParenType
	case Int:
		kind = IntType
	case Double:
		kind = DoubleType
	case Plus:
		kind = PlusType
	case Minus:
		kind = MinusType
	case Multi:
		kind = MultiType
	case Divide:
		kind = DivideType
	case Assign:
		kind = AssignType
	case Less:
		kind = LessType
	case Great:
		kind = GreatType
	case LessEq:
		kind = LessEqType
	case GreateEq:
		kind = GreateEqType
	case Equal:
		kind = EqualType
	case NotEq:
		kind = NotEqType
	case Space:
		kind = SpaceType
	default:
		kind = IdentType
	}
	return kind
}

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
