package onerest

import (
	"fmt"
)

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
