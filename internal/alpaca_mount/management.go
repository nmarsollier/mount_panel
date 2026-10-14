package alpaca_mount

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
)

// Management endpoints use fixed zero transaction ids (no client echo).

type apiVersionsResponse struct {
	Value               []int  `json:"Value"`
	ErrorNumber         int    `json:"ErrorNumber"`
	ErrorMessage        string `json:"ErrorMessage"`
	ClientTransactionID uint32 `json:"ClientTransactionID"`
	ServerTransactionID uint32 `json:"ServerTransactionID"`
}

type descriptionValue struct {
	ServerName          string `json:"ServerName"`
	Manufacturer        string `json:"Manufacturer"`
	ManufacturerVersion string `json:"ManufacturerVersion"`
	Location            string `json:"Location"`
}

type descriptionResponse struct {
	Value               descriptionValue `json:"Value"`
	ErrorNumber         int              `json:"ErrorNumber"`
	ErrorMessage        string           `json:"ErrorMessage"`
	ClientTransactionID uint32           `json:"ClientTransactionID"`
	ServerTransactionID uint32           `json:"ServerTransactionID"`
}

type configuredDevicesResponse struct {
	Value               []server.DeviceInfo `json:"Value"`
	ErrorNumber         int                 `json:"ErrorNumber"`
	ErrorMessage        string              `json:"ErrorMessage"`
	ClientTransactionID uint32              `json:"ClientTransactionID"`
	ServerTransactionID uint32              `json:"ServerTransactionID"`
}

// getAPIVersions returns the list of Alpaca API versions supported.
func getAPIVersions(c *gin.Context) {
	c.JSON(http.StatusOK, apiVersionsResponse{Value: []int{1}})
}

// getManagementDescription returns server metadata.
func getManagementDescription(c *gin.Context) {
	c.JSON(http.StatusOK, descriptionResponse{
		Value: descriptionValue{
			ServerName:          server.ServerName,
			Manufacturer:        server.ServerName,
			ManufacturerVersion: server.DriverVersion,
			Location:            "Embedded",
		},
	})
}

// getConfiguredDevices returns the list of devices exposed.
func getConfiguredDevices(c *gin.Context) {
	c.JSON(http.StatusOK, configuredDevicesResponse{
		Value: configuredDevices(),
	})
}
