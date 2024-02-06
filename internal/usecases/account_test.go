package usecases

import usecases_test "github.com/alexsibrin/runbot-auth/internal/usecases/mocks"

type Mocks struct {
	repo   usecases_test.MockIAccountRepo
	hasher usecases_test.MockIPasswordHasher
}

type TestCase struct {
	Name           string
	Request        string
	mocksFn        func(Mocks)
	ResponseStatus int
	ResponseBody   string
}
