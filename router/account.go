package router

import (
	"FD/util"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/go-session/session/v3"
	"golang.org/x/crypto/argon2"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if userId := util.LoginCheck(w, r); userId == nil {
		util.GlobalErr("already login", nil, 400, w)
		return
	}

	var loginData util.Login
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		resData, _ := json.Marshal(util.Res{
			Data: "data isn't json",
			Err:  true,
		})
		fmt.Fprint(w, string(resData))
	}

	var confirmData util.ConfirmLoginData
	var userId int
	err = db.QueryRow("SELECT user_id, password, salt FROM \"user\" WHERE login_id=$1", loginData.LoginId).
		Scan(&userId, &confirmData.Password, &confirmData.Salt)

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

	store, err := session.Start(ctx, w, r)
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: nil,
			Err:  true,
		})
		w.WriteHeader(500)
		fmt.Fprint(w, string(resData))
		return
	}

	store.Set("user_id", userId)
	err = store.Save()
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: "saving session error",
			Err:  true,
		})
		w.WriteHeader(500)
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpData util.SignUp
	err := json.NewDecoder(r.Body).Decode(&signUpData)
	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: "data isn't json",
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	if signUpData.ClubId == 0 ||
		len(signUpData.Email) < 2 ||
		len(signUpData.LoginId) < 2 ||
		len(signUpData.Password) < 2 ||
		len(signUpData.PhoneNum) < 2 ||
		len(signUpData.UserName) < 2 {
		resData, _ := json.Marshal(util.Res{
			Data: "not enough values",
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	salt := make([]byte, 32)
	rand.Read(salt)
	encryptedPwd := argon2.IDKey([]byte(signUpData.Password), salt, argonConfig.Time, argonConfig.Memory, argonConfig.Thread, argonConfig.KeyLen)

	_, err = db.Exec("INSERT INTO user (club_id, user_name, email, login_id, password, phone_num, salt) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		signUpData.ClubId, signUpData.UserName, signUpData.Email, signUpData.LoginId, hex.EncodeToString(encryptedPwd), signUpData.PhoneNum, hex.EncodeToString(salt))

	if err != nil {
		log.Println(err)
		resData, _ := json.Marshal(util.Res{
			Data: "cannot sign up",
			Err:  true,
		})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(resData))
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: "success",
		Err:  false,
	})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(resData))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if userId := util.LoginCheck(w, r); userId != nil {
		util.GlobalErr("didn't login", nil, 400, w)
		return
	}

	err := session.Destroy(ctx, w, r)
	if err != nil {
		util.GlobalErr("cannot logout", err, 500, w)
		return
	}

	resData, _ := json.Marshal(util.Res{
		Data: nil,
		Err:  false,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(resData))
}
