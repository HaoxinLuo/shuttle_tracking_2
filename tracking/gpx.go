package tracking

import (
//	"encoding/json"
//	"io/ioutil"
	"net/http"
//	"regexp"
//	"strings"
//	"time"
//  "strconv"
//  "os"

//	log "github.com/Sirupsen/logrus"
//	"gopkg.in/mgo.v2/bson"
)

// RoutesHandler finds all of the routes in the database
func (App *App) GPXHandler(w http.ResponseWriter, r *http.Request) {

  // Generate gpx file's context and save it to string gpxFile
	gpxFile := ""

  // ...

	// Send gpx file content to client as text
	WriteText(w, gpxFile)
}
