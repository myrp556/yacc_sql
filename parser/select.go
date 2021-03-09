package parser

import (
    "log"
)

func (*SelectStatement) stmt() {}

type SelectStatement struct {
    Fields Fields
    Sources Sources
    WhereCondition *ExprNode
}

func (stmt *SelectStatement) travel() {
    str := "SELECT " + stmt.Fields.str() + " FROM " + stmt.Sources.str()
    if stmt.WhereCondition != nil {
        str = str + " WHERE " + stmt.WhereCondition.str()
    }

    log.Println(str)
}


