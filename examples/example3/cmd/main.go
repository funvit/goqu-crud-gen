package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"example3/adapters/mysql"
	"example3/domain"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	dsn = "test_login:test_pass@tcp(localhost:53306)/example3"
)

var (
	logInfo = log.New(os.Stdout, "INF ", log.LstdFlags)
	logErr  = log.New(os.Stderr, "ERR ", log.LstdFlags)
)

func main() {
	fmt.Println("Example3")

	repo := mysql.NewUserRepo(dsn)
	err := repo.Connect(3 * time.Second)
	if err != nil {
		logErr.Fatalf("repo connect error: %s", err)
	}
	repo.SetMaxOpenConns(100)
	repo.SetConnMaxLifetime(5 * time.Minute)

	// migrate
	err = migrateUp(dsn)
	if err != nil {
		logErr.Fatalln(err)
	}
	defer func() {
		err = migrateDown(dsn)
		if err != nil {
			logErr.Fatalln("migrate down error:", err)
		}
	}()

	// add record
	var userId uuid.UUID
	{
		u := domain.User{
			Id:   uuid.New(),
			Name: "John",
			Account: &domain.Account{
				Login:        "johndoe",
				PasswordHash: "3492u8riehfvbd",
			},
		}

		err = repo.WithTran(context.Background(), func(ctx context.Context) error {
			return repo.Create(ctx, u)
		})
		if err != nil {
			logErr.Fatalln("Create error:", err)
		}
		logInfo.Printf("user record created: %+v", u)

		userId = u.Id
	}

	// get record
	{
		var u *domain.User
		err := repo.WithTran(context.Background(), func(ctx context.Context) error {
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

	// delete record
	{
		err := repo.WithTran(context.Background(), func(ctx context.Context) error {
			return repo.Delete(ctx, userId)
		})
		if err != nil {
			logErr.Fatalln("Delete error:", err)
		}

		logInfo.Printf("repo Delete succeed")
	}

	fmt.Println("End.")
}

func migrateUp(dsn string) error {
	c, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer func() { _ = c.Close() }()

	queries := []string{
		`
		CREATE TABLE user
		(
		    id varchar(36) not null,
		    name varchar(255),
		    
		    primary key (id),
		    unique key name_uniq (name)
		)`,
		`CREATE TABLE account
		(
			user_id varchar(36) not null,
			login varchar(255),
			pass varchar(255),

			primary key (user_id),
		    unique key login_uniq (login)
		)`,
	}

	for _, q := range queries {
		_, err = c.Exec(q)
		if err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
	}

	return nil
}

func migrateDown(dsn string) error {
	c, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}
	defer func() { _ = c.Close() }()

	queries := []string{
		`DROP TABLE user`,
		`DROP TABLE account`,
	}
	for _, q := range queries {
		_, err = c.Exec(q)
		if err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
	}

	return nil
}
