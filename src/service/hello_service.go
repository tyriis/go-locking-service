package service

import (
	"fmt"

	"github.com/tyriis/rest-go/src/dao"
)

type HelloService struct {
	dao *dao.HelloDAO
}

func NewHelloService(dao *dao.HelloDAO) *HelloService {
	return &HelloService{dao: dao}
}

func (s *HelloService) GetHelloMessage() string {
	return s.dao.RetrieveMessage()
}

func (s *HelloService) GetPersonalizedHello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
