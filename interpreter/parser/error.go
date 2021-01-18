package parser

import (
	"fmt"
)

type Error struct {
	Msg  string
	Line int // token pos line in source code
	Pos  int // token pos of source code
}

func (err Error) String() string {
	return fmt.Sprintf("%s at Line:%v Pos:%v",
		err.Msg, err.Line, err.Pos,
	)
}

func (p *Parser) AddError(format string, args ...interface{}) {
	err := Error{
		Msg:  fmt.Sprintf(format, args...),
		Line: p.l.CurrentLine(),
		Pos:  p.l.CurrentPosInLine(),
	}
	p.errors = append(p.errors, err)
}

// Errors return stored errors
func (p *Parser) Errors() []Error {
	return p.errors
}
