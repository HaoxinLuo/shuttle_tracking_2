package tracking

import (
	"net/http"
	"time"
  "strconv"

	"gopkg.in/mgo.v2/bson"
)

// GPXHandler finds all of the vehicles and their locations in the last 24 hours from the database
// Then "generate" gpx file by saving file's content to a string initially
func (App *App) GPXHandler(w http.ResponseWriter, r *http.Request) {

  // Start file with the "gpx" element
	gpxFile := "<gpx creator=\"RPI Shuttle Tracking v2\" version=\"1.0\">\n"

	// Store vehicles from database
	var vehicles []Vehicle
	// Query all Vehicles
	err := App.Vehicles.Find(bson.M{}).All(&vehicles)

	// Handle errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Find location updates in the last 24 hours for each vehicle
	for _, vehicle := range vehicles {
		// Store locations of vehicle
		var vehicleHistory []VehicleUpdate
		// Query all of the vehicle's locations in the last 24 hours
		err := App.Updates.Find(bson.M{"vehicleID": vehicle.VehicleID, "created": bson.M{"$lt": time.Now(), "$gt": time.Now().Add(time.Hour*-24)} }).Sort("-created").All(&vehicleHistory)

		// Add vehicle and its locations to gpx file only if updates were found for that vehicle
		if err == nil && len(vehicleHistory) > 0 {
			// Add "trk" and "trkseg" element to the file for vehicle
			gpxFile = gpxFile + "<trk>\n" + "<name>" + vehicle.VehicleID + "</name>\n" + "<trkseg>\n"
			// Add "trkpt" element to the file for each location of the vehicle
			for i := len(vehicleHistory)-1; i>=0; i-- {
				loc := vehicleHistory[i].Created
				gpxFile = gpxFile + "<trkpt lat=\"" + vehicleHistory[i].Lat + "\" lon=\"" + vehicleHistory[i].Lng + "\">\n"
				// Add "time" element inside the "trkpt" element
				gpxFile = gpxFile + "<time>" + convTime(loc)+ "</time>\n</trkpt>\n"
			}
			// Close "trkseg" and "trk" element
			gpxFile = gpxFile + "</trkseg>\n</trk>\n"
		}
	}

	// Close "gpx" element
	gpxFile = gpxFile + "</gpx>"

	// Send gpx file content to client as text
	WriteText(w, gpxFile)
}

// Convert "time.Time" object to a gpx formated string
func convTime(loc time.Time) (s string){
	s = strconv.Itoa(loc.Year())+"-"+
		strconv.Itoa(int(loc.Month()))+"-"+
		strconv.Itoa(loc.Day())+"T"+
		strconv.Itoa(loc.Hour())+":"+
		strconv.Itoa(loc.Minute())+":"+
		strconv.Itoa(loc.Second())+"Z"
	return s
}
