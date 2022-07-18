package postgres

import (
	"authservice/domain"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const tdbEnv = "AUTH_TEST_DB"

var tdb *Postgres

func restoreDB(tdb *Postgres) error {
	b, err := os.ReadFile(filepath.Join("testdata", "t.sql"))
	if err != nil {
		return err
	}

	return tdb.exec(context.Background(), string(b))
}

func TestMain(m *testing.M) {

	connstr, ok := os.LookupEnv(tdbEnv)
	if !ok {
		fmt.Fprintf(os.Stderr, "environment variable %q must be set\n", tdbEnv)
		os.Exit(m.Run()) // tests will be skipped
	}

	var err error
	tdb, err = New(connstr)
	if err != nil {
		log.Fatalf("db connection: %v", err)
	}
	defer tdb.Close()

	if err = restoreDB(tdb); err != nil {
		tdb.Close()
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestPostgres(t *testing.T) {
	if tdb == nil {
		t.Skip("no connection to test database, skipped...")
	}

	t.Run("AddUser", func(t *testing.T) {
		want := testUser3

		err := tdb.AddUser(context.Background(), want)
		if err != nil {
			t.Fatalf("AddUser() = error %v", err)
		}

		got, err := tdb.User(context.Background(), want.Login)
		if err != nil {
			t.Fatalf("AddUser() = error %v", err)
		}

		if got != want {
			t.Errorf("AddUser() = %#v, want %#v", got, want)
		}

	})

	t.Run("User", func(t *testing.T) {
		want := testUser1

		got, err := tdb.User(context.Background(), want.Login)
		if err != nil {
			t.Fatalf("User() = error %v", err)
		}

		if got != want {
			t.Errorf("User() = %#v, want %#v", got, want)
		}

	})

	t.Run("DisableUser", func(t *testing.T) {

		err := tdb.DisableUser(context.Background(), testUser3)
		if err != nil {
			t.Fatalf("DisableUser() = error %v", err)
		}

		got, err := tdb.User(context.Background(), testUser3.Login)
		if err != nil {
			t.Fatalf("User() = error %v", err)
		}

		if !got.IsDisabled {
			t.Errorf("DisableUser() = is disabled %t, want%t", got.IsDisabled, true)
		}

	})
	t.Run("Password", func(t *testing.T) {
		want := testUser1

		got, err := tdb.Password(context.Background(), want)
		if err != nil {
			t.Fatalf("Password() = error %v", err)
		}

		if got != want.Password {
			t.Errorf("Password() = %#v, want %#v", got, want)
		}

	})

	t.Run("UpdatePassword", func(t *testing.T) {
		want := testUser3
		want.Password.Hash = "new-password"

		err := tdb.UpdatePassword(context.Background(), want)
		if err != nil {
			t.Fatalf("UpdatePassword() = error %v", err)
		}

		got, err := tdb.Password(context.Background(), want)
		if err != nil {
			t.Fatalf("Password() = error %v", err)
		}

		if got != want.Password {
			t.Errorf("UpdatePassword() = %#v, want %#v", got, want)
		}

	})
}

var testUser1 = domain.User{
	Login:      "login1",
	CreatedAt:  1658141437,
	IsDisabled: false,
	Password: domain.Password{
		Hash:        "h1",
		GeneratedAt: 1658141437,
		IsActive:    true,
	},
}

var testUser2 = domain.User{
	Login:      "login2",
	CreatedAt:  1658141437,
	IsDisabled: false,
	Password: domain.Password{
		Hash:        "h2",
		GeneratedAt: 1658141437,
		IsActive:    true,
	},
}

var testUser3 = domain.User{
	Login:      "login3",
	CreatedAt:  1658141437,
	IsDisabled: false,
	Password: domain.Password{
		Hash:        "h3",
		GeneratedAt: 1658141437,
		IsActive:    true,
	},
}
