package sale

import "time"

// User represents a system sale with metadata for auditing and versioning.
type Sale struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Amount    float32   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

// UpdateFields represents the optional fields for updating a Sale.
// A nil pointer means “no change” for that field.
type UpdateFields struct {
	Status *string `json:"status"`
}

type Metadata struct {
	Quantity        int    		`json:"quantity"`
	Approved    	int    		`json:"approve"`
	Pending 		string 		`json:"pending"` 
	Rejected 		string 		`json:"rejected"`
	Total_amount    float32   	`json:"total_amount"`
}