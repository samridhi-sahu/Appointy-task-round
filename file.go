package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Resource struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BusinessHour struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	Quantity   int64  `json:"quantity"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type BlockHour struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type Appointment struct {
	Id         string `json:"id"`
	ResourceId string `json:"resource_id"`
	Quantity   int64  `json:"quantity"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

type Duration struct {
	Seconds int64 `json:"seconds"`
}

type ListBusinessHoursRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type ListBlockHoursRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type ListAppointmentRequest struct {
	ResourceId string `json:"resourceId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func TimeToString(tm time.Time) string {
	return tm.Format(time.RFC3339)
}

func StringToTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func main() {

	router := mux.NewRouter()

	// create /availability endpoint
	router.HandleFunc("/availability", availabilityHandler).Methods("GET")

	// Run server
	http.ListenAndServe(":8000", router)

	inputParam := map[string]interface{}{
		"resourceId": "res_2",
		"date":       "2023-08-05",
		"duration":   "30",
		"quantity":   "1",
	}

	resourceId := inputParam["resourceId"].(string)
	startTime := inputParam["date"].(string) + "T00:00:00Z"
	endTime := inputParam["date"].(string) + "T23:59:00Z"

	payload := map[string]interface{}{
		"resourceId": resourceId,
		"startTime":  startTime,
		"endTime":    endTime,
	}

	businesshours := apiCall("/business-hours", payload)
	blockhours := apiCall("/block-hours", payload)
	appointment := apiCall("/appointments", payload)

	var businesshoursMap []BusinessHour
	json.Unmarshal([]byte(businesshours), &businesshoursMap)

	var blockhoursMap []BlockHour
	json.Unmarshal([]byte(blockhours), &blockhoursMap)

	var appointmentMap []Appointment
	json.Unmarshal([]byte(appointment), &appointmentMap)

	for i := 0; i < len(businesshoursMap); i++ {
		startTime, _ := StringToTime(businesshoursMap[i].StartTime)
		endTime, _ := StringToTime(businesshoursMap[i].EndTime)

		duration, _ := time.ParseDuration(inputParam["duration"].(string) + "m")

		fmt.Println("Business Hours: ", i+1)
		fmt.Println("Start Time: ", startTime, "End Time: ", endTime)
		fmt.Println("Duration: ", duration)

		for j := startTime; j.Before(endTime); j = j.Add(duration) {

			fmt.Println("Slot: ", j, "to", j.Add(duration))

			for k := 0; k < len(blockhoursMap); k++ {
				blockStartTime, _ := StringToTime(blockhoursMap[k].StartTime)
				blockEndTime, _ := StringToTime(blockhoursMap[k].EndTime)

				if j.After(blockStartTime) && j.Before(blockEndTime) {
					fmt.Println("blocked")
					break
				}
			}

			fmt.Println("available")
		}

		fmt.Println("")
	}

}

func apiCall(endpoint string, payload map[string]interface{}) string {

	url := "http://api.internship.appointy.com:8000/v1"
	method := "GET"

	newurl := url + endpoint

	if payload != nil {
		newurl = newurl + "?"
		for key, value := range payload {
			newurl = newurl + key + "=" + value.(string) + "&"
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, newurl, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIzLTA4LTEwVDAwOjAwOjAwWiIsInVzZXJfaWQiOjMwMDF9.8pZMhoqZdBLqOKT0V7perD4vkoA347idSHVLaCcdefs")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
