package update_user_handler

type UpdateUserRequest struct {
	Nickname  *string `json:"nickname"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Bio       *string `json:"bio"`
}
