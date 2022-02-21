// main.go
package main

import (
	"database/sql"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Article - Our struct for all articles
// type Article struct {
// 	Id      string `json:"Id"`
// 	Title   string `json:"Title"`
// 	Desc    string `json:"desc"`
// 	Content string `json:"content"`
// 	Token   string `json:"token"`
// }

type Article struct {
	Id    string `json:"Id"`
	User  string `json:"User"`
	Token string `json:"token"`
}
type (
	StringInterfaceMap0 map[string]interface{}
	Event0              struct {
		InvoiceNo       string `json:"InvoiceNo"`
		TrackingID      string `json:"TrackingID"`
		Truck_type      string `json:"Truck_type"`
		Truck_plate     string `json:"Truck_plate"`
		Driver_name     string `json:"Driver_name"`
		Driver_Phone    string `json:"DriverTel"`
		Total_chest     string `json:"Total_chest"`
		Total_pack      string `json:"Total_pack"`
		To_created_date string `json:"To_created_date"`
	}
)

var (
	// scheduleTimes          = ""
	APIName                  = ""
	selectEventByIdQueryTOAT = `SELECT * FROM THPDTOATDB.VW_Tracking_TOAT WHERE InvoiceNo = ?`
)
var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to THPDApi!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func getDNSString(dbName, dbUser, dbPassword, conn string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=60s&readTimeout=60s",
		dbUser,
		dbPassword,
		conn,
		dbName)
}
func selectEventById(db *sql.DB, id string, event *Event0) error {
	row := db.QueryRow(selectEventByIdQueryTOAT, id)
	err := row.Scan(&event.InvoiceNo, &event.TrackingID, &event.Truck_type, &event.Truck_plate, &event.Driver_name, &event.Driver_Phone, &event.Total_chest, &event.Total_pack, &event.To_created_date)
	if err != nil {
		return err
	}
	return nil
}
func returnSingleArticle2(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	firstEvent := Event0{}
	err = selectEventById(db, article.Id, &firstEvent)
	if err != nil {
		panic(err)
	}

	resp := make(map[string]string)
	resp["InvoiceNo"] = firstEvent.InvoiceNo
	resp["TrackingID"] = firstEvent.TrackingID
	resp["Truck_type"] = firstEvent.Truck_type
	resp["Truck_plate"] = firstEvent.Truck_plate
	resp["Driver_name"] = firstEvent.Driver_name
	resp["Driver_Phone"] = firstEvent.Driver_Phone
	resp["Total_chest"] = firstEvent.Total_chest
	resp["Total_pack"] = firstEvent.Total_pack
	resp["Start_date"] = firstEvent.To_created_date

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	//fmt.Println(article.Id)
	//fmt.Fprintf(w, jsonResp)
	// for _, article := range Articles {
	// 	if article.Id == key {
	// 		json.NewEncoder(w).Encode(article)
	// 	}
	// }
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, article := range Articles {
		if article.Id == id {
			Articles = append(Articles[:index], Articles[index+1:]...)
		}
	}

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/TOATGetByInvoiceID", returnSingleArticle2).Methods("POST")
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":80", myRouter))
	//log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	Articles = []Article{
		// Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		// Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		Article{Id: "1", User: "User", Token: "xxx"},
		Article{Id: "2", User: "User", Token: "xxx"},
		Article{Id: "3", User: "User", Token: "xxx"},
	}
	handleRequests()
}
