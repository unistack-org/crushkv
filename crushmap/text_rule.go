package crushmap

import (
	"errors"
	"fmt"
	"strconv"
)

func ruleState(l *lex) stateFn {
	l.lexEmit(itemRuleBeg)

	l.lexTake(" \t")
	l.lexIgnore()

	r := l.lexPeek()
	switch r {
	case '{':
		l.lexNext()
		l.lexIgnore()
		return ruleIdentState
	case '#':
		l.lexErr(fmt.Sprintf("unexpected token %q", r))
		return l.lexPop()
	}

	return ruleNameState
}

func ruleIdentState(l *lex) stateFn {
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
			l.lexPush(ruleIdentState)
			return commentLineState
		case '}':
			l.lexNext()
			l.lexIgnore()
			l.lexEmit(itemRuleEnd)
			return l.lexPop()
		case '\n':
			l.lexNext()
			l.lexIgnore()
			return ruleIdentState
		default:
			l.lexNext()
		}
	}

	switch l.lexCurrent() {
	case "id":
		l.lexIgnore()
		return ruleRuleIDState
	case "ruleset":
		l.lexIgnore()
		return ruleRulesetState
	case "min_size":
		l.lexIgnore()
		return ruleMinSizeState
	case "max_size":
		l.lexIgnore()
		return ruleMaxSizeState
	case "type":
		l.lexIgnore()
		return ruleTypeState
	case "step":
		l.lexIgnore()
		l.lexEmit(itemRuleStepBeg)
		return ruleStepIdentState
	}
	return l.lexPop()
}

func ruleStepIdentState(l *lex) stateFn {
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

	switch l.lexCurrent() {
	case "set_chooseleaf_tries":
		l.lexIgnore()
		return ruleStepSetChooseleafTries
	case "set_choose_tries":
		l.lexIgnore()
		return ruleStepSetChooseTries
	case "take":
		l.lexIgnore()
		l.lexEmit(itemRuleStepTake)
		return ruleStepTake
	case "chooseleaf", "choose":
		l.lexEmit(itemRuleStepChoose)
		return ruleStepChoose
	case "emit":
		l.lexEmit(itemRuleStepEmit)
		return ruleStepEmit
	}

	return ruleIdentState
}

func ruleStepSetChooseleafTries(l *lex) stateFn {
	l.lexTake("0123456789")
	l.lexEmit(itemRuleStepSetChooseleafTries)
	return ruleIdentState
}

func ruleStepSetChooseTries(l *lex) stateFn {
	l.lexTake("0123456789")
	l.lexEmit(itemRuleStepSetChooseTries)
	return ruleIdentState
}

func ruleStepChoose(l *lex) stateFn {
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

	switch l.lexCurrent() {
	case "firstn":
		l.lexIgnore()
		return ruleStepChooseFirstN
	case "indep":
		l.lexIgnore()
		return ruleStepChooseIndep
	case "type":
		l.lexIgnore()
		return ruleStepChooseType
	}

	l.lexEmit(itemRuleStepEnd)
	return ruleIdentState
}

func ruleStepChooseFirstN(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	l.lexTake("0123456789")
	l.lexEmit(itemRuleStepChooseFirstN)
	return ruleStepChoose
}

func ruleStepChooseIndep(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	l.lexTake("0123456789")
	l.lexEmit(itemRuleStepChooseIndep)
	return ruleStepChoose
}

func ruleStepChooseType(l *lex) stateFn {
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

	l.lexEmit(itemRuleStepChooseType)
	return ruleStepChoose
}

func ruleStepEmit(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	l.lexEmit(itemRuleStepEnd)
	return ruleIdentState
}

func ruleStepTake(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

loop1:
	for {
		r := l.lexPeek()
		switch r {
		case ' ', '\n', '#', '\t':
			break loop1
		default:
			l.lexNext()
		}
	}

	l.lexEmit(itemRuleStepTakeType)

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
		return ruleStepTakeClass
	}

	l.lexEmit(itemRuleStepEnd)
	return ruleIdentState
}

func ruleStepTakeClass(l *lex) stateFn {
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

	l.lexEmit(itemRuleStepTakeClass)
	l.lexEmit(itemRuleStepEnd)
	return ruleIdentState
}

func ruleNameState(l *lex) stateFn {
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
	l.lexEmit(itemRuleName)
	return ruleIdentState
}

func ruleRulesetState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	for {
		r := l.lexPeek()
		if r == '\n' || r == ' ' || r == '#' {
			break
		}
		l.lexNext()
	}

	l.lexEmit(itemRuleRuleset)
	return ruleIdentState
}

func ruleRuleIDState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	for {
		r := l.lexPeek()
		if r == '\n' || r == ' ' || r == '#' {
			break
		}
		l.lexNext()
	}

	l.lexEmit(itemRuleID)
	return ruleIdentState
}

func ruleMinSizeState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	l.lexTake("0123456789")
	l.lexEmit(itemRuleMinSize)
	return ruleIdentState
}

func ruleMaxSizeState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()
	l.lexTake("0123456789")
	l.lexEmit(itemRuleMaxSize)
	return ruleIdentState
}

func ruleTypeState(l *lex) stateFn {
	l.lexTake(" \t")
	l.lexIgnore()

	for {
		r := l.lexPeek()
		if r == '\n' || r == ' ' || r == '#' {
			break
		}
		l.lexNext()
	}
	l.lexEmit(itemRuleType)
	return ruleIdentState
}

func (p *textParser) handleRule() (*Rule, error) {
	irule := &Rule{}

Loop:
	for {
		tok, done := p.l.lexNextToken()
		if done {
			break Loop
		}
		switch tok.itype {
		case itemEOF, itemRuleEnd:
			break Loop
		case itemComment:
			continue
		case itemRuleRuleset:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			irule.Ruleset = uint8(id)
		case itemRuleMinSize:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			irule.MinSize = uint8(id)
		case itemRuleMaxSize:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			irule.MaxSize = uint8(id)
		case itemRuleStepBeg:
			if istep, err := p.handleRuleStep(); err != nil {
				return nil, err
			} else {
				istep.Num = int32(len(irule.Steps))
				irule.Steps = append(irule.Steps, istep)
			}
		case itemRuleName:
			irule.Name = tok.ivalue
		case itemRuleType:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				switch tok.ivalue {
				case "replicated":
					irule.Type = ReplicatedPG
				case "erasure":
					irule.Type = ErasurePG
				default:
					return nil, errors.New("unknown rule type")
				}
			} else {
				irule.Type = uint8(id)
			}
		case itemRuleID:
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			irule.ID = uint8(id)
		}
	}
	if irule.ID != irule.Ruleset {
		irule.Ruleset = irule.ID
	}
	return irule, nil
}

func (p *textParser) handleRuleStep() (*Step, error) {
	istep := &Step{}

Loop:
	for {
		tok, done := p.l.lexNextToken()
		if done {
			break Loop
		}
		switch tok.itype {
		case itemEOF, itemRuleStepEnd:
			break Loop
		case itemComment:
			continue
		case itemRuleStepTake:
			istep.Op = "take"
			istep.Item = -1
		case itemRuleStepTakeType:
			istep.ItemName = tok.ivalue
		case itemRuleStepTakeClass:
			istep.ItemClass = tok.ivalue
		case itemRuleStepChoose:
			istep.Op = tok.ivalue
		case itemRuleStepChooseIndep:
			istep.Op = fmt.Sprintf("%s_%s", istep.Op, "indep")
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			istep.Num = int32(id)
		case itemRuleStepChooseFirstN:
			istep.Op = fmt.Sprintf("%s_%s", istep.Op, "firstn")
			id, err := strconv.Atoi(tok.ivalue)
			if err != nil {
				return nil, err
			}
			istep.Num = int32(id)
		case itemRuleStepChooseType:
			istep.ItemType = tok.ivalue
		case itemRuleStepEmit:
			istep.Op = "emit"
		}
	}
	return istep, nil
}
