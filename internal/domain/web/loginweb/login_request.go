package loginweb

type Request struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
