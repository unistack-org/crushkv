package crushmap

import (
	"errors"
	"fmt"
	"strconv"
)

func bucketState(l *lex) stateFn {
	l.lexEmit(itemBucketBeg)
	return bucketStartState
}

func bucketStartState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	r := l.lexPeek()
	switch r {
	case '{':
		l.lexNext()
		l.lexIgnore()
		return bucketIdentState
	case '#':
		l.lexErr(fmt.Sprintf("unexpected token %q", r))
		return l.lexPop()
	}

	return bucketNameState
}

func bucketNameState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop:
	for {
		r := l.lexPeek()
		switch r {
		case '{':
			l.lexNext()
			l.lexIgnore()
			break loop
		case ' ':
			break loop
		case '\n', '#':
			l.lexErr(fmt.Sprintf("unexpected token %q", r))
			return l.lexPop()
		default:
			l.lexNext()
		}
	}

	l.lexEmit(itemBucketName)
	return bucketIdentState
}

func bucketIDState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop1:
	for {
		r := l.lexPeek()
		switch r {
		case ' ', '\n', '#', '\t':
			break loop1
		}
		l.lexNext()
	}

	l.lexEmit(itemBucketID)

	l.lexTake(" \t")
	l.lexIgnore()

loop2:
	for {
		r := l.lexPeek()
		switch r {
		case ' ', '\n', '#', '\t':
			break loop2
		default:
			l.lexNext()
		}
	}

	switch l.lexCurrent() {
	case "class":
		l.lexIgnore()
		return bucketIDClassState
	}

	return bucketIdentState
}

func bucketIDClassState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop:
	for {
		r := l.lexPeek()
		switch r {
		case ' ', '\n', '#', '\t':
			break loop
		default:
			l.lexNext()
		}
	}

	l.lexEmit(itemBucketIDClass)

	return bucketIdentState
}

func bucketHashState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	l.lexTake("0123456789")
	l.lexEmit(itemBucketHash)

	return bucketIdentState
}

func bucketAlgState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	for {
		r := l.lexPeek()
		if r == '\n' || r == ' ' || r == '#' {
			break
		}
		l.lexNext()
	}

	l.lexEmit(itemBucketAlg)
	return bucketIdentState
}

func bucketItemState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop:
	for {
		r := l.lexPeek()
		switch r {
		case ' ':
			break loop
		case '\n', '#':
			l.lexErr(fmt.Sprintf("unexpected token %q", r))
			return l.lexPop()
		}
		l.lexNext()
	}

	l.lexEmit(itemBucketItemName)
	return bucketItemIdentState
}

func bucketItemIdentState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop:
	for {
		r := l.lexPeek()
		switch r {
		case ' ', '\n':
			break loop
		case '#':
			break loop
		default:
			l.lexNext()
		}
	}

	switch l.lexCurrent() {
	case "weight":
		l.lexIgnore()
		return bucketItemWeightState
	case "pos":
		l.lexIgnore()
		return bucketItemPosState
	}

	l.lexEmit(itemBucketItemEnd)
	return bucketIdentState
}

func bucketItemWeightState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	l.lexTake(".0123456789")
	l.lexEmit(itemBucketItemWeight)

	return bucketItemIdentState
}

func bucketItemPosState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	l.lexTake("0123456789")
	l.lexEmit(itemBucketItemPos)

	return bucketItemIdentState
}

func bucketIdentState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop:
	for {
		r := l.lexPeek()
		switch r {
		case ' ':
			break loop
		case '#':
			l.lexNext()
			l.lexIgnore()
			l.lexPush(bucketIdentState)
			return commentLineState
		case '}':
			l.lexNext()
			l.lexIgnore()
			l.lexEmit(itemBucketEnd)
			return l.lexPop()
		case '\n':
			l.lexNext()
			l.lexIgnore()
			return bucketIdentState
		default:
			l.lexNext()
		}
	}

	switch l.lexCurrent() {
	case "id":
		l.lexIgnore()
		return bucketIDState
	case "alg":
		l.lexIgnore()
		return bucketAlgState
	case "hash":
		l.lexIgnore()
		return bucketHashState
	case "item":
		l.lexIgnore()
		l.lexEmit(itemBucketItemBeg)
		return bucketItemState
	}

	return l.lexPop()
}

func (p *textParser) handleBucket(itype string) (*Bucket, error) {
	ibucket := &Bucket{TypeName: itype}

Loop:
	for {
		tok, done := p.l.lexNextToken()
		if done {
			break Loop
		}
		switch tok.itype {
		case itemEOF, itemBucketEnd:
			break Loop
		case itemComment:
			continue
		case itemBucketName:
			ibucket.Name = tok.ivalue
		case itemBucketIDClass:
			ibucket.IDClass = tok.ivalue
		case itemBucketID:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			ibucket.ID = int32(id)
		case itemBucketAlg:
			ibucket.Alg = tok.ivalue
		case itemBucketHash:
			if tok.ivalue == "0" {
				ibucket.Hash = "rjenkins1"
			} else {
				return nil, errors.New("invalid bucket hash")
			}
		case itemBucketItemBeg:
			item, err := p.handleBucketItem()
			if err != nil {
				return nil, err
			}
			ibucket.Items = append(ibucket.Items, item)
		}
	}

	return ibucket, nil
}

func (p *textParser) handleBucketItem() (*Item, error) {
	item := &Item{}

Loop:
	for {
		tok, done := p.l.lexNextToken()
		if done {
			break Loop
		}
		switch tok.itype {
		case itemEOF, itemBucketItemEnd:
			break Loop
		case itemComment:
			continue
		case itemBucketItemName:
			item.Name = tok.ivalue
		case itemBucketItemWeight:
			id, err := strconv.ParseFloat(tok.ivalue, 32)
			if err != nil {
				return nil, err
			}
			item.Weight = float32(id)
		case itemBucketItemPos:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			item.Pos = id
		}
	}
	return item, nil
}
