package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

var (
	num  int
	dnum int
)

func init() {
	num = 0
	dnum = 0
}

func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return err
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	wg := sync.WaitGroup{}
	wg.Add(len(rd))

	for n, fi := range rd {
		if fi.IsDir() {
			dnum++
			fullDir := pathname + string(os.PathSeparator) + fi.Name()
			// fmt.Println(fullDir)
			go func(i int) {
				Refresh(fullDir, &wg)
				wg.Done()
			}(n)
		} else {
			num++
			fullName := pathname + string(os.PathSeparator) + fi.Name()
			// fmt.Println(fullName)
			go AddKeyValueBatch(fullName, "file", &wg)
		}
	}
	wg.Wait()
	AddKeyValueByBucketName("count", "count", fmt.Sprintf("DIR %d FILES %d", dnum, num))
	s.Stop()
	return nil
}

func Refresh(pathname string, wg *sync.WaitGroup) error {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Errorln("read dir fail:", err)
		return err
	}
	wg.Add(len(rd))
	for _, fi := range rd {
		if fi.IsDir() {
			dnum++
			fullDir := pathname + string(os.PathSeparator) + fi.Name()
			go AddKeyValueBatch(fmt.Sprintf("%d %s", dnum, fi.Name()), fmt.Sprintf("D|%s", fullDir), wg)
			// log.Infoln("dir ", fullDir)
			fmt.Printf("\r Dir %d File %d Goroutine %d D|%s", dnum, num, runtime.NumGoroutine(), fi.Name())
			err = Refresh(fullDir, wg)
			if err != nil {
				log.Errorln("read dir fail:", err)
				return err
			}
		} else {
			num++
			fullName := pathname + string(os.PathSeparator) + fi.Name()
			go AddKeyValueBatch(fmt.Sprintf("%d %s", num, fi.Name()), fmt.Sprintf("F|%s", fullName), wg)
			// fmt.Printf("\r %d ; hhhhh", num)
			// log.Debugln("file ", fullName)
			fmt.Printf("\r Dir %d File Goroutine %d %d F|%s", dnum, num, runtime.NumGoroutine(), fi.Name())
		}
	}
	return nil
}
