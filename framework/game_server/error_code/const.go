package errcode

const (
	TimeLayout     = "15:04:05"
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
)

// bacErrorMap 错误码与错误字符串映射信息
var bacErrorMap map[int]string

// GetErrMsg 获取错误码对应的错误字符串
func GetErrMsg(errCode int) string {
	strMsg := bacErrorMap[ErrorUnknown]
	if strErr, ok := bacErrorMap[errCode]; ok {
		strMsg = strErr
	}

	return strMsg
}

const (
	MaxLoopCount = 1000 //查询UniqueId最大尝试次数

	ErrorInvalid     = -1   //无效小数值
	ErrorOk          = 0    //统一为没有错误
	ErrorSettleEmpty = 1    //结算注单为零
	ErrorUnknown     = 7000 //统一为内部错误
)

/*
	beego的HTTP Status限定在[511-999]
	下注接口返回给能力中心的错误限定在[7000-8999] 这些错误码已经融入到能力平台中
	其中[7001,7999]为能力中心 游戏服 客户端共同使用 一旦定义则不可改动
	[8000,8999]为游戏自己使用
*/

const (
	/* limit validate 相关错误*/

	/* begin 以下数据数值不可改动 已经与能力中心以及客户端协调好 预设置范围为[7001,7999] */

	ValidateErrorFailed               = iota + 7001 //统一校验错误
	ValidateErrorUserLimitNotExist                  //个人限红不存在
	ValidateErrorUserLimitGetFailed                 //获取个人限红失败
	ValidateErrorUserLimitMoreThanMax               //下注超过个人最大限红
	ValidateErrorUserLimitLessThanMin               //下注小于个人最小限红

	ValidateErrorRoomLimitGetFailed           //获取房间限红失败
	ValidateErrorRoomLimitMoreThanMax         //下注超过房间最大限红
	ValidateErrorRoomLimitLessThanMin         //下注小于房间最小限红
	ValidateErrorRoomLimitTotalBetMoreThanMax //超过房间总限红

	ValidateErrorNoMainPlayTypeBet      //主玩法没有下注
	ValidateErrorMoreThanMainPlayType50 //旁注玩法超过主玩法下注金额的50%
	ValidateErrorMoreThanMainPlayType20 //旁注玩法超过主玩法下注金额的20%

	/* end 以上数据数值不可改动 已经与能力中心以及客户端协调好 */

	ValidateErrorResultParseFailed //result 解析错误

)

const DBErrorOK = 0

/* 数据库相关错误 [8000, 8019]*/
const (
	DBErrorDuplicate = iota + 8000
	DBErrorNotOk     //非具体错误使用该错误码
)

const HttpStatusOK = 200 //Http正常返回码
/* 能力平台交互的Http相关错误 [8020, 8039]*/
const (
	HttpErrorOrderParse                = iota + 8020 //订单解析错误
	HttpErrorDataFailed                              //数据异常
	HttpErrorInvalidParam                            //无效的参数
	HttpErrorPlatformReply                           //平台中心返回错误码
	HttpErrorServerReply                             //服务器返回错误
	HttpErrorPlatformPost                            //向平台中心发送信息错误
	HttpErrorPlatFormBuildWorkerFailed               //向平台中心发送创建员工信息返回失败
)

/* redis相关错误 [8040, 8059]*/
const (
	RedisErrorSet            = iota + 8040 //redis set error
	RedisErrorGet                          //redis get error
	RedisErrorDelete                       //redis delete error
	RedisErrorLock                         //redis lock error
	RedisErrorNoCallbackFunc               //redis no callback function
	RedisErrorDataIsEmpty                  //redis data is empty error
	RedisErrorTTLNotValid                  //redis ttl not valid
)

/* json marshal unmarshal相关错误  [8060, 8069] */
const (
	JsonErrorMarshal = iota + 8060
	JsonErrorUnMarshal
)

/* 游戏相关错误 [8070, 8089] */

const (
	GameErrorWrongGameRoundStatus    = iota + 8070 //游戏状态错误
	GameErrorBetCancelFailed                       //取消订单失败
	GameErrorBetCancelTooFast                      //取消订单太快
	GameErrorBetTooFast                            //投注太快
	GameErrorGameRoundIdNotExist                   //局号不存在
	GameErrorUserIdNotExist                        //用户不存在
	GameErrorBalanceNotEnough                      //用户余额不足
	GameErrorBetParamIllegal                       //下注参数不合法
	GameErrorBetConfirmLater                       //投注确认慢
	GameErrorNoGameDBSaverRegistered               //数据库保存接口未注册
	GameErrorWrongGameId                           //错误的游戏Id
	GameErrorBettorNotExist                        //下注对象为注册
	GameErrorGameEventExist                        //游戏事件已存在

)

func init() {
	bacErrorMap = make(map[int]string, 32)
	bacErrorMap[ErrorOk] = "success"

	//初始化Http错误码字符串
	bacErrorMap[HttpErrorOrderParse] = "order parse to json error" //订单解析错误
	bacErrorMap[HttpErrorDataFailed] = "data abnormal"             //数据异常
	bacErrorMap[HttpErrorInvalidParam] = "invalid param"           //无效参数
	bacErrorMap[HttpErrorPlatformReply] = "platform reply error"   //平台中心返回错误
	bacErrorMap[HttpErrorServerReply] = "server reply error"       //服务器返回错误
	bacErrorMap[HttpErrorPlatformPost] = "post to platform error"  //向平台中心发送信息错误
	bacErrorMap[HttpErrorPlatFormBuildWorkerFailed] = "build worker failed"

	/* json marshal unmarshal相关错误*/
	bacErrorMap[JsonErrorMarshal] = "json data marshal error"
	bacErrorMap[JsonErrorUnMarshal] = "json data unmarshal error"

	//Redis相关错误
	bacErrorMap[RedisErrorSet] = "redis set error"
	bacErrorMap[RedisErrorGet] = "redis get error"
	bacErrorMap[RedisErrorDelete] = "redis delete error"
	bacErrorMap[RedisErrorLock] = "redis lock error"
	bacErrorMap[RedisErrorNoCallbackFunc] = "redis no callback function"
	bacErrorMap[RedisErrorDataIsEmpty] = "redis data is empty"
	bacErrorMap[RedisErrorTTLNotValid] = "redis ttl not valid"

	//下注校验相关错误
	bacErrorMap[ValidateErrorFailed] = "bet failed"                                                //统一校验错误
	bacErrorMap[ValidateErrorUserLimitNotExist] = "user limit not exist"                           //个人限红不存在
	bacErrorMap[ValidateErrorUserLimitGetFailed] = "user limit get failed"                         //获取个人限红失败
	bacErrorMap[ValidateErrorUserLimitMoreThanMax] = "user limit more than max"                    //下注超过个人最大限红
	bacErrorMap[ValidateErrorUserLimitLessThanMin] = "user limit less than min"                    //下注小于个人最小限红
	bacErrorMap[ValidateErrorRoomLimitGetFailed] = "room limit get failed"                         //获取房间限红失败
	bacErrorMap[ValidateErrorRoomLimitMoreThanMax] = "room limit more than max"                    //下注超过房间最大限红
	bacErrorMap[ValidateErrorRoomLimitLessThanMin] = "room limit less than min"                    //下注小于房间最小限红
	bacErrorMap[ValidateErrorRoomLimitTotalBetMoreThanMax] = "room total limit more than min"      //超过房间总限红
	bacErrorMap[ValidateErrorNoMainPlayTypeBet] = "no main bet"                                    //主玩法没有下注
	bacErrorMap[ValidateErrorMoreThanMainPlayType50] = "side bet more than 50 percent of main bet" //旁注玩法超过主玩法下注金额的50%
	bacErrorMap[ValidateErrorMoreThanMainPlayType20] = "side bet more than 20 percent of main bet" //旁注玩法超过主玩法下注金额的20%

	//游戏相关错误
	bacErrorMap[GameErrorWrongGameRoundStatus] = "wrong game round status"         //游戏状态错误
	bacErrorMap[GameErrorBetCancelFailed] = "bet cancel failed"                    //取消订单失败
	bacErrorMap[GameErrorBetCancelTooFast] = "bet cancel too fast"                 //取消订单太快
	bacErrorMap[GameErrorBetTooFast] = "bet too fast"                              //投注太快
	bacErrorMap[GameErrorGameRoundIdNotExist] = "game round id not exist in redis" //缓存中没有局信息
	bacErrorMap[GameErrorUserIdNotExist] = "user id not exist in redis"            //缓存中没有用户信息
	bacErrorMap[GameErrorBalanceNotEnough] = "balance not enough"                  //用户余额不足
	bacErrorMap[GameErrorBetParamIllegal] = "bet param illegal"                    //下注参数无效
	bacErrorMap[GameErrorBetConfirmLater] = "bet confirm later"                    //投注确认慢
	bacErrorMap[GameErrorNoGameDBSaverRegistered] = "no game db saver registered"  //数据库保存接口未注册
	bacErrorMap[GameErrorWrongGameId] = "wrong game id"                            //错误的游戏Id
	bacErrorMap[GameErrorBettorNotExist] = "bettor not exist"                      //下注对象不存在
	bacErrorMap[GameErrorGameEventExist] = "The game event exist"                  //游戏事件已存在

	//结算相关错误
	bacErrorMap[ValidateErrorResultParseFailed] = "result parse failed" //result 解析错误

	bacErrorMap[ErrorUnknown] = "game server error"
}
