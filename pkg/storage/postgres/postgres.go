package postgres

import (
	"authservice/pkg/storage"
	"context"
	"errors"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNoPassword = errors.New("user doesn't have password")
var ErrNoRows = pgx.ErrNoRows

// Postgres handles CRUD operation over the
// postgresql database
type Postgres struct {
	db *pgxpool.Pool
}

// New returns new *Postgres object
func New(connstr string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: pool}, pool.Ping(context.Background())
}

func (p *Postgres) AddUser(ctx context.Context, user storage.User) error {
	if user.Password.Hash == "" {
		return ErrNoPassword
	}
	sql := `
		INSERT 
			users(login, created_at, is_disabled)
		VALUES ($1, $2, FALSE);
		INSERT 
			passwords(login, hash, generated_at, is_active)
		VALUES ($1, $3, $4, TRUE);`

	return p.exec(ctx, sql,
		user.Login, user.CreatedAt, user.Password.Hash, user.Password.GeneratedAt)
}

func (p *Postgres) UpdatePassword(ctx context.Context, user storage.User) error {
	if user.Password.Hash == "" {
		return ErrNoPassword
	}
	sql := `
		UPDATE passwords
			SET is_active = FALSE
			WHERE login = $1; 
		INSERT 
			password(login, generated_at, is_active)
		VALUES ($1, $2, TRUE);`

	return p.exec(ctx, sql, user.Login,
		user.Password.Hash, user.Password.GeneratedAt)
}

func (p *Postgres) Password(ctx context.Context, user storage.User) (storage.Password, error) {
	sql := `
		SELECT p.hash, p.generated_at, p.is_active
		FROM passwords as p
		WHERE p.login = $1 AND p.is_active = TRUE;`

	var pwd storage.Password

	return pwd, p.db.QueryRow(ctx, sql).Scan(&pwd.Hash, &pwd.GeneratedAt, &pwd.IsActive)
}

// exec is helper function, runs
// *pgx.conn.Exec() in transaction
func (p *Postgres) exec(ctx context.Context, sql string, args ...any) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = p.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
