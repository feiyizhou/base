package utils

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	FileMode0755 = 0755
	FileMode0644 = 0644
)

type FileInfo struct {
	os.FileInfo
	Path string
}

// IsExist 文件是否存在，true-存在，false-不存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func CreateDir(path string) error {
	if !IsExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func AllFilesInDir(path string) ([]FileInfo, error) {
	var fileArr []FileInfo
	if !IsDir(path) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		fileArr = append(fileArr, FileInfo{
			FileInfo: fileInfo,
			Path:     path,
		})
		return fileArr, nil
	}
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileArr = append(fileArr, FileInfo{
			FileInfo: info,
			Path:     path,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileArr, err
}

func CountDirFiles(dirName string) int {
	if !IsExist(dirName) {
		return 0
	}
	if !IsDir(dirName) {
		return 0
	}
	var count int
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	})
	if err != nil {
		return 0
	}
	return count
}

// FileMD5ByPath count file md5
func FileMD5ByPath(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	return FileMD5(file)
}

func FileMD5(file *os.File) (string, error) {
	m := md5.New()
	if _, err := io.Copy(m, file); err != nil {
		return "", err
	}
	fileMd5 := fmt.Sprintf("%x", m.Sum(nil))
	return fileMd5, nil
}

func CreateFileByAllPath(fileName string) error {
	err := MkFileFullPathDir(fileName)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func MkFileFullPathDir(fileName string) error {
	localDir := filepath.Dir(fileName)
	err := os.MkdirAll(localDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create local dir %s failed: %v", localDir, err)
	}
	return nil
}

func Tar(src, dst, trimPrefix string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()
		path = strings.TrimPrefix(path, trimPrefix)
		fmt.Println(strings.TrimPrefix(path, string(filepath.Separator)))
		hdr.Name = strings.TrimPrefix(path, string(filepath.Separator))
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := io.Copy(tw, fr); err != nil {
			return err
		}
		return nil
	})
}

func Untar(src, dst string) error {
	fr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}
		dstPath := filepath.Join(dst, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if !IsExist(dstPath) && IsDir(dstPath) {
				if err := CreateDir(dstPath); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			if dir := filepath.Dir(dstPath); !IsExist(dir) {
				if err := CreateDir(dir); err != nil {
					return err
				}
			}
			file, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err = io.Copy(file, tr); err != nil {
				return err
			}
			fmt.Println(dstPath)
			file.Close()
		}
	}
}

// ReadFileToBytes 按字节读取文件
func ReadFileToBytes(path string) ([]byte, error) {
	if !IsExist(path) {
		return nil, errors.New(fmt.Sprintf("file %s is not exist", path))
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, err
}

// ReadFileToLines 按行读取文件为字符串数组
func ReadFileToLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var strArr []string
	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}
		strArr = append(strArr, string(lineBytes))
	}
	return strArr, err
}

// WriteToFile 追加内容到现存文件中
func WriteToFile(content any, path string) error {
	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, FileMode0755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(path, bytes, FileMode0644); err != nil {
		return err
	}
	return nil
}

// AppendStrToFile 在已经存在的文件中，添加字符串
func AppendStrToFile(path, content string) (err error) {
	if !IsExist(path) {
		err = fmt.Errorf("file is not exist, file path: %s ", path)
	} else {
		var file *os.File
		file, err = os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			return
		}
		defer file.Close()
		_, err = io.WriteString(file, content)
	}
	return
}
