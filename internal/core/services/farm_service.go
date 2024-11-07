package services

import (
	"0xFarms-backend/internal/adapters"
	"0xFarms-backend/internal/core/domain"
	"0xFarms-backend/internal/ports"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// FarmManagementSystemService handles all farm operations
type FarmManagementSystemService struct {
	db ports.MongoDB
}

// NewFarmManagementSystemService initializes a new farm management system
func NewFarmManagementSystemService(db ports.MongoDB) *FarmManagementSystemService {
	system := &FarmManagementSystemService{
		db: db,
	}
	return system
}

// // initializeCropSpecs sets up default crop specifications
// func (fms *FarmManagementSystemService) initializeCropSpecs() {
// 	fms.crops["lettuce"] = domain.CropSpecification{
// 		Name:               "Lettuce",
// 		OptimalPH:          6.5,
// 		OptimalHumidity:    65.0,
// 		GrowthPeriod:       45 * 24 * time.Hour,
// 		OptimalTemp:        23.0,
// 		NutrientNeeds:      0.8,
// 		ExpectedYieldPerM2: 4.5,
// 	}
// 	// Add more crops as needed
// }

// CreateFarm initializes a new vertical farm
func (fms *FarmManagementSystemService) CreateFarm(width, height float64, cropType string) (*domain.VerticalFarm, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("invalid dimensions")
	}

	cropSpec, exists := fms.getCropSpecification(cropType)
	if !exists {
		return nil, errors.New("unsupported crop type")
	}

	farm := &domain.VerticalFarm{
		ID:                   uuid.New().String(),
		Width:                width,
		Height:               height,
		TotalArea:            width * height,
		CropType:             cropType,
		PlantingDate:         time.Now(),
		EstimatedHarvestTime: time.Now().Add(cropSpec.GrowthPeriod),
		Owners:               make([]domain.Owner, 0),
		IoTData:              make([]domain.IoTReading, 0),
		Status:               "active",
		CurrentHealth:        100,
		LastUpdated:          time.Now(),
	}

	id, err := fms.db.CreateFarm(farm)
	if err != nil {
		return nil, err
	}
	farm.ID = id
	return farm, nil
}

// AddOwner adds a new owner to the farm
func (fms *FarmManagementSystemService) AddOwner(farmID, address string, shareSize float64) error {
	farm, err := fms.db.GetFarm(farmID)
	if err != nil {
		return err
	}

	// Calculate total existing shares
	var totalShares float64
	for _, owner := range farm.Owners {
		totalShares += owner.ShareSize
	}

	if totalShares+shareSize > 100 {
		return errors.New("ownership share exceeds 100%")
	}

	owner := domain.Owner{
		ID:        uuid.New().String(),
		Address:   address,
		ShareSize: shareSize,
		JoinedAt:  time.Now(),
	}

	farm.Owners = append(farm.Owners, owner)
	_, err = fms.db.UpdateFarm(farmID, farm)
	return err
}

// AddIoTReading adds a new IoT sensor reading and updates farm status
func (fms *FarmManagementSystemService) AddIoTReading(farmID string, reading domain.IoTReading) error {
	farm, err := fms.db.GetFarm(farmID)
	if err != nil {
		return err
	}

	cropSpec, ok := fms.getCropSpecification(farm.CropType)
	if !ok {
		return errors.New("crop specification not found")
	}

	// Calculate crop health based on optimal conditions
	healthScore := fms.calculateHealthScore(reading, cropSpec)
	reading.CropHealth = healthScore

	// Calculate expected yield based on health and area
	reading.ExpectedYield = fms.calculateExpectedYield(farm, healthScore, cropSpec)

	err = fms.db.AddIoTReading(farmID, &reading)
	if err != nil {
		return err
	}

	farm.IoTData = append(farm.IoTData, reading)
	farm.CurrentHealth = healthScore
	farm.LastUpdated = reading.Timestamp

	_, err = fms.db.UpdateFarm(farmID, farm)
	return err
}

// GetFarmStatus retrieves current farm status and analytics
func (fms *FarmManagementSystemService) GetFarmStatus(farmID string) (*domain.VerticalFarm, error) {
	return fms.db.GetFarm(farmID)
}

// calculateHealthScore determines crop health based on environmental conditions
func (fms *FarmManagementSystemService) calculateHealthScore(reading domain.IoTReading, spec domain.CropSpecification) int {
	phScore := 100 - math.Abs(reading.SoilPH-spec.OptimalPH)*10
	humidityScore := 100 - math.Abs(reading.Humidity-spec.OptimalHumidity)
	tempScore := 100 - math.Abs(reading.Temperature-spec.OptimalTemp)*2
	nutrientScore := (reading.NutrientLevel / spec.NutrientNeeds) * 100

	// Weighted average of all scores
	healthScore := (phScore*0.25 + humidityScore*0.25 + tempScore*0.25 + nutrientScore*0.25)

	return int(math.Max(0, math.Min(100, healthScore)))
}

// calculateExpectedYield estimates crop yield based on current conditions
func (fms *FarmManagementSystemService) calculateExpectedYield(farm *domain.VerticalFarm, health int, spec domain.CropSpecification) float64 {
	baseYield := farm.TotalArea * spec.ExpectedYieldPerM2
	healthFactor := float64(health) / 100.0
	return baseYield * healthFactor
}

// Example usage
func mains() {
	db, _ := adapters.NewMongoAdapter("")
	// Initialize the system
	fms := NewFarmManagementSystemService(db)

	// Create a new vertical farm
	farm, _ := fms.CreateFarm(10.0, 5.0, "lettuce")

	// Add an owner
	fms.AddOwner(farm.ID, "0x123abc...", 50.0)

	// Simulate IoT reading
	reading := domain.IoTReading{
		Timestamp:     time.Now(),
		SoilPH:        6.5,
		Humidity:      65.0,
		NutrientLevel: 0.8,
		Temperature:   23.0,
	}
	fms.AddIoTReading(farm.ID, reading)

	// Get farm status
	farmStatus, _ := fms.GetFarmStatus(farm.ID)
	_ = farmStatus // Use the status as needed
}

// getCropSpecification retrieves the crop specification from the database
func (fms *FarmManagementSystemService) getCropSpecification(cropType string) (domain.CropSpecification, bool) {
	// Fetch crop specification from the database
	cropSpec, err := fms.db.GetCropSpecification(cropType)
	if err != nil {
		return domain.CropSpecification{}, false
	}
	return cropSpec, true
}
