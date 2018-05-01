package domain

// UserUseCase ...
type UserUseCase struct {
	userRepository UserRepository
}

// NewUserUseCase ...
func NewUserUseCase(repository UserRepository) UserUseCase {
	return UserUseCase{
		userRepository: repository,
	}
}

// Register Creating Email Verification Token
func (useCase *UserUseCase) Register(newUser NewUser) (string, error) {
	return useCase.userRepository.StoreUser(newUser)
}

// ActivateAccount ...
func (useCase *UserUseCase) ActivateAccount(token string) error {
	user, err := useCase.userRepository.GetUserByToken(token)
	if err != nil {
		return err
	}
	user.IsActive = true
	user.IsEmailVerified = true
	if _, err := useCase.userRepository.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

// IsEmailUsed ...
func (useCase *UserUseCase) IsEmailUsed(email string) bool {
	_, err := useCase.userRepository.GetUserByEmail(email)
	return err == nil
}

// DoLogin ...
func (useCase *UserUseCase) DoLogin(email string, password string) (*JwtToken, error) {
	return useCase.userRepository.CreateJWTToken(email, password)
}

// RetrieveUserByToken ...
func (useCase *UserUseCase) RetrieveUserByToken(token string) (*User, error) {
	return useCase.userRepository.GetUserByToken(token)
}
