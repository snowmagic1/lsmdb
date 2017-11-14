package lsmdb

const (
	recComparer    = 1
	recJournalNum  = 2
	recNextFileNum = 3
	recSeqNum      = 4
	recCompPtr     = 5
	recDelTable    = 6
	recAddTable    = 7
	// 8 was used for large value refs
	recPrevJournalNum = 9
)

type atRecord struct {
	level int
	num   int64
	size  int64
	imin  internalKey
	imax  internalKey
}

type sessionRecord struct {
	hasRec     int
	addedTable []atRecord
}

func (sr *sessionRecord) addTableFile(level int, t *tFile) {
	sr.addTable(level, t.fd.Num, t.size, t.imin, t.imax)
}

func (sr *sessionRecord) addTable(level int, num, size int64, imin, imax internalKey) {
	sr.hasRec |= 1 << recAddTable
	r := atRecord{level, num, size, imin, imax}
	sr.addedTable = append(sr.addedTable, r)
}
