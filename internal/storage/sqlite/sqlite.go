package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/IlianBuh/Follow_Service/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Follower interface {
	Follow(context.Context, int, int) error
}
type Unfollower interface {
	Unfollow(context.Context, int, int) error
}
type FollowingsProvider interface {
	ListFollowers(context.Context, int) ([]int, error)
	ListFollowees(context.Context, int) ([]int, error)
}

type Storage struct {
	db *sql.DB
}

// New returns new instance of the repository. Database is stored by path 'path'
func New(path string) (*Storage, error) {
	const op = "sqlite.New"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS followings(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					follower INTEGER NOT NULL,
					followee INTEGER NOT NULL,
					
					CONSTRAINT unique_followings UNIQUE (follower, followee)
				);
				CREATE INDEX IF NOT EXISTS idx_follower ON followings(follower);
				CREATE INDEX IF NOT EXISTS idx_followee ON followings(followee);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

// Follow add new tuple into the database
func (s *Storage) Follow(ctx context.Context, src, target int) error {
	const op = "sqlite.Follow"

	prep, err := s.db.PrepareContext(ctx, `INSERT INTO followings(follower, followee) VALUES(?, ?)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer prep.Close()

	_, err = prep.ExecContext(ctx, src, target)
	if err != nil {
		var sqlerr sqlite3.Error
		if errors.As(err, &sqlerr) && errors.Is(sqlerr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return storage.ErrFollowing
		}

		return err
	}

	return nil
}

// Unfollow delete the tuple (src, target) from the database
func (s *Storage) Unfollow(ctx context.Context, src, target int) error {
	const op = "sqlite.Unfollow"

	prep, err := s.db.PrepareContext(ctx, `DELETE FROM followings WHERE follower=? AND followee=?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer prep.Close()

	_, err = prep.ExecContext(ctx, src, target)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ListFollowers returns lists of all followers of the user with uuid
func (s *Storage) ListFollowers(ctx context.Context, uuid int) ([]int, error) {
	const op = "sqlite.ListFollowers"

	prep, err := s.db.PrepareContext(ctx, `SELECT follower FROM followings WHERE followee=?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := make([]int, 0)
	var temp int
	for rows.Next() {
		err = rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		list = append(list, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}

// ListFollowees returns lists of all followees of the user with uuid
func (s *Storage) ListFollowees(ctx context.Context, uuid int) ([]int, error) {
	const op = "sqlite.ListFollowees"

	prep, err := s.db.PrepareContext(ctx, `SELECT followee FROM followings WHERE follower=?`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	list := make([]int, 0)
	var temp int
	for rows.Next() {
		err = rows.Scan(&temp)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		list = append(list, temp)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return list, nil
}
