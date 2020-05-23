package storage

import (
    "fmt"
    "time"
    "database/sql"
    "math/rand"

    _ "github.com/mattn/go-sqlite3"

    "github.com/jls83/crisgo/types"
)


type ResultStorage interface {
    Close() (err error)
    GetResultMapKey() types.ResKey
    GetValue(k types.ResKey) (types.ResValue, bool)
    InsertValue(v types.ResValue) types.ResKey
    // TODO: Add explicit `SetValue` method & endpoint for testing
}

// Section: Helper Methods
// This was cribbed from https://www.calhoun.io/creating-random-strings-in-go/
const charset = "abcdefghijklmnopqrstuvwxyz" +
                "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
                "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
    byteArray := make([]byte, length)
    for i:= range byteArray {
        byteArray[i] = charset[seededRand.Intn(len(charset))]
    }

    return string(byteArray)
}


// Section: LocalStorage
type LocalStorage struct {
    _innerStorage types.ResMap
}

func NewLocalStorage() *LocalStorage {
    localStorage := types.ResMap{}
    return &LocalStorage{localStorage}
}

func (s LocalStorage) Close() (err error) {
    // Since this is just in-memory, don't actually do anything
    return
}

func (s LocalStorage) GetResultMapKey() types.ResKey {
    return types.ResKey(StringWithCharset(5, charset))
}

func (s LocalStorage) GetValue(k types.ResKey) (types.ResValue, bool) {
    value, found := s._innerStorage[k]
    return value, found
}

func (s *LocalStorage) InsertValue(v types.ResValue) types.ResKey {
    // TODO: Add some error handling; I bet shit can get weird
    var resultKey types.ResKey
    hasKey := true

    // Loop until we have a key not already in the map
    for hasKey {
        resultKey = s.GetResultMapKey()
        _, hasKey = s._innerStorage[resultKey]
        fmt.Println(hasKey)
    }

    s._innerStorage[resultKey] = v

    return resultKey
}

// Section: SqliteStorage
const SQLITE_FILE_PATH = "crisgo.db"
const SQLITE_TABLE_NAME = "shortened_urls"

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

type SqliteStorage struct {
    _db *sql.DB
    databasePath string
    tableName string
}

func NewSqliteStorage(databasePath string, tableName string) *SqliteStorage {
    db, err := sql.Open("sqlite3", databasePath)
    checkErr(err)

    sqliteStorage := SqliteStorage{db, databasePath, tableName}

    return &sqliteStorage
}

func (s *SqliteStorage) CreateTable() error {
    createStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`" +
                             "(`id` INTEGER PRIMARY KEY AUTOINCREMENT," +
                             "`url` VARCHAR(1200)," +
                             "`shorten_key` VARCHAR(1200));", s.tableName)

    createStatement, err := s._db.Prepare(createStr)
    checkErr(err)

    _, err = createStatement.Exec()

    return err
}

func (s *SqliteStorage) Close() (err error) {
    return s._db.Close()
}

func (s *SqliteStorage) GetResultMapKey() types.ResKey {
    return types.ResKey(StringWithCharset(5, charset))
}

func (s *SqliteStorage) GetValue(k types.ResKey) (types.ResValue, bool) {
    selectStr := fmt.Sprintf("SELECT url FROM %s WHERE %s.shorten_key = %s LIMIT 1", s.tableName, s.tableName, k)
    selectRows, err := s._db.Query(selectStr)

    if err != nil {
        return types.ResValue(""), false
    }

    var resultValue types.ResValue

    for selectRows.Next() {
        err = selectRows.Scan(&resultValue)
        if err != nil {
            return types.ResValue(""), false
        }
    }

    return resultValue, true
}

func (s *SqliteStorage) InsertValue(v types.ResValue) types.ResKey {
    var resultKey types.ResKey
    hasKey := true

    // Loop until we have a key not already in the db
    for hasKey {
        resultKey = s.GetResultMapKey()
        selectStr := fmt.Sprintf("SELECT COUNT(*) FROM %s " +
                                 "WHERE %s.shorten_key = \"%s\" LIMIT 1;", s.tableName, s.tableName, resultKey)
        selectRows, err := s._db.Query(selectStr)
        checkErr(err)

        var resultKeyCount int
        for selectRows.Next() {
            // err = selectRows.Scan(&resultKeyCount)
            // checkErr(err)
            selectRows.Scan(&resultKeyCount)
        }
        if (resultKeyCount < 1) {
            hasKey = false
        }
    }

    insertStr := fmt.Sprintf("INSERT INTO %s (url, shorten_key)" +
                             "VALUES (?, ?)", s.tableName)

    insertStatement, err := s._db.Prepare(insertStr)
    checkErr(err)

    res, err := insertStatement.Exec(v, resultKey)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    fmt.Printf("Inserted %s; %d Rows affected", string(v), affect)

    return resultKey
}
