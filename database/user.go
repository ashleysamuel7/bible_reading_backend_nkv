package database

import (
	"bible_reading_backend_nkv/models"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (c Client) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := c.DB.WithContext(ctx).Create(user)
	return result.Error
}

func (c Client) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	result := c.DB.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (c Client) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := c.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (c Client) UpdateUser(ctx context.Context, id int, updates map[string]interface{}) error {
	result := c.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates)
	return result.Error
}

func (c Client) DeleteUser(ctx context.Context, id int) error {
	result := c.DB.WithContext(ctx).Delete(&models.User{}, id)
	return result.Error
}

func (c Client) VerifyPassword(ctx context.Context, email, password string) (*models.User, error) {
	user, err := c.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

