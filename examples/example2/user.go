package example2

type User struct {
	Id    int64  `db:"id,primary,auto"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
