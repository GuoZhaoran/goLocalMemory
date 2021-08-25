package setGroup

import (
	"container/list"
	"fmt"
	"math"
	"sync"
)

// memberUsedRecord is the member use record
type memberUsedRecord struct {
	member    interface{}
	usedTimes uint64
}

// setGroup Define the setGroup structure
type setGroup struct {
	memberIntMap          map[interface{}]uint64
	indexMemberSlice      []memberUsedRecord
	groupSetRelationalMap map[interface{}][][1024]uint64
	indexProduceCounter   uint64
	deleteIndexLinkedList *list.List
	mutex                 sync.RWMutex
}

// New create a setGroup instance
func New() *setGroup {
	return &setGroup{
		memberIntMap:          make(map[interface{}]uint64),
		groupSetRelationalMap: make(map[interface{}][][1024]uint64),
		deleteIndexLinkedList: list.New(),
	}
}

// Add add one member to one set
func (setGroup *setGroup) Add(setKey interface{}, member interface{}) {
	// Locking ensures secure write operations
	setGroup.mutex.Lock()
	defer setGroup.mutex.Unlock()

	// Check if member exists in setGroup, if not exists, retrieve a
	memberIndex := setGroup.generateMemberIndex(member)
	bitMapBucketIndex := memberIndex >> 16

	// Detect the existence of setKey
	setRelationalSlice, ok := setGroup.groupSetRelationalMap[setKey]
	if !ok {
		setRelationalSlice = [][1024]uint64{}
	}
	relationalSliceLen := uint64(len(setRelationalSlice))
	// capacity expansion the slice
	for i := bitMapBucketIndex; i >= relationalSliceLen; i-- {
		setRelationalSlice = append(setRelationalSlice, [1024]uint64{})
		if i == 0 {
			break
		}
	}
	setGroup.groupSetRelationalMap[setKey] = setRelationalSlice
	// find uint64 index location
	uint64ValuePosition := memberIndex & math.MaxUint16
	// Change the specified bit of uint64 to 1
	setRelationalSlice[bitMapBucketIndex][uint64ValuePosition>>6] |= 1 << (uint64ValuePosition & 63)
}

// generateMemberIndex : select the memberIndex
func (setGroup *setGroup) generateMemberIndex(member interface{}) uint64 {
	memberIndex, exists := setGroup.memberIntMap[member]
	if !exists {
		// first try to get from the list
		ele := setGroup.deleteIndexLinkedList.Front()
		if ele != nil {
			setGroup.deleteIndexLinkedList.Remove(ele)
			memberIndex = ele.Value.(uint64)
			setGroup.indexMemberSlice[memberIndex].member = member
			setGroup.indexMemberSlice[memberIndex].usedTimes++
		} else {
			// if list is empty, then take the value of counter and increment counter by itself
			memberIndex = setGroup.indexProduceCounter
			setGroup.indexProduceCounter++
			// add to index member slice
			var initUsedTime uint64 = 1
			newIndexMemberSlice := append(setGroup.indexMemberSlice, memberUsedRecord{
				member:    member,
				usedTimes: initUsedTime,
			})
			setGroup.indexMemberSlice = newIndexMemberSlice
		}
		// add member to memberIntMap
		setGroup.memberIntMap[member] = memberIndex
	} else {
		setGroup.indexMemberSlice[memberIndex].usedTimes++
	}

	return memberIndex
}

// Remove ：remove member from setKey
func (setGroup *setGroup) Remove(setKey interface{}, member interface{}) (setKeyExist bool, memberExist bool) {
	// Locking ensures secure write operations
	setGroup.mutex.Lock()
	defer setGroup.mutex.Unlock()

	//setKey not exist
	setRelationalSlice, ok := setGroup.groupSetRelationalMap[setKey]
	if !ok {
		return false, false
	}

	//determine if member exists
	memberIndex, exists := setGroup.memberIntMap[member]
	if !exists || setGroup.indexMemberSlice[memberIndex].usedTimes == 0 {
		return true, false
	}
	// find uint64 index location
	bitMapBucketIndex := memberIndex >> 16
	uint64ValuePosition := memberIndex & math.MaxUint16
	memberExist = ((setRelationalSlice[bitMapBucketIndex][uint64ValuePosition>>6] >> (uint64ValuePosition & 63)) & 1) == 1
	if !memberExist {
		return true, false
	}
	setRelationalSlice[bitMapBucketIndex][uint64ValuePosition>>6] &^= 1 << (uint64ValuePosition & 63)
	setGroup.indexMemberSlice[memberIndex].usedTimes--
	if setGroup.indexMemberSlice[memberIndex].usedTimes == 0 {
		setGroup.deleteIndexLinkedList.PushBack(memberIndex)
		delete(setGroup.memberIntMap, member)
	}

	return true, true
}

//Intersect calculation a group set intersection
func (setGroup *setGroup) Intersect(setKeys ...interface{}) []interface{} {
	setGroup.mutex.RLock()
	defer setGroup.mutex.RUnlock()

	setKeysLen := len(setKeys)
	setKeysRelational := make([][][1024]uint64, setKeysLen)
	var intersectResult []interface{}
	minIterationsOffset := math.MaxUint32
	for setKeyIndex, setKey := range setKeys {
		// judgment the set exists, if not, return []interface{}{}
		setMemberBitMap, setExist := setGroup.groupSetRelationalMap[setKey]
		if !setExist {
			return intersectResult
		}
		if setMemberBitMap == nil || len(setMemberBitMap) == 0 {
			return intersectResult
		}
		if len(setMemberBitMap) < minIterationsOffset {
			minIterationsOffset = len(setMemberBitMap)
		}
		setKeysRelational[setKeyIndex] = setMemberBitMap
	}

	// iterate from minIterationsOffset to 0 find the intersection
	for iter := minIterationsOffset - 1; iter >= 0; iter-- {
		for bitMapIndex := 1023; bitMapIndex >= 0; bitMapIndex-- {
			itemIterBitMap := uint64(math.MaxUint64)
			for _, setMemberBitMap := range setKeysRelational {
				itemIterBitMap &= setMemberBitMap[iter][bitMapIndex]
				if itemIterBitMap == 0 {
					break
				}
			}
			if itemIterBitMap > 0 {
				for nBit := 0; nBit < 63; nBit++ {
					if (itemIterBitMap>>nBit)&1 == 1 {
						intersectResult = append(intersectResult, setGroup.indexMemberSlice[iter*(1<<16)+bitMapIndex*64+nBit].member)
					}
				}
			}
		}
	}

	return intersectResult
}

//IsSetKey determine if setKey exists
func (setGroup *setGroup) IsSetKey(setKey interface{}) bool {
	setGroup.mutex.RLock()
	defer setGroup.mutex.RUnlock()
	_, ok := setGroup.groupSetRelationalMap[setKey]
	return ok
}

//IsSetMember determine if member exists in setKey
func (setGroup *setGroup) IsSetMember(setKey interface{}, member interface{}) (bool, bool) {
	setGroup.mutex.RLock()
	defer setGroup.mutex.RUnlock()
	setRelationalSlice, ok := setGroup.groupSetRelationalMap[setKey]
	if !ok {
		return false, false
	}
	memberIndex := setGroup.memberIntMap[member]
	bitMapBucketIndex := memberIndex >> 16
	uint64ValuePosition := memberIndex & math.MaxUint16

	return ((setRelationalSlice[bitMapBucketIndex][uint64ValuePosition>>6] >> (uint64ValuePosition & 63)) & 1) == 1, true
}

//FPrint format print the struct of setGroup
func (setGroup *setGroup) FPrint() {
	setGroup.mutex.RLock()
	defer setGroup.mutex.RUnlock()
	fmt.Println("----------- memberIntMap -----------")
	for k, v := range setGroup.memberIntMap {
		fmt.Println(k, v)
	}
	fmt.Println("----------- indexMemberSlice -----------")
	for k, v := range setGroup.indexMemberSlice {
		fmt.Println("index:", k, "member:", v.member, "used times:", v.usedTimes)
	}
	fmt.Println("----------- groupSetRelationalMap -----------")
	for k, v := range setGroup.groupSetRelationalMap {
		fmt.Println("set:", k, "bitmap:", v[0][0])
	}
	fmt.Println("----------- indexProduceCounter -----------", setGroup.indexProduceCounter)
	fmt.Println("----------- deleteIndexLinkedList -----------")
	for element := setGroup.deleteIndexLinkedList.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}
}
