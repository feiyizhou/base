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
)

var (
	CharArr = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	NumArr  = []byte("1234567890")
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
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
	str, err := InterfaceToStr(value)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(h, str)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// RandCharStr 生成随机字符串
func RandCharStr(n int) string {
	result := make([]byte, n)
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
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
func InterfaceToStr(value interface{}) (string, error) {
	// interface 转 string
	var str string
	if value == nil {
		return str, nil
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float64, reflect.Float32:
		return fmt.Sprintf("%f", v.Float()), nil
	case reflect.Int, reflect.Uint,
		reflect.Int8, reflect.Uint8,
		reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32,
		reflect.Int64, reflect.Uint64:
		return fmt.Sprintf("%d", v.Int()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool()), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes()), nil
		} else {
			return "", fmt.Errorf("unsupported type: %s", v.Kind())
		}
	default:
		bytes, _ := json.Marshal(value)
		return string(bytes), nil
	}
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
