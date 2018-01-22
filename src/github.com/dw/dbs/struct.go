package dbs

// This pacakage data structure about dbs_ prefix
import (
	"time"
)

const (
	// UserTable is the user table name in database
	UserTable = "dbs_user"
)

// User data structure
type User struct {
	UID         int64     `json:"uID" gorm:"column:uid"`
	Nickname    string    `json:"nickname" gorm:"column:nickname"`
	Mobile      string    `json:"mobile" gorm:"column:mobile"`
	Email       string    `json:"email" gorm:"column:email"`
	RegisterAt  time.Time `json:"registerAt"  gorm:"column:reg_time"`
	LastLoginAt time.Time `json:"lastLoginAt"  gorm:"column:last_login_time"`
}

// TableName return the table name in database for gorm using
func (s User) TableName() string {
	return UserTable
}
