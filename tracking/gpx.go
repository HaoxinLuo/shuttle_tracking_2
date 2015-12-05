package tracking

import (
//	"encoding/json"
//	"io/ioutil"
	"net/http"
//	"regexp"
	// "strings"
	"time"
  "strconv"
//  "os"
	// "log"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// RoutesHandler finds all of the routes in the database
func (App *App) GPXHandler(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now()
	prevTime := curTime
	prevTime.AddDate(0,-1,0)
	log.Infof("foo")
  // Generate gpx file's context and save it to string gpxFile
	gpxFile := "<gpx creator=\"RPI Shuttle Tracking v2\" version=\"1.0\">\n"


	var vehicles []Vehicle
	err := App.Vehicles.Find(bson.M{}).All(&vehicles)

	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	counter := 0
	// Find recent updates for each vehicle
	for _, vehicle := range vehicles {
		var vehicleHistory []VehicleUpdate
		// App.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID, "Created": {"$lt": time.Time(), "$gt": time.Time() - time.Duration(24)*time.Hour}}).Sort("-created").All(&vehicleHistory)
		App.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID, "created": bson.M{"$lt": time.Now(), "$gt": time.Now().Add(time.Hour*-24)} }).Sort("-created").All(&vehicleHistory)
		if len(vehicleHistory) > 0 {
			gpxFile = gpxFile + "<trk>\n" + "<name>" + vehicle.VehicleID + "</name>\n" + "<trkseg>\n"


			counter = counter + 1
			for _, location := range vehicleHistory {
				loc := location.Created
				gpxFile = gpxFile + "<trkpt lat=\"" + location.Lat + "\" lon=\"" + location.Lng + "\">\n"
				//gpxFile = gpxFile + "<time>" + location.Created.CFormat("%Y-%m-%dT%H:%M:%SZ")	 + "</time>\n</trkpt>\n"
				gpxFile = gpxFile + "<time>" + convTime(loc)+ "</time>\n</trkpt>\n"
			}
			gpxFile = gpxFile + "</trkseg>\n</trk>\n"
		}
	}

	// gpxFile := vehiclesHistory[0].Lat



  // ...
	gpxFile = gpxFile + "</gpx>"
	// Send gpx file content to client as text
	WriteText(w, gpxFile)
}

func convTime(loc time.Time) (s string){
	s = strconv.Itoa(loc.Year())+"-"+
		strconv.Itoa(int(loc.Month()))+"-"+
		strconv.Itoa(loc.Day())+"T"+
		strconv.Itoa(loc.Hour())+":"+
		strconv.Itoa(loc.Minute())+":"+
		strconv.Itoa(loc.Second())+"Z"
	return s
}
