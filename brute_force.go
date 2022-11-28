package bruteforce

import (
	"net"
	"strconv"
	"time"

	"github.com/muesli/cache2go"
)

const (
	BRUTE_TIME = 60 //1 minute
)

var (
	Bf_tables []*cache2go.CacheTable = make([]*cache2go.CacheTable, BF_PROTO_END)
)

type bf_counter struct {
	times_arr         []uint
	last_update_index int16
	last_update_time  time.Time
	total_count       uint
}

func insertIfNotExist(key interface{}, args ...interface{}) *cache2go.CacheItem {
	val := bf_counter{
		times_arr:         make([]uint, BRUTE_TIME),
		last_update_index: 0,
		last_update_time:  args[0].(time.Time),
	}
	item := cache2go.NewCacheItem(key, BRUTE_TIME*time.Second, &val)
	return item
}

func InitBruteForce() {
	for i := BF_PROTO_HTTP; i < BF_PROTO_END; i++ {
		Bf_tables[i] = cache2go.Cache(strconv.Itoa(i))
		Bf_tables[i].SetDataLoader(insertIfNotExist)
	}
}

func update_bf_item(bf *bf_counter, current time.Time) uint {
	dura := current.Sub(bf.last_update_time)
	bf.last_update_time = current
	if dura < 1*time.Second {
		bf.times_arr[bf.last_update_index]++
		bf.total_count++
	} else {
		idx := (bf.last_update_index + int16(dura.Seconds())) % BRUTE_TIME
		i := (bf.last_update_index + 1) % BRUTE_TIME
		for ; i != idx; i = (i + 1) % BRUTE_TIME {
			bf.total_count -= bf.times_arr[i]
			bf.times_arr[i] = 0
		}
		bf.total_count -= bf.times_arr[i]
		bf.last_update_index = idx
		bf.times_arr[i] = 1
		bf.total_count++
	}
	return bf.total_count
}

func BruteForceCheck(proto uint, ip *net.IP) (bool, error) {
	if proto >= BF_PROTO_END {
		return false, ErrProtocolIdInvalid
	}
	if Bf_setting[proto] == 0 {
		return false, nil
	}
	table := Bf_tables[proto]
	current_time := time.Now()
	item, err := table.Value(ip, current_time)
	if err != nil {
		return false, err
	}
	item.RWMutex.Lock()
	bf := item.Data().(*bf_counter)
	count := update_bf_item(bf, current_time)
	if count >= Bf_setting[proto] {
		bf.times_arr = make([]uint, BRUTE_TIME)
		bf.total_count = 0
		item.RWMutex.Unlock()
		return true, nil
	}
	item.RWMutex.Unlock()
	return false, nil
}
