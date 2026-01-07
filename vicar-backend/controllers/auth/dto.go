package auth

type registerDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginTotpDto struct {
	State string `json:"state"`
	Code  string `json:"code"`
}
