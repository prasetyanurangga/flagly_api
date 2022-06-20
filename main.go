package main

import (
     "math"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    swaggerFiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
   _ "github/prasetyanurangga/flagly_api/docs"
)

type Point struct {
     lat float64
     lng float64
}

type LatLong struct {
     LatFrom float64 `json:"lat_from"`
     LngFrom float64 `json:"lng_from"`
     LatTo float64 `json:"lat_to"`
     LngTo float64 `json:"lng_to"`
}

type Response struct {
     Data ResponseItem `json:"data"`
}

type ResponseItem struct {
     Distance int `json:"distance"`
     Bearing int `json:"bearing"`
     Direction string `json:"direction"`
}


const (
     // According to Wikipedia, the Earth's radius is about 6,371km
     EARTH_RADIUS = 6371
)

func (p *Point) GreatCircleDistance(p2 *Point) float64 {
     dLat := (p2.lat - p.lat) * (math.Pi / 180.0)
     dLon := (p2.lng - p.lng) * (math.Pi / 180.0)

     lat1 := p.lat * (math.Pi / 180.0)
     lat2 := p2.lat * (math.Pi / 180.0)

     a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
     a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

     a := a1 + a2

     c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

     // distance in KM
     return EARTH_RADIUS * c
}

func (p *Point) BearingTo(p2 *Point) float64 {

     dLon := (p2.lng - p.lng) * math.Pi / 180.0

     lat1 := p.lat * math.Pi / 180.0
     lat2 := p2.lat * math.Pi / 180.0

     y := math.Sin(dLon) * math.Cos(lat2)
     x := math.Cos(lat1)*math.Sin(lat2) -
          math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
     brngRaw := math.Atan2(y, x) 
     brng := math.Mod(((brngRaw * 180.0 / math.Pi) + 360), 360)
     return brng
}

func NewPoint(lat float64, lng float64) *Point {
     return &Point{lat: lat, lng: lng}
}

func getDirection(bearing float64) string {
     if (bearing >= 337.5 && bearing < 360) || (bearing >= 0 && bearing < 22.5) {
          return "N"
     } 

     if bearing >= 22.5 && bearing < 67.5 {
          return "NE"
     }

     if bearing >= 67.5 && bearing < 112.5 {
          return "E"
     }

     if bearing >= 112.5 && bearing < 157.5 {
          return "SE"
     }

     if bearing >= 157.5 && bearing < 202.5 {
          return "S"
     }

     if bearing >= 202.5 && bearing < 247.5 {
          return "SW"
     }

     if bearing >= 247.5 && bearing < 292.5 {
          return "W"
     }

     if bearing >= 292.5 && bearing < 337.5 {
          return "NW"
     }

     return "Error"
}

// @title Flagly API
// @version 1.0
// @description Calculate Distance, Bearing and Direction Compass.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url prasetyanurangga.github.io
// @contact.email angganurprasetya4@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host pintas.my.id
// @BasePath /
// @schemes https
func main() {

     router := gin.Default()
     router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"POST", "OPTIONS", "GET"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
     }))
     router.POST("/calculate_distance_bearing", Calculate)
     router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
     router.GET("/", func(c *gin.Context) {
          c.String(http.StatusOK, "Hello Word")
     })
     router.Run()
}

// @Summary Calculate Distance, Bearing and Direction Compass.
// @Description Calculate Distance, Bearing and Direction Compass by Coordinate
// @Tags "Flagly Api"
// @Accept json
// @Produce json
// @Param data body LatLong true "Lat Long Data"
// @Success 200 {object} Response
// @Router /calculate_distance_bearing [post]
func Calculate(c *gin.Context) {
   body := LatLong{}
   if err:=c.BindJSON(&body);err!=nil{
     c.JSON(200,Response{
     Data : ResponseItem{
          Distance : 0,
          Bearing : 0,
          Direction : "",
        },
   })
         
   }
   p := NewPoint(body.LatFrom, body.LngFrom)
   p2 := NewPoint(body.LatTo, body.LngTo)

     dist := p.GreatCircleDistance(p2)
     bearing := p.BearingTo(p2)
     direction := getDirection(bearing)
   c.Header("Content-Type", "application/json")
   c.JSON(200, Response{
     Data : ResponseItem{
          Distance : int(dist),
          Bearing : int(bearing),
          Direction : direction,
        },
   })
}