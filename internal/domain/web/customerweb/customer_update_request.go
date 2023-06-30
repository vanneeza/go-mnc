package customerweb

type UpdateRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Password string `json:"password"`
}
