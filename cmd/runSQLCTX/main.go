package main

import (
	"context"
	"database/sql"

	"github.com/deividroger/sqlc/internal/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type CourseDb struct {
	dbConn *sql.DB
	*db.Queries
}

func NewCourseDb(dbConn *sql.DB) *CourseDb {
	return &CourseDb{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (c *CourseDb) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price       float64
}
type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDb) CreateCourseAndCategory(ctx context.Context, argsCategory CategoryParams, argsCourse CourseParams) error {
	err := c.callTx(ctx, func(q *db.Queries) error {
		var err error
		err = q.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          argsCategory.ID,
			Name:        argsCategory.Name,
			Description: argsCategory.Description,
		})
		if err != nil {
			return err
		}
		err = q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          argsCourse.ID,
			Name:        argsCourse.Name,
			Description: argsCourse.Description,
			CategoryID:  argsCategory.ID,
			Price:       argsCourse.Price,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/course")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	courseArgs := CourseParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend course", Valid: true},
		Price:       10.98,
	}

	categoryArgs := CategoryParams{
		ID:          uuid.New().String(),
		Name:        "Backend",
		Description: sql.NullString{String: "Backend course", Valid: true},
	}

	courseDb := NewCourseDb(dbConn)
	err = courseDb.CreateCourseAndCategory(ctx, categoryArgs, courseArgs)

	if err != nil {
		panic(err)
	}

}
