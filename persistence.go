package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	dbFile = filepath.Join(homeDir(), ".scalingo-addons-tester")
)

type AddonDatabase []*AddonDatabaseEntry
type AddonDatabaseEntry struct {
	ID   string `json:"id"`
	Plan string `json:"plan"`
}

func (db *AddonDatabase) HasAddon(id string) (int, *AddonDatabaseEntry, bool) {
	for i, entry := range *db {
		if entry.ID == id {
			return i, entry, true
		}
	}
	return -1, nil, false
}

func saveAddonRef(id string, plan string) error {
	return doRefOperation(func(db *AddonDatabase) error {
		if _, addon, ok := db.HasAddon(id); ok {
			addon.Plan = plan
		} else {
			*db = append(*db, &AddonDatabaseEntry{
				ID: id, Plan: plan,
			})
		}
		return nil
	})
}

func deleteAddonRef(id string) error {
	return doRefOperation(func(db *AddonDatabase) error {
		if i, _, ok := db.HasAddon(id); ok {
			*db = append((*db)[:i], (*db)[i+1:]...)
		}
		return nil
	})
}

func getDB() (*AddonDatabase, error) {
	fd, err := os.OpenFile(dbFile, os.O_RDONLY, 0644)
	var db *AddonDatabase
	if err != nil {
		if os.IsNotExist(err) {
			db = &AddonDatabase{}
		} else {
			return nil, fmt.Errorf("Fail to open database file: %v", err)
		}
	}

	if !os.IsNotExist(err) {
		err = json.NewDecoder(fd).Decode(&db)
		if err != nil {
			return nil, fmt.Errorf("Fail to decode database file, please delete %v: %v", dbFile, err)
		}
		fd.Close()
	}
	return db, nil
}

func doRefOperation(f func(db *AddonDatabase) error) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	err = f(db)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(filepath.Join(homeDir(), ".scalingo-addons-tester"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	err = json.NewEncoder(fd).Encode(&db)
	if err != nil {
		return err
	}
	fd.Close()
	return nil
}

func rmDB() error {
	return os.Remove(dbFile)
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
