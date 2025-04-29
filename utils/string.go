package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	CharArr = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	NumArr  = []byte("1234567890")
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandUUIDStr() string {
	return uuid.NewString() // 输出类似 63719109-0989-49e9-b339-a7e6d04c3507
}

func DeleteSpecialChar(str string) string {
	if len(str) == 0 {
		return ""
	}
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\\", "")
	str = strings.ReplaceAll(str, "-", "")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "/", "")
	return str
}

func ValueMd5(value any) (string, error) {
	h := md5.New()
	_, err = io.WriteString(h, InterfaceToStr(value))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// RandCharStr 生成随机字符串
func RandCharStr(n int) string {
	result := make([]byte, n)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range n {
		result[i] = CharArr[rand.Int31()%int32(len(CharArr))]
	}
	return string(result)
}

// RandNumStr 生成随机数字串
func RandNumStr(n int) string {
	result := make([]byte, n)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range n {
		result[i] = NumArr[rand.Int31()%int32(len(NumArr))]
	}
	return string(result)
}

// InterfaceToStr interface转string
func InterfaceToStr(value interface{}) string {
	// interface 转 string
	var str string
	if value == nil {
		return ""
	}
	switch value := value.(type) {
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		str = strconv.FormatFloat(float64(value), 'f', -1, 32)
	case int:
		str = strconv.Itoa(value)
	case uint:
		str = strconv.FormatUint(uint64(value), 10)
	case int8:
		str = strconv.Itoa(int(value))
	case uint8:
		str = strconv.FormatUint(uint64(value), 10)
	case int16:
		str = strconv.Itoa(int(value))
	case uint16:
		str = strconv.FormatUint(uint64(value), 10)
	case int32:
		str = strconv.Itoa(int(value))
	case uint32:
		str = strconv.FormatUint(uint64(value), 10)
	case int64:
		str = strconv.FormatInt(value, 10)
	case uint64:
		str = strconv.FormatUint(value, 10)
	case string:
		str = value
	case []byte:
		str = string(value)
	default:
		bytes, _ := json.Marshal(value)
		return string(bytes)
	}
	return str
}

// SplitStrToSubStrArr 字符串根据长度切割为字符串数组
func SplitStrToSubStrArr(s string, l int) []string {
	var (
		subArr     [][]rune
		rl, rr, rs = 0, 0, []rune(s)
		strArr     []string
		ss         string
	)
	if int(rs[0]) == 34 {
		rs = rs[1:]
	}
	if int(rs[len(rs)-1]) == 34 {
		rs = rs[:len(rs)-1]
	}
	if len(rs) <= l {
		subArr = append(subArr, rs)
	} else {
		for {
			rl = rr
			rr += l
			if rr >= len(rs) {
				subArr = append(subArr, rs[rl:])
				break
			}
			subArr = append(subArr, rs[rl:rr])
		}
	}
	for _, subRune := range subArr {
		for _, r := range subRune {
			ss += string(r)
		}
		strArr = append(strArr, ss)
		ss = ""
	}
	return strArr
}

func ValueIsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// StrToUint64 string to uint64
func StrToUint64(str string) uint64 {
	v, _ := strconv.ParseUint(str, 10, 64)
	return v
}

// StrToUint64E string to uint64 with error
func StrToUint64E(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

func Base64Encode(content string) string {
	contentBytes := []byte(content)
	return base64.StdEncoding.EncodeToString(contentBytes)
}

func Base64Decode(content string) (string, error) {
	decodeContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}
	return string(decodeContent), err
}

func ObjToJsonStr(obj any) string {
	objBytes, _ := json.Marshal(obj)
	return string(objBytes)
}
