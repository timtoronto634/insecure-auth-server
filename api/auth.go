package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	SCOPE                 = "https://www.googleapis.com/auth/photoslibrary.readonly"
	AUTH_CODE_DURATION    = 300
	ACCESS_TOKEN_DURATION = 3600
)

type Client struct {
	id          string
	name        string
	redirectURL string
	secret      string
}

type info struct {
	id          string
	name        string
	password    string
	sub         string
	name_ja     string
	given_name  string
	family_name string
	locale      string
}

var clientInfo = Client{
	id:          "client_id_xxx",
	name:        "atsu",
	redirectURL: "localhost:3000",
	secret:      "hiro",
}

type Session struct {
	client                string
	state                 string
	scopes                string
	redirectUri           string
	code_challenge        string
	code_challenge_method string
}

type AuthCode struct {
	user         string
	clientId     string
	scopes       string
	redirect_uri string
	expires_at   int64
}

type TokenCode struct {
	user       string
	clientId   string
	scopes     string
	expires_at int64
}

func auth(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	requiredParameter := []string{"response_type", "client_id", "redirect_uri"} // scope, stateは一旦なし
	// 必須パラメータのチェック
	for _, v := range requiredParameter {
		if !query.Has(v) {
			log.Printf("%s is missing", v)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("invalid_request. %s is missing", v)))
			return
		}
	}
	// client id の一致確認
	if clientInfo.id != query.Get("client_id") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("client_id is not match"))
		return
	}
	// レスポンスタイプはいったん認可コードだけをサポート
	if "code" != query.Get("response_type") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("only support code"))
		return
	}
	sessionId := uuid.New().String()
	// セッションを保存しておく
	session := Session{
		client:                query.Get("client_id"),
		state:                 query.Get("state"),
		scopes:                query.Get("scope"),
		redirectUri:           query.Get("redirect_uri"),
		code_challenge:        query.Get("code_challenge"),
		code_challenge_method: query.Get("code_challenge_method"),
	}
	sessionList[sessionId] = session

	// CookieにセッションIDをセット
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionId,
	}
	http.SetCookie(w, cookie)

	// ログイン&権限認可の画面を戻す
	if err := templates["login"].Execute(w, struct {
		ClientId string
		Scope    string
	}{
		ClientId: session.client,
		Scope:    session.scopes,
	}); err != nil {
		log.Println(err)
	}
	log.Println("return login page...")

}

type userInfo struct {
	name     string
	password string
}

var user = userInfo{
	name:     "atsu",
	password: "hiro",
}

// 認可レスポンスを返す
func authCheck(w http.ResponseWriter, req *http.Request) {

	loginUser := req.FormValue("username")
	password := req.FormValue("password")

	if loginUser != user.name || password != user.password {
		w.Write([]byte("login failed"))
	} else {

		cookie, _ := req.Cookie("session")
		http.SetCookie(w, cookie)
		v, _ := sessionList[cookie.Value]

		authCodeString := uuid.New().String()
		authData := AuthCode{
			user:         loginUser,
			clientId:     v.client,
			scopes:       v.scopes,
			redirect_uri: v.redirectUri,
			expires_at:   time.Now().Unix() + 300,
		}
		// 認可コードを保存
		AuthCodeList[authCodeString] = authData

		log.Printf("auth code accepet : %s\n", authData)

		location := fmt.Sprintf("%s?code=%s&state=%s", v.redirectUri, authCodeString, v.state)
		w.Header().Add("Location", location)
		w.WriteHeader(302)

	}

}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	IdToken     string `json:"id_token,omitempty"`
}

// トークンを発行するエンドポイント
func token(w http.ResponseWriter, req *http.Request) {

	cookie, _ := req.Cookie("session")
	req.ParseForm()
	query := req.Form

	requiredParameter := []string{"grant_type", "code", "client_id", "redirect_uri"}
	// 必須パラメータのチェック
	for _, v := range requiredParameter {
		if !query.Has(v) {
			log.Printf("%s is missing", v)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("invalid_request. %s is missing\n", v)))
			return
		}
	}

	// 認可コードフローだけサポート
	if "authorization_code" != query.Get("grant_type") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid_request. not support type.\n")))
	}

	// 保存していた認可コードのデータを取得。なければエラーを返す
	v, ok := AuthCodeList[query.Get("code")]
	if !ok {
		log.Println("auth code isn't exist")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("no authrization code")))
	}

	// 認可リクエスト時のクライアントIDと比較
	if v.clientId != query.Get("client_id") {
		log.Println("client_id not match")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid_request. client_id not match.\n")))
	}

	// 認可リクエスト時のリダイレクトURIと比較
	if v.redirect_uri != query.Get("redirect_uri") {
		log.Println("redirect_uri not match")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid_request. redirect_uri not match.\n")))
	}

	// 認可コードの有効期限を確認
	if v.expires_at < time.Now().Unix() {
		log.Println("authcode expire")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid_request. auth code time limit is expire.\n")))
	}

	// clientシークレットの確認
	if clientInfo.secret != query.Get("client_secret") {
		log.Println("client_secret is not match.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("invalid_request. client_secret is not match.\n")))
	}

	// PKCEのチェック
	// clientから送られてきたverifyをsh256で計算&base64urlエンコードしてから
	// 認可リクエスト時に送られてきてセッションに保存しておいたchallengeと一致するか確認
	session := sessionList[cookie.Value]
	hash := sha256.Sum256([]byte(query.Get("code_verifier")))
	if session.code_challenge != base64.RawURLEncoding.EncodeToString(hash[:]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("PKCE check is err..."))
	}

	tokenString := uuid.New().String()
	expireTime := time.Now().Unix() + ACCESS_TOKEN_DURATION

	tokenInfo := TokenCode{
		user:       v.user,
		clientId:   v.clientId,
		scopes:     v.scopes,
		expires_at: expireTime,
	}
	TokenCodeList[tokenString] = tokenInfo
	// 認可コードを削除
	delete(AuthCodeList, query.Get("code"))

	tokenResp := TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expireTime,
	}
	resp, err := json.Marshal(tokenResp)
	if err != nil {
		log.Println("json marshal err")
	}

	log.Printf("token ok to client %s, token is %s", v.clientId, string(resp))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
