package captcha

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	defaultChars = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	defaultStr   = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

const (
	// default captcha attributes
	challengeNums = 4
	expiration    = 600
	CaptchaKey    = "a188528385820ec0e3107bd895201d8c"
	CookieName    = "captcha_id"
)

// Captcha struct
type Captcha struct {

	// captcha image width and height
	StdWidth  int
	StdHeight int

	// captcha chars nums
	ChallengeNums int

	// captcha key
	CaptchaKey string

	// captcha expiration seconds
	Expiration int
}

// generate rand chars with default chars
func (c *Captcha) genRandChars() []byte {
	return RandomCreateBytes(c.ChallengeNums, defaultChars...)
}

var verifyErr = errors.New("verify failed")

// verify from a request
func (c *Captcha) VerifyReq(req *http.Request, challenge string) error {
	cookie, err := req.Cookie(CookieName)
	if err != nil {
		return err
	}
	if !c.Verify(cookie.Value, challenge) {
		return verifyErr
	}
	return nil
}

// direct verify id and challenge string
// id == MD5(challenge+CaptchaKey)
func (c *Captcha) Verify(id string, challenge string) (success bool) {
	if len(challenge) == 0 || len(id) == 0 {
		return
	}

	return id == c.calcId(challenge)
}

func (c *Captcha) calcId(challenge string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(challenge + c.CaptchaKey))
	return base64.StdEncoding.EncodeToString(md5Ctx.Sum(nil))
}

func (c *Captcha) ServeHttp(w http.ResponseWriter) error {
	cand := c.genRandChars()
	str := make([]rune, len(cand))
	for i := range cand {
		str[i] = defaultStr[int(cand[i])]
	}
	img := NewImage(string(str), cand, c.StdWidth, c.StdHeight)
	cookie := &http.Cookie{
		Name:   CookieName,
		Value:  c.calcId(string(str)),
		Path:   "/",
		MaxAge: c.Expiration,
		Domain: ".codoon.com",
	}
	http.SetCookie(w, cookie)
	_, err := img.WriteTo(w)
	return err
}

// create a new captcha.Captcha
func NewCaptcha() *Captcha {
	cpt := &Captcha{}
	cpt.ChallengeNums = challengeNums
	cpt.Expiration = expiration
	cpt.StdWidth = StdWidth
	cpt.StdHeight = StdHeight
	cpt.CaptchaKey = CaptchaKey

	return cpt
}
