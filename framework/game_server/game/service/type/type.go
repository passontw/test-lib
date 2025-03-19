package types

import "strconv"

type GameId int64

/**
 * ToString
 * GameId转化为字符串
 */

func (id GameId) ToString() string {
	return strconv.FormatInt(int64(id), 10)
}

// PlayTypeId 玩法类型
type PlayTypeId int64

// ToString PlayTypeId转string类型
func (p PlayTypeId) ToString() string {
	return strconv.FormatInt(int64(p), 10)
}
