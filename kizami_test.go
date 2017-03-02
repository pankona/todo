package kokizami

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

type DBMock struct {
	DBInterface
	mockOpenDB      func() error
	mockClose       func()
	mockCreateTable func() error
	mockStart       func(desc string) (*kizami, error)
	mockEdit        func(id int, field, newValue string) (*kizami, error)
	mockCount       func() (int, error)
	mockList        func() ([]*kizami, error)
	mockStop        func(id int) error
	mockDelete      func(id int) error
}

func (db *DBMock) openDB() error {
	return db.mockOpenDB()
}

func (db *DBMock) close() {
	db.mockClose()
}

func (db *DBMock) createTable() error {
	return db.mockCreateTable()
}

func (db *DBMock) start(desc string) (*kizami, error) {
	return db.mockStart(desc)
}

func (db *DBMock) edit(id int, field, newValue string) (*kizami, error) {
	return db.mockEdit(id, field, newValue)
}

func (db *DBMock) count() (int, error) {
	return db.mockCount()
}

func (db *DBMock) list(start, end int) ([]*kizami, error) {
	return db.mockList()
}

func (db *DBMock) stop(id int) error {
	return db.mockStop(id)
}

func (db *DBMock) delete(id int) error {
	return db.mockDelete(id)
}

// default mock implementation
func genDefaultDBMock() *DBMock {
	return &DBMock{
		mockOpenDB: func() error {
			return nil
		},
		mockClose: func() {
		},
		mockCreateTable: func() error {
			return nil
		},
		mockStart: func(desc string) (*kizami, error) {
			return &kizami{desc: "test"}, nil
		},
		mockEdit: func(id int, field, newValue string) (*kizami, error) {
			return &kizami{desc: "edited"}, nil
		},
		mockCount: func() (int, error) {
			return 3, nil
		},
		mockList: func() ([]*kizami, error) {
			t := make([]*kizami, 0, 0)
			t = append(t, &kizami{desc: "test0"})
			t = append(t, &kizami{desc: "test1"})
			t = append(t, &kizami{desc: "test2"})
			return t, nil
		},
		mockStop: func(id int) error {
			return nil
		},
		mockDelete: func(id int) error {
			return nil
		},
	}
}

func TestInitializeNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Error("failed to initialize")
	}
}

func TestInitializeError(t *testing.T) {
	dbmock := genDefaultDBMock()
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}

	err := initialize(dbmock, "")
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
}

func TestStartNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatal("failed to initialize")
	}
	k, err := Start("test")
	if err != nil {
		t.Error("Start returned error")
	}
	if k == nil {
		t.Error("Start returned nil")
	}
	if k.Desc() != "test" {
		t.Error("Start returned unexpected value")
	}
}

func TestNormalWithDB(t *testing.T) {
	fp, err := ioutil.TempFile("", "tmp_")
	if err != nil {
		t.Error("failed to create temp file")
	}
	defer func() {
		err := os.Remove(fp.Name())
		if err != nil {
			t.Error("failed to remove temporary file")
		}
	}()

	err = initialize(nil, fp.Name())
	if err != nil {
		t.Fatal("failed to initialize")
	}
	k, err := Start("test")
	if err != nil {
		t.Error("Start returned error")
	}
	if k.ID() != 1 {
		t.Error("Start returned unexpected kizami instance")
	}
	if k.Desc() != "test" {
		t.Error("Start returned unexpected kizami instance")
	}
	l, err := List()
	if err != nil {
		t.Error("List returned error")
	}
	if len(l) != 1 {
		t.Error("unexpected list length")
	}

	k, err = Edit(1, "desc", "edited")
	if err != nil {
		t.Error("Edit returned error")
	}
	if k.ID() != 1 {
		t.Error("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Error("Edit returned unexpected kizami")
	}

	k, err = Edit(1, "started_at", "2010-01-02 03:04:05")
	if err != nil {
		t.Error("Edit returned error")
	}
	if k.ID() != 1 {
		t.Error("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Error("Edit returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("Edit returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "1970-01-01 00:00:00" {
		t.Error("Edit returned unexpected kizami")
	}

	k, err = Edit(1, "stopped_at", "2011-01-02 03:04:05")
	if err != nil {
		t.Error("Edit returned error")
	}
	if k.ID() != 1 {
		t.Error("Edit returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Error("Edit returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("Edit returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "2011-01-02 03:04:05" {
		t.Error("Edit returned unexpected kizami")
	}

	err = Stop(1)
	if err != nil {
		t.Error("Stop returned error")
	}
	if k.ID() != 1 {
		t.Error("Stop returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Error("Stop returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("Stop returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Error("Stop returned unexpected kizami")
	}

	k, err = Get(1)
	if err != nil {
		t.Error("Stop returned error")
	}
	if k.ID() != 1 {
		t.Error("Stop returned unexpected kizami")
	}
	if k.Desc() != "edited" {
		t.Error("Stop returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("Stop returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Error("Stop returned unexpected kizami")
	}

	err = Delete(1)
	if err != nil {
		t.Error("Delete returned error")
	}
	l, err = List()
	if err != nil {
		t.Error("Delete returned error")
	}
	if len(l) != 0 {
		t.Error("List returned unexpected result")
	}

	k, err = Start("test")
	if err != nil {
		t.Error("Delete returned error")
	}

	err = StopAll()
	if err != nil {
		t.Error("Delete returned error")
	}
	k, err = Edit(2, "started_at", "2010-01-02 03:04:05")
	if err != nil {
		t.Error("Get returned error")
	}
	if k.ID() != 2 {
		t.Error("Get returned unexpected kizami")
	}
	if k.Desc() != "test" {
		t.Error("Get returned unexpected kizami")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("Get returned unexpected kizami")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") == "1970-01-01 00:00:00" {
		t.Error("Get returned unexpected kizami")
	}

	k, err = Edit(2, "stopped_at", "2010-01-02 04:04:05")
	if err != nil {
		t.Error("Get returned error")
	}
	if k.ID() != 2 {
		t.Error("ID returned unexpected value")
	}
	if k.Desc() != "test" {
		t.Error("Desc returned unexpected value")
	}
	if k.StartedAt().Format("2006-01-02 15:04:05") != "2010-01-02 03:04:05" {
		t.Error("StartedAt returned unexpected value")
	}
	if k.StoppedAt().Format("2006-01-02 15:04:05") != "2010-01-02 04:04:05" {
		t.Error("StoppedAt returned unexpected value")
	}
	if k.String() != "2\ttest\t2010-01-02 03:04:05\t2010-01-02 04:04:05\t1h0m0s" {
		t.Error("Error returned unexpected value. actual =", k.String())
	}
}

func TestStartError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
	k, err := Start("test")
	if k != nil {
		t.Error("k is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is not nil but this is not expected")
	}

	// openDB goes success but Start goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockStart = func(desc string) (*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatal("failed to initialize")
	}
	k, err = Start("test")
	if k != nil {
		t.Error("k is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is nil but this is not expected")
	}
}

func TestEditNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatal("failed to initialize")
	}
	k, err := Edit(0, "desc", "edited")
	if err != nil {
		t.Error("Edit returned error")
	}
	if k == nil {
		t.Error("Edit returned nil")
	}
	if k.Desc() != "edited" {
		t.Error("edit returned unexpected value")
	}
}

func TestEditError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Error("Initialize succeeded but this is not expected")
	}
	k, err := Edit(0, "desc", "edited")
	if k != nil {
		t.Error("k is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is not nil but this is not expected")
	}

	// openDB goes success but Edit goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockEdit = func(id int, field, newValue string) (*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Fatal("failed to initialize")
	}
	k, err = Edit(0, "desc", "edited")
	if k != nil {
		t.Error("k is not nil but this is not expected")
	}
	if err == nil {
		t.Error("error is nil but this is not expected")
	}
}

func TestListNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Error("initialize failed")
	}
	ks, err := List()
	if err != nil {
		t.Error("List returned error")
	}
	if ks == nil {
		t.Error("List returned nil")
	}
	if ks[0].Desc() != "test0" {
		t.Error("List returned unexpected value")
	}
	if ks[1].Desc() != "test1" {
		t.Error("List returned unexpected value")
	}
	if ks[2].Desc() != "test2" {
		t.Error("List returned unexpected value")
	}
}

func TestListError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatal("initialize succeeded but this is not expected")
	}
	ks, err := List()
	if ks != nil {
		t.Error("list of kizami is not nil but this is not expected")
	}
	if err == nil {
		t.Error("err is nil but this is not expected")
	}

	// openDB goes success but List goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockList = func() ([]*kizami, error) {
		return nil, errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Error("Initialize failed")
	}
	ks, err = List()
	if ks != nil {
		t.Error("list of kizami is not nil but this is not expected")
	}
	if err == nil {
		t.Error("err is nil but this is not expected")
	}
}

func TestStopNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Error("failed to initialize")
	}
	err = Stop(0)
	if err != nil {
		t.Error("Stop returned error")
	}
}

func TestStopError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatal("initialize succeeded but this is not expected")
	}
	err = Stop(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}

	// openDB goes success but stop goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockStop = func(id int) error {
		return errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Error("Initialize failed")
	}
	err = Stop(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}
}

func TestDeleteNormal(t *testing.T) {
	dbmock := genDefaultDBMock()

	err := initialize(dbmock, "")
	if err != nil {
		t.Fatal("failed to initialize")
	}
	err = Delete(0)
	if err != nil {
		t.Error("Delete returned error")
	}
}

func TestDeleteError(t *testing.T) {
	dbmock := genDefaultDBMock()

	// openDB goes failure
	dbmock.mockOpenDB = func() error {
		return errors.New("error")
	}
	err := initialize(dbmock, "")
	if err == nil {
		t.Fatal("initialize succeeded but this is not expected")
	}
	err = Delete(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}

	// openDB goes success but stop goes failure
	dbmock.mockOpenDB = func() error {
		return nil
	}
	dbmock.mockDelete = func(id int) error {
		return errors.New("error")
	}
	err = initialize(dbmock, "")
	if err != nil {
		t.Error("Initialize failed")
	}
	err = Delete(0)
	if err == nil {
		t.Error("err is nil but this is not expected")
	}
}