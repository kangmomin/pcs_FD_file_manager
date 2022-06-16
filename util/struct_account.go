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

type SignUp struct {
	ClubId   int    `json:"club_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	LoginId  string `json:"login_id"`
	Password string `json:"password"`
	PhoneNum string `json:"phone_num"`
}

type UserList struct {
	UserId   int    `json:"user_id"`
	ClubId   int    `json:"club_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phone_num"`
}
