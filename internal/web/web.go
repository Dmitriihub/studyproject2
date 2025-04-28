package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/oapi-codegen/runtime"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	openapi_types "github.com/oapi-codegen/runtime/types"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/krisch/crm-backend/domain"
	"github.com/krisch/crm-backend/dto"
	"github.com/krisch/crm-backend/internal/app"
	"github.com/krisch/crm-backend/internal/configs"
	"github.com/krisch/crm-backend/internal/helpers"
	"github.com/krisch/crm-backend/internal/legalentities"
	"github.com/krisch/crm-backend/internal/web/olegalentities"
	"github.com/krisch/crm-backend/pkg/redis"

	validator "github.com/go-playground/validator/v10"
)

type Web struct {
	app     *app.App
	Options configs.Configs
	Router  *echo.Echo
	Port    int

	UUID string

	Now       string
	Version   string
	Tag       string
	BuildTime string
}

// GetLegalEntitiesUuidBankAccounts implements olegalentities.StrictServerInterface.
func (w *Web) GetLegalEntitiesUuidBankAccounts(ctx context.Context, request olegalentities.GetLegalEntitiesUuidBankAccountsRequestObject) (olegalentities.GetLegalEntitiesUuidBankAccountsResponseObject, error) {
	accounts, err := w.app.LegalEntities.GetAllBankAccounts(request.Uuid)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get bank accounts").SetInternal(err)
	}

	dtos := make([]olegalentities.BankAccountDTO, len(accounts))
	for i, acc := range accounts {
		dtos[i] = w.toBankAccountDTO(acc)
	}

	return olegalentities.GetLegalEntitiesUuidBankAccounts200JSONResponse(dtos), nil
}

// PostLegalEntitiesUuidBankAccounts implements olegalentities.StrictServerInterface.
func (w *Web) PostLegalEntitiesUuidBankAccounts(ctx context.Context, request olegalentities.PostLegalEntitiesUuidBankAccountsRequestObject) (olegalentities.PostLegalEntitiesUuidBankAccountsResponseObject, error) {
	// Проверяем существование юридического лица
	if _, err := w.app.LegalEntities.GetLegalEntityByUUID(request.Uuid); err != nil {
		if errors.Is(err, legalentities.ErrLegalEntityNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "legal entity not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to check legal entity").SetInternal(err)
	}

	domainAcc := w.toBankAccountDomain(*request.Body, request.Uuid)

	// Валидация перед созданием
	if err := ValidateBankAccount(&domainAcc); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := w.app.LegalEntities.CreateBankAccount(&domainAcc); err != nil {
		switch {
		case errors.Is(err, legalentities.ErrInvalidBankAccountData):
			return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case errors.Is(err, legalentities.ErrPrimaryAccountExists):
			return nil, echo.NewHTTPError(http.StatusConflict, "primary account already exists for this legal entity")
		default:
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to create bank account").SetInternal(err)
		}
	}

	return olegalentities.PostLegalEntitiesUuidBankAccounts201Response{}, nil
}

func ValidateBankAccount(bankAccount *domain.BankAccount) error {
	if bankAccount.BIC == "" {
		return fmt.Errorf("%w: BIC is required", legalentities.ErrInvalidBankAccountData)
	}

	if len(bankAccount.BIC) != 9 {
		return fmt.Errorf("%w: BIC must be 9 characters long", legalentities.ErrInvalidBankAccountData)
	}

	if bankAccount.BankName == "" {
		return fmt.Errorf("%w: bank name is required", legalentities.ErrInvalidBankAccountData)
	}

	if bankAccount.SettlementAccount == "" {
		return fmt.Errorf("%w: settlement account is required", legalentities.ErrInvalidBankAccountData)
	}

	if len(bankAccount.SettlementAccount) != 20 {
		return fmt.Errorf("%w: settlement account must be 20 characters long", legalentities.ErrInvalidBankAccountData)
	}

	return nil
}

// PutBankAccountsUuid implements olegalentities.StrictServerInterface.
func (w *Web) PutBankAccountsUuid(ctx context.Context, request olegalentities.PutBankAccountsUuidRequestObject) (olegalentities.PutBankAccountsUuidResponseObject, error) {
	// Получаем текущий счет для проверки LegalEntityID
	existingAcc, err := w.app.LegalEntities.GetBankAccount(request.Uuid)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "bank account not found")
	}

	domainAcc := w.toBankAccountDomain(*request.Body, existingAcc.LegalEntityID)
	domainAcc.UUID = request.Uuid

	if err := w.app.LegalEntities.UpdateBankAccount(&domainAcc); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to update bank account").SetInternal(err)
	}

	return olegalentities.PutBankAccountsUuid200Response{}, nil
}

// DeleteBankAccountsUuid implements olegalentities.StrictServerInterface.
// Это ваша текущая реализация для StrictServerInterface (оставьте как есть)
func (w *Web) DeleteBankAccountsUuid(ctx context.Context, request olegalentities.DeleteBankAccountsUuidRequestObject) (olegalentities.DeleteBankAccountsUuidResponseObject, error) {
	err := w.app.LegalEntities.DeleteBankAccount(request.Uuid)
	if err != nil {
		if errors.Is(err, legalentities.ErrBankAccountNotFound) {
			return nil, echo.NewHTTPError(http.StatusNotFound, "bank account not found")
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to delete bank account").SetInternal(err)
	}
	return olegalentities.DeleteBankAccountsUuid204Response{}, nil
}

// Добавьте этот метод-адаптер для Echo
func (w *Web) DeleteBankAccountsUuidEcho(ctx echo.Context) error {
	var uuid openapi_types.UUID
	if err := runtime.BindStyledParameter("simple", false, "uuid", ctx.Param("uuid"), &uuid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID format")
	}

	request := olegalentities.DeleteBankAccountsUuidRequestObject{Uuid: uuid}
	response, err := w.DeleteBankAccountsUuid(ctx.Request().Context(), request)
	if err != nil {
		return err
	}
	return response.VisitDeleteBankAccountsUuidResponse(ctx.Response())
}

// Типы для обработки DELETE запроса банковских счетов
type DeleteBankAccountsUuidRequestObject struct {
	Uuid uuid.UUID `json:"uuid"`
}

type DeleteBankAccountsUuidResponseObject interface {
	VisitDeleteBankAccountsUuidResponse(ctx context.Context, w http.ResponseWriter) error
}

type DeleteBankAccountsUuid204Response struct{}

func (response DeleteBankAccountsUuid204Response) VisitDeleteBankAccountsUuidResponse(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// DeleteLegalEntitiesUuid implements ofederation.StrictServerInterface.
//
//nolint:revive,stylecheck // метод сгенерирован автоматически и соответствует OpenAPI
func (a *Web) DeleteLegalEntitiesUuid(_ context.Context, req olegalentities.DeleteLegalEntitiesUuidRequestObject) (olegalentities.DeleteLegalEntitiesUuidResponseObject, error) {
	err := a.app.LegalEntities.DeleteLegalEntity(req.Uuid)
	if err != nil {
		return nil, err
	}

	return olegalentities.DeleteLegalEntitiesUuid204Response{}, nil
}

// PutLegalEntitiesUuid implements ofederation.StrictServerInterface.
//
//nolint:revive,stylecheck // метод сгенерирован автоматически и соответствует OpenAPI
func (a *Web) PutLegalEntitiesUuid(_ context.Context, req olegalentities.PutLegalEntitiesUuidRequestObject) (olegalentities.PutLegalEntitiesUuidResponseObject, error) {
	uuidParsed, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	updated := &domain.LegalEntity{
		UUID: uuidParsed, // Теперь типы совпадают
		Name: req.Body.Name,
	}

	err = a.app.LegalEntities.UpdateLegalEntity(updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update legal entity: %w", err)
	}

	return olegalentities.PutLegalEntitiesUuid200JSONResponse{
		Uuid:      lo.ToPtr(updated.UUID.String()), // Теперь UUID имеет метод String()
		Name:      updated.Name,
		CreatedAt: &updated.CreatedAt,
		UpdatedAt: &updated.UpdatedAt,
	}, nil
}

// GetLegalEntities implements ofederation.StrictServerInterface.
func (a *Web) GetLegalEntities(_ context.Context, _ olegalentities.GetLegalEntitiesRequestObject) (olegalentities.GetLegalEntitiesResponseObject, error) {
	entities, err := a.app.LegalEntities.GetAllLegalEntities()
	if err != nil {
		return nil, err
	}

	response := lo.Map(entities, func(e domain.LegalEntity, _ int) olegalentities.LegalEntityDTO {
		return olegalentities.LegalEntityDTO{
			Uuid:      lo.ToPtr(e.UUID.String()), // UUID теперь имеет метод String()
			Name:      e.Name,
			CreatedAt: &e.CreatedAt,
			UpdatedAt: &e.UpdatedAt,
		}
	})

	return olegalentities.GetLegalEntities200JSONResponse(response), nil
}

// PostLegalEntities implements olegalentities.StrictServerInterface.
func (a *Web) PostLegalEntities(ctx context.Context, req olegalentities.PostLegalEntitiesRequestObject) (olegalentities.PostLegalEntitiesResponseObject, error) {
	dto := req.Body

	// создаем сущность из DTO
	entity := &domain.LegalEntity{
		Name: dto.Name,
		Meta: domain.JSONB{},
	}

	// сохраняем в базе
	err := a.app.LegalEntities.CreateLegalEntity(entity)
	if err != nil {
		return nil, err
	}

	// создаем правильный ответ - LegalEntityCreateDTO, а не LegalEntityDTO
	resp := olegalentities.LegalEntityCreateDTO{
		Name: entity.Name,
		// обработай meta соответствующим образом
		Meta: nil, // или преобразуй entity.Meta в *map[string]interface{}
	}

	return olegalentities.PostLegalEntities201JSONResponse(resp), nil
}

// Преобразование domain -> DTO (для ответов API)
func (w *Web) toBankAccountDTO(domainAcc domain.BankAccount) olegalentities.BankAccountDTO {
	return olegalentities.BankAccountDTO{
		Uuid:                 domainAcc.UUID,
		Bik:                  domainAcc.BIC,
		Bank:                 domainAcc.BankName,
		Address:              &domainAcc.BankAddress,
		CorrespondentAccount: &domainAcc.CorrespondentAccount,
		CheckingAccount:      domainAcc.SettlementAccount,
		Currency:             &domainAcc.Currency,
		Comment:              &domainAcc.Comment,
		IsPrimary:            &domainAcc.IsPrimary,
	}
}

// Преобразование DTO -> domain (для входящих запросов)
func (w *Web) toBankAccountDomain(dto olegalentities.BankAccountDTO, legalEntityID uuid.UUID) domain.BankAccount {
	return domain.BankAccount{
		UUID:                 dto.Uuid,
		LegalEntityID:        legalEntityID,
		BIC:                  dto.Bik,
		BankName:             dto.Bank,
		BankAddress:          getString(dto.Address),
		CorrespondentAccount: getString(dto.CorrespondentAccount),
		SettlementAccount:    dto.CheckingAccount,
		Currency:             getString(dto.Currency),
		Comment:              getString(dto.Comment),
		IsPrimary:            getBool(dto.IsPrimary),
	}
}
func getString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func getBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func NewWeb(conf configs.Configs) *Web {
	name := helpers.FakeName()

	a, err := app.InitApp(name, conf.DB_CREDS, true, conf.REDIS_CREDS)
	if err != nil {
		logrus.Fatal(err)
	}

	return &Web{
		app:     a,
		Options: conf,
		Now:     helpers.DateNow(),
		UUID:    name,

		Port: conf.PORT,
	}
}

func (a *Web) Work(ctx context.Context, rds *redis.RDS) {
	a.app.Work(ctx, rds)
	a.app.Subscribe(ctx)
}

var upgrader = websocket.Upgrader{}

func hello(a *Web, _ *echo.Echo) func(c echo.Context) error {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				logrus.Error(err)
				continue
			}

			fmt.Printf("%s\n", msg)

			arr := strings.Split(string(msg), " ")

			if len(arr) < 2 {
				continue
			}

			search := domain.SearchUser{
				FederationUUID: uuid.MustParse(arr[0]),
				Search:         arr[1],
			}

			dmns, err := a.app.FederationService.SearchUserInDictionary(search)
			if err != nil {
				logrus.Error(err)
				continue
			}

			dtos := lo.Map(dmns, func(item domain.User, _ int) dto.UserDTO {
				return dto.NewUserDto(item, a.app.ProfileService)
			})

			jsn, err := json.Marshal(dtos)
			if err != nil {
				logrus.Error(err)
				continue
			}

			err = ws.WriteMessage(websocket.TextMessage, jsn)
			if err != nil {
				logrus.Error(err)
				continue
			}
		}
	}
}

func (a *Web) Init() *echo.Echo {
	e := echo.New()

	// Middlewares
	if a.Options.CORS_ENABLE {
		origins := strings.Split(a.Options.CORS_ALLOWED_ORIGINS, ",")

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     origins,
			AllowCredentials: a.Options.CORS_ALLOW_CREDENTIALS,
			AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
	}

	if a.Options.OTEL_ENABLE {
		e.Use(TraceMiddleware("crm", a.Options.OTEL_EXPORTER, a.Options.ENV, WithSkipper(middleware.DefaultSkipper)))
	}

	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("150M"))
	e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().RequestURI, "/comment") {
				return true
			}

			if strings.Contains(c.Request().RequestURI, "/profile/photo") {
				return true
			}

			return false
		},
		Limit: "2M",
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.BodyDump(LogMiddleware(a.app)))

	if a.Options.GZIP > 0 {
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: a.app.Options.GZIP,
			Skipper: func(c echo.Context) bool {
				return c.Request().RequestURI == "/metrics"
			},
		}))
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var notFoundErr dto.NotFoundError
		if errors.As(err, &notFoundErr) {
			//nolint
			c.JSON(http.StatusNotFound, RequestError{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
			return
		}

		if errors.Is(err, ErrUnauthorized) {
			//nolint
			c.JSON(http.StatusUnauthorized, RequestError{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		// check if error is known type to be handled differently
		var myErr *ValidationError
		if errors.As(err, &myErr) {
			//nolint
			c.JSON(http.StatusBadRequest, ValidationError{
				StatusCode: http.StatusBadRequest,
				Errors:     myErr.Errors,
			})
			return
		}

		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			message, err := httpError.Message.(string)
			if !err {
				message = "Unknown (not string) error"
			}

			//nolint
			c.JSON(http.StatusBadRequest, RequestError{
				StatusCode: httpError.Code,
				Message:    message,
			})
			return
		}

		//nolint
		c.JSON(http.StatusConflict, RequestError{
			StatusCode: http.StatusConflict,
			Message:    err.Error(),
		})

		e.DefaultHTTPErrorHandler(err, c)
	}

	// Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Global rate limiter
	if a.Options.RATE_LIMITER > 0 {
		rateMinimum := rate.Limit(a.Options.RATE_LIMITER)
		rateMaximum := a.app.Options.RATE_LIMITER * 2

		config := middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: rateMinimum, Burst: rateMaximum, ExpiresIn: 1 * time.Minute},
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				id := ctx.RealIP()
				return id, nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}

		e.Use(middleware.RateLimiterWithConfig(config))
	}

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			spew.Dump(values.Error)
			if values.Error != nil {
				msg := fmt.Sprintf("[error:%s] echo request error", values.Error.Error())

				logrus.WithFields(logrus.Fields{
					"uri":     values.URI,
					"status":  values.Status,
					"latency": values.Latency.Nanoseconds(),
					"ip":      values.RemoteIP,
				}).Error(msg)
			} else {
				msg := "request: " + values.URI

				logrus.WithFields(logrus.Fields{
					"uri":     values.URI,
					"status":  values.Status,
					"latency": values.Latency.Nanoseconds(),
					"ip":      values.RemoteIP,
				}).Info(msg)
			}

			return nil
		},
	}))

	// Routers
	initMetricsRoutes(a, e)
	initOpenAPIProfileRouters(a, e)
	initOpenAPIMainRouters(a, e)
	initOpenAPIFederationRouters(a, e)
	initOpenAPIProjectRouters(a, e)
	initOpenAPITaskRouters(a, e)
	initOpenAPIReminderRouters(a, e)
	initOpenAPIcatalogRouters(a, e)
	olegalentities.RegisterHandlersWithBaseURL(e, olegalentities.NewStrictHandler(a, nil), "/api")
	e.DELETE("/api/bank-accounts/:uuid", a.DeleteBankAccountsUuidEcho)
	e.File("/openapi.yaml", "./openapi.yaml", middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "pong")
	})

	e.GET("/ws", hello(a, e))

	e.GET("/seed", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		c.Response().WriteHeader(http.StatusOK)

		i := 0
		ch := make(chan string, 100)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("exception: %s", string(debug.Stack()))

					msg2 := fmt.Sprintf("id: %v\nevent: %s\ndata: {'msg':%s}\n\n", i, "seed", "error")
					fmt.Fprint(c.Response(), msg2)

					close(ch)
					return
				}
			}()

			usersCount := helpers.MustInt(c.QueryParam("usersCount"))
			projectsCount := helpers.MustInt(c.QueryParam("projectsCount"))
			cores := helpers.MustInt(c.QueryParam("cores"))
			tasksCountPerCore := helpers.MustInt(c.QueryParam("tasksCountPerCore"))
			batch := helpers.MustInt(c.QueryParam("batch"))

			err := a.app.Seed(ch, usersCount, projectsCount, cores, tasksCountPerCore, batch)
			if err != nil {
				logrus.Error(err)
			}
		}()

		for {
			// check chan close
			if v, ok := <-ch; ok {
				msg := v
				i++

				msg2 := fmt.Sprintf("id: %v\nevent: %s\ndata: {'msg':%s}\n\n", i, "seed", msg)
				fmt.Fprint(c.Response(), msg2)
				c.Response().Flush()
			} else {
				break
			}
		}

		return nil
	})

	e.DELETE("/bank-accounts/:uuid", func(c echo.Context) error {
		uuidStr := c.Param("uuid") // Получаем строку из параметра

		// Конвертируем строку в uuid.UUID
		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid UUID format"})
		}

		err = a.app.LegalEntities.DeleteBankAccount(uuid)
		if err != nil {
			if errors.Is(err, legalentities.ErrBankAccountNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "bank account not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/seed_task", func(c echo.Context) error {
		total := helpers.MustInt(c.QueryParam("total"))
		projectUUID := uuid.MustParse(c.QueryParam("project_uuid"))
		createdBy := c.QueryParam("created_by")
		randomImplemented := c.QueryParam("random_implemented") == "true"
		commentsMax := helpers.MustInt(c.QueryParam("comments_max"))

		if total > 1000 {
			return errors.New("total must be < 1000")
		}

		dmns, err := a.app.SeedTasks(c.Request().Context(), total, projectUUID, createdBy, randomImplemented, commentsMax)
		if err != nil {
			logrus.Error(err)
			return err
		}

		err = c.JSON(http.StatusOK, dmns)

		return err
	})

	a.Router = e

	strictHandler := olegalentities.NewStrictHandler(a, nil)
	olegalentities.RegisterHandlersWithBaseURL(e, strictHandler, "/api")
	return e
}

func (a *Web) Run() {
	go func() {
		if err := a.Router.Start(fmt.Sprintf(":%d", a.Port)); err != nil && errors.Is(err, http.ErrServerClosed) {
			a.Router.Logger.Fatal("🙏 shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Router.Shutdown(ctx); err != nil {
		a.Router.Logger.Fatal(err)
	}
}
