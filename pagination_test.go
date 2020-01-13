package pagination

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int
	Name string `gorm:"not null;size:100;unique"`
}

func Test_Pagination(t *testing.T) {
	db, err := gorm.Open("sqlite3", "example.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("example.db")
	db.AutoMigrate(&User{})
	count := 0
	db.Model(User{}).Count(&count)
	if count == 0 {
		for i := 0; i < 100; i++ {
			if err := db.Create(&User{
				Name: fmt.Sprintf("user-%v", i),
			}).Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	var users []User
	tests := []struct {
		name    string
		param   *Param
		model   interface{}
		wantErr bool
	}{
		{"Defaults", &Param{
			DB: db.Where("id > ?", 0),
		}, &users, false},
		{"Non-Defaults", &Param{
			DB:      db.Where("id > ?", 0),
			Page:    2,
			Limit:   10,
			OrderBy: []string{"name desc"},
		}, &users, false},
		{"Debug", &Param{
			DB:      db.Where("id > ?", 0),
			ShowSQL: true,
		}, &users, false},
		{"BadModel", &Param{
			DB: db.Where("id > ?", 0),
		}, &struct{ Dog string }{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Paging(tt.param, tt.model)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Paging() err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
