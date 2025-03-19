package gameserver

import (
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"os/signal"
	"sl.framework.com/async"
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/controller"
	"sl.framework.com/game_server/game/dao"
	"sl.framework.com/game_server/game/dao/redisdb"
	"sl.framework.com/game_server/game/dao/uiddb"
	"sl.framework.com/game_server/game/filter"
	"sl.framework.com/game_server/game/filter/common"
	"sl.framework.com/game_server/game/service"
	"sl.framework.com/game_server/game/service/base"
	"sl.framework.com/game_server/mq"
	"sl.framework.com/trace"
	"syscall"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

// 在init()函数中初始化日志模块 三方库使用方引入模块时候初始化
func init() {
	//初始化log 内部使用sync.Once保证只初始化一次
	trace.LoggerInit()
}

// beegoWebInit初始化并异步启动web
func beegoWebInit() {
	async.AsyncRunCoroutine(func() {
		//通用配置开启
		beego.BConfig.CopyRequestBody = true                                    //必须设置为true 否则http拿不到数据
		beego.BConfig.Listen.EnableAdmin = conf.ServerConf.BeegoCFG.AdminEnable //默认不启动8088端口 8088端口手动独立启动
		beego.BConfig.WebConfig.AutoRender = false                              //没有前端 设置为false
		beego.BConfig.RunMode = conf.ServerConf.BeegoCFG.RunMode                //显式设置开发模式为dev 打印路由耗时
		beego.BConfig.Listen.Graceful = conf.ServerConf.BeegoCFG.GraceEnable    //优雅关闭
		beego.BConfig.EnableErrorsRender = false                                //关闭默认错误页面
		beego.BConfig.RecoverPanic = true

		beego.BConfig.RecoverFunc = filter.RecoverPanic //重写 panic错误处理函数

		controller.RegisterRouter()
		//beego.InsertFilter("/*", beego.BeforeRouter, common.BeforeRouterCommonHandler)
		beego.InsertFilter("/*", beego.FinishRouter, common.ErrorHandler)
		//beego.InsertFilter("/*", beego.AfterExec, common.AfterExecCommonHandler)
		//服务端口 8080
		beego.Run(":8080")
	})

	async.AsyncRunCoroutine(func() {
		httpServer := beego.NewHttpSever()
		controller.RegisterHealthRouter(httpServer)

		//健康检查端口 8088
		httpServer.Run(":8088")
	})
}

/**
 * beegoStop
 * beego关闭时调用beego优雅关闭功能
 *
 * @param
 * @return
 */

func beegoStop() {
	// 设置超时时间为5秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//调用beego自带优雅关闭功能
	if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
		trace.Error("beegoStop Server forced to shutdown failed, err=%v", err.Error())
		return
	}
	trace.Info("beegoStop Server forced to shutdown success")
}

/*
	recoveryOnReboot 重启的时候恢复数据
*/

func recoveryOnReboot() {
}

/*
	FrameworkDataInit 服务数据初始化
*/

func FrameworkDataInit() {

}

// appInit 程序处理话 初始化结果以bool变量返回
func appInit() bool {
	//初始化Nocas读取配置
	if ok := conf.NacosClientInitOnce(conf.OnServerConf); !ok {
		trace.Error("appInit NacosClientInitOnce failed")
		return false
	}
	//读取到nacos配置以后设置日志打印级别
	trace.SetLevel(conf.GetLogLevel())

	// 初始化Redis
	if ok := redisdb.RedisClientInitOnce(); !ok {
		trace.Error("appInit RedisClientInitOnce failed")
		return false
	}

	//Database Orm初始化
	if err := dao.OrmInit(); nil != err {
		trace.Error("appInit OrmInit failed, error=%v", err.Error())
		return false
	}

	//初始化rocketmq
	if ok := mq.InitRocketManager(); !ok {
		trace.Error("appInit RocketMQInit failed")
		return false
	}

	//数据初始化
	FrameworkDataInit()

	//重启的时候恢复数据
	recoveryOnReboot()

	//启动内存管理任务
	base.GetCacheManager().RunLoopTask()

	//初始化uid db并设置server id
	uiddb.GetUniqueIdGeneratorInstance().SetUniqueId()

	//初始化api并启动监听端口
	beegoWebInit()

	return true
}

// Start 启动游戏服 供外部调用
func Start() {
	defer async.TryException() //保证运行时panic有日志输出
	service.ValidateService()

	//初始化log 内部使用sync.Once保证只初始化一次
	_ = trace.LoggerInit()
	//初始化App
	if !appInit() {
		trace.Error("main appInit error")
		return
	}
	trace.Info("main appInit success")

	waitOnSignal()
}

/**
 * waitOnSignal
 * 等待系统信号
 *
 * @param
 * @return
 */

func waitOnSignal() {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	for {
		select {
		case s := <-chSignal:
			trace.Notice("waitOnSignal receive os.signal:%v", s)
			graceStop()
			os.Exit(0)
		}
	}
}

/**
 * graceStop
 * 回收资源 实现程序优雅关闭
 *
 * @param
 * @return
 */

func graceStop() {
	trace.Info("graceStop start, current time=%v", time.Now().Format(timeLayout))

	//等待下注和结算任务处理完毕再关闭数据库等相关服务
	time.Sleep(time.Second * 5)

	//关闭beego 停止接收下注请求
	beegoStop()

	//关闭mq 停止接收结算消息
	mq.StopRocketManager()

	//关闭数据库
	/*
		beego orm 并不提供显式的关闭数据库连接的方法，通常依赖于数据库驱动的连接池管理
	*/

	//关闭缓存管理类

	//断开redis
	redisdb.RedisClientClose()
	trace.Info("graceStop done, current time=%v", time.Now().Format(timeLayout))
}
