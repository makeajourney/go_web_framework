package main

import (
	"fmt"
	"time"
)

type User struct {
	Id        string
	AddressId string
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

	// 8080 포트로 웹서버 구동
	s.Run(":8080")
}
