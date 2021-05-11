package api

import (
	"log"
	"net/http"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	js "github.com/lrweck/shawty/serializer/json"
	ms "github.com/lrweck/shawty/serializer/msgpack"
	sh "github.com/lrweck/shawty/shortener"
)

type RedirectHandler interface {
	Get(*fiber.Ctx) error
	Post(*fiber.Ctx) error
}

type handler struct {
	redirectService sh.RedirectService
}

func NewHandler(redirectService sh.RedirectService) RedirectHandler {
	return &handler{
		redirectService: redirectService,
	}
}

func setupResponse(f *fiber.Ctx, contentType string, statusCode int, body []byte) {
	f.Append("Content-Type", contentType)
	f.Status(statusCode)

	if err := f.Send(body); err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) sh.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.Redirect{}
	}
	return &js.Redirect{}
}

func (h *handler) Get(f *fiber.Ctx) error {
	code := f.Params("code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == sh.ErrRedirectNotFound {
			return fiber.ErrNotFound
		}
		return fiber.ErrInternalServerError
	}
	f.Redirect(redirect.URL, http.StatusMovedPermanently)
	return nil
}

func (h *handler) Post(f *fiber.Ctx) error {
	contentType := f.Get("Content-Type")
	requestBody := f.Body()
	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == sh.ErrRedirectInvalid {
			return fiber.ErrBadRequest
		}
		return fiber.ErrInternalServerError
	}

	responseBody, err := h.serializer(contentType).Encode(redirect)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	setupResponse(f, contentType, http.StatusCreated, responseBody)

	return nil
}
