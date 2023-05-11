package repositoryimpl

import (
	"context"
	"database/sql"

	"github.com/kenty51107/go-jwt/pkg/domain/model"
	"github.com/kenty51107/go-jwt/pkg/domain/repository"
)

type DBTX interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repositoryImpl struct {
    db DBTX
}

func NewRepository(db DBTX) repository.Repository {
    return &repositoryImpl{db: db}
}

func (ri *repositoryImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
    var lastInsertID int
    query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
    err := ri.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertID)
    if err != nil {
        return &model.User{}, err
    }

    user.ID = int64(lastInsertID)
    return user, nil
}

func (ri *repositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    u := model.User{}
    query := "SELECT id, username, email, password FROM users WHERE email = $1"
    err := ri.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
    if err != nil {
        return &model.User{}, nil
    }

    return &u, nil
}