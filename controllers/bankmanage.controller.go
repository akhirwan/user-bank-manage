package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"user-bank-manage/models"

	"github.com/gofiber/fiber/v2"
)

func FetchAllBanks(c *fiber.Ctx) error {
	result, err := models.FetchAllBanks()
	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return c.Status(400).JSON(http.StatusInternalServerError)
	}

	return c.Status(200).JSON(result)
}

func DetailBanks(c *fiber.Ctx) error {
	bank_id := c.FormValue("id")

	conv_id, err := strconv.Atoi(bank_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	result, err := models.DetailBanks(conv_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

func StoreBanks(c *fiber.Ctx) error {
	now := time.Now()

	bank_identifier := c.FormValue("bank_identifier")
	// bank_active := c.FormValue("bank_active")
	bank_active := "1"
	// bank_added_on := c.FormValue("bank_added_on")
	bank_added_on := now.String()
	// bank_deleted := c.FormValue("bank_deleted")
	bank_deleted := "0"

	result, err := models.StoreBanks(bank_identifier, bank_active, bank_added_on, bank_deleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

func UpdateBanks(c *fiber.Ctx) error {
	// now := time.Now()

	bank_id := c.FormValue("bank_id")
	bank_identifier := c.FormValue("bank_identifier")
	// bank_active := "1"
	// bank_added_on := c.FormValue("bank_added_on")
	// bank_added_on := now.String()
	// bank_deleted := "0"

	bank_active := c.FormValue("bank_active")
	// bank_deleted := c.FormValue("bank_deleted")

	conv_active, err := strconv.Atoi(bank_active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	conv_id, err := strconv.Atoi(bank_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	// result, err := models.UpdateBanks(conv_id, bank_identifier, bank_active, bank_added_on, bank_deleted)
	result, err := models.UpdateBanks(conv_id, bank_identifier, conv_active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

func DeleteBanks(c *fiber.Ctx) error {
	bank_id := c.FormValue("bank_id")

	conv_id, err := strconv.Atoi(bank_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	result, err := models.DeleteBanks(conv_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

type RawBank struct {
	BankId         int    `json:"bank_id" xml:"bank_id" form:"bank_id"`
	BankIdentifier string `json:"bank_identifier" xml:"bank_identifier" form:"bank_identifier"`
	BankActive     int    `json:"bank_active" xml:"bank_active" form:"bank_active"`
	BankDeleted    int    `json:"bank_deleted" xml:"bank_deleted" form:"bank_deleted"`
}

func StoreBanksJson(c *fiber.Ctx) error {
	b := new(RawBank)

	now := time.Now()

	if err := c.BodyParser(b); err != nil {
		return err
	}

	bank_identifier := b.BankIdentifier
	bank_active := "1"
	bank_added_on := now.String()
	bank_deleted := "0"

	result, err := models.StoreBanks(bank_identifier, bank_active, bank_added_on, bank_deleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

func UpdateBanksJson(c *fiber.Ctx) error {
	b := new(RawBank)

	if err := c.BodyParser(b); err != nil {
		return err
	}

	bank_id := b.BankId
	bank_identifier := b.BankIdentifier
	bank_active := b.BankActive

	// conv_id, err := strconv.Atoi(bank_id)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError)
	// }

	result, err := models.UpdateBanks(bank_id, bank_identifier, bank_active)
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	return c.JSON(result)
}

// func (c *Ctx) FormFile(key string) (*multipart.FileHeader, error)

func UploadImage(c *fiber.Ctx) error {
	// Get first file from form field "document":
	file, err := c.FormFile("file_image")
	if err != nil {
		return c.JSON(http.StatusInternalServerError)
	}

	// Save file to root directory:
	return c.SaveFile(file, fmt.Sprintf("./file-images/%s", file.Filename))
}

func ListObjectStorage(c *fiber.Ctx) error {
	result, err := models.ListObjectStorage()
	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return c.Status(400).JSON(http.StatusInternalServerError)
	}

	return c.Status(200).JSON(result)
}

func UploadObjectStorage(c *fiber.Ctx) error {
	file, err := c.FormFile("file_image")
	if err != nil {
		// return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return err
	}

	// Save file to root directory:
	// return c.SaveFile(file, fmt.Sprintf("./file-images/%s", file.Filename))
	result, err := models.UploadObjectStorage(file.Filename)
	if err != nil {
		return c.Status(404).JSON(http.StatusInternalServerError)
		// return err
	}

	return c.Status(200).JSON(result)
}

func MultipartObjectStorage(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		if token := form.Value["token"]; len(token) > 0 {
			// Get key value:
			fmt.Println(token[0])
		}

		// Get all files from "documents" key:
		files := form.File["file_image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			tmpDir := fmt.Sprintf("./file-images/%s", file.Filename)
			fmt.Println(tmpDir, os.TempDir())

			// if err := c.SaveFile(file, tmpDir); err != nil {
			// 	return err
			// }
			// // return c.Status(200).JSON(err)

			_, b, _, _ := runtime.Caller(0)
			basepath := filepath.Dir(b)

			fmt.Println(basepath)

			result, err := models.MultipartObjectStorage(file.Filename, os.TempDir())
			if err != nil {
				return err
			}
			return c.Status(200).JSON(result)
		}
		return c.Status(403).JSON(http.StatusInternalServerError)
	}
	return c.Status(400).JSON(http.StatusInternalServerError)
}
