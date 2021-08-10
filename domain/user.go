package domain

import (
    "fmt"
    "time"
)

type UserStatus string

const (
    UserStatusACTIVE  UserStatus = "ACTIVE"
    UserStatusBLOCKED UserStatus = "BLOCKED"
)

type User struct {
    ID                int64      `json:"id"`
    Username          string     `json:"username"`
    HashedPassword    string     `json:"hashed_password"`
    Status            UserStatus `json:"status"`
    FullName          string     `json:"full_name"`
    Email             string     `json:"email"`
    PasswordChangedAt time.Time  `json:"password_changed_at"`
    CreatedAt         time.Time  `json:"created_at"`
    UpdatedAt         time.Time  `json:"updated_at"`
}

func (e *UserStatus) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = UserStatus(s)
    case string:
        *e = UserStatus(s)
    default:
        return fmt.Errorf("unsupported scan type for UserStatus: %T", src)
    }
    return nil
}
