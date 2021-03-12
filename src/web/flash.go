package web

import (
	"encoding/base64"
	"net/http"
	"time"
)

func SetFlash(wr http.ResponseWriter, name string, message []byte) {
	cookie := &http.Cookie{Name: name, Value: EncodeBase64(message)}
	http.SetCookie(wr, cookie)
}

func GetFlash(wr http.ResponseWriter, req *http.Request, name string) ([]byte, error) {
	c, err := req.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		} else {
			return nil, err
		}
	}
	message, err := DecodeBase64(c.Value)
	if err != nil {
		return nil, err
	}
	delCookie := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(wr, delCookie)
	return message, nil
}

func EncodeBase64(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func DecodeBase64(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}
