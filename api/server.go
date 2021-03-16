package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"letaipays/entity"
	"sync"
)

type DBStorage interface {
	AddUserFullName(user_id int, fullname string) error
	AddUserNumberGph(user_id int, numberGph string) error
	AddUserNameDealer(user_id int, nameDealer string) error
	GetUser(user_id int) (u entity.User, err error)
	AddImsi(userId int, imsi string) error
	GetImsi(imsi string) (u entity.Imsi, err error)
	AddUserCity(userId int,city string) (error)
	GetAllData() (u []entity.DataAll, err error)
	GetAllDataWithUser(userId int) (u []entity.DataAll, err error)
}

type Server struct {
	dbStorage DBStorage
	log *logrus.Entry
	stop chan struct{}
	wg sync.WaitGroup
	gin *gin.Engine
	port string
}

func NewServer(storage DBStorage, port string) (*Server, error)  {
	s:= &Server{
		dbStorage: storage,
		log: logrus.WithField("subsystem", "server"),
		port: port,
	}

	s.stop = make(chan struct{})

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("api/v1")
	{
		v1.GET("reports", s.GetReports)
	}

	s.gin = r

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		for  {

			select {
			case <-s.stop:
				return
			default:
			}

			s.log.Info("start go")

			err := r.Run(fmt.Sprintf(":%v", s.port))
			if err != nil{
				s.log.WithError(err).Error("failed to start")
			}
		}

	}()


	return s, nil
}

//func (s *Server) Stop()  {
//	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
//	defer cancel()
//
//	close(s.stop)
//
//	err := s.gin.Run()
//}