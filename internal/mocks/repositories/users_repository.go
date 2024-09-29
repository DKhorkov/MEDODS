package mocks

type MockedUsersRepository struct {
}

func (repo *MockedUsersRepository) GetUserEmail(guid string) (string, error) {
	return "example@yandex.ru", nil
}
