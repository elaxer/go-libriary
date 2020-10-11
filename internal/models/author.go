package models

// Author ...
type Author struct {
	ID         int    `json:"id,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	Image      string `json:"image,omitempty"`
	Biography  string `json:"biography,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
}

// GetFields ...
func (a *Author) GetFields() []interface{} {
	return []interface{}{
		&a.ID,
		&a.FullName,
		&a.Image,
		&a.Biography,
		&a.FirstName,
		&a.LastName,
		&a.MiddleName,
	}
}
