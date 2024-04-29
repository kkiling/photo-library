package utils

import (
	"github.com/kkiling/photo-library/backend/api/internal/service/model"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testDuration struct {
	Duration *time.Duration `validate:"omitempty,duration-min=24h,duration-max=48h"`
}

func Test_validateDuration(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		validate := NewValidator()

		req := testDuration{
			Duration: lo.ToPtr(25 * time.Hour),
		}

		err := validate.Struct(req)
		require.NoError(t, err)
	})

	t.Run("max error", func(t *testing.T) {
		validate := NewValidator()

		req := testDuration{
			Duration: lo.ToPtr(50 * time.Hour),
		}

		err := validate.Struct(req)
		require.Error(t, err)
	})

	t.Run("min error", func(t *testing.T) {
		validate := NewValidator()

		req := testDuration{
			Duration: lo.ToPtr(1 * time.Hour),
		}

		err := validate.Struct(req)
		require.Error(t, err)
	})

	t.Run("null no error", func(t *testing.T) {
		validate := NewValidator()

		req := testDuration{
			Duration: nil,
		}

		err := validate.Struct(req)
		require.Nil(t, err)
	})
}

type testPagination struct {
	Name      string `validate:"required,min=3,max=128"`
	Paginator model.Pagination
}

func Test_validatePagination(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		validate := NewValidator()

		req := testPagination{
			Name: "name",
			Paginator: model.Pagination{
				Page:    1,
				PerPage: 10,
			},
		}

		err := validate.Struct(req)
		require.NoError(t, err)
	})

	t.Run("fail per page", func(t *testing.T) {
		validate := NewValidator()

		req := testPagination{
			Name: "name",
			Paginator: model.Pagination{
				Page:    1,
				PerPage: 10000,
			},
		}

		err := validate.Struct(req)
		require.Error(t, err)
	})
}
