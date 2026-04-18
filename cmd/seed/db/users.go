package db_seed

import (
	"context"

	sqlc "github.com/Mboukhal/SvGoPg/internal/db"
)

var users = []struct {
	name     string
	email    string
	password string
}{
	{
		name:     "Mohammed Boukhala",
		email:    "lios80466@gmail.com",
		password: "lios80466@gmail",
	},
}

func LoadUsers(q *sqlc.Queries, ctx context.Context) error {
	for _, user := range users {
		err := q.CreateUser(ctx, sqlc.CreateUserParams{
			Username:     user.name,
			Email:        user.email,
			PasswordHash: user.password,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
