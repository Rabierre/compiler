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
	IntLit
	DoubleLit
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
	IdentType
	CommentType
	ReturnType
	SemiColType
	CommaType
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
	case CmtSlash:
		return CommentType
	case Return:
		return ReturnType
	case SemiCol:
		return SemiColType
	case CommaLit:
		return CommaType
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
	case IdentType:
		return "ident"
	case CommentType:
		return "comment"
	case ReturnType:
		return "return"
	case SemiColType:
		return ";"
	case CommaType:
		return "comma"
	case EOFType:
		return "EOF"
	case IntLit:
		return "IntLit"
	case DoubleLit:
		return "DoubleLit"
	}
	return "NONE"
}

type CharType int

const (
	LETTER CharType = iota
	DIGIT
	DOT
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	SEMICOLON
	COMMA
	OTHER
)

func Kind(ch string) CharType {
	if ch == "" {
		return OTHER
	}

	c := []rune(ch)[0]
	switch {
	case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		return LETTER
	case '0' <= c && c <= '9':
		return DIGIT
	case c == '.':
		return DOT
	case c == '{':
		return LBRACE
	case c == '}':
		return RBRACE
	case c == '(':
		return LPAREN
	case c == ')':
		return RPAREN
	case c == ';':
		return SEMICOLON
	case c == ',':
		return COMMA
	default:
		return OTHER
	}
}

const (
	LowestPriority  = 0 // non-operators
	UnaryPriority   = 6
	HighestPriority = 7
)

func (t Token) Priority() int {
	switch t.kind {
	// case LOR:
	// 	return 1
	// case LAND:
	// 	return 2
	case EqType, NotEqType, LessType, LessEqType, GreatType, GreatEqType:
		return 3
	case PlusType, MinusType:
		return 4
	case MultiType, DivideType:
		return 5
	}
	return LowestPriority
}
