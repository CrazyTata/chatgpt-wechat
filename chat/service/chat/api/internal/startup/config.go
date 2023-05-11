package startup

import (
	"chat/service/chat/api/internal/config"
	"encoding/json"
	"fmt"
	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
	conf = dealSliceConfig(v, conf)
	// 直接反序列化到结构体中
	err = v.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("Unmarshal 的conf.Redis: %+v", conf.Redis)
	fmt.Printf("Unmarshal 的conf.Mysql: %s", conf.Embeddings.Milvus.Keywords[0])

	return
}

func dealSliceConfig(v *viper.Viper, conf config.Config) config.Config {
	for _, s := range v.GetStringSlice("redis") {
		var rc []config.Redis
		json.Unmarshal([]byte(s), &rc)
		for _, rSon := range rc {
			conf.RedisCache = append(conf.RedisCache, cache.NodeConf{
				RedisConf: redis.RedisConf{
					Host: rSon.Host,
					Pass: rSon.Pass,
				},
			})
		}
	}
	return conf
}
