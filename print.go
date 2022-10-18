package gutils

import "fmt"

type color int

const (
	Blank color = 30 + iota
	Red
	Green
	Yellow
	Blue
	Purple
	DarkGreen //深绿
	White
)

//彩色打印输出
func ColourPrint(msg interface{}, c color) {
	fmt.Printf("\x1b[%dm%v\x1b[0m\n", c, msg)
}

func PrintNum(num int) {

}
