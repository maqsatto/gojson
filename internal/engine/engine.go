package engine

type Engine interface {
	Select(table string, q Query) ([]Row, error)
	Insert(table string, row Row) error
	Update(table string, q Query, values Row) (int, error)
	Delete(table string, q Query) (int, error)
}
