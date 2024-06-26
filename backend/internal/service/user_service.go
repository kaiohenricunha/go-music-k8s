package service

import (
	"errors"
	"log"

	"github.com/kaiohenricunha/go-music-k8s/backend/internal/dao"
	"github.com/kaiohenricunha/go-music-k8s/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")

	// ErrUsernameOrEmailTaken is returned when a username or email is already taken.
	ErrUsernameOrEmailTaken = errors.New("username or email already taken")

	// ErrInvalidCredentials is returned when the username or password is incorrect.
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// UserService outlines the interface for user-related operations.
type UserService interface {
	ValidateUser(username, password string) (uint, bool)
	RegisterUser(user *model.User) error
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}

type userService struct {
	userDAO dao.MusicDAO
}

func NewUserService(userDAO dao.MusicDAO) UserService {
	return &userService{userDAO: userDAO}
}

// ValidateUser checks if the username and password are correct.
func (us *userService) ValidateUser(username, password string) (uint, bool) {
	log.Printf("Validating user: %s\n", username)
	user, err := us.userDAO.GetUserByUsername(username)
	if err != nil || user == nil {
		log.Printf("User not found: %s\n", username)
		return 0, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed for user: %s\n", username)
		return 0, false
	}

	return user.ID, true
}

// RegisterUser handles registering a new user with hashed password.
func (us *userService) RegisterUser(user *model.User) error {
	// Check if username already exists
	existingUser, _ := us.GetUserByUsername(user.Username)
	if existingUser != nil {
		if user.Email == existingUser.Email {
			return ErrUsernameOrEmailTaken
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create the user
	return us.userDAO.CreateUser(user)
}

// GetAllUsers retrieves all users.
func (us *userService) GetAllUsers() ([]model.User, error) {
	return us.userDAO.GetAllUsers()
}

// GetUserByUsername retrieves a user by their username.
func (us *userService) GetUserByUsername(username string) (*model.User, error) {
	return us.userDAO.GetUserByUsername(username)
}
