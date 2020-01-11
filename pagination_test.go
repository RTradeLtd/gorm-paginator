package pagination

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int
	UserName string `gorm:"not null;size:100;unique"`
}

func Test_Pagination(t *testing.T) {
	db, err := gorm.Open("sqlite3", "example.db")
	if err == nil {
		db.AutoMigrate(&User{})
		count := 0
		db.Model(User{}).Count(&count)
		if count == 0 {
			db.Create(User{ID: 1, UserName: "biezhi"})
			db.Create(User{ID: 2, UserName: "rose"})
			db.Create(User{ID: 3, UserName: "jack"})
			db.Create(User{ID: 4, UserName: "lili"})
			db.Create(User{ID: 5, UserName: "bob"})
			db.Create(User{ID: 6, UserName: "tom"})
			db.Create(User{ID: 7, UserName: "anny"})
			db.Create(User{ID: 8, UserName: "wat"})
			fmt.Println("Insert OK!")
		}
	} else {
		fmt.Println(err)
		return
	}

	var users []User

	Paging(&Param{
		DB:      db.Where("id > ?", 0),
		Page:    1,
		Limit:   3,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &users)
}
