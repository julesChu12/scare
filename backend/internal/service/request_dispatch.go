package service

import (
	"errors"
	"strings"

	"community-elderly-care-platform/internal/repository"
)

const (
	DispatchBasisServiceGeofence     = "service_geofence"
	DispatchBasisServiceNearest      = "service_nearest"
	DispatchBasisAddressManualReview = "service_address_manual_review"
)

var (
	ErrServiceLocationRequired = errors.New("service location required")
	ErrAddressRequired         = errors.New("address required")
)

type DispatchInput struct {
	Address          string
	SubmitLatitude   *float64
	SubmitLongitude  *float64
	ServiceLatitude  *float64
	ServiceLongitude *float64
	SourceStationID  *int64
}

type DispatchDecision struct {
	ResolvedAddress   string
	SubmitLatitude    float64
	SubmitLongitude   float64
	ServiceLatitude   float64
	ServiceLongitude  float64
	SourceStationID   int64
	AssignedStationID int64
	DispatchBasis     string
	NeedsManualVerify bool
}

func resolveDispatch(input DispatchInput, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService, geocodeSvc *GeocodeService) (*DispatchDecision, error) {
	address := strings.TrimSpace(input.Address)
	decision := &DispatchDecision{
		ResolvedAddress: address,
	}

	if err := assignCoordinatePair(input.SubmitLatitude, input.SubmitLongitude, &decision.SubmitLatitude, &decision.SubmitLongitude); err != nil {
		return nil, err
	}

	validSourceStationID, err := resolveSourceStationID(input.SourceStationID, stationRepo)
	if err != nil {
		return nil, err
	}
	decision.SourceStationID = validSourceStationID

	if address == "" {
		return nil, ErrAddressRequired
	}

	if geocodeSvc != nil {
		geo, geoErr := geocodeSvc.Geocode(address)
		if geoErr == nil && validCoordinate(geo.Latitude, geo.Longitude) {
			decision.ServiceLatitude = geo.Latitude
			decision.ServiceLongitude = geo.Longitude
			if geo.FormattedAddress != "" {
				decision.ResolvedAddress = geo.FormattedAddress
			}
		} else if geoErr != nil && !errors.Is(geoErr, ErrGeocodeNotFound) {
			return nil, geoErr
		}
	}

	hasServiceLocation := decision.ServiceLatitude != 0 && decision.ServiceLongitude != 0
	if hasServiceLocation {
		stationID, matched, matchErr := resolveAssignedStation(decision.ServiceLatitude, decision.ServiceLongitude, stationRepo, geofenceSvc)
		if matchErr != nil {
			return nil, matchErr
		}

		decision.AssignedStationID = stationID
		if matched {
			decision.DispatchBasis = DispatchBasisServiceGeofence
		} else {
			decision.DispatchBasis = DispatchBasisServiceNearest
		}
		return decision, nil
	}

	decision.DispatchBasis = DispatchBasisAddressManualReview
	decision.NeedsManualVerify = true
	return decision, nil
}

func assignCoordinatePair(lat, lng *float64, outLat, outLng *float64) error {
	if lat == nil && lng == nil {
		return nil
	}
	if lat == nil || lng == nil {
		return ErrInvalidRequest
	}
	if !validCoordinate(*lat, *lng) {
		return ErrInvalidRequest
	}

	*outLat = *lat
	*outLng = *lng
	return nil
}

func resolveSourceStationID(sourceStationID *int64, stationRepo *repository.StationRepository) (int64, error) {
	if sourceStationID == nil || *sourceStationID <= 0 || stationRepo == nil {
		return 0, nil
	}

	station, err := stationRepo.GetByID(*sourceStationID)
	if err != nil {
		return 0, nil
	}
	if station.Status != "active" {
		return 0, nil
	}
	return station.ID, nil
}

func resolveAssignedStation(lat, lng float64, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService) (int64, bool, error) {
	if geofenceSvc != nil {
		if stationID, matched := geofenceSvc.Match(lat, lng); matched {
			return stationID, true, nil
		}
	}

	if stationRepo == nil {
		return 0, false, ErrNoStation
	}

	stations, err := stationRepo.ListActive()
	if err != nil {
		return 0, false, err
	}

	nearest, err := nearestStationByHaversine(stations, lat, lng)
	if err != nil {
		return 0, false, err
	}
	return nearest.ID, false, nil
}
