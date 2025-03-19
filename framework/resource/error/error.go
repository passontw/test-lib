package err

const (
	ERR_OK                              = iota // 0 // 成功
	ERR_INVL_PARAM                             // 1 // 无效参数
	ERR_CS_DISCONN                             // 2 // 与中心服务器连接中断
	ERR_DS_DISCONN                             // 3 // 与数据源连接中断
	ERR_DB                                     // 4 // 数据库 I/O 失败
	ERR_DB_TIMEOUT                             // 5 // 数据库操作超时
	ERR_DBC_UNCONN                             // 6 // 数据库服务器未连接
	ERR_DBC_TIMEOUT                            // 7 // 数据库服务器超时
	ERR_DBC_FAIL_READ                          // 8 // 数据库服务器读失败
	ERR_DBC_FAIL_SEND                          // 9 // 数据库服务器发失败
	ERR_SQL_ERROR                              // 10 // SQL 发生异常
	ERR_NO_DBRECORD                            // 11 // 不存在数据库记录
	ERR_LOW_VER                                // 12 // 登录时，版本比服务器低
	ERR_PWD_ERROR                              // 13 // 密码错误
	ERR_NO_USER                                // 14 // 用户不存在
	ERR_TIMEOUT                                // 15 // 超时
	ERR_NO_LOGIN                               // 16 // 用户未登录
	ERR_INVL_USER_TYPE                         // 17 // 用户类型无效
	ERR_NOT_ON_TABLE                           // 18 // 未在台面就下注
	ERR_JETTON_SLIMIT                          // 19 // 超出位置限红
	ERR_JETTON_TLIMIT                          // 20 // 超出台面限红
	ERR_JETTON_TLIMIT_POT                      // 21 // 超出彩池总限红
	ERR_PERSONAL_LIMIT                         // 22 // 超个人盘口限额
	ERR_USER_FUNC_LIMIT                        // 23 // 试玩/真钱用户功能被限制
	ERR_USER_BET_LIMITED                       // 24 // 用户下注受限制
	ERR_NO_PLAYTYPE                            // 25 // 没有该玩法
	ERR_INVL_PLAYTYPE                          // 26 // 玩法无效
	ERR_NO_TABLE                               // 27 // 没有该台桌
	ERR_GMCODE_EXPIRED                         // 28 // 期号过期
	ERR_INVL_ROUND                             // 29 // 无效期号
	ERR_NO_GMCODE                              // 30 // 没有该期
	ERR_INVL_ORDER                             // 31 // 无效注单
	ERR_SEAT_OCCUPIED                          // 32 // 位置被占
	ERR_INVL_SEATNUM                           // 33 // 无效座位
	ERR_ON_SEAT                                // 34 // 已经坐在位置
	ERR_NOT_ON_SEAT                            // 35 // 不在座位上
	ERR_NO_FREE_SEAT                           // 36 // 没有空桌位
	ERR_INVL_GAME_STATUS                       // 37 // 游戏状态无效
	ERR_NO_SHOECODE                            // 38 // 该桌未起牌靴
	ERR_NOT_BET_TIME                           // 39 // 已过下注时间
	ERR_LESS_AMOUNT                            // 40 // 额度不够
	ERR_NO_DEALER                              // 41 // 没有荷官
	ERR_INVL_CARD_VAL                          // 42 // 无效的牌
	ERR_INVL_TABLEINFO                         // 43 // 桌子信息未正确从数据库加载
	ERR_ON_GAME                                // 44 // 游戏进行中，不能退出
	ERR_TABLE_NOTEMPTY                         // 45 // 该桌已经有人，不允许包桌
	ERR_TABLE_CONTRACTED                       // 46 // 该桌已经被包
	ERR_TOO_MANY_CARDS                         // 47 // 牌数过多
	ERR_INVL_3TH_CARD                          // 48 // 第 3 张发错
	ERR_EXIST_TABLE                            // 49 // 加新桌，但桌号已存在
	ERR_INVL_VIDEOADDR                         // 50 // 视频地址不正确
	ERR_TABLE_CLOSED                           // 51 // 该桌已经被关闭
	ERR_INCOMP_TABLESTATE                      // 52 // 桌状态不一致，不能共享一个荷官
	ERR_ROUND_ORDER_NO_CHANGED                 // 53 // 取消局单子，无单子取消
	ERR_UNKNOWN_GAMETYPE                       // 54 // 未支持游戏类别
	ERR_NO_BULLETIN                            // 55 // 没有该公告
	ERR_EXIST_VIDEO                            // 56 // 已存在视频
	ERR_USER_LIMITED                           // 57 // 用户受限制
	ERR_NO_GAMEVIDEO                           // 58 // 无视频对象
	ERR_INVL_VID                               // 59 // 无效视频
	ERR_NOPOWER                                // 60 // 没有权限视频包桌
	ERR_NO_CONTRACTED                          // 61 // 视频未被包桌，无法开局
	ERR_VIDEO_CONTRACTED                       // 62 // 该视频已经被包
	ERR_VIDEO_LOCKED                           // 63 // 该视频已经被锁定
	ERR_UNINITIALIZED                          // 64 // VIP 玩家未开局
	ERR_NONE_SHARED                            // 65 // 视频独占模式
	ERR_VIDEO_PRIVATE                          // 66 // 视频已独享
	ERR_THROWING_LIMIT                         // 67 // 飞牌次数受限
	ERR_CHANGE_SHOE_LIMIT                      // 68 // 换靴局数受限
	ERR_CHANGE_DEALER_LIMIT                    // 69 // 换荷官次数受限
	ERR_UNREACHABLE_LIMIT                      // 70 // 未达到最低额度
	ERR_NOT_SUBSCRIBED                         // 71 // 未订阅该视频
	ERR_INVL_CONTRACT_MODE                     // 72 // 未知包桌模式
	ERR_VIDEO_IN_BET                           // 73 // 下注中，不接收停止下注指令
	ERR_INVL_VIDEO_STATUS                      // 74 // 当前视频状态不支持此操作
	ERR_INVL_VIDEO_MODE                        // 75 // 当前包桌模式不支持此操作
	ERR_NOT_MATCH_CONDITION                    // 76 // 未满足包桌条件
	ERR_VIDEO_RESERVED                         // 77 // 视频已经被预约了
	ERR_INVL_COMMAND                           // 78 // 无效命令号
	ERR_LIMIT_NOT_MATCH                        // 79 // 未匹配到合适的限红
	ERR_PLAYERS_LIMIT                          // 80 // 视频旁注人数超限
	ERR_TABLE_RESERVED                         // 81 // 桌子被旁注预订
	ERR_INVL_PLATFORM_ID                       // 82 // 无效游戏平台 ID
	ERR_PLATFORM_NOT_EMPTY                     // 83 // 平台不为空
	ERR_BANKER_PLAYER_PAIR_JETTON_LIMIT        // 84 // 庄对闲对下注总和超出限红上线
	ERR_NO_JACKPOT                             // 85 // 无奖池
	ERR_NOT_LED_VIDEO                          // 90 // 视频不是 LED 桌
	ERR_VIP_ENTER_LIMITED                      // 91 // 试玩账户不能包桌和入座下注，选择旁观下注
	ERR_LED_ENTER_LIMITED                      // 92 // 试玩账户不能入座下注，选择旁观下注
	ERR_EXIST_INS_BANKER_BET                   // 101 // 本局已投注过庄保险玩法
	ERR_EXIST_INS_PLAYER_BET                   // 102 // 本局已投注过闲保险玩法
	ERR_LOW_CREDIT_LIMITED                     // 103 // 竞眯百家乐下注超过 50 万需要 5 倍额度
	ERR_RECORD_EXIST                           // 120 // 数据库表里存在相同记录
	ERR_NO_ZHUBO_SEAT                          // 121 // 无主播座位
	ERR_ANCHOR_NOT_CURTABLE                    // 122 // 主播不是当前桌台的主播
	ERR_ANCHOT_NOTONLINE                       // 123 // 主播不上播，不允许主播账号进入游戏
	ERR_INVL_CURRENCY                          // 130 // 暂不支持的货币类型
	ERR_INVL_CLIENTVER                         // 140 // 客户端版本过低
	ERR_USER_DISABLE                           // 150 // 账户已禁用
	ERR_AS_DATA                         = 9001 // 9001 // 活动数据错误
)
