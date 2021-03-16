package api

import "github.com/gin-gonic/gin"

func (s *Server) GetReports(ctx *gin.Context) {

	report, err := s.dbStorage.GetAllData()

	if err != nil{
		ctx.JSON(500, report)
		return
	}

	ctx.JSON(200, report)
}
