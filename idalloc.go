package idalloc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

//发号池
type Pooler interface {
	init()
	Debug(bool)
	BootAutoIncre(uint64)
	Gen() (uint64, error)
	SyncCacheAll() (bool, error)
	openCacheFile(string) (*os.File, error)
}

//发号器
type Pool struct {
	Type string `发号器类型`
}

//缓存ID
var idallocId map[string]uint64

//缓存刷新时间
var idallocTimeout map[string]int64

//Debug 模式
var isDebug bool = false

//启动自增步长
var bootAutoIncre uint64 = 1000

//同步磁盘间隔
var idallocSyncDuration int64 = 1

//缓存文件
const filePathPrefix = "./idalloc_"

func (self *Pool) init() {
	if idallocId == nil {
		idallocId = make(map[string]uint64)
		idallocTimeout = make(map[string]int64)
	}
}

//是否为Debug模式
func (self *Pool) Debug(b bool) {
	isDebug = b
}

//启动自增
func (self *Pool) BootAutoIncre(n uint64) {
	bootAutoIncre = n
}

//生成自增ID
func (self *Pool) Gen() (uint64, error) {
	if idallocId == nil {
		self.init()
	}
	var file *os.File
	var err error
	var filePath string

	if idallocId[self.Type] == 0 {
		filePath = filePathPrefix + self.Type
		if isDebug {
			fmt.Println("[Debug]", "Cache file is", filePath)
		}
		file, err = openCacheFile(filePath)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		reader := bufio.NewReader(file)
		con, _ := reader.ReadString('\n') //读取一行
		id, _ := strconv.ParseUint(con, 10, 64)
		//防止ID回流
		id += bootAutoIncre
		idallocId[self.Type] = id
	}

	idallocId[self.Type]++

	if idallocTimeout[self.Type] < time.Now().Unix()-idallocSyncDuration {
		//无句柄
		if file == nil {
			file, err = openCacheFile(filePath)
			defer file.Close()
			if err != nil {
				fmt.Println("[Error]", "Open cache file failed", err)
			}
		}
		idallocTimeout[self.Type] = time.Now().Unix()

		if file != nil {
			filePath = filePathPrefix + self.Type
			id_str := strconv.FormatUint(idallocId[self.Type], 10)
			if isDebug {
				fmt.Println("[Debug]", "save "+filePath, "is", id_str, "to file")
			}
			_, err = file.WriteAt([]byte(id_str), 0)
			if err != nil {
				fmt.Println("[Error]", "Save cache file failed", err, id_str)
			}
		} else {
			fmt.Println("[Warn]", "Need save cache file, but `file` is nill")
		}

	}

	return idallocId[self.Type], nil
}

//同步所有缓存到磁盘
func (self *Pool) SyncCacheAll() (bool, error) {

	for key, value := range idallocId {
		filePath := filePathPrefix + key
		file, err := openCacheFile(filePath)
		defer file.Close()
		if err != nil {
			fmt.Println("[Error]", "Open cache file failed", err)
		}
		id_str := strconv.FormatUint(value, 10)
		fmt.Println("save "+filePath, "is", id_str, "to file")
		_, err = file.WriteAt([]byte(id_str), 0)
		if err != nil {
			fmt.Println("[Error]", "Save cache file failed", err, id_str)
		}
	}

	return true, nil
}

//打开缓存文件
func openCacheFile(filePath string) (file *os.File, err error) {

	file, err = os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil && os.IsNotExist(err) {
		if isDebug {
			fmt.Println("[Debug]", "File", filePath, "is not Exist, now create it")
		}
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return file, nil
}
