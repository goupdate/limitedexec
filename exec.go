package limitedexec

import (
	"sync"
	"time"
)

/*
	Exec(K, N, F)
	Monitors to ensure that function F identified by key K doesn't execute more than once within duration N, skipping calls.

	Usage example: A frequently recurring identical event that clutters logs:
	12:12: came banned ip xx
	12:13: came banned ip xx
	12:14: came banned ip xx
	13:15: came banned ip xx
	14:16: came banned ip xx

	limitedexec.Exec("banned "+ip, time.Minute(), func(){ log.Print("came banned ip") })

	The above will result in:
	12:12: came banned ip xx
	13:15: came banned ip xx
	14:16: came banned ip xx
*/

var check_time time.Time //clean up
var limited sync.Map
var mutex sync.Mutex

/*
	Run function "fun" for the key "k" not often then "t":
	fun-k1
	fun-k2
	...<t>...
	fun-k1
	fun-k2
	fun-k3
	....<t>...
	fun-k3
*/
func Exec(k interface{}, t time.Duration, fun func()) {
	allow_run := true
	mutex.Lock()
	old, ok := limited.Load(k)
	if ok {
		if time.Now().Before(old.(time.Time)) {
			allow_run = false
		}
	}
	mutex.Unlock()

	if !allow_run {
		return
	}

	limited.Store(k, time.Now().Add(t))

	mutex.Lock()
	if time.Now().After(check_time) {
		check_time = time.Now().Add(time.Minute)
		go func() {
			mutex.Lock()
			limited.Range(func(k interface{}, old interface{}) bool {
				if time.Now().After(old.(time.Time)) {
					limited.Delete(k)
				}
				return true
			})
			mutex.Unlock()
		}()
	}
	mutex.Unlock()

	fun()
}
