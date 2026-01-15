package models

import (
    "database/sql"
    "errors"
    "sync"
)

type ATK struct {
    ID        int    `json:"id"`
    Nama      string `json:"nama"`
    Jenis     string `json:"jenis"`
    Qty       int    `json:"qty"`
    IsDeleted int    `json:"id_deleted"` // 0 = aktif, 1 = soft deleted
}

var mu sync.Mutex

// DB is set by InitDB; if nil, the in-memory store is used.
var DB *sql.DB

var DatabaseATK = []ATK{
    {ID: 1, Nama: "Pulpen Gel", Jenis: "Alat Tulis", Qty: 10, IsDeleted: 0},
}

func InitDB(db *sql.DB) error {
    if db == nil {
        return errors.New("nil db")
    }
    DB = db
    
    // create table if not exists with id_deleted column for soft delete
    _, err := DB.Exec(`CREATE TABLE IF NOT EXISTS mst_atk (
        id INT AUTO_INCREMENT PRIMARY KEY,
        nama VARCHAR(255),
        jenis VARCHAR(255),
        qty INT,
        id_deleted INT DEFAULT 0
    )`)
    if err != nil {
        return err
    }
    
    // Migration: Add id_deleted column if it doesn't exist (for existing tables)
    _, err = DB.Exec(`ALTER TABLE mst_atk ADD COLUMN IF NOT EXISTS id_deleted INT DEFAULT 0`)
    // Ignore error if column already exists (MySQL doesn't support IF NOT EXISTS for columns in all versions)
    // For MySQL 5.x compatibility:
    DB.Exec(`ALTER TABLE mst_atk ADD COLUMN id_deleted INT DEFAULT 0`)
    
    return nil
}

func Get() ([]ATK, error) {
    if DB != nil {
        // Only get non-deleted items (id_deleted = 0)
        rows, err := DB.Query("SELECT id, nama, jenis, qty, id_deleted FROM mst_atk WHERE id_deleted = 0")
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        var out []ATK
        for rows.Next() {
            var a ATK
            if err := rows.Scan(&a.ID, &a.Nama, &a.Jenis, &a.Qty, &a.IsDeleted); err != nil {
                return nil, err
            }
            out = append(out, a)
        }
        return out, nil
    }

    mu.Lock()
    defer mu.Unlock()
    // Filter non-deleted items in memory store
    var out []ATK
    for _, atk := range DatabaseATK {
        if atk.IsDeleted == 0 {
            out = append(out, atk)
        }
    }
    return out, nil
}

func Post(atk ATK) (ATK, error) {
    if DB != nil {
        // Default id_deleted to 0 for new items
        res, err := DB.Exec("INSERT INTO mst_atk (nama, jenis, qty, id_deleted) VALUES (?, ?, ?, 0)", atk.Nama, atk.Jenis, atk.Qty)
        if err != nil {
            return ATK{}, err
        }
        id64, err := res.LastInsertId()
        if err == nil {
            atk.ID = int(id64)
        }
        atk.IsDeleted = 0
        return atk, nil
    }

    mu.Lock()
    defer mu.Unlock()
    atk.ID = 1
    if len(DatabaseATK) > 0 {
        atk.ID = DatabaseATK[len(DatabaseATK)-1].ID + 1
    }
    atk.IsDeleted = 0
    DatabaseATK = append(DatabaseATK, atk)
    return atk, nil
}

func Update(id int, updated ATK) (ATK, bool, error) {
    if DB != nil {
        // Only update non-deleted items
        res, err := DB.Exec("UPDATE mst_atk SET nama=?, jenis=?, qty=? WHERE id=? AND id_deleted=0", updated.Nama, updated.Jenis, updated.Qty, id)
        if err != nil {
            return ATK{}, false, err
        }
        n, _ := res.RowsAffected()
        if n == 0 {
            return ATK{}, false, nil
        }
        updated.ID = id
        updated.IsDeleted = 0
        return updated, true, nil
    }

    mu.Lock()
    defer mu.Unlock()
    for i := range DatabaseATK {
        if DatabaseATK[i].ID == id && DatabaseATK[i].IsDeleted == 0 {
            updated.ID = id
            updated.IsDeleted = 0
            DatabaseATK[i] = updated
            return updated, true, nil
        }
    }
    return ATK{}, false, nil
}

func Delete(id int) (bool, error) {
    if DB != nil {
        // Soft delete: set id_deleted = 1 instead of actual DELETE
        res, err := DB.Exec("UPDATE mst_atk SET id_deleted=1 WHERE id=? AND id_deleted=0", id)
        if err != nil {
            return false, err
        }
        n, _ := res.RowsAffected()
        return n > 0, nil
    }

    mu.Lock()
    defer mu.Unlock()
    for i := range DatabaseATK {
        if DatabaseATK[i].ID == id && DatabaseATK[i].IsDeleted == 0 {
            // Soft delete in memory
            DatabaseATK[i].IsDeleted = 1
            return true, nil
        }
    }
    return false, nil
}

// Restore: mengembalikan data yang sudah di-soft delete
func Restore(id int) (bool, error) {
    if DB != nil {
        res, err := DB.Exec("UPDATE mst_atk SET id_deleted=0 WHERE id=? AND id_deleted=1", id)
        if err != nil {
            return false, err
        }
        n, _ := res.RowsAffected()
        return n > 0, nil
    }

    mu.Lock()
    defer mu.Unlock()
    for i := range DatabaseATK {
        if DatabaseATK[i].ID == id && DatabaseATK[i].IsDeleted == 1 {
            DatabaseATK[i].IsDeleted = 0
            return true, nil
        }
    }
    return false, nil
}

// GetDeleted: ambil data yang sudah di-soft delete
func GetDeleted() ([]ATK, error) {
    if DB != nil {
        rows, err := DB.Query("SELECT id, nama, jenis, qty, id_deleted FROM mst_atk WHERE id_deleted = 1")
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        var out []ATK
        for rows.Next() {
            var a ATK
            if err := rows.Scan(&a.ID, &a.Nama, &a.Jenis, &a.Qty, &a.IsDeleted); err != nil {
                return nil, err
            }
            out = append(out, a)
        }
        return out, nil
    }

    mu.Lock()
    defer mu.Unlock()
    var out []ATK
    for _, atk := range DatabaseATK {
        if atk.IsDeleted == 1 {
            out = append(out, atk)
        }
    }
    return out, nil
}

// HardDelete: hapus permanen dari database (gunakan dengan hati-hati!)
func HardDelete(id int) (bool, error) {
    if DB != nil {
        res, err := DB.Exec("DELETE FROM mst_atk WHERE id=?", id)
        if err != nil {
            return false, err
        }
        n, _ := res.RowsAffected()
        return n > 0, nil
    }

    mu.Lock()
    defer mu.Unlock()
    for i := range DatabaseATK {
        if DatabaseATK[i].ID == id {
            DatabaseATK = append(DatabaseATK[:i], DatabaseATK[i+1:]...)
            return true, nil
        }
    }
    return false, nil
}
