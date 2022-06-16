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

func (s SignUp) IsValidLen() (isOk bool) {
	if s.ClubId == 0 ||
		len(s.Email) < 2 ||
		len(s.LoginId) < 2 ||
		len(s.Password) < 2 ||
		len(s.PhoneNum) < 2 ||
		len(s.UserName) < 2 {
		return false
	}
	return true
}

type UserList struct {
	UserId   int    `json:"user_id"`
	ClubId   int    `json:"club_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	PhoneNum string `json:"phone_num"`
}

type ApplyAdmin struct {
	UserId   int    `json:"user_id"`
	ClubId   int    `json:"club_id"`
	UserName string `json:"user_name"`
}
