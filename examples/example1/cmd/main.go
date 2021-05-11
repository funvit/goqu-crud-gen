package main

import (
	"context"
	"example1/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	dsn   = "test_login:test_pass@tcp(localhost:53306)/example1"
	table = "user"
)

func main() {
	fmt.Println("Example1")

	repo := model.NewUserRepo(dsn)
	err := repo.Connect(3 * time.Second)
	if err != nil {
		log.Fatalln("repo connect error: %s", err)
	}

	// migrate
	err = migrateUp(dsn, table)
	if err != nil {
		log.Fatalln(err)
	}
	defer migrateDown(dsn, table)

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
			log.Fatalln("Create error:", err)
		}
		log.Printf("user record created: %+v", u)

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
			log.Fatalln("Get error: %w", err)
		}
		if u == nil {
			log.Fatalln("user not exists")
		}

		log.Printf("repo Get user result: %+v", u)
	}

	fmt.Println("End.")
}

func migrateUp(dsn, table string) error {
	c, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer c.Close()

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
	defer c.Close()

	_, err = c.Exec(fmt.Sprintf(`DROP TABLE %s`, table))
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}
