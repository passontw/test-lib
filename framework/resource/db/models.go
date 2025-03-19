package db

import (
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	// Register the models with Beego ORM
	orm.RegisterModel(
		new(VideoInfo),
		new(SubscriberInfo),
		new(HttpPostRequests),
		new(ApiInfo),
		new(GmcodeMapping),
	)
}
