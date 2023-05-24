package logic

import (
	"chat/common/util"
	"chat/common/wecom"
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
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
)

const EMBEDDING_MODEMBEDDING = "ARTICLE"
const MaxToken = 4000

type CustomerChatLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	model      string
	baseHost   string
	basePrompt string
	message    string
	postModel  string
}

func NewCustomerChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CustomerChatLogic {
	return &CustomerChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
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

	template["#direct"] = CustomerCommendDirect{}
	template["#voice"] = CustomerCommendVoice{}
	template["#help"] = CustomerCommendHelp{}
	template["#system"] = CustomerCommendSystem{}
	template["#clear"] = CustomerCommendClear{}
	template["#about"] = CustomerCommendAbout{}

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
		sendToUser(req.OpenKfID, "", req.CustomerID, "系统错误:未能读取到音频信息", l.svcCtx.Config)
		return false
	}

	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).
		WithModel(l.model).
		WithBaseHost(l.svcCtx.Config.OpenAi.Host).
		WithOrigin(l.svcCtx.Config.OpenAi.Origin).
		WithEngine(l.svcCtx.Config.OpenAi.Engine).
		WithPostModel(l.postModel)

	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	var cli openai.Speaker
	if l.svcCtx.Config.Speaker.Company == "" {
		l.svcCtx.Config.Speaker.Company = "openai"
	}
	switch l.svcCtx.Config.Speaker.Company {
	case "openai":
		logx.Info("使用openai音频转换")
		cli = c
	case "ali":
		logx.Info("使用阿里云音频转换")
		//s, err := voice.NewSpeakClient(
		//	l.svcCtx.Config.Speaker.AliYun.AccessKeyId,
		//	l.svcCtx.Config.Speaker.AliYun.AccessKeySecret,
		//	l.svcCtx.Config.Speaker.AliYun.AppKey,
		//)
		//if err != nil {
		// sendToUser(req.OpenKfID, "", req.CustomerID, "阿里云音频转换初始化失败:"+err.Error(), l.svcCtx.Config)
		//	return false
		//}
		//msg = strings.Replace(msg, ".mp3", ".amr", -1)
		//cli = s
	default:
		sendToUser(req.OpenKfID, "", req.CustomerID, "系统错误:未知的音频转换服务商", l.svcCtx.Config)
		return false
	}

	txt, err := cli.SpeakToTxt(msg)
	if txt == "" || err != nil {
		logx.Info("openai转换错误", err.Error())
		sendToUser(req.OpenKfID, "", req.CustomerID, "系统错误:音频信息转换错误", l.svcCtx.Config)
		return false
	}
	l.message = txt
	return true
}

type CustomerCommendClear struct{}

func (p CustomerCommendClear) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	// 清理上下文
	openai.NewUserContext(
		openai.GetUserUniqueID(req.CustomerID, req.OpenKfID),
	).Clear()
	sendToUser(req.OpenKfID, "", req.CustomerID, "记忆清除完成:来开始新一轮的chat吧", l.svcCtx.Config)
	return false
}

func (l *CustomerChatLogic) CustomerChat(req *types.CustomerChatReq) (resp *types.CustomerChatReply, err error) {

	l.setModelNameV2().setBasePrompt().setBaseHost()

	// 确认消息没有被处理过
	originInfo, err := l.svcCtx.ChatModel.FindOneByQuery(context.Background(),
		l.svcCtx.ChatModel.RowBuilder().Where(squirrel.Eq{"message_id": req.MsgID}).Where(squirrel.Eq{"user": req.CustomerID}),
	)
	// 消息已处理
	if err != nil {
		return &types.CustomerChatReply{}, err
	}

	if originInfo != nil && originInfo.Id > 0 {
		return &types.CustomerChatReply{
			Message: "ok",
		}, nil
	}

	embeddingEnable := true
	embeddingMode := EMBEDDING_MODEMBEDDING
	var baseScore float32
	var baseTopK int
	var clearContextTime int64
	var postModel string
	//get prompt
	customerConfigPo, err := l.svcCtx.CustomerConfigModel.FindOneByQuery(context.Background(),
		l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_id": req.OpenKfID}),
	)
	if err != nil {
		return nil, err
	}
	if customerConfigPo != nil {
		if customerConfigPo.Prompt != "" {
			l.basePrompt = customerConfigPo.Prompt
		}
		embeddingEnable = customerConfigPo.EmbeddingEnable
		embeddingMode = customerConfigPo.EmbeddingMode
		if customerConfigPo.Score.Valid {
			baseScore = float32(customerConfigPo.Score.Float64)
		}
		if customerConfigPo.TopK != 0 {
			baseTopK = int(customerConfigPo.TopK)
		}
		clearContextTime = customerConfigPo.ClearContextTime
		postModel = customerConfigPo.PostModel
	}
	l.setPostModel(postModel)

	// 指令匹配， 根据响应值判定是否需要去调用 openai 接口了
	proceed, _ := l.FactoryCommend(req)
	if !proceed {
		return
	}
	if l.message != "" {
		req.Msg = l.message
	}

	// openai client
	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).
		WithModel(l.model).
		WithBaseHost(l.baseHost).
		WithOrigin(l.svcCtx.Config.OpenAi.Origin).
		WithEngine(l.svcCtx.Config.OpenAi.Engine).WithPostModel(postModel)
	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	// context
	collection := openai.NewUserContext(
		openai.GetUserUniqueID(req.CustomerID, req.OpenKfID),
	).WithPrompt(l.basePrompt).WithModel(l.model).WithClient(c)

	//判断是否需要清除聊天记录
	if clearContextTime > 0 {
		duration := time.Duration(clearContextTime) * time.Minute
		formattedTime := time.Now().Add(-duration).Format("2006-01-02 15:04:05")
		clearStatus, err1 := l.CheckClearContext(context.Background(), req.OpenKfID, req.CustomerID, formattedTime)
		if err1 != nil {
			util.Error("CustomerChatLogic:CustomerChat:读取 stream 失败:" + err1.Error())
			return nil, err1
		}
		if clearStatus {
			collection.Clear()
			collection.Messages = []openai.ChatModelMessage{}
			collection.Summary = []openai.ChatModelMessage{}
		}
	}

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
				result := milvusService.SearchFromQA(embedding, baseTopK)
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
					go sendToUser(req.OpenKfID, "", req.CustomerID, "正在思考中，也许您还想知道"+"\n\n"+tempMessage, l.svcCtx.Config)
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
				result := milvusService.SearchFromArticle(embedding, baseTopK)

				for _, item := range result {
					fmt.Println("text:", item.CnText)
					fmt.Println("score:", item.Score)
					if item.Score < baseScore {
						continue
					}
					embeddingData = append(embeddingData, EmbeddingData{
						text: item.CnText,
					})
				}

				if len(embeddingData) > 0 {
					messageText += "When given CONTEXT you answer questions using only that information,and you always format your output in markdown.Answer with chinese.\n\n"
					messageText += "CONTEXT:"
					for _, chat := range embeddingData {
						messageText += chat.text + "\n\n"
						if c.WithModel(l.model).WithBaseHost(l.baseHost).GetNumTokens(messageText) > MaxToken {
							break
						}
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

				messageText, err := c.ChatStream(prompts, channel)

				if err != nil {
					util.Error("CustomerChatLogic:CustomerChat:读取 stream 失败:" + err.Error())
					logx.Error("读取 stream 失败：", err.Error())
					sendToUser(req.OpenKfID, "", req.CustomerID, "系统拥挤，稍后再试~"+err.Error(), l.svcCtx.Config)
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
				go l.InsertWechatUser(req.CustomerID)
			}()

			var rs []rune
			// 加快初次响应的时间 后续可改为阶梯式（用户体验好）
			first := true
			for {
				s, ok := <-channel
				if !ok {
					// 数据接受完成
					if len(rs) > 0 {
						go sendToUser(req.OpenKfID, "", req.CustomerID,
							string(rs)+"\n--------------------------------\n"+req.Msg,
							l.svcCtx.Config,
						)
					}
					return
				}
				rs = append(rs, []rune(s)...)

				if first && len(rs) > 50 && strings.Contains(s, "\n\n") {
					go sendToUser(req.OpenKfID, "", req.CustomerID, strings.TrimRight(string(rs), "\n\n"), l.svcCtx.Config)
					rs = []rune{}
					first = false
				} else if len(rs) > 200 && strings.Contains(s, "\n\n") {
					go sendToUser(req.OpenKfID, "", req.CustomerID, strings.TrimRight(string(rs), "\n\n"), l.svcCtx.Config)
					rs = []rune{}
				}
			}
		}

		messageText, err := c.Chat(prompts)

		if err != nil {
			util.Error("CustomerChatLogic:CustomerChat:error:" + err.Error())
			sendToUser(req.OpenKfID, "", req.CustomerID, "系统错误:"+err.Error(), l.svcCtx.Config)
			return
		}

		// 然后把数据 发给对应的客户
		go sendToUser(req.OpenKfID, "", req.CustomerID, messageText, l.svcCtx.Config)
		collection.Set("", messageText, true)
		_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
			User:       req.CustomerID,
			OpenKfId:   req.OpenKfID,
			MessageId:  req.MsgID,
			ReqContent: req.Msg,
			ResContent: messageText,
		})
		go l.InsertWechatUser(req.CustomerID)
	}()

	return &types.CustomerChatReply{
		Message: "ok",
	}, nil
}

func (l *CustomerChatLogic) setModelNameV2() (ls *CustomerChatLogic) {
	m := l.svcCtx.Config.WeCom.Model
	if m == "" || (m != openai.TextModel && m != openai.ChatModel && m != openai.ChatModelNew && m != openai.ChatModel4) {
		m = openai.TextModel
	}
	l.svcCtx.Config.WeCom.Model = m
	l.model = m
	return l
}

func (l *CustomerChatLogic) InsertWechatUser(user string) {
	ctx := context.Background()
	if user == "" {
		return
	}
	//先查询下
	wechatUserPo, err := l.svcCtx.WechatUserModel.FindOneByQuery(ctx,
		l.svcCtx.WechatUserModel.RowBuilder().Where(squirrel.Eq{"user": user}),
	)
	if err != nil {
		fmt.Printf("InsertWechatUser FindOneByQuery error: %v", err)
		return
	}
	//已存在就不需要重复去查询并且插入了
	if wechatUserPo != nil && wechatUserPo.Nickname != "" {
		return
	}

	res, err := wecom.GetCustomer([]string{user})
	if err != nil {
		return
	}
	if len(res) > 0 {
		wechatInfo := res[0]
		wechatModel := &model.WechatUser{
			User:     user,
			Nickname: wechatInfo.Nickname,
			Avatar:   wechatInfo.Avatar,
			Gender:   wechatInfo.Gender,
			Unionid:  wechatInfo.Unionid,
		}
		if wechatUserPo != nil && wechatUserPo.Id > 0 {
			wechatUserPo.User = user
			wechatUserPo.Nickname = wechatInfo.Nickname
			wechatUserPo.Avatar = wechatInfo.Avatar
			wechatUserPo.Gender = wechatInfo.Gender
			wechatUserPo.Unionid = wechatInfo.Unionid
			wechatUserPo.UpdatedAt = time.Now()
			err = l.svcCtx.WechatUserModel.Update(ctx, wechatUserPo)
			return
		}
		_, err = l.svcCtx.WechatUserModel.Insert(ctx, wechatModel)
	}
	return
}

func (l *CustomerChatLogic) CheckClearContext(ctx context.Context, openKfId, userId, formattedTime string) (bool, error) {

	chatPo, err := l.svcCtx.ChatModel.FindAll(ctx,
		l.svcCtx.ChatModel.RowBuilder().Where(squirrel.Eq{"open_kf_id": openKfId}).Where(squirrel.Eq{"user": userId}).Where(" created_at > ?", formattedTime),
	)
	if err != nil {
		return false, err
	}
	if len(chatPo) > 0 {
		return false, nil
	}
	return true, nil
}

// CustomerCommendSystem 查看系统信息
type CustomerCommendSystem struct{}

func (p CustomerCommendSystem) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	tips := fmt.Sprintf(
		"系统信息\n系统版本为：%s \nmodel 版本为：%s \n系统基础设定：%s \n",
		l.svcCtx.Config.SystemVersion,
		l.model,
		l.basePrompt,
	)
	sendToUser(req.OpenKfID, "", req.CustomerID, tips, l.svcCtx.Config)
	return false
}

// CustomerCommendHelp 查看所有指令
type CustomerCommendHelp struct{}

func (p CustomerCommendHelp) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	tips := fmt.Sprintf(
		"支持指令：\n\n%s\n%s\n%s\n",
		"基础模块🕹️\n\n#help       查看所有指令",
		"#system 查看会话系统信息",
		"#clear 清空当前会话的数据",
	)
	sendToUser(req.OpenKfID, "", req.CustomerID, tips, l.svcCtx.Config)
	return false
}

type CustomerCommendAbout struct{}

func (p CustomerCommendAbout) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	sendToUser(req.OpenKfID, "", req.CustomerID, "https://github.com/whyiyhw/chatgpt-wechat", l.svcCtx.Config)
	return false
}

type CustomerCommendDirect struct{}

func (p CustomerCommendDirect) customerExec(l *CustomerChatLogic, req *types.CustomerChatReq) bool {
	msg := strings.Replace(req.Msg, "#direct:", "", -1)
	sendToUser(req.OpenKfID, "", req.CustomerID, msg, l.svcCtx.Config)
	return false
}

func (l *CustomerChatLogic) setPostModel(postModel string) (ls *CustomerChatLogic) {
	var m string
	if "" != postModel {
		m = postModel
	}
	if m == "" || (m != openai.TextModel && m != openai.ChatModel && m != openai.ChatModelNew && m != openai.ChatModel4) {
		m = openai.ChatModel
	}
	l.postModel = m
	return l
}
