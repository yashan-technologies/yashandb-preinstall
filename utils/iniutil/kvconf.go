package iniutil

import (
	"fmt"
	"os"
	"strings"

	"preinstall/commons/consts"

	"git.yasdb.com/go/yasutil/fs"
)

const (
	new_line = "\n"
	equal    = "="
	space    = " "
)

var (
	_commentChar = "#" // 注释符
)

type DataLine struct {
	StartLine int
	EndLine   int
	Key       string
	Value     string
}

type KeyData struct {
	RawLines     []string
	Filepath     string
	Data         *DataLine
	DefaultQuote string
}

func LoadKeyData(filePath string, key string) (*KeyData, error) {
	if !fs.IsFileExist(filePath) {
		return nil, fmt.Errorf("file not exist: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), new_line)

	kd := &KeyData{
		Filepath:     filePath,
		DefaultQuote: "\"",
		RawLines:     lines,
	}
	for i, line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, _commentChar) {
			continue
		}

		if strings.HasPrefix(line, key+equal) {
			kd.Data = &DataLine{
				Key:       key,
				StartLine: i,
			}
			value := strings.TrimPrefix(line, strings.Split(line, equal)[0]+equal)

			if len(value) == 0 {
				kd.Data.EndLine = kd.Data.StartLine
				continue
			}

			quote := kd.getQuote(value)

			// 没有引号包裹，则认为是单行
			if len(quote) == 0 {
				kd.Data.EndLine = kd.Data.StartLine
				kd.Data.Value = value
				continue
			}

			// 如果有引号包裹，则根据引号类型查找结束行
			if !strings.HasSuffix(value, quote) || strings.HasSuffix(value, `\`+quote) {
				value += new_line
				for j := i + 1; j < len(lines); j++ {
					// 如果当前行以引号结尾，且不是转义引号，则认为找到了结束行
					if strings.HasSuffix(lines[j], quote) && !strings.HasSuffix(lines[j], `\`+quote) {
						value += lines[j]
						i = j
						kd.Data.EndLine = j
						break
					}
					value += lines[j] + new_line
				}
			} else {
				kd.Data.EndLine = kd.Data.StartLine
			}
			kd.Data.Value = value
			return kd, nil
		}
	}
	return kd, nil
}

func (kd *KeyData) Valid() bool {
	return kd != nil && kd.Data != nil &&
		kd.Data.StartLine >= 0 && kd.Data.EndLine >= kd.Data.StartLine
}

func (kd *KeyData) UpdateTo(filePath, value string) error {
	if !kd.Valid() {
		return fmt.Errorf("invalid kvdata")
	}

	kd.RawLines[kd.Data.StartLine] = kd.Data.Key + equal + kd.Data.Value
	if kd.Data.StartLine != kd.Data.EndLine {
		kd.RawLines = append(kd.RawLines[:kd.Data.StartLine+1], kd.RawLines[kd.Data.EndLine+1:]...)
	}

	content := strings.Join(kd.RawLines, new_line)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (kd *KeyData) Update(value string) error {
	return kd.UpdateTo(kd.Filepath, value)
}

func (kd *KeyData) AppendTo(filePath, key string, values ...string) error {
	if len(values) == 0 {
		return nil
	}

	// 如果是更新当前Key，则直接更新
	if kd.Valid() && kd.Data.Key == key {
		kd.appendValues(values...)
		return kd.UpdateTo(filePath, kd.Data.Value)
	}

	// 如果是更新不存在的Key，则直接追加
	kd.RawLines = append(kd.RawLines, key+equal+fmt.Sprintf("%s%s%s", kd.DefaultQuote, strings.Join(values, space), kd.DefaultQuote))
	content := strings.Join(kd.RawLines, new_line)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (kd *KeyData) Append(key string, values ...string) error {
	return kd.AppendTo(kd.Filepath, key, values...)
}

func (kd *KeyData) BackupTo(filePath string) error {
	backupFile := filePath + consts.BakupExt()
	if err := fs.CopyFile(filePath, backupFile); err != nil {
		return err
	}
	return nil
}

func (kd *KeyData) Backup() error {
	return kd.BackupTo(kd.Filepath)
}

func (kd *KeyData) getQuote(value string) string {
	if strings.HasPrefix(value, `"""`) {
		return `"""`
	} else if strings.HasPrefix(value, `"`) {
		return `"`
	} else if strings.HasPrefix(value, `'''`) {
		return `'''`
	} else if strings.HasPrefix(value, `'`) {
		return `'`
	} else if strings.HasPrefix(value, "`") {
		return "`"
	}
	return ""
}

func (kd *KeyData) appendValues(values ...string) {
	quote := kd.getQuote(kd.Data.Value)
	isOriginValueEmpty := len(strings.TrimPrefix(strings.TrimSuffix(kd.Data.Value, quote), quote)) == 0

	if isOriginValueEmpty {
		valueString := strings.Join(values, space)
		if len(quote) == 0 {
			kd.Data.Value = fmt.Sprintf("%s%s%s", kd.DefaultQuote, valueString, kd.DefaultQuote)
			return
		}
		kd.Data.Value = quote + valueString + quote
		return
	}

	var targetFileds []string
	matchedValueIndexes := make(map[int]struct{})
	originValueUnquote := strings.TrimPrefix(strings.TrimSuffix(kd.Data.Value, quote), quote)
	fields := strings.Fields(originValueUnquote)
	for _, f := range fields {
		field := f
		for i, v := range values {
			prefix := v
			// 如果是KV类型的Value，则只取Key
			if strings.Contains(v, equal) {
				prefix = strings.Split(v, equal)[0] + equal
			}
			if strings.HasPrefix(f, prefix) {
				field = v
				matchedValueIndexes[i] = struct{}{}
				break
			}
		}
		targetFileds = append(targetFileds, field)
	}

	// 如果没有匹配到任何值，则直接追加
	for i, v := range values {
		if _, ok := matchedValueIndexes[i]; !ok {
			targetFileds = append(targetFileds, v)
		}
	}

	valueString := strings.Join(targetFileds, space)
	if len(quote) == 0 {
		kd.Data.Value = fmt.Sprintf("%s%s%s", kd.DefaultQuote, valueString, kd.DefaultQuote)
		return
	}
	kd.Data.Value = quote + valueString + quote
}
