package logger

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
)

var (
	Red         = color.Red.Render
	Cyan        = color.Cyan.Render
	Yellow      = color.Yellow.Render
	White       = color.White.Render
	Blue        = color.Blue.Render
	Purple      = color.Style{color.Magenta, color.OpBold}.Render
	LightRed    = color.Style{color.Red, color.OpBold}.Render
	LightGreen  = color.Style{color.Green, color.OpBold}.Render
	LightWhite  = color.Style{color.White, color.OpBold}.Render
	LightCyan   = color.Style{color.Cyan, color.OpBold}.Render
	LightYellow = color.Style{color.Yellow, color.OpBold}.Render
)

var (
	defaultLevel = LevelWarning
	noWrite      int
)

func SetLevel(l Level) {
	defaultLevel = l
}

func getCallerInfo(skip int) (info string) {
	_, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		info = "runtime.Caller() failed"
	}

	fileName := path.Base(file)
	return fmt.Sprintf("%s line:%d", fileName, lineNo)
}

func log(l Level, w int, detail string) {
	switch LogLevel {
	case 0:
		SetLevel(0)
	case 1:
		SetLevel(1)
	case 2:
		SetLevel(2)
	case 3:
		SetLevel(3)
	case 4:
		SetLevel(4)
	case 5:
		SetLevel(5)
	}

	if l > defaultLevel {
		return
	}

	ClearProgressBar()

	if NoColor {
		fmt.Println(clean(detail))
		return
	} else {
		fmt.Println(detail)
	}

	if noWrite == 0 {
		writeLogFile(clean(detail), OutputFileName)
	}

	if l == LevelFatal {
		os.Exit(0)
	}
}

func Fatal(detail string) {
	noWrite = 1
	log(LevelFatal, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightRed("×"), LightWhite("]"), detail))
}

func Error(detail string) {
	noWrite = 1
	log(LevelError, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightRed("×"), LightWhite("]"), detail))
}

func ErrorWithContext(errorMsg, url string) {
	noWrite = 1
	log(LevelError, noWrite, fmt.Sprintf("%s%s%s %s%s%s 访问 %s 时遇到错误，重试中", LightWhite("["), LightRed("×"), LightWhite("]"), LightWhite("["), LightWhite(errorMsg), LightWhite("]"), url))
}

func Info(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightGreen("√"), LightWhite("]"), detail))
}

func Warning(detail string) {
	noWrite = 1
	log(LevelWarning, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightYellow("!"), LightWhite("]"), detail))
}

func WarningWithContext(errorMsg, url string) {
	noWrite = 1
	log(LevelWarning, noWrite, fmt.Sprintf("%s%s%s %s%s%s 访问 %s 时遇到错误，重试中", LightWhite("["), LightYellow("!"), LightWhite("]"), LightWhite("["), LightWhite(errorMsg), LightWhite("]"), url))
}

func Debug(detail string) {
	noWrite = 1
	log(LevelDebug, noWrite, fmt.Sprintf("%s%s%s %s%s%s %s", LightWhite("["), LightWhite("?"), LightWhite("]"), LightWhite("["), Yellow(getCallerInfo(2)), LightWhite("]"), detail))
}

func Verbose(detail string) {
	noWrite = 1
	log(LevelVerbose, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightCyan("i"), LightWhite("]"), detail))
}

func Success(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("%s", detail))
}

func Failed(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("%s%s%s %s", LightWhite("["), LightRed("×"), LightWhite("]"), detail))
}

func Common(detail string) {
	noWrite = 1
	log(LevelInfo, noWrite, fmt.Sprintf("%s", detail))
}

func getTime() string {
	return time.Now().Format("15:04:05")
}

func getDate() string {
	return time.Now().Format("2006.1.2")
}

func DebugError(err error) bool {
	if err != nil {
		pc, _, line, _ := runtime.Caller(1)
		Debug(fmt.Sprintf("%s%s%s",
			White(runtime.FuncForPC(pc).Name()),
			LightWhite(fmt.Sprintf(" line:%d ", line)),
			White(err)))
		return true
	}
	return false
}

func clean(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

func writeLogFile(result string, filename string) {
	var text = []byte(result + "\n")
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write(text)
	fl.Close()
	if err != nil {
		fmt.Printf("Write %s error, %v\n", filename, err)
	}
}

func ShowProgressBar(current, total int, prefix string) {
	percent := float64(current) / float64(total) * 100
	barLength := 50
	filledLength := int(float64(barLength) * float64(current) / float64(total))

	bar := strings.Repeat("█", filledLength) + strings.Repeat("░", barLength-filledLength)

	if NoColor {
		fmt.Printf("\r%s: [%s] %.1f%% (%d/%d)", prefix, bar, percent, current, total)
	} else {
		fmt.Printf("\r%s: [%s%s] %.1f%% (%d/%d)",
			LightCyan(prefix),
			LightWhite(strings.Repeat("█", filledLength)),
			strings.Repeat("░", barLength-filledLength),
			percent, current, total)
	}

	if current == total {
		fmt.Println()
	}
}

func ClearProgressBar() {
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
}

func UpdateProgress(current, total int, message string) {
	ShowProgressBar(current, total, message)
}

func ProgressWithColor(current, total int, prefix string, showETA bool) {
	percent := float64(current) / float64(total) * 100
	barLength := 40
	filledLength := int(float64(barLength) * float64(current) / float64(total))

	filledBar := LightWhite(strings.Repeat("█", filledLength))
	emptyBar := strings.Repeat("░", barLength-filledLength)

	etaInfo := ""
	if showETA && current > 0 {
		etaInfo = " | ETA: 计算中"
	}

	if NoColor {
		fmt.Printf("\r%s: [%s] %.1f%% (%d/%d)%s",
			prefix, strings.Repeat("█", filledLength)+strings.Repeat("░", barLength-filledLength),
			percent, current, total, etaInfo)
	} else {
		fmt.Printf("\r%s: [%s%s] %s (%d/%d)%s",
			LightCyan(prefix),
			filledBar,
			emptyBar,
			LightWhite(fmt.Sprintf("%.1f%%", percent)),
			current, total,
			LightYellow(etaInfo))
	}

	if current == total {
		fmt.Println()
	}
}
