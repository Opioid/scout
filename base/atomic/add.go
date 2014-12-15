package atomic

import (
	"sync/atomic"
	"unsafe"
	"math"
)

func AddFloat32(addr *float32, delta float32) (new float32) {
	unsafeAddr := (*uint32)(unsafe.Pointer(addr))

	for {
		oldValue := math.Float32bits(*addr)
		new       = *addr + delta
		newValue := math.Float32bits(new)

		if atomic.CompareAndSwapUint32(unsafeAddr, oldValue, newValue) {
			return
		}
	}
}