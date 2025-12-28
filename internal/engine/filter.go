package engine

func match(row Row, conds []Condition) bool {
	for _, cond := range conds {
		v, ok := row[cond.Field]
		if !ok {
			return false
		}
		switch cond.Op {
		case "=":
			if v != cond.Value {
				return false
			}
		default:
			return false
		}
	}
	return true
}
