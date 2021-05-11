package model

//go:generate goqu-crud-gen -model User -table user -dialect mysql -g
type User struct {
	Id    int64  `db:"id,primary,auto"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
