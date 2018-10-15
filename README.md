## Bikestore Api Service

An api with endpoints for getting bikestores.
Calls the Google PLACES Api and Geocode API and then returns
bikestores within 2000m of desired adress.
The main dependency that needs to be installed and isn't part of normal go-package is gorilla.mux, to install it:
>go get github.com/gorilla/mux

The api runs locally when you start it with 
>go run main.go


### End points
By default all api calls are made towwards
>http://localhost:8081/

# Test
There is a test end point that I have used to just get bikestores for SergelsTorg. It is reached at
>http://localhost:8081/test

# Get all bikestores for a location
To get all nearby bikestores of a 2000m radius of a location
http://localhost:8081/{location}
Where location is an adress or a place. Meant to be for adress but SergelsTorg also works.
**IMPORTANT** to not that all locations and adresses has to use + instead of space. Also all url-encoded type of
input will not work. If you want to search for Sergels Torg. You have to search for Sergels+Torg. This is
a bug that I discovered in the final stages of the api. It's an underlying issue with Gorilla Mux and how it handles
the first variable in the url. It is a big limitation and a way of handling this would have been to use query instead.
(Perhaps future improvements).

# Get a specific bikestore for a location
To get a specific bikestore within a 2000m radius of a location
>http://localhost:8081/{location}/{bikestore}
Where the bikestore is the name of the bikestore and location is adress or place. The bikestores name can have url-encoded name.
In fact to search for it, one has to use a normal space. To see what names you can search for, you can call
the bikestores for a location first.


