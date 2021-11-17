// +build !sqlite

package devicedb

// ###########################################################################
// ###########################################################################
// Database
// ###########################################################################
// ###########################################################################

// SqliteConnection ...
type SqliteConnection struct {
	SQLConnection
	path string
}

// ###########################################################################

// NewDatabaseSqlite ...
func NewDatabaseSqlite(connectString string, tableName string, clientName string) (*SqliteConnection, bool) {
	panic("Error: SQLite not compiled in this version!")
}
