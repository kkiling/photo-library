package model

const (
	// PerPageDefault дефолтное значение для пагинации
	PerPageDefault = 25
	// PerPageMax максимально-возможное кол-во записей на странице
	PerPageMax = 250
	// PageMax максимальная страница
	PageMax = 1000
)

// Pagination пагинация
type Pagination struct {
	Page    uint64 `validate:"gte=0,lte=1000"`
	PerPage uint64 `validate:"gte=1,lte=250"`
}

// GetLimit возвращает количество строк на странице
func (g *Pagination) GetLimit() int32 {
	if g.PerPage == 0 || g.PerPage > PerPageMax {
		return PerPageDefault
	}

	return int32(g.PerPage)
}

// GetOffset возвращает номер строки, с которой надо начинать выборку
func (g *Pagination) GetOffset() int32 {
	if g.Page == 0 {
		return 0
	}

	return int32(g.Page-1) * g.GetLimit()
}
