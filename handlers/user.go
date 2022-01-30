package handlers

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sajjadanwar0/clubhouse-clone/ent/user"
	"github.com/sajjadanwar0/clubhouse-clone/middleware"
	"github.com/sajjadanwar0/clubhouse-clone/utils"
	"net/http"
)

func (r registerRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Firstname, validation.Required, validation.Length(3, 20)),
		validation.Field(&r.Lastname, validation.Required, validation.Length(3, 20)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 12)),
	)
}

func (h *Handler) UserRegister(ctx *fiber.Ctx) error {
	var request registerRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		err = ctx.Status(http.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "Invalid Json",
		})

		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
	}

	if err = request.validate(); err != nil {
		ctx.Status(http.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": err,
		})
		return nil
	}

	exist, _ := h.Client.User.Query().Where(user.Email(request.Email)).Only(ctx.Context())

	if exist != nil {
		ctx.Status(http.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "The email is already taken",
		})
	}

	hashPassword, err := utils.HashPassword(request.Password)

	if err != nil {
		fmt.Errorf("failed hash user password: ", err)
		return nil
	}

	_, err = h.Client.User.Create().
		SetEmail(request.Email).
		SetFirstName(request.Firstname).
		SetLastName(request.Lastname).
		SetEmail(request.Email).
		SetAvatar(request.Avatar).
		SetPassword(hashPassword).
		Save(ctx.Context())

	if err != nil {
		ctx.Status(http.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "registered error",
		})
		return nil
	}

	_ = ctx.Status(http.StatusCreated).JSONP(fiber.Map{
		"error":   false,
		"message": "registered successfully",
	})

	return nil
}

func (h *Handler) UserLogin(ctx *fiber.Ctx) error {

	var request loginRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		err = ctx.Status(http.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "Invalid JSON",
		})
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		}
		return nil
	}

	u, err := h.Client.User.Query().Where(user.Email(request.Email)).Only(ctx.Context())

	if err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "Invalid user",
		})
		return nil
	}

	if err = utils.ComparePassword(request.Password, u.Password); err != nil {
		_ = ctx.Status(fiber.StatusBadRequest).JSONP(fiber.Map{
			"error":   true,
			"message": "Invalid credentials",
		})
		return nil
	}

	token, err := middleware.ClaimToken(u.ID)
	if err != nil {
		utils.Errorf("token generation error :", err)
		return nil
	}

	response := map[string]interface{}{
		"firstname": u.FirstName,
		"lastname":  u.LastName,
		"email":     u.Email,
		"avatar":    u.Avatar,
	}
	_ = ctx.Status(fiber.StatusOK).JSONP(fiber.Map{
		"error": false,
		"data":  response,
		"token": token,
	})

	return nil
}

func (h *Handler) MeQuery(ctx *fiber.Ctx) error {

	userId, err := middleware.GetUserIdFromContext(ctx)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return nil
	}

	uid, _ := uuid.Parse(userId)
	u, err := h.Client.User.Query().Where(user.ID(uid)).Only(ctx.Context())

	if err != nil {
		ctx.Status(fiber.StatusNotFound).JSONP(fiber.Map{
			"error":   true,
			"message": "cannot find the user",
		})
		return nil
	}

	response := map[string]interface{}{
		"firstname": u.FirstName,
		"lastname":  u.LastName,
		"email":     u.Email,
		"avatar":    u.Avatar,
	}

	_ = ctx.Status(fiber.StatusOK).JSONP(fiber.Map{
		"error": false,
		"data":  response,
	})

	return nil
}
