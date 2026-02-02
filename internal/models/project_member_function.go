package models

// ProjectMemberFunction lie un membre du projet à une fonction (many-to-many).
// Un membre peut avoir plusieurs fonctions (ex. Dev + Lead, Testeur + Chef de projet).
// Table: project_member_functions
type ProjectMemberFunction struct {
	ProjectMemberID  uint `gorm:"primaryKey" json:"project_member_id"`
	ProjectFunctionID uint `gorm:"primaryKey" json:"project_function_id"`

	ProjectMember  *ProjectMember  `gorm:"foreignKey:ProjectMemberID" json:"-"`
	ProjectFunction *ProjectFunction `gorm:"foreignKey:ProjectFunctionID" json:"-"`
}

// TableName spécifie le nom de la table
func (ProjectMemberFunction) TableName() string {
	return "project_member_functions"
}
