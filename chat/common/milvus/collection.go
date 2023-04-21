package milvus

type Articles struct {
	ID     int64 //snow
	Name   string
	EnText string
	CnText string
	Vector []float32
	Text   string
	Score  float32
}
