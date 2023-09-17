package paginate

import (
	"fmt"
)

type Paginate struct {
	NextPage       uint
	TotalPages     uint
	EntriesPerPage uint
}

func (pg *Paginate) IsFirstPage() bool {
	return pg.CurrentPage() <= 0
}

func (pg *Paginate) IsLastPage() bool {
	return pg.NextPage >= pg.TotalPages
}

func (pg *Paginate) CurrentPage() uint {
	return pg.NextPage - 1
}

func (pg *Paginate) PreviousPage() uint {
	return pg.CurrentPage() - 1
}

func (pg *Paginate) GetPageCounter() string {
	if pg.TotalPages == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%d", pg.CurrentPage()+1, pg.TotalPages)
}
