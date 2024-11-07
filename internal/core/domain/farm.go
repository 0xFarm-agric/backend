package domain

import "time"

// Owner represents a stakeholder in the farm
type Owner struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	ShareSize float64   `json:"shareSize"` // Percentage of ownership
	JoinedAt  time.Time `json:"joinedAt"`
}

// IoTReading represents a single data point from IoT sensors
type IoTReading struct {
	Timestamp     time.Time `json:"timestamp"`
	SoilPH        float64   `json:"soilPH"`
	Humidity      float64   `json:"humidity"`
	NutrientLevel float64   `json:"nutrientLevel"`
	CropHealth    int       `json:"cropHealth"`    // Scale of 1-100
	ExpectedYield float64   `json:"expectedYield"` // in kgs
	Temperature   float64   `json:"temperature"`   // in Celsius
}

// VerticalFarm represents a single vertical farming unit
type VerticalFarm struct {
	ID                   string       `json:"id"`
	Width                float64      `json:"width"`     // in meters
	Height               float64      `json:"height"`    // in meters
	TotalArea            float64      `json:"totalArea"` // in square meters
	CropType             string       `json:"cropType"`
	PlantingDate         time.Time    `json:"plantingDate"`
	EstimatedHarvestTime time.Time    `json:"estimatedHarvestTime"`
	Owners               []Owner      `json:"owners"`
	IoTData              []IoTReading `json:"iotData"`
	Status               string       `json:"status"` // active, harvested, maintenance
	CurrentHealth        int          `json:"currentHealth"`
	LastUpdated          time.Time    `json:"lastUpdated"`
}

// CropSpecification contains default parameters for different crops
type CropSpecification struct {
	Name               string        `json:"name"`
	OptimalPH          float64       `json:"optimalPH"`
	OptimalHumidity    float64       `json:"optimalHumidity"`
	GrowthPeriod       time.Duration `json:"growthPeriod"`
	OptimalTemp        float64       `json:"optimalTemp"`
	NutrientNeeds      float64       `json:"nutrientNeeds"`
	ExpectedYieldPerM2 float64       `json:"expectedYieldPerM2"`
}
