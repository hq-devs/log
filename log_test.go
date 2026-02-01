package log

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	// 1. 准备测试参数
	testDir := "./test_logs"
	testFile := "test.log"
	logPath := testDir + "/" + testFile

	// 清理旧数据
	defer os.RemoveAll(testDir)

	// 2. 初始化日志器
	// 注意：参数需要对应你 SetLogger 的定义（dir, file, unit, level, count, size, compressT）
	err := SetLogger(testDir, testFile, "day", "debug", 7, 100, 1)
	if err != nil {
		t.Fatalf("Failed to set logger: %v", err)
	}

	// 3. 触发日志记录
	// 这里是关键：记录这个行号，稍后验证输出是否匹配
	Debug("this is a debug message", "key1", "value1")
	Infof("hello %s", "world")

	// 4. 读取文件内容并验证
	file, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundSource := false
	for scanner.Scan() {
		line := scanner.Text()
		t.Logf("Output Line: %s", line)

		// 验证是否包含调用者的文件信息
		// 如果你是在 log_test.go 中运行，日志里应该包含 "log_test.go"
		if strings.Contains(line, "log_test.go") {
			foundSource = true
		}

		// 验证 KV 对是否生效
		if strings.Contains(line, "key1=value1") {
			t.Log("Successfully found custom attributes")
		}
	}

	if !foundSource {
		t.Error("Log source (file:line) not found or incorrect! It might be pointing to the wrapper file instead of the caller.")
	}
}
