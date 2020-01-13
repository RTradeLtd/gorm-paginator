package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Param is used to configure the paging query
type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

// Paginator is the response returned by Paging
// this should be fed to a http request response.
// If using gin-gonic, you should use c.JSON(200, Paginator)
type Paginator struct {
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

// Paging is used to return a paged result from the database
func Paging(p *Param, result interface{}) (*Paginator, error) {
	db := p.DB
	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}
	var (
		paginator     = new(Paginator)
		count, offset int
	)
	if err := db.Model(result).Count(&count).Error; err != nil {
		return nil, err
	}
	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}
	if err := db.Limit(p.Limit).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}
	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page
	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))
	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}
	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return paginator, nil
}
