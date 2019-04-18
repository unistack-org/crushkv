package crushmap

import (
	"encoding/binary"
	"io"
)

type tunables struct {
	// new block
	ChooseLocalTries         uint32 `json:"choose_local_tries,omitempty"`
	ChooseLocalFallbackTries uint32 `json:"choose_local_fallback_tries,omitempty"`
	ChooseTotalTries         uint32 `json:"choose_total_tries,omitempty"`
	// new block must be equal 1
	ChooseleafDescendOnce uint32 `json:"chooseleaf_descend_once,omitempty"`
	// new block must be equal 1
	ChooseleafVaryR uint8 `json:"chooseleaf_vary_r,omitempty"`
	// new block must be equal 1
	StrawCalcVersion uint8 `json:"straw_calc_version,omitempty"`
	// new block must be equal ??
	AllowedBucketAlgs uint32 `json:"allowed_bucket_algs,omitempty"`
	// new block must be equal 1
	ChooseleafStable uint8 `json:"chooseleaf_stable,omitempty"`
}

func legacyTunable() tunables {
	return tunables{
		ChooseLocalTries:         2,
		ChooseLocalFallbackTries: 5,
		ChooseTotalTries:         19,
		ChooseleafDescendOnce:    0,
		ChooseleafVaryR:          0,
		ChooseleafStable:         0,
		AllowedBucketAlgs:        CrushLegacyAllowedBucketAlgs,
		StrawCalcVersion:         0,
	}
}

func (p *binaryParser) handleTunable() (map[string]interface{}, error) {
	var err error

	itunables := make(map[string]interface{})
	tune := legacyTunable()

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseLocalTries)
	if err != nil {
		return nil, err
	}
	itunables["choose_local_tries"] = tune.ChooseLocalTries

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseTotalTries)
	if err != nil {
		return nil, err
	}
	itunables["choose_total_tries"] = tune.ChooseTotalTries

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseLocalFallbackTries)
	if err != nil {
		return nil, err
	}
	itunables["choose_local_fallback_tries"] = tune.ChooseLocalFallbackTries

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseleafDescendOnce)
	if err != nil && err == io.EOF {
		return itunables, nil
	} else if err != nil {
		return nil, err
	}
	itunables["chooseleaf_descend_once"] = tune.ChooseleafDescendOnce

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseleafVaryR)
	if err != nil && err == io.EOF {
		return itunables, nil
	} else if err != nil {
		return nil, err
	}
	itunables["chooseleaf_vary_r"] = tune.ChooseleafVaryR

	err = binary.Read(p.r, binary.LittleEndian, &tune.StrawCalcVersion)
	if err != nil && err == io.EOF {
		return itunables, nil
	} else if err != nil {
		return nil, err
	}
	itunables["straw_calc_version"] = tune.StrawCalcVersion

	err = binary.Read(p.r, binary.LittleEndian, &tune.AllowedBucketAlgs)
	if err != nil && err == io.EOF {
		return itunables, nil
	} else if err != nil {
		return nil, err
	}
	itunables["allowed_bucket_algs"] = tune.AllowedBucketAlgs

	err = binary.Read(p.r, binary.LittleEndian, &tune.ChooseleafStable)
	if err != nil && err == io.EOF {
		return itunables, nil
	} else if err != nil {
		return nil, err
	}
	itunables["chooseleaf_stable"] = tune.ChooseleafStable

	return itunables, nil
}
