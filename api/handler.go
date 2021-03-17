package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *Server) GetReports(ctx *gin.Context) {

	report, err := s.dbStorage.GetAllData()

	if err != nil{
		ctx.JSON(500, report)
		return
	}

	ctx.JSON(200, report)
}


func (s *Server) UpdateStateIMSI(context *gin.Context) {
	_, err := strconv.ParseInt(context.Param("imsi"), 10, 64)
	if err != nil {
		context.JSON(200,  gin.H{
			"state": false,
			"messages": "",
		})
		s.log.WithError(err).Error("Bad Request")
		return
	}

	if  len(context.Param("imsi")) != 19 || context.Param("imsi")[0:7] != "8970127" {
		context.JSON(200,  gin.H{
			"state": false,
			"messages": "",
		})
		s.log.WithError(errors.New("failed len imsi")).Error("Bad Request")
		return
	}


	if _, err := s.dbStorage.GetImsi(context.Param("imsi")); err != nil{
		context.JSON(200,  gin.H{
			"state": false,
			"messages": "not found imsi",
		})
		s.log.WithError(errors.New("not found imsi")).Error("")
		return
	}


	switch context.Param("state") {
	case "disapprove":
		err := s.dbStorage.UpdateStateIMSI(context.Param("imsi"), 2)
		if err != nil{
			context.JSON(200,  gin.H{
				"state": false,
				"messages": "",
			})
			s.log.WithError(err).Error("failed update db")
			return
		}
	case "aproved":
		err := s.dbStorage.UpdateStateIMSI(context.Param("imsi"), 1)
		if err != nil{
			context.JSON(200,  gin.H{
				"state": false,
				"messages": "",
			})
			s.log.WithError(err).Error("failed update db")
			return
		}

	default:
		context.JSON(200,  gin.H{
			"state": false,
			"messages": "",
		})
		return

	}

	context.JSON(200,  gin.H{
		"state": true,
		"messages": "success update state in " + context.Param("imsi"),
	})
}
