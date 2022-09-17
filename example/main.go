package main

import (
	"context"
	"github.com/bigmate/bsql"
	"github.com/bigmate/bsql/example/models"
	"github.com/bigmate/bsql/example/repository"
	"github.com/bigmate/bsql/example/services/profile"
	"github.com/bigmate/bsql/example/services/user"
	"github.com/jmoiron/sqlx"
	"log"
)

const (
	driverName = ""
	dataSource = ""
)

func main() {
	db, _ := sqlx.Open(driverName, dataSource)
	trx, factory := bsql.New(db, nil)

	userSrv := user.NewService(trx,
		repository.NewUserRepository(factory),
		profile.NewService(trx, repository.NewProfileRepository(factory)),
	)

	created, err := userSrv.Create(context.Background(), models.User{Username: "tsoding"})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(created)
}
