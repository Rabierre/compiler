BNF description for LL(>=1) grammars

```
Program ::= DeclList ?
DeclList ::= ( VarDecl | FunctionDecl ) DeclList ?
FunctionDecl ::= Type identifier "(" FieldList ? ")" CompoundStmt
FieldList ::= Field ( "," FieldList ) ?
Field ::= Type identifier
VarDecl ::= Type IdentList
Type ::= "int"
       | "double"
IdentList ::= identifier ( "=" Expr ) ? ( "," IdentList ) ?
Stmt ::= ForStmt
       | Expr
       | IfStmt
       | CompoundStmt
       | "return" Expr ?
ForStmt ::= "for" "(" OptExpr ";" OptExpr ";" OptExpr ")" CompoundStmt
OptExpr ::= Expr ?
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
Comment ::= "//" Text
Text ::= a-z
         | A-Z
         | 0-9
```