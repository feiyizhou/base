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
	"time"
)

var (
	CharArr = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	NumArr  = []byte("1234567890")
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ValueMd5(value interface{}) string {
	h := md5.New()
	_, err := io.WriteString(h, InterfaceToStr(value))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
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
	for i := 0; i < n; i++ {
		result[i] = NumArr[rand.Int31()%int32(len(NumArr))]
	}
	return string(result)
}

// InterfaceToStr interface转string
func InterfaceToStr(value interface{}) string {
	// interface 转 string
	var str string
	if value == nil {
		return str
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		str = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		str = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		str = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		str = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		str = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		str = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		str = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		str = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		str = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		str = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		str = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		str = strconv.FormatUint(it, 10)
	case string:
		str = value.(string)
	case []byte:
		str = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		str = string(newValue)
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
