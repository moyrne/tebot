package models

// QGroup QQ Group
type QGroup struct {
	ID   int    `json:"id"`
	QGID int    `json:"qgid"`
	Name string `json:"name"`
}

func (g QGroup) TableName() string {
	return "q_group"
}
