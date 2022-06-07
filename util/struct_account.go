package util

type Login struct {
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
}

type ConfirmLoginData struct {
	Password   string `json:"comfirm_password"`
	Salt       string `json:"confirm_salt"`
	DecodeSalt []byte `json:"decode_password"`
}

type ArgonConfig struct {
	Time   uint32
	Memory uint32
	Thread uint8
	KeyLen uint32
}
