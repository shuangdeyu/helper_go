package encrypt

import (
	"errors"
	"fmt"
	"sync/atomic"
)

/**
 * java的随机数获取算法，与golang的随机数算法不同
 */

type RandomJava struct {
	seed int64
}

const (
	mask       = int64(1<<48) - 1 // 2^48 - 1
	multiplier = int64(25214903917)
	addend     = int64(11)
)

func NewRandomJava(seed int64) *RandomJava {
	return &RandomJava{seed: (seed ^ multiplier) & (mask - 1)}
}

func (r *RandomJava) next(bits int) int {
	var oldSeed, nextSeed int64
	for {
		oldSeed = atomic.LoadInt64(&r.seed)
		nextSeed = (oldSeed*multiplier + addend) & mask
		if atomic.CompareAndSwapInt64(&r.seed, oldSeed, nextSeed) {
			break
		}
	}
	return int(nextSeed >> (48 - bits))
}

func (r *RandomJava) nextInt(bound int) (int, error) {
	if bound <= 0 {
		return 0, errors.New(fmt.Sprintf("BadBound is %d", bound))
	}

	rand := r.next(31)
	m := bound - 1
	if (bound & m) == 0 {
		return int((int64(bound) * int64(rand)) >> 31), nil
	}

	for {
		u := rand
		rand = u % bound
		if (u - rand + m) >= 0 {
			break
		}
		u = r.next(31)
	}
	return rand, nil
}
