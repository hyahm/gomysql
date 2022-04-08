package gomysql

type Counter interface {
	uint8 | uint16 | uint | int | uint32 | uint64 |
		int32 | int64 | int8 | int16 |
		float32 | float64
}
