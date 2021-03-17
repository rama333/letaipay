package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		context.JSON(http.StatusBadRequest, "state: false")
		return
	}

	if  len(context.Param("imsi")) != 19 || context.Param("imsi")[0:7] != "8970127" {
		context.JSON(http.StatusBadRequest, "state: false")
		return
	}

	switch context.Param("state") {
	case "disapprove":
		s.dbStorage.UpdateStateIMSI(context.Param("imsi"), 2)
	case "aproved":
		s.dbStorage.UpdateStateIMSI(context.Param("imsi"), 1)

	default:
		context.JSON(http.StatusBadRequest, "state: false")

	}


	context.JSON(200,  "state:true")
}
