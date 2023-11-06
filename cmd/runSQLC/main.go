package main

import (
	"context"
	"database/sql"

	"github.com/deividroger/sqlc/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/course")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   uuid.New().String(),
		Name: "Backend",
		Description: sql.NullString{
			String: "Backend course",
			Valid:  true,
		},
	})

	if err != nil {
		panic(err)
	}
	categories, err := queries.ListCategories(ctx)

	for _, category := range categories {
		println(category.ID, category.Name, category.Description.String)
	}

}
