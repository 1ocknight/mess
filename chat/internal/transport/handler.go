package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/1ocknight/mess/chat/internal/domain"
	httpdto "github.com/1ocknight/mess/shared/dto/http"
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

func (h *Handler) AddChat(c *gin.Context) {
	var req httpdto.AddChatRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	chat, err := h.domain.AddChat(c.Request.Context(), req.SecondSubjectID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, ChatMetadataModelToDTO(chat))
}

func (h *Handler) GetChatBySubjectID(c *gin.Context) {
	id := c.Param("subject_id")
	if id == "" {
		h.sendError(c, InvalidRequestError)
		return
	}

	chat, err := h.domain.GetChatMetadataBySubjectID(c.Request.Context(), id)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, ChatMetadataModelToDTO(chat))
}

func (h *Handler) GetChatByID(c *gin.Context) {
	var err error

	sChatID := c.Param("chat_id")
	if sChatID == "" {
		h.sendError(c, InvalidRequestError)
		return
	}
	chatID, err := strconv.Atoi(sChatID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	chat, err := h.domain.GetChatMetadataByID(c.Request.Context(), chatID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, ChatMetadataModelToDTO(chat))
}

func (h *Handler) GetChats(c *gin.Context) {
	sLimit := c.Query("limit")
	sBefore := c.Query("before")
	sAfter := c.Query("after")

	filter, err := MakeChatPaginationFilter(sLimit, sBefore, sAfter)
	if err != nil {
		h.sendError(c, err)
		return
	}

	chatsMetadata, err := h.domain.GetChatsMetadata(c.Request.Context(), filter)
	if err != nil {
		h.sendError(c, err)
		return
	}
	resChats := ChatsMetadataModelToDTO(chatsMetadata)

	c.JSON(http.StatusOK, resChats)
}

func (h *Handler) GetMessages(c *gin.Context) {
	sChatID := c.Param("chat_id")
	if sChatID == "" {
		h.sendError(c, InvalidRequestError)
		return
	}

	chatID, err := strconv.Atoi(sChatID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	sLimit := c.Query("limit")
	sBefore := c.Query("before")
	sAfter := c.Query("after")

	filter, err := MakeMessagePaginationFilter(sLimit, sBefore, sAfter)
	if err != nil {
		h.sendError(c, err)
		return
	}

	messages, err := h.domain.GetMessages(c.Request.Context(), chatID, filter)
	if err != nil {
		h.sendError(c, err)
		return
	}
	resMessages := MessagesModelToMessageDTO(messages)

	c.JSON(http.StatusOK, resMessages)
}

func (h *Handler) AddMessage(c *gin.Context) {
	sChatID := c.Param("chat_id")
	if sChatID == "" {
		h.sendError(c, InvalidRequestError)
		return
	}

	chatID, err := strconv.Atoi(sChatID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	var req httpdto.AddMessageRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	mess, err := h.domain.SendMessage(c.Request.Context(), chatID, req.Content)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusCreated, MessageModelToMessageDTO(mess))
}

func (h *Handler) UpdateMessage(c *gin.Context) {
	sMessageID := c.Param("message_id")
	if sMessageID == "" {
		h.sendError(c, InvalidRequestError)
		return
	}

	messageID, err := strconv.Atoi(sMessageID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	var req httpdto.UpdateMessageRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	mess, err := h.domain.UpdateMessage(c.Request.Context(), messageID, req.Content, req.Version)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.JSON(http.StatusOK, MessageModelToMessageDTO(mess))
}

func (h *Handler) UpdateLastRead(c *gin.Context) {
	sChatID := c.Param("chat_id")
	if sChatID == "" {
		h.sendError(c, InvalidRequestError)
		return
	}

	chatID, err := strconv.Atoi(sChatID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	var req httpdto.UpdateLastReadRequest
	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, err)
		return
	}

	_, err = h.domain.UpdateLastRead(c.Request.Context(), chatID, req.MessageID)
	if err != nil {
		h.sendError(c, err)
		return
	}

	c.Status(http.StatusOK)
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
