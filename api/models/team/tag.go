package team

// TagRequest is the body for POST /api/team/tags.
type TagRequest struct {
	Name string `json:"name"`
}

// TagResponse is the shape returned by both POST and GET /api/team/tags -
// see Services::TeamTags#create/#list in customer-console.
type TagResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
