package http

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	product_repo "github.com/maximmihin/as25/internal/dal/repos/products"
	products_models "github.com/maximmihin/as25/internal/dal/repos/products/models"
	pvz_repo "github.com/maximmihin/as25/internal/dal/repos/pvz"
	pvz_models "github.com/maximmihin/as25/internal/dal/repos/pvz/models"
	receptions_repo "github.com/maximmihin/as25/internal/dal/repos/receptions"
	receptions_models "github.com/maximmihin/as25/internal/dal/repos/receptions/models"
	dal_types "github.com/maximmihin/as25/internal/dal/types"
	"github.com/maximmihin/as25/internal/logger"
	"github.com/maximmihin/as25/internal/services/accesscontrol"
)

//go:generate go tool oapi-codegen -config oapi_config.yaml ../../../api/api.yaml

type Server struct {
	JwtSecret string

	PvzRepo       *pvz_repo.Repo
	ReceptionRepo *receptions_repo.Repo
	ProductRepo   *product_repo.Repo
}

func (s Server) PostDummyLogin(c *fiber.Ctx) error {

	ctx := c.UserContext()

	var req PostDummyLoginJSONBody
	if err := c.BodyParser(&req); err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(400, "Invalid request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return logger.WrapError(ctx,
			fiber.NewError(400, err.Error()),
		)
	}

	jwtToken, err := accesscontrol.NewDummyJWT(s.JwtSecret, string(req.Role))
	if err != nil {
		if errors.Is(err, accesscontrol.ErrInvalidRole) {
			return logger.WrapErrorWithMessage(ctx, "the error passed through the request validator",
				fiber.NewError(400, "Invalid role"),
			)
		}
		return logger.WrapErrorWithMessage(ctx, fmt.Sprintf("unexpected error from accesscontrol.NewDummyJWT: %v", err),
			fiber.NewError(500),
		)
	}

	return c.Status(200).JSON(Token(jwtToken))
}

func (s Server) PostPvz(c *fiber.Ctx) error {

	ctx := c.UserContext()

	var req PostPvzJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(400, "Invalid request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return logger.WrapError(ctx,
			fiber.NewError(400, err.Error()),
		)
	}

	newPvz := pvz_models.Pvz{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		LocatedAt: dal_types.FtCity(req.City),
	}

	if err := s.PvzRepo.Create(ctx, newPvz); err != nil {
		if errors.Is(err, pvz_repo.ErrUnavailableCityType) {
			return logger.WrapErrorWithMessage(ctx,
				"triggered enum constraint (city), but validation but the validation above did not work", // TODO this must be marked as internal error
				fiber.NewError(400))
		}
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(500))
	}

	return c.Status(201).JSON(PVZ{
		City:             PVZAvailableCity(newPvz.LocatedAt),
		Id:               newPvz.ID,
		RegistrationDate: newPvz.CreatedAt,
	})
}

func (s Server) PostReceptions(c *fiber.Ctx) error {

	ctx := c.UserContext()

	var req PostReceptionsJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(400, "Invalid request body"),
		)
	}

	newReceptions := receptions_models.Reception{
		ID:        uuid.New(),
		PvzID:     req.PvzId,
		Status:    receptions_models.ReceptionProgressInProgress,
		CreatedAt: time.Now(),
	}

	if err := s.ReceptionRepo.Create(ctx, newReceptions); err != nil {
		if errors.Is(err, receptions_repo.ErrNonexistentPvzId) {
			return logger.WrapErrorWithMessage(ctx, err.Error(),
				fiber.NewError(400, "Pvz with such id does not exist"))
		}
		if errors.Is(err, receptions_repo.ErrThereAreOpenReceptions) {
			return logger.WrapErrorWithMessage(ctx, err.Error(),
				fiber.NewError(400, "There are open receptions"))
		}
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark us internal error
			fiber.NewError(500))
	}
	return c.Status(201).JSON(Reception{
		Id:        newReceptions.ID,
		PvzId:     newReceptions.PvzID,
		Status:    ReceptionStatus(newReceptions.Status),
		StartDate: newReceptions.CreatedAt,
	})

}

func (s Server) PostProducts(c *fiber.Ctx) error {

	ctx := c.UserContext()

	var req PostProductsJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(400, "Invalid request body"),
		)
	}

	if err := req.Validate(); err != nil {
		return logger.WrapError(ctx,
			fiber.NewError(400, err.Error()),
		)
	}

	tx, err := s.ReceptionRepo.Begin(ctx)
	if err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark us internal error
			fiber.NewError(500))
	}
	defer tx.Rollback(ctx)

	activeRec, err := s.ReceptionRepo.WithTx(tx).GetActiveByPVZ(ctx, req.PvzId)
	if err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark us internal error
			fiber.NewError(500))
	}

	if activeRec == nil {
		return logger.WrapError(ctx,
			fiber.NewError(400, "no reception \"in_progress\""))
	}

	// TODO add it to log?
	newProduct := products_models.Product{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Type:        products_models.ProductsType(req.Type),
		ReceptionID: activeRec.ID,
	}

	if err = s.ProductRepo.WithTx(tx).Create(ctx, newProduct); err != nil {
		if errors.Is(err, product_repo.ErrNonexistentReceptionId) {
			return logger.WrapErrorWithMessage(ctx,
				fmt.Sprintf("reception id was found previous in this tx, but fk constraing work: %v", err), // TODO mark us internal error
				fiber.NewError(500))
		}
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark us internal error
			fiber.NewError(500))
	}

	if err := tx.Commit(ctx); err != nil {
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark us internal error
			fiber.NewError(500))
	}

	return c.Status(201).JSON(Product{
		AddDate:     newProduct.CreatedAt,
		Id:          newProduct.ID,
		ReceptionId: newProduct.ReceptionID,
		Type:        ProductType(newProduct.Type),
	})

}

func (s Server) PostPvzPvzIdDeleteLastProduct(c *fiber.Ctx, pvzId openapi_types.UUID) error {

	ctx := c.UserContext()

	if err := s.ProductRepo.DeleteLastInPvz(ctx, pvzId); err != nil {
		if errors.Is(err, product_repo.ErrNothingToDelete) {
			return logger.WrapErrorWithMessage(ctx, err.Error(),
				fiber.NewError(400, "Nothing to delete"))
		}
		return logger.WrapErrorWithMessage(ctx, err.Error(), // TODO mark as internal error
			fiber.NewError(500))
	}
	c.Status(204)
	return nil
}

func (s Server) PostPvzPvzIdCloseLastReception(c *fiber.Ctx, pvzId openapi_types.UUID) error {

	ctx := c.UserContext()

	rec, err := s.ReceptionRepo.CloseActive(ctx, pvzId)
	if err != nil {
		if errors.Is(err, receptions_repo.ErrNothingToClose) {
			return logger.WrapErrorWithMessage(ctx, err.Error(),
				fiber.NewError(400, "Nothing to close")) // TODO type extended
		}
		return logger.WrapErrorWithMessage(ctx, err.Error(),
			fiber.NewError(500))
	}
	return c.Status(200).JSON(Reception{
		Id:        rec.ID,
		PvzId:     rec.PvzID,
		StartDate: rec.CreatedAt,
		Status:    ReceptionStatus(rec.Status),
	})
}

// TODO implement me
func (s Server) GetPvz(c *fiber.Ctx, params GetPvzParams) error {

	ctx := c.UserContext()

	if err := params.WithDefaults().Validate(); err != nil {
		return logger.WrapErrorWithMessage(ctx, "invalid request params"+err.Error(),
			fiber.NewError(400, err.Error()))
	}

	return c.Status(400).JSON(map[string]any{
		"message": "GetPvz not implemented yet",
	})
}

// TODO implement me
func (s Server) PostRegister(c *fiber.Ctx) error {
	return c.Status(400).JSON(map[string]any{
		"message": "PostRegister not implemented yet",
	})
}

// TODO implement me
func (s Server) PostLogin(c *fiber.Ctx) error {
	return c.Status(400).JSON(map[string]any{
		"message": "PostRegister not implemented yet",
	})
}
