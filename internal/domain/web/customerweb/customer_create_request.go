package customerweb

type CreateRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Password string `json:"password"`
}
