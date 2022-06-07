package router

import (
	"FD/util"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData util.Login
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
	}

	var confirmData util.ConfirmLoginData
	err = db.QueryRow("SELECT password, salt FROM user WHERE login_id=?", loginData.LoginId).
		Scan(&confirmData.Password, &confirmData.Salt)

	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: "id error",
			Err:  false,
		})
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, string(resData))
		return
	}

	confirmData.DecodeSalt, _ = hex.DecodeString(confirmData.Salt)

	encodedPwd := hex.EncodeToString(
		argon2.IDKey([]byte(loginData.Password), confirmData.DecodeSalt,
			argonConfig.Time, argonConfig.Memory, argonConfig.Thread, argonConfig.KeyLen))

	if encodedPwd != confirmData.Password {
		w.WriteHeader(http.StatusUnauthorized)
		resData, _ := json.Marshal(util.Res{
			Data: "password error",
			Err:  false,
		})
		fmt.Fprint(w, string(resData))
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: "login sucess",
		Err:  false,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
