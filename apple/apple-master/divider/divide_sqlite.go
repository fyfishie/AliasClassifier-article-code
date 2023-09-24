/*
 * @Author: fyfishie
 * @Date: 2023-05-07:12
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-10:08
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package divider

import (
	"aliasParseMaster/geocode"
	"aliasParseMaster/idgen"
	"aliasParseMaster/lib"
	"bufio"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Divider struct {
	IPList         []string
	bufferLen      int
	db             *sql.DB
	stmt           *sql.Stmt
	locationIDGen  *idgen.Generator
	locationDisMap map[lib.LocationID]int
	//map[country\|city][]IP{}
	// CityIPMap  map[string][]lib.IP
}

func NewDivider(BufferLen int) *Divider {
	if BufferLen < 100 {
		BufferLen = 100
	}
	return &Divider{
		bufferLen:      BufferLen,
		locationIDGen:  idgen.NewIdGen().Run(),
		locationDisMap: make(map[string]int),
	}
}

func (d *Divider) Divide_sql(rdpath string) error {
	s := geocode.NewXdbSearcher()
	s.SearchInit("./ip2location.xdb")
	buffer := []string{}
	rfi, err := os.OpenFile(rdpath, os.O_RDONLY, 0000)
	if err != nil {
		return err
	}
	defer rfi.Close()

	inChan := make(chan map[lib.IP]lib.LocationID)
	err = d.prepareDB()
	d.runDB(inChan)

	scaner := bufio.NewScanner(rfi)
	if err != nil {
		logrus.Errorf("error while prepare database, err = %v\n", err.Error())
		return err
	}

	count := 0
	for {
		for i := 0; i < d.bufferLen && scaner.Scan(); i++ {
			count++
			if count%1000 == 0 {
				fmt.Printf("count: %v\n", count)
			}
			buffer = append(buffer, scaner.Text())
		}
		if len(buffer) == 0 {
			break
		}
		locations, err := s.GetCitysIDByXdb(buffer)
		if err != nil {
			logrus.Errorf("error while request city info, err = %v\n", err.Error())
			continue
		}
		inChan <- locations
		buffer = []string{}
	}
	close(inChan)
	return nil
}

func (d *Divider) cleanDB() error {
	if d.stmt != nil {
		d.stmt.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
	return os.Remove("./divider.db")
}
func (d *Divider) prepareDB() error {
	d.cleanDB()
	db, err := sql.Open("sqlite3", "./divider.db")
	if err != nil {
		logrus.Errorf("error while open sqlite db, err = %v\n", err.Error())
		return err
	}
	d.db = db
	sqlStmt := `
	create table ip_city_map (ip text not null primary key, city text);
	delete from ip_city_map;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logrus.Errorf("%q: %s\n", err, sqlStmt)
		d.cleanDB()
		return err
	}
	stmt, err := db.Prepare("insert into ip_city_map(ip, city) values(?, ?)")
	if err != nil {
		logrus.Errorf("error while prepare stmt for insert item, err = %v\n", err.Error())
		d.cleanDB()
		return err
	}
	d.stmt = stmt
	return nil
}

func (d *Divider) runDB(inChan chan map[lib.IP]lib.LocationID) {
	go func() {
		for m := range inChan {
			for ip, locID := range m {
				d.stmt.Exec(ip, locID)
			}
		}
		d.cleanDB()
	}()
}
