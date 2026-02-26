package database

import (
	"ToDoList/internal/entity"
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type PostgresRepo struct {
	conn *pgx.Conn
}

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{
		conn: conn,
	}
}

func InitDatabase(ctx context.Context, connstr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connstr)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(connstr); err != nil {
		return nil, err
	}

	return conn, nil
}

func runMigrations(databaseURL string) error {

	m, err := migrate.New(
		"file://internal/infrastructure/migrations",
		databaseURL,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func (p *PostgresRepo) AddTask(ctx context.Context, t *entity.Task) error {

	query := `
	INSERT INTO tasks (name, text, time_add, is_done, is_important) 
	VALUES ($1, $2, $3, $4, $5)
	`

	return p.conn.QueryRow(ctx, query, t.Name, t.Text, t.Time_add, t.Is_done, t.Is_important).Scan(&t.Id)
}

func (p *PostgresRepo) DeleteTask(ctx context.Context, id int64) error {
	query := `
	DELETE FROM tasks WHERE id = $1
	`

	_, err := p.conn.Exec(ctx, query, id)
	return err
}

func (p *PostgresRepo) Update(ctx context.Context, t *entity.Task) error {
	query := `
	UPDATE tasks
	SET 
	name = $1,
	text = $2,
	time_add = $3,
	is_done = $4,
	time_done = $5,
	is_important = $6
	WHERE id = $7	
	`

	_, err := p.conn.Exec(ctx, query, t.Name, t.Text, t.Time_add, t.Is_done, t.Time_done, t.Is_important, t.Id)
	return err
}

func (p *PostgresRepo) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	query := `
	SELECT id, name, text, time_add, is_done, time_done, is_important
	FROM tasks
	`

	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var t entity.Task

		if err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Text,
			&t.Time_add,
			&t.Is_done,
			&t.Time_done,
			&t.Is_important,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (p *PostgresRepo) GetTaskById(ctx context.Context, id int64) (*entity.Task, error) {
	query := `
	SELECT id, name, text, time_add, is_done, time_done, is_important
	FROM tasks
	WHERE id = $1
	`
	row := p.conn.QueryRow(ctx, query, id)

	var t entity.Task
	if err := row.Scan(
		&t.Id,
		&t.Name,
		&t.Text,
		&t.Time_add,
		&t.Is_done,
		&t.Time_done,
		&t.Is_important,
	); err != nil {
		return nil, err
	}

	return &t, nil
}

func (p *PostgresRepo) GetAllTasksPages(ctx context.Context, n int) (map[int][]entity.Task, error) {

	var total int
	if err := p.conn.QueryRow(ctx, `SELECT COUNT(*) FROM tasks`).Scan(&total); err != nil {
		return nil, err
	}
	fmt.Println(total)
	tasksPages := make(map[int][]entity.Task)

	var pages int
	if total%n == 0 {
		pages = total / n
	} else {
		pages = (total / n) + 1
	}

	for i := 0; i < pages; i++ {
		query := `
		SELECT id, name, text, time_add, is_done, time_done
		FROM tasks ORDER BY id ASC
		LIMIT $1 OFFSET $2
		`

		rows, err := p.conn.Query(ctx, query, n, n*i)
		defer rows.Close()
		if err != nil {
			return nil, err
		}

		var tasks []entity.Task

		for rows.Next() {
			var t entity.Task
			if err := rows.Scan(
				&t.Id,
				&t.Name,
				&t.Text,
				&t.Time_add,
				&t.Is_done,
				&t.Time_done,
			); err != nil {
				return nil, err
			}
			tasks = append(tasks, t)
		}
		tasksPages[i] = tasks
	}

	return tasksPages, nil
}
