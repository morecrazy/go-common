package common

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func IsStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

//check the item whether in object
//object must be slice
func HasItem(obj interface{}, item interface{}) bool {
	for i := 0; i < reflect.ValueOf(obj).Len(); i++ {
		if reflect.DeepEqual(reflect.ValueOf(obj).Index(i).Interface(), item) {
			return true
		}
	}
	return false
}

//Be careful to use, from,to must be pointer
func DumpStruct(to interface{}, from interface{}) {
	fromv := reflect.ValueOf(from)
	tov := reflect.ValueOf(to)
	if fromv.Kind() != reflect.Ptr || tov.Kind() != reflect.Ptr {
		return
	}

	from_val := reflect.Indirect(fromv)
	to_val := reflect.Indirect(tov)

	for i := 0; i < from_val.Type().NumField(); i++ {
		fdi_from_val := from_val.Field(i)
		fd_name := from_val.Type().Field(i).Name
		fdi_to_val := to_val.FieldByName(fd_name)

		if fdi_to_val.IsValid() && fdi_from_val.Type() == fdi_to_val.Type() {
			fdi_to_val.Set(fdi_from_val)
		}
	}
}

func DumpList(to interface{}, from interface{}) {
	val_from := reflect.ValueOf(from)
	val_to := reflect.ValueOf(to)

	if val_from.Type().Kind() == reflect.Slice && val_to.Type().Kind() == reflect.Slice &&
		val_from.Len() == val_to.Len() {
		for i := 0; i < val_from.Len(); i++ {
			DumpStruct(val_to.Index(i).Addr().Interface(), val_from.Index(i).Addr().Interface())
		}
	}
}

func ParseToLocalTime(str string) time.Time {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", str, time.Now().Location())
	return tm
}

func TimeToLocalFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ValidateLessEqChar(str string, n int) bool {
	return utf8.RuneCountInString(str) <= n
}

func ValidateNum(str string) bool {
	_, err := strconv.ParseInt(str, 10, 64)
	return err == nil
}

func ValidateEmail(str string) bool {
	reg := "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}" +
		"-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]" +
		"|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)" +
		"*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]" +
		"|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09" +
		"\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20" +
		"|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-" +
		"\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-" +
		"\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*" +
		"([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x" +
		"{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}" +
		"-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\" +
		"x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

	Re := regexp.MustCompile(reg)
	return Re.MatchString(str)
}

func ValidatePhone(str string) bool {
	if !ValidateNum(str) {
		return false
	}

	l := len(str)
	return l >= 6 && l < 11
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func WritePid() error {
	pid_fp, err := os.OpenFile("./server.pid", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open pid file failed[%s]\n", err)
		return err
	}
	defer pid_fp.Close()

	pid := os.Getpid()

	pid_fp.WriteString(strconv.Itoa(pid))
	return nil
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start) + len(start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func SendSmsChina(pszMobis, pszMsg string) (int, string, error) {
	var userId = "JKDX01"
	var password = "2015abc"
	var iMobiCount = "1"

	var url = "http://114.67.48.66:5132/MWGate/wmgw.asmx/MongateCsSpSendSmsNew?"
	var params = map[string]string{
		"userId":     userId,
		"password":   password,
		"pszMobis":   pszMobis,
		"pszMsg":     pszMsg,
		"iMobiCount": iMobiCount,
		"pszSubPort": "*",
	}

	code, result, err := SMSSendRequest("GET", url, params)
	return code, result, err

}

func SendSmsInternational(pszMobis, pszMsg string) (int, string, error) {
	var userId = "JKDX03"
	var password = "codoon20150902"
	var iMobiCount = "1"

	var url = "http://114.67.48.66:5132/MWGate/wmgw.asmx/MongateCsSpSendSmsNew?"
	var params = map[string]string{
		"userId":     userId,
		"password":   password,
		"pszMobis":   pszMobis,
		"pszMsg":     pszMsg,
		"iMobiCount": iMobiCount,
		"pszSubPort": "*",
	}

	code, result, err := SMSSendRequest("GET", url, params)
	return code, result, err
}

func SendMail(email, password, host, toemail, sender, receiver, subject, body string) error {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	from := mail.Address{sender, email}
	to := mail.Address{receiver, toemail}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	err := smtp.SendMail(
		host+":25",
		auth,
		email,
		[]string{to.Address},
		[]byte(message),
	)
	return err

}

func SendGroupMail(email, password, sender, host string, toemail []string, subject, body string) error {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	var tos []string
	from := mail.Address{sender, email}

	for i := range toemail {
		to := mail.Address{"", toemail[i]}
		tos = append(tos, to.Address)
	}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = strings.Join(tos, ",")
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	err := smtp.SendMail(
		host+":25",
		auth,
		email,
		toemail,
		[]byte(message),
	)
	return err

}

type DecimalismConfusion struct {
	Move    int64
	StartId int64
}

func InitDecimalismConfusion(move, startId int64) (*DecimalismConfusion, error) {
	if move <= 0 {
		return nil, errors.New("move error")
	}

	d := DecimalismConfusion{
		Move:    1,
		StartId: startId,
	}

	var i int64 = 0
	for ; i < move; i++ {
		d.Move = d.Move * 10
	}

	return &d, nil
}

func (d *DecimalismConfusion) sign(id int64) int64 {
	var signId int64

	stringId := strconv.FormatInt(id, 10)
	for i := 0; i < len(stringId); i++ {
		k := id / 10
		if k == 0 {
			signId = signId + id
			break
		}
		ii := id - k*10
		signId = signId + ii
		id = k
	}

	return signId % d.Move
}

func (d *DecimalismConfusion) EncodeId(id int64) int64 {
	if id < d.StartId {
		return id
	}

	encodeId := id * d.Move
	encodeId = encodeId + d.sign(id)
	return encodeId
}

func (d *DecimalismConfusion) DecodeId(id int64) (int64, error) {
	if id < d.StartId {
		return id, nil
	}

	var decodeId int64
	decodeId = id / d.Move
	signId := id - decodeId*d.Move
	if signId != d.sign(decodeId) {
		fmt.Println(decodeId, signId, d.sign(decodeId))
		return id, errors.New("decode error")
	}

	return decodeId, nil
}

func RedirectCoreDump(dump_file *os.File) {
	//syscall.Dup2(int(dump_file.Fd()), 2)
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
