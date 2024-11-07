package handlers

import "0xFarms-backend/internal/core/services"

type FarmHandler struct {
	farmService *services.FarmManagementSystemService
}

// NewCommitHandler creates a new instance of CommitHandler with the given services
func NewFarmHandler(farmService *services.FarmManagementSystemService) *FarmHandler {
	return &FarmHandler{
		farmService: farmService,
	}
}
