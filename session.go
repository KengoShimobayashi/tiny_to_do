package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"
)

const cookieSessionId = "sessionId"

func ensureSession(w http.ResponseWriter, r *http.Request)(string, error) {
	c, err := r.Cookie(cookieSessionId)

	// エラー処理
	// リクエストにsessionIdがCookieとして存在しない場合、セッションを開始
	if err == http.ErrNoCookie {
		sessionId, err := startSession(w)
		return sessionId, err
	}

	// CookieとしてsessionIdが存在する場合、その値を返す
	if err == nil{
		sessionId := c.Value
		return sessionId, nil
	}

	return "", nil
}

func startSession(w http.ResponseWriter)(string, error){
	sessionId, err := makeSessionId()

	// エラー処理
	if err != nil{
		return "", err
	}

	cookie := &http.Cookie{
		Name: cookieSessionId,
		Value: sessionId,
		Expires: time.Now().Add(1800 * time.Second),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return sessionId, nil
}

func makeSessionId() (string, error){
	randBytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil{
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randBytes), nil
}