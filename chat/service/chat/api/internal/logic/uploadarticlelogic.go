package logic

import (
	"chat/common/milvus"
	"chat/common/openai"
	"chat/common/util"
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"net/http"
	"os"
	"strings"
	"time"

	"chat/service/chat/api/internal/svc"
	"chat/service/chat/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadArticleLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	model      string
	baseHost   string
	basePrompt string
}

const MAX_UPLOAD_SIZE = 2000

func NewUploadArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadArticleLogic {
	return &UploadArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (f *UploadArticleLogic) UploadArticle(req *types.UploadArticleHandlerReq, r *http.Request) (resp *types.UploadArticleHandlerReply, err error) {
	ctx := context.Background()
	// 1. parse input , type multipart/form-data
	r.ParseMultipartForm(1000)

	// 2. retrieve file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Error retrieving file from form-data %v", err)
		return
	}
	defer file.Close()

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

	//get data
	var rows [][]string
	if strings.Contains(fileName, "csv") {
		//CSV文件
		rows, err = util.GetCSVDataByPath(fileName)
	} else {
		rows, err = util.GetExcelDataByPath(fileName)
	}
	if err != nil {
		return nil, err
	}
	//check data
	baseData, err := f.checkPreview(ctx, rows)

	if err != nil {
		return nil, err
	}
	fmt.Println("check file success")
	//format data
	data, err := f.formatData(ctx, baseData)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败,请检查文件内容是否正确")
	}
	fmt.Println("format file success")

	//save data
	err = f.SaveData(data)

	fmt.Println("save file success")

	return
}

func (f *UploadArticleLogic) checkPreview(ctx context.Context, rows [][]string) ([]milvus.Articles, error) {

	if len(rows) <= 1 {
		return nil, fmt.Errorf("文件内容为空")
	}
	if len(rows) > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("超过最大上传数量:%d", MAX_UPLOAD_SIZE)
	}
	rows1 := rows[1:]
	var ret []milvus.Articles
	var names []string
	for _, v := range rows1 {
		names = append(names, v[0]+v[1])
	}

	fmt.Println(names)

	existInfo, err := f.QueryArticleByName(ctx, names)
	if err != nil {
		return nil, err
	}

outerLoop:
	for _, vv := range rows1 {
		if len(existInfo) > 0 {
			for _, vvv := range existInfo {
				if vvv == vv[0]+vv[1] {
					continue outerLoop
				}
			}
		}
		fi := milvus.Articles{}
		fi.Name = vv[0] + vv[1]
		fi.EnText = vv[2]
		fi.CnText = vv[3]
		ret = append(ret, fi)
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("文件内容为空")
	}

	return ret, nil
}

func (f *UploadArticleLogic) formatData(ctx context.Context, baseData []milvus.Articles) (ret []milvus.Articles, err error) {

	c := openai.NewChatClient(f.svcCtx.Config.OpenAi.Key).WithModel(f.model).WithBaseHost(f.baseHost)
	if f.svcCtx.Config.Proxy.Enable {
		c = c.WithHttpProxy(f.svcCtx.Config.Proxy.Http).WithSocks5Proxy(f.svcCtx.Config.Proxy.Socket5)
	}
	for _, v := range baseData {
		// Create a new Node with a Node number of 1
		node, errNode := snowflake.NewNode(1)
		if errNode != nil {
			return nil, errNode
		}
		parts, vectorErr := f.DealDataToVector(ctx, c, v.Name)
		if vectorErr != nil {
			fmt.Printf("vector error : %v", vectorErr)
		}
		fi := milvus.Articles{}
		// Generate a snowflake ID.
		fi.ID = node.Generate().Int64()
		fi.Name = v.Name
		fi.EnText = v.EnText
		fi.CnText = v.CnText
		for idx, vv := range parts {
			fi.Vector[idx] = float32(vv)
		}
		fmt.Println(fi.Name)
		ret = append(ret, fi)
	}

	return
}

func (f *UploadArticleLogic) SaveData(message []milvus.Articles) (err error) {

	//数据库没有
	milvusService, err := milvus.InitMilvus(f.svcCtx.Config.Embeddings.Milvus.Host, f.svcCtx.Config.Embeddings.Milvus.Username, f.svcCtx.Config.Embeddings.Milvus.Password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer milvusService.CloseClient()
	err = milvusService.Save(message, milvus.ARTICLE_COLLECTION)
	return
}

func (f *UploadArticleLogic) DealDataToVector(ctx context.Context, c *openai.ChatClient, message string) ([]float64, error) {

	// 把中文转换成向量
	res, err := c.CreateOpenAIEmbeddings(message)
	if err != nil {
		return nil, err
	}
	embedding := res.Data[0].Embedding
	return embedding, err
}

func (f *UploadArticleLogic) QueryArticleByName(ctx context.Context, names []string) (result []string, err error) {
	//数据库没有
	milvusService, err := milvus.InitMilvus(f.svcCtx.Config.Embeddings.Milvus.Host, f.svcCtx.Config.Embeddings.Milvus.Username, f.svcCtx.Config.Embeddings.Milvus.Password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer milvusService.CloseClient()

	result, err = milvusService.QueryArticleByName(ctx, names)
	return
}
