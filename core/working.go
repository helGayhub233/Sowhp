package core

import (
	dire "Sowhp/concert"
	log "Sowhp/concert/logger"
	"Sowhp/scripts"
	"errors"
	"flag"
	"fmt"
	"sync"
)

const Banner = `
╔═══════════════════════════════════════════════════╗
║                                                   ║
║    ███████╗ ██████╗ ██╗    ██╗██╗  ██╗██████╗     ║
║    ██╔════╝██╔═══██╗██║    ██║██║  ██║██╔══██╗    ║
║    ███████╗██║   ██║██║ █╗ ██║███████║██████╔╝    ║
║    ╚════██║██║   ██║██║███╗██║██╔══██║██╔═══╝     ║
║    ███████║╚██████╔╝╚███╔███╔╝██║  ██║██║         ║
║    ╚══════╝ ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═╝╚═╝         ║
║                                                   ║
║              Website Screenshot Tool              ║
║                 Revision & 修改版                 ║
╚═══════════════════════════════════════════════════╝

`

type Config struct {
	FilePath string
	LogLevel int
}

type App struct {
	config      *Config
	resultMap   map[string]map[string][]string
	arrayMap    map[string][]string
	count       int
	countResult int
	mu          sync.Mutex
}

func NewApp() *App {
	return &App{
		config:      &Config{},
		resultMap:   make(map[string]map[string][]string),
		arrayMap:    make(map[string][]string),
		count:       0,
		countResult: 0,
	}
}

func (app *App) parseFlags() error {
	flag.StringVar(&app.config.FilePath, "f", "", "指定包含URL列表的文本文件路径（必需参数）\n\t\t示例: -f /path/to/urls.txt（每行一个地址）")
	flag.IntVar(&app.config.LogLevel, "log", 3, "设置日志输出详细程度（可选参数，默认值: 3）\n\t\t级别说明: 1=错误 2=警告 3=信息 4=调试\n\t\t示例: -log 4")
	print(Banner)
	flag.Parse()

	if app.config.FilePath == "" {
		return errors.New("文件路径不能为空")
	}

	log.LogLevel = app.config.LogLevel
	return nil
}

func Run() error {
	app := NewApp()
	if err := app.parseFlags(); err != nil {
		flag.Usage()
		log.Error(err.Error())
		return err
	}

	log.Debug(fmt.Sprintf("当前输入路径为：%s", app.config.FilePath))
	return app.run()
}

func isValidPath(path string) bool {
	return len(path) > 0
}

func (app *App) run() error {
	urls := scripts.FindTextUrl(app.config.FilePath)
	if len(urls) == 0 {
		return errors.New("未能从文件中获取到有效的 URL 列表")
	}

	log.Debug(fmt.Sprint("已获取 URL 列表：", urls))
	resultName := fmt.Sprintf("result_%s", scripts.GetTimeStrin())
	total := len(urls)

	if err := app.createDirectories(resultName); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := app.processURLs(urls, resultName, total); err != nil {
		return fmt.Errorf("处理截图失败: %w", err)
	}

	app.resultMap[resultName] = app.arrayMap
	log.Info(fmt.Sprintf("处理完成，成功截图 %d 个网站", app.countResult))
	if err := scripts.CreateHtml(app.resultMap); err != nil {
		return fmt.Errorf("生成报告失败: %w", err)
	}

	return nil
}

func (app *App) createDirectories(resultName string) error {
	dire.MkdirResport()
	dire.Dir_mk(fmt.Sprintf("./result/%s", resultName))
	dire.Dir_mk(fmt.Sprintf("./result/%s/%s", resultName, "data"))
	return nil
}

func (app *App) processURLs(urls []string, resultName string, total int) error {
	maxConcurrency := 5
	if total < maxConcurrency {
		maxConcurrency = total
	}

	urlChan := make(chan string, len(urls))
	resultChan := make(chan ScreenshotResult, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go app.screenshotWorker(&wg, urlChan, resultChan, resultName)
	}

	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		app.mu.Lock()
		app.count++

		log.ClearProgressBar()

		if result.Success {
			app.countResult++
			log.Common(fmt.Sprintf("%s %s", log.LightGreen("[√]"), result.URL))
			app.arrayMap[result.URL] = []string{result.Data[1], result.Data[2], result.Data[3], result.Data[4]}
		} else {
			log.Common(fmt.Sprintf("%s %s - %s", log.LightRed("[×]"), result.URL, result.Error))
			app.arrayMap[result.URL] = []string{"无标题", "连接失败", "data/", result.Error}
		}

		log.ShowProgressBar(app.count, total, "执行进度")
		app.mu.Unlock()
	}

	return nil
}

type ScreenshotResult struct {
	URL     string
	Success bool
	Data    []string
	Error   string
}

func (app *App) screenshotWorker(wg *sync.WaitGroup, urlChan <-chan string, resultChan chan<- ScreenshotResult, resultName string) {
	defer wg.Done()

	for url := range urlChan {
		urlResultList := scripts.SmartScreenshot(url, resultName)

		result := ScreenshotResult{
			URL: url,
		}

		if len(urlResultList) == 0 {
			result.Success = false
			result.Error = "截图失败或网络超时"
		} else {
			result.Success = true
			result.Data = urlResultList
		}

		resultChan <- result
	}
}
