package milvus

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	ARTICLE_COLLECTION_NAME = "articles"
	ARTICLE_DIMENSION       = 128
)

func Search(films []float64, addr string) []QA {
	// setup context for client creation, use 8 seconds here
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	c, err := client.NewGrpcClient(ctx, addr)
	if err != nil {
		// handling error and exit, to make example simple here
		log.Fatal("failed to connect to milvus:", err.Error())
	}
	// in a main func, remember to close the client
	defer func(c client.Client) {
		_ = c.Close()
	}(c)

	// here is the collection name we use in this example
	collectionName := `q_a_demo`
	// load collection with async=false
	err = c.LoadCollection(ctx, collectionName, false)
	if err != nil {
		log.Fatal("failed to load collection:", err.Error())
	}
	log.Println("load collection completed")

	var searchFilm []float32
	for i, film := range films {
		if i > 1023 {
			break
		}
		searchFilm = append(searchFilm, float32(film))
	}
	vector := entity.FloatVector(searchFilm[:])
	// Use flat search param
	sp, err := entity.NewIndexIvfFlatSearchParam(5)
	if err != nil {
		log.Fatal("fail to create flat search param:", err.Error())
	}
	sr, err := c.Search(
		ctx, collectionName,
		[]string{},
		"",
		[]string{"ID", "Q", "A"},
		[]entity.Vector{vector},
		"Vector",
		entity.L2,
		4,
		sp,
	)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}

	fmt.Println(sr)

	var qas []QA
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
	defer func(c client.Client, ctx context.Context, collName string) {
		_ = c.ReleaseCollection(ctx, collName)
	}(c, ctx, collectionName)

	return qas
}

func Save(addr string, userName string, password string, films []Articles, collectionName string) (err error) {
	// setup context for client creation, use 8 seconds here
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	//c, err := client.NewGrpcClient(ctx, addr)
	c, err := client.NewDefaultGrpcClientWithAuth(ctx, addr, userName, password)
	if err != nil {
		// handling error and exit, to make example simple here
		fmt.Printf("failed to connect to milvus: %v ", err)
		return
	}
	// in a main func, remember to close the client
	defer func(c client.Client) {
		_ = c.Close()
	}(c)

	has, err := c.HasCollection(ctx, collectionName)
	if err != nil {
		fmt.Printf("failed to check whether collection exists: %v+\n", err)
		return
	}
	if has {
		// collection with same name exist, clean up mess
		_ = c.DropCollection(ctx, collectionName)
	}

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
					entity.TypeParamDim: strconv.Itoa(ARTICLE_DIMENSION),
				},
			},
		},
	}

	err = c.CreateCollection(ctx, schema, 1) // only 1 shard
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}

	id := make([]int64, 0, len(films))
	name := make([]string, 0, len(films))
	enText := make([]string, 0, len(films))
	cnText := make([]string, 0, len(films))
	vector := make([][]float32, 0, len(films))
	for _, film := range films {
		id = append(id, film.ID)
		name = append(name, film.Name)
		enText = append(enText, film.EnText)
		cnText = append(cnText, film.CnText)
		vector = append(vector, film.Vector)
	}

	idColumn := entity.NewColumnInt64("id", id)
	nameColumn := entity.NewColumnVarChar("name", name)
	enTextColumn := entity.NewColumnVarChar("en_text", enText)
	cnTextColumn := entity.NewColumnVarChar("cn_text", cnText)
	vectorColumn := entity.NewColumnFloatVector("vector", ARTICLE_DIMENSION, vector)

	// insert into default partition
	ret, err := c.Insert(ctx, collectionName, "", idColumn, nameColumn, enTextColumn, cnTextColumn, vectorColumn)

	if err != nil {
		fmt.Printf("failed to insert film data: %v", err)
		return
	}
	fmt.Printf("insert completed ret: %v", ret)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err = c.Flush(ctx, collectionName, false)
	if err != nil {
		fmt.Printf("failed to flush collection: %v", err)
	}
	return
}

func SearchFromArticle(embeddings []float64, addr, username, password string) (articles []Articles) {
	collectionName := ARTICLE_COLLECTION_NAME
	dimension := ARTICLE_DIMENSION
	// setup context for client creation, use 8 seconds here
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	c, err := client.NewDefaultGrpcClientWithAuth(ctx, addr, username, password)
	if err != nil {
		// handling error and exit, to make example simple here
		log.Fatal("failed to connect to milvus:", err.Error())
	}
	// in a main func, remember to close the client
	defer func(c client.Client) {
		_ = c.Close()
	}(c)

	// load collection with async=false
	err = c.LoadCollection(ctx, collectionName, false)
	if err != nil {
		log.Fatal("failed to load collection:", err.Error())
	}
	//log.Println("load collection completed")

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
		log.Fatal("fail to create flat search param:", err.Error())
	}
	sr, err := c.Search(
		ctx, collectionName,
		[]string{},
		"",
		[]string{"id", "cn_text"},
		[]entity.Vector{vector},
		"vector",
		entity.L2,
		3,
		sp,
	)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
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
			err = errors.New("result field not math")
			log.Fatal("result field not math")
		}
		for i := 0; i < result.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			text, err := textColumn.ValueByIdx(i)
			if err != nil {
				log.Fatal(err.Error())
			}
			article := new(Articles)
			article.ID = id
			article.CnText = text
			article.Score = result.Scores[i]
			articles = append(articles, *article)
		}
	}
	// clean up
	defer func(c client.Client, ctx context.Context, collName string) {
		_ = c.ReleaseCollection(ctx, collName)
	}(c, ctx, collectionName)

	return
}
