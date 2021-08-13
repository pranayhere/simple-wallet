package constant

import "time"

const (
    AccessTokenDuration = 15 * time.Minute
    SymmetricKey        = "12345678901234567890123456789012"
)

const (
    AuthorizationHeaderKey  = "authorization"
    AuthorizationTypeBearer = "bearer"
    AuthorizationPayloadKey = "authorization_payload"
)
