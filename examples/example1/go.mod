module example1

go 1.14

require (
	github.com/doug-martin/goqu/v9 v9.12.0
	github.com/funvit/goqu-crud-gen v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.3
)

// this replace forces example to use current goqu-crud-gen branch
replace github.com/funvit/goqu-crud-gen => ../..
