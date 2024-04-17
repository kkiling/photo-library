package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/kkiling/photo-library/backend/api/internal/service/serviceerr"
)

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
	Page    uint64
	PerPage uint64
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

func (g *Pagination) Validate() error {
	validate := validator.New()

	if err := validate.Var(g.Page, fmt.Sprintf("gte=%d,lte=%d", 0, PageMax)); err != nil {
		return serviceerr.MakeErr(err, "invalid page")
	}

	if err := validate.Var(g.PerPage, fmt.Sprintf("gte=%d,lte=%d", 1, PerPageMax)); err != nil {
		return serviceerr.MakeErr(err, "invalid perPage")
	}

	return nil
}
