package main

import "testing"
import "github.com/shewinder/pixiv"

func init() {
	pixiv.InitAuth("5_wLcosaJG103dcOR_ES8ybX3NTwKVxEjH7nFVF9YRA")
}

func TestIllustDetail(t *testing.T) {
	pid := "75101768"
	illust, _ := pixiv.IllustDetail(pid)
	t.Log(illust.Tags)
}

