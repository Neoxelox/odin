package view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/pkg/model"
	"github.com/neoxelox/odin/pkg/payload"
	"github.com/neoxelox/odin/pkg/usecase/auth"
	"github.com/neoxelox/odin/pkg/usecase/otp"
)

type AuthView struct {
	class.View
	otpCreator otp.CreatorUsecase
	authLogger auth.LoggerUsecase
}

func NewAuthView(configuration internal.Configuration, logger core.Logger, otpCreator otp.CreatorUsecase,
	authLogger auth.LoggerUsecase) *AuthView {
	return &AuthView{
		View:       *class.NewView(configuration, logger),
		otpCreator: otpCreator,
		authLogger: authLogger,
	}
}

func (self *AuthView) PostLoginStart(ctx echo.Context) error {
	request := &payload.PostLoginStartRequest{}
	response := &payload.PostLoginStartResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		newOTP, err := self.otpCreator.Create(ctx.Request().Context(), request.Phone, model.OTPType.SMS)
		switch {
		case err == nil:
			response.ID = newOTP.ID
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrAlreadySent().Is(err):
			return ExcOTPAlreadySent.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *AuthView) PostLoginEnd(ctx echo.Context) error {
	request := &payload.PostLoginEndRequest{}
	response := &payload.PostLoginEndResponse{}
	return self.Handle(ctx, class.Endpoint{
		Request: request,
	}, func() error {
		accessToken, user, err := self.authLogger.Login(ctx.Request().Context(), request.ID, request.Code, model.SessionMetadata{
			IP:         ctx.RealIP(),
			Device:     ctx.Request().Host,
			ApiVersion: "v1",
		})
		switch {
		case err == nil:
			response.AccessToken = accessToken
			if user != nil {
				response.User = &payload.User{
					ID:       user.ID,
					Phone:    user.Phone,
					Name:     user.Name,
					Email:    user.Email,
					Picture:  user.Picture,
					Birthday: user.Birthday,
				}
			}
			return ctx.JSON(http.StatusOK, response)
		case otp.ErrInvalid().Is(err):
			return internal.ExcInvalidRequest.Cause(err)
		case otp.ErrMaxAttempts().Is(err):
			return ExcOTPMaxAttempts.Cause(err)
		case otp.ErrWrongCode().Is(err):
			return ExcOTPWrongCode.Cause(err)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}

func (self *AuthView) PostLogout(ctx echo.Context) error {
	requestSession := RequestSession(ctx)
	response := &payload.PostLogoutResponse{}
	return self.Handle(ctx, class.Endpoint{}, func() error {
		err := self.authLogger.Logout(ctx.Request().Context(), *requestSession)
		switch {
		case err == nil:
			return ctx.JSON(http.StatusOK, response)
		default:
			return internal.ExcServerGeneric.Cause(err)
		}
	})
}
