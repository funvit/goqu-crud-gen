package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"example1/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	dsn   = "test_login:test_pass@tcp(localhost:53306)/example1"
	table = "user"
)

var (
	logInfo = log.New(os.Stdout, "INF ", log.LstdFlags)
	logErr  = log.New(os.Stderr, "ERR ", log.LstdFlags)
)

func main() {
	fmt.Println("Example1")

	repo := model.NewUserRepo(dsn)
	err := repo.Connect(3 * time.Second)
	if err != nil {
		logErr.Fatalf("repo connect error: %s", err)
	}

	// migrate
	err = migrateUp(dsn, table)
	if err != nil {
		logErr.Fatalln(err)
	}
	defer func() {
		err = migrateDown(dsn, table)
		if err != nil {
			logErr.Fatalln("migrate down error:", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// add record
	var userId int64
	{
		u := model.User{
			Id:    0,
			Name:  "John",
			Email: "john.doe@email.com",
		}

		err = repo.WithTran(ctx, func(ctx context.Context) error {
			return repo.Create(ctx, &u)
		})
		if err != nil {
			logErr.Fatalln("Create error:", err)
		}
		logInfo.Printf("user record created: %+v", u)

		userId = u.Id
	}

	// get record
	{
		var u *model.User
		err := repo.WithTran(ctx, func(ctx context.Context) error {
			u, err = repo.Get(ctx, userId)
			return err
		})
		if err != nil {
			logErr.Fatalln("Get error:", err)
		}
		if u == nil {
			logErr.Fatalln("user not exists")
		}

		logInfo.Printf("repo Get user result: %+v", u)
	}

	fmt.Println("End.")
}

func migrateUp(dsn, table string) error {
	c, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer func() { _ = c.Close() }()

	_, err = c.Exec(fmt.Sprintf(`
		CREATE TABLE %s 
		(
		    id bigint not null auto_increment,
		    name varchar(255),
		    email varchar(255),
		    
		    primary key (id),
		    unique key name_uniq (name)
		)`,
		table,
	))
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}

func migrateDown(dsn, table string) error {
	c, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer func() { _ = c.Close() }()

	_, err = c.Exec(fmt.Sprintf(`DROP TABLE %s`, table))
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}
