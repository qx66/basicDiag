package video

import (
	"fmt"
	"github.com/imkira/go-libav/avformat"
)

func GetVideo() {
	ctx, err := avformat.NewContextForInput()
	if err != nil {
		return
	}
	
	ctx.OpenInput("https://startops-static.oss-cn-hangzhou.aliyuncs.com/video/mao1.mov", nil, nil)
	ctx.FindStreamInfo(nil)
	
	for _, stream := range ctx.Streams() {
		
		fmt.Println("Index: ", stream.Index())
		fmt.Println("StartTime: ", stream.StartTime())
		fmt.Println("Duration: ", stream.Duration())
		fmt.Println("AverageFrameRate: ", stream.AverageFrameRate())
	}
	
}
