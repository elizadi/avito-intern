package main

import (
	"avito/internal/app"
	"avito/internal/domain"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddUserToSegmentRequest struct {
	SlugsAdd    []string
	SlugsDelete []string
	ID          uint64
}

func main() {
	router := gin.Default()
	uc, err := app.Execute()
	if err != nil {
		// добавить логгер
		fmt.Println("cannot create usecase")
		return
	}
	s := server{uc: uc}
	router.POST("/segment/:slug", s.CreateSegment)
	router.DELETE("/segment/:slug", s.DeleteSegment)
	router.GET("/segment", s.GetSegments)
	router.POST("/addUserToSegment", s.AddUserToSegment)
	router.GET("/activeUserSegments/:id", s.GetActiveUserSegments)
	router.GET("/downloadOperations/:year/:month", s.GetOperations)

	router.Run("0.0.0.0:8080")
}

type server struct {
	uc domain.UseCase
}

// нужно обрабатывать ошибки и возвращать свои
func (s *server) CreateSegment(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrEmptyParameter.Error())
		return
	}
	err := s.uc.CreateSegment(slug)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.String(http.StatusOK, "Success")
}

func (s *server) DeleteSegment(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrEmptyParameter.Error())
		return
	}
	err := s.uc.DeleteSegment(slug)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.String(http.StatusOK, "Success")
}

func (s *server) GetSegments(ctx *gin.Context) {
	segments, err := s.uc.GetSegments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, segments)
}

func (s *server) AddUserToSegment(ctx *gin.Context) {
	var req AddUserToSegmentRequest
	ctx.BindJSON(&req)
	err := s.uc.AddUserToSegment(req.SlugsAdd, req.SlugsDelete, req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.String(http.StatusOK, "Success")
}

func (s *server) GetActiveUserSegments(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrEmptyParameter.Error())
		return
	}
	new_id, _ := strconv.ParseUint(string(id), 10, 64)
	segment, err := s.uc.GetActiveUserSegments(new_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, segment)
}

func (s *server) GetOperations(ctx *gin.Context) {
	year := ctx.Param("year")
	month := ctx.Param("month")
	if year == "" || month == "" {
		ctx.JSON(http.StatusBadRequest, domain.ErrEmptyParameter.Error())
		return
	}
	new_year, err := strconv.ParseUint(string(year), 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "year must be an integer")
		return
	}
	new_month, err := strconv.ParseUint(string(month), 10, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "month must be an integer")
		return
	}
	operations, err := s.uc.GetOperations(int(new_year), int(new_month))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	w := ctx.Writer
	filename := fmt.Sprintf("operations%s%s.csv", year, month)
	ctx.Writer.Header().Set("attachment", fmt.Sprintf("filename=%s", filename))
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"UserID",
		"SegmentSlug",
		"Operation",
		"Data",
	}
	if err := writer.Write(header); err != nil {
		return
	}

	for _, operation := range operations {
		var csvRow []string
		csvRow = append(csvRow, fmt.Sprint(operation.UserID), operation.SegmentSlug, operation.Operation, operation.Data.String())
		if err := writer.Write(csvRow); err != nil {
			return
		}
	}
}

// id uint64
