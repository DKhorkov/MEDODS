package mocks

type MockedUsersRepository struct {
}

func (repo *MockedUsersRepository) GetUserEmail(guid string) (string, error) {
	return "alexqwerty35@yandex.ru", nil
}
