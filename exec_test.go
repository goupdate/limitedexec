package limitedexec

import (
	"testing"
	"time"
)

var a = 0
var b = 0

func incA() {
	a++
}

func incB() {
	b++
}

func Test_Exec(t *testing.T) {

	Exec(1, time.Second, incA) //a+1
	Exec(1, time.Second, incA) //-
	Exec(1, time.Second, incA) //-

	Exec("2", time.Second, incA) //a+1

	Exec(3, time.Second*2, incB) //b+1

	time.Sleep(time.Second) //wait

	Exec(1, time.Second, incA) //a+1
	Exec(1, time.Second, incA) //-
	Exec(1, time.Second, incA) //-

	Exec(3, time.Second*2, incB) //-

	if a != 3 || b != 1 {
		t.Fatalf("err: a:%d b:%d", a, b)
	}

	time.Sleep(time.Millisecond * 2100)

	check_time = check_time.Add(-time.Minute)

	Exec("runme", time.Minute, incA) //a+1

	time.Sleep(time.Millisecond * 100) //give some time for clean up

	mutex.Lock()
	limited_len := 0
	limited.Range(func(k interface{}, v interface{}) bool {
		limited_len++
		t.Logf("limited: %v %v\n", k, v.(time.Time).String())
		return true
	})
	mutex.Unlock()

	if limited_len != 1 {
		t.Fatalf("limited len: %d\n now: %s", limited_len, time.Now().String())
	}
	if _, ok := limited.Load("runme"); !ok {
		t.Fatal("limited runme")
	}
}
