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

var idalloc_id map[string]uint64
var idalloc_timeout map[string]int64

//同步磁盘间隔
var idalloc_sync_duration int64 = 1

const filePathPrefix = "./idalloc_"

func (self *Idalloc) init() {
	if idalloc_id == nil {
		idalloc_id = make(map[string]uint64)
		idalloc_timeout = make(map[string]int64)
	}
}

func (self *Idalloc) Gen() (uint64, error) {
	if idalloc_id == nil {
		self.init()
	}
	var file *os.File
	var err error
	filePath := filePathPrefix + self.Type
	fmt.Println("Cache file is", filePath)
	if idalloc_id[self.Type] == 0 {
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
		id += 1000
		idalloc_id[self.Type] = id
	}

	idalloc_id[self.Type]++
	id_str := strconv.FormatUint(idalloc_id[self.Type], 10)

	if idalloc_timeout[self.Type] < time.Now().Unix()-idalloc_sync_duration {
		//无句柄
		if file == nil {
			file, err = openCacheFile(filePath)
			defer file.Close()
			if err != nil {
				fmt.Println("Open cache file failed", err)
			}
		}
		idalloc_timeout[self.Type] = time.Now().Unix()

		if file != nil {
			fmt.Println("save "+filePath, "is", id_str, "to file")
			_, err = file.WriteAt([]byte(id_str), 0)
			if err != nil {
				fmt.Println("Save cache file failed", err, id_str)
			}
		} else {
			fmt.Println("Need save cache file, but `file` is nill")
		}

	}
	// time.Sleep(5 * time.Second)

	return idalloc_id[self.Type], nil
}

func (self *Idalloc) SyncCacheAll() (bool, error) {

	for key, value := range idalloc_id {
		filePath := filePathPrefix + key
		file, err := openCacheFile(filePath)
		defer file.Close()
		if err != nil {
			fmt.Println("Open cache file failed", err)
		}
		id_str := strconv.FormatUint(value, 10)
		fmt.Println("save "+filePath, "is", id_str, "to file")
		_, err = file.WriteAt([]byte(id_str), 0)
		if err != nil {
			fmt.Println("Save cache file failed", err, id_str)
		}
	}

	return true, nil
}

func openCacheFile(filePath string) (file *os.File, err error) {

	file, err = os.OpenFile(filePath, os.O_RDWR, 0)
	fmt.Println(file)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("File", filePath, "is not Exist, now create it")
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return file, nil
}
