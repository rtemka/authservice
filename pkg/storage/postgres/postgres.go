// package postgres conforms the contract
// of the authentication service.
package postgres

import (
	"authservice/domain"
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ErrNoPassword - error occurring when user
// was provided to DB without password.
var ErrNoPassword = errors.New("user doesn't have password")
var ErrNoRows = pgx.ErrNoRows

// Postgres handles CRUD operation over the
// postgresql database.
type Postgres struct {
	db *pgxpool.Pool
}

type statement struct {
	sql  string
	args []any
}

// New returns new *Postgres object.
func New(connstr string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: pool}, pool.Ping(context.Background())
}

// Close closes all connections to DB.
func (p *Postgres) Close() {
	p.db.Close()
}

// AddUser adds user and user password to DB,
// if no password is provided then ErrNoPassword will be returned.
func (p *Postgres) AddUser(ctx context.Context, user domain.User) error {
	if user.Password.Hash == "" {
		return ErrNoPassword
	}
	iu := statement{
		sql: `
		INSERT INTO 
			users(login, created_at, is_disabled)
		VALUES ($1, $2, FALSE);`,
		args: []any{user.Login, user.CreatedAt},
	}
	ip := statement{
		sql: `
		INSERT INTO
			passwords(user_login, hash, generated_at, is_active)
		VALUES ($1, $2, $3, TRUE);`,
		args: []any{user.Login, user.Password.Hash, user.Password.GeneratedAt},
	}

	return p.execBatch(ctx, iu, ip)
}

// User fetches user's info from the DB. Information also
// includes user's current password.
func (p *Postgres) User(ctx context.Context, login string) (domain.User, error) {
	sql := `
		SELECT 
			u.login, u.created_at, u.is_disabled, 
			p.hash, p.generated_at, p.is_active
		FROM users as u
		JOIN passwords as p ON u.login = p.user_login
		WHERE u.login = $1 AND p.is_active = TRUE;
	`
	var u domain.User

	err := p.db.QueryRow(ctx, sql, login).Scan(&u.Login, &u.CreatedAt, &u.IsDisabled,
		&u.Password.Hash, &u.Password.GeneratedAt, &u.Password.IsActive)
	if err != nil {
		return u, err
	}

	return u, nil
}

// DisableUser sets user's status to disabled.
func (p *Postgres) DisableUser(ctx context.Context, user domain.User) error {
	sql := `
		UPDATE users
			SET is_disabled = TRUE
		WHERE login = $1;`
	return p.exec(ctx, sql, user.Login)
}

// UpdatePassword updates user's password and invalidates
// old one if any.
func (p *Postgres) UpdatePassword(ctx context.Context, user domain.User) error {
	if user.Password.Hash == "" {
		return ErrNoPassword
	}

	us := statement{
		sql: `
			UPDATE passwords
				SET is_active = FALSE
			WHERE user_login = $1;`,
		args: []any{user.Login},
	}
	is := statement{
		sql: `
			INSERT INTO
				passwords(user_login, hash, generated_at, is_active)
			VALUES ($1, $2, $3, TRUE);`,
		args: []any{user.Login, user.Password.Hash, user.Password.GeneratedAt},
	}

	return p.execBatch(ctx, us, is)
}

// Password fetches user's current password from the DB.
func (p *Postgres) Password(ctx context.Context, user domain.User) (domain.Password, error) {
	sql := `
		SELECT p.hash, p.generated_at, p.is_active
		FROM passwords as p
		WHERE p.user_login = $1 AND p.is_active = TRUE;`

	var pwd domain.Password

	return pwd, p.db.QueryRow(ctx, sql, user.Login).
		Scan(&pwd.Hash, &pwd.GeneratedAt, &pwd.IsActive)
}

// exec is helper function, runs
// *pgx.conn.Exec() in transaction.
func (p *Postgres) exec(ctx context.Context, sql string, args ...any) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = p.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// execBatch is helper function, runs
// multiple batch queries *pgx.conn.Exec()
// in one transaction. This queries must be ones
// that don't return results (INSERT, UPDATE, DELETE).
func (p *Postgres) execBatch(ctx context.Context, stmts ...statement) error {

	b := new(pgx.Batch)

	for i := range stmts {
		b.Queue(stmts[i].sql, stmts[i].args...)
	}

	return p.db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return tx.SendBatch(ctx, b).Close()
	})
}
