package grant

type Response struct {
	AccessToken string `json:"accessToken"`
}

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
