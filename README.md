BNF description for LL(>=1) grammars

```
Program ::= DeclList ?
DeclList ::= ( VarDecl | FunctionDecl ) DeclList ?
FunctionDecl ::= Type identifier "(" ArgList ? ")" CompoundStmt
ArgList ::= Arg ( "," ArgList ) ?
Arg ::= Type identifier
VarDecl ::= Type IdentList
Type ::= "int"
       | "float"
IdentList ::= identifier ( "=" Expr ) ? ( "," IdentList ) ?
Stmt ::= ForStmt
       | WhileStmt
       | Expr
       | IfStmt
       | CompoundStmt
       | "return" Expr ?
ForStmt ::= "for" "(" Expr ";" OptExpr ";" OptExpr ")" CompoundStmt
OptExpr ::= Expr ?
WhileStmt ::= "while" "(" Expr ")" CompoundStmt
IfStmt ::= "if" "(" Expr ")" CompoundStmt ElsePart
ElsePart ::= ( "else" CompoundStmt ) ?
CompoundStmt ::= "{" VarDeclList ? StmtList ? "}"
VarDeclList ::= VarDecl VarDeclList ?
StmtList ::= Stmt StmtList ?
Expr ::= identifier "=" Expr
       | Rvalue
Rvalue ::= Mag ( Compare Rvalue ) ?
Compare ::= "=="
          | "<"
          | ">"
          | "<="
          | ">="
          | "!="
Mag ::= Term ( AddSub Mag ) ?
AddSub ::= "+"
         | "-"
Term ::= Factor ( MulDiv Term ) ?
MulDiv ::= "*"
         | "/"
Factor ::= "(" Expr ")"
         | AddSub Factor
         | identifier "(" ExprList ? ")"
         | identifier
         | number
         | string
ExprList ::= Expr ( "," ExprList ) ?
```