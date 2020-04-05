/**
 * @Author: victor
 * @Description:
 * @File:  redislock
 * @Version: 1.0.0
 * @Date: 2020/4/5 9:48 上午
 */

package redislock

/*
func main() {
	fmt.Println("start")
	DefaultTimeout := 10
	conn, err := redis.Dial("tcp", "localhost:6379")

	lock, ok, err := TryLockWithTimeout(database.GetRedisConn(), "test", "1", time.Minute)
	if err != nil {
		log.Fatal("Error while attempting lock")
	}
	if !ok {
		log.Fatal("bug")
	}
	lock.AddTimeout(100)

	time.Sleep(time.Duration(DefaultTimeout) * time.Second)
	fmt.Println("end")
	defer lock.Unlock()
}
*/

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// DefaultTimeout is the duration for which the lock is valid
const DefaultTimeout = 10 * time.Minute

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

// Lock represents a held lock.
type Lock struct {
	resource string
	token    string
	conn     redis.Conn
	timeout  time.Duration
}

func (lock *Lock) tryLock() (ok bool, err error) {
	status, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int64(lock.timeout/time.Second), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return status == "OK", nil
}

// Unlock releases the lock. If the lock has timed out, it silently fails without error.
func (lock *Lock) Unlock() (err error) {
	_, err = unlockScript.Do(lock.conn, lock.key(), lock.token)
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// TryLock attempts to acquire a lock on the given resource in a non-blocking manner.
// The lock is valid for the duration specified by DefaultTimeout.
func TryLock(conn redis.Conn, resource string, token string) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resource, token, DefaultTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout time.Duration) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, token, conn, timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}