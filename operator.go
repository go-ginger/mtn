package mtn

func (p *Parser) generateOperator(op string) *string {
	c, exists := p.operators[op]
	if exists {
		query := c
		return &query
	}
	return nil
}
