package route

import (
	"time"

	"github.com/KalinduGandara/crm-system/api/controller"
	"github.com/KalinduGandara/crm-system/bootstrap"
	"github.com/KalinduGandara/crm-system/db/mongo"
	"github.com/KalinduGandara/crm-system/domain"
	"github.com/KalinduGandara/crm-system/repository"
	"github.com/KalinduGandara/crm-system/usecase"
	"github.com/gin-gonic/gin"
)

func NewTaskRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTaskRepository(db, domain.CollectionTask)
	tc := &controller.TaskController{
		TaskUsecase: usecase.NewTaskUsecase(tr, timeout),
	}
	group.GET("/task", tc.Fetch)
	group.POST("/task", tc.Create)
}
