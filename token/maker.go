package token

import "time"

type Maker interface {
    // CreateToken creates a new token for specified username and duration
    CreateToken(userID int64, duration time.Duration) (string, error)

    // VerifyToken checks if the token is valid or not
    VerifyToken(token string) (*Payload, error)
}
