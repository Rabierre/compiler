package main

type TokenType int

const (
	IfType TokenType = iota
	ElseType
	ForType
	FuncType
	LParenType
	RParenType
	LBraceType
	RBraceType
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
	GreatEqType
	EqType
	NotEqType
	SpaceType
	IdentType
	CommentType
	ReturnType
	EOFType
)

func KeywordType(token string) TokenType {
	switch token {
	case If:
		return IfType
	case Else:
		return ElseType
	case For:
		return ForType
	case Func:
		return FuncType
	case LParen:
		return LParenType
	case RParen:
		return RParenType
	case LBrace:
		return LBraceType
	case RBrace:
		return RBraceType
	case Int:
		return IntType
	case Double:
		return DoubleType
	case Plus:
		return PlusType
	case Minus:
		return MinusType
	case Multi:
		return MultiType
	case Divide:
		return DivideType
	case Assign:
		return AssignType
	case Less:
		return LessType
	case Great:
		return GreatType
	case LessEq:
		return LessEqType
	case GreateEq:
		return GreatEqType
	case Equal:
		return EqType
	case NotEq:
		return NotEqType
	case Space:
		return SpaceType
	case CmtSlash:
		return CommentType
	case Return:
		return ReturnType
	}
	return IdentType
}

type Token struct {
	val  string
	kind TokenType
}

func (t TokenType) String() string {
	switch t {
	case IfType:
		return "if"
	case ElseType:
		return "else"
	case ForType:
		return "for"
	case FuncType:
		return "func"
	case LParenType:
		return "lParen"
	case RParenType:
		return "rParen"
	case LBraceType:
		return "lBrace"
	case RBraceType:
		return "rBrace"
	case IntType:
		return "int"
	case DoubleType:
		return "double"
	case PlusType:
		return "plus"
	case MinusType:
		return "minus"
	case MultiType:
		return "multi"
	case DivideType:
		return "divide"
	case AssignType:
		return "assign"
	case LessType:
		return "less"
	case GreatType:
		return "greate"
	case LessEqType:
		return "lessEq"
	case GreatEqType:
		return "greatEq"
	case EqType:
		return "eq"
	case NotEqType:
		return "notEq"
	case SpaceType:
		return "space"
	case IdentType:
		return "ident"
	case CommentType:
		return "comment"
	case ReturnType:
		return "return"
	case EOFType:
		return "EOF"
	}
	return "WRONG TYPE"
}

type CharType int

const (
	LETTER CharType = iota
	DIGIT
	DOUBLE_QUOTE
	LBRACE
	RBRACE
	OTHER
)

func Kind(ch string) CharType {
	switch ch {
	case "a", "b", "c", "d", "e",
		"f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o",
		"p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z":
		return LETTER
	case "0", "1", "2", "3", "4",
		"5", "6", "7", "8", "9":
		return DIGIT
	case ".":
		return DOUBLE_QUOTE
	case "{":
		return LBRACE
	case "}":
		return RBRACE
	default:
		return OTHER
	}
}
