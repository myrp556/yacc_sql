package parser

import (
    "log"
    "strconv"
)

type lex struct {
    input string
    scanner *Scanner
    Query Query
}

func NewLex(input string) *lex {
    return &lex{
        input: input,
        scanner: newScanner(input),
    }
}

func (l *lex) Lex(lval *yySymType) int {
    typ, str := l.scanner.nextToken()
    log.Println(typ, str)

    switch typ {
    case 0:
        return 0
    case INTEGER:
        lval.int64, _ = strconv.ParseInt(str, 10, 64)
    case FLOAT:
        lval.float64, _ = strconv.ParseFloat(str, 64)
    //case STRING, NAME:
    //    lex.str = str
    case EQ, NEQ, LT, LTE, GT, GTE, AND, OR, ADD, DEC:
        lval.int = typ
    }
    lval.str = str

    return typ
}

func (l *lex) Error(err string) {
    log.Fatal(err)
}

func Parse(l *lex) Statements {
    yyParse(l)
    return l.Query.Statements
}
