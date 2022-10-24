package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
)

func main() {
	//var loc baidu.Locations
	//进行截图
	//filename := Begin()
	//if filename == "" {
	//	return
	//}
	//ocr识别
	//data := baidu.AccurateBasic("zawu/test2.png")
	//for _, v := range data.WordsResult {
	//	if find := strings.Contains(v.Words, "迭嘉"); find {
	//		loc = v.Location
	//	}
	//}
	//fmt.Println(loc)
	robotgo.Move(1091, 1041)
	//fmt.Println(robotgo.GetMousePos())
}

func Begin() (filename string) {
	bitmap := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(bitmap)
	img := robotgo.ToImage(bitmap)
	err := robotgo.Save(img, fmt.Sprintf("zawu/lalaimg.png"))
	if err != nil {
		fmt.Println("screen is ok")
		return ""
	}
	return "lalaimg.png"
}
