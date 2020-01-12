package mtn

func (p *Parser) Initialize() {
	p.operators = map[string]string{
		"$and": " AND ",
		"$or":  " OR ",
	}
	p.conditions = map[string]string{
		"$lt":  " < ",
		"$lte": " <= ",
		"$gt":  " > ",
		"$gte": " >= ",
		"$ne":  " <> ",
	}
}
