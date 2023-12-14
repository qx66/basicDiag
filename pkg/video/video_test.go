package video

import "testing"

func TestGetVideo(t *testing.T) {
	videoUrl := "https://startops-static.oss-cn-hangzhou.aliyuncs.com/video/mao1.mov"
	GetVideo(videoUrl)
}
