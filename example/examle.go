package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hyahm/gomysql"
)

type Category struct {
	Id        int64  `json:"id" db:"id,omitempty"`
	English   string `json:"english" db:"english"` // 分类英文名， 文件夹命名 唯一索引
	Name      string `json:"name" db:"name"`
	ServerId  int64  `json:"server_id" db:"server_id"` // 服务器信息id
	SubCateId []int  `json:"subcate" db:"subcate"`
	Uids      []int  `json:"uids" db:"uids"`
}

type Postparam struct {
	ID            int64  `json:"id" db:"id,omitempty"`
	Duration      int    `json:"duration"  db:"duration"`
	WebpCount     int    `json:"webp_count"  db:"webp_count"`
	PlayUrl       string `json:"play_url"  db:"play_url"`
	DownloadUrl   string `json:"download_url"  db:"download_url"`
	ThumbLongview string `json:"thumb_longview"  db:"thumb_longview"`
	Thumbnail     string `json:"thumbnail"  db:"thumbnail"`
	ThumbVer      string `json:"thumb_ver"  db:"thumb_ver"`
	ThumbHor      string `json:"thumb_hor"  db:"thumb_hor"`
	Cover         string `json:"cover"  db:"cover"`
	Webp          string `json:"webp"  db:"webp"`
	Preview       string `json:"preview"  db:"preview"`
	ThumbSeries   []int  `json:"thumb_series"  db:"thumb_series"`
}

var (
	conf = &gomysql.Sqlconfig{
		Host:         "192.168.50.71",
		Port:         3306,
		UserName:     "cander",
		Password:     "123456",
		DbName:       "shop",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
)

func main() {
	db, err := conf.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	cate := &Category{}
	err = db.Select(&cate, fmt.Sprintf("select * from category where id=%d", 51))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cate.SubCateId)
	fmt.Println(cate.Id)
	fmt.Println(cate.ServerId)
	fmt.Println(cate.Name)
	fmt.Println(cate.English)

	post := &Postparam{}
	err = db.Select(&post, "select * from postparam where id=118186")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(post.ThumbSeries)
}
