package crushmap

import (
	"encoding/binary"
)

type binaryRuleStep struct {
	Op   CrushRuleOpType
	Arg1 int32
	Arg2 int32
}

type binaryRuleMask struct {
	Ruleset uint8
	Type    uint8
	MinSize uint8
	MaxSize uint8
}

type binaryRule struct {
	Len   uint32
	Mask  binaryRuleMask
	Steps []binaryRuleStep
}

func (p *binaryParser) handleRule() (*Rule, error) {
	var err error
	irule := &Rule{}

	var yes uint32
	err = binary.Read(p.r, binary.LittleEndian, &yes)
	if err != nil {
		return nil, err
	}
	if yes == 0 {
		return nil, nil
	}

	var rule binaryRule
	err = binary.Read(p.r, binary.LittleEndian, &rule.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Read(p.r, binary.LittleEndian, &rule.Mask)
	if err != nil {
		return nil, err
	}

	//	rule.Steps = make([]RuleStep, rule.Len)
	for i := uint32(0); i < rule.Len; i++ {
		var step binaryRuleStep
		istep := &Step{}
		err = binary.Read(p.r, binary.LittleEndian, &step)
		if err != nil {
			return nil, err
		}
		istep.Op = step.Op.String()
		switch step.Op {
		case CrushRuleChooseFirstN, CrushRuleChooseIndep, CrushRuleChooseleafFirstN, CrushRuleChooseleafIndep:
			istep.Num = step.Arg1
			istep.ItemTypeID = step.Arg2 //TYPE!!!
		case CrushRuleTake:
			istep.ItemTypeID = step.Arg1
			//	case CrushRuleEmit:

			//	default:
			//		panic(step.Op.String())
		}
		/*
		   Op        string `json:"op"`
		   Item      int    `json:"item,omitempty"`
		   ItemName  string `json:"item_name,omitempty"`
		   ItemClass string `json:"item_class,omitempty"`
		   Num       int    `json:"num,omitempty"`
		   ItemType  string `json:"type,omitempty"`
		*/
		irule.Steps = append(irule.Steps, istep)
	}
	irule.ID = rule.Mask.Ruleset
	irule.Ruleset = rule.Mask.Ruleset
	irule.MinSize = rule.Mask.MinSize
	irule.MaxSize = rule.Mask.MaxSize
	irule.Type = rule.Mask.Type
	return irule, nil
}
