package engine

type JSONEngine struct {
	st *JSONStorage
}

func NewJSONEngine(st *JSONStorage) *JSONEngine {
	return &JSONEngine{st: st}
}

func (e *JSONEngine) Select(table string, q Query) ([]Row, error) {
	e.st.mu.RLock()
	defer e.st.mu.RUnlock()

	rows, ok := e.st.data[table]
	if !ok {
		return nil, ErrTableNotFound
	}

	var output []Row
	for _, row := range rows {
		if match(row, q.Where) {
			output = append(output, row)
		}
	}
	return output, nil
}

func (e *JSONEngine) Insert(table string, row Row) error {
	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	e.st.data[table] = append(e.st.data[table], row)
	return e.st.save()
}

func (e *JSONEngine) Update(table string, q Query, values Row) (int, error) {
	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	rows := e.st.data[table]
	affected := 0
	for _, row := range rows {
		if match(row, q.Where) {
			for key, value := range values {
				row[key] = value
			}
			affected++
		}
	}
	return affected, nil
}

func (e *JSONEngine) Delete(table string, q Query) (int, error) {
	e.st.mu.Lock()
	defer e.st.mu.Unlock()

	rows := e.st.data[table]
	var output []Row
	affected := 0

	for _, row := range rows {
		if match(row, q.Where) {
			affected++
			continue
		}
		output = append(output, row)
	}
	e.st.data[table] = output
	return affected, e.st.save()
}
