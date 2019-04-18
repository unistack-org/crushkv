package crushmap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	Magic = uint32(0x00010000)
)

type CrushRuleOpType uint32

const (
	CrushRuleNoop CrushRuleOpType = iota
	CrushRuleTake
	CrushRuleChooseFirstN
	CrushRuleChooseIndep
	CrushRuleEmit
	CrushRuleChooseleafFirstN
	CrushRuleChooseleafIndep

	CrushRuleSetChooseTries
	CrushRuleSetChooseleafTries
	CrushRuleSetChooseLocalTries
	CrushRuleSetChooseLocalFallbackTries
	CrushRuleSetChooseleafVaryR
	CrushRuleSetChooseleafStable
)

var (
	crushRuleOpTypeStringMap = map[CrushRuleOpType]string{
		CrushRuleNoop:                        "noop",
		CrushRuleTake:                        "take",
		CrushRuleChooseFirstN:                "choose firstn",
		CrushRuleChooseIndep:                 "choose indep",
		CrushRuleEmit:                        "emit",
		CrushRuleChooseleafFirstN:            "choose_leaf firstn",
		CrushRuleChooseleafIndep:             "choose_leaf indep",
		CrushRuleSetChooseTries:              "set_choose_tries",
		CrushRuleSetChooseleafTries:          "set_chooseleaf_tries",
		CrushRuleSetChooseLocalTries:         "set_choose_local_tries",
		CrushRuleSetChooseLocalFallbackTries: "set_choose_local_fallback_tries",
		CrushRuleSetChooseleafVaryR:          "set_choose_leaf_vary_r",
		CrushRuleSetChooseleafStable:         "set_choose_leaf_stable",
	}
	crushRuleOpStringTypeMap = map[string]CrushRuleOpType{
		"noop":                            CrushRuleNoop,
		"take":                            CrushRuleTake,
		"choose firstn":                   CrushRuleChooseFirstN,
		"choose indep":                    CrushRuleChooseIndep,
		"emit":                            CrushRuleEmit,
		"choose_leaf firstn":              CrushRuleChooseleafFirstN,
		"choose_leaf indep":               CrushRuleChooseleafIndep,
		"set choose_tries":                CrushRuleSetChooseTries,
		"set chooseleaf_tries":            CrushRuleSetChooseleafTries,
		"set choose_local_tries":          CrushRuleSetChooseLocalTries,
		"set choose_local_fallback_tries": CrushRuleSetChooseLocalFallbackTries,
		"set choose_leaf_vary_r":          CrushRuleSetChooseleafVaryR,
		"set choose_leaf_stable":          CrushRuleSetChooseleafStable,
	}
)

func (t CrushRuleOpType) String() string {
	op, ok := crushRuleOpTypeStringMap[t]
	if !ok {
		op = "invalid"
	}
	return op
}

type CrushRuleStep struct {
	Op   CrushRuleOpType
	Arg1 int32
	Arg2 int32
}

type binaryParser struct {
	r io.Reader
	w io.Writer
}

func (cmap *Map) DecodeBinary(data []byte) error {
	var err error

	var magic uint32
	p := &binaryParser{r: bytes.NewBuffer(data)}

	err = binary.Read(p.r, binary.LittleEndian, &magic)
	if err != nil {
		return err
	} else if magic != Magic {
		return fmt.Errorf("invalid magic: %0x != %0x", magic, Magic)
	}

	var (
		maxBuckets int32
		maxRules   uint32
		maxDevices int32
	)

	err = binary.Read(p.r, binary.LittleEndian, &maxBuckets)
	if err != nil {
		return err
	}
	err = binary.Read(p.r, binary.LittleEndian, &maxRules)
	if err != nil {
		return err
	}
	err = binary.Read(p.r, binary.LittleEndian, &maxDevices)
	if err != nil {
		return err
	}

	for i := int32(0); i < maxBuckets; i++ {
		ibucket, err := p.handleBucket()
		if err != nil {
			return err
		}
		if ibucket == nil {
			continue
		}
		cmap.Buckets = append(cmap.Buckets, ibucket)
	}

	for i := uint32(0); i < maxRules; i++ {
		irule, err := p.handleRule()
		if err != nil {
			return err
		}
		cmap.Rules = append(cmap.Rules, irule)
	}

	itypes, err := p.handleType()
	if err != nil {
		return err
	}
	cmap.Types = itypes

	btypes := make(map[int32]string, len(itypes))
	for _, t := range itypes {
		btypes[t.ID] = t.Name
	}

	bnames := make(map[int32]string)
	itypes, err = p.handleType()
	if err != nil {
		return err
	}
	for _, t := range itypes {
		bnames[t.ID] = t.Name
	}

	rnames := make(map[int32]string)
	itypes, err = p.handleType()
	if err != nil {
		return err
	}
	for _, t := range itypes {
		rnames[t.ID] = t.Name
	}

	var ok bool
	for _, bucket := range cmap.Buckets {
		if bucket != nil {
			if bucket.TypeName, ok = btypes[int32(bucket.TypeID)]; !ok {
				return fmt.Errorf("unknown type id: %d", bucket.TypeID)
			}
			for _, item := range bucket.Items {
				if item.Name, ok = bnames[int32(bucket.ID)]; !ok {
					return fmt.Errorf("unknown type id: %d", bucket.ID)
				}
			}
		}
	}

	itypes, err = p.handleType()
	if err != nil {
		return err
	}

	for _, rule := range cmap.Rules {
		if rule.Name, ok = rnames[int32(rule.Ruleset)]; !ok {
			return fmt.Errorf("unknown type id: %d", rule.ID)
		}
		for _, step := range rule.Steps {
			switch step.Op {
			default:
			case CrushRuleChooseFirstN.String(), CrushRuleChooseIndep.String(), CrushRuleChooseleafFirstN.String(), CrushRuleChooseleafIndep.String():
				if step.ItemType, ok = btypes[step.ItemTypeID]; !ok {
					return fmt.Errorf("unknown type id: %d", step.ItemTypeID)
				}
			}
		}
	}

	itunables, err := p.handleTunable()
	if err != nil {
		return err
	}
	cmap.Tunables = itunables

	cmap.rulesSort()
	cmap.bucketsSort()
	return nil
}
