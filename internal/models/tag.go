package models

var (
	TagCreatableFields = []string{"name"}
	TagUpdatableFields = []string{"name"}
)

// Tag ...
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetFields ...
func (t *Tag) GetFields() []interface{} {
	return []interface{}{
		&t.ID,
		&t.Name,
	}
}
