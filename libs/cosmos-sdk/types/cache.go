package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	maxAccInMap        = 100000
	deleteAccCount     = 10000
	maxStorageInMap    = 10000000
	deleteStorageCount = 1000000

	FlagMultiCache         = "multi-cache"
	MaxAccInMultiCache     = "multi-cache-acc"
	MaxStorageInMultiCache = "multi-cache-storage"
	UseCache               bool
)

type account interface {
	Copy() interface{}
	GetAddress() AccAddress
	SetAddress(AccAddress) error
	GetPubKey() crypto.PubKey
	SetPubKey(crypto.PubKey) error
	GetAccountNumber() uint64
	SetAccountNumber(uint64) error
	GetSequence() uint64
	SetSequence(uint64) error
	GetCoins() Coins
	SetCoins(Coins) error
	SpendableCoins(blockTime time.Time) Coins
	String() string
}

type storageWithCache struct {
	Value  []byte
	Dirty  bool
	Delete bool
}

type accountWithCache struct {
	Acc      account
	Gas      uint64
	Bz       []byte
	IsDirty  bool
	ISDelete bool
}

type codeWithCache struct {
	Code    []byte
	IsDirty bool
}

type Cache struct {
	mu sync.Mutex

	useCache  bool
	parent    *Cache
	gasConfig types.GasConfig

	dirtyStorageMap map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache
	readStorageMap  map[ethcmn.Address]map[ethcmn.Hash][]byte

	dirtyaccMap map[ethcmn.Address]*accountWithCache
	readaccMap  map[ethcmn.Address]*accountWithCache

	dirtycodeMap map[ethcmn.Hash]*codeWithCache
	readcodeMap  map[ethcmn.Hash][]byte
}

func initCacheParam() {
	UseCache = viper.GetBool(FlagMultiCache)

	if data := viper.GetInt(MaxAccInMultiCache); data != 0 {
		maxAccInMap = data
		deleteAccCount = maxAccInMap / 10
	}

	if data := viper.GetInt(MaxStorageInMultiCache); data != 0 {
		maxStorageInMap = data
		deleteStorageCount = maxStorageInMap / 10
	}
}

func NewChainCache() *Cache {
	initCacheParam()
	return NewCache(nil, UseCache)
}

func NewCache(parent *Cache, useCache bool) *Cache {
	return &Cache{
		mu: sync.Mutex{},

		useCache: useCache,
		parent:   parent,

		dirtyStorageMap: make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache, 0),
		readStorageMap:  make(map[ethcmn.Address]map[ethcmn.Hash][]byte, 0),

		dirtyaccMap: make(map[ethcmn.Address]*accountWithCache, 0),
		readaccMap:  make(map[ethcmn.Address]*accountWithCache, 0),

		dirtycodeMap: make(map[ethcmn.Hash]*codeWithCache),
		readcodeMap:  make(map[ethcmn.Hash][]byte),
		gasConfig:    types.KVGasConfig(),
	}

}

func (c *Cache) UseCache() bool {
	return !c.skip()
}

func (c *Cache) skip() bool {
	if c == nil || !c.useCache {
		return true
	}
	return false
}

func (c *Cache) DeleteCacheAccount(acc ethcmn.Address) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.readaccMap, acc)
}

func (c *Cache) UpdateAccount(addr AccAddress, acc account, bz []byte, isDirty bool, isDelete bool) {
	if c.skip() {
		return
	}
	ethAddr := ethcmn.BytesToAddress(addr.Bytes())

	tt := &accountWithCache{
		Acc:      acc,
		IsDirty:  isDirty,
		ISDelete: isDelete,
		Bz:       bz,
		Gas:      types.Gas(len(bz))*c.gasConfig.ReadCostPerByte + c.gasConfig.ReadCostFlat,
	}

	c.mu.Lock()
	if !isDirty {
		c.setReadAccount(ethAddr, acc, bz, tt.Gas)
	} else {
		c.dirtyaccMap[ethAddr] = tt
	}
	c.mu.Unlock()
}

func (c *Cache) UpdateStorage(addr ethcmn.Address, key ethcmn.Hash, value []byte, isDirty bool, isDelete bool) {
	//fmt.Println("uuuuuu", addr.String(), key.String(), isDirty, isDelete)
	if c.skip() {
		fmt.Println("skip----")
		return
	}

	c.mu.Lock()
	if isDirty {
		if _, ok := c.dirtyStorageMap[addr]; !ok {
			c.dirtyStorageMap[addr] = make(map[ethcmn.Hash]*storageWithCache, 0)
		}

		c.dirtyStorageMap[addr][key] = &storageWithCache{
			Value:  value,
			Dirty:  isDirty,
			Delete: isDelete,
		}
	} else {
		c.setReadStorage(addr, key, value)
	}
	c.mu.Unlock()
}

func (c *Cache) UpdateCode(key []byte, value []byte, isdirty bool) {
	if c.skip() {
		return
	}
	hash := ethcmn.BytesToHash(key)
	c.mu.Lock()
	if isdirty {
		c.dirtycodeMap[hash] = &codeWithCache{
			Code:    value,
			IsDirty: isdirty,
		}
	} else {
		c.SetReadCode(hash, value)
	}

	c.mu.Unlock()
}

func (c *Cache) GetAccount(addr ethcmn.Address) (account, uint64, []byte, bool) {
	if c.skip() {
		return nil, 0, nil, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if data, ok := c.dirtyaccMap[addr]; ok {
		return data.Acc, data.Gas, data.Bz, ok
	}

	if data, ok := c.readaccMap[addr]; ok {
		return data.Acc, data.Gas, data.Bz, ok
	}

	if c.parent != nil {
		acc, gas, bz, ok := c.parent.GetAccount(addr)
		return acc, gas, bz, ok
	}
	return nil, 0, nil, false
}

func (c *Cache) setReadAccount(addr ethcmn.Address, acc account, bz []byte, gas uint64) {
	c.readaccMap[addr] = &accountWithCache{
		Acc:     acc,
		Gas:     gas,
		Bz:      bz,
		IsDirty: false,
	}
}

func (c *Cache) setReadStorage(addr ethcmn.Address, key ethcmn.Hash, value []byte) {
	if _, ok := c.readStorageMap[addr]; !ok {
		c.readStorageMap[addr] = make(map[ethcmn.Hash][]byte)
	}
	c.readStorageMap[addr][key] = value
}

func (c *Cache) SetReadCode(hash ethcmn.Hash, value []byte) {
	c.readcodeMap[hash] = value
}
func (c *Cache) GetStorage(addr ethcmn.Address, key ethcmn.Hash) ([]byte, bool) {
	if c.skip() {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, hasAddr := c.dirtyStorageMap[addr]; hasAddr {
		data, hasKey := c.dirtyStorageMap[addr][key]
		if hasKey {
			return data.Value, hasKey
		}
	} else {
		c.dirtyStorageMap[addr] = make(map[ethcmn.Hash]*storageWithCache)
	}

	if _, hasAddr := c.readStorageMap[addr]; hasAddr {
		if data, hasKey := c.readStorageMap[addr][key]; hasKey {
			return data, true
		}
	}

	if c.parent != nil {
		value, ok := c.parent.GetStorage(addr, key)
		return value, ok
	}
	return nil, false
}

func (c *Cache) GetCode(key []byte) ([]byte, bool) {
	if c.skip() {
		return nil, false
	}

	hash := ethcmn.BytesToHash(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	if data, ok := c.dirtycodeMap[hash]; ok {
		return data.Code, ok
	}

	if data, ok := c.readcodeMap[hash]; ok {
		return data, ok
	}
	if c.parent != nil {
		code, ok := c.parent.GetCode(hash.Bytes())
		return code, ok
	}
	return nil, false
}

func (c *Cache) GetDirtyAcc() map[ethcmn.Address]*accountWithCache {
	return c.dirtyaccMap
}

func (c *Cache) GetDirtyCode() map[ethcmn.Hash]*codeWithCache {
	return c.dirtycodeMap
}

func (c *Cache) GetDirtyStorage() map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache {
	return c.dirtyStorageMap
}
func (c *Cache) Write(updateDirty bool, printLog bool) {
	if c.skip() {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.parent == nil {
		return
	}

	c.writeStorage(updateDirty, printLog)
	c.writeAcc(updateDirty)
	c.writeCode(updateDirty)
}

func (c *Cache) writeStorage(updateDirty bool, printLog bool) {
	for addr, storages := range c.dirtyStorageMap {
		if _, ok := c.parent.dirtyStorageMap[addr]; !ok {
			c.parent.dirtyStorageMap[addr] = make(map[ethcmn.Hash]*storageWithCache, 0)
		}

		for key, v := range storages {
			if updateDirty {
				if printLog {
					if addr.String() == "0xd90838EC67025E2b92d814AAe244Ac4ed889994D" {
						fmt.Println("addr", addr.String(), key.String(), hex.EncodeToString(v.Value))
					}

				}
				c.parent.dirtyStorageMap[addr][key] = v
			}
		}
	}

	for addr, storages := range c.readStorageMap {
		if _, ok := c.parent.readStorageMap[addr]; !ok {
			c.parent.readStorageMap[addr] = make(map[ethcmn.Hash][]byte, 0)
		}

		for key, v := range storages {
			if updateDirty {
				//if addr.String() == "0xadf4916d11F352a2748e19F3056428639313F6E1" {
				//fmt.Println("writeStorage", addr.String(), key.String(), hex.EncodeToString(v.value))
				//}
				c.parent.readStorageMap[addr][key] = v
			}
		}
	}

	c.dirtyStorageMap = make(map[ethcmn.Address]map[ethcmn.Hash]*storageWithCache)
	c.readStorageMap = make(map[ethcmn.Address]map[ethcmn.Hash][]byte)
}

func (c *Cache) writeAcc(updateDirty bool) {
	for addr, v := range c.dirtyaccMap {
		if updateDirty {
			c.parent.dirtyaccMap[addr] = v
		}
	}

	for addr, v := range c.readaccMap {
		if updateDirty {
			c.parent.readaccMap[addr] = v
		}
	}
	c.dirtyaccMap = make(map[ethcmn.Address]*accountWithCache)
	c.readaccMap = make(map[ethcmn.Address]*accountWithCache)
}

func (c *Cache) writeCode(updateDirty bool) {
	for hash, v := range c.dirtycodeMap {
		if updateDirty {
			c.parent.dirtycodeMap[hash] = v
		}
	}
	for hash, v := range c.readcodeMap {
		if updateDirty {
			c.parent.readcodeMap[hash] = v
		}
	}
	c.dirtycodeMap = make(map[ethcmn.Hash]*codeWithCache)
	c.readcodeMap = make(map[ethcmn.Hash][]byte)
}

func (c *Cache) IsConflict(newCache *Cache) bool {

	c.mu.Lock()
	defer c.mu.Unlock()
	//fmt.Println("readStorageMap", len(newCache.readaccMap), len(newCache.readStorageMap), len(newCache.readcodeMap))

	for acc, v := range newCache.readaccMap {
		if data, ok := c.dirtyaccMap[acc]; ok && data.IsDirty {
			if !bytes.Equal(v.Bz, data.Bz) {
				fmt.Println("conflict-acc", acc.String())
				return true
			}
		}
	}

	for acc, ss := range newCache.readStorageMap {
		preSS, ok := c.dirtyStorageMap[acc]
		if !ok {
			continue
		}
		for kk, vv := range ss {
			if pp, ok1 := preSS[kk]; ok1 && pp.Dirty {
				if !bytes.Equal(pp.Value, vv) {
					fmt.Println("conflict-storage", acc.String(), kk.String(), "now", hex.EncodeToString(pp.Value), "read", hex.EncodeToString(vv))
					return true
				}
			}

		}
	}

	for acc, code := range newCache.readcodeMap {
		if data, ok := c.dirtycodeMap[acc]; ok && data.IsDirty {
			if !bytes.Equal(code, data.Code) {
				fmt.Println("conflict-code", acc.String())
				return true
			}
		}
	}
	return false
}

func (c *Cache) storageSize() int {
	lenStorage := 0
	for _, v := range c.dirtyStorageMap {
		lenStorage += len(v)
	}
	return lenStorage
}

//func (c *Cache) TryDelete(logger log.Logger, height int64) {
//	if c.skip() {
//		return
//	}
//	if height%100 == 0 {
//		c.logInfo(logger, "null")
//	}
//
//	lenStorage := c.storageSize()
//	if len(c.accMap) < maxAccInMap && lenStorage < maxStorageInMap {
//		return
//	}
//
//	deleteMsg := ""
//	if len(c.accMap) >= maxAccInMap {
//		deleteMsg += fmt.Sprintf("Acc:Deleted Before:%d", len(c.accMap))
//		cnt := 0
//		for key := range c.accMap {
//			delete(c.accMap, key)
//			cnt++
//			if cnt > deleteAccCount {
//				break
//			}
//		}
//	}
//
//	if lenStorage >= maxStorageInMap {
//		deleteMsg += fmt.Sprintf("Storage:Deleted Before:len(contract):%d, len(storage):%d", len(c.storageMap), lenStorage)
//		cnt := 0
//		for key, value := range c.storageMap {
//			cnt += len(value)
//			delete(c.storageMap, key)
//			if cnt > deleteStorageCount {
//				break
//			}
//		}
//	}
//	if deleteMsg != "" {
//		c.logInfo(logger, deleteMsg)
//	}
//}

//func (c *Cache) logInfo(logger log.Logger, deleteMsg string) {
//	nowStats := fmt.Sprintf("len(acc):%d len(contracts):%d len(storage):%d", len(c.accMap), len(c.storageMap), c.storageSize())
//	logger.Info("MultiCache", "deleteMsg", deleteMsg, "nowStats", nowStats)
//}

func (c *Cache) GetParent() *Cache {
	return c.parent
}

func (c *Cache) Print(printLog bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//fmt.Println("size::::", len(c.dirtyaccMap), len(c.dirtyStorageMap), len(c.dirtycodeMap))
	//fmt.Println("size::::", len(c.readaccMap), len(c.readStorageMap), len(c.readcodeMap))

	if !printLog {
		return
	}
	for acc, v := range c.dirtyaccMap {
		fmt.Println("acc:", acc.String(), v.IsDirty)
	}
	for acc, v := range c.dirtyStorageMap {
		for kk, vv := range v {
			fmt.Println("storage:", acc.String(), kk.String(), hex.EncodeToString(vv.Value), vv.Dirty)
		}
	}
	for acc, _ := range c.dirtycodeMap {
		fmt.Println("code:", acc.String())
	}
	//for acc, v := range c.accMap {
	//	fmt.Println("acc", acc.String(), v.isDirty)
	//}
	//
	//for acc, v := range c.storageMap {
	//	fmt.Println("storage", acc.String(), v)
	//}
	//for acc, v := range c.codeMap {
	//	fmt.Println("code", acc.String(), v.isDirty)
	//}
}
