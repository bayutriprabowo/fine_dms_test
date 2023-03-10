package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"enigmacamp.com/fine_dms/model"
	"enigmacamp.com/fine_dms/repository"
)

func IsValidTag(tagName string) bool {
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(tagName) {
		return false
	}

	words := strings.Fields(tagName)
	if len(words) > 16 {
		return false
	}
	return true
}

func IsValidInput(user *model.User) error {
	if user.Username == "" {
		return errors.New("Username cannot be empty")
	}
	if user.Email == "" {
		return errors.New("Email cannot be empty")
	}
	if user.Password == "" {
		return errors.New("Password cannot be empty")
	}

	return nil
}

func IsDuplicateUsername(userRepo repository.UserRepository, user *model.User) error {
	existingUsers, err := userRepo.SelectUser()
	if err != nil {
		return err
	}

	for _, existingUser := range existingUsers {
		if existingUser.Username == user.Username {
			return fmt.Errorf("username %s already exists", user.Username)
		}
	}

	return nil
}

func IsDuplicateEmail(userRepo repository.UserRepository, user *model.User) error {
	existingUsers, err := userRepo.SelectUser()
	if err != nil {
		return err
	}

	for _, existingUser := range existingUsers {
		if existingUser.Email == user.Email {
			return fmt.Errorf("email %s already exists", user.Email)
		}
	}

	return nil
}
