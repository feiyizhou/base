package sensitive_words

import (
	"bufio"
	"io"
	"os"
)

type SensitiveWordValidator struct {
	trie *TrieTree
}

func NewSensWordValidator() *SensitiveWordValidator {
	return &SensitiveWordValidator{trie: NewTrieTree()}
}

// Load 加载敏感词
func (sw *SensitiveWordValidator) Load(rd io.Reader) error {
	buf := bufio.NewReader(rd)
	for {
		keywords, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		sw.trie.Add(string(keywords))
	}
	return nil
}

// LoadFromFile 加载敏感词库文件
func (sw *SensitiveWordValidator) LoadFromFile(path string) error {
	fs, err := os.Open(path)
	defer fs.Close()
	if err != nil {
		return nil
	}
	return sw.Load(fs)
}

// AddWord  添加敏感词
func (sw *SensitiveWordValidator) AddWord(word string) {
	sw.trie.Add(word)
}

// AddWords 添加敏感词数组
func (sw *SensitiveWordValidator) AddWords(words ...string) {
	sw.trie.AddWords(words...)
}

// Replace 过滤敏感词为*
func (sw *SensitiveWordValidator) Replace(input string) string {
	return sw.trie.Replace(input)
}

// Find 查找敏感词,找到第一个就退出
func (sw *SensitiveWordValidator) Find(input string) (sensitive bool, keyword string) {
	return sw.trie.Find(input)
}

// Check 是否包含敏感词
func (sw *SensitiveWordValidator) Check(input string) (sensitive bool) {
	sensitive, _ = sw.trie.Find(input)
	return sensitive
}

// FindAll 查找ALL敏感词
func (sw *SensitiveWordValidator) FindAll(input string) (sensitive bool, results []string) {
	return sw.trie.FindAll(input)
}

// FindAny 找到N个敏感词才退出
func (sw *SensitiveWordValidator) FindAny(input string, count int) (sensitive bool, results []string) {
	return sw.trie.FindAny(input, count)
}
