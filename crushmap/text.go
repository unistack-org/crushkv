package crushmap

import (
	"errors"
	"fmt"
)

type textParser struct {
	l *lex
}

func identState(l *lex) stateFn {

loop:
	for {
		r := l.lexPeek()
		switch r {
		case ' ':
			break loop
		case EOFRune:
			l.lexErr("unexpected end of crush map")
			return nil
		default:
			l.lexNext()
		}
	}

	switch l.lexCurrent() {
	case "device":
		l.lexIgnore()
		l.lexPush(topState)
		return deviceState
	case "type":
		l.lexIgnore()
		l.lexPush(topState)
		return typeState
	case "rule":
		l.lexIgnore()
		l.lexPush(topState)
		return ruleState
	case "tunable":
		l.lexIgnore()
		l.lexPush(topState)
		return tunableState
	}

	l.lexPush(topState)
	return bucketState
}

func topState(l *lex) stateFn {
	for {
		r := l.lexPeek()
		switch r {
		case ' ':
			l.lexNext()
			l.lexIgnore()
		case '\n':
			l.lexNext()
			l.lexIgnore()
		case EOFRune:
			l.lexEmit(itemEOF)
			return nil
		case '#':
			l.lexNext()
			l.lexIgnore()
			l.lexPush(topState)
			return commentLineState
		default:
			return identState
		}
	}
}

func (cmap *Map) DecodeText(data []byte) error {
	mapItems := make(map[string]int32)
	p := &textParser{l: lexNew(string(data), topState)}
	p.l.lexStart()

loop:
	for {
		tok, done := p.l.lexNextToken()
		if done {
			break loop
		}
		switch tok.itype {
		case itemEOF:
			break loop
		case itemComment:
			continue
		case itemTunableBeg:
			if itunekey, ituneval, err := p.handleTunable(); err != nil {
				return err
			} else {
				cmap.Tunables[itunekey] = ituneval
			}
		case itemDeviceBeg:
			if idevice, err := p.handleDevice(); err != nil {
				return err
			} else {
				mapItems[idevice.Name] = idevice.ID
				cmap.Devices = append(cmap.Devices, idevice)
			}
		case itemTypeBeg:
			if itype, err := p.handleType(); err != nil {
				return err
			} else {
				mapItems[itype.Name] = itype.ID
				cmap.Types = append(cmap.Types, itype)
			}
		case itemRuleBeg:
			if irule, err := p.handleRule(); err != nil {
				return err
			} else {
				cmap.Rules = append(cmap.Rules, irule)
			}
		case itemBucketBeg:
			if ibucket, err := p.handleBucket(tok.ivalue); err != nil {
				return err
			} else {
				mapItems[ibucket.Name] = ibucket.ID
				cmap.Buckets = append(cmap.Buckets, ibucket)
			}
		default:
			return fmt.Errorf("error: %s\n", tok.ivalue)
		}
	}

	// The lexer stops on the first error and closes the items channel; a
	// closed channel (done) without an itemEOF token means the input ended
	// mid-construct. Reading l.err after the channel is closed is safe.
	if p.l.err != nil {
		return p.l.err
	}

	for idx := range cmap.Buckets {
		id, ok := mapItems[cmap.Buckets[idx].TypeName]
		if !ok {
			return fmt.Errorf("invalid bucket type: %s", cmap.Buckets[idx].TypeName)
		}
		cmap.Buckets[idx].TypeID = CrushBucketType(id)

		for _, it := range cmap.Buckets[idx].Items {
			if it.Name == "" {
				continue
			}
			iid, ok := mapItems[it.Name]
			if !ok {
				return fmt.Errorf("bucket %s: unknown item %s", cmap.Buckets[idx].Name, it.Name)
			}
			it.ID = iid
		}
	}

	if len(cmap.Buckets) == 0 && len(cmap.Rules) == 0 {
		return errors.New("empty crush map")
	}

	cmap.rulesSort()
	cmap.bucketsSort()
	return nil
}
