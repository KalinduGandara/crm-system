package route

import (
	"time"

	"github.com/KalinduGandara/crm-system/api/controller"
	"github.com/KalinduGandara/crm-system/bootstrap"
	"github.com/KalinduGandara/crm-system/db"
	"github.com/KalinduGandara/crm-system/domain"
	"github.com/KalinduGandara/crm-system/repository"
	"github.com/KalinduGandara/crm-system/usecase"
	"github.com/gin-gonic/gin"
)

func NewLoginRouter(env *bootstrap.Env, timeout time.Duration, db db.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	lc := &controller.LoginController{
		LoginUsecase: usecase.NewLoginUsecase(ur, timeout),
		Env:          env,
	}
	group.POST("/login", lc.Login)
}
