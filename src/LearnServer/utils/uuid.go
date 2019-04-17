package utils

import (
	"crypto/rand"
	"fmt"
)

type StrList []string

type Column struct {
	Title     string `json:"title"`
	DataIndex string `json:"dataIndex"`
}

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func MakeColumns(cs ...StrList) []Column {
	result := make([]Column, len(cs))
	for i, c := range cs {
		result[i] = Column{
			Title:     c[0],
			DataIndex: c[1],
		}
	}
	return result
}
