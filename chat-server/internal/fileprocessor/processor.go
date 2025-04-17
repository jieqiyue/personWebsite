package fileprocessor

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// 草稿文件目录
	draftDir = "public/markdown/draftArticles"
	// 图片目标目录
	imageDir = "public/images"
	// 图片目标子目录
	imageArtDir = "public/images/art"
	// 图片链接正则表达式 - 匹配Markdown中的图片引用格式 ![alt](path/to/image)
	imgRegexPattern = `!\[(.*?)\]\((.*?)\)`
)

// ProcessDrafts 处理草稿目录中的所有Markdown文件
func ProcessDrafts() error {
	log.Println("开始处理草稿文件夹中的Markdown文件...")

	// 检查草稿目录是否存在
	if _, err := os.Stat(draftDir); os.IsNotExist(err) {
		log.Printf("草稿目录 %s 不存在，无需处理", draftDir)
		return nil
	}

	// 扫描草稿文件夹中的Markdown文件
	mdFiles, err := ScanMarkdownFiles(draftDir)
	if err != nil {
		return fmt.Errorf("扫描Markdown文件失败: %w", err)
	}

	if len(mdFiles) == 0 {
		log.Println("草稿文件夹中没有Markdown文件，无需处理")
		return nil
	}

	log.Printf("找到 %d 个Markdown草稿文件", len(mdFiles))

	// 确保图片目录存在
	if err := ensureDirectoryExists(imageDir); err != nil {
		return fmt.Errorf("确保图片目录存在失败: %w", err)
	}

	// 确保图片art子目录存在
	if err := ensureDirectoryExists(imageArtDir); err != nil {
		return fmt.Errorf("确保图片art子目录存在失败: %w", err)
	}

	// 处理每个Markdown文件
	for _, mdFile := range mdFiles {
		if err := ProcessMarkdownFile(mdFile); err != nil {
			log.Printf("处理文件 %s 失败: %v", mdFile, err)
			// 继续处理其他文件
			continue
		}
		log.Printf("成功处理文件 %s", mdFile)
	}

	log.Println("草稿文件处理完成")
	return nil
}

// ScanMarkdownFiles 扫描目录获取所有Markdown文件
func ScanMarkdownFiles(dir string) ([]string, error) {
	var mdFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理.md文件
		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			mdFiles = append(mdFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mdFiles, nil
}

// ProcessMarkdownFile 处理单个Markdown文件
func ProcessMarkdownFile(filePath string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	fileContent := string(content)

	// 提取图片引用
	imgRefs, err := ExtractImageReferences(fileContent)
	if err != nil {
		return fmt.Errorf("提取图片引用失败: %w", err)
	}

	if len(imgRefs) == 0 {
		log.Printf("文件 %s 中没有图片引用", filePath)
		return nil
	}

	// 创建图片的映射，用于替换
	imgMapping := make(map[string]string)

	// 处理每个图片引用
	for _, imgRef := range imgRefs {
		// 处理图片路径，移动图片
		oldPath := imgRef.Path

		// 如果是相对路径，转换为相对于草稿文件的绝对路径
		if !filepath.IsAbs(oldPath) && !strings.HasPrefix(oldPath, "/") {
			oldPath = filepath.Join(filepath.Dir(filePath), oldPath)
		}

		// 生成新的图片路径和文件名
		newFileName := generateImageFileName(filepath.Base(oldPath))

		// 使用指定格式的引用路径 (/images/art/filename.ext)
		newRelativePath := fmt.Sprintf("/images/art/%s", newFileName)

		// 图片的实际保存路径
		newAbsolutePath := filepath.Join(imageArtDir, newFileName)

		// 移动图片文件
		if err := MoveImageFile(oldPath, newAbsolutePath); err != nil {
			log.Printf("移动图片 %s 失败: %v, 跳过此图片", oldPath, err)
			continue
		}

		// 添加到映射中
		imgMapping[imgRef.OriginalText] = fmt.Sprintf("![%s](%s)", imgRef.Alt, newRelativePath)
	}

	// 更新Markdown内容
	updatedContent := UpdateMarkdownContent(fileContent, imgMapping)

	// 写回文件
	if err := os.WriteFile(filePath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("写入更新后的文件失败: %w", err)
	}

	return nil
}

// ImageReference 表示Markdown中的图片引用
type ImageReference struct {
	OriginalText string // 原始文本
	Alt          string // 图片描述
	Path         string // 图片路径
}

// ExtractImageReferences 从Markdown内容中提取图片引用
func ExtractImageReferences(content string) ([]ImageReference, error) {
	var imgRefs []ImageReference

	re, err := regexp.Compile(imgRegexPattern)
	if err != nil {
		return nil, fmt.Errorf("编译正则表达式失败: %w", err)
	}

	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			imgRefs = append(imgRefs, ImageReference{
				OriginalText: match[0],
				Alt:          match[1],
				Path:         match[2],
			})
		}
	}

	return imgRefs, nil
}

// MoveImageFile 移动图片文件到新位置
func MoveImageFile(sourcePath, destPath string) error {
	// 检查源文件是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("源图片文件不存在: %w", err)
	}

	// 确保目标目录存在
	if err := ensureDirectoryExists(filepath.Dir(destPath)); err != nil {
		return err
	}

	// 打开源文件
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer sourceFile.Close()

	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 同步文件内容到磁盘
	if err = destFile.Sync(); err != nil {
		return fmt.Errorf("同步文件内容失败: %w", err)
	}

	log.Printf("成功将图片从 %s 移动到 %s", sourcePath, destPath)
	return nil
}

// UpdateMarkdownContent 更新Markdown内容中的图片引用
func UpdateMarkdownContent(content string, imgMapping map[string]string) string {
	updatedContent := content

	// 替换所有图片引用
	for oldText, newText := range imgMapping {
		updatedContent = strings.Replace(updatedContent, oldText, newText, -1)
	}

	return updatedContent
}

// generateImageFileName 生成新的图片文件名
func generateImageFileName(originalName string) string {
	// 获取文件扩展名
	ext := filepath.Ext(originalName)

	// 文件名部分（不含扩展名）
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 添加时间戳防止重名
	timestamp := time.Now().Format("20060102150405")

	// 组合新文件名: 原名-时间戳.扩展名
	return fmt.Sprintf("%s-%s%s", nameWithoutExt, timestamp, ext)
}

// ensureDirectoryExists 确保目录存在，如果不存在则创建
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		log.Printf("创建目录: %s", dir)
	}
	return nil
}
