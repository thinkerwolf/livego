package httpopera

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thinkerwolf/livego/configure"
	"github.com/thinkerwolf/livego/utils"
)

/**
 * @api /api/record/folders
 * @apiParam start 分页开始
 * @apiParam limit 长度
 * @apiParam sort  排序字段
 * @apiParam order 顺序
 * @apiParam q 查询字符串
 */
func RecordFolders(ctx *gin.Context) {
	ffmpegCfg := configure.RtmpServercfg.Ffmpeg
	mpegPath := ffmpegCfg.Dir_path

	qs := ctx.Query("q")

	folders := make([]interface{}, 0)
	if len(mpegPath) > 0 {

		visit := func(folders *[]interface{}, filter utils.FileFilter) filepath.WalkFunc {
			return func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if path == mpegPath {
					return nil
				}
				if !info.IsDir() {
					return nil
				}

				folder := path[len(mpegPath):]
				if filter(qs, folder, info) {
					modTime := info.ModTime().Local().Format("20060102150405")
					*folders = append(*folders, map[string]interface{}{
						"folder": folder,
						"date":   modTime,
					})
				}

				return filepath.SkipDir
			}
		}

		filepath.Walk(mpegPath, visit(&folders, filter))
	}
	pr := utils.NewPageResult(folders)
	pr = pr.Sort(ctx.Query("sort"), ctx.Query("order"))

	start := 0
	limit := 200
	if ctx.Query("start") != "" {
		start, _ = strconv.Atoi(ctx.Query("start"))
	}
	if ctx.Query("limit") != "" {
		l, err := strconv.Atoi(ctx.Query("limit"))
		if err == nil {
			limit = l
		}
	}
	pr = pr.Slice(limit, start)
	ctx.IndentedJSON(http.StatusOK, pr)
}

/**
 * @api /api/record/files 获取文件夹下所有的录像文件
 * @apiParam folder 文件夹
 * @apiParam q 查询字符串
 * @apiParam start 分页开始
 * @apiParam limit 分页大小
 * @apiParam order 排序
 * @apiParam sort 排序字段
 */
func RecordFiles(ctx *gin.Context) {
	ffmpegCfg := configure.RtmpServercfg.Ffmpeg
	mpegPath := ffmpegCfg.Dir_path

	folder := ctx.Query("folder")
	dirPath := filepath.Join(ffmpegCfg.Dir_path, folder)
	qs := ctx.Query("q")
	files := make([]interface{}, 0)
	if mpegPath != "" {
		visit := func(files *[]interface{}, filter utils.FileFilter) filepath.WalkFunc {
			return func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if path == dirPath {
					return nil
				}
				if !info.IsDir() {
					fil := path[len(mpegPath):]
					if filter(qs, fil, info) {
						modTime := info.ModTime().Local().Format("20060102150405")
						*files = append(*files, map[string]interface{}{
							"file": fil,
							"date": modTime,
						})
					}
				}
				return nil
			}
		}

		filepath.Walk(dirPath, visit(&files, filter))
	}

	pr := utils.NewPageResult(files)

	pr = pr.Sort(ctx.Query("sort"), ctx.Query("order"))

	start := 0
	limit := 200
	if ctx.Query("start") != "" {
		start, _ = strconv.Atoi(ctx.Query("start"))
	}
	if ctx.Query("limit") != "" {
		l, err := strconv.Atoi(ctx.Query("limit"))
		if err == nil {
			limit = l
		}
	}
	pr = pr.Slice(limit, start)

	ctx.IndentedJSON(http.StatusOK, pr)
}

// 过滤func
func filter(q string, path string, info os.FileInfo) bool {
	if len(q) <= 0 {
		return true
	}
	keyValues := strings.Split(q, ",")
	kv := make(map[string]string, 0)
	for _, keyValue := range keyValues {
		ps := strings.Split(keyValue, "=")
		if len(ps) > 1 {
			kv[ps[0]] = ps[1]
		}
	}

	log.Printf("query kv %v", kv)

	if kv["sd"] != "" {
		date, err := time.Parse("20060102150405", kv["sd"])
		if err != nil {
			log.Println(err)
		} else {
			if info.ModTime().After(date) == false {
				return false
			}
		}
	}

	if kv["ed"] != "" {
		date, err := time.Parse("20060102150405", kv["ed"])
		if err != nil {
			log.Println(err)
		} else {
			if info.ModTime().Before(date) == false {
				return false
			}
		}
	}

	if kv["name"] != "" {
		reg := regexp.MustCompile(kv["name"])
		if len(reg.FindStringSubmatch(path)) <= 0 {
			return false
		}
	}
	return true
}
