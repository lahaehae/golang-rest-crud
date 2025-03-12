package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/lahaehae/crud_project/internal/models"
	"github.com/lahaehae/crud_project/internal/repository/mocks"
	"github.com/lahaehae/crud_project/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T){
	mockRepo := new(mocks.UserRepo) //repository mock

	userService := service.NewUserService(mockRepo)

	testUser := &models.User{
		Id: 1,
		Name: "Test User",
		Email: "test@example.com",
		Balance: 1000,
	}

	mockRepo.On("GetUser", mock.Anything, int64(1)).Return(testUser, nil)

	t.Run("successfuly get user", func(t *testing.T){
		user, err := userService.GetUser(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.Id, user.Id)
		assert.Equal(t, testUser.Name, user.Name)
		fmt.Println(user)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetUser", mock.Anything, int64(2)).Return(nil, errors.New("user not found"))
		user, err := userService.GetUser(context.Background(), 2)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}