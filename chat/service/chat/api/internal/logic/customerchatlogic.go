package logic

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"chat/common/milvus"
	"chat/common/openai"
	"chat/common/redis"
	"chat/common/wecom"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const EMBEDDING_MODEmbedding = "ARTICLE"

type CustomerChatLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	model      string
	baseHost   string
	basePrompt string
	message    string
}

func NewCustomerChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CustomerChatLogic {
	return &CustomerChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CustomerChatLogic) CustomerChat(req *types.CustomerChatReq) (resp *types.CustomerChatReply, err error) {

	l.setModelName().setBasePrompt().setBaseHost()

	// 确认消息没有被处理过
	_, err = l.svcCtx.ChatModel.FindOneByQuery(context.Background(),
		l.svcCtx.ChatModel.RowBuilder().Where(squirrel.Eq{"message_id": req.MsgID}).Where(squirrel.Eq{"user": req.CustomerID}),
	)
	// 消息已处理
	if err == nil {
		return &types.CustomerChatReply{
			Message: "ok",
		}, nil
	}

	// 指令匹配， 根据响应值判定是否需要去调用 openai 接口了
	proceed, _ := l.FactoryCommend(req)
	if !proceed {
		return
	}
	if l.message != "" {
		req.Msg = l.message
	}

	// openai client
	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).WithModel(l.model).WithBaseHost(l.baseHost)
	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	// context
	collection := openai.NewUserContext(
		openai.GetUserUniqueID(req.CustomerID, req.OpenKfID),
	).WithPrompt(l.basePrompt).WithClient(c)

	// 然后 把 消息 发给 openai
	go func() {
		// 去通过 embeddings 进行数据匹配
		type EmbeddingData struct {
			Q string `json:"q"`
			A string `json:"a"`
		}
		var embeddingData []EmbeddingData
		// 为了避免 embedding 的冷启动问题，对问题进行缓存来避免冷启动, 先简单处理
		matchEmbeddings := len(l.svcCtx.Config.Embeddings.Milvus.Keywords) == 0
		for _, keyword := range l.svcCtx.Config.Embeddings.Milvus.Keywords {
			if strings.Contains(req.Msg, keyword) {
				matchEmbeddings = true
			}
		}

		if l.svcCtx.Config.Embeddings.Enable && matchEmbeddings {
			milvusService, err := milvus.InitMilvus(l.svcCtx.Config.Embeddings.Milvus.Host, l.svcCtx.Config.Embeddings.Milvus.Username, l.svcCtx.Config.Embeddings.Milvus.Password)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer milvusService.CloseClient()
			// md5 this req.MSG to key
			key := md5.New()
			_, _ = io.WriteString(key, req.Msg)
			keyStr := fmt.Sprintf("%x", key.Sum(nil))
			type EmbeddingCache struct {
				Embedding []float64 `json:"embedding"`
			}
			embeddingRes, err := redis.Rdb.Get(context.Background(), fmt.Sprintf(redis.EmbeddingsCacheKey, keyStr)).Result()
			if err == nil {
				tmp := new(EmbeddingCache)
				_ = json.Unmarshal([]byte(embeddingRes), tmp)

				result := milvusService.SearchFromQA(tmp.Embedding)
				tempMessage := ""
				for _, qa := range result {
					if qa.Score > 0.3 {
						continue
					}
					if len(embeddingData) < 2 {
						embeddingData = append(embeddingData, EmbeddingData{
							Q: qa.Q,
							A: qa.A,
						})
					} else {
						tempMessage += qa.Q + "\n"
					}
				}
				if tempMessage != "" {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "正在思考中，也许您还想知道"+"\n\n"+tempMessage)
				}
			} else {
				go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "正在为您搜索相关数据")
				res, err := c.CreateOpenAIEmbeddings(req.Msg)
				if err == nil {
					fmt.Println(res.Data)
					fmt.Println(l.svcCtx.Config.Embeddings)
					embedding := res.Data[0].Embedding
					// 去将其存入 redis
					embeddingCache := EmbeddingCache{
						Embedding: embedding,
					}
					redisData, err := json.Marshal(embeddingCache)
					if err == nil {
						redis.Rdb.Set(context.Background(), fmt.Sprintf(redis.EmbeddingsCacheKey, keyStr), string(redisData), -1*time.Second)
					}
					// 将 embedding 数据与 milvus 数据库 内的数据做对比响应前3个相关联的数据
					result := milvusService.SearchFromQA(embedding)

					tempMessage := ""
					for _, qa := range result {
						if qa.Score > 0.3 {
							continue
						}
						if len(embeddingData) < 2 {
							embeddingData = append(embeddingData, EmbeddingData{
								Q: qa.Q,
								A: qa.A,
							})
						} else {
							tempMessage += qa.Q + "\n"
						}
					}
					if tempMessage != "" {
						go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "正在思考中，也许您还想知道"+"\n\n"+tempMessage)
					}
				}
			}
		}

		// 基于 summary 进行补充
		messageText := ""
		for _, chat := range embeddingData {
			collection.Set(chat.Q, chat.A, false)
		}
		collection.Set(req.Msg, "", false)

		prompts := collection.GetChatSummary()
		if l.svcCtx.Config.Response.Stream {
			channel := make(chan string, 100)
			go func() {
				messageText, err := c.WithModel(l.model).WithBaseHost(l.baseHost).ChatStream(prompts, channel)
				if err != nil {
					logx.Error("读取 stream 失败：", err.Error())
					wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统拥挤，稍后再试~"+err.Error())
					return
				}
				collection.Set("", messageText, true)
				// 再去插入数据
				_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
					User:       req.CustomerID,
					OpenKfId:   req.OpenKfID,
					MessageId:  req.MsgID,
					ReqContent: req.Msg,
					ResContent: messageText,
				})
			}()

			var rs []rune
			// 加快初次响应的时间 后续可改为阶梯式（用户体验好）
			first := true
			for {
				s, ok := <-channel
				if !ok {
					// 数据接受完成
					if len(rs) > 0 {
						go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, string(rs))
					}
					return
				}
				rs = append(rs, []rune(s)...)

				if first && len(rs) > 50 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
					first = false
				} else if len(rs) > 200 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
				}
			}
		}

		messageText, err := c.WithModel(l.model).WithBaseHost(l.baseHost).Chat(prompts)

		if err != nil {
			wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统错误:"+err.Error())
			return
		}

		// 然后把数据 发给对应的客户
		go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, messageText)
		collection.Set("", messageText, true)
		_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
			User:       req.CustomerID,
			OpenKfId:   req.OpenKfID,
			MessageId:  req.MsgID,
			ReqContent: req.Msg,
			ResContent: messageText,
		})
	}()

	return &types.CustomerChatReply{
		Message: "ok",
	}, nil
}

func (l *CustomerChatLogic) setModelName() (ls *CustomerChatLogic) {
	m := l.svcCtx.Config.WeCom.Model
	if m == "" || (m != openai.TextModel && m != openai.ChatModel && m != openai.ChatModelNew && m != openai.ChatModel4) {
		m = openai.TextModel
	}
	l.svcCtx.Config.WeCom.Model = m
	return l
}

func (l *CustomerChatLogic) setBasePrompt() (ls *CustomerChatLogic) {
	p := l.svcCtx.Config.WeCom.BasePrompt
	if p == "" {
		p = "你是 ChatGPT, 一个由 OpenAI 训练的大型语言模型, 你旨在回答并解决人们的任何问题，并且可以使用多种语言与人交流。\n"
	}
	l.basePrompt = p
	return l
}

func (l *CustomerChatLogic) setBaseHost() (ls *CustomerChatLogic) {
	if l.svcCtx.Config.OpenAi.Host == "" {
		l.svcCtx.Config.OpenAi.Host = "https://api.openai.com"
	}
	l.baseHost = l.svcCtx.Config.OpenAi.Host
	return l
}

func (l *CustomerChatLogic) FactoryCommend(req *types.CustomerChatReq) (proceed bool, err error) {
	template := make(map[string]CustomerTemplateData)
	//当 message 以 # 开头时，表示是特殊指令
	if !strings.HasPrefix(req.Msg, "#") {
		return true, nil
	}

	template["#voice"] = CustomerCommendVoice{}
	template["#clear"] = CustomerCommendClear{}

	for s, data := range template {
		if strings.HasPrefix(req.Msg, s) {
			proceed = data.customerExec(l, req)
			return proceed, nil
		}
	}

	return true, nil
}

type CustomerTemplateData interface {
	customerExec(svcCtx *CustomerChatLogic, req *types.CustomerChatReq) (proceed bool)
}

type CustomerCommendVoice struct{}

func (p CustomerCommendVoice) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	msg := strings.Replace(req.Msg, "#voice:", "", -1)
	if msg == "" {
		wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统错误:未能读取到音频信息")
		return false
	}

	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key)

	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	txt, err := c.SpeakToTxt(msg)

	if txt == "" || err != nil {
		logx.Info("openai转换错误", err.Error())
		wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统错误:音频信息转换错误")
		return false
	}
	// 语音识别成功
	//wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "语音识别成功:\n\n"+txt+"\n\n系统正在思考中...")

	l.message = txt
	return true
}

type CustomerCommendClear struct{}

func (p CustomerCommendClear) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	// 清理上下文
	openai.NewUserContext(
		openai.GetUserUniqueID(req.CustomerID, req.OpenKfID),
	).Clear()
	wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "记忆清除完成:来开始新一轮的chat吧")
	return false
}

func (l *CustomerChatLogic) CustomerChatV2(req *types.CustomerChatReq) (resp *types.CustomerChatReply, err error) {

	l.setModelName().setBasePrompt().setBaseHost()

	// 确认消息没有被处理过
	_, err = l.svcCtx.ChatModel.FindOneByQuery(context.Background(),
		l.svcCtx.ChatModel.RowBuilder().Where(squirrel.Eq{"message_id": req.MsgID}).Where(squirrel.Eq{"user": req.CustomerID}),
	)
	// 消息已处理
	if err == nil {
		return &types.CustomerChatReply{
			Message: "ok",
		}, nil
	}

	// 指令匹配， 根据响应值判定是否需要去调用 openai 接口了
	proceed, _ := l.FactoryCommend(req)
	if !proceed {
		return
	}
	if l.message != "" {
		req.Msg = l.message
	}

	// openai client
	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).WithModel(l.model).WithBaseHost(l.baseHost)
	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	//get prompt
	promptPo, err := l.svcCtx.CustomerPromptModel.FindOneByQuery(context.Background(),
		l.svcCtx.CustomerPromptModel.RowBuilder().Where(squirrel.Eq{"kf_id": req.OpenKfID}),
	)
	fmt.Println(promptPo.Prompt)

	if err != nil {
		return nil, err
	}
	if promptPo != nil && promptPo.Prompt != "" {
		l.basePrompt = promptPo.Prompt
	}
	// context
	collection := openai.NewUserContext(
		openai.GetUserUniqueID(req.CustomerID, req.OpenKfID),
	).WithPrompt(l.basePrompt).WithClient(c)
	embeddingEnable := true
	embeddingMode := EMBEDDING_MODEmbedding

	// 然后 把 消息 发给 openai
	go func() {
		// 基于 summary 进行补充
		messageText := ""
		if embeddingEnable {

			// md5 this req.MSG to key
			key := md5.New()
			_, _ = io.WriteString(key, req.Msg)
			keyStr := fmt.Sprintf("%x", key.Sum(nil))
			type EmbeddingCache struct {
				Embedding []float64 `json:"embedding"`
			}
			milvusService, err := milvus.InitMilvus(l.svcCtx.Config.Embeddings.Milvus.Host, l.svcCtx.Config.Embeddings.Milvus.Username, l.svcCtx.Config.Embeddings.Milvus.Password)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer milvusService.CloseClient()
			embeddingRes, err := redis.Rdb.Get(context.Background(), fmt.Sprintf(redis.EmbeddingsCacheKey, keyStr)).Result()
			var embedding []float64
			if err == nil {
				tmp := new(EmbeddingCache)
				_ = json.Unmarshal([]byte(embeddingRes), tmp)
				embedding = tmp.Embedding
			} else {
				res, err := c.CreateOpenAIEmbeddings(req.Msg)
				if err == nil {
					embedding = res.Data[0].Embedding
					// 去将其存入 redis
					embeddingCache := EmbeddingCache{
						Embedding: embedding,
					}
					redisData, err := json.Marshal(embeddingCache)
					if err == nil {
						redis.Rdb.Set(context.Background(), fmt.Sprintf(redis.EmbeddingsCacheKey, keyStr), string(redisData), -1*time.Second)
					}
				}
			}

			if embeddingMode == "QA" {
				// 去通过 embeddings 进行数据匹配
				type EmbeddingData struct {
					Q string `json:"q"`
					A string `json:"a"`
				}
				var embeddingData []EmbeddingData
				result := milvusService.SearchFromQA(embedding)
				tempMessage := ""
				for _, qa := range result {
					if qa.Score > 0.3 {
						continue
					}
					if len(embeddingData) < 2 {
						embeddingData = append(embeddingData, EmbeddingData{
							Q: qa.Q,
							A: qa.A,
						})
					} else {
						tempMessage += qa.Q + "\n"
					}
				}
				if tempMessage != "" {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "正在思考中，也许您还想知道"+"\n\n"+tempMessage)
				}
				for _, chat := range embeddingData {
					collection.Set(chat.Q, chat.A, false)
				}
				collection.Set(req.Msg, "", false)
			} else if embeddingMode == "ARTICLE" {
				//如果是article模式，清理掉上下午，因为文章内容可能会很长
				collection.Clear()
				collection.Messages = []openai.ChatModelMessage{}
				collection.Summary = []openai.ChatModelMessage{}
				// 去通过 embeddings 进行数据匹配
				type EmbeddingData struct {
					text string `json:"text"`
				}
				var embeddingData []EmbeddingData
				result := milvusService.SearchFromArticle(embedding)
				for _, qa := range result {

					fmt.Println("text:", qa.CnText)
					fmt.Println("score:", qa.Score)
					if len(embeddingData) < 1 {
						embeddingData = append(embeddingData, EmbeddingData{
							text: qa.CnText,
						})
						go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统语料:"+qa.CnText)

					}
				}

				if len(embeddingData) > 0 {
					messageText += "When given CONTEXT you answer questions using only that information,and you always format your output in markdown.Answer with chinese.\n\n"
					messageText += "CONTEXT:"
					for _, chat := range embeddingData {
						messageText += chat.text + "\n\n"
					}
					messageText += "USER QUESTION:" + req.Msg
					collection.Set(messageText, "", false)
				} else {
					collection.Set(req.Msg, "", false)
				}
			} else {
				collection.Set(req.Msg, "", false)
			}
		}

		collection.Set(req.Msg, "", false)

		prompts := collection.GetChatSummary()
		if l.svcCtx.Config.Response.Stream {
			channel := make(chan string, 100)
			go func() {
				messageText, err := c.WithModel(l.model).WithBaseHost(l.baseHost).ChatStream(prompts, channel)
				if err != nil {
					logx.Error("读取 stream 失败：", err.Error())
					wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统拥挤，稍后再试~"+err.Error())
					return
				}
				collection.Set("", messageText, true)
				// 再去插入数据
				_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
					User:       req.CustomerID,
					OpenKfId:   req.OpenKfID,
					MessageId:  req.MsgID,
					ReqContent: req.Msg,
					ResContent: messageText,
				})
			}()

			var rs []rune
			// 加快初次响应的时间 后续可改为阶梯式（用户体验好）
			first := true
			for {
				s, ok := <-channel
				if !ok {
					// 数据接受完成
					if len(rs) > 0 {
						go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, string(rs))
					}
					return
				}
				rs = append(rs, []rune(s)...)

				if first && len(rs) > 50 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
					first = false
				} else if len(rs) > 200 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
				}
			}
		}

		messageText, err = c.WithModel(l.model).WithBaseHost(l.baseHost).Chat(prompts)

		if err != nil {
			wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, "系统错误:"+err.Error())
			return
		}

		// 然后把数据 发给对应的客户
		go wecom.SendCustomerChatMessage(req.OpenKfID, req.CustomerID, messageText)
		collection.Set("", messageText, true)
		_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
			User:       req.CustomerID,
			OpenKfId:   req.OpenKfID,
			MessageId:  req.MsgID,
			ReqContent: req.Msg,
			ResContent: messageText,
		})
	}()

	return &types.CustomerChatReply{
		Message: "ok",
	}, nil
}
