package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TATAROmangol/mess/profile/internal/domain"
	"github.com/TATAROmangol/mess/profile/internal/model"
	"github.com/TATAROmangol/mess/profile/pkg/dto"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	domain domain.Service
}

func NewHandler(domain domain.Service) *Handler {
	return &Handler{
		domain: domain,
	}
}

func (h *Handler) GetProfile(c *gin.Context) {
	id := c.Param("id")

	var profile *model.Profile
	var url string
	var err error

	if id == "" {
		profile, url, err = h.domain.GetCurrentProfile(c.Request.Context())
		if err != nil {
			h.sendError(c, err)
			return
		}
	}

	if id != "" {
		profile, url, err = h.domain.GetProfileFromSubjectID(c.Request.Context(), id)
		if err != nil {
			h.sendError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		SubjectID: profile.SubjectID,
		Alias:     profile.Alias,
		AvatarURL: url,
		Version:   profile.Version,
	})
}

func (h *Handler) GetProfiles(c *gin.Context) {
	alias := c.Query("alias")
	sLimit := c.Query("limit")
	before := c.Query("before")
	after := c.Query("after")

	if after != "" && before != "" {
		h.sendError(c, fmt.Errorf("%w, wait after or before null", InvalidRequestError))
	}

	var limit int
	var err error
	if sLimit != "" {
		limit, err = strconv.Atoi(sLimit)
		if err != nil {
			h.sendError(c, fmt.Errorf("%w, atoi: %w", InvalidRequestError, err))
		}
	}

	filter := domain.ProfilePaginationFilter{
		Limit: limit,
	}

	if after != "" {
		filter.Direction = domain.DirectionAfter
		filter.LastSubjectID = &after
	}

	if before != "" {
		filter.Direction = domain.DirectionBefore
		filter.LastSubjectID = &before
	}

	profiles, urls, err := h.domain.GetProfilesFromAlias(c.Request.Context(), alias, &filter)
	if err != nil {
		h.sendError(c, err)
		return
	}

	if len(profiles) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	res := make([]*dto.ProfileResponse, 0, len(profiles))
	for _, profile := range profiles {
		res = append(res, &dto.ProfileResponse{
			SubjectID: profile.SubjectID,
			Alias:     profile.Alias,
			AvatarURL: urls[profile.SubjectID],
			Version:   profile.Version,
		})
	}

	c.JSON(http.StatusOK, dto.ProfilesResponse{
		Profiles: res,
	})
}

func (h *Handler) AddProfile(c *gin.Context) {
	var req *dto.AddProfileRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	profile, url, err := h.domain.AddProfile(c.Request.Context(), req.Alias)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ProfileResponse{
		SubjectID: profile.SubjectID,
		Alias:     profile.Alias,
		AvatarURL: url,
		Version:   profile.Version,
	})
}

func (h *Handler) UpdateProfileMetadata(c *gin.Context) {
	var req *dto.UpdateProfileMetadataRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	profile, url, err := h.domain.UpdateProfileMetadata(c.Request.Context(), req.Version, req.Alias)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		SubjectID: profile.SubjectID,
		Alias:     profile.Alias,
		AvatarURL: url,
		Version:   profile.Version,
	})
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	url, err := h.domain.UploadAvatar(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UploadAvatarResponse{
		UploadURL: url,
	})
}

func (h *Handler) DeleteAvatar(c *gin.Context) {
	profile, url, err := h.domain.DeleteAvatar(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		SubjectID: profile.SubjectID,
		Alias:     profile.Alias,
		AvatarURL: url,
		Version:   profile.Version,
	})
}

func (h *Handler) DeleteProfile(c *gin.Context) {
	profile, url, err := h.domain.DeleteProfile(c.Request.Context())
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		SubjectID: profile.SubjectID,
		Alias:     profile.Alias,
		AvatarURL: url,
		Version:   profile.Version,
	})
}

func (h *Handler) sendError(c *gin.Context, err error) {
	var code int

	if errors.Is(err, InvalidRequestError) {
		code = http.StatusBadRequest
	}

	if errors.Is(err, domain.ErrNotFound) {
		code = http.StatusNoContent
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	c.AbortWithError(code, err)
}
