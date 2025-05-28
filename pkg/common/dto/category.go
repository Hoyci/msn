package dto

type Category struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Icon  string        `json:"icon"`
	Users *int          `json:"users,omitempty"`
	Subs  []Subcategory `json:"subs,omitempty"`
}

type Subcategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
}
