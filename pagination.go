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

// Paged is the response returned by Paging
// this should be fed to a http request response.
// If using gin-gonic, you should use c.JSON(200, Paged)
type Paged struct {
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
func Paging(p *Param, result interface{}) (*Paged, error) {
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
		paged         = new(Paged)
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
	paged.TotalRecord = count
	paged.Records = result
	paged.Page = p.Page
	paged.Offset = offset
	paged.Limit = p.Limit
	paged.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))
	if p.Page > 1 {
		paged.PrevPage = p.Page - 1
	} else {
		paged.PrevPage = p.Page
	}
	if p.Page == paged.TotalPage {
		paged.NextPage = p.Page
	} else {
		paged.NextPage = p.Page + 1
	}
	return paged, nil
}
