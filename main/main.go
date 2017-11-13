package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

const VerifyMessage = "verified"

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		// url prefix가 "/login", "public/index.html"이면 auth를 체크하지 않음
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// "X_AUTH" 쿠키 값이 없으면 "/login"으로 이동
			c.Redirect("/login")
			return
		} else if err != nil {
			// err handling
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			// cookie 값으로 인증이 확인되면 다음 핸들러로 넘어감
			next(c)
			return
		}

		// "/login"으로 이동
		c.Redirect("/login")
	}
}

// auth token 확인
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

func main() {
	s := NewServer()

	// '/' 경로로 접속했을 때 처리할 핸들러 함수 지정
	s.HandleFunc("GET", "/", func(c *Context) {
		// 'welcome!' 문자열을 화면에 출력
		c.RenderTemplate("/public/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", "/users/:id", func(c *Context) {
		if c.Params["id"] == "0" {
			panic("id is zero")
		}
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
	})

	s.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *Context) {
		u := User{c.Params["user_id"].(string), c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	s.HandleFunc("POST", "/users", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, c.Params)
	})

	s.HandleFunc("GET", "/login", func(c *Context) {
		// "login.html" 렌더링
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "required login"})
	})

	s.HandleFunc("POST", "/login", func(c *Context) {
		// login 정보를 확인하여 쿠키에 인증 토근값 기록
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		// id와 password가 맞지 않으면 다시 "/login" 페이지 렌더링
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다."})
	})

	s.Use(AuthHandler)

	// 8080 포트로 웹서버 구동
	s.Run(":8080")
}

func CheckLogin(username, password string) bool {
	// login 처리
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)

	return username == USERNAME && password == PASSWORD
}

// auth token 생성
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return " "
	}
	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}
