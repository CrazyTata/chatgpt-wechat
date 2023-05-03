package startup

import (
	"chat/service/chat/api/internal/config"
	"fmt"
	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
)

const ApolloAppId = "chatgpt-wechat"
const ApolloHost = "localhost:8080"
const ApolloNamespace = "application"

func LoadConfig() (conf config.Config, err error) {
	remote.SetAppID(ApolloAppId)
	v := viper.New()
	v.SetConfigType("prop") // 根据namespace实际格式设置对应type
	err = v.AddRemoteProvider("apollo", ApolloHost, ApolloNamespace)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v.AllSettings())
	// 直接反序列化到结构体中
	err = v.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Unmarshal 的conf.Redis: %+v", conf.Redis)
	fmt.Printf("Unmarshal 的conf.Mysql: %s", conf.Mysql.DataSource)

	return
}
