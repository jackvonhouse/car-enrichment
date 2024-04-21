package car

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackvonhouse/car-enrichment/internal/dto"
	"strings"
)

func (r Repository) where(
	builder sq.SelectBuilder,
	filter dto.Filter,
) sq.SelectBuilder {

	builder = r.whereName(builder, filter.OwnerName)
	builder = r.whereSurname(builder, filter.OwnerSurname)
	builder = r.wherePatronymic(builder, filter.OwnerPatronymic)
	builder = r.whereRegNumber(builder, filter.RegNum)
	builder = r.whereMark(builder, filter.Mark)
	builder = r.whereModel(builder, filter.Model)
	builder = r.whereYear(builder, filter.Year)

	return builder
}

func (r Repository) whereName(
	builder sq.SelectBuilder,
	name string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(name)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"owner.name": fmt.Sprintf("%%%s%%", name),
		},
	)
}

func (r Repository) whereSurname(
	builder sq.SelectBuilder,
	surname string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(surname)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"owner.surname": fmt.Sprintf("%%%s%%", surname),
		},
	)
}

func (r Repository) wherePatronymic(
	builder sq.SelectBuilder,
	patronymic string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(patronymic)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"owner.patronymic": fmt.Sprintf("%%%s%%", patronymic),
		},
	)
}

func (r Repository) whereRegNumber(
	builder sq.SelectBuilder,
	regNumber string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(regNumber)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"car.regNum": fmt.Sprintf("%%%s%%", regNumber),
		},
	)
}

func (r Repository) whereMark(
	builder sq.SelectBuilder,
	mark string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(mark)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"car.mark": fmt.Sprintf("%%%s%%", mark),
		},
	)
}

func (r Repository) whereModel(
	builder sq.SelectBuilder,
	model string,
) sq.SelectBuilder {

	if len(strings.TrimSpace(model)) == 0 {
		return builder
	}

	return builder.Where(
		sq.Like{
			"car.model": fmt.Sprintf("%%%s%%", model),
		},
	)
}

func (r Repository) whereYear(
	builder sq.SelectBuilder,
	year int,
) sq.SelectBuilder {

	if year == 0 {
		return builder
	}

	return builder.Where(
		sq.Eq{
			"car.year": year,
		},
	)
}
