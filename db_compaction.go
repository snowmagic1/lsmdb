package lsmdb

type cCmd interface {
	ack(err error)
}

type cAuto struct {
	ackC chan<- error
}

func (r cAuto) ack(err error) {
	if r.ackC != nil {
		r.ackC <- err
	}
}

func (db *DB) compTrigger(compC chan<- cCmd) {
	select {
	case compC <- cAuto{}:
	default:
	}
}

func (db *DB) compTriggerWait(compC chan<- cCmd) (err error) {
	ch := make(chan error)
	defer close(ch)

	//send
	select {
	case compC <- cAuto{ch}:
	case <-db.closeC:
		return ErrClosed
	}

	// wait
	select {
	case err = <-ch:
	case <-db.closeC:
		return ErrClosed
	}

	return err
}

func (db *DB) mCompaction() {
	var x cCmd

	for {
		select {
		case x = <-db.mcompCmdC:
			switch x.(type) {
			case cAuto:
				db.memCompaction()
				x.ack(nil)
				x = nil
			default:
				panic("unknown command")
			}
		case <-db.closeC:
			return
		}
	}
}

func (db *DB) memCompaction() {
	mdb := db.getFrozenMem()
	if mdb == nil {
		return
	}
	defer mdb.decref()

	if mdb.Len() == 0 {
		db.dropFrozenMem()
		return
	}

	// pause table compaction

	var (
		rec        = &sessionRecord{}
		flushLevel int
	)

	// generate table
	flushLevel, err := db.s.flushMemdb(rec, mdb.DB)

	if err != nil || flushLevel == 0 {

	}

	return
}
