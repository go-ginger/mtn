package mtn

import (
	"fmt"
	dl "github.com/go-ginger/dl/query"
	"github.com/go-ginger/models"
	gm "github.com/go-ginger/models"
	"strings"
)

type Parser struct {
	dl.IParser

	conditions map[string]string
	operators  map[string]string
}

func (p *Parser) getKey(key string, params map[string]interface{}) (result string) {
	ind := 1
	k2 := key
	for {
		if _, ok := params[k2]; !ok {
			break
		}
		ind++
		k2 = key + fmt.Sprintf("%s_%d", key, ind)
	}
	result = k2
	return
}

func (p *Parser) iterate(data map[string]interface{}, existingParams map[string]interface{},
	temp *string, extra ...interface{}) (queryItems []string, params map[string]interface{}) {
	queryItems = []string{}
	params = existingParams
	if params == nil {
		params = map[string]interface{}{}
	}
	var prefix interface{}
	if extra != nil && len(extra) > 0 {
		prefix = extra[0]
	}
	for k, v := range data {
		paramKey := p.getKey(k, params)
		if prefix != nil {
			k = fmt.Sprintf("%v%s", prefix, k)
		}
		op := p.generateOperator(k)
		if op != nil {
			d2, success := v.([]interface{})
			if success {
				var opQueryItems []string
				for _, d := range d2 {
					d3, success := d.(map[string]interface{})
					if success {
						q, p := p.iterate(d3, params, nil)
						for pk, pv := range p {
							params[pk] = pv
						}
						opQueryItems = append(opQueryItems, q...)
					}
				}
				query := "(" + strings.Join(opQueryItems, ") "+*op+" (") + ")"
				queryItems = append(queryItems, query)
			}
		} else {
			var condition *string
			if temp != nil {
				condition = p.generateCondition(k, *temp, "$"+paramKey)
				if condition != nil {
					queryItems = append(queryItems, *condition)
					params[paramKey] = v
				}
			}
			if condition == nil {
				if iv, ok := v.(map[string]interface{}); ok {
					q, p := p.iterate(iv, params, &k)
					for pk, pv := range p {
						params[pk] = pv
					}
					queryItems = append(queryItems, q...)
				} else if iv, ok := v.(models.Filters); ok {
					q, p := p.iterate(iv, params, &k)
					for pk, pv := range p {
						params[pk] = pv
					}
					queryItems = append(queryItems, q...)

				} else {
					queryItems = append(queryItems, k+"=$"+paramKey)
					params[paramKey] = v
				}
			}
		}
	}
	if queryItems != nil {
		queryItems = []string{"(" + strings.Join(queryItems, ") AND (") + ")"}
	}
	return queryItems, params
}

func (p *Parser) Parse(request gm.IRequest, extra ...interface{}) (result dl.IParseResult) {
	req := request.GetBaseRequest()
	parts, params := p.iterate(*req.Filters, nil, nil, extra...)
	var q string
	if parts != nil && len(parts) > 0 {
		if len(parts) > 1 {
			q = "(" + strings.Join(parts, ") AND (") + ")"
		} else {
			q = parts[0]
		}
	}
	result = &ParseResult{
		Query:  q,
		Params: params,
	}
	return
}
