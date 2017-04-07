package main

type LoginGroup int8 // 登录组别

const (
	LOGIN_GROUP_D LoginGroup = 1 // 项目d
	LOGIN_GROUP_E LoginGroup = 2
	LOGIN_GROUP_F LoginGroup = 3
	LOGIN_GROUP_G LoginGroup = 4
)

// 玩家登录参数
type LoginIn struct {
	User  []byte     // 平台帐号，最大长度100字符
	Group LoginGroup // 组别
}

type LoginStatus int8 // 登录状态

const (
	LOGIN_STATUS_SUCCEED    LoginStatus = 1 // 登录成功
	LOGIN_STATUS_FIRST_TIME LoginStatus = 2 // 首次登录
)

// 玩家登录返回
type LoginOut struct {
	Status       LoginStatus // 登录返回结果
	PlayerId     int64       // 玩家在游戏中的ID
	LastDistance int64       // 上次距离
	MaxDistance  int64       // 最远的距离
}

type DistanceIn struct {
	Distance  int64 // 距离
	AliveTime int32 // 存活时间
	Golds     int32 // 金币数量
}

type DistanceOut struct{}

type RankPlayer struct {
	BestGroup     LoginGroup // 组别
	User          string     // 名字
	Num           int32      // 名次
	Distance      int64      // 距离
	BestAliveTime int32      // 存活时间
	BestGolds     int32      // 金币数量
}

type RankIn struct{}

type RankOut struct {
	List []RankPlayer
	Rank int32
}

//
// 通知关闭Session
//
type NotifyCloseSessionOut struct {
}
