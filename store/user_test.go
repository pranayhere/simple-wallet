package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domains"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func createRandomUser(t *testing.T) domains.User {
    userRepo := store.NewUserRepo(testDb)

    hashedPassword, err := util.HashPassword(util.RandomString(6))
    require.NoError(t, err)

    args := store.CreateUserParams{
        Username:       util.RandomUser(),
        Status:         domains.UserStatusACTIVE,
        FullName:       util.RandomUser(),
        Email:          util.RandomEmail(),
        HashedPassword: hashedPassword,
    }

    user, err := userRepo.CreateUser(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, user)

    require.Equal(t, args.Username, user.Username)
    require.Equal(t, args.Status, user.Status)
    require.Equal(t, args.Email, user.Email)
    require.Equal(t, args.HashedPassword, user.HashedPassword)
    require.Equal(t, args.FullName, user.FullName)

    require.NotZero(t, user.CreatedAt)

    return user
}

func TestCreateUser(t *testing.T) {
    createRandomUser(t)
}

func TestGetUser(t *testing.T) {
    userRepo := store.NewUserRepo(testDb)
    user1 := createRandomUser(t)
    user2, err := userRepo.GetUser(context.Background(), user1.Username)
    require.NoError(t, err)
    require.NotEmpty(t, user2)

    require.Equal(t, user1.Username, user2.Username)
    require.Equal(t, user1.Email, user2.Email)
    require.Equal(t, user1.Status, user2.Status)
    require.Equal(t, user1.HashedPassword, user2.HashedPassword)
    require.Equal(t, user1.FullName, user2.FullName)

    require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
    require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserStatus(t *testing.T) {
    userRepo := store.NewUserRepo(testDb)
    user1 := createRandomUser(t)

    args := store.UpdateUserStatusParams{
        ID:     user1.ID,
        Status: domains.UserStatusBLOCKED,
    }
    user2, err := userRepo.UpdateUserStatus(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, user2)

    require.Equal(t, user1.Username, user2.Username)
    require.Equal(t, user1.Email, user2.Email)
    require.Equal(t, domains.UserStatusBLOCKED, user2.Status)
    require.Equal(t, user1.HashedPassword, user2.HashedPassword)
    require.Equal(t, user1.FullName, user2.FullName)

    require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
    require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}