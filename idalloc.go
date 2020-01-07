package idalloc

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//Pooler 发号器池
type Pooler interface {
	init()
	genderate() (uint64, error)
	cache() (bool, error)
	Run(string, chan<- uint64)
}

//Pool 发号器池
type Pool struct{}

//缓存ID
var idallocID map[string]uint64

//缓存刷新时间
var idallocTimeout map[string]int64

//Debug 模式
var isDebug bool = false

//启动自增步长
var bootAutoIncre uint64 = 1000

//同步磁盘间隔
var idallocSyncDuration int64 = 1

//缓存文件
const cacheFilePrefix = "./idalloc_"

//Debug 是否为Debug模式
func (pl *Pool) Debug(b bool) {
	isDebug = b
}

//BootAutoIncre 设置启动自增步长
func (pl *Pool) BootAutoIncre(n uint64) {
	bootAutoIncre = n
}

//初始化
func (pl *Pool) init() {
	if idallocID == nil {
		idallocID = make(map[string]uint64)
		idallocTimeout = make(map[string]int64)
	}
}

//Run 开启一个发号器协程，通道只能出
func (pl *Pool) Run(t string, out chan<- uint64) {
	go func() {
		for {
			i, err := pl.genderate(t)
			if nil != err {
				panic("Generate id error")
			}
			out <- i
		}

	}()

	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		log.Println(`[Before Quite]Sync cache ` + t)
		pl.cache(t)
		os.Exit(0)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
}

//genderate 生成自增ID
func (pl *Pool) genderate(t string) (uint64, error) {
	if idallocID == nil {
		pl.init()
	}
	var file *os.File
	var err error
	var cf string
	cf = cacheFilePrefix + t

	if idallocID[t] == 0 {

		if isDebug {
			log.Println("[Debug]", "Cache file is", cf)
		}
		file, err = openCacheFile(cf)
		defer file.Close()
		if err != nil {
			return 0, errors.New("Open cache file failed")
		}

		reader := bufio.NewReader(file)
		con, _ := reader.ReadString('\n') //读取一行
		id, _ := strconv.ParseUint(con, 10, 64)
		//防止ID回流
		id += bootAutoIncre
		idallocID[t] = id
	}

	idallocID[t]++

	if idallocTimeout[t] < time.Now().Unix()-idallocSyncDuration {
		//无句柄
		if file == nil {
			file, err = openCacheFile(cf)
			defer file.Close()
			if err != nil {
				return 0, errors.New("Open cache file failed")
			}
		}
		idallocTimeout[t] = time.Now().Unix()

		if file != nil {
			s := strconv.FormatUint(idallocID[t], 10)
			if isDebug {
				log.Println("[Debug]", "save "+cf, "is", s, "to file")
			}
			_, err = file.WriteAt([]byte(s), 0)
			if err != nil {
				panic("Save cache file failed" + err.Error())
			}
		} else {
			// fmt.Println("[Warn]", "Need save cache file, but `file` is nill")
			panic("Need save cache file, but `file` is nill")
		}

	}

	return idallocID[t], nil
}

//cache 同步缓存到磁盘
func (pl *Pool) cache(t string) (bool, error) {

	cf := cacheFilePrefix + t
	file, err := openCacheFile(cf)
	defer file.Close()
	if err != nil {
		return false, errors.New("Open cache file failed")
	}
	s := strconv.FormatUint(idallocID[t], 10)
	_, err = file.WriteAt([]byte(s), 0)
	if err != nil {
		panic("Save cache file failed. " + err.Error())
	}
	return true, nil
}

//openCacheFile 打开缓存文件
func openCacheFile(p string) (file *os.File, err error) {

	file, err = os.OpenFile(p, os.O_RDWR, 0)
	if err != nil && os.IsNotExist(err) {
		if isDebug {
			log.Println("[Debug]", "File", p, "is not Exist, now create it")
		}
		file, err = os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
	}
	return file, nil
}
