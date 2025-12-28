package engine

import (
	"cmp"
	"sort"
)

func match(row Row, conds []Condition) (bool, error) {
	for _, c := range conds {
		v, ok := row[c.Field]
		if !ok {
			return false, nil
		}
		ok2, err := compare(v, c.Op, c.Value)
		if err != nil {
			return false, err
		}
		if !ok2 {
			return false, nil
		}
	}
	return true, nil
}

func compare(a any, op string, b any) (bool, error) {
	// numbers: json.Unmarshal => float64
	if af, ok := toFloat(a); ok {
		bf, ok := toFloat(b)
		if !ok {
			return false, ErrTypeMismatch
		}
		switch op {
		case "=", "==":
			return af == bf, nil
		case "!=":
			return af != bf, nil
		case ">":
			return af > bf, nil
		case ">=":
			return af >= bf, nil
		case "<":
			return af < bf, nil
		case "<=":
			return af <= bf, nil
		default:
			return false, ErrInvalidOp
		}
	}

	// strings
	if as, ok := a.(string); ok {
		bs, ok := b.(string)
		if !ok {
			return false, ErrTypeMismatch
		}
		switch op {
		case "=", "==":
			return as == bs, nil
		case "!=":
			return as != bs, nil
		case ">":
			return as > bs, nil
		case ">=":
			return as >= bs, nil
		case "<":
			return as < bs, nil
		case "<=":
			return as <= bs, nil
		default:
			return false, ErrInvalidOp
		}
	}

	// bool
	if ab, ok := a.(bool); ok {
		bb, ok := b.(bool)
		if !ok {
			return false, ErrTypeMismatch
		}
		switch op {
		case "=", "==":
			return ab == bb, nil
		case "!=":
			return ab != bb, nil
		default:
			return false, ErrInvalidOp
		}
	}

	// fallback: only equality
	switch op {
	case "=", "==":
		return a == b, nil
	case "!=":
		return a != b, nil
	default:
		return false, ErrInvalidOp
	}
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case uint64:
		return float64(t), true
	default:
		return 0, false
	}
}

func applySort(rows []Row, sortBy string, desc bool) {
	if sortBy == "" {
		return
	}

	sort.Slice(rows, func(i, j int) bool {
		ai, aok := rows[i][sortBy]
		aj, bok := rows[j][sortBy]
		if !aok && !bok {
			return false
		}
		if !aok {
			return !desc
		}
		if !bok {
			return desc
		}

		// numeric first
		if af, ok := toFloat(ai); ok {
			bf, ok2 := toFloat(aj)
			if !ok2 {
				return !desc
			}
			if desc {
				return af > bf
			}
			return af < bf
		}

		// string next
		as, ok := ai.(string)
		bs, ok2 := aj.(string)
		if ok && ok2 {
			c := cmp.Compare(as, bs)
			if desc {
				return c > 0
			}
			return c < 0
		}

		// fallback: do nothing stable-ish
		return !desc
	})
}
