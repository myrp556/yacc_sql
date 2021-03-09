%{
package parser

import (
    //"fmt"
)

func addStatement(yylex yyLexer, stmt Statement) {
    yylex.(*lex).Query.Statements = append(yylex.(*lex).Query.Statements, stmt)
}

%}

%union {
    stmt        Statement
    stmts       Statements
    selectStmt  *SelectStatement

    str         string
    query       Query
    field       *Field
    fields      Fields
    source       *Source
    sources      Sources

    int         int
    int64       int64
    float64     float64
    bool        bool

    expr        *ExprNode
}

%token <str>    SELECT FROM WHERE AS GROUP BY ORDER LIMIT CREATE
%token <str>    UPDATE INSERT INTO DELETE
%token <str>    STAR
%token <int>    EQ NEQ LT LTE GT GTE
%token <int>    LEFTC, RIGHTC
%token <str>    NAME
%token <int64>  INTEGER
%token <float64> FLOAT
%token <str>    STRING
%token <str>    COMMA SEMICOLON

%left <int>     OR
%left <int>     AND
%left <int>     ADD DEC
%left <int>     TIME DIV MOD

%type <stmts>       QUERIES
%type <stmt>        STATEMENT
%type <selectStmt>  SELECT_STATEMENT
%type <field>       FIELD
%type <fields>      FIELDS
%type <source>      SOURCE
%type <sources>     SOURCES
%type <expr>        CONDITION WHERE_CONDITION
%type <expr>        CONDITION_VAR 
%type <int>         OPERATOR

//%start main

%%

QUERIES: 
    STATEMENT {
        addStatement(yylex, $1)
    }|
    STATEMENT SEMICOLON {
        addStatement(yylex, $1)
    }|
    STATEMENT SEMICOLON QUERIES {
        addStatement(yylex, $1)
    }

STATEMENT:
    SELECT_STATEMENT {
        $$ = $1
    }

SELECT_STATEMENT:
    SELECT FIELDS FROM SOURCES WHERE_CONDITION {
        sele := &SelectStatement{
            Fields: $2,
            Sources: $4,
            WhereCondition: $5,
        }

        $$ = sele
    }

FIELDS:
    FIELD {
        $$ = []*Field{$1}
    }|
    FIELD COMMA FIELDS {
        $$ = append($3, $1)
    }

FIELD:
    NAME {
        $$ = &Field{Name: $1}
    }|
    STAR {
        $$ = &Field{Name: $1}
    }

SOURCES:
    SOURCE {
        $$ = []*Source{$1}
    }|
    SOURCE COMMA SOURCES {
        $$ = append($3, $1)
    }

SOURCE:
     NAME {
        $$ = &Source{Name: $1}
     }

WHERE_CONDITION:
    WHERE CONDITION {
        $$ = $2
    }|
     {
        $$ = nil
    }

CONDITION:
    LEFTC CONDITION RIGHTC {
        $$ = $2
    }|
    CONDITION_VAR OPERATOR CONDITION_VAR {
        $$ = &ExprNode{Type: BinaryNode, Left: $1, Op: $2, Right: $3}
    }|
    CONDITION AND CONDITION {
        $$ = &ExprNode{Type: BinaryNode, Left: $1, Op: $2, Right: $3}
    }|
    CONDITION OR CONDITION {
        $$ = &ExprNode{Type: BinaryNode, Left: $1, Op: $2, Right: $3}
    }

OPERATOR:
    EQ {
        $$ = $1
    }|
    NEQ {
        $$ = $1
    }|
    LT {
        $$ = $1
    }|
    LTE {
        $$ = $1
    }|
    GT {
        $$ = $1
    }|
    GTE {
        $$ = $1
    }

CONDITION_VAR:
    NAME {
        $$ = &ExprNode{Type: FieldNode, Name: $1}
    }|
    STRING {
        $$ = &ExprNode{Type: StringNode, StrVal: $1}
    }|
    INTEGER {
        $$ = &ExprNode{Type: IntegerNode, IntVal: $1}
    }|
    FLOAT {
        $$ = &ExprNode{Type: FloatNode, FloVal: $1}
    }

%%
