package life

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// logger is a simple function to disable logs if Game.printlog is set to false
func logger(a ...interface{}) {
	if PRINTLOG {
		fmt.Println(strings.Trim(fmt.Sprintf("%v", a), "[]"))
	}
}

// flip a d2 to get a random true of false value
func d2(seed int64) bool {
	s1 := rand.NewSource(time.Now().UnixNano() + seed)
	r1 := rand.New(s1)
	return r1.Intn(2) == 1
}
