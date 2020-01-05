package idalloc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Idalloc struct {
	Type string
}

//缓存ID
var idalloc_id map[string]uint64

//缓存刷新时间
var idalloc_timeout map[string]int64

//Debug 模式
var is_debug bool = false

//启动自增步长
var boot_auto_incre uint64 = 1000

//同步磁盘间隔
var idalloc_sync_duration int64 = 1

//缓存文件
const filePathPrefix = "./idalloc_"

func (self *Idalloc) Debug(b bool) {
	is_debug = b
}

func (self *Idalloc) init() {
	if idalloc_id == nil {
		idalloc_id = make(map[string]uint64)
		idalloc_timeout = make(map[string]int64)
	}
}

func (self *Idalloc) BootAutoIncre(n uint64) {
	boot_auto_incre = n
}

func (self *Idalloc) Gen() (uint64, error) {
	if idalloc_id == nil {
		self.init()
	}
	var file *os.File
	var err error
	var filePath string

	if idalloc_id[self.Type] == 0 {
		filePath = filePathPrefix + self.Type
		if is_debug {
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
		id += boot_auto_incre
		idalloc_id[self.Type] = id
	}

	idalloc_id[self.Type]++

	if idalloc_timeout[self.Type] < time.Now().Unix()-idalloc_sync_duration {
		//无句柄
		if file == nil {
			file, err = openCacheFile(filePath)
			defer file.Close()
			if err != nil {
				fmt.Println("[Error]", "Open cache file failed", err)
			}
		}
		idalloc_timeout[self.Type] = time.Now().Unix()

		if file != nil {
			filePath = filePathPrefix + self.Type
			id_str := strconv.FormatUint(idalloc_id[self.Type], 10)
			if is_debug {
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

	return idalloc_id[self.Type], nil
}

func (self *Idalloc) SyncCacheAll() (bool, error) {

	for key, value := range idalloc_id {
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

func openCacheFile(filePath string) (file *os.File, err error) {

	file, err = os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil && os.IsNotExist(err) {
		if is_debug {
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
