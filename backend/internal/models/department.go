package models

// Department groups users by business function.
type Department struct {
	BaseModel
	Name        string `gorm:"size:150;uniqueIndex" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	ManagerID   *uint  `json:"managerId"`
	Manager     *User  `json:"manager"`
}
