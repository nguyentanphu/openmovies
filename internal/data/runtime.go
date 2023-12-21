package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (rt Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", rt)
	quoted := strconv.Quote(jsonValue)
	return []byte(quoted), nil
}
