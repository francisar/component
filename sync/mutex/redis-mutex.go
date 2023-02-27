package mutex

import (
	"context"
	"github.com/francisar/component/nosql/redis"
	"time"
)

const lockService = "RedisMutexService"

func newLockError(lockKey *LockKey) *lockError {
	lockErr := lockError{
		key: lockKey.Key,
		lockService: lockService,
	}
	return &lockErr
}

func newUnLockError(key string) *unLockError {
	unLockErr := unLockError{
		key: key,
		lockService: lockService,
	}
	return &unLockErr
}

type redisMutexService struct {
	client redis.Client
}

func NewRedisMutexServiceFromOps(ops *redis.Option) LockService {
	redisClient := redis.NewRedisClientFromOps(ops)
	service := redisMutexService{
		client: redisClient,
	}
	return &service
}

func NewRedisMutexServiceFromClient(client redis.Client) LockService {
	service := redisMutexService{
		client: client,
	}
	return &service
}

func (r *redisMutexService)Lock(lockKey *LockKey) *Status  {
	ctx := context.Background()
	valStr := lockKey.Meta
	resp := r.client.SetNX(ctx, lockKey.Key, valStr, lockKey.TimeOut) //返回执行结果
	lockStatus := Status{}
	lockSuccess, err := resp.Result()
	lockErr := newLockError(lockKey)
	if lockSuccess {
		lockStatus.Success = lockSuccess
		lockStatus.Err = nil
		lockStatus.Meta = lockKey.Meta
		lockStatus.ExpireAt = time.Now().Add(lockKey.TimeOut)
	} else {
		getResp, getErr := r.client.Get(ctx, lockKey.Key).Result()
		timeOut := r.client.TTL(ctx, lockKey.Key).Val()
		if getErr != nil {
			lockStatus.Success = false
			lockStatus.Err = lockErr.WrapError(getErr).WrapMsg("get lock info from db failed")
		} else if valStr == getResp {
			lockStatus.Success = true
			lockStatus.Err = nil
			lockStatus.Meta = valStr
			lockStatus.ExpireAt = time.Now().Add(timeOut)
			if timeOut < 5 {
				r.client.Expire(ctx, lockKey.Key, lockKey.TimeOut - timeOut)
				lockStatus.ExpireAt = time.Now().Add(lockKey.TimeOut)
			}
		} else {
			lockStatus.Success = false
			lockStatus.Err = lockErr.WrapError(err).WrapMsg("lock failed")
			lockStatus.Meta = getResp
			lockStatus.ExpireAt = time.Now().Add(timeOut)
		}
	}
	return &lockStatus
}

func (r *redisMutexService)Unlock(lockKey *LockKey) (bool, error)  {
	ctx := context.Background()
	unLockErr := newUnLockError(lockKey.Key)
	getResp, err := r.client.Get(ctx, lockKey.Key).Result()
	if err != nil {
		return false, unLockErr.WrapError(err).WrapMsg("get lock info failed")
	}
	if getResp == lockKey.Meta {
		_, err := r.client.Del(ctx, lockKey.Key).Result()
		if err != nil {
			return false, unLockErr.WrapError(err).WrapMsg("release lock failed")
		}
		return true, nil
	}
	return false, unLockErr.WrapMsg("I don't have the lock")
}