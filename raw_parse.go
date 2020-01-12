package mtn

import (
	"fmt"
	"strings"
)

func (p *Parser) generateRawCondition(op, key string, value interface{}) *string {
	c, exists := p.conditions[op]
	if exists {
		query := key
		query += c
		query += value.(string)
		return &query
	}
	return nil
}

func (p *Parser) iterateRaw(data map[string]interface{}, temp *string) []string {
	var queryItems []string
	for k, v := range data {
		op := p.generateOperator(k)
		if op != nil {
			queryItems = append(queryItems, *op)
		} else if temp != nil {
			condition := p.generateRawCondition(k, *temp, v)
			if condition != nil {
				queryItems = append(queryItems, *condition)
			}
		} else {
			switch v.(type) {
			case string:
				queryItems = append(queryItems, k+" = '"+v.(string)+"'")
			default:
				queryItems = append(queryItems, k+" = "+fmt.Sprintf("%v", v))
			}
		}
	}
	return queryItems
}

func (p *Parser) RawParse(query interface{}) interface{} {
	parts := p.iterateRaw(query.(map[string]interface{}), nil)
	result := "(" + strings.Join(parts, ") AND (") + ")"
	return result
}
