package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

//Google API key
var key = "AIzaSyAFbIU6eUJI63qktYdl3im41NSAckq5J14"

//Apiresponse is a structure made for the api response, so we may parse it
type Apiresponse struct {
	Results       []Result `json:"results"`
	NextPageToken string   `json:"next_page_token"`
}

//Location struct to save locations  coordinats as strings and json easily
type Location struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

//Result is a support struct for parsing results-array inside of the api response
type Result struct {
	//ID   string `json:"id"`
	Name string `json:"name"`
	//PlaceID  string `json:"place_id"`
	Vicinity string `json:"vicinity"`
}

//Shop is a struct made to save a type of shop and its results
type Shop struct {
	Results []Result `json:"results"`
}

//GeoCodeResponse is a struct to parse json from google geocode api
type GeoCodeResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64
				Lng float64
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

var bikestores []Shop

//CallGoogleNearbyPlaces calls googles api with two input parameters
func CallGoogleNearbyPlaces(lat string, lng string, location string) []byte {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	var resMars []byte
	var ResultArray []Result
	var nextToken string

	response, err := client.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + lat + "," + lng + "&radius=2000&type=bicycle_store&key=" + key)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(response.Body)
		structure := new(Apiresponse)
		json.Unmarshal(data, structure)
		nextToken = structure.NextPageToken
		theArray := structure.Results

		for _, loc := range theArray {
			ResultArray = append(ResultArray, loc)
		}
		response.Body.Close()
		response.Close = true
		for len(nextToken) != 0 {
			time.Sleep(10 * time.Second)
			response, err := client.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?&key=AIzaSyAFbIU6eUJI63qktYdl3im41NSAckq5J14&pagetoken=" + string(nextToken))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				structure := new(Apiresponse)
				json.Unmarshal(data, structure)
				nextToken = structure.NextPageToken
				theArray := structure.Results
				for _, loc := range theArray {
					ResultArray = append(ResultArray, loc)
				}
				response.Body.Close()
				response.Close = true
			}
		}
		shopItem := new(Shop)
		shopItem.Results = ResultArray
		resShop, error := json.Marshal(shopItem)
		if error != nil {
			println("Did not work")
		}
		resMars = resShop
	}
	return resMars
}

//CallGoogleGeoAPI is for calling GeoCoding
func CallGoogleGeoAPI(location string) (string, string) {
	var returnlat string
	var returnlng string
	response, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + location + "&key=" + key)
	if err != nil {
		fmt.Printf("Error %s\n", err)
		fmt.Printf("Failed to get bikestores for location %s perhaps you formated it wrong. Replace space with plus sign", location)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		geostruct := new(GeoCodeResponse)
		json.Unmarshal(data, geostruct)
		lat := geostruct.Results[0].Geometry.Location.Lat
		lng := geostruct.Results[0].Geometry.Location.Lng
		returnlat = fmt.Sprintf("%f", lat)
		returnlng = fmt.Sprintf("%f", lng)
		response.Body.Close()
	}
	return returnlat, returnlng
}

//CallGeoAPIorReadFile just checks if there is a file first
func CallGeoAPIorReadFile(place string) (string, string) {
	var a string
	var b string

	locfile, errfile := ioutil.ReadFile("Location" + place + ".json")
	if errfile != nil {
		a, b = CallGoogleGeoAPI(place)
		locdone := new(Location)
		locdone.Lat = a
		locdone.Lng = b
		writingfile, errfile := os.Create("Location" + place + ".json")
		if errfile != nil {
			println("File could not be created")
		} else {
			done, error := json.Marshal(locdone)
			if error != nil {
				println("Did not work")
			}
			writingfile.Write(done)
			writingfile.Close()
		}
	} else {
		//Read location file
		locstruc := new(Location)
		json.Unmarshal(locfile, locstruc)
		a = locstruc.Lat
		b = locstruc.Lng
	}
	return a, b
}

//GetBicycleStoresEndpoint returns all the bicycle stores close to Segels Torg
func GetBicycleStoresEndpoint(w http.ResponseWriter, req *http.Request) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	var ResultArray []Result
	var nextToken string

	response, err := client.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=59.332356,18.064545&radius=2000&type=bicycle_store&key=" + key)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		writingfile, errfile := os.Create("Sergels+Torg.json")
		if errfile != nil {
			println("File could not be created")
		}

		data, _ := ioutil.ReadAll(response.Body)
		structure := new(Apiresponse)
		json.Unmarshal(data, structure)
		nextToken = structure.NextPageToken
		theArray := structure.Results
		//fmt.Println(theArray)
		for _, loc := range theArray {
			res, _ := json.Marshal(loc)
			ResultArray = append(ResultArray, loc)
			fmt.Print(string(res))
		}
		response.Body.Close()
		response.Close = true
		for len(nextToken) != 0 {
			time.Sleep(60 * time.Second)
			response, err := client.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?&key=AIzaSyAFbIU6eUJI63qktYdl3im41NSAckq5J14&pagetoken=" + string(nextToken))
			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				structure := new(Apiresponse)
				json.Unmarshal(data, structure)
				nextToken = structure.NextPageToken
				theArray := structure.Results
				for _, loc := range theArray {
					res, _ := json.Marshal(loc)
					ResultArray = append(ResultArray, loc)
					fmt.Print(string(res))
				}
				fmt.Println(string(nextToken))
				response.Body.Close()
				response.Close = true
				println("Next token len : " + nextToken)
			}
		}
		shopItem := new(Shop)
		shopItem.Results = ResultArray
		resMars, error := json.Marshal(shopItem)
		if error != nil {
			println("Did not work")
		}

		writingfile.Write(resMars)
		writingfile.Close()
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.Encode(ResultArray)
	}
}

//GetSpecificBikestoreLocation  fetches a specific store
func GetSpecificBikestoreLocation(w http.ResponseWriter, req *http.Request) {

	//This is a big weakness, the variable has to be input "Sergels+Torg"
	//Making routes with url.query would solve this but no time
	params := mux.Vars(req)
	if len(params) == 0 {
		error := fmt.Sprintf("Failed to get bikestores for location perhaps you formated it wrong. Replace space with plus sign")
		json.NewEncoder(w).Encode(error)
		return
	}
	var file []byte
	location := params["location"]
	//This function does both things since location RARELY ever changes coordinates (never?)
	a, b := CallGeoAPIorReadFile(location)
	bikefile, errfile := ioutil.ReadFile(location + ".json")
	if errfile != nil {
		file = CallGoogleNearbyPlaces(a, b, location)
		writingfile, errfile := os.Create(location + ".json")
		if errfile != nil {
			println("File could not be created")
		} else {
			writingfile.Write(file)
			writingfile.Close()
		}
	} else {
		filestruc := new(Shop)
		json.Unmarshal(bikefile, filestruc)
		file = bikefile
	}
	fileShop := new(Shop)
	json.Unmarshal(file, fileShop)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(fileShop.Results)
}

//GetSpecificBikestore gets locations and nearby places but also searches for a store
func GetSpecificBikestore(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	if len(params) == 0 {
		error := fmt.Sprintf("Failed to get bikestores for location perhaps you formated it wrong. Replace space with plus sign")
		json.NewEncoder(w).Encode(error)
		return
	}
	notFound := Result{"Name not found", ""}
	empty, error := json.Marshal(notFound)
	if error != nil {
		println("parsing went wrong notFound")
	}
	location := params["location"]
	var shopFound = empty
	var file []byte

	//This function does both things since location RARELY ever changes coordinates (never?)
	a, b := CallGeoAPIorReadFile(location)
	bikefile, errfile := ioutil.ReadFile(location + ".json")
	if errfile != nil {
		file = CallGoogleNearbyPlaces(a, b, location)
		writingfile, errfile := os.Create(location + ".json")
		if errfile != nil {
			println("File could not be created")
		} else {
			writingfile.Write(file)
			writingfile.Close()
		}
	} else {
		filestruc := new(Shop)
		json.Unmarshal(bikefile, filestruc)
		file = bikefile
	}

	fileShop := new(Shop)
	json.Unmarshal(file, fileShop)

	for _, loc := range fileShop.Results {
		if loc.Name == params["store"] {
			mrshed, err := json.Marshal(loc)
			if err != nil {
				println("parsing error")
			}
			res := new(Result)
			json.Unmarshal(mrshed, res)
			d, en := json.Marshal(res)
			if en != nil {
				println("parsing error")
			}
			shopFound = d
		}
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(string(shopFound))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/test", GetBicycleStoresEndpoint).Methods("GET")
	router.HandleFunc("/bicyclestores/{location}/", GetSpecificBikestoreLocation).Methods("GET")
	router.HandleFunc("/bicyclestores/{location}/{store}", GetSpecificBikestore).Methods("GET")
	//router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	//router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", router))
}
