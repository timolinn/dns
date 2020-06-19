package handlers

import (
	"errors"
	"math"
)

// SectorID is a unique identifier for different
// sectors in the galaxy
var SectorID float64 = 1

// System represents a type that may request DNS service
type System string

const (
	Drone      System = "drone"
	Ship              = "vessel"
	UltraDrone        = "ultradrones"
)

var (
	ErrUnknownSystemType = errors.New("invalid system type: requires 'drone', 'vessel' or 'ultradrone'")
)

// CoordsVelocity expected request data
type CoordsVelocity struct {
	X   float64 `json:"x,string" validate:"required"`
	Y   float64 `json:"y,string" validate:"required"`
	Z   float64 `json:"z,string" validate:"required"`
	Vel float64 `json:"vel,string" validate:"required"`
}

// Navigator describes a contract for providing
// navigation service to multiple kinds of systems
type Navigator interface {
	Solve(CoordsVelocity, System) (float64, error)
	Response(float64, System) map[string]float64
}

// SectorNavigator provides navigation functionality
// it implements Navigator interface
type SectorNavigator struct {
	SectorID float64
}

// NewSectorNavigator constructor a new Navigator type for a sector
func NewSectorNavigator() Navigator {
	return &SectorNavigator{
		SectorID: SectorID,
	}
}

// Solve computes the navigation puzzle
func (sn *SectorNavigator) Solve(cv CoordsVelocity, system System) (float64, error) {
	switch system {
	case Drone, Ship:
		result := (cv.X * sn.SectorID) + (cv.Y * sn.SectorID) + (cv.Z * sn.SectorID) + cv.Vel*sn.SectorID
		// round result to two decimal places
		return (math.Round(result*100) / 100), nil
	default:
		return 0, ErrUnknownSystemType
	}
}

// Response constucts a response map based on systemType
func (sn *SectorNavigator) Response(data float64, systemType System) map[string]float64 {
	resp := make(map[string]float64)
	switch systemType {
	case Ship:
		resp["location"] = data
	default:
		resp["loc"] = data
	}
	return resp
}
