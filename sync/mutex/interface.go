package mutex


type LockService interface {
	Lock(lockKey *LockKey) *Status
	Unlock(lockKey *LockKey) (bool, error)
}