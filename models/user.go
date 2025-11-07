package models

import "time"

type User struct {
	ID              int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email           string    `gorm:"column:email;uniqueIndex;not null;size:255" json:"email"`
	Password        string    `gorm:"column:password;not null;size:255" json:"-"`
	FirstName       string    `gorm:"column:first_name;not null;size:255" json:"first_name"`
	LastName        string    `gorm:"column:last_name;not null;size:255" json:"last_name"`
	Age             int       `gorm:"column:age;not null" json:"age"`
	BelieverCategory int      `gorm:"column:believer_category;not null" json:"believer_category"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName overrides the default pluralized table name
func (User) TableName() string {
	return "users"
}

// GetFullName returns the concatenated first and last name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

