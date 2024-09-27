package mocks

//
// import (
//	"time"
//
//	"github.com/DKhorkov/medods/entities"
//
//	customerrors "github.com/DKhorkov/medods/internal/errors"
//)
//
// type MockedAuthRepository struct {
//	UsersStorage map[int]*entities.User
//}
//
// func (repo *MockedAuthRepository) RegisterUser(userData entities.RegisterUserDTO) (int, error) {
//	var user entities.User
//	user.Email = userData.Credentials.Email
//	user.ID = len(repo.UsersStorage) + 1
//	user.CreatedAt = time.Now()
//	user.UpdatedAt = time.Now()
//
//	repo.UsersStorage[user.ID] = &user
//	return user.ID, nil
//}
//
// func (repo *MockedAuthRepository) GetUserByID(id int) (*entities.User, error) {
//	user := repo.UsersStorage[id]
//	if user != nil {
//		return user, nil
//	}
//
//	return nil, &customerrors.UserNotFoundError{}
//}
//
// func (repo *MockedAuthRepository) GetAllUsers() ([]*entities.User, error) {
//	var users []*entities.User
//	for _, user := range repo.UsersStorage {
//		users = append(users, user)
//	}
//
//	return users, nil
//}
//
// func (repo *MockedAuthRepository) GetUserByEmail(email string) (*entities.User, error) {
//	for _, user := range repo.UsersStorage {
//		if user.Email == email {
//			return user, nil
//		}
//	}
//
//	return nil, &customerrors.UserNotFoundError{}
//}
