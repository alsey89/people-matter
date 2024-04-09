package auth

import (
	"github.com/alsey89/hrms/internal/schema"
)

func SignupService(email string, password string) (*schema.User, error) {

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, fmt.Errorf("auth.s.signup: %w", err)
	// }

	// newUser := schema.User{
	// 	FirstName: "New",
	// 	LastName:  "User",
	// 	Email:     email,
	// 	Password:  string(hashedPassword),
	// 	Role:      "user",
	// }

	// // createdUser, err := CreateNewUser(&newUser) //* this adds the ID to newUser
	// // if err != nil {
	// // 	return nil, fmt.Errorf("auth.s.signup: %w", err)
	// // }

	// return createdUser, nil
	return nil, nil
}

func SigninService(email string, password string) (*schema.User, error) {

	// user, err := GetUserByEmail(email)
	// if err != nil {
	// 	return nil, fmt.Errorf("auth.s.signin: %w", err)
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	if err == bcrypt.ErrMismatchedHashAndPassword {
	// 		return nil, fmt.Errorf("auth.s.signin: %w", ErrInvalidCredentials)
	// 	}
	// 	return nil, fmt.Errorf("auth.s.signin: %w", err)
	// }

	return nil, nil
}
