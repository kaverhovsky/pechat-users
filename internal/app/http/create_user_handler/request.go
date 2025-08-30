package create_user_handler

type createUserRequest struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}
