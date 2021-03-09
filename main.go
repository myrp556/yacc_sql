package main

import (
    "github.com/myrp556/yacc_parse/parser"
)

func main() {
    lex := parser.NewLex("SELECT * FROM table1, table2 WHERE a>1 and (b<2 OR c=1)")
    statements := parser.Parse(lex)
    parser.Travel(statements)
}
