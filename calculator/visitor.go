package calculator

type Visitor interface {
	assignStmt(stmt *assignStmt)
	//cmpExpr(expr *cmpExpr)
	//arithExpr(expr *arithExpr)
	//termExpr(expr *termExpr)
	//powerExpr(expr *powerExpr)
	binaryExpr(expr *binaryExpr)
	subExpr(expr *subExpr)
	unaryExpr(expr *unaryExpr)
	number(num *number)
	identifier(idn *identifier)
}
