package crushmap

import (
	"encoding/binary"
	"fmt"
	"math"
)

type binaryBucket interface {
	BucketID() int32
	BucketType() CrushBucketType
	BucketAlg() CrushAlgType
	BucketHash() CrushBucketHashType
	BucketWeight() float32
	BucketSize() uint32
}

func (b *BucketUniform) BucketID() int32 {
	return b.ID
}
func (b *BucketUniform) BucketType() CrushBucketType {
	return b.Type
}
func (b *BucketUniform) BucketAlg() CrushAlgType {
	return b.Alg
}
func (b *BucketUniform) BucketHash() CrushBucketHashType {
	return b.Hash
}
func (b *BucketUniform) BucketWeight() float32 {
	return b.Weight
}
func (b *BucketUniform) BucketSize() uint32 {
	return b.Size
}
func (b *BucketList) BucketID() int32 {
	return b.ID
}
func (b *BucketList) BucketType() CrushBucketType {
	return b.Type
}
func (b *BucketList) BucketAlg() CrushAlgType {
	return b.Alg
}
func (b *BucketList) BucketHash() CrushBucketHashType {
	return b.Hash
}
func (b *BucketList) BucketWeight() float32 {
	return b.Weight
}
func (b *BucketList) BucketSize() uint32 {
	return b.Size
}
func (b *BucketTree) BucketID() int32 {
	return b.ID
}
func (b *BucketTree) BucketType() CrushBucketType {
	return b.Type
}
func (b *BucketTree) BucketAlg() CrushAlgType {
	return b.Alg
}
func (b *BucketTree) BucketHash() CrushBucketHashType {
	return b.Hash
}
func (b *BucketTree) BucketWeight() float32 {
	return b.Weight
}
func (b *BucketTree) BucketSize() uint32 {
	return b.Size
}
func (b *BucketStraw) BucketID() int32 {
	return b.ID
}
func (b *BucketStraw) BucketType() CrushBucketType {
	return b.Type
}
func (b *BucketStraw) BucketAlg() CrushAlgType {
	return b.Alg
}
func (b *BucketStraw) BucketHash() CrushBucketHashType {
	return b.Hash
}
func (b *BucketStraw) BucketWeight() float32 {
	return b.Weight
}
func (b *BucketStraw) BucketSize() uint32 {
	return b.Size
}

func (b *BucketStraw2) BucketID() int32 {
	return b.ID
}
func (b *BucketStraw2) BucketType() CrushBucketType {
	return b.Type
}
func (b *BucketStraw2) BucketAlg() CrushAlgType {
	return b.Alg
}
func (b *BucketStraw2) BucketHash() CrushBucketHashType {
	return b.Hash
}
func (b *BucketStraw2) BucketWeight() float32 {
	return b.Weight
}
func (b *BucketStraw2) BucketSize() uint32 {
	return b.Size
}

type binaryBucketHeader struct {
	ID     int32
	Type   CrushBucketType
	Alg    CrushAlgType
	Hash   CrushBucketHashType
	Weight float32
	Size   uint32
}

type binaryBucketCommon struct {
	binaryBucketHeader
	Items []int32
}

type BucketUniform struct {
	binaryBucketCommon
	ItemWeight float32
}

type BucketList struct {
	binaryBucketCommon
	ItemWeights []float32
	SumWeights  []float32
}

type BucketTree struct {
	binaryBucketCommon
	NumNodes    uint8
	NodeWeights []float32
}

type BucketStraw struct {
	binaryBucketCommon
	ItemWeights []float32
	Straws      []uint32
}

type BucketStraw2 struct {
	binaryBucketCommon
	ItemWeights []float32
}

func (b *binaryBucketHeader) String() string {
	return fmt.Sprintf("id: %d, type: %s, alg: %s, hash: %s, weight: %f, size: %d",
		b.ID, b.Type, b.Alg, b.Hash, b.Weight, b.Size)
}

func (p *binaryParser) handleBucket() (*Bucket, error) {
	var err error
	ibucket := &Bucket{}
	var bucket binaryBucket

	var alg uint32
	err = binary.Read(p.r, binary.LittleEndian, &alg)
	if err != nil {
		return nil, err
	}

	if CrushAlgType(alg) == CrushAlgInvalid {
		return nil, nil
	}

	bucketHeader := binaryBucketHeader{}
	err = binary.Read(p.r, binary.LittleEndian, &bucketHeader.ID)
	if err != nil {
		return nil, err
	}
	err = binary.Read(p.r, binary.LittleEndian, &bucketHeader.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Read(p.r, binary.LittleEndian, &bucketHeader.Alg)
	if err != nil {
		return nil, err
	}
	err = binary.Read(p.r, binary.LittleEndian, &bucketHeader.Hash)
	if err != nil {
		return nil, err
	}

	var weight uint32
	err = binary.Read(p.r, binary.LittleEndian, &weight)
	if err != nil {
		return nil, err
	}
	bucketHeader.Weight = math.Float32frombits(weight)
	err = binary.Read(p.r, binary.LittleEndian, &bucketHeader.Size)
	if err != nil {
		return nil, err
	}

	bucketCommon := binaryBucketCommon{binaryBucketHeader: bucketHeader}
	bucketCommon.Items = make([]int32, bucketHeader.Size)
	for i := uint32(0); i < bucketHeader.Size; i++ {
		err = binary.Read(p.r, binary.LittleEndian, &bucketCommon.Items[i])
		if err != nil {
			return nil, err
		}
	}
	switch bucketHeader.Alg {
	case CrushAlgUniform:
		bucketUniform := &BucketUniform{
			binaryBucketCommon: bucketCommon,
		}
		var itemWeight uint32
		err = binary.Read(p.r, binary.LittleEndian, &itemWeight)
		if err != nil {
			return nil, err
		}
		bucketUniform.ItemWeight = math.Float32frombits(itemWeight)
		bucket = bucketUniform
	case CrushAlgList:
		itemWeights := make([]uint32, bucketHeader.Size)
		sumWeights := make([]uint32, bucketHeader.Size)
		bucketList := &BucketList{
			binaryBucketCommon: bucketCommon,
			ItemWeights:        make([]float32, bucketHeader.Size),
			SumWeights:         make([]float32, bucketHeader.Size),
		}
		for i := uint32(0); i <= bucketHeader.Size; i++ {
			err = binary.Read(p.r, binary.LittleEndian, &itemWeights[i])
			if err != nil {
				return nil, err
			}
			bucketList.ItemWeights[i] = math.Float32frombits(itemWeights[i])
			err = binary.Read(p.r, binary.LittleEndian, &sumWeights[i])
			if err != nil {
				return nil, err
			}
			bucketList.SumWeights[i] = math.Float32frombits(sumWeights[i])
		}
		bucket = bucketList
	case CrushAlgTree:
		bucketTree := &BucketTree{
			binaryBucketCommon: bucketCommon,
		}
		err = binary.Read(p.r, binary.LittleEndian, &bucketTree.NumNodes)
		if err != nil {
			return nil, err
		}
		nodeWeights := make([]uint32, bucketTree.NumNodes*4)
		bucketTree.NodeWeights = make([]float32, bucketTree.NumNodes*4)
		err = binary.Read(p.r, binary.LittleEndian, &nodeWeights)
		if err != nil {
			return nil, err
		}
		for i := 0; i < int(bucketTree.NumNodes*4); i++ {
			bucketTree.NodeWeights[i] = math.Float32frombits(nodeWeights[i])
		}
		bucket = bucketTree
	case CrushAlgStraw:
		itemWeights := make([]uint32, (bucketHeader.Size)*4)
		bucketStraw := &BucketStraw{
			binaryBucketCommon: bucketCommon,
			Straws:             make([]uint32, (bucketHeader.Size)*4),
			ItemWeights:        make([]float32, (bucketHeader.Size)*4),
		}
		for i := uint32(0); i < bucketHeader.Size; i++ {
			err = binary.Read(p.r, binary.LittleEndian, &itemWeights[i])
			if err != nil {
				return nil, err
			}
			bucketStraw.ItemWeights[i] = math.Float32frombits(itemWeights[i])
			err = binary.Read(p.r, binary.LittleEndian, &bucketStraw.Straws[i])
			if err != nil {
				return nil, err
			}
		}
		bucket = bucketStraw
	case CrushAlgStraw2:
		itemWeights := make([]uint32, (bucketHeader.Size+1)*4)
		bucketStraw2 := &BucketStraw2{
			binaryBucketCommon: bucketCommon,
			ItemWeights:        make([]float32, (bucketHeader.Size+1)*4),
		}
		err = binary.Read(p.r, binary.LittleEndian, &itemWeights)
		if err != nil {
			return nil, err
		}
		for i := uint32(0); i < (bucketHeader.Size+1)*4; i++ {
			bucketStraw2.ItemWeights[i] = math.Float32frombits(itemWeights[i])
		}
		bucket = bucketStraw2
	}
	ibucket.ID = bucketHeader.ID
	ibucket.Alg = bucketHeader.Alg.String()
	ibucket.Hash = bucketHeader.Hash.String()
	ibucket.TypeID = bucketHeader.Type
	ibucket.Weight = bucketHeader.Weight
	ibucket.Size = bucketHeader.Size
	_ = bucket
	return ibucket, nil
}
