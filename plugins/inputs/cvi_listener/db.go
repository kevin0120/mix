package cvi_listener

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"github.com/masami10/rush/utils/cvi"
)

type CURDB struct {
	File string
	//db   *sql.DB
}

//func (cur *CURDB) Init() (error) {
//	db, err := sql.Open("sqlite3", cur.File)
//	if err != nil {
//		return err
//	}
//
//	cur.db = db
//}

func (cur *CURDB) PreSave(sn string, workorder_id int, screw_id string) (error) {
	db, err := sql.Open("sqlite3", cur.File)
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into curs(sn, workorderid, screw_id, has_result, saved, update_time, result_data) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	date, time := cvi.GetDateTime()
	stmt.Exec(sn, workorder_id, screw_id, false, false, date + " " + time, "")
	tx.Commit()

	return nil
}

func (cur *CURDB) UpdateResultData(sn string, workorder_id int, screw_id string, result string) (error) {
	db, err := sql.Open("sqlite3", cur.File)
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update curs set result_data = ?, has_result = ? where sn = ? and workorder_id = ? and screw_id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	date, time := cvi.GetDateTime()
	stmt.Exec(result, true, sn, workorder_id, screw_id)
	tx.Commit()

	return nil
}

func (cur *CURDB) Saved(sn string, workorder_id int, screw_id string) (error) {
	db, err := sql.Open("sqlite3", cur.File)
	if err != nil {
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update curs set saved = ? where sn = ? and workorder_id = ? and screw_id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	date, time := cvi.GetDateTime()
	stmt.Exec(true, sn, workorder_id, screw_id)
	tx.Commit()

	return nil
}

