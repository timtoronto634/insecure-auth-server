package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var sessionList = make(map[string]Session)
var templates = make(map[string]*template.Template)
var AuthCodeList = make(map[string]AuthCode)
var TokenCodeList = make(map[string]TokenCode)

func Serve() {
	var err error
	templates["login"], err = template.ParseFiles("login.html")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("start oauth server on localhost:8081...")
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/authcheck", authCheck)
	http.HandleFunc("/token", token)
	http.HandleFunc("/certs", certs)
	http.HandleFunc("/userinfo", userinfo)
	err = http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func certs(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(makeJWK())
}

func userinfo(w http.ResponseWriter, req *http.Request) {
	h := req.Header.Get("Authorization")
	tmp := strings.Split(h, " ")

	// トークンがあるか確認
	v, ok := TokenCodeList[tmp[1]]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("token is wrong.\n")))
		return
	}

	// トークンの有効期限が切れてないか
	if v.expires_at < time.Now().Unix() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("token is expire.\n")))
		return
	}

	// スコープが正しいか、openid profileで固定
	if v.scopes != "openid profile" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("scope is not permit.\n")))
		return
	}

	// ユーザ情報を返す
	var m = map[string]interface{}{
		"sub":         user.sub,
		"name":        user.name_ja,
		"given_name":  user.given_name,
		"family_name": user.family_name,
		"locale":      user.locale,
	}
	buf, _ := json.MarshalIndent(m, "", "  ")
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}
