package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GlobalErr(data string, err error, code int, w http.ResponseWriter) {
	log.Println(err)
	resData, _ := json.Marshal(Res{
		Data: data,
		Err:  true,
	})
	w.WriteHeader(code)
	fmt.Fprint(w, string(resData))
}

func SessionErr(w http.ResponseWriter) {
	resData, _ := json.Marshal(Res{
		Data: "session error",
		Err:  true,
	})
	w.WriteHeader(500)
	fmt.Fprint(w, string(resData))
}

func LoginErr(w http.ResponseWriter) {
	resData, _ := json.Marshal(Res{
		Data: "need login",
		Err:  true,
	})
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, string(resData))
}
