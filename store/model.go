package store

type chunkRow struct {
	ID       int64
	DocName  string
	Text     string
	Page     int
	Position int
}
