package mysql

import (
	"github.com/go-sql-driver/mysql"
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

func KeyDuplicatedError(err error) bool {
	if err == nil {
		return false
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}

func MD5String(data string) string {
        return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func MD5Sum(data string) []byte {
        val := md5.Sum([]byte(data))
        return val[:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MD5High(md5sum []byte) int64 {
        return int64(binary.BigEndian.Uint64(md5sum[:8]))
}

func MD5Low(md5sum []byte) int64 {
        return int64(binary.BigEndian.Uint64(md5sum[8:]))
}
