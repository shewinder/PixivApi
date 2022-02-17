package main

import (
	"fmt"
	"os"
	"pixiv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RefreshToken string `yaml:"refresh-token"`
}

func parse(path string) *Config {
	file, err := os.ReadFile(path)
	config := &Config{}
	if err == nil {
		yaml.NewDecoder(strings.NewReader(string(file))).Decode(config)
		if err != nil {
			log.Fatal("配置文件不合法!", err)
		}
	} else {
		os.WriteFile("config.yml", []byte("refresh-token: xxxxxx"), 0o644)
		fmt.Println("配置文件已生成， 请填写配置后重新启动")
		os.Exit(0)
	}
	return config
}

func main() {
	log.Info("Pixiv APi is runing")
	gin.SetMode(gin.ReleaseMode)
	conf := parse("config.yml")
	//fmt.Println(conf.RefreshToken)
	pixiv.InitAuth(conf.RefreshToken)
	go func() {
		ticker := time.NewTicker(3000 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			log.Info("refreshing token")
			pixiv.RefreshToken() // 定期刷新，避免token过期
		}
	}()
	router := gin.Default()

	router.GET("/pixiv/rank", getIllustRanking)
	router.GET("/pixiv/illust_detail", getIllustDetail)
	router.GET("/pixiv/following", getIllustFollow)
	router.GET("/pixiv/user", getUserIllusts)
	router.GET("/pixiv/:file", getPidFile)

	router.Run(":9500")
}
