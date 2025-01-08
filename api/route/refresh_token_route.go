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

func NewRefreshTokenRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	rtc := &controller.RefreshTokenController{
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(ur, timeout),
		Env:                 env,
	}
	group.POST("/refresh", rtc.RefreshToken)
}
