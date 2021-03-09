package parser

import (
    "strings"
    "bytes"
    "log"
)

const (
    UNKNOWN_TYPE int = 0
)

func newScanner(input string) *Scanner {
    log.Println("new scanner", input)
    return &Scanner{
        pos: 0,
        query: input,
    }
}

type Scanner struct {
    pos int
    query string
}

func (scanner *Scanner) nextToken() (int, string) {
    var ch rune
scan:
    ch, eof := scanner.popCh()
    if eof {
        return 0, ""
    }

    switch {
    case isSpace(ch):
        scanner.readSpace()
        goto scan

    case ch=='.' || isDigit(ch):
        scanner.unpopCh()
        return scanner.readNumber()

    case ch=='\'' || ch=='"':
        scanner.unpopCh()
        return scanner.readString()

    case ch=='+':
        return ADD, "+"

    case ch=='-':
        /*
        nextCh, _ := scanner.nextCh()
        if nextCh=='.' || isDigit(nextCh) {
            scanner.unpopCh()
            return scanner.readNumber()
        }
        */
        return DEC, "-"

    case ch=='*':
        return STAR, "*"

    case ch=='<':
        nextCh, _ := scanner.nextCh()
        if nextCh == '=' {
            scanner.popCh()
            return LTE, "<="
        } else {
            if nextCh == '>' {
                scanner.popCh()
                return NEQ, "<>"
            }
        }
        return LT, "<"

    case ch=='>':
        nextCh, _ := scanner.nextCh()
        if nextCh == '=' {
            scanner.popCh()
            return GTE, ">="
        }
        return GT, ">"

    case ch=='=':
        return EQ, "="

    case ch=='(':
        return LEFTC, "("

    case ch==')':
        return RIGHTC, ")"

    case ch==';':
        return SEMICOLON, ";"

    case ch==',':
        return COMMA, ","

    default:
        scanner.unpopCh()
        return scanner.readToken()
    }

}

func (scanner *Scanner) popCh() (rune, bool) {
    if scanner.pos < len(scanner.query) {
        ch := rune(scanner.query[scanner.pos])
        scanner.pos++
        return ch, false
    }
    return rune(' '), true
}

func (scanner *Scanner) nextCh() (rune, bool) {
    if scanner.pos < len(scanner.query) {
        return rune(scanner.query[scanner.pos]), false
    }
    return rune(' '), true
}

func (scanner *Scanner) unpopCh() {
    scanner.pos --
}

func (scanner *Scanner) readSpace() {
    for scanner.pos < len(scanner.query) {
        ch, eof := scanner.nextCh()
        if !eof && isSpace(ch) {
            scanner.popCh()
        } else {
            break
        }
    }
}

func (scanner *Scanner) readNumber() (int, string) {
    typ := INTEGER
    var buf bytes.Buffer

    ch, eof := scanner.nextCh()
    if !eof && ch=='-' {
        scanner.popCh()
        buf.WriteRune(ch)
    }
    for scanner.pos < len(scanner.query) {
        ch, eof := scanner.nextCh()
        if eof || isSpace(ch) {
            break
        }
        if !eof && !isDigit(ch) && ch!='.' {
            break
        }
        if ch == '.' {
            typ = FLOAT
        }
        scanner.popCh()
        buf.WriteRune(ch)
    }

    return typ, string(buf.Bytes())
}

func (scanner *Scanner) readString() (int, string) {
    ch0, eof := scanner.nextCh()
    if eof || (ch0!='\'' && ch0!='"') {
        return UNKNOWN_TYPE, ""
    }

    var buf bytes.Buffer
    scanner.popCh()
    for scanner.pos < len(scanner.query) {
        ch, eof := scanner.popCh()
        if ch == ch0 {
            return STRING, string(buf.Bytes())
        } else {
            if eof {
                return UNKNOWN_TYPE, string(buf.Bytes())
            } else {
                buf.WriteRune(ch)
            }
        }
    }

    return UNKNOWN_TYPE, string(buf.Bytes())
}

func (scanner *Scanner) readToken() (int, string) {
    var buf bytes.Buffer
    for scanner.pos < len(scanner.query) {
        ch, eof := scanner.popCh()
        if eof || isSpace(ch) {
            break
        }
        if !eof && !isLetter(ch) && !isDigit(ch) && ch!='_' {
            scanner.unpopCh()
            break
        }
        buf.WriteRune(ch)
    }

    str := string(buf.Bytes())
    upper := strings.ToUpper(str)
    switch upper {
    case "SELECT":
        return SELECT, upper
    case "FROM":
        return FROM, upper
    case "WHERE":
        return WHERE, upper
    case "AS":
        return AS, upper
    case "GROUP":
        return GROUP, upper
    case "BY":
        return BY, upper
    case "ORDER":
        return ORDER, upper
    case "LIMIT":
        return LIMIT, upper
    case "CREATE":
        return CREATE, upper
    case "UPDATE":
        return UPDATE, upper
    case "INSERT":
        return INSERT, upper
    case "INTO":
        return INTO, upper
    case "DELETE":
        return DELETE, upper
    case "AND":
        return AND, upper
    case "OR":
        return OR, upper
    }

    return NAME, str
}

func isSpace(ch rune) bool {
    return ch==' ' || ch =='\t' || ch == '\n'
}

func isDigit(ch rune) bool {
    return ch>='0' && ch<='9'
}

func isLetter(ch rune) bool {
    return ('a'<=ch&&ch<='z') || ('A'<=ch&&ch<='Z')
}


