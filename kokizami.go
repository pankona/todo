package kokizami

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pankona/kokizami/models"
	"github.com/xo/xoutil"
)

type Kokizami struct {
	DBPath string
}

func (k *Kokizami) execWithDB(f func(db models.XODB) error) error {
	models.XOLog = func(s string, p ...interface{}) {
		fmt.Printf("-------------------------------------\nQUERY: %s\n  VAL: %v\n", s, p)
	}

	conn, err := sql.Open("sqlite3", k.DBPath)
	if err != nil {
		return err
	}
	defer func() {
		e := conn.Close()
		if e != nil {
			log.Printf("failed to close DB connection: %v", e)
		}
	}()

	return f(conn)
}

func mustParse(format, value string) time.Time {
	t, err := time.Parse(format, value)
	if err != nil {
		panic(err)
	}
	return t
}

func (k *Kokizami) Start(desc string) (*models.Kizami, error) {
	var ki *models.Kizami
	return ki, k.execWithDB(func(db models.XODB) error {
		entry := &models.Kizami{
			Desc:      desc,
			StartedAt: xoutil.SqTime{time.Now()},
			StoppedAt: xoutil.SqTime{mustParse("2006-01-02 15:04:05", "1970-01-01 00:00:00")},
		}
		err := entry.Insert(db)
		if err != nil {
			return err
		}
		ki, err = models.KizamiByID(db, entry.ID)
		return err
	})
}

func (k *Kokizami) Get(id int) (*models.Kizami, error) {
	var ki *models.Kizami
	var err error
	return ki, k.execWithDB(func(db models.XODB) error {
		ki, err = models.KizamiByID(db, id)
		return err
	})
}

func (k *Kokizami) Edit(ki *models.Kizami) (*models.Kizami, error) {
	return ki, k.execWithDB(func(db models.XODB) error {
		a, err := models.KizamiByID(db, ki.ID)
		if err != nil {
			return err
		}

		*a = *ki

		err = a.Update(db)
		if err != nil {
			return err
		}
		ki, err = models.KizamiByID(db, ki.ID)
		return err
	})
}

func (k *Kokizami) Stop(id int) error {
	return k.execWithDB(func(db models.XODB) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		ki.StoppedAt = xoutil.SqTime{time.Now()}
		return ki.Update(db)
	})
}

func (k *Kokizami) StopAll() error {
	return k.execWithDB(func(db models.XODB) error {
		t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")
		if err != nil {
			return err
		}
		ks, err := models.KizamisByStoppedAt(db, xoutil.SqTime{t})
		if err != nil {
			return err
		}
		now := time.Now()
		for i := range ks {
			ks[i].StoppedAt = xoutil.SqTime{now}
			if err := ks[i].Update(db); err != nil {
				return err
			}
		}
		return nil
	})
}

func (k *Kokizami) Delete(id int) error {
	return k.execWithDB(func(db models.XODB) error {
		ki, err := models.KizamiByID(db, id)
		if err != nil {
			return err
		}
		return ki.Delete(db)
	})
}

func (k *Kokizami) List() ([]*models.Kizami, error) {
	var ks []*models.Kizami
	return ks, k.execWithDB(func(db models.XODB) error {
		var err error
		ks, err = models.AllKizami(db)
		return err
	})
}
