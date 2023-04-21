package milvus

type QA struct {
	ID    int64
	Q     string
	A     string
	Score float32
}

type Articles struct {
	ID     int64 //snow
	Name   string
	EnText string
	CnText string
	Vector [ARTICLE_VECTOR_DIMENSION]float32
	Score  float32
}
