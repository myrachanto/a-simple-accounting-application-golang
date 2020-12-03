package support


import (
				"fmt"
				"strconv"
)

func str() {

				var uval uint64 = 1 * (1 << 20)

				str := strconv.FormatUint(uval, 10)

				fmt.Println(str) // uint64 in string format

}