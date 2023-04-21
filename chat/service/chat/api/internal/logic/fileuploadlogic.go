package logic

import (
	"chat/common/milvus"
	"chat/common/openai"
	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"time"
)

const MAX_UPLOAD_SIZE = 2000

type FileUploadLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	model      string
	baseHost   string
	basePrompt string
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (f *FileUploadLogic) UploadArticle(ctx context.Context, req *types.FileUploadHandlerReq, r *http.Request) (resp *types.FileUploadHandlerReply, err error) {

	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(10)

	// 2. retrieve file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Error retrieving file from form-data %v", err)
		return
	}
	defer file.Close()
	//Todo 可以加入文件名校验逻辑
	dataBuf := make([]byte, handler.Size)
	_, err = file.Read(dataBuf)
	if err != nil {
		fmt.Printf("read file error %v", err)
		return
	}
	err = os.Mkdir("./temp-files", 0755)
	if err != nil {
		fmt.Printf("Cannot create a file when that file already exists %v \n ", err)
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	fileName := "./temp-files/" + "article_" + nowTime + handler.Filename
	fileHandle, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("create file error %v", err)
		return
	}
	_, err = fileHandle.Write(dataBuf)
	if err != nil {
		fmt.Printf("write file error %v", err)
		return
	}
	fmt.Println("upload file success")

	baseData, err := f.checkPreview(ctx, fileName)
	if err != nil {
		return nil, err
	}

	for _, r := range baseData {
		err = f.DealData(context.Background(), r)
		if err != nil {
			fmt.Printf("write file error %v", err)
		}
	}
	return
}

func (f *FileUploadLogic) checkPreview(ctx context.Context, path string) ([]milvus.Articles, error) {
	var err error

	fileHandler, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	rows := fileHandler.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	data, err := f.getData(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败,请检查文件内容是否正确")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("文件内容为空")
	}

	if len(data) > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("超过最大上传数量:%d", MAX_UPLOAD_SIZE)
	}
	return data, nil
}

func (f *FileUploadLogic) getData(ctx context.Context, bts [][]string) (ret []milvus.Articles, err error) {
	if len(bts) <= 1 {
		return nil, nil
	}
	data := bts[1:]
	for _, v := range data {
		if len(v) != 4 {
			continue
		}
		// Create a new Node with a Node number of 1
		node, errNode := snowflake.NewNode(1)
		if errNode != nil {
			return nil, errNode
		}

		// Generate a snowflake ID.
		ID := node.Generate().Int64()

		ret = append(ret, milvus.Articles{
			ID:     ID,
			Name:   v[0],
			EnText: v[2],
			CnText: v[3],
		})
	}

	return
}

func (f *FileUploadLogic) DealData(ctx context.Context, message milvus.Articles) (err error) {
	// openai client
	c := openai.NewChatClient(f.svcCtx.Config.OpenAi.Key).WithModel(f.model).WithBaseHost(f.baseHost)
	if f.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(f.svcCtx.Config.Proxy.Http).WithSocks5Proxy(f.svcCtx.Config.Proxy.Socket5)
	}

	// 把中文转换成向量
	res, err := c.CreateOpenAIEmbeddings(message.CnText)
	if err != nil {
		return
	}
	embedding := res.Data[0].Embedding
	film32 := make([]float32, milvus.ARTICLE_VECTOR_DIMENSION)
	for _, v := range embedding {
		film32 = append(film32, float32(v)) // 向量指支持float32，不支持float64
	}

	message.Vector = film32
	fmt.Println("create open ai embeddings success")

	//数据库没有
	milvusService, err := milvus.InitMilvus(f.svcCtx.Config.Embeddings.Milvus.Host, f.svcCtx.Config.Embeddings.Milvus.Username, f.svcCtx.Config.Embeddings.Milvus.Password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer milvusService.CloseClient()
	err = milvusService.Save([]milvus.Articles{message}, milvus.ARTICLE_COLLECTION)
	return
}
