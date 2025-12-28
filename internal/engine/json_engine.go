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

	// filter
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		ok, err := match(r, q.Where)
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, r)
		}
	}

	// sort
	applySort(out, q.SortBy, q.Desc)

	// offset/limit
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

	// copy rows so caller can't mutate shared state
	copied := make([]Row, 0, len(out))
	for _, r := range out {
		nr := Row{}
		for k, v := range r {
			nr[k] = v
		}
		copied = append(copied, nr)
	}
	return copied, nil
}

func (e *JSONEngine) Insert(table string, row Row) error {
	if table == "" || row == nil {
		return ErrBadRequest
	}

	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	e.st.data[table] = append(e.st.data[table], row)
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
	for _, r := range rows {
		ok, err := match(r, q.Where)
		if err != nil {
			return 0, err
		}
		if ok {
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

	out := make([]Row, 0, len(rows))
	affected := 0

	for _, r := range rows {
		ok, err := match(r, q.Where)
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
	return affected, e.st.saveLockedAtomic()
}
