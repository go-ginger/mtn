package mtn

import dl "github.com/go-ginger/dl/query"

type ParseResult struct {
	dl.IParseResult

	Query  string
	Params map[string]interface{}
}

func (r *ParseResult) GetQuery() (query interface{}) {
	return r.Query
}

func (r *ParseResult) GetParams() (params interface{}) {
	return r.Params
}
