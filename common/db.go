package common

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"os"
	"path"
)

// GetDbPath returns string with path to database
func GetDbPath(dataDir string, name string) string {
	return path.Join(dataDir, "db", name)
}

// OpenDatabase opens (creates if it does not exist) a LevelDB database with provided name
func OpenDatabase(dataDir string, name string, options *opt.Options) (*leveldb.DB, error) {
	dbPath := GetDbPath(dataDir, name)

	db, err := leveldb.OpenFile(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// DeleteDatabase deletes LevelDB database files from hdd
func DeleteDatabase(dataDir string, name string) error {
	dbPath := GetDbPath(dataDir, name)
	err := os.RemoveAll(dbPath)
	if err != nil {
		return err
	}
	return nil
}
