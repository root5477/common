package common

import (
	"fmt"
	"github.com/gofrs/uuid"
	"strconv"
	"strings"
	"github.com/xormplus/xorm"
	"xorm.io/core"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"sync"
)

type UUID struct {
	Id   int64  `json:"id" xorm:"id"`
	Uuid string `json:"uuid" xorm:"uuid"`
}

var PdnsSqlEngine *xorm.Engine
var wg sync.WaitGroup

func init() {
	var err error
	PdnsSqlEngine, err = xorm.NewEngine("mysql", "root:123456@tcp(10.226.133.105:3358)/cdmstest?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Println("create mysql engine failed, err is:", err)
		return
	}
	errPing := PdnsSqlEngine.Ping()
	if errPing != nil {
		fmt.Println("ping test failed, err is:", errPing)
		return
	}
	fmt.Println("connect to mysql success!")
	PdnsSqlEngine.ShowSQL(true)
	PdnsSqlEngine.SetMapper(core.GonicMapper{})
}

var t1 time.Time

//func main() {
//	t1 = time.Now()
//	wg.Add(3)
//	go InsertUUid(2000)
//	go InsertUUid(2000)
//	go InsertUUid(2000)
//	go InsertUUid(2000)
//	wg.Wait()
//}

func InsertUUid(num int)  {
	defer wg.Done()
	for i := 1; i <= num; i ++ {
		id := GetUUID10()
		uid := UUID {
			Uuid:id,
		}
		_, err := PdnsSqlEngine.Insert(&uid)
		if err != nil {
			fmt.Printf("insetr %v failed, err is:%v\n", uid, err)
			return
		}
	}
	fmt.Println(time.Since(t1))
}

func GetUUID10() string {

	letterPool := []string {"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
		"t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9"}
	id := uuid.Must(uuid.NewV4())

	idStr := fmt.Sprintf("%v", id)
	idStr = strings.ReplaceAll(idStr, "-", "")

	id10 := ""
	for i := 0; i < 10; i++ {
		tmpSub := idStr[i*3 : i*3+3]
		tmpSub16 := fmt.Sprintf("%x", tmpSub)
		tmpSub16Int, err := strconv.Atoi(tmpSub16)
		if err != nil {
			fmt.Println("Atoi failed, err:", err)
			return ""
		}
		id10 += letterPool[tmpSub16Int % 36]
	}
	return id10
}
