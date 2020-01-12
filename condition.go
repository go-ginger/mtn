package mtn

func (p *Parser) generateCondition(op, key string, value interface{}) *string {
	c, exists := p.conditions[op]
	if exists {
		query := key
		query += c
		query += value.(string)
		return &query
	}
	return nil
}
