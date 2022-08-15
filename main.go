package main

import (
	"context"
	"emailAvailability/core"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/webhooks", func(ctx *fiber.Ctx) error {
		if challenge := ctx.Query("challenge"); challenge != "" {
			return ctx.Status(fiber.StatusOK).SendString(challenge)
		}

		return ctx.SendStatus(fiber.StatusBadRequest)
	})

	app.Post("/webhooks", func(ctx *fiber.Ctx) error {
		ctx.Accepts("application/json")

		signature := ctx.GetReqHeaders()["X-Nylas-Signature"]
		if err := core.CheckSignature(os.Getenv("NYLAS_CLIENT_SECRET"), signature, ctx.Body()); err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		w := new(core.WebhookRequest)
		if e := json.Unmarshal(ctx.Body(), w); e != nil {
			log.Println("Error parsing webhook body into WebhookRequest", e)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		w.LogInfo()

		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Post("/send", func(ctx *fiber.Ctx) error {
		ctx.Accepts("application/json")

		sendRes, injectErr := core.InjectAvailabilityAndSendMessage(ctx.Body(), ctx.GetReqHeaders())

		if injectErr != nil {
			res := core.MessageResponse{
				Success:      false,
				ErrorMessage: injectErr.Error(),
			}

			if reqErr, ok := injectErr.(*core.RequestError); ok && reqErr.StatusCode != 0 {
				ctx.Status(reqErr.StatusCode)
			} else if errors.Is(reqErr, context.DeadlineExceeded) {
				res.ErrorMessage = "API timeout calling Nylas"
			} else {
				ctx.Status(fiber.StatusInternalServerError)
			}

			b, _ := json.Marshal(res)
			return ctx.Send(b)
		}

		b, _ := json.Marshal(core.MessageResponse{
			Success: true,
			Data:    sendRes,
		})
		return ctx.Status(fiber.StatusOK).Send(b)
	})

	log.Fatal(app.Listen(":8000"))
}
