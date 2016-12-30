// YFTool
package YFTool

import (
	log "YFTool/RpcLogs"
	"crypto/md5"
	"database/sql"
	"fmt"
	"hash/crc32"
	"runtime"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

func CheckError(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Panic(file, ":", runtime.FuncForPC(funcName).Name(), ":", line, " ", err)
		} else {
			log.Panic(err)
		}
	}
}

func AsyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}

func PrintRecover() {
	log.Println("PrintRecover")
	if e := recover(); e != nil {
		i := 0
		stackErrInfo := "StackErrInfo begin:"
		for {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				stackErrInfo += "\n"
				stackErrInfo += fmt.Sprintln("[", i, "]", file, ":", runtime.FuncForPC(funcName).Name(), ":", line)
				i++
			} else {
				stackErrInfo += "\n" + fmt.Sprintln("[errInfo]", e)
				log.Println(stackErrInfo)
				break
			}
		}
	}
}

func GetStringRemainder(id string, number uint32) uint32 {
	md5Handler := md5.New()
	md5Handler.Write([]byte(id))
	data := md5Handler.Sum(nil)
	return crc32.ChecksumIEEE(data) % number
}

func GetIntRemainder(id int, number int) int {
	if id < 0 {
		id = -id
	}
	return id % number
}

func NewRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 480 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

func GetMysqlDB(Name string) (*sql.DB, error) {
	db, err := sql.Open("mysql", Name)
	if err != nil {
		log.Error(err.Error())
	} else {
		err = db.Ping()
		if err != nil {
			// do something here
			log.Error("db is err:", err)
		}
	}
	return db, err
}

func CloseMysqlDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
