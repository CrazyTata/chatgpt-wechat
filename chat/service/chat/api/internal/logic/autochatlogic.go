package logic

import (
	"chat/common/openai"
	"chat/common/util"
	"chat/common/wecom"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"chat/service/chat/model"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

type AutoChatLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	model             string
	baseHost          string
	basePrompt        string
	message           string
	agentSecret       string
	customerChatLogic *CustomerChatLogic
}

func NewAutoChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AutoChatLogic {
	return &AutoChatLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		customerChatLogic: NewCustomerChatLogic(ctx, svcCtx),
	}
}

func (l *AutoChatLogic) AutoChat(req *types.DoGenerateActiveChatContentRequest) (resp *types.DoGenerateActiveChatContentReply, err error) {
	userId, kfId, agentId, err := NewChatRecordLogic(l.ctx, l.svcCtx).FormatCondition(req.UserNickname, req.KfName, req.ChatType)
	if err != nil {
		return nil, err
	}
	if agentId != 0 {
		err = l.Chat(userId, req.ContextMessage, req.Prompt, agentId)

	} else if kfId != "" {
		// Create a new Node with a Node number of 1
		node, errNode := snowflake.NewNode(1)
		if errNode != nil {
			return nil, errNode
		}

		err = l.CustomerChat(userId, kfId, req.ContextMessage, req.Prompt, node.Generate().String())

	}
	return &types.DoGenerateActiveChatContentReply{Message: "ok"}, nil
}

func (l *AutoChatLogic) CustomerChat(customerID, openKfID, message, newPromote, messageID string) (err error) {

	l.setModelName("").setBasePrompt("").setBaseHost()

	// openai client
	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).
		WithModel(l.model).
		WithBaseHost(l.baseHost).
		WithOrigin(l.svcCtx.Config.OpenAi.Origin).
		WithEngine(l.svcCtx.Config.OpenAi.Engine)
	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	//get prompt
	customerConfigPo, err := l.svcCtx.CustomerConfigModel.FindOneByQuery(context.Background(),
		l.svcCtx.CustomerConfigModel.RowBuilder().Where(squirrel.Eq{"kf_id": openKfID}),
	)

	if err != nil {
		return err
	}
	if customerConfigPo != nil {
		if customerConfigPo.Prompt != "" {
			l.basePrompt = customerConfigPo.Prompt
		}
	}

	// context
	collection := openai.NewUserContext(
		openai.GetUserUniqueID(customerID, openKfID),
	).WithPrompt(l.basePrompt).WithModel(l.model).WithClient(c)
	//
	//collection.Clear()
	//collection.Messages = []openai.ChatModelMessage{}
	//collection.Summary = []openai.ChatModelMessage{}

	// 然后 把 消息 发给 openai
	go func() {
		// 基于 summary 进行补充
		messageText := ""

		collection.Set(message, "", false)
		prompts1 := collection.GetChatSummary()
		prompts := collection.GetOtherChatSummary(prompts1, newPromote)
		if l.svcCtx.Config.Response.Stream {
			channel := make(chan string, 100)
			go func() {

				messageText, err := c.ChatStream(prompts, channel)

				if err != nil {
					logx.Error("读取 stream 失败：", err.Error())
					wecom.SendCustomerChatMessage(openKfID, customerID, "系统拥挤，稍后再试~"+err.Error())
					return
				}
				collection.Set("", messageText, true)
				// 再去插入数据
				_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
					User:       customerID,
					OpenKfId:   openKfID,
					MessageId:  messageID,
					ReqContent: message,
					ResContent: messageText,
				})
				go l.customerChatLogic.InsertWechatUser(customerID)
			}()

			var rs []rune
			// 加快初次响应的时间 后续可改为阶梯式（用户体验好）
			first := true
			for {
				s, ok := <-channel
				if !ok {
					// 数据接受完成
					if len(rs) > 0 {
						go wecom.SendCustomerChatMessage(openKfID, customerID, string(rs))
					}
					return
				}
				rs = append(rs, []rune(s)...)

				if first && len(rs) > 50 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(openKfID, customerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
					first = false
				} else if len(rs) > 200 && strings.Contains(s, "\n\n") {
					go wecom.SendCustomerChatMessage(openKfID, customerID, strings.TrimRight(string(rs), "\n\n"))
					rs = []rune{}
				}
			}
		}

		messageText, err := c.Chat(prompts)

		if err != nil {
			util.Error("AutoChatLogic:CustomerChat:error:" + err.Error())
			wecom.SendCustomerChatMessage(openKfID, customerID, "系统错误:"+err.Error())
			return
		}

		// 然后把数据 发给对应的客户
		go wecom.SendCustomerChatMessage(openKfID, customerID, messageText)
		collection.Set("", messageText, true)
		_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
			User:       customerID,
			OpenKfId:   openKfID,
			MessageId:  messageID,
			ReqContent: message,
			ResContent: messageText,
		})
		go l.customerChatLogic.InsertWechatUser(customerID)
	}()

	return nil
}

func (l *AutoChatLogic) Chat(userId, message, newPromote string, agentId int64) (err error) {

	var prompt, baseModel, agentSecret string

	//get config
	applicationConfigPo, err := l.svcCtx.ApplicationConfigModel.FindOneByQuery(context.Background(),
		l.svcCtx.ApplicationConfigModel.RowBuilder().Where(squirrel.Eq{"agent_id": agentId}),
	)

	if err != nil {
		return err
	}
	if applicationConfigPo != nil {
		prompt = applicationConfigPo.BasePrompt
		baseModel = applicationConfigPo.Model
		agentSecret = applicationConfigPo.AgentSecret

	}

	l.setModelName(baseModel).setBasePrompt(prompt).setAgentSecret(agentSecret).setBaseHost()

	// openai client
	c := openai.NewChatClient(l.svcCtx.Config.OpenAi.Key).
		WithModel(l.model).
		WithBaseHost(l.baseHost).
		WithOrigin(l.svcCtx.Config.OpenAi.Origin).
		WithEngine(l.svcCtx.Config.OpenAi.Engine)
	if l.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(l.svcCtx.Config.Proxy.Http).WithSocks5Proxy(l.svcCtx.Config.Proxy.Socket5)
	}

	// context
	collection := openai.NewUserContext(
		openai.GetUserUniqueID(userId, strconv.FormatInt(agentId, 10)),
	).WithPrompt(l.basePrompt).WithModel(l.model).WithClient(c)

	//collection.Clear()
	//collection.Messages = []openai.ChatModelMessage{}
	//collection.Summary = []openai.ChatModelMessage{}

	go func() {
		// 基于 summary 进行补充
		messageText := ""
		collection.Set(message, "", false)
		if l.model == openai.TextModel {
			messageText, err = c.Completion(collection.GetCompletionSummary())
			collection.Set("", messageText, true)
		} else {

			prompts1 := collection.GetChatSummary()
			prompts := collection.GetOtherChatSummary(prompts1, newPromote)

			if l.svcCtx.Config.Response.Stream {
				channel := make(chan string, 100)
				go func() {
					messageText, err := c.ChatStream(prompts, channel)
					if err != nil {
						errInfo := err.Error()
						if strings.Contains(errInfo, "maximum context length") {
							errInfo += "\n 请使用 #clear 清理所有上下文"
						}
						util.Error("AutoChatLogic:chat:error:" + errInfo)
						sendToUser(agentId, agentSecret, userId, "系统错误:"+errInfo, l.svcCtx.Config)
						return
					}
					collection.Set("", messageText, true)
					// 再去插入数据
					_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
						AgentId:    agentId,
						User:       userId,
						ReqContent: message,
						ResContent: messageText,
					})
				}()

				var rs []rune
				first := true
				for {
					s, ok := <-channel
					if !ok {
						// 数据接受完成
						if len(rs) > 0 {
							go sendToUser(agentId, agentSecret, userId, string(rs), l.svcCtx.Config)
						}
						return
					}
					rs = append(rs, []rune(s)...)

					if first && len(rs) > 50 && strings.Contains(s, "\n\n") {
						go sendToUser(agentId, agentSecret, userId, strings.TrimRight(string(rs), "\n\n"), l.svcCtx.Config)
						rs = []rune{}
						first = false
					} else if len(rs) > 100 && strings.Contains(s, "\n\n") {
						go sendToUser(agentId, agentSecret, userId, strings.TrimRight(string(rs), "\n\n"), l.svcCtx.Config)
						rs = []rune{}
					}
				}
			}

			messageText, err = c.Chat(prompts)
		}

		if err != nil {
			errInfo := err.Error()
			if strings.Contains(errInfo, "maximum context length") {
				errInfo += "\n 请使用 #clear 清理所有上下文"
			}
			util.Error("AutoChatLogic:chat:error:" + errInfo)
			sendToUser(agentId, agentSecret, userId, "系统错误:"+errInfo, l.svcCtx.Config)
			return
		}

		// 把数据 发给微信用户
		go sendToUser(agentId, agentSecret, userId, messageText, l.svcCtx.Config)

		collection.Set("", messageText, true)
		// 再去插入数据
		_, _ = l.svcCtx.ChatModel.Insert(context.Background(), &model.Chat{
			AgentId:    agentId,
			User:       userId,
			ReqContent: message,
			ResContent: messageText,
		})
	}()

	return
}

func (l *AutoChatLogic) setBaseHost() (ls *AutoChatLogic) {
	if l.svcCtx.Config.OpenAi.Host == "" {
		l.svcCtx.Config.OpenAi.Host = "https://api.openai.com"
	}
	l.baseHost = l.svcCtx.Config.OpenAi.Host
	return l
}

func (l *AutoChatLogic) setModelName(baseModel string) (ls *AutoChatLogic) {
	m := l.svcCtx.Config.WeCom.Model
	if "" != baseModel {
		m = baseModel
	}
	if m == "" || (m != openai.TextModel && m != openai.ChatModel && m != openai.ChatModelNew && m != openai.ChatModel4) {
		m = openai.TextModel
	}
	l.model = m
	return l
}

func (l *AutoChatLogic) setBasePrompt(prompt string) (ls *AutoChatLogic) {
	p := l.svcCtx.Config.WeCom.BasePrompt
	if prompt != "" {
		p = prompt
	}
	if p == "" {
		p = "你是 ChatGPT, 一个由 OpenAI 训练的大型语言模型, 你旨在回答并解决人们的任何问题，并且可以使用多种语言与人交流。\n"
	}
	l.basePrompt = p
	return l
}

func (l *AutoChatLogic) setAgentSecret(agentSecret string) (ls *AutoChatLogic) {
	l.agentSecret = agentSecret
	return l
}
