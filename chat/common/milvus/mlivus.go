package milvus

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"log"
)

const (
	QA_COLLECTION            = "q_a_demo"
	QA_VECTOR_DIMENSION      = 1024
	ARTICLE_COLLECTION       = "articles"
	ARTICLE_VECTOR_DIMENSION = 1024
)

type Milvus struct {
	client client.Client
	ctx    context.Context
}

func InitMilvus(addr, username, password string) (milvus *Milvus, err error) {
	milvus = new(Milvus)
	ctx := context.Background()
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
	sp, err := entity.NewIndexIvfFlatSearchParam(10)
	if err != nil {
		return
	}
	sr, err = m.client.Search(
		m.ctx, collectionName,
		[]string{},
		"",
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

func (m Milvus) SearchFromQA(films []float64) (qas []QA) {
	sr, err := m.search(QA_COLLECTION, films, QA_VECTOR_DIMENSION, []string{"ID", "Q", "A"}, "Vector", 4)
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

func (m Milvus) SearchFromArticle(embeddings []float64) (articles []Articles) {
	sr, err := m.search(ARTICLE_COLLECTION, embeddings, ARTICLE_VECTOR_DIMENSION, []string{"id", "cn_text"}, "vector", 1)
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
