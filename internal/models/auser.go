package models

// AUser Admin User
type AUser struct {
}

func (u AUser) TableName() string {
	return "a_user"
}
