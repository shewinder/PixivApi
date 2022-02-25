package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shewinder/pixiv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getIllustRanking(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().AddDate(0, 0, -1).Format("2006-1-2")) // 日期默认为昨天
	mode := c.DefaultQuery("mode", "day")                                           // mode默认day
	offset := c.DefaultQuery("offset", "30")
	m, err := pixiv.IllustRanking(mode, date, offset)
	if err != nil {
		log.Errorf("%v", err)
	}
	c.JSON(http.StatusOK, m)
}

func getIllustDetail(c *gin.Context) {
	pid := c.Query("illust_id")
	m, err := pixiv.IllustDetail(pid)
	if err != nil {
		log.Errorf("%v", err)
	}
	c.JSON(http.StatusOK, m)
}

func getIllustFollow(c *gin.Context) {
	restrict := c.DefaultQuery("restrict", "public")
	m, err := pixiv.IllustFollow(restrict)
	if err != nil {
		log.Errorf("%v", err)
	}
	c.JSON(http.StatusOK, m)
}

func getUserIllusts(c *gin.Context) {
	offset := c.DefaultQuery("offset", "0")
	type_ := c.DefaultQuery("type", "illust")
	userId := c.Query("user_id")
	m, err := pixiv.UserIllusts(userId, offset, type_)
	if err != nil {
		log.Errorf("%v", err)
	}
	c.JSON(http.StatusOK, m)
}

func getPidFile(c *gin.Context) {
	_, err := os.Stat("image")
	if !os.IsExist(err) {
		os.Mkdir("image", os.ModePerm)
	}
	file := c.Param("file")
	_, err = os.Lstat(fmt.Sprintf("image/%s", file))
	if err != nil {
		log.Infof("%v dose not exist, downloading first", file)
		pid, index := strings.Split(file, "-")[0], strings.Split(file, "-")[1]
		index = strings.Split(index, ".")[0]
		_, err = strconv.Atoi(pid)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "illust_id is not a number"})
			c.Abort()
		}
		i, err := strconv.Atoi(index)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "index is not a number"})
			c.Abort()
		}
		illust, err := pixiv.IllustDetail(pid)
		if err != nil {
			log.Errorf("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			c.Abort()
		}
		if i < 1 || i > illust.PageCount {
			c.JSON(http.StatusNotFound, gin.H{"error": "index out of range"})
			c.Abort()
		}
		var url string
		if i == 1 {
			url = illust.MetaSinglePage.OriginalImageURL
		} else {
			url = illust.MetaPages[i-1].ImageUrls.Original
		}
		url = strings.Replace(url, "i.pximg.net", "pixiv.shewinder.win", 1)
		log.Infof("downloading %v", url)
		path := fmt.Sprintf("image/%s", file)
		resp, err := http.Get(url)
		if err != nil {
			log.Warning("download failed", err)
		}
		defer resp.Body.Close()
		pix, _ := ioutil.ReadAll(resp.Body)
		out, _ := os.Create(path)
		defer out.Close()
		io.Copy(out, bytes.NewReader(pix))
		log.Infof("%v downloaded", file)
	}
	c.File(fmt.Sprintf("image/%s", file))
}
