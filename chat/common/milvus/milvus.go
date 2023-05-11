package milvus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"log"
	"strconv"
	"time"
)

const (
	QA_COLLECTION            = "q_a_demo"
	QA_VECTOR_DIMENSION      = 1024
	ARTICLE_COLLECTION       = "articles"
	ARTICLE_VECTOR_DIMENSION = 1536
)

type Milvus struct {
	client client.Client
	ctx    context.Context
}

func InitMilvus(addr, username, password string) (milvus *Milvus, err error) {
	milvus = new(Milvus)
	ctx, _ := context.WithTimeout(context.Background(), 360*time.Second)
	milvus.ctx = ctx
	if username != "" {
		milvus.client, err = client.NewDefaultGrpcClientWithAuth(ctx, addr, username, password)
	} else {
		milvus.client, err = client.NewGrpcClient(ctx, addr)
	}
	return
}

func (m Milvus) CloseClient() {
	m.client.Close()
}

func (m Milvus) search(collectionName string, embeddings []float64, dimension int, fields []string, vectorField string, topK int) (sr []client.SearchResult, err error) {
	// load collection with async=false
	err = m.client.LoadCollection(m.ctx, collectionName, false)
	if err != nil {
		return
	}
	var searchEmbedding []float32

	for i, embedding := range embeddings {
		if i >= dimension {
			break
		}
		searchEmbedding = append(searchEmbedding, float32(embedding))
	}
	vector := entity.FloatVector(searchEmbedding[:])
	// Use flat search param
	var expr string
	sp, err := entity.NewIndexIvfFlatSearchParam(10)
	if err != nil {
		return
	}
	sr, err = m.client.Search(
		m.ctx, collectionName,
		[]string{},
		expr,
		fields,
		[]entity.Vector{vector},
		vectorField,
		entity.L2,
		topK,
		sp,
	)
	return
}

func (m Milvus) clearUp(collectionName string) {
	_ = m.client.ReleaseCollection(m.ctx, collectionName)
}

func (m Milvus) SearchFromQA(films []float64, topK int) (qas []QA) {
	sr, err := m.search(QA_COLLECTION, films, QA_VECTOR_DIMENSION, []string{"ID", "Q", "A"}, "Vector", topK)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, result := range sr {
		var idColumn *entity.ColumnInt64
		var qColumn *entity.ColumnVarChar
		var aColumn *entity.ColumnVarChar
		for _, field := range result.Fields {
			if field.Name() == "ID" {
				c, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
			if field.Name() == "Q" {
				q, ok := field.(*entity.ColumnVarChar)
				if ok {
					qColumn = q
				}
			}
			if field.Name() == "A" {
				a, ok := field.(*entity.ColumnVarChar)
				if ok {
					aColumn = a
				}
			}
		}
		if idColumn == nil {
			log.Fatal("result field not math")
		}
		for i := 0; i < result.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			q, err := qColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			a, err := aColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			qa := new(QA)
			qa.ID = id
			qa.Q = q
			qa.A = a
			qa.Score = result.Scores[i]
			qas = append(qas, *qa)
		}
	}
	// clean up
	defer m.clearUp(QA_COLLECTION)
	return
}

func (m Milvus) Save(films []Articles, collectionName string) (err error) {
	has, err := m.client.HasCollection(m.ctx, collectionName)
	if err != nil {
		fmt.Printf("failed to check whether collection exists: %v+\n", err)
		return
	}
	if !has {
		schema := &entity.Schema{
			CollectionName: collectionName,
			Description:    "this is the ashley collection for insert and search",
			AutoID:         false,
			Fields: []*entity.Field{
				{
					Name:       "id",
					DataType:   entity.FieldTypeInt64, // int64 only for now
					PrimaryKey: true,
					AutoID:     false,
				},
				{
					Name:     "name",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						entity.TypeParamMaxLength: "100",
					},
				},
				{
					Name:     "en_text",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						entity.TypeParamMaxLength: "10000",
					},
				},
				{
					Name:     "cn_text",
					DataType: entity.FieldTypeVarChar,
					TypeParams: map[string]string{
						entity.TypeParamMaxLength: "10000",
					},
				},
				{
					Name:     "vector",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						entity.TypeParamDim: strconv.Itoa(ARTICLE_VECTOR_DIMENSION),
					},
				},
			},
		}
		err = m.client.CreateCollection(m.ctx, schema, 1) // only 1 shard
		if err != nil {
			log.Fatal("failed to create collection:", err.Error())
		}
	}
	id := make([]int64, 0, len(films))
	name := make([]string, 0, len(films))
	enText := make([]string, 0, len(films))
	cnText := make([]string, 0, len(films))
	vector := make([][]float32, 0, len(films))
	for idx, film := range films {
		id = append(id, film.ID)
		name = append(name, film.Name)
		enText = append(enText, film.EnText)
		cnText = append(cnText, film.CnText)
		vector = append(vector, films[idx].Vector[:]) // prevent same vector
	}

	idColumn := entity.NewColumnInt64("id", id)
	nameColumn := entity.NewColumnVarChar("name", name)
	enTextColumn := entity.NewColumnVarChar("en_text", enText)
	cnTextColumn := entity.NewColumnVarChar("cn_text", cnText)
	vectorColumn := entity.NewColumnFloatVector("vector", ARTICLE_VECTOR_DIMENSION, vector)

	// insert into default partition
	ret, err := m.client.Insert(m.ctx, collectionName, "", idColumn, nameColumn, enTextColumn, cnTextColumn, vectorColumn)

	if err != nil {
		fmt.Printf("failed to insert film data: %v", err)
		return
	}
	fmt.Printf("insert completed result: %v", ret)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err = m.client.Flush(ctx, collectionName, false)
	if err != nil {
		fmt.Printf("failed to flush collection: %v", err)
	}
	return
}

func (m Milvus) SearchFromArticle(embeddings []float64, topK int) (articles []Articles) {
	sr, err := m.search(ARTICLE_COLLECTION, embeddings, ARTICLE_VECTOR_DIMENSION, []string{"id", "cn_text"}, "vector", topK)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(sr)
	for _, result := range sr {

		var idColumn *entity.ColumnInt64
		var textColumn *entity.ColumnVarChar
		for _, field := range result.Fields {
			if field.Name() == "id" {
				c, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
			if field.Name() == "cn_text" {
				q, ok := field.(*entity.ColumnVarChar)
				if ok {
					textColumn = q
				}
			}
		}
		if idColumn == nil {
			fmt.Println("result field not math")
			return
		}
		for i := 0; i < result.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			text, err := textColumn.ValueByIdx(i)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			article := new(Articles)
			article.ID = id
			article.CnText = text
			article.Score = result.Scores[i]
			articles = append(articles, *article)
		}
	}
	// clean up
	defer m.clearUp(ARTICLE_COLLECTION)
	return
}

func (m Milvus) DeleteCollection(collectionName string) (err error) {
	has, err := m.client.HasCollection(m.ctx, collectionName)
	if err != nil {
		fmt.Printf("failed to check whether collection exists: %v+\n", err)
		return
	}
	if has {
		// collection with same name exist, clean up mess
		_ = m.client.DropCollection(m.ctx, collectionName)
	}
	return
}

func (m Milvus) QueryArticleByName(ctx context.Context, names []string) (result []string, err error) {
	// load collection with async=false
	err = m.client.LoadCollection(m.ctx, ARTICLE_COLLECTION, false)
	if err != nil {
		return
	}
	nameStr, _ := json.Marshal(names)

	queryResult, err := m.client.Query(ctx, ARTICLE_COLLECTION, nil, "name in "+string(nameStr), []string{"name"})
	if err != nil {
		return
	}
	fmt.Printf("%#v\n", queryResult)
	for _, v := range queryResult {
		if v.Name() != "name" {
			continue
		}
		var nameColumn *entity.ColumnVarChar

		c, ok := v.(*entity.ColumnVarChar)
		if ok {
			nameColumn = c
		}

		result = nameColumn.Data()
		break
	}
	fmt.Println(result)
	return
}
