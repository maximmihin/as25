// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package http

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for PVZAvailableCity.
const (
	Казань         PVZAvailableCity = "Казань"
	Москва         PVZAvailableCity = "Москва"
	СанктПетербург PVZAvailableCity = "Санкт-Петербург"
)

// Defines values for PVZUserRole.
const (
	Employee  PVZUserRole = "employee"
	Moderator PVZUserRole = "moderator"
)

// Defines values for ProductType.
const (
	Обувь       ProductType = "обувь"
	Одежда      ProductType = "одежда"
	Электроника ProductType = "электроника"
)

// Defines values for ReceptionStatus.
const (
	Close      ReceptionStatus = "close"
	InProgress ReceptionStatus = "in_progress"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// PVZ defines model for PVZ.
type PVZ struct {
	City             PVZAvailableCity   `json:"city"`
	Id               openapi_types.UUID `json:"id"`
	RegistrationDate time.Time          `json:"registrationDate"`
}

// PVZAvailableCity defines model for PVZAvailableCity.
type PVZAvailableCity string

// PVZUserRole defines model for PVZUserRole.
type PVZUserRole string

// Product defines model for Product.
type Product struct {
	AddDate     time.Time          `json:"addDate"`
	Id          openapi_types.UUID `json:"id"`
	ReceptionId openapi_types.UUID `json:"receptionId"`
	Type        ProductType        `json:"type"`
}

// ProductType defines model for ProductType.
type ProductType string

// Reception defines model for Reception.
type Reception struct {
	Id        openapi_types.UUID `json:"id"`
	PvzId     openapi_types.UUID `json:"pvzId"`
	StartDate time.Time          `json:"startDate"`
	Status    ReceptionStatus    `json:"status"`
}

// ReceptionStatus defines model for Reception.Status.
type ReceptionStatus string

// Token defines model for Token.
type Token = string

// User defines model for User.
type User struct {
	Email openapi_types.Email `json:"email"`
	Id    openapi_types.UUID  `json:"id"`
	Role  PVZUserRole         `json:"role"`
}

// N400BadRequest defines model for 400BadRequest.
type N400BadRequest = Error

// N401UnauthorizedResponse defines model for 401UnauthorizedResponse.
type N401UnauthorizedResponse = Error

// N403ForbiddenResponse defines model for 403ForbiddenResponse.
type N403ForbiddenResponse = Error

// PostDummyLoginJSONBody defines parameters for PostDummyLogin.
type PostDummyLoginJSONBody struct {
	Role PVZUserRole `json:"role"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
}

// PostProductsJSONBody defines parameters for PostProducts.
type PostProductsJSONBody struct {
	PvzId openapi_types.UUID `json:"pvzId"`
	Type  ProductType        `json:"type"`
}

// GetPvzParams defines parameters for GetPvz.
type GetPvzParams struct {
	// StartDate Начальная дата диапазона
	StartDate *time.Time `form:"startDate,omitempty" json:"startDate,omitempty"`

	// EndDate Конечная дата диапазона
	EndDate *time.Time `form:"endDate,omitempty" json:"endDate,omitempty"`

	// Page Номер страницы
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Limit Количество элементов на странице
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostPvzJSONBody defines parameters for PostPvz.
type PostPvzJSONBody struct {
	City PVZAvailableCity `json:"city"`
}

// PostReceptionsJSONBody defines parameters for PostReceptions.
type PostReceptionsJSONBody struct {
	PvzId openapi_types.UUID `json:"pvzId"`
}

// PostRegisterJSONBody defines parameters for PostRegister.
type PostRegisterJSONBody struct {
	Email    openapi_types.Email `json:"email"`
	Password string              `json:"password"`
	Role     PVZUserRole         `json:"role"`
}

// PostDummyLoginJSONRequestBody defines body for PostDummyLogin for application/json ContentType.
type PostDummyLoginJSONRequestBody PostDummyLoginJSONBody

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostProductsJSONRequestBody defines body for PostProducts for application/json ContentType.
type PostProductsJSONRequestBody PostProductsJSONBody

// PostPvzJSONRequestBody defines body for PostPvz for application/json ContentType.
type PostPvzJSONRequestBody PostPvzJSONBody

// PostReceptionsJSONRequestBody defines body for PostReceptions for application/json ContentType.
type PostReceptionsJSONRequestBody PostReceptionsJSONBody

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody PostRegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Получение тестового токена
	// (POST /dummyLogin)
	PostDummyLogin(c *fiber.Ctx) error
	// Авторизация пользователя
	// (POST /login)
	PostLogin(c *fiber.Ctx) error
	// Добавление товара в текущую приемку (только для сотрудников ПВЗ)
	// (POST /products)
	PostProducts(c *fiber.Ctx) error
	// Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
	// (GET /pvz)
	GetPvz(c *fiber.Ctx, params GetPvzParams) error
	// Создание ПВЗ (только для модераторов)
	// (POST /pvz)
	PostPvz(c *fiber.Ctx) error
	// Закрытие последней открытой приемки товаров в рамках ПВЗ
	// (POST /pvz/{pvzId}/close_last_reception)
	PostPvzPvzIdCloseLastReception(c *fiber.Ctx, pvzId openapi_types.UUID) error
	// Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
	// (POST /pvz/{pvzId}/delete_last_product)
	PostPvzPvzIdDeleteLastProduct(c *fiber.Ctx, pvzId openapi_types.UUID) error
	// Создание новой приемки товаров (только для сотрудников ПВЗ)
	// (POST /receptions)
	PostReceptions(c *fiber.Ctx) error
	// Регистрация пользователя
	// (POST /register)
	PostRegister(c *fiber.Ctx) error
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

type MiddlewareFunc fiber.Handler

// PostDummyLogin operation middleware
func (siw *ServerInterfaceWrapper) PostDummyLogin(c *fiber.Ctx) error {

	return siw.Handler.PostDummyLogin(c)
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(c *fiber.Ctx) error {

	return siw.Handler.PostLogin(c)
}

// PostProducts operation middleware
func (siw *ServerInterfaceWrapper) PostProducts(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.PostProducts(c)
}

// GetPvz operation middleware
func (siw *ServerInterfaceWrapper) GetPvz(c *fiber.Ctx) error {

	var err error

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPvzParams

	var query url.Values
	query, err = url.ParseQuery(string(c.Request().URI().QueryString()))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for query string: %w", err).Error())
	}

	// ------------- Optional query parameter "startDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "startDate", query, &params.StartDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter startDate: %w", err).Error())
	}

	// ------------- Optional query parameter "endDate" -------------

	err = runtime.BindQueryParameter("form", true, false, "endDate", query, &params.EndDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter endDate: %w", err).Error())
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", query, &params.Page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter page: %w", err).Error())
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", query, &params.Limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter limit: %w", err).Error())
	}

	return siw.Handler.GetPvz(c, params)
}

// PostPvz operation middleware
func (siw *ServerInterfaceWrapper) PostPvz(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.PostPvz(c)
}

// PostPvzPvzIdCloseLastReception operation middleware
func (siw *ServerInterfaceWrapper) PostPvzPvzIdCloseLastReception(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", c.Params("pvzId"), &pvzId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter pvzId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.PostPvzPvzIdCloseLastReception(c, pvzId)
}

// PostPvzPvzIdDeleteLastProduct operation middleware
func (siw *ServerInterfaceWrapper) PostPvzPvzIdDeleteLastProduct(c *fiber.Ctx) error {

	var err error

	// ------------- Path parameter "pvzId" -------------
	var pvzId openapi_types.UUID

	err = runtime.BindStyledParameterWithOptions("simple", "pvzId", c.Params("pvzId"), &pvzId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Errorf("Invalid format for parameter pvzId: %w", err).Error())
	}

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.PostPvzPvzIdDeleteLastProduct(c, pvzId)
}

// PostReceptions operation middleware
func (siw *ServerInterfaceWrapper) PostReceptions(c *fiber.Ctx) error {

	c.Context().SetUserValue(BearerAuthScopes, []string{})

	return siw.Handler.PostReceptions(c)
}

// PostRegister operation middleware
func (siw *ServerInterfaceWrapper) PostRegister(c *fiber.Ctx) error {

	return siw.Handler.PostRegister(c)
}

// FiberServerOptions provides options for the Fiber server.
type FiberServerOptions struct {
	BaseURL     string
	Middlewares []MiddlewareFunc
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router fiber.Router, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, FiberServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router fiber.Router, si ServerInterface, options FiberServerOptions) {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	for _, m := range options.Middlewares {
		router.Use(fiber.Handler(m))
	}

	router.Post(options.BaseURL+"/dummyLogin", wrapper.PostDummyLogin)

	router.Post(options.BaseURL+"/login", wrapper.PostLogin)

	router.Post(options.BaseURL+"/products", wrapper.PostProducts)

	router.Get(options.BaseURL+"/pvz", wrapper.GetPvz)

	router.Post(options.BaseURL+"/pvz", wrapper.PostPvz)

	router.Post(options.BaseURL+"/pvz/:pvzId/close_last_reception", wrapper.PostPvzPvzIdCloseLastReception)

	router.Post(options.BaseURL+"/pvz/:pvzId/delete_last_product", wrapper.PostPvzPvzIdDeleteLastProduct)

	router.Post(options.BaseURL+"/receptions", wrapper.PostReceptions)

	router.Post(options.BaseURL+"/register", wrapper.PostRegister)

}
