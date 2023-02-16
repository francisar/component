package mutex

import "fmt"

type lockError struct {
	key string
	lockService string
	err error
	msg string
}

func (e *lockError) Error() string  {
	msg := fmt.Sprintf("LockError lock key:%s lockService:%s %s", e.key, e.lockService, e.msg)
	if e.err != nil {
		msg = fmt.Sprintf("%s Wrap Error:%s", msg, e.err.Error())
	}
	return msg
}

func (e *lockError)WrapError(err error) *lockError {
	e.err = err
	return e
}

func (e *lockError)WrapMsg(msg string) *lockError {
	e.msg = msg
	return e
}

type unLockError struct {
	key string
	lockService string
	err error
	msg string
}

func (e *unLockError) Error() string  {
	msg := fmt.Sprintf("UnLockError lock key:%s lockService:%s %s", e.key, e.lockService, e.msg)
	if e.err != nil {
		msg = fmt.Sprintf("%s Wrap Error:%s", msg, e.err.Error())
	}
	return msg
}

func (e *unLockError)WrapError(err error) *unLockError {
	e.err = err
	return e
}

func (e *unLockError)WrapMsg(msg string) *unLockError {
	e.msg = msg
	return e
}