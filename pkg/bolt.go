package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	bolt "github.com/coreos/bbolt"
)

var (
	Db   *bolt.DB
	Mmap map[string]string
	err  error
)

const demo string = `[{
	"id": 0,
	"index": [0],
	"label": "1",
	"children": [{
		"id": 1,
		"index": [0, 0],
		"label": "1-1",
		"children": [{
			"id": 2,
			"index": [0, 0, 0],
			"label": "1-1-1",
			"children": null
		}]
	}, {
		"id": 3,
		"index": [0, 1],
		"label": "1-2",
		"children": [{
			"id": 4,
			"index": [0, 1, 0],
			"label": "1-2-1",
			"children": null
		}]
	}]
}, {
	"id": 9,
	"index": [1],
	"label": "2",
	"children": [{
		"id": 5,
		"index": [1, 0],
		"label": "2-1",
		"children": [{
			"id": 6,
			"index": [1, 0, 0],
			"label": "2-1-1",
			"children": null
		}]
	}, {
		"id": 7,
		"index": [1, 1],
		"label": "2-2",
		"children": [{
			"id": 8,
			"index": [1, 1, 0],
			"label": "2-2-1",
			"children": null
		}]
	}]
}]`

var (
	dbname string
	db     string = "files"
	count  string = "count"
	home   string
)

func Delete() error {
	//delete file
	err = os.Remove(dbname)

	// err = DeleteBucket(db)
	return err
}

func init() {

	home, err = Home()
	if err != nil {
		panic(err)
	}
	dbname = fmt.Sprintf("%s/%s", home, ".search.db")

	Mmap = map[string]string{}
	Db, err = bolt.Open(dbname, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	//init tables

	err = CreateBucket(db)
	if err != nil {
		log.Debug(err.Error())
	}
	err = CreateBucket(count)
	if err != nil {
		log.Debug(err.Error())
	}
	// err = CreateBucket(data)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	// panic(err)
	// }
	// err = CreateBucket(usertable)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// //初始化user表
	// us, _ := GetValueByBucketName(usertable, username)
	// if len(us) == 0 {
	// 	log.Warning("初始化用户表数据")
	// 	AddKeyValueByBucketName(usertable, username, password)
	// }
	// //初始化索引数据 即创建
	// ds, _ := GetValueByBucketName(data, data)
	// if len(ds) == 0 {
	// 	log.Warning("初始化索引数据")
	// 	AddKeyValueByBucketName(data, data, string(demo))
	// }
	GetAllTables()
	// log.Println(Mmap)
	log.Debug("init db logging success")
}

func CreateBucket(tablename string) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(tablename))
		if err != nil {
			return err
		}
		//如果没有index表和index key则立
		return nil
	})
	AddTables(tablename)
	return err
}

func DeleteBucket(tablename string) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(tablename))
		if err != nil {
			return fmt.Errorf("delete bucket: %s ", err.Error())
		}
		return nil
	})
	DeleteTables(tablename)
	return err
}

func AddTables(tablename string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Put([]byte(tablename), []byte(tablename))
		Mmap[string(tablename)] = string(tablename)
		return err
	})
}

func GetAllTables() (map[string]string, error) {
	tmp := map[string]string{}
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp[string(k)] = string(v)
			Mmap[string(k)] = string(v)
		}
		return nil
	})
	return tmp, err
}

func GetAllByTables(name string) (map[string]string, error) {
	tmp := map[string]string{}
	err := Db.View(func(tx *bolt.Tx) error {
		if _, ok := Mmap[name]; !ok {
			return errors.New(fmt.Sprintf("%s not exist", name))
		}
		b := tx.Bucket([]byte(name))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp[string(k)] = string(v)
			Mmap[string(k)] = string(v)
		}
		return nil
	})
	return tmp, err
}

func DeleteTables(tablename string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Delete([]byte(tablename))
		delete(Mmap, tablename)
		return err
	})
}

func SearchAll(key string) error {
	return Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(db))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			// if strings.Contains(string(k), key) {
			// 	fmt.Printf("key=%s, value=%s\n", k, v)
			// }
			matched, _ := regexp.MatchString(key, string(k))

			if matched {
				fmt.Printf("key=%s, value=%s\n", k, v)
			}
		}

		return nil
	})
}

func SearchPrefix(key string) error {
	return Db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(db)).Cursor()

		prefix := []byte(key)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
}

func AddKeyValueBatch(key, value string, wg *sync.WaitGroup) error {
	// log.Println(Mmap)
	return Db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Put([]byte(key), []byte(value))
		defer wg.Done()
		return err
	})
}

func AddKeyValue(key, value string) error {
	// log.Println(Mmap)
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func AddKeyValueByBucketName(table, key, value string) error {
	// fmt.Println(Mmap)
	if _, ok := Mmap[table]; !ok {
		fmt.Println(fmt.Printf("%s is not exist\n", table))
		CreateBucket(table)
	}
	// log.Println(Mmap)
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func AddKeyValueByBucketNameAuto(table, key, value string) error {

	return Db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		var err error
		if _, ok := Mmap[table]; !ok {
			b, err = tx.CreateBucket([]byte(table))
			if err != nil {
				return err
			}
		} else {
			b = tx.Bucket([]byte(table))
			if err != nil {
				return err
			}
		}

		err = b.Put([]byte(key), []byte(value))
		return err
	})
}

func GetValueByBucketName(table, key string) ([]byte, error) {
	var value []byte
	var err error
	Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if _, ok := Mmap[table]; ok {
			value = b.Get([]byte(key))
		} else {
			err = errors.New(fmt.Sprintf("table %s not exist", table))
		}
		return err
	})
	return value, err
}

func DeleteKeyValueByBucket(table, key string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Delete([]byte(key))
		return err
	})
}
