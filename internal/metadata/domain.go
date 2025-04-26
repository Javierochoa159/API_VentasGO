package metadata

// Metadata represents a system sale with metadata for auditing and versioning.
type Metadata struct {
	Quantity     int     `json:"quantity"`
	Approved     int     `json:"approve"`
	Pending      int     `json:"pending"`
	Rejected     int     `json:"rejected"`
	Total_amount float32 `json:"total_amount"`
}
