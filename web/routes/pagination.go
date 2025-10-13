package routes

import (
	"strconv"

	"github.com/Lekuruu/go-puush/internal/app"
)

const UploadsPerPage = 60

type PaginationData struct {
	CurrentPage int
	TotalPages  int
	PerPage     int
	Start       int
	End         int
}

func (p *PaginationData) HasPrevious() bool {
	return p.CurrentPage > 1
}

func (p *PaginationData) HasNext() bool {
	return p.CurrentPage < p.TotalPages
}

func (p *PaginationData) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}

func (p *PaginationData) Limit() int {
	return p.PerPage
}

func NewPaginationData(currentPage, totalItems, perPage int) *PaginationData {
	if perPage <= 0 {
		perPage = UploadsPerPage
	}
	if currentPage <= 0 {
		currentPage = 1
	}

	totalPages := (totalItems + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}
	if currentPage > totalPages {
		currentPage = totalPages
	}

	start := (currentPage-1)*perPage + 1
	end := start + perPage - 1
	if end > totalItems {
		end = totalItems
	}
	if totalItems == 0 {
		start = 0
		end = 0
	}

	return &PaginationData{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		PerPage:     perPage,
		Start:       start,
		End:         end,
	}
}

func GetPageFromQuery(ctx *app.Context) int {
	currentPage := ctx.Request.URL.Query().Get("page")
	if currentPage == "" {
		currentPage = "1"
	}
	currentPageInt, err := strconv.Atoi(currentPage)
	if err != nil {
		currentPageInt = 1
	}
	return currentPageInt
}
