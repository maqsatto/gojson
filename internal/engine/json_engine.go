package engine

type JSONEngine struct {
	st *JSONStorage
}

func NewJSONEngine(st *JSONStorage) *JSONEngine {
	return &JSONEngine{st: st}
}

func (e *JSONEngine) Select(table string, q Query) ([]Row, error) {
	e.st.mu.RLock()
	rows, ok := e.st.data[table]
	e.st.mu.RUnlock()
	if !ok {
		return nil, ErrTableNotFound
	}

	if idVal, ok := extractIDEq(q); ok {
		e.st.mu.RLock()
		idxMap := e.st.idIndex[table]
		i, ok2 := idxMap[idVal]
		if ok2 && i >= 0 && i < len(e.st.data[table]) {
			one := e.st.data[table][i]
			e.st.mu.RUnlock()
			ok3, err := matchQuery(one, q)
			if err != nil {
				return nil, err
			}
			if !ok3 {
				return []Row{}, nil
			}
			return []Row{copyRow(one)}, nil
		}
		e.st.mu.RUnlock()
		return []Row{}, nil
	}

	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		ok, err := matchQuery(r, q)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, r)
		}
	}

	applySort(out, q.SortBy, q.Desc)

	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Offset > len(out) {
		out = []Row{}
	} else {
		out = out[q.Offset:]
	}
	if q.Limit > 0 && q.Limit < len(out) {
		out = out[:q.Limit]
	}

	copied := make([]Row, 0, len(out))
	for _, r := range out {
		copied = append(copied, copyRow(r))
	}
	return copied, nil
}

func (e *JSONEngine) Insert(table string, row Row) error {
	if table == "" || row == nil {
		return ErrBadRequest
	}

	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	if _, has := row["id"]; !has {
		e.st.maxID[table] = e.st.maxID[table] + 1
		row["id"] = e.st.maxID[table]
	} else {
		idf, ok := toFloat(row["id"])
		if !ok {
			return ErrTypeMismatch
		}
		if _, exists := e.st.idIndex[table][idf]; exists {
			return ErrDuplicateID
		}
		if idf > e.st.maxID[table] {
			e.st.maxID[table] = idf
		}
	}

	e.st.data[table] = append(e.st.data[table], row)
	idf, _ := toFloat(row["id"])
	if e.st.idIndex[table] == nil {
		e.st.idIndex[table] = make(map[float64]int)
	}
	e.st.idIndex[table][idf] = len(e.st.data[table]) - 1

	return e.st.saveLockedAtomic()
}

func (e *JSONEngine) Update(table string, q Query, set Row) (int, error) {
	if table == "" || set == nil {
		return 0, ErrBadRequest
	}

	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	rows, ok := e.st.data[table]
	if !ok {
		return 0, ErrTableNotFound
	}

	affected := 0

	if idVal, ok := extractIDEq(q); ok {
		if i, ok2 := e.st.idIndex[table][idVal]; ok2 && i >= 0 && i < len(rows) {
			r := rows[i]
			ok3, err := matchQuery(r, q)
			if err != nil {
				return 0, err
			}
			if !ok3 {
				return 0, nil
			}
			delete(set, "id")
			for k, v := range set {
				r[k] = v
			}
			affected = 1
			e.st.data[table] = rows
			return affected, e.st.saveLockedAtomic()
		}
		return 0, nil
	}

	for _, r := range rows {
		ok, err := matchQuery(r, q)
		if err != nil {
			return 0, err
		}
		if ok {
			delete(set, "id")
			for k, v := range set {
				r[k] = v
			}
			affected++
		}
	}

	e.st.data[table] = rows
	return affected, e.st.saveLockedAtomic()
}

func (e *JSONEngine) Delete(table string, q Query) (int, error) {
	if table == "" {
		return 0, ErrBadRequest
	}

	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	rows, ok := e.st.data[table]
	if !ok {
		return 0, ErrTableNotFound
	}

	if idVal, ok := extractIDEq(q); ok {
		if i, ok2 := e.st.idIndex[table][idVal]; ok2 && i >= 0 && i < len(rows) {
			r := rows[i]
			ok3, err := matchQuery(r, q)
			if err != nil {
				return 0, err
			}
			if !ok3 {
				return 0, nil
			}
			last := len(rows) - 1
			rows[i] = rows[last]
			rows = rows[:last]

			e.st.data[table] = rows
			e.st.rebuildIndexLocked()
			return 1, e.st.saveLockedAtomic()
		}
		return 0, nil
	}

	out := make([]Row, 0, len(rows))
	affected := 0
	for _, r := range rows {
		ok, err := matchQuery(r, q)
		if err != nil {
			return 0, err
		}
		if ok {
			affected++
			continue
		}
		out = append(out, r)
	}

	e.st.data[table] = out
	e.st.rebuildIndexLocked()
	return affected, e.st.saveLockedAtomic()
}

func copyRow(r Row) Row {
	nr := Row{}
	for k, v := range r {
		nr[k] = v
	}
	return nr
}

func extractIDEq(q Query) (float64, bool) {
	// только из Where (простого) для быстрого пути
	for _, c := range q.Where {
		if c.Field == "id" && (c.Op == "=" || c.Op == "==") {
			idf, ok := toFloat(c.Value)
			return idf, ok
		}
	}
	return 0, false
}
