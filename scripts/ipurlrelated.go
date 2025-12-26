package scripts

import (
	log "Sowhp/concert/logger"
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func extractDomainAndIP(url string) (domain string) {
	s := strings.TrimPrefix(url, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "..", "__")
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	s = reg.ReplaceAllString(s, "")
	
	return s
}

func FindTextUrl(filepath string) []string {
	if filepath == "" {
		log.Error("文件路径不能为空")
		return []string{}
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Error(fmt.Sprintf("无法打开文件 %s: %v", filepath, err))
		return []string{}
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Warning(fmt.Sprintf("关闭文件失败: %v", closeErr))
		}
	}()

	urls := []string{}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		urlList := extractURL(line)
		for _, url := range urlList {
			if url != "" {
				urls = append(urls, url)
				log.Debug(fmt.Sprintf("第 %d 行提取到地址: %s", lineNum, url))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error(fmt.Sprintf("读取 %s 时发生错误: %v", filepath, err))
		return []string{}
	}

	log.Info(fmt.Sprintf("文件成功提取到 %d 个地址, 开始执行...", len(urls)))
	return urls
}

func extractURL(line string) []string {

	var urlTmpList []string
	var https string = "https://"

	if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
		urllist := append(urlTmpList, line)
		return urllist
	} else if IsIPAddress(line) || IsDomainName(line) {
		line1 := https + line
		urllist := append(urlTmpList, line1)
		return urllist
	} else if IsIPAddressWithPort(line) || IsDomainNameWithPort(line) {
		line1 := https + line
		urllist := append(urlTmpList, line1)
		return urllist
	}

	return []string{}
}

func visitURL(url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
	}
}

func GetTimeStrin() string {
	currentTime := time.Now()
	dateString := currentTime.Format("20060102")

	baseDir := "./result"
	counter := 1

	for {
		resultName := fmt.Sprintf("%s%04d", dateString, counter)
		resultPath := fmt.Sprintf("%s/result_%s", baseDir, resultName)

		if _, err := os.Stat(resultPath); os.IsNotExist(err) {
			return resultName
		}
		counter++
	}
}

func GetUrlStatusCodeAndResponse(url string) (string, string) {
	if url == "" {
		log.Warning("URL为空，无法获取状态码")
		return "N/A", "URL为空"
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	resp, err := client.Get(url)
	if err != nil {
		errStr := err.Error()
		var statusCode, responseContent string

		if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
			statusCode = "TIMEOUT"
			responseContent = fmt.Sprintf("请求超时: %v", err)
		} else if strings.Contains(errStr, "connection refused") {
			statusCode = "CONNECTION_REFUSED"
			responseContent = fmt.Sprintf("连接被拒绝: %v", err)
		} else if strings.Contains(errStr, "no such host") {
			statusCode = "DNS_ERROR"
			responseContent = fmt.Sprintf("DNS解析失败: %v", err)
		} else if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") {
			statusCode = "SSL_ERROR"
			responseContent = fmt.Sprintf("SSL证书错误: %v", err)
		} else {
			statusCode = "ERROR"
			responseContent = fmt.Sprintf("请求失败: %v", err)
		}

		log.WarningWithContext(fmt.Sprintf("%v", err), url)

		time.Sleep(1 * time.Second)
		resp, retryErr := client.Get(url)
		if retryErr != nil {
			log.ErrorWithContext(fmt.Sprintf("%v", retryErr), url)
			return statusCode, responseContent
		}
		defer resp.Body.Close()
	} else {
		defer resp.Body.Close()
	}

	statusCode := strconv.Itoa(resp.StatusCode)
	responseBuilder := strings.Builder{}

	responseBuilder.WriteString(fmt.Sprintf("%s %d %s\n", resp.Proto, resp.StatusCode, resp.Status))

	for name, values := range resp.Header {
		for _, value := range values {
			responseBuilder.WriteString(fmt.Sprintf("%s: %s\n", name, value))
		}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err == nil && len(body) > 0 {
		responseBuilder.WriteString("\n--- Response Body (Preview) ---\n")
		responseBuilder.WriteString(string(body))
		if len(body) == 1024 {
			responseBuilder.WriteString("\n... (truncated)")
		}
	}

	log.Debug(fmt.Sprintf("URL %s 状态码: %s", url, statusCode))
	return statusCode, responseBuilder.String()
}

func GetUrlStatusCode(url string) string {
	statusCode, _ := GetUrlStatusCodeAndResponse(url)
	return statusCode
}

func SmartScreenshot(URL string, resultName string) []string {

	result := ChromeScreenshot(URL, resultName)
	if len(result) > 0 {
		return result
	}

	if strings.HasPrefix(URL, "https://") {
		httpURL := strings.Replace(URL, "https://", "http://", 1)
		log.Info(fmt.Sprintf("访问 %s 失败，正在尝试 HTTP 请求", httpURL))
		return ChromeScreenshot(httpURL, resultName)
	}

	return []string{}
}

func ChromeScreenshot(URL string, resultName string) []string {
	if URL == "" {
		log.Error("URL不能为空")
		return []string{}
	}

	if resultName == "" {
		log.Error("结果名称不能为空")
		return []string{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("ignore-ssl-errors", true),
		chromedp.Flag("ignore-certificate-errors-spki-list", true),
		chromedp.Flag("ignore-certificate-errors-skip-list", true),
		chromedp.Flag("allow-running-insecure-content", true),
		chromedp.Flag("disable-ssl-verification", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-default-apps", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	var pageTitle string
	var screenshot []byte

	executeScreenshot := func() error {
		return chromedp.Run(browserCtx,

			chromedp.Emulate(device.Reset),

			chromedp.EmulateViewport(1920, 1080),

			visitURL(URL),

			chromedp.Sleep(3*time.Second),

			chromedp.ActionFunc(func(ctx context.Context) error {

				chromedp.WaitVisible(`body`, chromedp.ByQuery).Do(ctx)
				return nil
			}),

			chromedp.Evaluate(`document.title || 'No Title'`, &pageTitle),

			chromedp.CaptureScreenshot(&screenshot),
		)
	}

	err := executeScreenshot()
	if err != nil {

		errStr := err.Error()
		if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "ssl") {
			log.WarningWithContext("SSL Certificate Error", URL)
		} else if strings.Contains(errStr, "deadline exceeded") || strings.Contains(errStr, "timeout") {
			log.WarningWithContext("Connect Timeout", URL)
		} else {
			log.WarningWithContext(fmt.Sprintf("%v", err), URL)
		}

		time.Sleep(2 * time.Second)
		err = executeScreenshot()
		if err != nil {
			log.ErrorWithContext(fmt.Sprintf("%v", err), URL)
			return []string{}
		}
		log.Info(fmt.Sprintf("访问 %s 重试成功", URL))
	}

	urlName := extractDomainAndIP(URL)
	if urlName == "" {
		urlName = "unknown"
	}

	photoName := fmt.Sprintf("data/%s-%s.png", urlName, resultName)
	resultPath := fmt.Sprintf("./result/%s/%s", resultName, photoName)

	if err := os.WriteFile(resultPath, screenshot, 0644); err != nil {
		log.Error(fmt.Sprintf("保存截图文件失败 %s: %v", resultPath, err))
		return []string{}
	}

	log.Debug(fmt.Sprintf("截图保存成功: %s", resultPath))

	statusCode, responseContent := GetUrlStatusCodeAndResponse(URL)

	infoArray := []string{URL, pageTitle, statusCode, photoName, responseContent}
	log.Debug(fmt.Sprintf("URL %s 处理完成，标题: %s，状态码: %s", URL, pageTitle, statusCode))

	return infoArray
}
