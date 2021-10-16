package models

import (
	"fmt"
)

// This is the type used for any standard ID field — namely, `serial NOT NULL PRIMARY KEY`
type IDType = int64

type SelectQuery struct {
	query      string
	args       []interface{}
	Limit      int
	Offset     int
	OrderBy    string
	Descending bool
}

func (p *SelectQuery) SetQuery(query string, args ...interface{}) {
	p.query = query
	p.args = args
}

func (p SelectQuery) Query() string {
	orderByClause := ""
	if p.OrderBy != "" {
		orderByClause += ` ORDER BY ` + p.OrderBy
		if p.Descending {
			orderByClause += ` DESC `
		}
	}
	if p.Limit > 0 {
		return p.query + orderByClause + fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(p.args)+1, len(p.args)+2)
	} else {
		return p.query + orderByClause
	}
}

func (p SelectQuery) Args() []interface{} {
	if p.Limit > 0 {
		return append(p.args, p.Limit, p.Offset)
	} else {
		return p.args
	}
}
