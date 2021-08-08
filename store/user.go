package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domains"
)

type UserRepo interface {
    CreateUser(ctx context.Context, arg CreateUserParams) (domains.User, error)
    GetUser(ctx context.Context, id int64) (domains.User, error)
    GetUserByUsername(ctx context.Context, username string) (domains.User, error)
    UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) (domains.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepo(client *sql.DB) UserRepo {
    return &userRepository{
        db: client,
    }
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    status,
    full_name,
    email
) values (
$1, $2, $3, $4, $5
) RETURNING id, username, hashed_password, status, full_name, email, password_changed_at, created_at, updated_at
`

type CreateUserParams struct {
    Username       string             `json:"username"`
    HashedPassword string             `json:"hashed_password"`
    Status         domains.UserStatus `json:"status"`
    FullName       string             `json:"full_name"`
    Email          string             `json:"email"`
}

func (q *userRepository) CreateUser(ctx context.Context, arg CreateUserParams) (domains.User, error) {
    row := q.db.QueryRowContext(ctx, createUser,
        arg.Username,
        arg.HashedPassword,
        arg.Status,
        arg.FullName,
        arg.Email,
    )
    var i domains.User
    err := row.Scan(
        &i.ID,
        &i.Username,
        &i.HashedPassword,
        &i.Status,
        &i.FullName,
        &i.Email,
        &i.PasswordChangedAt,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getUserByUsername = `-- name: getUserByUsername :one
SELECT id, username, hashed_password, status, full_name, email, password_changed_at, created_at, updated_at from users
where username = $1 LIMIT 1
`

func (q *userRepository) GetUserByUsername(ctx context.Context, username string) (domains.User, error) {
    row := q.db.QueryRowContext(ctx, getUserByUsername, username)
    var i domains.User
    err := row.Scan(
        &i.ID,
        &i.Username,
        &i.HashedPassword,
        &i.Status,
        &i.FullName,
        &i.Email,
        &i.PasswordChangedAt,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getUser = `-- name: getUser :one
SELECT id, username, hashed_password, status, full_name, email, password_changed_at, created_at, updated_at from users
where id = $1 LIMIT 1
`

func (q *userRepository) GetUser(ctx context.Context, id int64) (domains.User, error) {
    row := q.db.QueryRowContext(ctx, getUser, id)
    var i domains.User
    err := row.Scan(
        &i.ID,
        &i.Username,
        &i.HashedPassword,
        &i.Status,
        &i.FullName,
        &i.Email,
        &i.PasswordChangedAt,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const updateUserStatus = `-- name: UpdateUserStatus :one
UPDATE users
set Status = $1
where id = $2
RETURNING id, username, hashed_password, status, full_name, email, password_changed_at, created_at, updated_at
`

type UpdateUserStatusParams struct {
    Status domains.UserStatus `json:"status"`
    ID     int64              `json:"id"`
}

func (q *userRepository) UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) (domains.User, error) {
    row := q.db.QueryRowContext(ctx, updateUserStatus, arg.Status, arg.ID)
    var i domains.User
    err := row.Scan(
        &i.ID,
        &i.Username,
        &i.HashedPassword,
        &i.Status,
        &i.FullName,
        &i.Email,
        &i.PasswordChangedAt,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}
