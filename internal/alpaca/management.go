package alpaca

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Management endpoints use fixed zero transaction ids (no client echo).

type apiVersionsResponse struct {
	Value               []int  `json:"Value"`
	ErrorNumber         int    `json:"ErrorNumber"`
	ErrorMessage        string `json:"ErrorMessage"`
	ClientTransactionID uint32 `json:"ClientTransactionID"`
	ServerTransactionID uint32 `json:"ServerTransactionID"`
}

type serverDescriptionValue struct {
	ServerName          string `json:"ServerName"`
	Manufacturer        string `json:"Manufacturer"`
	ManufacturerVersion string `json:"ManufacturerVersion"`
	Location            string `json:"Location"`
}

type descriptionResponse struct {
	Value               serverDescriptionValue `json:"Value"`
	ErrorNumber         int                    `json:"ErrorNumber"`
	ErrorMessage        string                 `json:"ErrorMessage"`
	ClientTransactionID uint32                 `json:"ClientTransactionID"`
	ServerTransactionID uint32                 `json:"ServerTransactionID"`
}

type configuredDevice struct {
	DeviceName   string `json:"DeviceName"`
	DeviceType   string `json:"DeviceType"`
	DeviceNumber int    `json:"DeviceNumber"`
	UniqueID     string `json:"UniqueID"`
}

type configuredDevicesResponse struct {
	Value               []configuredDevice `json:"Value"`
	ErrorNumber         int                `json:"ErrorNumber"`
	ErrorMessage        string             `json:"ErrorMessage"`
	ClientTransactionID uint32             `json:"ClientTransactionID"`
	ServerTransactionID uint32             `json:"ServerTransactionID"`
}

// getAPIVersions returns the list of Alpaca API versions supported.
func getAPIVersions(c *gin.Context) {
	c.JSON(http.StatusOK, apiVersionsResponse{Value: []int{1}})
}

// getManagementDescription returns server metadata.
func getManagementDescription(c *gin.Context) {
	c.JSON(http.StatusOK, descriptionResponse{
		Value: serverDescriptionValue{
			ServerName:          serverName,
			Manufacturer:        serverName,
			ManufacturerVersion: driverVersion,
			Location:            "Embedded",
		},
	})
}

// getConfiguredDevices returns the list of devices exposed (one telescope).
func getConfiguredDevices(c *gin.Context) {
	c.JSON(http.StatusOK, configuredDevicesResponse{
		Value: []configuredDevice{{
			DeviceName:   deviceName,
			DeviceType:   deviceType,
			DeviceNumber: deviceNumber,
			UniqueID:     uniqueID,
		}},
	})
}
