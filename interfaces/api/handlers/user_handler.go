package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	user, err := h.userService.Register(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Registration failed", err)
	}

	// Get user stats for register response
	topicCount, videoCount, _ := h.userService.GetUserStats(c.Context(), user.ID)
	userResponse := dto.UserToUserResponse(user, topicCount, videoCount, false, false)
	return utils.SuccessResponse(c, "User registered successfully", userResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	token, user, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Login failed", err)
	}

	// Get user stats for login response
	topicCount, videoCount, _ := h.userService.GetUserStats(c.Context(), user.ID)
	loginResponse := &dto.LoginResponse{
		Token: token,
		User:  *dto.UserToUserResponse(user, topicCount, videoCount, false, false),
	}
	return utils.SuccessResponse(c, "Login successful", loginResponse)
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	profile, err := h.userService.GetProfileWithStats(c.Context(), user.ID, &user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	_, err = h.userService.UpdateProfile(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Profile update failed", err)
	}

	// Get updated profile with stats
	profile, err := h.userService.GetProfileWithStats(c.Context(), user.ID, &user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	return utils.SuccessResponse(c, "Profile updated successfully", profile)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.userService.DeleteUser(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User deletion failed", err)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	offsetStr := c.Query("offset", "0")
	limitStr := c.Query("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid offset parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid limit parameter")
	}

	users, total, err := h.userService.ListUsers(c.Context(), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	userResponses := make([]dto.UserResponseAdmin, len(users))
	for i, user := range users {
		topicCount, videoCount, _ := h.userService.GetUserStats(c.Context(), user.ID)
		userResponses[i] = *dto.UserToUserResponseAdmin(user, topicCount, videoCount)
	}

	response := &dto.UserListResponse{
		Users: userResponses,
		Meta:  dto.NewPaginationMeta(total, offset, limit),
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", response)
}
