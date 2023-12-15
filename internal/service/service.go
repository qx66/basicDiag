package service

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qx66/basicDiag/internal/biz"
	"go.uber.org/zap"
)

type Service struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewService(logger *zap.Logger) (*Service, error) {
	db, err := sql.Open("sqlite3", "sqlite3.db")
	if err != nil {
		return nil, err
	}
	
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS diag (
	    id varchar(50) not null primary key ,
	    value text
	    );
`
	
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	
	return &Service{
		db:     db,
		logger: logger,
	}, nil
}

func (service *Service) RecordReport(c *gin.Context) {
	var req biz.BasicDiagResult
	err := c.ShouldBind(&req)
	if err != nil {
		service.logger.Error(
			"请求参数异常",
			zap.Error(err),
		)
		c.JSON(500, gin.H{"errCode": 500, "errMsg": "请求参数异常"})
		return
	}
	
	r, err := json.Marshal(&req)
	if err != nil {
		service.logger.Error(
			"序列化失败",
			zap.Error(err),
		)
		c.JSON(500, gin.H{"errCode": 500, "errMsg": "序列化失败"})
		return
	}
	
	id := uuid.NewString()
	_, err = service.db.Exec("insert into diag (`id`, `value`) VALUES (?, ?)", id, r)
	
	if err != nil {
		service.logger.Error(
			"记录数据失败",
			zap.Error(err),
		)
		c.JSON(500, gin.H{"errCode": 500, "errMsg": "记录数据失败"})
		return
	}
	
	service.logger.Info(
		"记录数据成功",
		zap.String("id", id),
	)
	c.JSON(200, gin.H{"errCode": 0, "errMsg": "ok", "id": id})
}

func (service *Service) GetReport(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"errCode": 400, "errMsg": "请输出id参数"})
		return
	}
	
	//rows, err := service.db.Query("SELECT * FROM diag")
	rows, err := service.db.Query("SELECT * FROM diag where `id` = ?", id)
	
	if err != nil {
		service.logger.Error(
			"查询失败",
			zap.Error(err),
		)
		c.JSON(500, gin.H{"errCode": 500, "errMsg": "查询失败"})
		return
	}
	
	r := make(map[string]biz.BasicDiagResult)
	
	for rows.Next() {
		var id string
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			service.logger.Error(
				"扫描数据失败",
				zap.Error(err),
			)
			c.JSON(500, gin.H{"errCode": 500, "errMsg": "查询失败"})
			return
		}
		var tr biz.BasicDiagResult
		err = json.Unmarshal([]byte(value), &tr)
		if err != nil {
			c.JSON(500, gin.H{"errCode": 500, "errMsg": "查询失败"})
			return
		}
		
		r[id] = tr
	}
	
	c.JSON(200, gin.H{"errCode": 0, "errMsg": "ok", "value": r})
	return
}
