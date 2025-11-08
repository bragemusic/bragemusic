package database

import (
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

func normalizeForCompare(s string) string {
	var b strings.Builder
	for _, r := range s {
		// Keep only letters, digits
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return strings.TrimSpace(b.String())
}

func upgradeConn(conn *sqlx.Conn) error {
	err := conn.Conn.Raw(func(driverConn any) error {
		sqliteConn := driverConn.(*sqlite3.SQLiteConn)
		// Register our custom normalize() SQL function
		return sqliteConn.RegisterFunc("normalize", normalizeForCompare, true)
	})
	return err
}
