package storage

import (
    "fmt"
    "database/sql"

    _ "github.com/mattn/go-sqlite3"

    "github.com/jls83/crisgo/types"
)

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
    selectStr := fmt.Sprintf("SELECT url FROM %s WHERE %s.shorten_key = \"%s\" LIMIT 1", s.tableName, s.tableName, k)
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
