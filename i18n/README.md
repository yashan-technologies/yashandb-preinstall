# 多语言支持 (Multi-language Support)

本工具支持中文和英文两种语言。

## 使用方法

### 命令行参数

使用 `-l` 或 `--lang` 参数指定语言：

```bash
# 使用中文（默认）
./preinstall -l zh

# 使用英文
./preinstall -l en

```

### 支持的语言代码

- `zh` - 中文 (Chinese)
- `en` - 英文 (English)  

## 语言包文件

语言包使用TOML格式存储在 `i18n/locales/` 目录下：

- `zh.toml` - 中文翻译
- `en.toml` - 英文翻译

## 添加新的翻译

如需添加新的翻译文本：

1. 在对应语言的TOML文件中添加翻译键值对
2. 在代码中使用 `i18n.T("key")` 调用翻译
3. 对于带参数的翻译，使用 `fmt.Sprintf(i18n.T("key"), args...)`

### 示例

在语言包中添加：

```toml
[example]
hello = "你好，%s！"
```

在代码中使用：

```go
import "preinstall/i18n"

msg := fmt.Sprintf(i18n.T("example.hello"), "世界")
```

## 技术实现

- 使用 `github.com/nicksnyder/go-i18n/v2` 库实现国际化
- 语言包使用TOML格式，便于阅读和维护
- 通过Go的embed特性将语言包嵌入到二进制文件中
- 支持运行时动态切换语言

## 注意事项

1. 默认语言为中文 (zh)
2. 如果指定的语言不存在，会回退到中文
3. 所有命令行帮助信息使用英文，以保持国际化标准
4. 程序输出信息根据 `-l` 参数动态切换语言
