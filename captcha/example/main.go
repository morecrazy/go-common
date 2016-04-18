package main

import (
	"backend/common/captcha"
	"net/http"
	"strings"
)

var cpt = captcha.NewCaptcha()

func GenCaptcha(w http.ResponseWriter, req *http.Request) {
	cpt.ServeHttp(w)
}

func Verify(w http.ResponseWriter, req *http.Request) {
	if !strings.Contains(req.Host, "codoon.com") {
		w.Write([]byte("请使用codoon.com域名访问"))
		return
	}
	req.ParseForm()
	cap := req.FormValue("captcha")
	if err := cpt.VerifyReq(req, cap); err != nil {
		w.Write([]byte("验证失败"))
	} else {
		w.Write([]byte("验证成功"))
	}
}

var tpl = `<!DOCTYPE html><html><body>
<form action="/verify" method="post">
	<img onclick="this.src='/captcha?v='+Math.random()" class="captcha-img" src="/captcha"></a>
	<input name="captcha" type="text">
	<input value="submit" type="submit">
</form></body></html>
`

func Test(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(tpl))
}

func main() {
	http.HandleFunc("/captcha", GenCaptcha)
	http.HandleFunc("/verify", Verify)
	http.HandleFunc("/", Test)
	http.ListenAndServe(":8001", nil)
}
