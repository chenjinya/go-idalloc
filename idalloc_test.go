package idalloc

import "testing"

import "os"

func TestCacheFile(t *testing.T) {

	cf := cacheFilePrefix + "test"
	file, err := openCacheFile(cf)
	defer file.Close()
	if err != nil {
		t.Error(`func openCacheFile error`)
	}

}

func TestRun(t *testing.T) {
	cf := "./idalloc_test"
	ida := Pool{}
	ida.BootAutoIncre(0)
	idc := make(chan uint64)
	os.Remove(cf)
	ida.Run("test", idc)
	if 1 != <-idc {
		t.Error(`Run goroutine idalloc id error`)
	}
	os.Remove(cf)
}
