package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"tes.go/m/database"
	"tes.go/m/model"
)

func CreateData(res *fiber.Ctx) error {
	db := database.DB.Db

	data := new(model.Data)

	form, err := res.MultipartForm()

	if err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed parsing input"})
	}

	var title string
	var desc string
	var activityDate string
	var linkMom string
	var location string
	var actor string

	if titleVal := form.Value["title"]; len(titleVal) > 0 {
		title = titleVal[0]
	}

	if descVal := form.Value["description"]; len(descVal) > 0 {
		desc = descVal[0]
	}

	if activityDateVal := form.Value["activity_date"]; len(activityDateVal) > 0 {
		activityDate = activityDateVal[0]
	}

	if linkMomVal := form.Value["link_mom"]; len(linkMomVal) > 0 {
		linkMom = linkMomVal[0]
	}

	if locationVal := form.Value["location"]; len(locationVal) > 0 {
		location = locationVal[0]
	}

	if actorVal := form.Value["actor"]; len(actorVal) > 0 {
		actor = actorVal[0]
	}

	file, err := res.FormFile("image")
	if err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed parsing image"})
	}

	src, err := file.Open()
	if err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed parsing image"})
	}
	defer src.Close()
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s", timestamp, file.Filename)

	dst, err := os.Create("./static/" + strings.Replace(fileName, " ", "_", -1))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed saving image"})
	}

	formattedDate, error := time.Parse("02-01-2006", activityDate)

	if error != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Wrong date format"})
	}

	data.ActivityDate = formattedDate
	data.Title = title
	data.Description = desc
	data.LinkMom = linkMom
	data.Location = location
	data.Actor = actor
	data.Image = res.Hostname() + "/images/" + strings.Replace(fileName, " ", "_", -1)

	err = db.Create(data).Error

	if err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed creating data"})
	}

	return res.Status(200).JSON(fiber.Map{"status": "success", "message": "User created successfully", "data": data})

}

func GetAllData(res *fiber.Ctx) error {
	db := database.DB.Db
	var datas []model.Data

	page, err := strconv.Atoi(res.Query("page"))
	if err != nil {
		page = 1 // Use default page if invalid
	}

	pageSize, err := strconv.Atoi(res.Query("pageSize"))
	if err != nil {
		pageSize = 10 // Use default page size if invalid
	}

	offset := (page - 1) * pageSize

	db.Limit(pageSize).Offset(offset).Find(&datas)

	if len(datas) == 0 {
		return res.Status(404).JSON(fiber.Map{"status": "error", "message": "not found"})
	}
	fmt.Println(datas)

	var count int64
	if err := db.Model(&model.Data{}).Count(&count).Error; err != nil {
		return res.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	paginationInfo := fiber.Map{
		"total":    count,
		"page":     page,
		"pageSize": pageSize,
	}

	responseData := make([]json.RawMessage, 0)

	for _, data := range datas {
		jsonBytes, err := data.MarshalJSON()
		if err != nil {
			return res.Status(500).JSON(fiber.Map{"status": "error", "message": "Error marshaling data"})
		}
		responseData = append(responseData, jsonBytes)
	}

	return res.Status(200).JSON(fiber.Map{"status": "success", "message": "Found", "data": responseData, "paginationInfo": paginationInfo})

}
