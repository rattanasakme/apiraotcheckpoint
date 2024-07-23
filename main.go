// main.go
package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nfnt/resize"

	"github.com/gorilla/mux"
)

var db *sql.DB // pointer to sql.open which returns a sql.db
var orderid int64
var connPrd = "prd-olivedb.mysql.database.azure.com"
var connuat = "uat-olivedb.mysql.database.azure.com"

var connraot = "10.99.80.51"
var passprd = "!10<>Oms!"

//var conn = connPrd

var tmsapiprd = "https://tms-api.promptsong.co/"
var tmsapiuat = "https://uattms-api.promptsong.co/"

//var tmsapi = tmsapiuat
//var tmsapi = tmsapiprd

var conn = connraot           //connuat //conntest
var userlogin = "nta"         //connuat //admin
var passlogin = "N@xtech!234" //"!10<>Oms!"
var dblogin = "raotdb"        //"THPDMPDB"

// var conn = connuat          //connuat //conntest
// var userlogin = "admin"     //connuat //admin
// var passlogin = "!10<>Oms!" //"!10<>Oms!"
// var dblogin = "THPDMPDB"    //"THPDMPDB"

var tmsapi = tmsapiuat

var raottokenrefresh = ""

// var conn = connPrd
// var tmsapi = tmsapiprd

// Article - Our struct for all articles
//
//	type Article struct {
//		Id      string `json:"Id"`
//		Title   string `json:"Title"`
//		Desc    string `json:"desc"`
//		Content string `json:"content"`
//		Token   string `json:"token"`
//	}
type KTBApproveJson struct {
	User     string  `json:"user"`
	Password string  `json:"password"`
	ComCode  string  `json:"comCode"`
	ProdCode string  `json:"prodCode"`
	Command  string  `json:"command"`
	BankCode int     `json:"bankCode"`
	BankRef  string  `json:"bankRef"`
	DateTime string  `json:"dateTime"`
	EffDate  string  `json:"effDate"`
	Amount   float64 `json:"amount"`
	Channel  string  `json:"channel"`
	CusName  string  `json:"cusName"`
	Ref1     string  `json:"ref1"`
	Ref2     string  `json:"ref2"`
	Ref3     string  `json:"ref3"`
	Ref4     string  `json:"ref4"`
}
type MobileVersion struct {
	MobileType string `json:"MobileType"`
	VersionNow string `json:"VersionNow"`
}

type SSOAccessToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

type PWAApi struct {
	MobileType    string `json:"MobileType"`
	VersionNow    string `json:"VersionNow"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
}
type PWALogin struct {
	UserID        string `json:"UserID"`
	PassWord      string `json:"PassWord"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	CUser         string `json:"CUser"`
}
type PWASendWeight struct {
	WHName           string `json:"WHName"`
	Weight           string `json:"Weight"`
	CTime            string `json:"CTime"`
	CUser            string `json:"CUser"`
	CLocal           string `json:"CLocal"`
	Authorization    string `json:"Authorization"`
	Channel          string `json:"Channel"`
	Token            string `json:"Token"`
	EType            string `json:"EType"`
	SearchTxt        string `json:"SearchTxt"`
	SearchStatus     string `json:"SearchStatus"`
	SearchCustomname string `json:"SearchCustomname"`
}

type GetRaotUser struct {
	UserName      string `json:"UserName"`
	Position      string `json:"Position"`
	Department    string `json:"Department"`
	ResponseGroup string `json:"ResponseGroup"`
	Email         string `json:"Email"`
	Permission    string `json:"Permission"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	EType         string `json:"EType"`
}

type RAOTRole struct {
	Roles []string `json:"roles"`
}

type PWAResponseLocation struct {
	ResponseLocationName string `json:"ResponseLocationName"`
	CUser                string `json:"CUser"`
	Authorization        string `json:"Authorization"`
	Channel              string `json:"Channel"`
	Token                string `json:"Token"`
	EType                string `json:"EType"`
}

type PWAWidget struct {
	UserID        string `json:"UserID"`
	IMG           string `json:"IMG"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	EType         string `json:"EType"`
}

type PWASendCheckInOut struct {
	TransportSubID string `json:"TransportSubID"`
	License        string `json:"License"`
	Status         string `json:"Status"`
	CustomerName   string `json:"CustomerName"`
	StatusDate     string `json:"StatusDate"`
	LAT            string `json:"LAT"`
	LNG            string `json:"LNG"`
	Authorization  string `json:"Authorization"`
	Channel        string `json:"Channel"`
	Token          string `json:"Token"`
	EType          string `json:"EType"`
}

type PWASendDashboard struct {
	StartDT       string `json:"StartDT"`
	StopDT        string `json:"StopDT"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
}
type PWASendCheckpointStart struct {
	WHName                 string `json:"WHName"`
	PositionName           string `json:"PositionName"`
	StartDT                string `json:"StartDT"`
	StopDT                 string `json:"StopDT"`
	StDategps              string `json:"StDategps"`
	EtDategps              string `json:"EtDategps"`
	CheckPointID           string `json:"CheckPointID"`
	CheckPointLocationName string `json:"CheckPointLocationName"`
	CheckPointResponse     string `json:"CheckPointResponse"`
	CUser                  string `json:"CUser"`
	User                   string `json:"User"`
	Authorization          string `json:"Authorization"`
	Channel                string `json:"Channel"`
	Token                  string `json:"Token"`
	EType                  string `json:"EType"`
	SearchTxt              string `json:"SearchTxt"`
	SearchCustomname       string `json:"SearchCustomname"`
	SearchDT               string `json:"SearchDT"`
	SearchRubberType       string `json:"SearchRubberType"`
	SearchResponse         string `json:"SearchResponse"`
	SearchLocationName     string `json:"SearchLocationName"`
}
type PWASendCheckpointStartJob struct {
	JobID         string `json:"JobID"`
	CessID        string `json:"CessID"`
	TranID        string `json:"TranID"`
	CLocal        string `json:"CLocal"`
	CUser         string `json:"CUser"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	EType         string `json:"EType"`
}
type PWASendConsent struct {
	ConsentMsg    string `json:"ConsentMsg"`
	CLocal        string `json:"CLocal"`
	CUser         string `json:"CUser"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	EType         string `json:"EType"`
}

type PWASetCalculate struct {
	TruckID         string `json:"TruckID"`
	TruckRubberType string `json:"TruckRubberType"`
	ContainerType   string `json:"ContainerType"`
	ContainerTypeID string `json:"ContainerTypeID"`
	TruckWeight     string `json:"TruckWeight"`
	TrailWeight     string `json:"TrailWeight"`
	TrailCount      string `json:"TrailCount"`
	ContainerWeight string `json:"ContainerWeight"`
	ContainerCount  string `json:"ContainerCount"`
	BoxWeight       string `json:"BoxWeight"`
	BoxCount        string `json:"BoxCount"`
	Authorization   string `json:"Authorization"`
	Channel         string `json:"Channel"`
	Token           string `json:"Token"`
	CUser           string `json:"CUser"`
	EType           string `json:"EType"`
}
type PWASaveWeight struct {
	TranSubID     string `json:"TranSubID"`
	Weight        string `json:"Weight"`
	WeightCalcEst string `json:"WeightCalcEst"`
	CalcType      string `json:"CalcType"`
	WeightID      string `json:"WeightID"`
	IMGBase64     string `json:"IMGBase64"`
	Authorization string `json:"Authorization"`
	Channel       string `json:"Channel"`
	Token         string `json:"Token"`
	EType         string `json:"EType"`
	CUser         string `json:"CUser"`
	SearchTxt     string `json:"SearchTxt"`
}
type PWASaveLicense struct {
	TranSubID       string `json:"TranSubID"`
	TruckLicenseID  string `json:"TruckLicenseID"`
	TruckDistinct   string `json:"TruckDistinct"`
	TruckType       string `json:"TruckType"`
	TruckTypeDetail string `json:"TruckTypeDetail"`
	Authorization   string `json:"Authorization"`
	Channel         string `json:"Channel"`
	CUser           string `json:"CUser"`
	Token           string `json:"Token"`
	EType           string `json:"EType"`
}
type PWASendLocation struct {
	LocationID       string `json:"LocationID"`
	LocationName     string `json:"LocationName"`
	LocationResponse string `json:"LocationResponse"`
	LocationDetail   string `json:"LocationDetail"`
	LocationGPS      string `json:"LocationGPS"`
	Authorization    string `json:"Authorization"`
	Channel          string `json:"Channel"`
	Token            string `json:"Token"`
	EType            string `json:"EType"`
}
type IMGGeneratedLPR struct {
	RecognizedData struct {
		LicenseNumber string  `json:"license_number"`
		Score         float64 `json:"score"`
		Filename      string  `json:"filename"`
		Timestamp     string  `json:"timestamp"`
		Version       int     `json:"version"`
	} `json:"recognized_data"`
}

type IMGGenerated struct {
	ParsedResults []struct {
		TextOverlay struct {
			Lines []struct {
				LineText string `json:"LineText"`
				Words    []struct {
					WordText string  `json:"WordText"`
					Left     float64 `json:"Left"`
					Top      float64 `json:"Top"`
					Height   float64 `json:"Height"`
					Width    float64 `json:"Width"`
				} `json:"Words"`
				MaxHeight float64 `json:"MaxHeight"`
				MinTop    float64 `json:"MinTop"`
			} `json:"Lines"`
			HasOverlay bool `json:"HasOverlay"`
		} `json:"TextOverlay"`
		TextOrientation   string `json:"TextOrientation"`
		FileParseExitCode int    `json:"FileParseExitCode"`
		ParsedText        string `json:"ParsedText"`
		ErrorMessage      string `json:"ErrorMessage"`
		ErrorDetails      string `json:"ErrorDetails"`
	} `json:"ParsedResults"`
	OCRExitCode                  int    `json:"OCRExitCode"`
	IsErroredOnProcessing        bool   `json:"IsErroredOnProcessing"`
	ProcessingTimeInMilliseconds string `json:"ProcessingTimeInMilliseconds"`
	SearchablePDFURL             string `json:"SearchablePDFURL"`
}
type OMISEResponse struct {
	Object          string `json:"object"`
	ID              string `json:"id"`
	Location        string `json:"location"`
	Amount          int    `json:"amount"`
	Net             int    `json:"net"`
	Fee             int    `json:"fee"`
	FeeVat          int    `json:"fee_vat"`
	Interest        int    `json:"interest"`
	InterestVat     int    `json:"interest_vat"`
	FundingAmount   int    `json:"funding_amount"`
	RefundedAmount  int    `json:"refunded_amount"`
	TransactionFees struct {
		FeeFlat string `json:"fee_flat"`
		FeeRate string `json:"fee_rate"`
		VatRate string `json:"vat_rate"`
	} `json:"transaction_fees"`
	PlatformFee struct {
		Fixed      interface{} `json:"fixed"`
		Amount     interface{} `json:"amount"`
		Percentage interface{} `json:"percentage"`
	} `json:"platform_fee"`
	Currency        string      `json:"currency"`
	FundingCurrency string      `json:"funding_currency"`
	IP              interface{} `json:"ip"`
	Refunds         struct {
		Object   string        `json:"object"`
		Data     []interface{} `json:"data"`
		Limit    int           `json:"limit"`
		Offset   int           `json:"offset"`
		Total    int           `json:"total"`
		Location string        `json:"location"`
		Order    string        `json:"order"`
		From     time.Time     `json:"from"`
		To       time.Time     `json:"to"`
	} `json:"refunds"`
	Link        interface{} `json:"link"`
	Description string      `json:"description"`
	Metadata    struct {
	} `json:"metadata"`
	Card struct {
		Object             string      `json:"object"`
		ID                 string      `json:"id"`
		Livemode           bool        `json:"livemode"`
		Location           interface{} `json:"location"`
		Deleted            bool        `json:"deleted"`
		Street1            interface{} `json:"street1"`
		Street2            interface{} `json:"street2"`
		City               string      `json:"city"`
		State              interface{} `json:"state"`
		PhoneNumber        interface{} `json:"phone_number"`
		PostalCode         string      `json:"postal_code"`
		Country            string      `json:"country"`
		Financing          string      `json:"financing"`
		Bank               string      `json:"bank"`
		Brand              string      `json:"brand"`
		Fingerprint        string      `json:"fingerprint"`
		FirstDigits        interface{} `json:"first_digits"`
		LastDigits         string      `json:"last_digits"`
		Name               string      `json:"name"`
		ExpirationMonth    int         `json:"expiration_month"`
		ExpirationYear     int         `json:"expiration_year"`
		SecurityCodeCheck  bool        `json:"security_code_check"`
		TokenizationMethod interface{} `json:"tokenization_method"`
		CreatedAt          time.Time   `json:"created_at"`
	} `json:"card"`
	Source                   interface{} `json:"source"`
	Schedule                 interface{} `json:"schedule"`
	Customer                 interface{} `json:"customer"`
	Dispute                  interface{} `json:"dispute"`
	Transaction              string      `json:"transaction"`
	FailureCode              interface{} `json:"failure_code"`
	FailureMessage           interface{} `json:"failure_message"`
	Status                   string      `json:"status"`
	AuthorizeURI             string      `json:"authorize_uri"`
	ReturnURI                string      `json:"return_uri"`
	CreatedAt                time.Time   `json:"created_at"`
	PaidAt                   time.Time   `json:"paid_at"`
	ExpiresAt                time.Time   `json:"expires_at"`
	ExpiredAt                interface{} `json:"expired_at"`
	ReversedAt               interface{} `json:"reversed_at"`
	ZeroInterestInstallments bool        `json:"zero_interest_installments"`
	Branch                   interface{} `json:"branch"`
	Terminal                 interface{} `json:"terminal"`
	Device                   interface{} `json:"device"`
	Authorized               bool        `json:"authorized"`
	Capturable               bool        `json:"capturable"`
	Capture                  bool        `json:"capture"`
	Disputable               bool        `json:"disputable"`
	Livemode                 bool        `json:"livemode"`
	Refundable               bool        `json:"refundable"`
	Reversed                 bool        `json:"reversed"`
	Reversible               bool        `json:"reversible"`
	Voided                   bool        `json:"voided"`
	Paid                     bool        `json:"paid"`
	Expired                  bool        `json:"expired"`
}
type OMISEPayStructure struct {
	Object       string `json:"object"`
	ID           string `json:"id"`
	Livemode     bool   `json:"livemode"`
	Location     string `json:"location"`
	Used         bool   `json:"used"`
	ChargeStatus string `json:"charge_status"`
	Card         struct {
		Object             string      `json:"object"`
		ID                 string      `json:"id"`
		Livemode           bool        `json:"livemode"`
		Location           interface{} `json:"location"`
		Deleted            bool        `json:"deleted"`
		Street1            interface{} `json:"street1"`
		Street2            interface{} `json:"street2"`
		City               string      `json:"city"`
		State              interface{} `json:"state"`
		PhoneNumber        interface{} `json:"phone_number"`
		PostalCode         string      `json:"postal_code"`
		Country            string      `json:"country"`
		Financing          string      `json:"financing"`
		Bank               string      `json:"bank"`
		Brand              string      `json:"brand"`
		Fingerprint        string      `json:"fingerprint"`
		FirstDigits        interface{} `json:"first_digits"`
		LastDigits         string      `json:"last_digits"`
		Name               string      `json:"name"`
		ExpirationMonth    int         `json:"expiration_month"`
		ExpirationYear     int         `json:"expiration_year"`
		SecurityCodeCheck  bool        `json:"security_code_check"`
		TokenizationMethod interface{} `json:"tokenization_method"`
		CreatedAt          time.Time   `json:"created_at"`
	} `json:"card"`
	CreatedAt time.Time `json:"created_at"`
}
type KTBReponseApproveJson struct {
	TranxID  string  `json:"tranxId"`
	BankRef  string  `json:"bankRef"`
	RespCode int     `json:"respCode"`
	RespMsg  string  `json:"respMsg"`
	Balance  float64 `json:"balance"`
	CusName  string  `json:"cusName"`
	Info     string  `json:"info"`
	Print1   string  `json:"print1"`
	Print2   string  `json:"print2"`
	Print3   string  `json:"print3"`
	Print4   string  `json:"print4"`
	Print5   string  `json:"print5"`
	Print6   string  `json:"print6"`
	Print7   string  `json:"print7"`
}

type KTBChargeJson struct {
	User     string  `json:"user"`
	Password string  `json:"password"`
	ComCode  string  `json:"comCode"`
	ProdCode string  `json:"prodCode"`
	Command  string  `json:"command"`
	BankCode int     `json:"bankCode"`
	BankRef  string  `json:"bankRef"`
	TranxID  string  `json:"tranxId"`
	DateTime string  `json:"dateTime"`
	EffDate  string  `json:"effDate"`
	Amount   float64 `json:"amount"`
	Channel  string  `json:"channel"`
	CusName  string  `json:"cusName"`
	Ref1     string  `json:"ref1"`
	Ref2     string  `json:"ref2"`
	Ref3     string  `json:"ref3"`
	Ref4     string  `json:"ref4"`
}
type KTBChargeResponseJson struct {
	TranxID  string  `json:"tranxId"`
	BankRef  string  `json:"bankRef"`
	RespCode int     `json:"respCode"`
	RespMsg  string  `json:"respMsg"`
	Balance  float64 `json:"balance"`
	CusName  string  `json:"cusName"`
	Info     string  `json:"info"`
	Print1   string  `json:"print1"`
	Print2   string  `json:"print2"`
	Print3   string  `json:"print3"`
	Print4   string  `json:"print4"`
	Print5   string  `json:"print5"`
	Print6   string  `json:"print6"`
	Print7   string  `json:"print7"`
}
type PaymentChannel struct {
	Amount         int      `json:"amount"`
	CurrencyCode   string   `json:"currencyCode"`
	Description    string   `json:"description"`
	InvoiceNo      string   `json:"invoiceNo"`
	MerchantID     string   `json:"merchantID"`
	PaymentChannel []string `json:"paymentChannel"`
	Tokenize       bool     `json:"tokenize"`
}

type LineStruceture struct {
	Destination string `json:"destination"`
	Events      []struct {
		Type    string `json:"type"`
		Message struct {
			Type string `json:"type"`
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"message"`
		WebhookEventID  string `json:"webhookEventId"`
		DeliveryContext struct {
			IsRedelivery bool `json:"isRedelivery"`
		} `json:"deliveryContext"`
		Timestamp int64 `json:"timestamp"`
		Source    struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		ReplyToken string `json:"replyToken"`
		Mode       string `json:"mode"`
	} `json:"events"`
}

type Article struct {
	Id    string `json:"Id"`
	User  string `json:"User"`
	Token string `json:"token"`
	Type  string `json:"Type"`
	Date  string `json:"Date"`
}
type APIIot struct {
	Api_key string `json:"api_key"`
	Sensor  string `json:"sensor"`
	Value1  string `json:"value1"`
	Value2  string `json:"value2"`
	Value3  string `json:"value3"`
}

type RateExtra struct {
	RateID    int    `json:"RateID"`
	RateName  string `json:"RateName"`
	RatePrice string `json:"RatePrice"`
}
type Citicen struct {
	CiticenId string `json:"CiticenId"`
	MobileId  string `json:"MobileId"`
}
type Zipcode struct {
	PostId     string `json:"PostId"`
	Tumbon     string `json:"Tumbon"`
	TumbonEN   string `json:"TumbonEN"`
	District   string `json:"District"`
	DistrictEN string `json:"DistrictEN"`
	Province   string `json:"Province"`
	GPS        string `json:"GPS"`
}
type (
	StringInterfaceMap7 map[string]interface{}
	EventBooking        struct {
		BookingID         string `json:"booking_id"`
		BookingStatus     string `json:"booking_status"`
		Price             string `json:"price"`
		PriceUnit         string `json:"price_unit"`
		TruckPlate        string `json:"truck_plate"`
		DriverName        string `json:"driver_name"`
		Oid               string `json:"oid"`
		Company           string `json:"company"`
		Address           string `json:"address"`
		Contact           string `json:"contact"`
		Phone             string `json:"phone"`
		Email             string `json:"email"`
		Rating            string `json:"rating"`
		BookingScore      string `json:"booking_score"`
		Time              string `json:"time"`
		PositionLatitude  string `json:"position_latitude"`
		PositionLongitude string `json:"position_longitude"`
		PositionTime      string `json:"position_time"`
	}
)
type (
	StringInterfaceMap77     map[string]interface{}
	EventBookingMatchAlready struct {
		JobID           string `json:"JobID"`
		MobileID        string `json:"MobileID"`
		Receivename     string `json:"Receivename"`
		ReceiveAddress  string `json:"ReceiveAddress"`
		ReceiveTumbon   string `json:"ReceiveTumbon"`
		ReceiveDistrict string `json:"ReceiveDistrict"`
		ReceiveProvince string `json:"ReceiveProvince"`
		ReceiveZipcode  string `json:"ReceiveZipcode"`
		ReceivePhoneNo  string `json:"ReceivePhoneNo"`
		Status          string `json:"Status"`
		PickupStartDt   string `json:"PickupStartDt"`
		TrackingID      string `json:"TrackingID"`
		CarID           string `json:"CarID"`
		DriverName      string `json:"DriverName"`
		DriverPhoneNo   string `json:"DriverPhoneNo"`
		Price           string `json:"Price"`
		CustomerSelect  string `json:"CustomerSelect"`
		WarehouseName   string `json:"WarehouseName"`
		WHAddressTH     string `json:"WHAddressTH"`
	}
)
type TruckRAOT struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}
type TruckRAOTDetail []struct {
	ID                int       `json:"id"`
	LicenseNumberTh   string    `json:"license_number_th"`
	LicenseProvinceTh string    `json:"license_province_th"`
	LicenseNumberMy   string    `json:"license_number_my"`
	RegistrationType  string    `json:"registration_type"`
	TruckType         string    `json:"truck_type"`
	WeightKg          string    `json:"weight_kg"`
	GpsBoxID          string    `json:"gps_box_id"`
	GpsProviderName   string    `json:"gps_provider_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Company           string    `json:"company"`
	Status            string    `json:"status"`
}
type TruckRAOTDetailTrader []struct {
	ID                int       `json:"id"`
	LicenseNumberTh   string    `json:"license_number_th"`
	LicenseProvinceTh string    `json:"license_province_th_name"`
	LicenseNumberMy   string    `json:"license_number_my"`
	RegistrationType  string    `json:"registration_type"`
	TruckType         string    `json:"truck_type"`
	WeightKg          string    `json:"weight_kg"`
	GpsBoxID          string    `json:"gps_box_id"`
	GpsProviderName   string    `json:"gps_company_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Company           string    `json:"company_name"`
	Status            string    `json:"status"`
}
type TrailersRAOTDetailTrader []struct {
	ID                int       `json:"id"`
	LicenseNumberTh   string    `json:"license_number"`
	LicenseProvinceTh string    `json:"license_province_th_name"`
	LicenseNumberMy   string    `json:"license_number_my"`
	RegistrationType  string    `json:"registration_type"`
	TruckType         string    `json:"truck_type"`
	WeightKg          string    `json:"weight_kg"`
	GpsBoxID          string    `json:"gps_box_id"`
	GpsProviderName   string    `json:"gps_company_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Company           string    `json:"company_name"`
	Status            string    `json:"status"`
}
type DriverLoadboard struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		PostRes struct {
			Success  bool          `json:"success"`
			Status   string        `json:"status"`
			Message  string        `json:"message"`
			Error    []interface{} `json:"error"`
			Shipment struct {
				ShipmentID       string `json:"shipment_id"`
				ShipmentNo       string `json:"shipment_no"`
				ShipmentType     string `json:"shipment_type"`
				MatchMode        string `json:"match_mode"`
				MatchPriority2   string `json:"match_priority_2"`
				MatchPriority3   string `json:"match_priority_3"`
				AutoMatchTime    string `json:"auto_match_time"`
				CustomerRating   string `json:"customer_rating"`
				RecipientName    string `json:"recipient_name"`
				Owner            string `json:"owner"`
				OwnerContact     string `json:"owner_contact"`
				OwnerTel         string `json:"owner_tel"`
				OwnerRating      int    `json:"owner_rating"`
				OwnerRatingCount int    `json:"owner_rating_count"`
				Broker           struct {
					UUID    string `json:"uuid"`
					Title   string `json:"title"`
					Address string `json:"address"`
					Contact string `json:"contact"`
					Phone   string `json:"phone"`
					LogoURL string `json:"logo_url"`
				} `json:"broker"`
				Carrier struct {
					UUID         string `json:"uuid"`
					Title        string `json:"title"`
					Address      string `json:"address"`
					Contact      string `json:"contact"`
					Phone        string `json:"phone"`
					LogoURL      string `json:"logo_url"`
					Rating       string `json:"rating"`
					RatingCount  string `json:"rating_count"`
					BookingScore string `json:"booking_score"`
				} `json:"carrier"`
				TruckPlate             string  `json:"truck_plate"`
				DriverName             string  `json:"driver_name"`
				DriverImgURL           string  `json:"driver_img_url"`
				SignatureURL           string  `json:"signature_url"`
				Age                    string  `json:"age"`
				PickupDate             string  `json:"pickup_date"`
				DeliveryDate           string  `json:"delivery_date"`
				OriginLocation         string  `json:"origin_location"`
				OriginAddress          string  `json:"origin_address"`
				OriginSubdistrict      string  `json:"origin_subdistrict"`
				OriginDistrict         string  `json:"origin_district"`
				OriginProvince         string  `json:"origin_province"`
				OriginPostal           string  `json:"origin_postal"`
				OriginLatitude         string  `json:"origin_latitude"`
				OriginLongitude        string  `json:"origin_longitude"`
				DestinationLocation    string  `json:"destination_location"`
				DestinationAddress     string  `json:"destination_address"`
				DestinationSubdistrict string  `json:"destination_subdistrict"`
				DestinationDistrict    string  `json:"destination_district"`
				DestinationProvince    string  `json:"destination_province"`
				DestinationPostal      string  `json:"destination_postal"`
				DestinationLatitude    string  `json:"destination_latitude"`
				DestinationLongitude   string  `json:"destination_longitude"`
				JobCount               string  `json:"job_count"`
				Trip                   string  `json:"trip"`
				Weight                 float64 `json:"weight"`
				TruckType              string  `json:"truck_type"`
				TruckTypeTitle         string  `json:"truck_type_title"`
				PriceStandard          string  `json:"price_standard"`
				PriceExpress           string  `json:"price_express"`
				PriceFestival          string  `json:"price_festival"`
				PriceAsking            string  `json:"price_asking"`
				PriceOffer             string  `json:"price_offer"`
				PriceUnit              string  `json:"price_unit"`
				PriceRateWeight        string  `json:"price_rate_weight"`
				PriceRateDistance      string  `json:"price_rate_distance"`
				PaymentType            string  `json:"payment_type"`
				BoxAmount              string  `json:"box_amount"`
				DeliveryRemark         string  `json:"delivery_remark"`
				StatusID               string  `json:"status_id"`
				Status                 string  `json:"status"`
				Remark                 string  `json:"remark"`
				CancelReason           string  `json:"cancel_reason"`
				JobData                []struct {
					RefOrderNo       string `json:"ref_order_no"`
					ConNo            string `json:"con_no"`
					Type             string `json:"type"`
					PickupDate       string `json:"pickup_date"`
					DeliveryDate     string `json:"delivery_date"`
					BoxAmount        string `json:"box_amount"`
					TotalWeight      int    `json:"total_weight"`
					TotalCbm         int    `json:"total_cbm"`
					Endpoint         string `json:"endpoint"`
					EndpointUID      string `json:"endpoint_uid"`
					EndpointOid      string `json:"endpoint_oid"`
					NodeID           string `json:"node_id"`
					TransportPrice   string `json:"transport_price"`
					SName            string `json:"s_name"`
					SAddress         string `json:"s_address"`
					SSubdistrict     string `json:"s_subdistrict"`
					SDistrict        string `json:"s_district"`
					SProvince        string `json:"s_province"`
					SLat             string `json:"s_lat"`
					SLon             string `json:"s_lon"`
					SZipcode         string `json:"s_zipcode"`
					STel             string `json:"s_tel"`
					SEmail           string `json:"s_email"`
					SContact         string `json:"s_contact"`
					SRefCode         string `json:"s_ref_code"`
					SRefLink         string `json:"s_ref_link"`
					RName            string `json:"r_name"`
					RAddress         string `json:"r_address"`
					RSubdistrict     string `json:"r_subdistrict"`
					RDistrict        string `json:"r_district"`
					RProvince        string `json:"r_province"`
					RLat             string `json:"r_lat"`
					RLon             string `json:"r_lon"`
					RZipcode         string `json:"r_zipcode"`
					RTel             string `json:"r_tel"`
					REmail           string `json:"r_email"`
					RContact         string `json:"r_contact"`
					RRefCode         string `json:"r_ref_code"`
					RRefLink         string `json:"r_ref_link"`
					CName            string `json:"c_name"`
					CAddress         string `json:"c_address"`
					CSubdistrict     string `json:"c_subdistrict"`
					CDistrict        string `json:"c_district"`
					CProvince        string `json:"c_province"`
					CZipcode         string `json:"c_zipcode"`
					CTel             string `json:"c_tel"`
					CEmail           string `json:"c_email"`
					CContact         string `json:"c_contact"`
					CRefCode         string `json:"c_ref_code"`
					ServiceCode      string `json:"service_code"`
					CodType          string `json:"cod_type"`
					CodPrice         string `json:"cod_price"`
					TransportCompany string `json:"transport_company"`
					CompanyCode      string `json:"company_code"`
					GroupRefID       string `json:"group_ref_id"`
					SignatureURL     string `json:"signature_url"`
					SignatureName    string `json:"signature_name"`
					Remark           string `json:"remark"`
					Status           string `json:"status"`
					Created          string `json:"created"`
					CreatedUID       string `json:"created_uid"`
					Removed          string `json:"removed"`
					RemovedUID       string `json:"removed_uid"`
					ProductDetail    []struct {
						RefID       string `json:"ref_id"`
						Name        string `json:"name"`
						Sku         string `json:"sku"`
						Code        string `json:"code"`
						Qty         string `json:"qty"`
						Unit        string `json:"unit"`
						Width       string `json:"width"`
						Length      string `json:"length"`
						Height      string `json:"height"`
						Weight      string `json:"weight"`
						GrossWeight string `json:"gross_weight"`
						TotalWeight string `json:"total_weight"`
						TotalCbm    string `json:"total_cbm"`
						Status      string `json:"status"`
						Created     string `json:"created"`
						Changed     string `json:"changed"`
					} `json:"product_detail"`
				} `json:"job_data"`
				BookingCount int `json:"booking_count"`
				Booking      []struct {
					BookingID         string `json:"booking_id"`
					BookingStatus     string `json:"booking_status"`
					Price             string `json:"price"`
					PriceUnit         string `json:"price_unit"`
					TruckPlate        string `json:"truck_plate"`
					DriverName        string `json:"driver_name"`
					UUID              string `json:"uuid"`
					Company           string `json:"company"`
					Address           string `json:"address"`
					Contact           string `json:"contact"`
					Phone             string `json:"phone"`
					Email             string `json:"email"`
					LogoURL           string `json:"logo_url"`
					Rating            string `json:"rating"`
					BookingScore      string `json:"booking_score"`
					Time              string `json:"time"`
					PositionLatitude  string `json:"position_latitude"`
					PositionLongitude string `json:"position_longitude"`
					PositionTime      string `json:"position_time"`
					Oid               string `json:"oid"`
				} `json:"booking"`
			} `json:"shipment"`
		} `json:"post_res"`
		To struct {
			ToID                    string      `json:"to_id"`
			Oid                     string      `json:"oid"`
			OwnerID                 string      `json:"owner_id"`
			Owner                   string      `json:"owner"`
			OwnerRating             string      `json:"owner_rating"`
			OwnerRatingCount        string      `json:"owner_rating_count"`
			MatchMode               string      `json:"match_mode"`
			MatchTime               interface{} `json:"match_time"`
			MatchPriority2          string      `json:"match_priority_2"`
			MatchPriority3          string      `json:"match_priority_3"`
			ToNumber                string      `json:"to_number"`
			TransportType           string      `json:"transport_type"`
			NodeID                  string      `json:"node_id"`
			TransportStatus         string      `json:"transport_status"`
			CarrierUUID             interface{} `json:"carrier_uuid"`
			CarrierName             interface{} `json:"carrier_name"`
			TruckID                 string      `json:"truck_id"`
			TruckType               string      `json:"truck_type"`
			TruckPlate              interface{} `json:"truck_plate"`
			TailPlateID             string      `json:"tail_plate_id"`
			TailPlate               string      `json:"tail_plate"`
			Cooling                 string      `json:"cooling"`
			Gps                     string      `json:"gps"`
			Lift                    string      `json:"lift"`
			DriverID                string      `json:"driver_id"`
			DriverName              string      `json:"driver_name"`
			DriverPhone             string      `json:"driver_phone"`
			DriverMail              string      `json:"driver_mail"`
			DatePickupPlan          string      `json:"date_pickup_plan"`
			DateDeliveryPlan        string      `json:"date_delivery_plan"`
			DatePickupActual        string      `json:"date_pickup_actual"`
			DateDeliveryActual      string      `json:"date_delivery_actual"`
			LocationOriginID        string      `json:"location_origin_id"`
			LocationDestinationID   string      `json:"location_destination_id"`
			OriginAddress           string      `json:"origin_address"`
			DestinationAddress      string      `json:"destination_address"`
			LocationOriginName      string      `json:"location_origin_name"`
			LocationDestinationName string      `json:"location_destination_name"`
			SourceType              string      `json:"source_type"`
			OperationType           string      `json:"operation_type"`
			Quantity                string      `json:"quantity"`
			Weight                  string      `json:"weight"`
			Volume                  string      `json:"volume"`
			Remark                  string      `json:"remark"`
			StandardPrice           string      `json:"standard_price"`
			AskingPrice             string      `json:"asking_price"`
			OfferPrice              string      `json:"offer_price"`
			TransportPrice          string      `json:"transport_price"`
			InsurancePrice          string      `json:"insurance_price"`
			PriceRateWeight         interface{} `json:"price_rate_weight"`
			PriceRateDistance       interface{} `json:"price_rate_distance"`
			PriceUnit               string      `json:"price_unit"`
			AssignDate              string      `json:"assign_date"`
			AcceptDate              string      `json:"accept_date"`
			OriginArriveDate        string      `json:"origin_arrive_date"`
			DestinationArriveDate   string      `json:"destination_arrive_date"`
			CompleteDate            string      `json:"complete_date"`
			AFailedDerveryDate      string      `json:"a_failed_dervery_date"`
			ASortingDate            string      `json:"a_sorting_date"`
			AScanReceiveDate        string      `json:"a_scan_receive_date"`
			CanceledDate            string      `json:"canceled_date"`
			UID                     string      `json:"uid"`
			ChangedUID              string      `json:"changed_uid"`
			Created                 string      `json:"created"`
			Changed                 string      `json:"changed"`
			Gate                    interface{} `json:"gate"`
			RealTruckID             string      `json:"real_truck_id"`
			RealDriverID            string      `json:"real_driver_id"`
			CashFlag                string      `json:"cash_flag"`
			GrossWeight             string      `json:"gross_weight"`
			GrossVolume             string      `json:"gross_volume"`
			EmergencyRemark         interface{} `json:"emergency_remark"`
			ToTrack                 string      `json:"to_track"`
			TotalDuration           string      `json:"total_duration"`
			TotalDistance           string      `json:"total_distance"`
			DriverName2             string      `json:"driver_name2"`
			RealDriverID2           string      `json:"real_driver_id2"`
			TempRate                string      `json:"temp_rate"`
			RefNumber               interface{} `json:"ref_number"`
			Line                    string      `json:"line"`
			PaymentType             string      `json:"payment_type"`
			PaymentOption           string      `json:"payment_option"`
			IsPost                  string      `json:"is_post"`
			IsPaid                  string      `json:"is_paid"`
			FromLoad                string      `json:"from_load"`
			LbID                    string      `json:"lb_id"`
			WithdrawID              string      `json:"withdraw_id"`
			CusRatingTruck          string      `json:"cus_rating_truck"`
			CusRatingTruckComment   interface{} `json:"cus_rating_truck_comment"`
			TruckRatingCus          string      `json:"truck_rating_cus"`
			TruckRatingCusComment   string      `json:"truck_rating_cus_comment"`
			ShipmentNo              string      `json:"shipment_no"`
		} `json:"to"`
	} `json:"data"`
}

// type DriverLoadboard struct {
// 	Success bool   `json:"success"`
// 	Msg     string `json:"msg"`
// 	Data    struct {
// 		PostRes struct {
// 			Success  bool   `json:"success"`
// 			Status   string `json:"status"`
// 			Message  string `json:"message"`
// 			Shipment struct {
// 				Sid            string `json:"sid"`
// 				ShipmentNo     string `json:"shipment_no"`
// 				ShipmentType   string `json:"shipment_type"`
// 				MatchMode      string `json:"match_mode"`
// 				MatchPriority2 string `json:"match_priority_2"`
// 				MatchPriority3 string `json:"match_priority_3"`
// 				AutoMatchTime  string `json:"auto_match_time"`
// 				CustomerRating string `json:"customer_rating"`
// 				RecipientName  string `json:"recipient_name"`
// 				Broker         struct {
// 					Oid     string `json:"oid"`
// 					Title   string `json:"title"`
// 					Address string `json:"address"`
// 					Contact string `json:"contact"`
// 					Phone   string `json:"phone"`
// 				} `json:"broker"`
// 				Carrier                []interface{} `json:"carrier"`
// 				TruckPlate             string        `json:"truck_plate"`
// 				DriverName             string        `json:"driver_name"`
// 				DriverImgURL           string        `json:"driver_img_url"`
// 				SignatureURL           string        `json:"signature_url"`
// 				Age                    string        `json:"age"`
// 				PickupDate             string        `json:"pickup_date"`
// 				DeliveryDate           string        `json:"delivery_date"`
// 				OriginLocation         string        `json:"origin_location"`
// 				OriginAddress          string        `json:"origin_address"`
// 				OriginSubdistrict      string        `json:"origin_subdistrict"`
// 				OriginDistrict         string        `json:"origin_district"`
// 				OriginProvince         string        `json:"origin_province"`
// 				OriginPostal           string        `json:"origin_postal"`
// 				OriginLatitude         string        `json:"origin_latitude"`
// 				OriginLongitude        string        `json:"origin_longitude"`
// 				DestinationLocation    string        `json:"destination_location"`
// 				DestinationAddress     string        `json:"destination_address"`
// 				DestinationSubdistrict string        `json:"destination_subdistrict"`
// 				DestinationDistrict    string        `json:"destination_district"`
// 				DestinationProvince    string        `json:"destination_province"`
// 				DestinationPostal      string        `json:"destination_postal"`
// 				DestinationLatitude    string        `json:"destination_latitude"`
// 				DestinationLongitude   string        `json:"destination_longitude"`
// 				JobCount               string        `json:"job_count"`
// 				Trip                   string        `json:"trip"`
// 				Weight                 float64       `json:"weight"`
// 				TruckType              string        `json:"truck_type"`
// 				TruckTypeTitle         string        `json:"truck_type_title"`
// 				PriceStandard          string        `json:"price_standard"`
// 				PriceExpress           string        `json:"price_express"`
// 				PriceFestival          string        `json:"price_festival"`
// 				PriceAsking            string        `json:"price_asking"`
// 				PriceOffer             string        `json:"price_offer"`
// 				PriceUnit              string        `json:"price_unit"`
// 				PriceRateWeight        string        `json:"price_rate_weight"`
// 				PriceRateDistance      string        `json:"price_rate_distance"`
// 				PaymentType            string        `json:"payment_type"`
// 				BoxAmount              string        `json:"box_amount"`
// 				DeliveryRemark         string        `json:"delivery_remark"`
// 				StatusID               string        `json:"status_id"`
// 				Status                 string        `json:"status"`
// 				Remark                 string        `json:"remark"`
// 				CancelReason           interface{}   `json:"cancel_reason"`
// 				BookingCount           int           `json:"booking_count"`
// 				JobData                []struct {
// 					RefOrderNo       string `json:"ref_order_no"`
// 					ConNo            string `json:"con_no"`
// 					SName            string `json:"s_name"`
// 					SAddress         string `json:"s_address"`
// 					SSubdistrict     string `json:"s_subdistrict"`
// 					SDistrict        string `json:"s_district"`
// 					SProvince        string `json:"s_province"`
// 					SLat             string `json:"s_lat"`
// 					SLon             string `json:"s_lon"`
// 					SZipcode         string `json:"s_zipcode"`
// 					STel             string `json:"s_tel"`
// 					SEmail           string `json:"s_email"`
// 					SContact         string `json:"s_contact"`
// 					SRefCode         string `json:"s_ref_code"`
// 					SRefLink         string `json:"s_ref_link"`
// 					RName            string `json:"r_name"`
// 					RAddress         string `json:"r_address"`
// 					RSubdistrict     string `json:"r_subdistrict"`
// 					RDistrict        string `json:"r_district"`
// 					RProvince        string `json:"r_province"`
// 					RLat             string `json:"r_lat"`
// 					RLon             string `json:"r_lon"`
// 					RZipcode         string `json:"r_zipcode"`
// 					RTel             string `json:"r_tel"`
// 					REmail           string `json:"r_email"`
// 					RContact         string `json:"r_contact"`
// 					RRefCode         string `json:"r_ref_code"`
// 					RRefLink         string `json:"r_ref_link"`
// 					CName            string `json:"c_name"`
// 					CAddress         string `json:"c_address"`
// 					CSubdistrict     string `json:"c_subdistrict"`
// 					CDistrict        string `json:"c_district"`
// 					CProvince        string `json:"c_province"`
// 					CZipcode         string `json:"c_zipcode"`
// 					CTel             string `json:"c_tel"`
// 					CEmail           string `json:"c_email"`
// 					CContact         string `json:"c_contact"`
// 					CRefCode         string `json:"c_ref_code"`
// 					ServiceCode      string `json:"service_code"`
// 					CodType          string `json:"cod_type"`
// 					TransportCompany string `json:"transport_company"`
// 					CompanyCode      string `json:"company_code"`
// 					GroupRefID       string `json:"group_ref_id"`
// 					PickupDate       string `json:"pickup_date"`
// 					DeliveryDate     string `json:"delivery_date"`
// 					BoxAmount        string `json:"box_amount"`
// 					TotalWeight      int    `json:"total_weight"`
// 					TotalCbm         int    `json:"total_cbm"`
// 					Endpoint         string `json:"endpoint"`
// 					EndpointUID      string `json:"endpoint_uid"`
// 					EndpointOid      string `json:"endpoint_oid"`
// 					NodeID           string `json:"node_id"`
// 					TransportPrice   string `json:"transport_price"`
// 					ProductDetail    []struct {
// 						RefID       string `json:"ref_id"`
// 						Name        string `json:"name"`
// 						Sku         string `json:"sku"`
// 						Code        string `json:"code"`
// 						Qty         string `json:"qty"`
// 						Unit        string `json:"unit"`
// 						Width       string `json:"width"`
// 						Length      string `json:"length"`
// 						Height      string `json:"height"`
// 						Weight      string `json:"weight"`
// 						GrossWeight string `json:"gross_weight"`
// 						TotalWeight string `json:"total_weight"`
// 						TotalCbm    string `json:"total_cbm"`
// 					} `json:"product_detail"`
// 				} `json:"job_data"`
// 				Owner        string `json:"owner"`
// 				OwnerContact string `json:"owner_contact"`
// 				OwnerTel     string `json:"owner_tel"`
// 			} `json:"shipment"`
// 			Booking []struct {
// 				BookingID         string `json:"booking_id"`
// 				BookingStatus     string `json:"booking_status"`
// 				Price             string `json:"price"`
// 				PriceUnit         string `json:"price_unit"`
// 				TruckPlate        string `json:"truck_plate"`
// 				DriverName        string `json:"driver_name"`
// 				Oid               string `json:"oid"`
// 				Company           string `json:"company"`
// 				Address           string `json:"address"`
// 				Contact           string `json:"contact"`
// 				Phone             string `json:"phone"`
// 				Email             string `json:"email"`
// 				Rating            string `json:"rating"`
// 				BookingScore      string `json:"booking_score"`
// 				Time              string `json:"time"`
// 				PositionLatitude  string `json:"position_latitude"`
// 				PositionLongitude string `json:"position_longitude"`
// 				PositionTime      string `json:"position_time"`
// 			} `json:"booking"`
// 			BookingCount    int           `json:"booking_count"`
// 			MatchedBackhaul []interface{} `json:"matched_backhaul"`
// 			Backhaul        []interface{} `json:"backhaul"`
// 			ShipmentUpdate  []interface{} `json:"shipment_update"`
// 			ShipmentImage   []interface{} `json:"shipment_image"`
// 		} `json:"post_res"`
// 		To struct {
// 			ShipmentNo string `json:"shipment_no"`
// 		} `json:"to"`
// 	} `json:"data"`
// }

//	type DriverLoadboard struct {
//		Success bool   `json:"success"`
//		Msg     string `json:"msg"`
//		Data    struct {
//			Success  bool   `json:"success"`
//			Status   string `json:"status"`
//			Message  string `json:"message"`
//			Shipment struct {
//				ShipmentNo     string `json:"shipment_no"`
//				ShipmentType   string `json:"shipment_type"`
//				MatchMode      string `json:"match_mode"`
//				MatchPriority2 string `json:"match_priority_2"`
//				MatchPriority3 string `json:"match_priority_3"`
//				AutoMatchTime  string `json:"auto_match_time"`
//				CustomerRating string `json:"customer_rating"`
//				RecipientName  string `json:"recipient_name"`
//				Broker         struct {
//					Oid     string `json:"oid"`
//					Title   string `json:"title"`
//					Address string `json:"address"`
//					Contact string `json:"contact"`
//					Phone   string `json:"phone"`
//				} `json:"broker"`
//				Owner   string `json:"owner"`
//				Carrier struct {
//					Oid     string `json:"oid"`
//					Title   string `json:"title"`
//					Address string `json:"address"`
//					Contact string `json:"contact"`
//					Phone   string `json:"phone"`
//				} `json:"carrier"`
//				TruckPlate             string      `json:"truck_plate"`
//				DriverName             string      `json:"driver_name"`
//				DriverImgURL           string      `json:"driver_img_url"`
//				SignatureURL           string      `json:"signature_url"`
//				Age                    string      `json:"age"`
//				PickupDate             string      `json:"pickup_date"`
//				DeliveryDate           string      `json:"delivery_date"`
//				OriginLocation         string      `json:"origin_location"`
//				OriginAddress          string      `json:"origin_address"`
//				OriginSubdistrict      string      `json:"origin_subdistrict"`
//				OriginDistrict         string      `json:"origin_district"`
//				OriginProvince         string      `json:"origin_province"`
//				OriginPostal           string      `json:"origin_postal"`
//				OriginLatitude         string      `json:"origin_latitude"`
//				OriginLongitude        string      `json:"origin_longitude"`
//				DestinationLocation    string      `json:"destination_location"`
//				DestinationAddress     string      `json:"destination_address"`
//				DestinationSubdistrict string      `json:"destination_subdistrict"`
//				DestinationDistrict    string      `json:"destination_district"`
//				DestinationProvince    string      `json:"destination_province"`
//				DestinationPostal      string      `json:"destination_postal"`
//				DestinationLatitude    string      `json:"destination_latitude"`
//				DestinationLongitude   string      `json:"destination_longitude"`
//				JobCount               string      `json:"job_count"`
//				Trip                   string      `json:"trip"`
//				Weight                 int         `json:"weight"`
//				TruckType              string      `json:"truck_type"`
//				TruckTypeTitle         string      `json:"truck_type_title"`
//				PriceStandard          string      `json:"price_standard"`
//				PriceExpress           string      `json:"price_express"`
//				PriceFestival          string      `json:"price_festival"`
//				PriceAsking            string      `json:"price_asking"`
//				PriceOffer             string      `json:"price_offer"`
//				PriceUnit              string      `json:"price_unit"`
//				PriceRateWeight        string      `json:"price_rate_weight"`
//				PriceRateDistance      string      `json:"price_rate_distance"`
//				PaymentType            string      `json:"payment_type"`
//				BoxAmount              string      `json:"box_amount"`
//				DeliveryRemark         string      `json:"delivery_remark"`
//				StatusID               string      `json:"status_id"`
//				Status                 string      `json:"status"`
//				Remark                 string      `json:"remark"`
//				CancelReason           interface{} `json:"cancel_reason"`
//				BookingCount           int         `json:"booking_count"`
//				JobData                []struct {
//					RefOrderNo       string `json:"ref_order_no"`
//					ConNo            string `json:"con_no"`
//					SName            string `json:"s_name"`
//					SAddress         string `json:"s_address"`
//					SSubdistrict     string `json:"s_subdistrict"`
//					SDistrict        string `json:"s_district"`
//					SProvince        string `json:"s_province"`
//					SLat             string `json:"s_lat"`
//					SLon             string `json:"s_lon"`
//					SZipcode         string `json:"s_zipcode"`
//					STel             string `json:"s_tel"`
//					SEmail           string `json:"s_email"`
//					SContact         string `json:"s_contact"`
//					SRefCode         string `json:"s_ref_code"`
//					SRefLink         string `json:"s_ref_link"`
//					RName            string `json:"r_name"`
//					RAddress         string `json:"r_address"`
//					RSubdistrict     string `json:"r_subdistrict"`
//					RDistrict        string `json:"r_district"`
//					RProvince        string `json:"r_province"`
//					RLat             string `json:"r_lat"`
//					RLon             string `json:"r_lon"`
//					RZipcode         string `json:"r_zipcode"`
//					RTel             string `json:"r_tel"`
//					REmail           string `json:"r_email"`
//					RContact         string `json:"r_contact"`
//					RRefCode         string `json:"r_ref_code"`
//					RRefLink         string `json:"r_ref_link"`
//					CName            string `json:"c_name"`
//					CAddress         string `json:"c_address"`
//					CSubdistrict     string `json:"c_subdistrict"`
//					CDistrict        string `json:"c_district"`
//					CProvince        string `json:"c_province"`
//					CZipcode         string `json:"c_zipcode"`
//					CTel             string `json:"c_tel"`
//					CEmail           string `json:"c_email"`
//					CContact         string `json:"c_contact"`
//					CRefCode         string `json:"c_ref_code"`
//					ServiceCode      string `json:"service_code"`
//					CodType          string `json:"cod_type"`
//					TransportCompany string `json:"transport_company"`
//					CompanyCode      string `json:"company_code"`
//					GroupRefID       string `json:"group_ref_id"`
//					PickupDate       string `json:"pickup_date"`
//					DeliveryDate     string `json:"delivery_date"`
//					BoxAmount        string `json:"box_amount"`
//					TotalWeight      int    `json:"total_weight"`
//					TotalCbm         int    `json:"total_cbm"`
//					Endpoint         string `json:"endpoint"`
//					EndpointUID      string `json:"endpoint_uid"`
//					EndpointOid      string `json:"endpoint_oid"`
//					NodeID           string `json:"node_id"`
//					TransportPrice   string `json:"transport_price"`
//					ProductDetail    []struct {
//						RefID       string `json:"ref_id"`
//						Name        string `json:"name"`
//						Sku         string `json:"sku"`
//						Code        string `json:"code"`
//						Qty         string `json:"qty"`
//						Unit        string `json:"unit"`
//						Width       string `json:"width"`
//						Length      string `json:"length"`
//						Height      string `json:"height"`
//						Weight      string `json:"weight"`
//						GrossWeight string `json:"gross_weight"`
//						TotalWeight string `json:"total_weight"`
//						TotalCbm    string `json:"total_cbm"`
//					} `json:"product_detail"`
//				} `json:"job_data"`
//			} `json:"shipment"`
//			Booking []struct {
//				BookingID         string `json:"booking_id"`
//				BookingStatus     string `json:"booking_status"`
//				Price             string `json:"price"`
//				PriceUnit         string `json:"price_unit"`
//				TruckPlate        string `json:"truck_plate"`
//				DriverName        string `json:"driver_name"`
//				Oid               string `json:"oid"`
//				Company           string `json:"company"`
//				Address           string `json:"address"`
//				Contact           string `json:"contact"`
//				Phone             string `json:"phone"`
//				Email             string `json:"email"`
//				Rating            string `json:"rating"`
//				BookingScore      string `json:"booking_score"`
//				Time              string `json:"time"`
//				PositionLatitude  string `json:"position_latitude"`
//				PositionLongitude string `json:"position_longitude"`
//				PositionTime      string `json:"position_time"`
//			} `json:"booking"`
//			ShipmentUpdate struct {
//				Num40 struct {
//					ShipmentStatus     string `json:"shipment_status"`
//					ActionTime         string `json:"action_time"`
//					PositionLat        string `json:"position_lat"`
//					PositionLon        string `json:"position_lon"`
//					PositionTime       string `json:"position_time"`
//					DistributionCenter string `json:"distribution_center"`
//					Status             string `json:"status"`
//					UID                string `json:"uid"`
//					Created            string `json:"created"`
//					Changed            string `json:"changed"`
//				} `json:"40"`
//				Num59 struct {
//					ShipmentStatus     string `json:"shipment_status"`
//					ActionTime         string `json:"action_time"`
//					PositionLat        string `json:"position_lat"`
//					PositionLon        string `json:"position_lon"`
//					PositionTime       string `json:"position_time"`
//					DistributionCenter string `json:"distribution_center"`
//					Status             string `json:"status"`
//					UID                string `json:"uid"`
//					Created            string `json:"created"`
//					Changed            string `json:"changed"`
//				} `json:"59"`
//			} `json:"shipment_update"`
//			ShipmentImage []interface{} `json:"shipment_image"`
//		} `json:"data"`
//	}
type DriverLoadboard0 struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    struct {
		Success  bool   `json:"success"`
		Status   string `json:"status"`
		Message  string `json:"message"`
		Shipment struct {
			ShipmentNo     string `json:"shipment_no"`
			ShipmentType   string `json:"shipment_type"`
			MatchMode      string `json:"match_mode"`
			MatchPriority2 string `json:"match_priority_2"`
			MatchPriority3 string `json:"match_priority_3"`
			AutoMatchTime  string `json:"auto_match_time"`
			CustomerRating string `json:"customer_rating"`
			RecipientName  string `json:"recipient_name"`
			Broker         struct {
				Oid     string `json:"oid"`
				Title   string `json:"title"`
				Address string `json:"address"`
				Contact string `json:"contact"`
				Phone   string `json:"phone"`
			} `json:"broker"`
			Owner   string `json:"owner"`
			Carrier struct {
				Oid     string `json:"oid"`
				Title   string `json:"title"`
				Address string `json:"address"`
				Contact string `json:"contact"`
				Phone   string `json:"phone"`
			} `json:"carrier"`
			TruckPlate             string      `json:"truck_plate"`
			DriverName             string      `json:"driver_name"`
			DriverImgURL           string      `json:"driver_img_url"`
			SignatureURL           string      `json:"signature_url"`
			Age                    string      `json:"age"`
			PickupDate             string      `json:"pickup_date"`
			DeliveryDate           string      `json:"delivery_date"`
			OriginLocation         string      `json:"origin_location"`
			OriginAddress          string      `json:"origin_address"`
			OriginSubdistrict      string      `json:"origin_subdistrict"`
			OriginDistrict         string      `json:"origin_district"`
			OriginProvince         string      `json:"origin_province"`
			OriginPostal           string      `json:"origin_postal"`
			OriginLatitude         string      `json:"origin_latitude"`
			OriginLongitude        string      `json:"origin_longitude"`
			DestinationLocation    string      `json:"destination_location"`
			DestinationAddress     string      `json:"destination_address"`
			DestinationSubdistrict string      `json:"destination_subdistrict"`
			DestinationDistrict    string      `json:"destination_district"`
			DestinationProvince    string      `json:"destination_province"`
			DestinationPostal      string      `json:"destination_postal"`
			DestinationLatitude    string      `json:"destination_latitude"`
			DestinationLongitude   string      `json:"destination_longitude"`
			JobCount               string      `json:"job_count"`
			Trip                   string      `json:"trip"`
			Weight                 int         `json:"weight"`
			TruckType              string      `json:"truck_type"`
			TruckTypeTitle         string      `json:"truck_type_title"`
			PriceStandard          string      `json:"price_standard"`
			PriceExpress           string      `json:"price_express"`
			PriceFestival          string      `json:"price_festival"`
			PriceAsking            string      `json:"price_asking"`
			PriceOffer             string      `json:"price_offer"`
			PriceUnit              string      `json:"price_unit"`
			PriceRateWeight        string      `json:"price_rate_weight"`
			PriceRateDistance      string      `json:"price_rate_distance"`
			PaymentType            string      `json:"payment_type"`
			BoxAmount              string      `json:"box_amount"`
			DeliveryRemark         string      `json:"delivery_remark"`
			StatusID               string      `json:"status_id"`
			Status                 string      `json:"status"`
			Remark                 string      `json:"remark"`
			CancelReason           interface{} `json:"cancel_reason"`
			BookingCount           int         `json:"booking_count"`
			JobData                []struct {
				RefOrderNo       string `json:"ref_order_no"`
				ConNo            string `json:"con_no"`
				SName            string `json:"s_name"`
				SAddress         string `json:"s_address"`
				SSubdistrict     string `json:"s_subdistrict"`
				SDistrict        string `json:"s_district"`
				SProvince        string `json:"s_province"`
				SLat             string `json:"s_lat"`
				SLon             string `json:"s_lon"`
				SZipcode         string `json:"s_zipcode"`
				STel             string `json:"s_tel"`
				SEmail           string `json:"s_email"`
				SContact         string `json:"s_contact"`
				SRefCode         string `json:"s_ref_code"`
				SRefLink         string `json:"s_ref_link"`
				RName            string `json:"r_name"`
				RAddress         string `json:"r_address"`
				RSubdistrict     string `json:"r_subdistrict"`
				RDistrict        string `json:"r_district"`
				RProvince        string `json:"r_province"`
				RLat             string `json:"r_lat"`
				RLon             string `json:"r_lon"`
				RZipcode         string `json:"r_zipcode"`
				RTel             string `json:"r_tel"`
				REmail           string `json:"r_email"`
				RContact         string `json:"r_contact"`
				RRefCode         string `json:"r_ref_code"`
				RRefLink         string `json:"r_ref_link"`
				CName            string `json:"c_name"`
				CAddress         string `json:"c_address"`
				CSubdistrict     string `json:"c_subdistrict"`
				CDistrict        string `json:"c_district"`
				CProvince        string `json:"c_province"`
				CZipcode         string `json:"c_zipcode"`
				CTel             string `json:"c_tel"`
				CEmail           string `json:"c_email"`
				CContact         string `json:"c_contact"`
				CRefCode         string `json:"c_ref_code"`
				ServiceCode      string `json:"service_code"`
				CodType          string `json:"cod_type"`
				TransportCompany string `json:"transport_company"`
				CompanyCode      string `json:"company_code"`
				GroupRefID       string `json:"group_ref_id"`
				PickupDate       string `json:"pickup_date"`
				DeliveryDate     string `json:"delivery_date"`
				BoxAmount        string `json:"box_amount"`
				TotalWeight      int    `json:"total_weight"`
				TotalCbm         int    `json:"total_cbm"`
				Endpoint         string `json:"endpoint"`
				EndpointUID      string `json:"endpoint_uid"`
				EndpointOid      string `json:"endpoint_oid"`
				NodeID           string `json:"node_id"`
				TransportPrice   string `json:"transport_price"`
				ProductDetail    []struct {
					RefID       string `json:"ref_id"`
					Name        string `json:"name"`
					Sku         string `json:"sku"`
					Code        string `json:"code"`
					Qty         string `json:"qty"`
					Unit        string `json:"unit"`
					Width       string `json:"width"`
					Length      string `json:"length"`
					Height      string `json:"height"`
					Weight      string `json:"weight"`
					GrossWeight string `json:"gross_weight"`
					TotalWeight string `json:"total_weight"`
					TotalCbm    string `json:"total_cbm"`
				} `json:"product_detail"`
			} `json:"job_data"`
		} `json:"shipment"`
		Booking []struct {
			BookingID         string `json:"booking_id"`
			BookingStatus     string `json:"booking_status"`
			Price             string `json:"price"`
			PriceUnit         string `json:"price_unit"`
			TruckPlate        string `json:"truck_plate"`
			DriverName        string `json:"driver_name"`
			Oid               string `json:"oid"`
			Company           string `json:"company"`
			Address           string `json:"address"`
			Contact           string `json:"contact"`
			Phone             string `json:"phone"`
			Email             string `json:"email"`
			Rating            string `json:"rating"`
			BookingScore      string `json:"booking_score"`
			Time              string `json:"time"`
			PositionLatitude  string `json:"position_latitude"`
			PositionLongitude string `json:"position_longitude"`
			PositionTime      string `json:"position_time"`
		} `json:"booking"`
		ShipmentUpdate struct {
			Num40 struct {
				ShipmentStatus     string `json:"shipment_status"`
				ActionTime         string `json:"action_time"`
				PositionLat        string `json:"position_lat"`
				PositionLon        string `json:"position_lon"`
				PositionTime       string `json:"position_time"`
				DistributionCenter string `json:"distribution_center"`
				Status             string `json:"status"`
				UID                string `json:"uid"`
				Created            string `json:"created"`
				Changed            string `json:"changed"`
			} `json:"40"`
			Num59 struct {
				ShipmentStatus     string `json:"shipment_status"`
				ActionTime         string `json:"action_time"`
				PositionLat        string `json:"position_lat"`
				PositionLon        string `json:"position_lon"`
				PositionTime       string `json:"position_time"`
				DistributionCenter string `json:"distribution_center"`
				Status             string `json:"status"`
				UID                string `json:"uid"`
				Created            string `json:"created"`
				Changed            string `json:"changed"`
			} `json:"59"`
		} `json:"shipment_update"`
		ShipmentImage []interface{} `json:"shipment_image"`
	} `json:"data"`
}
type DriverLoadboard2 struct {
	Data struct {
		Booking []struct {
			Address           string `json:"address"`
			BookingID         string `json:"booking_id"`
			BookingScore      string `json:"booking_score"`
			BookingStatus     string `json:"booking_status"`
			Company           string `json:"company"`
			Contact           string `json:"contact"`
			DriverName        string `json:"driver_name"`
			Email             string `json:"email"`
			Oid               string `json:"oid"`
			Phone             string `json:"phone"`
			PositionLatitude  string `json:"position_latitude"`
			PositionLongitude string `json:"position_longitude"`
			PositionTime      string `json:"position_time"`
			Price             string `json:"price"`
			PriceUnit         string `json:"price_unit"`
			Rating            string `json:"rating"`
			Time              string `json:"time"`
			TruckPlate        string `json:"truck_plate"`
		} `json:"booking"`
		Message  string `json:"message"`
		Shipment struct {
			Age           string `json:"age"`
			AutoMatchTime string `json:"auto_match_time"`
			BookingCount  int64  `json:"booking_count"`
			BoxAmount     string `json:"box_amount"`
			Broker        struct {
				Address string `json:"address"`
				Contact string `json:"contact"`
				Oid     string `json:"oid"`
				Phone   string `json:"phone"`
				Title   string `json:"title"`
			} `json:"broker"`
			CancelReason interface{} `json:"cancel_reason"`
			Carrier      struct {
				Address string `json:"address"`
				Contact string `json:"contact"`
				Oid     string `json:"oid"`
				Phone   string `json:"phone"`
				Title   string `json:"title"`
			} `json:"carrier"`
			CustomerRating         string `json:"customer_rating"`
			DeliveryDate           string `json:"delivery_date"`
			DeliveryRemark         string `json:"delivery_remark"`
			DestinationAddress     string `json:"destination_address"`
			DestinationDistrict    string `json:"destination_district"`
			DestinationLatitude    string `json:"destination_latitude"`
			DestinationLocation    string `json:"destination_location"`
			DestinationLongitude   string `json:"destination_longitude"`
			DestinationPostal      string `json:"destination_postal"`
			DestinationProvince    string `json:"destination_province"`
			DestinationSubdistrict string `json:"destination_subdistrict"`
			DriverImgURL           string `json:"driver_img_url"`
			DriverName             string `json:"driver_name"`
			JobCount               string `json:"job_count"`
			JobData                []struct {
				BoxAmount     string `json:"box_amount"`
				CAddress      string `json:"c_address"`
				CContact      string `json:"c_contact"`
				CDistrict     string `json:"c_district"`
				CEmail        string `json:"c_email"`
				CName         string `json:"c_name"`
				CProvince     string `json:"c_province"`
				CRefCode      string `json:"c_ref_code"`
				CSubdistrict  string `json:"c_subdistrict"`
				CTel          string `json:"c_tel"`
				CZipcode      string `json:"c_zipcode"`
				CodType       string `json:"cod_type"`
				CompanyCode   string `json:"company_code"`
				ConNo         string `json:"con_no"`
				DeliveryDate  string `json:"delivery_date"`
				Endpoint      string `json:"endpoint"`
				EndpointOid   string `json:"endpoint_oid"`
				EndpointUID   string `json:"endpoint_uid"`
				GroupRefID    string `json:"group_ref_id"`
				NodeID        string `json:"node_id"`
				PickupDate    string `json:"pickup_date"`
				ProductDetail []struct {
					Code        string `json:"code"`
					GrossWeight string `json:"gross_weight"`
					Height      string `json:"height"`
					Length      string `json:"length"`
					Name        string `json:"name"`
					Qty         string `json:"qty"`
					RefID       string `json:"ref_id"`
					Sku         string `json:"sku"`
					TotalCbm    string `json:"total_cbm"`
					TotalWeight string `json:"total_weight"`
					Unit        string `json:"unit"`
					Weight      string `json:"weight"`
					Width       string `json:"width"`
				} `json:"product_detail"`
				RAddress         string `json:"r_address"`
				RContact         string `json:"r_contact"`
				RDistrict        string `json:"r_district"`
				REmail           string `json:"r_email"`
				RLat             string `json:"r_lat"`
				RLon             string `json:"r_lon"`
				RName            string `json:"r_name"`
				RProvince        string `json:"r_province"`
				RRefCode         string `json:"r_ref_code"`
				RRefLink         string `json:"r_ref_link"`
				RSubdistrict     string `json:"r_subdistrict"`
				RTel             string `json:"r_tel"`
				RZipcode         string `json:"r_zipcode"`
				RefOrderNo       string `json:"ref_order_no"`
				SAddress         string `json:"s_address"`
				SContact         string `json:"s_contact"`
				SDistrict        string `json:"s_district"`
				SEmail           string `json:"s_email"`
				SLat             string `json:"s_lat"`
				SLon             string `json:"s_lon"`
				SName            string `json:"s_name"`
				SProvince        string `json:"s_province"`
				SRefCode         string `json:"s_ref_code"`
				SRefLink         string `json:"s_ref_link"`
				SSubdistrict     string `json:"s_subdistrict"`
				STel             string `json:"s_tel"`
				SZipcode         string `json:"s_zipcode"`
				ServiceCode      string `json:"service_code"`
				TotalCbm         int64  `json:"total_cbm"`
				TotalWeight      int64  `json:"total_weight"`
				TransportCompany string `json:"transport_company"`
				TransportPrice   string `json:"transport_price"`
			} `json:"job_data"`
			MatchMode         string `json:"match_mode"`
			MatchPriority2    string `json:"match_priority_2"`
			MatchPriority3    string `json:"match_priority_3"`
			OriginAddress     string `json:"origin_address"`
			OriginDistrict    string `json:"origin_district"`
			OriginLatitude    string `json:"origin_latitude"`
			OriginLocation    string `json:"origin_location"`
			OriginLongitude   string `json:"origin_longitude"`
			OriginPostal      string `json:"origin_postal"`
			OriginProvince    string `json:"origin_province"`
			OriginSubdistrict string `json:"origin_subdistrict"`
			Owner             string `json:"owner"`
			PaymentType       string `json:"payment_type"`
			PickupDate        string `json:"pickup_date"`
			PriceAsking       string `json:"price_asking"`
			PriceExpress      string `json:"price_express"`
			PriceFestival     string `json:"price_festival"`
			PriceOffer        string `json:"price_offer"`
			PriceRateDistance string `json:"price_rate_distance"`
			PriceRateWeight   string `json:"price_rate_weight"`
			PriceStandard     string `json:"price_standard"`
			PriceUnit         string `json:"price_unit"`
			RecipientName     string `json:"recipient_name"`
			Remark            string `json:"remark"`
			ShipmentNo        string `json:"shipment_no"`
			ShipmentType      string `json:"shipment_type"`
			SignatureURL      string `json:"signature_url"`
			Status            string `json:"status"`
			StatusID          string `json:"status_id"`
			Trip              string `json:"trip"`
			TruckPlate        string `json:"truck_plate"`
			TruckType         string `json:"truck_type"`
			TruckTypeTitle    string `json:"truck_type_title"`
			Weight            int64  `json:"weight"`
		} `json:"shipment"`
		ShipmentImage  []interface{} `json:"shipment_image"`
		ShipmentUpdate struct {
			Four0 struct {
				ActionTime         string `json:"action_time"`
				Changed            string `json:"changed"`
				Created            string `json:"created"`
				DistributionCenter string `json:"distribution_center"`
				PositionLat        string `json:"position_lat"`
				PositionLon        string `json:"position_lon"`
				PositionTime       string `json:"position_time"`
				ShipmentStatus     string `json:"shipment_status"`
				Status             string `json:"status"`
				UID                string `json:"uid"`
			} `json:"40"`
			Five9 struct {
				ActionTime         string `json:"action_time"`
				Changed            string `json:"changed"`
				Created            string `json:"created"`
				DistributionCenter string `json:"distribution_center"`
				PositionLat        string `json:"position_lat"`
				PositionLon        string `json:"position_lon"`
				PositionTime       string `json:"position_time"`
				ShipmentStatus     string `json:"shipment_status"`
				Status             string `json:"status"`
				UID                string `json:"uid"`
			} `json:"59"`
		} `json:"shipment_update"`
		Status  string `json:"status"`
		Success bool   `json:"success"`
	} `json:"data"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type (
	StringInterfaceMapMobile map[string]interface{}
	EventMobile              struct {
		MerchantID     string `json:"MerchantID"`
		MerchantName   string `json:"MerchantName"`
		Email          string `json:"Email"`
		ContactName    string `json:"ContactName"`
		PhoneNumber    string `json:"PhoneNumber"`
		MemberDiscount string `json:"MemberDiscount"`
	}
)
type AccountBank struct {
	TrackingID             string
	AccountID              string
	AccountName            string
	AccType                string
	PhoneNo                string
	Email                  string
	BankName               string
	BankAccNo              string
	BankAccName            string
	BankType               string
	PromptPayID            string
	UserType               string
	PrimaryTransferAccount string
	AmountPrice            string
	ComCode                string
	TaxID                  string
	BankCode               string
	Corporation            string
	ReceiptAddress         string
	ReceiptSubDistrict     string
	ReceiptDistrict        string
	ReceiptProvince        string
	ReceiptZipcode         string
}
type ReOrder struct {
	JobID                     string
	JobName                   string
	MerchantID                string
	ProjectName               string
	MerchantWHID              string
	PostCodeWH                string
	PostCodeTHPD              string
	JobDmsID                  string
	JobDMSNet                 string
	JobVolumeWeightNet        string
	JobQty                    string
	JobWeightID               string
	JobWeightNet              string
	JobGrossWeightNet         string
	JobProductType            string
	JobType                   string
	JobStatus                 string
	JobStatusName             string
	JobDriverID               string
	JobDriverName             string
	JobTruckID                string
	JobTruckPlate             string
	JobOpenDT                 string
	JobCallReceiveDT          string
	JobPickupDT               string
	JobAssignDT               string
	JobAppointmentDT          string
	JobBranchReceiveDT        string
	JobBranchAssignDT         string
	JobPendingDT              string
	JobSuspendDT              string
	JobCancelDT               string
	JobCompleteDT             string
	JobRateNet                string
	LTLDiscount               string
	ContactName               string
	ContactPhone              string
	Remark                    string
	IsActual                  string
	IsRequestPrint            string
	MsgToDriver               string
	IsSurvey                  string
	TruckTypeID               string
	FTLStartPrice             string
	FTLAddServices            string
	FTLAddSeason              string
	FTLStandardPrice          string
	FTLSpecialPrice           string
	FTLExpectPrice            string
	FTLFinalPrice             string
	FTLRateOilKM              string
	FTLDistanceKM             string
	FTLOriginGPS              string
	FTLDestinationGPS         string
	FTLAddOnName              string
	FTLDiscount               string
	FTLBidType                string
	FTLInsuranceFee           string
	FTLInsurancePrice         string
	DropPointCount            string
	TMSSendStatus             string
	TMSSendDT                 string
	JobCustomerRead           string
	JobOpenRead               string
	JobAssignRead             string
	IsBookmark                string
	JobUpdateDT               string
	JobCreateDT               string
	Customer_Po_1             string
	Customer_Po               string
	Order_No                  string
	TrackingID                string
	UserID                    string
	WarehouseID               string
	Customer_No               string
	Customer_Name             string
	Receive_name              string
	Receive_Address           string
	Receive_Tel               string
	Receive_Country           string
	Receive_City              string
	Receive_State             string
	Receive_Postal_Code       string
	Receive_GPS               string
	Item                      string
	Item_Description          string
	Product_Category          string
	Original_Product_Category string
	Qty                       string
	UnitType                  string
	Weight                    string
	Width                     string
	Length                    string
	Height                    string
	Ref_1                     string
	Ref_2                     string
	Ref_3                     string
	MessageToSeller           string
	THPDRatePrice             string
	Saved                     string
	PickupStatus              string
	PickupTruckID             string
	PickupTruckPlate          string
	PickupDriverID            string
	PickupDriverName          string
	PickupDT                  string
	IsCOD                     string
	AmountCOD                 string
	IsInsurance               string
	AmountPrice               string
	ItemNo                    string
	DropPointNo               string
	ItemStatus                string
	PreDeliveryID             string
	PreDeliveryDT             string
	RemarkID                  string
	RemarkDesc                string
	IMG1                      string
}
type AmountCOD struct {
	TrackingID string
	AmountCOD  string
}
type Appversion struct {
	Appid          string
	AppType        string
	Version        string
	LastestDT      string
	LocationGPS    string
	Locationdetail string
}

type AppCheckPointversion struct {
	Appid        string
	AppType      string
	Version      string
	LastestDT    string
	TkenExpiredt string
	CustomDetail string
	TranStatus   string
	ShipStartDT  string
	NetWeight    string
	Notifier     string
}
type AppGetIMG struct {
	TransportSPID string
	imgBase64     string
	uploadDT      string
	LastestDT     string
	TkenExpiredt  string
	CustomDetail  string
	TranStatus    string
	ShipStartDT   string
	NetWeight     string
}

type AppGetCalc struct {
	calcID          string
	rubberType      string
	containerType   string
	truckWeight     string
	trailWeight     string
	tankweight      string
	containerWeight string
	boxWeight       string
	totalcontainer  string
	totalbox        string
	containerTypeID string
}
type AppCheckPointTranDetailversion struct {
	TransportSPID            string
	TransportID              string
	CustomName               string
	CusAddress               string
	CusPhone                 string
	ContainerType            string
	ContainerNo              string
	ContainerSideNo          string
	TruckTypeID              string
	TruckLicenseNo           string
	TruckLicenseCountry      string
	TruckLicenseTrail        string
	TruckLicenseTrailCountry string
	ToCustomDetail           string
	CustomID                 string
	TransferDT               string
	NetWeight                string
	GrossWeight              string
	ShippingName             string
	TotalBox                 string
	RubberType               string
	GrossWeightOnSite        string
	EstWeightOnSite          string
	Calctype                 string
	TimeOnSite               string
	Transubstatus            string
	WeightID                 string
}

// type AppToken struct {
// 	Appid        string
// 	AppType      string
// 	Token        string
// 	LastestDT    string
// 	TkenExpiredt string
// }

type AppToken struct {
	Appid            string
	AppType          string
	Token            string
	LastestDT        string
	TkenExpiredt     string
	Name             string
	Surename         string
	Department       string
	Locationresponse string
	Permission       string
	Position         string
	KeycloakToken    string
	KeyclockRole     string
	User             string
	Role             string
	Discrepancy      string
	APIKey           string
	Locationdetail   string
}

type AppReturnUserPositionRAOT struct {
	DBID                     string
	Positionname             string
	Isdashboard              string
	Ismanageweight           string
	Ismanagecheckpoint       string
	Ischeckpoint             string
	Ischeckpointlist         string
	Iscumtomerdata           string
	Istracking               string
	Isreportlicense          string
	Isreporttruckchecked     string
	Isreportweightchecked    string
	Ismanageemployee         string
	Ismanageposition         string
	Ismanagelocationresponse string
	Ismanagegroup            string
	Iscondition              string
	Ismanagecalculate        string
	Isreporttrader           string

	Iditbyuser string
}
type AppReturnResponseLocationRAOT struct {
	DBID                 string
	LocationResponsename string
}

type AppSetUserPositionRAOT struct {
	DBID                     string
	Positionname             string
	Isdashboard              string
	Ismanageweight           string
	Ismanagecheckpoint       string
	Ischeckpoint             string
	Ischeckpointlist         string
	Iscumtomerdata           string
	Istracking               string
	Isreportlicense          string
	Isreporttruckchecked     string
	Isreportweightchecked    string
	Ismanageemployee         string
	Ismanageposition         string
	Ismanagelocationresponse string
	Ismanagegroup            string
	Ismanagecalculate        string
	Isreporttrader           string
	Iscondition              string
	Iditbyuser               string
	UserID                   string
	Token                    string
	EType                    string
	Authorization            string
	Channel                  string
}

type AppReturnTruckEntry struct {
	Configdetail string
	Ref1         string
	Ref2         string
	Ref3         string
	Ref4         string
}

type AppReturnUserRAOT struct {
	DBID             string
	Username         string
	Name             string
	Surename         string
	Position         string
	Department       string
	Locationresponse string
	Grouporganize    string
	Email            string
	Permission       string
	Lastlogindt      string
	Logincounter     string
	KeycloakUserRole string
}
type Widget struct {
	IMG string
}
type Dashboard struct {
	DMY                string
	Cntall             string
	Cntsuccess         string
	Cntincheckpoint    string
	Cntnotincheckpoint string
}

type TruckNoti struct {
	MSGID       string
	MSGTxt      string
	MSGStatus   string
	MSGRef1     string
	MSGDT       string
	TransportID string

	TransportSPID string
	CustomName    string
}
type TruckTrader struct {
	ID                  string
	License_number_th   string
	License_province_th string
	License_number_my   string
	Registration_type   string
	Ttruck_type         string
	Weight_kg           string
	Gps_box_id          string
	Gps_provider_name   string
	Company             string
	Status              string
}

type AppReturnRepTrader struct {
	CustomName string
}

type AppReturnRepWeightChecked struct {
	TransferDT                      string
	RubberType                      string
	CustomName                      string
	CessID                          string
	TransportID                     string
	TransportSPID                   string
	LicenseMY                       string
	TruckLicenseNo                  string
	TruckLicenseCountry             string
	TruckLPRlicenseOverwrite        string
	TruckLPRlicenseCountryOverwrite string
	TruckTypeID                     string
	Checkpointname                  string
	Checkpointcode                  string
	Name                            string
	Locationresponse                string
	Grouporganize                   string
	Transubstatus                   string
	ArrivalDT                       string
	NetWeight                       string
	EstWeightOnSite                 string
	Netallweight                    string
	SumEstWeightOnSite              string
}

type AppReturnRubberType struct {
	RubberType string
}

type AppReturnRepTraderTruck struct {
	TransportSPID            string
	TransportID              string
	CustomName               string
	CusCitizenID             string
	CusAddress               string
	CusPhone                 string
	ContainerType            string
	ContainerNo              string
	ContainerSideNo          string
	TruckTypeID              string
	TruckLicenseNo           string
	TruckLicenseCountry      string
	TruckLicenseTrail        string
	TruckLicenseTrailCountry string
	ToCustomDetail           string
	CustomID                 string
	TransferDT               string
	NetWeight                string
	GrossWeight              string
	ShippingName             string
	TotalBox                 string
	RubberType               string

	CreateDT          string
	ModifyDT          string
	Transubstatus     string
	GrossWeightOnSite string
	EstWeightOnSite   string

	Calctype                        string
	TimeOnSite                      string
	WeightID                        string
	UserID                          string
	TruckLPRlicenseOverwrite        string
	TruckLPRlicenseCountryOverwrite string
	TruckLPRTypeOverwrite           string
	TruckLPRTypeTruckOverwrite      string
	GPSProvider                     string
	GPSBoxID                        string
	TruckLicenseNoMY                string
}
type AppReturnTrader struct {
	DBID             string
	Trader_tax_id    string
	Trader_branch_no string
	Trader_name      string
	Trader_name_en   string
	Trader_type      string
	Address          string
	Subdistrict      string
	District         string
	Province         string
	Postal_code      string
	Phone_number     string
	Fax_number       string
	Email            string
}

type AppReturnCheckpointStart struct {
	CPID             string
	CheckpointCode   string
	CheckpointName   string
	StartDate        string
	StopDate         string
	WeightId         string
	LocationResponse string
	CreateDate       string
	JobID            string
	HaveJob          string
}
type AppReturnConsent struct {
	MsgID      string
	ConsentMsg string
}
type AppReturnJobID struct {
	JobID    string
	WeightId string
}

type AppWeightConfig struct {
	Id           int `json:"id"`
	AppType      string
	Token        string
	LastestDT    string
	TkenExpiredt string
	LocationGPS  string
}

type AppConfigTrader struct {
	Id               int `json:"id"`
	DBID             string
	Trader_tax_id    string
	Trader_branch_no string
	Trader_name      string
	Trader_name_en   string
	Trader_type      string
	Address          string
	Subdistrict      string
	District         string
	Province         string
	Postal_code      string
	Phone_number     string
	Fax_number       string
	Email            string
	Trader_update_dt string
	Trader_create_dt string
	Update_dt        string
	Create_dt        string

	Token        string
	LastestDT    string
	TkenExpiredt string
}

type AppGetJobDetail struct {
	JobID       string
	CessID      string
	TransportID string
}
type AppConfigCheckpointStart struct {
	Id               int `json:"id"`
	CPID             string
	CheckpointCode   string
	CheckpointName   string
	StartDate        string
	StopDate         string
	WeightId         string
	LocationResponse string
	JobID            string
	Token            string
	LastestDT        string
	TkenExpiredt     string
	HaveJob          string
}

type AppConfigCheckpointCalculate struct {
	Id              int `json:"id"`
	CalcID          string
	RubberType      string
	ContainerType   string
	TruckWeight     string
	TrailWeight     string
	TrailCount      string
	ContainerWeight string
	ContainerCount  string
	BoxWeight       string
	BoxCount        string
	Token           string
	LastestDT       string
	TkenExpiredt    string
}

type AutoGPSGenerated struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Error     string `json:"error"`
	Data      struct {
		VehicleID         string `json:"vehicle_id"`
		GpsID             string `json:"gps_id"`
		GpsUnitID         string `json:"gps_unit_id"`
		LicenseNum        string `json:"license_num"`
		LicenseProvince   string `json:"license_province"`
		VehicleType       int    `json:"vehicle_type"`
		AffiliatedUnit    string `json:"affiliated_unit"`
		ResponsiblePerson string `json:"responsible_person"`
		Remark            string `json:"remark"`
		Gps               struct {
			GpsID      string    `json:"gps_id"`
			Imei       string    `json:"imei"`
			Lat        string    `json:"lat"`
			Lng        string    `json:"lng"`
			Speed      string    `json:"speed"`
			Battery1   int       `json:"battery1"`
			Battery2   int       `json:"battery2"`
			Acc        int       `json:"acc"`
			Vbat       int       `json:"vbat"`
			Conn       int       `json:"conn"`
			Sleep      int       `json:"sleep"`
			Temp       int       `json:"temp"`
			Trigger    string    `json:"trigger"`
			CreateDate time.Time `json:"create_date"`
			ModifyDate time.Time `json:"modify_date"`
		} `json:"gps"`
	} `json:"data"`
}

type AppCheckPointConfig struct {
	Id           int `json:"id"`
	AppType      string
	Token        string
	LastestDT    string
	TkenExpiredt string
	CustomDetail string
	TranStatus   string
	StatusCD     string
	ShipStartDT  string
	NetWeight    string
	Notifier     string
}

type AppLicense struct {
	Id      int `json:"id"`
	License string
	LAT     string
	LNG     string
}
type AppReturnIMG struct {
	Id            int `json:"id"`
	TransportSPID string
	IMGBase64     string
	UploadDT      string
}

type AppReturnCalc struct {
	Id              int `json:"id"`
	CalcID          string
	RubberType      string
	ContainerType   string
	TruckWeight     string
	TrailWeight     string
	Tankweight      string
	ContainerWeight string
	BoxWeight       string
	Totalcontainer  string
	Totalbox        string
	ContainerTypeID string
}
type AppCheckPointTranDetailConfig struct {
	Id                       int `json:"id"`
	TransportSPID            string
	TransportID              string
	CustomName               string
	CusAddress               string
	CusPhone                 string
	ContainerType            string
	ContainerNo              string
	ContainerSideNo          string
	TruckTypeID              string
	TruckLicenseNo           string
	TruckLicenseCountry      string
	TruckLicenseTrail        string
	TruckLicenseTrailCountry string
	ToCustomDetail           string
	CustomID                 string
	TransferDT               string
	NetWeight                string
	GrossWeight              string
	ShippingName             string
	TotalBox                 string
	RubberType               string
	GrossWeightOnSite        string
	EstWeightOnSite          string
	Calctype                 string
	TimeOnSite               string
	Transubstatus            string
	WeightID                 string
	Statusflag               string
	GPSLAT                   string
	GPSLNG                   string
	GPSLastUpdateDT          string
}
type Mobile struct {
	MobileID           string
	Password           string
	TypeEvent          string
	Lang               string
	BG                 string
	LoginType          string
	AccName            string
	AccBankName        string
	AccNo              string
	AccPromptPay       string
	AccTransferPrimary string
	AccType            string
}
type TruckType struct {
	TruckID          string
	TruckName        string
	TruckDesc        string
	TruckSize        string
	RateServiceStart string
	RateOilPerKM     string
	IsUse            string
}
type MerchantRep struct {
	MerchantID string
	RepType    string
	SDT        string
	EDT        string
}
type Merchant struct {
	MerchantID       string
	Address1         string
	Address2         string
	LocationGPS      string
	WarehouseID      string
	WarehouseName    string
	SubDistrict      string
	District         string
	ProvinceName     string
	PostCode         string
	ContactName      string
	PhoneNumber      string
	EmailContact     string
	StartWorkingHour string
	EndWorkingHour   string
}
type Coupon struct {
	CouponName     string
	CouponCode     string
	CouponDiscount string
	CouponExpireDT string
}

type PickupTime struct {
	TimeType string
	TimeName string
}
type RateAddOn struct {
	DBID       int
	AddonName  string
	AddonDesc  string
	AddonPrice string
	IsUse      int
	ModifyDt   string
	CreateDT   string
}
type MobileOTP struct {
	MobileID    string
	OTPPassword string
}
type PaymentQRCode struct {
	QRCodetxt string
	ExpDT     string
}
type RepPayment struct {
	PMerchantID string
	Ftypes      string
	SumAmount   string
}
type RepSummary struct {
	UserID       string
	Status_Name  string
	CStatus_Name string
}
type ChecksendPrice struct {
	Jobid         string
	TrackingID    string
	THPDRatePrice string
	AmountCOD     string
	FTLFinalPrice string
}
type ChecksPrice struct {
	Jobid      string
	TrackingID string
	Price      string
	PayDT      string
}
type UrlAdvise struct {
	MsgAlert string
	StartDT  string
	EndDT    string
}
type MobileID struct {
	MobileID        string
	SendSMSFirstJob int
	TrackingID      string
}
type Checksendprice struct {
	Jobid         string
	TrackingID    string
	THPDRatePrice string
	AmountCOD     string
	FTLFinalPrice string
	BankRef       string
}
type AuthDevice struct {
	Cid              int
	ChannelName      string
	ChannelSecretKey string
}
type DriverTracking struct {
	TrackingID string
	DriverID   string
	DriverName string
	Plate      string
	FinalPrice string
}
type DriverJobMasterBooking struct {
	TrackingID    string
	JobDriverID   string
	JobDriverName string
	JobTruckID    string
	JobTruckPlate string
}
type Get2c2ppaymentToken struct {
	TrackingID          string
	Amount              string
	PhoneNumber         string
	PaymentType         string
	ComCode             string
	TransportationPrice string
	ProductPrice        string
	NumberOfPieces      string
	PaytoMerchant       string
}
type GetPostPrintOnline struct {
	PostCode   string
	CountItems int
}
type GetCancel struct {
	Shipment_no string
	Remark      string
}
type GetAccount struct {
	Shipment_no string
	Remark      string
}

type TblTMSItemEvents struct {
	DBID        string `json:"DBID"`
	TrackingID  string `json:"TrackingID"`
	Status_Name string `json:"Status_Name"`
	CreateDT    string `json:"CreateDT"`
	PaymentFlag string `json:"PaymentFlag"`
	OfferPrice  string `json:"OfferPrice"`
}

//	type response1 struct {
//		MobileID       string
//		OMSOrderStruct []string
//	}
type OMSOrderStruct struct {
	JobID           string `json:"JobID"`
	MobileID        string `json:"MobileID"`
	Token           string `json:"token"`
	ReceiveName     string `json:"ReceiveName"`
	ReceiveAddress  string `json:"ReceiveAddress"`
	ReceiveTumbon   string `json:"ReceiveTumbon"`
	ReceiveDistrict string `json:"ReceiveDistrict"`
	ReceiveProvince string `json:"ReceiveProvince"`
	ReceivePhoneNo  string `json:"ReceivePhoneNo"`
	ReceiveZipcode  string `json:"ReceiveZipcode"`
	SenderZipcode   string `json:"SenderZipcode"`
	SendPrice       string `json:"SendPrice"`
	PickupStartDt   string `json:"PickupStartDt"`
	DeliveryEndDt   string `json:"DeliveryEndDt"`
	PaymentFlag     string `json:"PaymentFlag"`
	PaymentDetail   string `json:"PaymentDetail"`
	IMG1            string `json:"IMG1"`
	IMG2            string `json:"IMG2"`
	IMG3            string `json:"IMG3"`
	IMG4            string `json:"IMG4"`
	CreateDt        string `json:"Date"`
	JobDesc         string `json:"JobDesc"`
	JobType         string `json:"JobType"`
	MerchantID      string `json:"MerchantID"`
	WarehouseID     string `json:"WarehouseID"`
}
type (
	StringInterfaceMap0 map[string]interface{}
	Event0              struct {
		InvoiceNo       string `json:"InvoiceNo"`
		TrackingID      string `json:"TrackingID"`
		Truck_type      string `json:"Truck_type"`
		Truck_plate     string `json:"Truck_plate"`
		Driver_name     string `json:"FullName"`
		Driver_Phone    string `json:"DriverTel"`
		Total_chest     string `json:"Total_chest"`
		Total_pack      string `json:"Total_pack"`
		CustomerNo      string `json:"CustomerNo"`
		InvoiceDate     string `json:"InvoiceDate"`
		To_created_date string `json:"To_created_date"`
	}
)
type (
	StringInterfaceMap5noIMG map[string]interface{}
	Event5noIMG              struct {
		JobID           string `json:"JobID"`
		MobileID        string `json:"MobileID"`
		ReceiveName     string `json:"ReceiveName"`
		ReceiveAddress  string `json:"ReceiveAddress"`
		ReceiveTumbon   string `json:"ReceiveTumbon"`
		ReceiveDistrict string `json:"ReceiveDistrict"`
		ReceiveProvince string `json:"ReceiveProvince"`
		ReceiveZipcode  string `json:"ReceiveZipcode"`
		ReceivePhoneNo  string `json:"ReceivePhoneNo"`
		SenderZipcode   string `json:"SenderZipcode"`
		SendPrice       string `json:"SendPrice"`
		PickupStartDt   string `json:"PickupStartDt"`
		DeliveryEndDt   string `json:"DeliveryEndDt"`
		PaymentFlag     string `json:"PaymentFlag"`
		PaymentDetail   string `json:"PaymentDetail"`
		CreateDt        string `json:"CreateDt"`
		JobDesc         string `json:"JobDesc"`
		JobType         string `json:"JobType"`
		WarehouseID     string `json:"WarehouseID"`
		MerchantID      string `json:"MerchantID"`
		Status          string `json:"Status"`
		TrackingID      string `json:"TrackingID"`
		IMG1            string `json:"IMG1"`
		IMG2            string `json:"IMG2"`
		IMG3            string `json:"IMG3"`
		IMG4            string `json:"IMG4"`
	}
)
type (
	StringInterfaceMapPostPrint map[string]interface{}
	EventPostPrint              struct {
		THPDID         string `json:"THPDID"`
		LblContainerID string `json:"lblContainerID"`
		THPDName       string `json:"THPDName"`
	}
)

type (
	StringInterfaceMapDriverTrack map[string]interface{}
	EventDriverTrack              struct {
		JobID           string `json:"JobID"`
		TrackingID      string `json:"TrackingID"`
		CarID           string `json:"CarID"`
		DriverName      string `json:"DriverName"`
		DriverPhoneNo   string `json:"DriverPhoneNo"`
		Company         string `json:"Company"`
		Price           string `json:"Price"`
		CustomerSelect  string `json:"CustomerSelect"`
		CustomConfirmDT string `json:"CustomConfirmDT"`
		GetBeforeDT     string `json:"GetBeforeDT"`
	}
)

// respok["respDesc"] = result.respDesc
// respok["amount"] = fmt.Sprintf("%f", result.amount) //QR
// respok["refCode"] = result.refCode                  //QR
// respok["invoiceNo"] = result.invoiceNo
// respok["agentCode"] = result.agentCode
// respok["cardType"] = result.cardType
// respok["transactionDatetime"] = result.transactionDatetime
type (
	StringInterfaceMapDriverTrackPayment map[string]interface{}
	EventPayment                         struct {
		respDesc            string `json:"respDesc"`
		amount              string `json:"amount"`
		refCode             string `json:"refCode"`
		invoiceNo           string `json:"invoiceNo"`
		agentCode           string `json:"agentCode"`
		cardType            string `json:"cardType"`
		cardNo              string `json:"cardNo"`
		transactionDatetime string `json:"transactionDatetime"`
		phoneNo             string `json:"phoneNo"`
		refID               string `json:"refID"`
	}
)

type (
	StringInterfaceMap55 map[string]interface{}
	Event55              struct {
		JobID       string `json:"JobID"`
		TrackingID  string `json:"TrackingID"`
		Customer_Po string `json:"Customer_Po"`
	}
)
type (
	StringInterfaceMap5 map[string]interface{}
	Event5              struct {
		JobID           string `json:"JobID"`
		MobileID        string `json:"MobileID"`
		ReceiveName     string `json:"ReceiveName"`
		ReceiveAddress  string `json:"ReceiveAddress"`
		ReceiveTumbon   string `json:"ReceiveTumbon"`
		ReceiveDistrict string `json:"ReceiveDistrict"`
		ReceiveProvince string `json:"ReceiveProvince"`
		ReceiveZipcode  string `json:"ReceiveZipcode"`
		ReceivePhoneNo  string `json:"ReceivePhoneNo"`
		SenderZipcode   string `json:"SenderZipcode"`
		SendPrice       string `json:"SendPrice"`
		PickupStartDt   string `json:"PickupStartDt"`
		DeliveryEndDt   string `json:"DeliveryEndDt"`
		PaymentFlag     string `json:"PaymentFlag"`
		PaymentDetail   string `json:"PaymentDetail"`
		IMG1            string `json:"IMG1"`
		IMG2            string `json:"IMG2"`
		IMG3            string `json:"IMG3"`
		IMG4            string `json:"IMG4"`
		CreateDt        string `json:"CreateDt"`
		JobDesc         string `json:"JobDesc"`
		JobType         string `json:"JobType"`
		WarehouseID     string `json:"WarehouseID"`
		MerchantID      string `json:"MerchantID"`
		Status          string `json:"Status"`
	}
)
type (
	StringInterfaceMap1 map[string]interface{}
	Event1              struct {
		MerchantName string `json:"MerchantName"`
		Email        string `json:"Email"`
	}
)

type Payload struct {
	Stuff Data
}
type Data struct {
	Job []Event5
}

var (
	// Options
	//flagAlg     = flag.String("alg", "HS256", algHelp())
	flagKey     = flag.String("key", "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", "path to key file or '-' to read from stdin")
	flagCompact = flag.Bool("compact", false, "output compact JSON")
	flagDebug   = flag.Bool("debug", false, "print out all kinds of debug data")
	//flagClaims  = make(ArgList)
	//flagHead    = make(ArgList)

	// Modes - exactly one of these is required
	flagSign   = flag.String("sign", "", "path to claims object to sign, '-' to read from stdin, or '+' to use only -claim args")
	flagVerify = flag.String("verify", "", "path to JWT token to verify or '-' to read from stdin")
	flagShow   = flag.String("show", "", "path to JWT file or '-' to read from stdin")
)

type Payload2 struct {
	Payload  string `json:"payload"`
	RespDesc string `json:"respDesc"`
	RespCode string `json:"respCode"`
}
type PaymentQRStruct struct {
	Type              string `json:"type"`
	ExpiryTimer       string `json:"expiryTimer"`
	ExpiryDescription string `json:"expiryDescription"`
	Data              string `json:"data"`
	ChannelCode       string `json:"channelCode"`
	RespCode          string `json:"respCode"`
	RespDesc          string `json:"respDesc"`
}
type User struct {
	paymentToken      string
	ExpDT             string
	channelCode       string
	respDesc          string
	expiryDescription string
}
type RespCode struct {
	respDesc            string
	invoiceNo           string
	amount              float64
	refCode             string
	agentCode           string
	cardType            string
	transactionDatetime string
	phoneNo             string
	cardNo              string
	refID               int64
}
type User2 struct {
	paymentToken  string
	webPaymentUrl string
	channelCode   string
	respDesc      string
}

var DB *sql.DB

var (
	APIName = ""
	//selectEventByIdQueryTOAT    = `SELECT * FROM THPDTOATDB.VW_TOAT_Track WHERE InvoiceNo = ? and InvoiceNo <> '' and Truck_type <> ''  limit 1`
	//selectEventByCusIdQueryTOAT = `SELECT * FROM THPDTOATDB.VW_TOAT_Track WHERE CustomerNo = ? and InvoiceDate = ? and InvoiceNo <> '' and Truck_type <> ''  limit 1`
	selectEventByToken  = `SELECT Email,MerchantName FROM THPDMPDB.tblMerchant WHERE MerchantID = ? and APIAuthenKey = ? limit 1`
	insertOrderItem     = `INSERT INTO THPDMPDB.tblMobileOMSOrder(JobID,MobileID,ReceiveName,ReceiveAddress,ReceiveTumbon,ReceiveDistrict,ReceiveProvince,ReceiveZipcode,ReceivePhoneNo,SenderZipcode,SendPrice,PickupStartDt,DeliveryEndDt,PaymentFlag,PaymentDetail,CreateDt,JobDesc,JobType,WarehouseID,MerchantID,Status) values (  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),?,?,?,?,? )`
	insertOrderItemDesc = `INSERT INTO THPDMPDB.tblMobileOMSOrderDesc(JobID,MobileID,IMG1,IMG2,IMG3,IMG4,CreateDt) values ( ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )`
	insertGetDriver     = `INSERT INTO THPDMPDB.tblMobileOMSJobDriverBooking(JobID, TrackingID, CarID, DriverName, DriverPhoneNo, Company, Price, CustomerSelect, CustomConfirmDT, GetBeforeDT,CreateDT) values ( ?, ?, ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )`
	insertGetPayment    = `INSERT INTO THPDMPDB.tblPayment( Refid,TrackingID, MerchantID, PhoneNumber, Amount, PaymentType, CreditCardNo, PaymentComplete, PaymentCreateDT, PaymentModifyDT,RefCode) values ( ?, ?, ?, ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),? )`

	selectOrderOMS = `SELECT * FROM THPDMPDB.tblMobileOMSOrder WHERE MobileID = ? `

	//selectEventByIdQueryTOAT    = `CALL THPDTOATDB.SPGetTrackbyInvoiceID (?)`
	//selectEventByCusIdQueryTOAT = `CALL THPDTOATDB.SPGetTrackbyCustomerID (?,?)`
)

type Employee struct {
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Address map[string]string `json:"address,omitempty"`
}
type THSMS struct {
	Message string   `json:"message"`
	Msisdn  []string `json:"msisdn"`
	Sender  string   `json:"sender"`
}
type LineMSG struct {
	ReplyToken string   `json:"replyToken"`
	Messages   string   `json:"messages"`
	Msisdn     []string `json:"msisdn"`
}
type Con_no struct {
	con_no string `json:"con_no"`
}

type company struct {
	replyToken string `json:"replyToken"`
	messages   []message
}

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

var Articles []Article
var OMSOrderStructs []OMSOrderStruct
var Zipcodes []Zipcode
var Mobiles []Mobile
var Merchants []Merchant
var MerchantReps []MerchantRep
var DriverTrackings []DriverTracking
var DriverJobMasterBookings []DriverJobMasterBooking
var Get2c2ppaymentTokens []Get2c2ppaymentToken
var GetCancels []GetCancel
var GetAccounts []AccountBank
var GetKTBApproveJson []KTBApproveJson
var GetMobileVersion []MobileVersion

//var GetPostPrintOnlines []GetPostPrintOnline

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("headers: %v\n", r.Header)

	//	_, err := io.Copy(os.Stdout, r.Body)
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}
	//reqBody, _ := ioutil.ReadAll(r.Body)

	var article LineStruceture // line bot
	json.Unmarshal(reqBody, &article)

	ReplyToken := article.Events[0].ReplyToken
	MsgText := article.Events[0].Message.Text
	EventType := article.Events[0].Type

	// fmt.Printf("ReplyToken: \n", ReplyToken)
	// fmt.Printf("MsgText: \n", MsgText)

	matched, err := regexp.MatchString(`FT([0-9]+)DB`, MsgText)
	bookingmatch := 0
	////fmt.Println(matched)

	if !matched {
		matched, err = regexp.MatchString(`LT([0-9]+)DB`, MsgText)
		//fmt.Println(matched)
		//fmt.Println(err)
	}
	if !matched {
		matched, err = regexp.MatchString(`GD([0-9]+)TH`, MsgText)
		//fmt.Println(matched)
		//fmt.Println(err)
	}
	if !matched {
		matched, err = regexp.MatchString(`2TO([0-9]+)`, MsgText)
		if err == nil {

			bookingmatch = 1

		}
		//fmt.Println(matched)
		//fmt.Println(err)
	}
	if !matched {
		matched, err = regexp.MatchString(`B([0-9]+)`, MsgText)
		if err == nil {

			bookingmatch = 1

		}
		////fmt.Println(bookingmatch)
		////fmt.Println(matched)
		////fmt.Println(err)
	}

	//var bearerToken = "YBAnQHA776xhZVp5IiPGWG556tuo5hU3h0Ke0afe2LBi2pCAymiz0fgLcjLAEeWQZ+9+oJkjH1RFYyA676vOmk78/CCb7Bgns1eSm7GRrGf7GKlwGhp944byiEQZbhV5X1QXpAQMXSM0nN/zusI6yAdB04t89/1O/w1cDnyilFU=" line oongang
	var bearerToken = "OuTmwv2qDq9DJQw4+nIF6nYUv3WQt2dZFeKt3sxMyJY5x+X5wnSliG/vzkH/LSqOb1MIaTNPpIGo/aAiYxiArKY+HzWm8P+oNdO6sL1/16GygzZu6WO84ZMKz1v6xzAhlOJwrL+uVC02wWRixEFBVgdB04t89/1O/w1cDnyilFU="
	var myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"ระบบบริหารแพลตฟอร์มขนส่ง THPD DE ยินดีให้บริการ "},{"type":"text","text":"ท่านสามารถค้นหาสิ่งของโดยพิมพ์หมายเลข TrackID 13 หลัก"}]}`)

	if matched {

		chkdigit := 0
		chkdigit2 := 0

		if bookingmatch == 0 {

			// Check Digit
			chars := []rune(MsgText)
			//strArray := strings.Fields(MsgText)
			//var tx =
			t1x, err := strconv.Atoi(string(chars[2]))
			t2x, err := strconv.Atoi(string(chars[3]))
			t3x, err := strconv.Atoi(string(chars[4]))
			t4x, err := strconv.Atoi(string(chars[5]))
			t5x, err := strconv.Atoi(string(chars[6]))
			t6x, err := strconv.Atoi(string(chars[7]))
			t7x, err := strconv.Atoi(string(chars[8]))
			t8x, err := strconv.Atoi(string(chars[9]))
			t9x, err := strconv.Atoi(string(chars[10]))

			if err != nil {
				panic(err)
			}

			t1 := t1x * 8
			t2 := t2x * 6
			t3 := t3x * 4

			t4 := t4x * 2
			t5 := t5x * 3
			t6 := t6x * 5

			t7 := t7x * 9
			t8 := t8x * 7

			t9 := t9x * 1

			tall := (t1 + t2 + t3 + t4 + t5 + t6 + t7 + t8) % 11

			//chkdigit := 0

			if tall == 0 {
				chkdigit = 5
			} else if tall == 1 {
				chkdigit = 0
			} else {
				chkdigit = 11 - int(tall)
			}
			chkdigit2 = int(t9)

		}

		if int(chkdigit) != chkdigit2 {

			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" หมายเลข TrackID ของท่านไม่ถูกต้อง! "}]}`)

			//return
		} else {

			//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
			//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
			//dns := getDNSString("THPDMPDB", "admin", "!10<>Oms!", "uat-olivedb.mysql.database.azure.com")

			// dns := getDNSString(dblogin, userlogin, passlogin, conn)
			// db, err := sql.Open("mysql", dns)

			// if err != nil {
			// 	panic(err)
			// }
			// err = db.Ping()
			// if err != nil {
			// 	panic(err)
			// }

			if DB.Ping() != nil {
				connectDb()
			}
			if DB.Stats().OpenConnections != 0 {
				//fmt.Println(DB.Stats().OpenConnections)
			} else {
				connectDb()
			}
			db := DB
			defer db.Close()
			// err := db.Ping()
			// if err != nil {
			// 	connectDb()

			// }

			count := 0
			ress2, err := db.Query("SELECT  DBID,TrackingID,Status_Name,CreateDT,IFNULL(PaymentFlag,'N')PaymentFlag,IFNULL(OfferPrice,'0')OfferPrice FROM THPDMPDB.VW_OMSMobileJobPayment  where (trackingid = '" + MsgText + "' or DBID ='" + MsgText + "') and  createdt > CURRENT_TIME()-30 order by createdt desc ")

			if err == nil {

				for ress2.Next() {
					var event TblTMSItemEvents

					//JobID := ress2.Scan(&event.JobID)
					err := ress2.Scan(&event.DBID, &event.TrackingID, &event.Status_Name, &event.CreateDT, &event.PaymentFlag, &event.OfferPrice)

					if err != nil {
						panic(err)
					}
					paymentMassage := "ยังไม่ชำระเงินค่าขนส่ง"
					if event.PaymentFlag == "Success" {
						paymentMassage = "ชำระเงินค่าขนส่งแล้วจำนวน " + event.OfferPrice
					} else {

					}

					dt := strings.Split(event.CreateDT, "T")[0]
					dt = strings.Split(dt, "-")[2] + "-" + strings.Split(dt, "-")[1] + "-" + strings.Split(dt, "-")[0]

					tt := strings.Split(event.CreateDT, "T")[1]
					tt = strings.Replace(tt, "Z", "", -1)

					if event.PaymentFlag != "Success" {

						myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะจัดส่งล่าสุด : ` + event.Status_Name + ` วันที่ ` + dt + ` เวลา ` + tt + ` "},{"type":"text","text":" การชำระเงิน : ` + paymentMassage + `"},{
						"type": "template",
						"altText": "This is a buttons template",
						"template": {
							"type": "buttons",
							"thumbnailImageUrl": "https://pbs.twimg.com/media/Dp9bLnWVsAATqao.jpg",
							"imageAspectRatio": "rectangle",
							"imageSize": "cover",
							"imageBackgroundColor": "#FFFFFF",
							"title": "การชำระเงิน",
							"text": "โปรดเลือกวิธีการชำระเงิน",
							"defaultAction": {
								"type": "uri",
								"label": "เข้าสู่เว็บไซต์พร้อมสั่ง",
								"uri": "http://example.com/page/123"
							},
							"actions": [
								{
									"type": "uri",
									"label": "ชำระผ่าน Credit Card",
									"uri": "https://oms-report.promptsong.co/?Credit=` + MsgText + `"  
									
								},
								{
								  "type": "uri",
								  "label": "ชำระผ่าน QRCode",
								  "uri": "https://oms-report.promptsong.co/?QRCode=` + MsgText + `" 
								}
							]
						}
					  }	 ]}`)
					}
					//https://uat-oms-report.azurewebsites.net/?QRCode
					if event.PaymentFlag == "Success" {

						myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะจัดส่งล่าสุด : ` + event.Status_Name + ` วันที่ ` + dt + ` เวลา ` + tt + ` "},{"type":"text","text":" การชำระเงิน : ` + paymentMassage + `"}	 ]}`)
					}

					count = 1
					break
				}

				defer ress2.Close()
				err = ress2.Close()

				result1 := strings.Index(MsgText, "DB0")

				if result1 == -1 {
					result1 = strings.Index(MsgText, "2TO")
				}

				if result1 >= 0 {

					ress22, err := db.Query("SELECT PaymentID DBID,TrackingID ,'' Status_Name ,CreateDT ,'N' PaymentFlag,FinalPrice  OfferPrice   FROM thpdmpdb.tblpaymentrequest WHERE TrackingID = '" + MsgText + "' and  TrackingID not in (SELECT Ref2 FROM thpdmpdb.tblpaymentktbcharge)  order by CreateDT desc ")
					if err == nil {
						for ress22.Next() {
							var event2 TblTMSItemEvents
							err := ress22.Scan(&event2.DBID, &event2.TrackingID, &event2.Status_Name, &event2.CreateDT, &event2.PaymentFlag, &event2.OfferPrice)

							if err != nil {
								panic(err)
							}
							paymentMassage := "ยังไม่ชำระเงินค่าขนส่ง"
							if event2.PaymentFlag == "Success" {
								paymentMassage = "ชำระเงินค่าขนส่งแล้วจำนวน " + event2.OfferPrice
							} else {

							}

							dt := strings.Split(event2.CreateDT, "T")[0]
							dt = strings.Split(dt, "-")[2] + "-" + strings.Split(dt, "-")[1] + "-" + strings.Split(dt, "-")[0]

							tt := strings.Split(event2.CreateDT, "T")[1]
							tt = strings.Replace(tt, "Z", "", -1)

							if event2.PaymentFlag != "Success" {

								myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะจัดส่งล่าสุด : ` + event2.Status_Name + ` วันที่ ` + dt + ` เวลา ` + tt + ` "},{"type":"text","text":" การชำระเงิน : ` + paymentMassage + `"},{
								"type": "template",
								"altText": "This is a buttons template",
								"template": {
									"type": "buttons",
									"thumbnailImageUrl": "https://pbs.twimg.com/media/Dp9bLnWVsAATqao.jpg",
									"imageAspectRatio": "rectangle",
									"imageSize": "cover",
									"imageBackgroundColor": "#FFFFFF",
									"title": "การชำระเงิน",
									"text": "โปรดเลือกวิธีการชำระเงิน",
									"defaultAction": {
										"type": "uri",
										"label": "เข้าสู่เว็บไซต์พร้อมสั่ง",
										"uri": "http://example.com/page/123"
									},
									"actions": [
										{
											"type": "uri",
											"label": "ชำระผ่าน Credit Card",
											"uri": "https://oms-report.promptsong.co/?Credit=` + MsgText + `"  
											
										},
										{
										  "type": "uri",
										  "label": "ชำระผ่าน QRCode",
										  "uri": "https://oms-report.promptsong.co/?QRCode=` + MsgText + `" 
										}
									]
								}
							  }	 ]}`)
							}

						}
					}

					//myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะล่าสุด : ไม่พบสถานะการจัดส่ง "}]}`)

				}

				if count == 0 && result1 < 0 {
					myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะล่าสุด : ไม่พบสถานะการจัดส่ง "}]}`)

				}

			}
			defer db.Close()

			url := "https://api.line.me/v2/bot/message/reply"
			//req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(myJsonString))
			if err != nil {
				panic(err)
			}

			req.Header.Add("Authorization", "Bearer "+bearerToken)
			req.Header.Add("Content-Type", "application/json")
			//fmt.Println(req.Header)

			//req.Body.Read(json_data)
			//Send req using http Client
			client := &http.Client{}
			resp2, err := client.Do(req)

			if resp2 != nil {
				////fmt.Println(resp2)
				//panic(err)
				myJsonString = []byte(`{"replyToken":"` + ReplyToken + `"}`)
			}
			if err == nil {

			}

		}

		//	myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะนำจ่ายสำเร็จ วันที่ 10/6/2565 12:00:00"}]}`)

	} else {

		if EventType != "postback" {
			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[
					{
						"type": "template",
						"altText": "This is a buttons template",
						"template": {
							"type": "buttons",
							"thumbnailImageUrl": "https://uat-oms-report.azurewebsites.net/images/B2.png",
							"imageAspectRatio": "rectangle",
							"imageSize": "cover",
							"imageBackgroundColor": "#FFFFFF",
							"title": "โปรดตอบแบบสอบถามเพื่อให้เรารู้จักท่านยิ่งดียิ่งขึ้น",
							"text": "กดทำแบบสอบถาม",
							"defaultAction": {
								"type": "uri",
								"label": "เข้าสู่เว็บไซต์พร้อมสั่ง",
								"uri": "https://docs.google.com/forms/d/1c3asJE0eHvlsPPxmqh6nZU9m0LPkR-_GZHYo_ydIgTY/viewform?edit_requested=true"
							},
							"actions": [
								{
									"type": "uri",
									"label": "เข้าสู่เว็บไซต์พร้อมสั่ง",
									"uri": "https://www.promptsung.co"
								},
								{
								  "type": "uri",
								  "label": "เข้าสู่เว็บไซต์พร้อมส่ง",
								  "uri": "https://www.promptsong.co"
								},
								{
								  "type": "postback",
								  "label": "รายการสินค้าที่น่าสนใจ",
								  "data": "action=add&itemid=123"
								}
							]
						}
					  }	
					]}`)
		} else {

			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[
				{
					"type": "template",
					"altText": "this is a carousel template",
					"template": {
						"type": "carousel",
						"columns": [
							{
							  "thumbnailImageUrl": "https://uat-api-marketplace-frontend.azurewebsites.net/user-content/42324ab2db88408490684088d36a1ce7.jpg",
							  "imageBackgroundColor": "#FFFFFF",
							  "title": "มะม่วงเขียวสามรส มะม่วงมัน กรอบ หวาน  ",
							  "text": "ชุดละ 70 บาท (มีแค่ 200 ชุด.เท่านั้น) สั่งเลย",
							  "defaultAction": {
								  "type": "uri",
								  "label": "View detail",
								  "uri": "https://www.promptsung.co/"
							  },"actions": [
								{
									"type": "uri",
									"label": "View detail",
									"uri": "https://www.promptsung.co/product/default/7e81b39c-2030-42bb-b114-08da75edfdce"
								}
							]
							 
							},
							{
								"thumbnailImageUrl": "https://uat-api-marketplace-frontend.azurewebsites.net/user-content/7aff36fd238148779173d9d9fe8034e2.jpg",
								"imageBackgroundColor": "#FFFFFF",
								"title": "ทุเรียนหมอนทองเนื้อล้วน (ราคาต่อขีดจ้า) หอมหวานมัน กรอบนอกนุ่มใน ",
								"text": "โลละ 140 บาท",
								"defaultAction": {
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/111"
								},"actions": [
						
								  {
									  "type": "uri",
									  "label": "View detail",
									  "uri": "http://example.com/page/111"
								  }
							  ]
							   
							  },
							{
							  "thumbnailImageUrl": "https://backend.tops.co.th/media/catalog/product/0/2/0213967000000_mc2021.jpg",
							  "imageBackgroundColor": "#000000",
							  "title": "ลำไยเชียงใหม่ ",
							  "text": "โลละ 50 บาท สั่งเลย!!",
							  "defaultAction": {
								  "type": "uri",
								  "label": "View detail",
								  "uri": "http://example.com/page/222"
							  },
							  "actions": [
					
								{
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/111"
								}
							]
							},
							{
								"thumbnailImageUrl": "https://www.thaihealth.or.th/data/content/2017/08/38295/cms/newscms_thaihealth_c_ijnoquw34567.jpg",
								"imageBackgroundColor": "#000000",
								"title": "มังคุด ",
								"text": "โลละ 30 บาท สั่งเลย!!",
								"defaultAction": {
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/222"
								},
								"actions": [
					
									{
										"type": "uri",
										"label": "View detail",
										"uri": "http://example.com/page/111"
									}
								]
							},
							{
								"thumbnailImageUrl": "https://s.isanook.com/ca/0/ui/182/910427/06190_002.jpg",
								"imageBackgroundColor": "#000000",
								"title": "เงาะโรงเรียน ",
								"text": "โลละ 30 บาท สั่งเลย!!",
								"defaultAction": {
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/222"
								},
								"actions": [
					
									{
										"type": "uri",
										"label": "View detail",
										"uri": "http://example.com/page/111"
									}
								]
							},
							{
								"thumbnailImageUrl": "https://inwfile.com/s-fy/tj0xx9.jpg",
								"imageBackgroundColor": "#000000",
								"title": "ลองกอง ",
								"text": "โลละ 30 บาท สั่งเลย!!",
								"defaultAction": {
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/222"
								},
								"actions": [
					
									{
										"type": "uri",
										"label": "View detail",
										"uri": "http://example.com/page/111"
									}
								]
							}
						],
						"imageAspectRatio": "rectangle",
						"imageSize": "cover"
					}
				  }

			]}`)
		}

		url := "https://api.line.me/v2/bot/message/reply"
		//req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(myJsonString))
		if err != nil {
			panic(err)
		}

		req.Header.Add("Authorization", "Bearer "+bearerToken)
		req.Header.Add("Content-Type", "application/json")
		//fmt.Println(req.Header)

		//req.Body.Read(json_data)
		//Send req using http Client
		client := &http.Client{}
		resp2, err := client.Do(req)

		if resp2 != nil {
			////fmt.Println(resp2)
			//panic(err)
			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `"}`)
		}
		if err == nil {

		}

	}

}

func homePage(w http.ResponseWriter, r *http.Request) {

	t := time.Now()
	//fmt.Println(t.String())
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "WELCOME TO CHECKPOINT-API version 1.0.4 "+userlogin)
	//fmt.Println("Endpoint Hit: homePage"

	// version V 1.0.2
	// j  ecit for sent update driver to DX and stamp update mobileorder
	// k  edit for driver offer price
	// m insert phone to tblpaymentrequest
	// M edit payment LTL api name OMSMobileSelectOrderFromBooking
	// n edit update auto from backend
	// N edit checkprice
	// o update webhook O update vieworder loadboard
	// p update cancel
	// P update payment report
	// q update
	// Q update edit problem auto matching driver in MobileGetDriver
	// r edit first time login insert to tblmobileomsmemberconfig
	// R keep log smsotp use
	// s edit OMSCheckPrice
	// S add OMSGetBankAccount
	// t edit webhook line
	// T update api UAT from Chitman
	// u getCOD
	// U set ThaiQR When call Create QRCode
	// v update uattme to tms
	// V update tms
	// w update 2c2p golive
	// W update 2c2p to Sandbox
	// y update ktb gen qrcode
	// Y update ktb json structure
	// z update del warehouse update isenable = 0
	// Z update PromptSong SMS Sendername
	// V 1.0.3a  update ktb approve
	// 1.0.3A update getdriver insert or update แก้ไขเรื่องเบิ้ล
	// 1.0.3b update ktb not approve
	// 1.0.3B update FntGetpaymentToken for market place
	// 1.0.3c update login tikky not accbank
	// 1.0.3d update db to prd.
	// 1.0.3e update payment qrcode to 2c2p.
	// 1.0.3f update app version update.
	// 1.0.3F update Auth All.
	// 1.0.3g update 2C2P Production.
	// 1.0.3G update QRCode Direct link Production.
	// 1.0.3h update check ktb user pass before db open.
	// 1.0.3H update pass ktb .
	// 1.0.3i update select carrier  tms.
	// 1.0.3I keep log count.
	// 1.0.3j Add BankAccName in OMSGetBankAccount.
	// 1.0.3J Add Response case dup, invalid price in OMSGetBankAccount.
	// 1.0.3k edit case panic.
	// 1.0.3K Add Log to tblmobileapilog from ktbapprove.
	// V 1.0.3K5S  add Show Alert Carrer Get Job
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint Hit: returnAllArticles")
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

func Callback(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)

	w.Write(reqBody)
	// var article GetPostPrintOnline
	// json.Unmarshal(reqBody, &article)

	// //dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// //dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// //var m MyStruct
	// typefind := article.PostCode
	// //Counti := strconv.Itoa(article.CountItems)

	// //t := time.Now() //It will return time.Time object with current timestamp
	// //fmt.Printf("time.Time %s\n", t)

	// //tUnix := t.Unix()
	// fmt.Printf("timeUnix: \n", typefind)
	// //str := strconv.FormatInt(tUnix, 10)

	// ress2, err := db.Query("SELECT THPDID,THPDName,lblContainerID FROM THPDDB.tblTHPDPrintConfig WHERE THPDID = '" + typefind + "' ")

	// if err == nil {

	// 	boxes := []EventPostPrint{}

	// 	//boxes = append(boxes, BoxData{Width: 5, Height: 30})

	// 	resp := make(map[string]string)
	// 	resp["THPDID"] = article.PostCode

	// 	for ress2.Next() {
	// 		var event EventPostPrint
	// 		//JobID := ress2.Scan(&event.JobID)
	// 		err := ress2.Scan(&event.THPDID, &event.THPDName, &event.LblContainerID)

	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		boxes = append(boxes, EventPostPrint{THPDID: event.THPDID, LblContainerID: event.LblContainerID, THPDName: event.THPDName})
	// 	}

	// 	b, _ := json.Marshal(boxes)

	// 	defer ress2.Close()
	// 	err = ress2.Close()
	// 	//jsonResp, err := json.Marshal(b)
	// 	if err != nil {
	// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 		return
	// 	}
	// 	w.Write(b)
	//counter = 0

	//}
	//return fmt.Sprintf("%x", b)
}
func GenlblPrint(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article GetPostPrintOnline
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
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
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.PostCode
	Counti := strconv.Itoa(article.CountItems)

	ress, err := db.Query("UPDATE THPDDB.tblTHPDPrintConfig SET lblContainerID = lblContainerID + " + Counti + "  WHERE THPDID = '" + article.PostCode + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	}

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: \n", typefind)
	//str := strconv.FormatInt(tUnix, 10)

	ress2, err := db.Query("SELECT THPDID,THPDName,lblContainerID FROM THPDDB.tblTHPDPrintConfig WHERE THPDID = '" + typefind + "' ")

	if err == nil {

		boxes := []EventPostPrint{}

		//boxes = append(boxes, BoxData{Width: 5, Height: 30})

		resp := make(map[string]string)
		resp["THPDID"] = article.PostCode

		for ress2.Next() {
			var event EventPostPrint
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.THPDID, &event.THPDName, &event.LblContainerID)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, EventPostPrint{THPDID: event.THPDID, LblContainerID: event.LblContainerID, THPDName: event.THPDName})
		}

		b, _ := json.Marshal(boxes)

		defer ress2.Close()
		err = ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return //ost
		}
		w.Write(b)
		//counter = 0

	}
	//return fmt.Sprintf("%x", b)
}
func tokenGenerator(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	rand.Read(b)
	//str := fmt.Sprintf("%x", b)
	//w.Write(str)
	resp := make(map[string]string)
	resp["token"] = fmt.Sprintf("%x", b)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	//return fmt.Sprintf("%x", b)
}

// PRD 20s time out
func getDNSString(dbName, dbUser, dbPassword, conn string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=20s&readTimeout=20s",
		dbUser,
		dbPassword,
		conn,
		dbName)
}

// UAT Use 60s time out
// func getDNSString(dbName, dbUser, dbPassword, conn string) string {
// 	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=60s&readTimeout=60s",
// 		dbUser,
// 		dbPassword,
// 		conn,
// 		dbName)
// }

// func selectEventById(db *sql.DB, id string, event *Event0) error {
// 	row := db.QueryRow(selectEventByIdQueryTOAT, id)
// 	err := row.Scan(&event.InvoiceNo, &event.TrackingID, &event.Truck_type, &event.Truck_plate, &event.Driver_name, &event.Driver_Phone, &event.Total_chest, &event.Total_pack, &event.CustomerNo, &event.InvoiceDate, &event.To_created_date)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func selectEventByCusId(db *sql.DB, id string, dt string, event *Event0) error {
// 	row := db.QueryRow(selectEventByCusIdQueryTOAT, id, dt)
// 	err := row.Scan(&event.InvoiceNo, &event.TrackingID, &event.Truck_type, &event.Truck_plate, &event.Driver_name, &event.Driver_Phone, &event.Total_chest, &event.Total_pack, &event.CustomerNo, &event.InvoiceDate, &event.To_created_date)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func selectEventByCusToken(db *sql.DB, id string, token string, event *Event1) error {
	row := db.QueryRow(selectEventByToken, id, token)
	err := row.Scan(&event.Email, &event.MerchantName)
	if err != nil {
		return err
	}
	return nil
}

func selectOMSOrder(db *sql.DB, Mobileid string, token string, event *Event5) error {
	row := db.QueryRow(selectOrderOMS, Mobileid)
	err := row.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode,
		&event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.IMG1,
		&event.IMG2, &event.IMG3, &event.IMG4, &event.CreateDt, &event.JobDesc)
	if err != nil {
		return err
	}
	return nil
}
func insertOrderItems(db *sql.DB, event Event5) (int64, error) {
	res, err := db.Exec(insertOrderItem, event.JobID, event.MobileID, event.ReceiveName, event.ReceiveAddress, event.ReceiveTumbon, event.ReceiveDistrict, event.ReceiveProvince, event.ReceiveZipcode,
		event.ReceivePhoneNo, event.SenderZipcode, event.SendPrice, event.PickupStartDt, event.DeliveryEndDt, event.PaymentFlag, event.PaymentDetail,
		event.JobDesc, event.JobType, event.WarehouseID, event.MerchantID, event.Status)
	if err != nil {
		return 0, err
	}
	lid1, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	res2, err := db.Exec(insertOrderItemDesc, event.JobID, event.MobileID, event.IMG1,
		event.IMG2, event.IMG3, event.IMG4)
	if err != nil {
		return 0, err
	}
	lid2, err := res2.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid1 + lid2, nil

}

func insertPayment(db *sql.DB, event EventPayment) (int64, error) {
	res, err := db.Exec(insertGetPayment, event.refID, event.invoiceNo, "", event.phoneNo, event.amount, event.cardType, event.cardNo, event.respDesc, event.transactionDatetime, event.refCode)
	if err != nil {
		return 0, err
	}
	lid1, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid1, nil

}
func insertGetDriverTrack(db *sql.DB, event EventDriverTrack) (int64, error) {
	res, err := db.Exec(insertGetDriver, event.JobID, event.TrackingID, event.CarID, event.DriverName, event.DriverPhoneNo, event.Company, event.Price, event.CustomerSelect)
	if err != nil {
		return 0, err
	}
	lid1, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lid1, nil

}

// func returnSingleArticle2(w http.ResponseWriter, r *http.Request) {
// 	//vars := mux.Vars(r)
// 	//key := vars["id"]

// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	var article Article
// 	json.Unmarshal(reqBody, &article)

// 	dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
// 	db, err := sql.Open("mysql", dns)

// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	typefind := article.Type

// 	firstEvent2 := Event1{}
// 	err = selectEventByCusToken(db, article.User, article.Token, &firstEvent2)
// 	if err != nil {
// 		resperr := make(map[string]string)
// 		resperr["errmsg"] = err.Error()

// 		jsonResp, err := json.Marshal(resperr)
// 		if err != nil {
// 			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 		}
// 		w.Write(jsonResp)
// 		return
// 	}

// 	//datett := article.Date
// 	if typefind == "InvoiceNo" {
// 		////fmt.Println(article.Id)

// 		firstEvent := Event0{}
// 		err = selectEventById(db, article.Id, &firstEvent)
// 		if err != nil {
// 			panic(err)
// 		}

// 		resp := make(map[string]string)
// 		resp["InvoiceNo"] = firstEvent.InvoiceNo
// 		resp["TrackingID"] = firstEvent.TrackingID
// 		resp["Truck_type"] = firstEvent.Truck_type
// 		resp["Truck_plate"] = firstEvent.Truck_plate
// 		resp["Driver_name"] = firstEvent.Driver_name
// 		resp["Driver_Phone"] = firstEvent.Driver_Phone
// 		resp["Total_chest"] = firstEvent.Total_chest
// 		resp["Total_pack"] = firstEvent.Total_pack
// 		resp["CustomerNo"] = firstEvent.CustomerNo
// 		resp["InvoiceDate"] = firstEvent.InvoiceDate
// 		resp["Start_date"] = firstEvent.To_created_date

// 		jsonResp, err := json.Marshal(resp)
// 		if err != nil {
// 			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 		}
// 		w.Write(jsonResp)
// 	}
// 	if typefind == "CustomerNo" {
// 		////fmt.Println(article.Id)
// 		firstEvent := Event0{}
// 		err = selectEventByCusId(db, article.Id, article.Date, &firstEvent)
// 		if err != nil {
// 			panic(err)
// 		}

// 		resp := make(map[string]string)
// 		resp["InvoiceNo"] = firstEvent.InvoiceNo
// 		resp["TrackingID"] = firstEvent.TrackingID
// 		resp["Truck_type"] = firstEvent.Truck_type
// 		resp["Truck_plate"] = firstEvent.Truck_plate
// 		resp["Driver_name"] = firstEvent.Driver_name
// 		resp["Driver_Phone"] = firstEvent.Driver_Phone
// 		resp["Total_chest"] = firstEvent.Total_chest
// 		resp["Total_pack"] = firstEvent.Total_pack
// 		resp["CustomerNo"] = firstEvent.CustomerNo
// 		resp["InvoiceDate"] = firstEvent.InvoiceDate
// 		resp["Start_date"] = firstEvent.To_created_date

// 		jsonResp, err := json.Marshal(resp)
// 		if err != nil {
// 			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			return
// 		}
// 		w.Write(jsonResp)

// 	}
// 	// if datett == "" {
// 	// 	//fmt.Println(article.Id)
// 	// }

// 	////fmt.Println(article.Id)
// 	//fmt.Fprintf(w, jsonResp)
// 	// for _, article := range Articles {
// 	// 	if article.Id == key {
// 	// 		json.NewEncoder(w).Encode(article)
// 	// 	}
// 	// }
// }

type Epoch int64
type MyStruct struct {
	Date Epoch
}

// func OMSMobileCreateOrder(w http.ResponseWriter, r *http.Request) {
// 	//vars := mux.Vars(r)
// 	//key := vars["id"]

// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	var article OMSOrderStruct
// 	json.Unmarshal(reqBody, &article)

// 	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
// 	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

// 	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
// 	// db, err := sql.Open("mysql", dns)

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// err = db.Ping()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer db.Close()

// 	if DB.Ping() != nil {
// 		connectDb()
// 	}
// 	if DB.Stats().OpenConnections != 0 {
// 		//fmt.Println(DB.Stats().OpenConnections)
// 	} else {
// 		connectDb()
// 	}
// 	db := DB
// 	defer db.Close()
// 	// err := db.Ping()
// 	// if err != nil {
// 	// 	connectDb()

// 	// }
// 	//var m MyStruct
// 	typefind := article.MobileID

// 	//t := time.Now() //It will return time.Time object with current timestamp
// 	//fmt.Printf("time.Time %s\n", t)

// 	//tUnix := t.Unix()
// 	//fmt.Printf("timeUnix: %d\n", tUnix)
// 	//str := strconv.FormatInt(tUnix, 10)

// 	if typefind != "" {
// 		//var unixTime int64 = int64(m.Date)

// 		event5 := Event5{
// 			JobID:           string(article.JobID),
// 			MobileID:        string(article.MobileID),
// 			ReceiveName:     string(article.ReceiveName),
// 			ReceiveAddress:  string(article.ReceiveAddress),
// 			ReceiveTumbon:   string(article.ReceiveTumbon),
// 			ReceiveDistrict: string(article.ReceiveDistrict),
// 			ReceiveProvince: string(article.ReceiveProvince),
// 			ReceivePhoneNo:  string(article.ReceivePhoneNo),
// 			ReceiveZipcode:  string(article.ReceiveZipcode),
// 			SenderZipcode:   string(article.SenderZipcode),
// 			SendPrice:       string(article.SendPrice),
// 			PickupStartDt:   string(article.PickupStartDt),
// 			DeliveryEndDt:   string(article.DeliveryEndDt),
// 			PaymentFlag:     string(article.PaymentFlag),
// 			PaymentDetail:   string(article.PaymentDetail),
// 			IMG1:            string(article.IMG1),
// 			IMG2:            string(article.IMG2),
// 			IMG3:            string(article.IMG3),
// 			IMG4:            string(article.IMG4),
// 			JobDesc:         string(article.JobDesc),
// 			JobType:         string(article.JobType),
// 			MerchantID:      string(article.MerchantID),
// 			WarehouseID:     string(article.WarehouseID),
// 			Status:          "รอดำเนินการ",
// 			//CreateDt:        string(article.CreateDt),
// 		}
// 		insertedId, err := insertOrderItems(db, event5)
// 		fmt.Println(insertedId)
// 		if err != nil {
// 			//panic(err)
// 			resperr := make(map[string]string)
// 			resperr["errmsg"] = err.Error()

// 			jsonResp, err := json.Marshal(resperr)
// 			if err != nil {
// 				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			}
// 			w.Write(jsonResp)
// 			return
// 		} else {
// 			respok := make(map[string]string)
// 			respok["Success"] = "Insert Success"
// 			jsonResp, err := json.Marshal(respok)
// 			if err != nil {
// 				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			}
// 			w.Write(jsonResp)
// 			return
// 		}
// 		////fmt.Println(insertedId)
// 	}

// 	//return

// }
// func OMSMobileCreateOrderWithAuth(w http.ResponseWriter, r *http.Request) {
// 	//vars := mux.Vars(r)
// 	//key := vars["id"]

// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	var article OMSOrderStruct
// 	json.Unmarshal(reqBody, &article)

// 	reqHeader := r.Header["Authorization"]
// 	////fmt.Println("Header", reqHeader)

// 	reqHeaderChannel := r.Header["Channel"]
// 	////fmt.Println("Header", reqHeaderChannel)

// 	aa := ChkAuth(reqHeader, reqHeaderChannel)
// 	//fmt.Println("aa", aa)

// 	if aa <= 0 {
// 		respok := make(map[string]string)
// 		//respok["err"] = "Existing Invoice Number"

// 		respok["QRcodetxt"] = ""
// 		respok["ExpDT"] = ""                                  //QR
// 		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

// 		jsonResp, err := json.Marshal(respok)

// 		if err != nil {
// 			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			return
// 		}
// 		w.Write(jsonResp)
// 		return
// 	}

// 	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
// 	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

// 	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
// 	// db, err := sql.Open("mysql", dns)

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// err = db.Ping()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer db.Close()

// 	if DB.Ping() != nil {
// 		connectDb()
// 	}
// 	if DB.Stats().OpenConnections != 0 {
// 		//fmt.Println(DB.Stats().OpenConnections)
// 	} else {
// 		connectDb()
// 	}
// 	db := DB
// 	defer db.Close()
// 	// err := db.Ping()
// 	// if err != nil {
// 	// 	connectDb()

// 	// }
// 	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileCreateOrderWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
// 	defer ress99.Close()
// 	if err != nil {
// 		panic(err)
// 	}
// 	//var m MyStruct
// 	typefind := article.MobileID

// 	//t := time.Now() //It will return time.Time object with current timestamp
// 	//fmt.Printf("time.Time %s\n", t)

// 	//tUnix := t.Unix()
// 	//fmt.Printf("timeUnix: %d\n", tUnix)
// 	//str := strconv.FormatInt(tUnix, 10)

// 	if typefind != "" {
// 		//var unixTime int64 = int64(m.Date)

// 		event5 := Event5{
// 			JobID:           string(article.JobID),
// 			MobileID:        string(article.MobileID),
// 			ReceiveName:     string(article.ReceiveName),
// 			ReceiveAddress:  string(article.ReceiveAddress),
// 			ReceiveTumbon:   string(article.ReceiveTumbon),
// 			ReceiveDistrict: string(article.ReceiveDistrict),
// 			ReceiveProvince: string(article.ReceiveProvince),
// 			ReceivePhoneNo:  string(article.ReceivePhoneNo),
// 			ReceiveZipcode:  string(article.ReceiveZipcode),
// 			SenderZipcode:   string(article.SenderZipcode),
// 			SendPrice:       string(article.SendPrice),
// 			PickupStartDt:   string(article.PickupStartDt),
// 			DeliveryEndDt:   string(article.DeliveryEndDt),
// 			PaymentFlag:     string(article.PaymentFlag),
// 			PaymentDetail:   string(article.PaymentDetail),
// 			IMG1:            string(article.IMG1),
// 			IMG2:            string(article.IMG2),
// 			IMG3:            string(article.IMG3),
// 			IMG4:            string(article.IMG4),
// 			JobDesc:         string(article.JobDesc),
// 			JobType:         string(article.JobType),
// 			MerchantID:      string(article.MerchantID),
// 			WarehouseID:     string(article.WarehouseID),
// 			Status:          "รอดำเนินการ",
// 			//CreateDt:        string(article.CreateDt),
// 		}
// 		insertedId, err := insertOrderItems(db, event5)
// 		//fmt.Println(insertedId)
// 		if err != nil {
// 			//panic(err)
// 			resperr := make(map[string]string)
// 			resperr["errmsg"] = err.Error()

// 			jsonResp, err := json.Marshal(resperr)
// 			if err != nil {
// 				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			}
// 			w.Write(jsonResp)
// 			return
// 		} else {
// 			respok := make(map[string]string)
// 			respok["Success"] = "Insert Success"
// 			jsonResp, err := json.Marshal(respok)
// 			if err != nil {
// 				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
// 			}
// 			w.Write(jsonResp)
// 			return
// 		}
// 		////fmt.Println(insertedId)
// 	}

// 	//return

// }

// type Bird struct {
// 	Species     string
// 	Description string
// }

type BoxData struct {
	Width  int
	Height int
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func OMSMobileConnect(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileConnect','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0
		drivername := strings.Split(article.DriverName, ",")[0]
		drivercompany := strings.Split(article.DriverName, ",")[1]

		ress, err := db.Query("UPDATE THPDMPDB.tblJobMaster SET JobDriverID = '" + article.DriverID + "',JobDriverName = '" + drivername + "', JobTruckPlate = '" + article.Plate + "' , FTLFinalPrice= '" + article.FinalPrice + "' , JobAssignDT = NOW() , jobupdatedt = NOW()  WHERE JobID in (SELECT JobID FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + article.TrackingID + "' ) ")
		defer ress.Close()
		if err != nil {
			panic(err)
		} else {

			/// ยิง api เลือกคนขับ
			UpdateAPIDriverToLoadboard(db, article.TrackingID, article.DriverID, drivercompany)

			// ยิงเสริม
			UpdateAPIDriverToLoadboard2(db, article.TrackingID, article.DriverID, drivercompany)

			ress2, err2 := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET CustomerSelect = '1', CustomConfirmDT = NOW() WHERE  TrackingID = '" + article.TrackingID + "' and DriverName = '" + drivername + "' ")
			defer ress2.Close()
			if err2 != nil {
				panic(err)
			}
			ress3, err3 := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'รอขนส่งรับสินค้า' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + article.TrackingID + "')")
			defer ress3.Close()
			if err3 != nil {
				panic(err)
			}

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID
			b, _ := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
		}

	}

}
func OMSMobileUpdateDriverJobMaster(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0
		drivername := strings.Split(article.DriverName, ",")[0]
		drivercompany := strings.Split(article.DriverName, ",")[1]

		ress, err := db.Query("UPDATE THPDMPDB.tblJobMaster SET JobDriverID = '" + article.DriverID + "',JobDriverName = '" + drivername + "', JobTruckPlate = '" + article.Plate + "' , FTLFinalPrice= '" + article.FinalPrice + "' , JobAssignDT = NOW() , jobupdatedt = NOW()  WHERE JobID in (SELECT JobID FROM THPDMPDB.tblOrderMaster WHERE JobID = '" + article.TrackingID + "' ) ")
		defer ress.Close()
		if err != nil {
			panic(err)
		} else {

			/// ยิง api เลือกคนขับ
			UpdateAPIDriverToLoadboard(db, article.TrackingID, article.DriverID, drivercompany)

			// ยิงเสริม
			UpdateAPIDriverToLoadboard2(db, article.TrackingID, article.DriverID, drivercompany)

			ress2, err2 := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET CustomerSelect = '1', CustomConfirmDT = NOW() WHERE  TrackingID = '" + article.TrackingID + "' and DriverName = '" + drivername + "' ")
			defer ress2.Close()
			if err2 != nil {
				panic(err)
			}
			ress3, err3 := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'Booking' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE JobID = '" + article.TrackingID + "')")
			defer ress3.Close()
			if err3 != nil {
				panic(err)
			}

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID
			b, _ := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
		}

	}

}
func OMSMobileUpdateDriverJobMasterWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileUpdateDriverJobMasterWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0
		drivername := strings.Split(article.DriverName, ",")[0]
		drivercompany := strings.Split(article.DriverName, ",")[1]

		ress, err := db.Query("UPDATE THPDMPDB.tblJobMaster SET JobDriverID = '" + article.DriverID + "',JobDriverName = '" + drivername + "', JobTruckPlate = '" + article.Plate + "' , FTLFinalPrice= '" + article.FinalPrice + "' , JobAssignDT = NOW() , jobupdatedt = NOW()  WHERE JobID in (SELECT JobID FROM THPDMPDB.tblOrderMaster WHERE JobID = '" + article.TrackingID + "' ) ")
		defer ress.Close()
		if err != nil {
			panic(err)
		} else {

			// ยิงเสริม
			//UpdateAPIDriverToLoadboard2(article.TrackingID, article.DriverID, drivercompany)

			ress2, err2 := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET CustomerSelect = '1', CustomConfirmDT = NOW() WHERE  TrackingID = '" + article.TrackingID + "' and DriverName = '" + drivername + "' ")
			defer ress2.Close()
			if err2 != nil {
				panic(err)
			}
			ress3, err3 := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'Booking' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE JobID = '" + article.TrackingID + "')")
			defer ress3.Close()
			if err3 != nil {
				panic(err)
			}

			/// ยิง api เลือกคนขับ
			UpdateAPIDriverToLoadboard(db, article.TrackingID, article.DriverID, drivercompany)

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID
			b, _ := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
		}

	}

}
func OMSMobileGetTrcukSize(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article TruckType
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TruckID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT  TruckID, TruckName, TruckDesc, TruckSize, RateServiceStart, OilStandardPrice as RateOilPerKM, IsUse FROM THPDDB.VW_TruckType ")

		if err == nil {

			boxes := []TruckType{}

			resp := make(map[string]string)
			resp["TruckID"] = article.TruckID

			for ress2.Next() {
				var event TruckType
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.TruckID, &event.TruckName, &event.TruckDesc, &event.TruckSize, &event.RateServiceStart, &event.RateOilPerKM, &event.IsUse)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, TruckType{TruckID: event.TruckID, TruckName: event.TruckName, TruckDesc: event.TruckDesc, TruckSize: event.TruckSize, RateServiceStart: event.RateServiceStart, RateOilPerKM: event.RateOilPerKM, IsUse: event.IsUse})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}
func OMSMobileGetTrcukSizeWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article TruckType
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	//var m MyStruct
	typefind := article.TruckID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	theTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 100, time.Local)
	//showchannel := "OMSMobileGetDriver"
	d := theTime.Format("2006-1-2")
	apiname := "OMSMobileGetTrcukSizeWithAuth_" + d

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('" + apiname + "','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT  TruckID, TruckName, TruckDesc, TruckSize, RateServiceStart, OilStandardPrice as RateOilPerKM, IsUse FROM THPDDB.VW_TruckType ")

		if err == nil {

			boxes := []TruckType{}

			resp := make(map[string]string)
			resp["TruckID"] = article.TruckID

			for ress2.Next() {
				var event TruckType
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.TruckID, &event.TruckName, &event.TruckDesc, &event.TruckSize, &event.RateServiceStart, &event.RateOilPerKM, &event.IsUse)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, TruckType{TruckID: event.TruckID, TruckName: event.TruckName, TruckDesc: event.TruckDesc, TruckSize: event.TruckSize, RateServiceStart: event.RateServiceStart, RateOilPerKM: event.RateOilPerKM, IsUse: event.IsUse})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}
func OMSMobileGetCoupon(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT CouponName,CouponCode,CouponDiscount,CouponExpireDT FROM THPDMPDB.tblCoupons  ")

		if err == nil {

			boxes := []Coupon{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Coupon
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.CouponName, &event.CouponCode, &event.CouponDiscount, &event.CouponExpireDT)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Coupon{CouponName: event.CouponName, CouponCode: event.CouponCode, CouponDiscount: event.CouponDiscount, CouponExpireDT: event.CouponExpireDT})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}
func OMSMobileGetCouponWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileGetCouponWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT CouponName,CouponCode,CouponDiscount,CouponExpireDT FROM THPDMPDB.tblCoupons  ")

		if err == nil {

			boxes := []Coupon{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Coupon
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.CouponName, &event.CouponCode, &event.CouponDiscount, &event.CouponExpireDT)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Coupon{CouponName: event.CouponName, CouponCode: event.CouponCode, CouponDiscount: event.CouponDiscount, CouponExpireDT: event.CouponExpireDT})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}

func OMSMobileGetPickupTime(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	// reqBody, _ := ioutil.ReadAll(r.Body)
	// var article Merchant
	// json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT TimeType, TimeName  FROM THPDDB.tblpickuptimeslot  ")

		if err == nil {

			boxes := []PickupTime{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event PickupTime
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.TimeType, &event.TimeName)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, PickupTime{TimeType: event.TimeType, TimeName: event.TimeName})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}

func OMSMobileGetRateAddOn(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT DBID, AddonName, AddonDesc, AddonPrice, IsUse, ModifyDt, CreateDT FROM THPDDB.tblRateAddon  ")

		if err == nil {

			boxes := []RateAddOn{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event RateAddOn
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.DBID, &event.AddonName, &event.AddonDesc, &event.AddonPrice, &event.IsUse, &event.ModifyDt, &event.CreateDT)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, RateAddOn{DBID: event.DBID, AddonName: event.AddonName, AddonDesc: event.AddonDesc, AddonPrice: event.AddonPrice, IsUse: event.IsUse, ModifyDt: event.ModifyDt, CreateDT: event.CreateDT})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}

func OMSMobileDelWareHouse(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0
		ress3, err3 := db.Query("UPDATE THPDMPDB.tblMerchantWH SET IsEnable = 0 WHERE WarehouseID = '" + article.WarehouseID + "'  ")
		defer ress3.Close()
		if err3 != nil {
			panic(err3)
		}

		if err == nil {

			ress2, err := db.Query("SELECT MerchantID, Address1, Address2, IFNULL(LocationGPS,'0,0')LocationGPS, WarehouseID, IFNULL(SubDistrict,'')SubDistrict, IFNULL(District,'0,0')District,ProvinceName,PostCode FROM THPDMPDB.tblMerchantWH WHERE MerchantID = '" + article.MerchantID + "' and IsEnable = 1 ")

			boxes := []Merchant{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Merchant
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.Address1, &event.Address2, &event.LocationGPS, &event.WarehouseID, &event.SubDistrict, &event.District, &event.ProvinceName, &event.PostCode)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

			// respok := make(map[string]string)
			// respok["Success"] = "true"
			// respok["WarehouseID"] = article.WarehouseID

			// jsonResp, err := json.Marshal(respok)

			// if err != nil {
			// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// }
			// w.Write(jsonResp)
			// return

		}

	}

}
func OMSMobileDelWareHouseWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileDelWareHouseWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0
		ress3, err3 := db.Query("UPDATE THPDMPDB.tblMerchantWH SET IsEnable = 0 WHERE WarehouseID = '" + article.WarehouseID + "'  ")
		defer ress3.Close()
		if err3 != nil {
			panic(err3)
		}

		if err == nil {

			ress2, err := db.Query("SELECT MerchantID, Address1, Address2, IFNULL(LocationGPS,'0,0')LocationGPS, WarehouseID, IFNULL(SubDistrict,'')SubDistrict, IFNULL(District,'0,0')District,ProvinceName,PostCode FROM THPDMPDB.tblMerchantWH WHERE MerchantID = '" + article.MerchantID + "' and IsEnable = 1 ")

			boxes := []Merchant{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Merchant
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.Address1, &event.Address2, &event.LocationGPS, &event.WarehouseID, &event.SubDistrict, &event.District, &event.ProvinceName, &event.PostCode)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

			// respok := make(map[string]string)
			// respok["Success"] = "true"
			// respok["WarehouseID"] = article.WarehouseID

			// jsonResp, err := json.Marshal(respok)

			// if err != nil {
			// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// }
			// w.Write(jsonResp)
			// return

		}

	}

}
func OMSMobileGetWareHouse(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT MerchantID, Address1, Address2, IFNULL(LocationGPS,'0,0')LocationGPS, WarehouseID, IFNULL(SubDistrict,'')SubDistrict, IFNULL(District,'0,0')District,ProvinceName,PostCode FROM THPDMPDB.tblMerchantWH WHERE MerchantID = '" + article.MerchantID + "' and IsEnable = 1 ")

		if err == nil {

			boxes := []Merchant{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Merchant
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.Address1, &event.Address2, &event.LocationGPS, &event.WarehouseID, &event.SubDistrict, &event.District, &event.ProvinceName, &event.PostCode)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}
func OMSMobileGetWareHouseWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	//var m MyStruct
	typefind := article.MerchantID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	theTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 100, time.Local)
	//showchannel := "OMSMobileGetDriver"
	d := theTime.Format("2006-1-2")
	apiname := "OMSMobileGetWareHouseWithAuth_" + typefind + "_" + d

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('" + apiname + "','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	// strOTP = substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {

		//counter := 0

		ress2, err := db.Query("SELECT MerchantID, Address1, Address2, IFNULL(LocationGPS,'0,0')LocationGPS, WarehouseID,WarehouseName, IFNULL(SubDistrict,'')SubDistrict, IFNULL(District,'0,0')District,ProvinceName,PostCode, ContactName,PhoneNumber,IFNULL(EmailContact,'')EmailContact,IFNULL(StartWorkingHour,'')StartWorkingHour,IFNULL(EndWorkingHour,'')EndWorkingHour FROM THPDMPDB.tblMerchantWH WHERE MerchantID = '" + article.MerchantID + "' and IsEnable = 1 ")

		if err == nil {

			boxes := []Merchant{}

			resp := make(map[string]string)
			resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event Merchant
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.Address1, &event.Address2, &event.LocationGPS, &event.WarehouseID, &event.WarehouseName, &event.SubDistrict, &event.District, &event.ProvinceName, &event.PostCode, &event.ContactName, &event.PhoneNumber, &event.EmailContact, &event.StartWorkingHour, &event.EndWorkingHour)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, WarehouseName: event.WarehouseName, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode, ContactName: event.ContactName, PhoneNumber: event.PhoneNumber, EmailContact: event.EmailContact, StartWorkingHour: event.StartWorkingHour, EndWorkingHour: event.EndWorkingHour})

			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)

		}

	}

}
func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
func CancelJobToLoadboard(trackID string, remark string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte

	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	if trackID == "" {
		return "0"
	}

	payload := strings.NewReader(`{"shipment_no": "` + trackID + `" ,"remark": "` + remark + `"}`)

	method := "POST"
	url := tmsapi + "API/BrokerGateways/cancelPostLoad"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	req.Header.Add("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		//fmt.Println(err)
		return string(err.Error())
	}
	//fmt.Println(string(body))

	// sqlstr := "INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  "
	// //fmt.Println(string(sqlstr))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	// ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET PriceUnit = '1'  WHERE TrackingID = '" + trackID + "'  ")
	// defer ress.Close()
	// if err != nil {
	// 	panic(err)
	// } else {
	/// SMS ให้คนขับ ///

	SendMessageToDriver(trackID, "0850270971", "ยกเลิกการจัดส่ง", remark)

	//}
	if err == nil {
	}

	return string("ok")
}
func CancelJobToLoadboard2(trackID string, remark string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte

	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	if trackID == "" {
		return "0"
	}

	payload := strings.NewReader(`{"shipment_no": "` + trackID + `" ,"remark": "` + remark + `"}`)

	method := "POST"
	url := "https://demo-api.dxplace.com/API/BrokerGateways/cancelPostLoad"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	req.Header.Add("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		//fmt.Println(err)
		return string(err.Error())
	}
	//fmt.Println(string(body))

	// sqlstr := "INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  "
	// //fmt.Println(string(sqlstr))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	// ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET PriceUnit = '1'  WHERE TrackingID = '" + trackID + "'  ")
	// defer ress.Close()
	// if err != nil {
	// 	panic(err)
	// } else {
	/// SMS ให้คนขับ ///

	//SendMessageToDriver(trackID, "0850270971", "ยกเลิกการจัดส่ง", remark)

	//}
	if err == nil {
	}

	return string("ok")
}
func UpdateAPIDriverToLoadboard(db *sql.DB, trackID string, carrier_id string, drivername string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	payload := strings.NewReader(`{"con_no": "` + trackID + `" ,"carrier_id": "` + carrier_id + `"}`)

	method := "POST"
	url := tmsapi + "API/BrokerGateways/selectCarrier"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	req.Header.Add("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		//fmt.Println(err)
		return string(err.Error())
	}
	//fmt.Println(string(body))

	// sqlstr := "INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  "
	// //fmt.Println(string(sqlstr))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET PriceUnit = '1'  WHERE TrackingID = '" + trackID + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	} else {
		/// SMS ให้คนขับ ///
		//SendMessageToDriver(trackID, "0850270971", carrier_id, drivername)

	}
	if err == nil {
	}

	return string("ok")
}
func UpdateAPIDriverToLoadboard2(db *sql.DB, trackID string, carrier_id string, drivername string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	payload := strings.NewReader(`{"con_no": "` + trackID + `" ,"carrier_id": "` + carrier_id + `"}`)

	method := "POST"
	url := "https://demo-api.dxplace.com/API/BrokerGateways/selectCarrier"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	req.Header.Add("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		//fmt.Println(err)
		return string(err.Error())
	}
	//fmt.Println(string(body))

	// sqlstr := "INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  "
	// //fmt.Println(string(sqlstr))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	// ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET PriceUnit = '1'  WHERE TrackingID = '" + trackID + "'  ")
	// defer ress.Close()
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	/// SMS ให้คนขับ ///

	// 	//SendMessageToDriver(trackID, "0850270971", carrier_id, drivername)

	// }
	// if err == nil {
	// }
	return string("ok")
}

//

func PreparePayload(trackID string) string {

	// userData := map[string]interface{}{"merchantID": "JT04", "invoiceNo": "AX000002956TX", "description": "AX000002956TX", "amount": 10000.00, "currencyCode": "THB"}
	// //access := L.Sign(userData, "CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2", 1)

	// accessToken, err := Sign(userData, "CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2", 1) // data -> secretkey env name -> expiredAt
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	// 	os.Exit(1)
	// }

	// return string(accessToken)
	return ""
}

// func Sign(Data map[string]interface{}, SecrePublicKeyEnvName string, ExpiredAt time.Duration) (string, error) {

// 	expiredAt := time.Now().Add(time.Duration(time.Second) * ExpiredAt).Unix()

// 	//jwtSecretKey := GodotEnv(SecrePublicKeyEnvName)
// 	jwtSecretKey := SecrePublicKeyEnvName
// 	// metadata for your jwt
// 	claims := jwtv4.MapClaims{}
// 	claims["expiredAt"] = expiredAt
// 	// claims["authorization"] = true

// 	for i, v := range Data {
// 		claims[i] = v
// 	}

// 	to := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, claims)
// 	accessToken, err := to.SignedString([]byte(jwtSecretKey))

// 	fmt.Print(accessToken)

// 	if err != nil {
// 		//logrus.Error(err.Error())
// 		return accessToken, err
// 	}

// 	return accessToken, nil
// }
// func ConnectDB() {

// 	dns := getDNSString(dblogin, userlogin, passlogin, conn)
// 	db, err := sql.Open("mysql", dns)

// 	//db, err := sql.Open("mysql", "username:password@/dbname")
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	DB = db

// }

func ChkAuth(Authorization []string, Channel []string) int {

	//db := ConnectDB

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//st := db.Stats
	////fmt.Println(st)

	result1 := strings.Index(Channel[0], "OMSDEWeb")
	result2 := strings.Index(Channel[0], "OMSDEMobile")
	result3 := strings.Index(Channel[0], "LineOADE")
	result4 := strings.Index(Channel[0], "TMSDxMobile")
	result5 := strings.Index(Channel[0], "MarketPlaceDEWeb")
	result6 := strings.Index(Channel[0], "OMSDEMobileAndroid")
	result7 := strings.Index(Channel[0], "MarketPlaceDEMobile")
	result8 := strings.Index(Channel[0], "OMSDEMobileService")
	result9 := strings.Index(Channel[0], "CESSRAOT")

	resultAuth1 := strings.Index(Authorization[0], "bG9hZGJvYXJkVXNyOjU0MzQ1Mw==")
	resultAuth2 := strings.Index(Authorization[0], "QUNDT1VOVEBUSFBEQEFudDEyMzQ1")
	resultAuth3 := strings.Index(Authorization[0], "T0ZGSUNFUkBUSFBEQE9mZjEyMzQ1")
	resultAuth4 := strings.Index(Authorization[0], "QURNSU5AVEhQREBBcHAxMjM0NQ==")

	passed := ""
	if result1 > -1 || result2 > -1 || result3 > -1 || result4 > -1 || result5 > -1 || result6 > -1 || result7 > -1 || result8 > -1 || result9 > -1 {

		if resultAuth1 > -1 || resultAuth2 > -1 || resultAuth3 > -1 || resultAuth4 > -1 {
			passed = "ok"
		} else {
			return 0
		}
	} else {
		return 0
	}
	fmt.Println(passed)

	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	//st := db.Stats
	////fmt.Println(st)

	if err != nil {

		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//db := DB
	//defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()
	// }

	//showchannel := "ChkAuth_" + Channel[0]
	// ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('" + showchannel + "','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	ress99, err := db.Query("INSERT INTO raotdb.tblauthchannel (Cid, ChannelName, ChannelSecretKey,UseCount,LastUseDT) Values (UNIX_TIMESTAMP(),'" + Channel[0] + "','" + Authorization[0] + "','0' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastUseDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), UseCount = UseCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	ress2, err := db.Query("SELECT Cid, ChannelName, ChannelSecretKey FROM raotdb.tblauthchannel  WHERE ChannelName = '" + Channel[0] + "'  ")
	defer ress2.Close()
	returnval := 0

	if err != nil {
		return returnval

	}

	for ress2.Next() {

		var event AuthDevice
		//JobID := ress2.Scan(&event.JobID)
		err = ress2.Scan(&event.Cid, &event.ChannelName, &event.ChannelSecretKey)

		if err != nil {
			panic(err)
		}
		//mobileid = event.MobileID
		//SendSMSFirstJob = event.SendSMSFirstJob

		// ress, err := db.Query("UPDATE  raotdb.tblAuthChannel Set UseCount = UseCount + 1, LastUseDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  WHERE ChannelName = '" + Channel[0] + "'  ")
		// defer ress.Close()
		// if err != nil {
		// 	panic(err)
		// } else {

		// }

		// ress, err := db.Query("UPDATE  raotdb.tblAuthChannel Set UseCount = UseCount + 1, LastUseDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  WHERE ChannelName = '" + Channel[0] + "'  ")
		// defer ress.Close()
		// if err != nil {
		// 	panic(err)
		// } else {

		// }

		return event.Cid
		// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

	}

	err = ress2.Close()

	return returnval
}
func PrintStuff(request *http.Request) {
	request.ParseForm()
	time.Sleep(time.Second * 3)
	log.Println(request.PostForm)
}
func AddHandler(respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json")
	go PrintStuff(request)
	respWriter.Write([]byte("Thanks\n"))
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to my web server!</h1>"))
}
func getToken(length int) string {
	token := ""
	//codeAlphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//codeAlphabet += "abcdefghijklmnopqrstuvwxyz"
	codeAlphabet := "0123456789"

	for i := 0; i < length; i++ {
		token += string(codeAlphabet[cryptoRandSecure(int64(len(codeAlphabet)))])
	}
	return token
}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}
func delegate(article KTBApproveJson, jsonstr string, w http.ResponseWriter, ch chan<- int) {
	// do some heavy calculations first
	// present the result (in the original code, the image)
	//fmt.Fprint(w, "hello")

	resp := make(map[string]string)
	resp["bankRef"] = article.BankRef

	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	bc, _ := json.Marshal(boxes1)

	w.Write(bc)

	ch <- 1

	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: tr}
	// _, err := client.Get("https://golang.org/")
	// if err != nil {
	// 	//fmt.Println(err)
	// }

	//Amount := fmt.Sprintf("%v", article.Amount)
	//BankCode := fmt.Sprintf("%v", article.BankCode)

	s := getToken(9)

	User := fmt.Sprintf("%v", article.User)
	Password := fmt.Sprintf("%v", article.Password)

	if User != "THPDKTB" {
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
		boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 101, RespMsg: "Invalid username/password", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}
	if Password != "!10ktb@Thpd!" {
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
		boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 101, RespMsg: "Invalid username/password", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}

	//db.Exec("INSERT INTO ListItems(items) Values(?)" , (r.Form["item"][0])) /*error*/
	if db == nil {
		//log.Println("DB is nill")

		dns := getDNSString(dblogin, userlogin, passlogin, conn)
		db2, err := sql.Open("mysql", dns)

		db = db2
		//db, err := sql.Open("mysql", "root:abc@/abc")
		if err != nil {
			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
		err = db.Ping()
		if err != nil {
			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction Ping", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
	}
	//time.Sleep(100 * time.Millisecond)
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// // err = db.Ping()
	// // if err != nil {
	// // 	panic(err)
	// // }
	// defer db.Close()

	orderid++
	//iorderid := int(orderid)
	//ia := strconv.Itoa(iorderid)

	//value, _ := rand.Int(rand.Reader, big.NewInt(100))
	//num := value
	//bigstr := value.String()

	//return string(b)
	t := time.Now()
	tUnix := t.Unix()
	//bb := strconv.FormatInt(tUnix, 10) + strconv.Itoa(aa)
	bb := strconv.FormatInt(tUnix, 10)

	bb = bb + s

	strtUnix := bb

	// ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + article.Ref2 + "' , '" + strings.ReplaceAll("BankKTBApprove", "'", `\'`) + "', '" + strings.ReplaceAll(string(jsonstr), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	// defer ress3.Close()
	// if err2 != nil {
	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction APILog", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// 	return
	// 	//panic(err)
	// }

	// haveitem := 0

	// Tkid := substr(article.Ref2, 0, 1)
	// price := "0"

	// if Tkid == "B" {
	// 	//jobid := substr(article.Ref2, 4, 10)
	// 	jobid := substr(article.Ref2, 0, 10)
	// 	ress2, err := db.Query("SELECT Jobid, TrackingID, THPDRatePrice, AmountCOD, FTLFinalPrice FROM thpdmpdb.vw_omschecksendprice  where Jobid = '" + jobid + "'  ")
	// 	defer ress2.Close()
	// 	if err == nil {

	// 		var event Checksendprice
	// 		//JobID := ress2.Scan(&event.JobID)
	// 		for ress2.Next() {

	// 			err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice)
	// 			price = event.FTLFinalPrice
	// 			haveitem = 1

	// 			break

	// 		}
	// 	}
	// } else {
	// 	//obid := substr(article.Ref2, 4, 13)
	// 	//jobid := substr(article.Ref2, 0, 13)
	// 	jobid := article.Ref2
	// 	ress2, err := db.Query("SELECT * FROM thpdmpdb.vw_omschecksendprice  where TrackingID = '" + jobid + "'  ")
	// 	defer ress2.Close()
	// 	if err == nil {

	// 		var event Checksendprice

	// 		for ress2.Next() {
	// 			err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice)
	// 			//price = event.THPDRatePrice
	// 			price = event.FTLFinalPrice
	// 			if price == "" {
	// 				price = event.THPDRatePrice
	// 			}

	// 			haveitem = 1

	// 			break
	// 		}
	// 	}
	// }

	// tt := int(article.Amount)
	// s2 := strconv.Itoa(tt)

	// fstr, err := strconv.ParseFloat(price, 32)
	// intVar := int(fstr)

	// if err != nil {
	// 	// tUnix := now.Unix()
	// 	// strtUnix := strconv.FormatInt(tUnix, 10)
	// 	// //fmt.Println(strtUnix)

	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
	// 	boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 108, RespMsg: "Invalid price or amount", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// 	return
	// }

	// price2 := strconv.Itoa(intVar)

	// time.Sleep(100 * time.Millisecond)

	// //now := time.Now()
	// ////fmt.Println("Today : ", now.Format(time.ANSIC))

	// // wrong way to convert nano to millisecond

	// ////fmt.Println(strtUnix + strconv.Itoa(aa))

	// if haveitem == 1 {

	// 	if s2 != price2 || price == "" {

	// 		resp := make(map[string]string)
	// 		resp["bankRef"] = article.BankRef

	// 		//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
	// 		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 108, RespMsg: "Invalid price or amount", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 		b, _ := json.Marshal(boxes1)

	// 		w.Write(b)
	// 		ch <- 1
	// 		return
	// 	}

	// 	ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblpaymentktbapprove ( tranxId, ComCode, ProdCode, Command, BankCode, BankRef, DateTime, EffDate, Amount, Channel, CusName, Ref1, Ref2, Ref3, Ref4,CreateDT) Values (" + strtUnix + ",'" + article.ComCode + "','" + article.ProdCode + "','" + article.Command + "','" + BankCode + "','" + article.BankRef + "','" + article.DateTime + "','" + article.EffDate + "','" + Amount + "','" + article.Channel + "','" + article.CusName + "','" + article.Ref1 + "','" + article.Ref2 + "','" + article.Ref3 + "','" + article.Ref4 + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))")
	// 	defer ress3.Close()
	// 	if err2 != nil {

	// 		// defer db.Close()
	// 		// connectDb()
	// 		// time.Sleep(100 * time.Millisecond)

	// 		resp := make(map[string]string)
	// 		resp["bankRef"] = article.BankRef

	// 		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 		b, _ := json.Marshal(boxes1)

	// 		w.Write(b)
	// 		ch <- 1
	// 		return
	// 		//panic(err)
	// 	}
	// 	//defer ress3.Close()
	// 	err = ress3.Close()

	// 	//boxes := []KTBReponseApproveJson{}
	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
	// 	boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// } else {
	// 	//boxes := []KTBReponseApproveJson{}

	// 	// tUnix := t.Unix()
	// 	// strtUnix := strconv.FormatInt(tUnix, 10)
	// 	// //fmt.Println(strtUnix)

	// 	ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblpaymentktbapprove ( tranxId, ComCode, ProdCode, Command, BankCode, BankRef, DateTime, EffDate, Amount, Channel, CusName, Ref1, Ref2, Ref3, Ref4,CreateDT) Values (" + strtUnix + ",'" + article.ComCode + "','" + article.ProdCode + "','Not Approve','" + BankCode + "','" + article.BankRef + "','" + article.DateTime + "','" + article.EffDate + "','" + Amount + "','" + article.Channel + "','" + article.CusName + "','" + article.Ref1 + "','" + article.Ref2 + "','" + article.Ref3 + "','" + article.Ref4 + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))")
	// 	defer ress3.Close()
	// 	if err2 != nil {

	// 		resp := make(map[string]string)
	// 		resp["bankRef"] = article.BankRef

	// 		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 		b, _ := json.Marshal(boxes1)

	// 		w.Write(b)
	// 		ch <- 1
	// 		return
	// 		//panic(err)
	// 	}
	// 	//defer ress3.Close()
	// 	err = ress3.Close()

	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 104, RespMsg: "Invalid reference", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// }
	resp = make(map[string]string)
	resp["bankRef"] = article.BankRef

	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	boxes1 = KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	bc, _ = json.Marshal(boxes1)

	w.Write(bc)

	ch <- 1
}
func delegatetrue(article KTBApproveJson, jsonstr string, w http.ResponseWriter, ch chan<- int) {
	// do some heavy calculations first
	// present the result (in the original code, the image)
	//fmt.Fprint(w, "hello")

	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: tr}
	// _, err := client.Get("https://golang.org/")
	// if err != nil {
	// 	////fmt.Println(err)
	// }

	s := getToken(9)

	//reqBody, _ := ioutil.ReadAll(r.Body)
	// var article KTBApproveJson
	// json.Unmarshal(reqBody, &article)

	//KTBReponseApproveJson
	// insert payment request ktb
	//KTBApproveJson

	//jsonResp, err := json.Marshal(b)
	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	return
	// }

	//SendOTPToCustomer(strOTP, event.PhoneNumber)

	//var s string = strconv.FormatFloat(f, 'E', -1, 32)
	Amount := fmt.Sprintf("%v", article.Amount)
	BankCode := fmt.Sprintf("%v", article.BankCode)

	User := fmt.Sprintf("%v", article.User)
	Password := fmt.Sprintf("%v", article.Password)

	if User != "THPDKTB" {
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
		boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 101, RespMsg: "Invalid username/password", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}
	if Password != "!10ktb@Thpd!" {
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
		boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 101, RespMsg: "Invalid username/password", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}

	// if db == nil {
	// 	//log.Println("DB is nill")
	// 	connectDb()

	// 	time.Sleep(200 * time.Millisecond)
	// 	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// 	// db2, err := sql.Open("mysql", dns)

	// 	// db = db2
	// 	// //db, err := sql.Open("mysql", "root:abc@/abc")
	// 	// if err != nil {

	// 	// resp := make(map[string]string)
	// 	// resp["bankRef"] = article.BankRef

	// 	// //boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	// boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction DB Closed", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	// b, _ := json.Marshal(boxes1)

	// 	// w.Write(b)
	// 	// ch <- 1
	// 	// return

	// }
	if db == nil {
		//log.Println("DB is nill")

		dns := getDNSString(dblogin, userlogin, passlogin, conn)
		db2, err := sql.Open("mysql", dns)

		db = db2
		//db, err := sql.Open("mysql", "root:abc@/abc")
		if err != nil {
			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
		err = db.Ping()
		if err != nil {
			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction Ping", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
	}
	//time.Sleep(100 * time.Millisecond)
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// defer db.Close()
	// if err != nil {
	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// 	return
	// 	//panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction Ping", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// 	return
	// 	//panic(err)
	// }

	//jsonstr := string(reqBody)
	////fmt.Println(jsonstr)
	//time.Sleep(100 * time.Millisecond)

	// ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + article.Ref2 + "' , '" + strings.ReplaceAll("BankKTBApprove", "'", `\'`) + "', '" + strings.ReplaceAll(string(jsonstr), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	// defer ress3.Close()
	// if err2 != nil {

	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	ch <- 1
	// 	return
	// 	//panic(err)
	// }

	// ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('GetKTBApprove','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	// defer ress99.Close()
	// if err != nil {
	// 	resp := make(map[string]string)
	// 	resp["bankRef"] = article.BankRef

	// 	//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
	// 	boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

	// 	b, _ := json.Marshal(boxes1)

	// 	w.Write(b)
	// 	return
	// 	//panic(err)
	// }

	haveitem := 0

	t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//Tkid := substr(article.Ref2, 4, 1)
	Tkid := substr(article.Ref2, 0, 1)
	price := "0"
	BankRefCharge := ""

	if Tkid == "B" {
		//jobid := substr(article.Ref2, 4, 10)
		jobid := substr(article.Ref2, 0, 10)
		ress2, err := db.Query("SELECT Jobid, TrackingID, THPDRatePrice, AmountCOD, FTLFinalPrice, BankRef FROM thpdmpdb.vw_omschecksendprice  where Jobid = '" + jobid + "'  ")
		defer ress2.Close()
		if err == nil {

			var event Checksendprice
			//JobID := ress2.Scan(&event.JobID)
			for ress2.Next() {

				err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice, &event.BankRef)
				price = event.FTLFinalPrice
				BankRefCharge = event.BankRef
				haveitem = 1

				break

			}
		}
	} else {
		//obid := substr(article.Ref2, 4, 13)
		//jobid := substr(article.Ref2, 0, 13)
		jobid := article.Ref2
		ress2, err := db.Query("SELECT * FROM thpdmpdb.vw_omschecksendprice  where TrackingID = '" + jobid + "'  ")
		defer ress2.Close()
		if err == nil {

			var event Checksendprice

			for ress2.Next() {
				err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice, &event.BankRef)
				//price = event.THPDRatePrice
				price = event.FTLFinalPrice
				BankRefCharge = event.BankRef
				if price == "" {
					price = event.THPDRatePrice
				}

				haveitem = 1

				break
			}
		}
	}

	if BankRefCharge != "" {

		// ชำระเงินแล้ว
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
		boxes1 := KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 106, RespMsg: "Transaction number duplicate", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}

	tt := int(article.Amount)
	s2 := strconv.Itoa(tt)

	fstr, err := strconv.ParseFloat(price, 32)
	intVar := int(fstr)
	if err != nil {
		tUnix := t.Unix()
		strtUnix := strconv.FormatInt(tUnix, 10)
		//fmt.Println(strtUnix)

		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 108, RespMsg: "Invalid price or amount", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
		return
	}

	price2 := strconv.Itoa(intVar)

	//now := time.Now()

	//fmt.Println("Today : ", t.Format(time.ANSIC))

	// wrong way to convert nano to millisecond

	//time.Sleep(100 * time.Millisecond)
	//nano := t.Nanosecond()
	//aa := nano / 1000000
	//aa := nano / 100000
	//aa := nano

	// orderid++
	// iorderid := int(orderid)
	// ia := strconv.Itoa(iorderid)

	// if len(ia) == 1 {
	// 	ia = "0000000" + ia
	// } else if len(ia) == 2 {
	// 	ia = "000000" + ia
	// } else if len(ia) == 3 {
	// 	ia = "00000" + ia
	// } else if len(ia) == 4 {
	// 	ia = "0000" + ia
	// } else if len(ia) == 5 {
	// 	ia = "000" + ia
	// } else if len(ia) == 6 {
	// 	ia = "00" + ia
	// } else if len(ia) == 7 {
	// 	ia = "0" + ia
	// }

	//value, _ := rand.Int(rand.Reader, big.NewInt(100))
	//num := value
	//bigstr := value.String()

	//return string(b)

	tUnix := t.Unix()
	//bb := strconv.FormatInt(tUnix, 10) + strconv.Itoa(aa)
	bb := strconv.FormatInt(tUnix, 10)

	bb = bb + s

	strtUnix := bb

	//jobid := article.Ref2
	// ress22, err := db.Query("SELECT * FROM thpdmpdb.tblpaymentktbcharge  where Ref2 = '" + jobid + "'  ")
	// defer ress22.Close()
	// if err == nil {

	// 	for ress22.Next() {
	// 		/// ชำระเงินแล้ว
	// 		resp := make(map[string]string)
	// 		resp["bankRef"] = article.BankRef

	// 		//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
	// 		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 106, RespMsg: "Transaction number duplicate", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"}

	// 		b, _ := json.Marshal(boxes1)

	// 		w.Write(b)
	// 		ch <- 1
	// 		return

	// 	}
	// }

	if haveitem == 1 {

		if s2 != price2 || price == "" {

			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
			boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 108, RespMsg: "Invalid price or amount", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
		}

		ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblpaymentktbapprove ( tranxId, ComCode, ProdCode, Command, BankCode, BankRef, DateTime, EffDate, Amount, Channel, CusName, Ref1, Ref2, Ref3, Ref4,CreateDT) Values (" + strtUnix + ",'" + article.ComCode + "','" + article.ProdCode + "','" + article.Command + "','" + BankCode + "','" + article.BankRef + "','" + article.DateTime + "','" + article.EffDate + "','" + Amount + "','" + article.Channel + "','" + article.CusName + "','" + article.Ref1 + "','" + article.Ref2 + "','" + article.Ref3 + "','" + article.Ref4 + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))")
		defer ress3.Close()
		if err2 != nil {

			// defer db.Close()
			// connectDb()
			// time.Sleep(100 * time.Millisecond)

			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
		//defer ress3.Close()
		err = ress3.Close()

		//boxes := []KTBReponseApproveJson{}
		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
	} else {
		//boxes := []KTBReponseApproveJson{}

		// tUnix := t.Unix()
		// strtUnix := strconv.FormatInt(tUnix, 10)
		// //fmt.Println(strtUnix)

		ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblpaymentktbapprove ( tranxId, ComCode, ProdCode, Command, BankCode, BankRef, DateTime, EffDate, Amount, Channel, CusName, Ref1, Ref2, Ref3, Ref4,CreateDT) Values (" + strtUnix + ",'" + article.ComCode + "','" + article.ProdCode + "','Not Approve','" + BankCode + "','" + article.BankRef + "','" + article.DateTime + "','" + article.EffDate + "','" + Amount + "','" + article.Channel + "','" + article.CusName + "','" + article.Ref1 + "','" + article.Ref2 + "','" + article.Ref3 + "','" + article.Ref4 + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))")
		defer ress3.Close()
		if err2 != nil {

			resp := make(map[string]string)
			resp["bankRef"] = article.BankRef

			//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
			boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 100, RespMsg: "Unable to process transaction", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

			b, _ := json.Marshal(boxes1)

			w.Write(b)
			ch <- 1
			return
			//panic(err)
		}
		//defer ress3.Close()
		err = ress3.Close()

		resp := make(map[string]string)
		resp["bankRef"] = article.BankRef

		//boxes = append(boxes, KTBReponseApproveJson{TranxID: "", BankRef: article.BankRef, RespCode: 0, RespMsg: "Not Success", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""})
		boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 104, RespMsg: "Invalid reference", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: ""}

		b, _ := json.Marshal(boxes1)

		w.Write(b)
		ch <- 1
	}

	//fmt.Fprint(w, b)

	// resp := make(map[string]string)
	// resp["bankRef"] = article.BankRef

	// //boxes = append(boxes, KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"})
	// boxes1 := KTBReponseApproveJson{TranxID: strtUnix, BankRef: article.BankRef, RespCode: 0, RespMsg: "Successful", Balance: article.Amount, CusName: article.CusName, Info: "", Print1: "ไปรษณ๊ย์ไทยดิสทริบ้วชั่น ขอขอบพระคุณ"}

	// b, _ := json.Marshal(boxes1)

	// w.Write(b)
	//ch <- 1
	//ch <- 1
}

func connectDb() {
	//socket : var/run/mysqld/mysqld.sock
	/* connection string examples :
	db, err := sql.Open("mysql", "user:password@/dbname")
	user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true

	TCP using default port (3306) on localhost:
	  user:password@tcp/dbname?charset=utf8mb4,utf8&sys_var=esc%40ped

	Use the default protocol (tcp) and host (localhost:3306):
	  user:password@/dbname

	No Database preselected:
	  user:password@/
	*/

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db2, err := sql.Open("mysql", dns)

	DB = db2
	//db, err := sql.Open("mysql", "root:abc@/abc")
	if err != nil {
		return
		//panic(err)
	}
	err = DB.Ping()
	if err != nil {
		return
		//panic(err)
	}

	// ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('DBConnect','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println("DB: ", db)

	checkErr(err)

	// Open doesn't open a connection. Validate DSN data:

}
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	} else {
		log.Println(err)
	}
}

// func GetKTBApprove(w http.ResponseWriter, r *http.Request) {

// 	//time.Sleep(100 * time.Millisecond)

// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	var article KTBApproveJson
// 	json.Unmarshal(reqBody, &article)

// 	jsonstr := string(reqBody)
// 	////fmt.Println(jsonstr)

// 	// tr := &http.Transport{
// 	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	// }
// 	// client := &http.Client{Transport: tr}
// 	// _, err := client.Get("https://golang.org/")
// 	// if err != nil {
// 	// 	//fmt.Println(err)
// 	// }

// 	//go f(article, jsonstr, w, r)

// 	ch := make(chan int)
// 	go delegatetrue(article, jsonstr, w, ch)
// 	<-ch

// }

func DoneAsync() chan int {
	r := make(chan int)
	//fmt.Println("Warming up ...")
	go func() {
		time.Sleep(3 * time.Second)
		r <- 1
		//fmt.Println("Done ...")
	}()
	return r
}
func DoneAsync2() chan int {
	r := make(chan int)
	//fmt.Println("Warming up ...")
	go func() {
		time.Sleep(3 * time.Second)
		r <- 1
		//fmt.Println("Done ...")
	}()
	return r
}

func Sign(Data map[string]interface{}, SecrePublicKeyEnvName string, ExpiredAt time.Duration) (string, error) {

	expiredAt := time.Now().Add(time.Duration(time.Second) * ExpiredAt).Unix()
	//sexpiredAt := strconv.FormatInt(expiredAt, 10)
	//jwtSecretKey := GodotEnv(SecrePublicKeyEnvName)
	jwtSecretKey := SecrePublicKeyEnvName
	// metadata for your jwt
	claims := jwt.MapClaims{}
	claims["expiredAt"] = expiredAt
	// claims["authorization"] = true

	for i, v := range Data {
		claims[i] = v
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(jwtSecretKey))

	fmt.Print(accessToken)

	if err != nil {
		//logrus.Error(err.Error())
		return accessToken, err
	}

	//value := claims["USERID"]
	//strvalue := fmt.Sprintf("%v", value)

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// ress, err := db.Query("UPDATE raotdb.tbluser SET token = '" + accessToken + "', tokenExpiredt = '" + string(sexpiredAt) + "' WHERE username = '" + strvalue + "'")
	// defer ress.Close()
	// if err != nil {
	// 	panic(err)
	// }

	return accessToken, nil
}
func SignRAOT(Data map[string]interface{}, SecrePublicKeyEnvName string, ExpiredAt time.Duration) (string, error) {

	expiredAt := time.Now().Add(time.Duration(time.Second) * ExpiredAt).Unix()
	sexpiredAt := strconv.FormatInt(expiredAt, 10)
	//jwtSecretKey := GodotEnv(SecrePublicKeyEnvName)
	jwtSecretKey := SecrePublicKeyEnvName
	// metadata for your jwt
	claims := jwt.MapClaims{}
	claims["expiredAt"] = expiredAt
	// claims["authorization"] = true

	for i, v := range Data {
		claims[i] = v
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(jwtSecretKey))

	fmt.Print(accessToken)

	if err != nil {
		//logrus.Error(err.Error())
		return accessToken, err
	}

	value := claims["USERID"]
	strvalue := fmt.Sprintf("%v", value)

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ress, err := db.Query("UPDATE raotdb.tbluser SET token = '" + accessToken + "', tokenExpiredt = '" + string(sexpiredAt) + "',logincounter = logincounter + 1  WHERE username = '" + strvalue + "' and tokenExpiredt <> '" + string(sexpiredAt) + "' ")
	defer ress.Close()
	if err != nil {
		panic(err)
	}

	return accessToken, nil
}

func SendOTPToCustomer(OTPCode string, CustomerPhone string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)

	var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	var emp THSMS
	emp.Sender = "THPD-OM-OTP" //SMSOTP
	emp.Msisdn = []string{CustomerPhone}
	emp.Message = "หมายเลข OTP : " + OTPCode + "\n   "

	var data []byte
	// Convert struct to jsonข
	data, _ = json.MarshalIndent(emp, "", "    ")

	url := "https://thsms.com/api/send-sms"
	//req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println(req.Header)

	//req.Body.Read(json_data)
	//Send req using http Client
	client := &http.Client{}
	resp2, err := client.Do(req)

	if resp2 != nil {
		//panic(err)
	}
	if err == nil {
		/// update tracking send message already
		// ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET SendSMSFirstJob = 1 WHERE TrackingID = '" + trackingID + "'")
		// defer ress.Close()
		// if err != nil {

		// 	panic(err)
		// } else {
		// 	ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'ผู้ขนส่งรับงานแล้ว' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + trackingID + "')")
		// 	defer ress2.Close()
		// 	if err != nil {
		// 		panic(err)
		// 	} else {

		// 	}

		// }
		//panic(err)
	}
	return string("OK")
}
func SendMessageToDriver(trackingID string, CustomerPhone string, CarID string, DriverName string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)

	var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	var emp THSMS
	emp.Sender = "PromptSong"
	emp.Msisdn = []string{CustomerPhone}
	emp.Message = "OMS DE : เลข Track: " + trackingID + "\n  ผู้ขนส่ง: " + DriverName + "\n  รหัสผู้ขนส่ง: " + CarID + "\n ลูกค้าได้เลือกใช้บริการท่านแล้ว \n กรุณาเข้ารับให้ตรงตามวันรับสินค้า!"

	var data []byte
	// Convert struct to json
	data, _ = json.MarshalIndent(emp, "", "    ")

	url := "https://thsms.com/api/send-sms"
	//req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println(req.Header)

	//req.Body.Read(json_data)
	//Send req using http Client
	client := &http.Client{}
	resp2, err := client.Do(req)

	if resp2 != nil {
		//panic(err)
	}
	if err == nil {
		/// update tracking send message already
		// ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET SendSMSFirstJob = 1 WHERE TrackingID = '" + trackingID + "'")
		// defer ress.Close()
		// if err != nil {

		// 	panic(err)
		// } else {
		// 	ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'ผู้ขนส่งรับงานแล้ว' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + trackingID + "')")
		// 	defer ress2.Close()
		// 	if err != nil {
		// 		panic(err)
		// 	} else {

		// 	}

		// }
		//panic(err)
	}
	return string("OK")
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func OMSRepSummaryReport(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MerchantRep
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	//defer db.Close()
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	MerchantID := article.MerchantID
	RepType := article.RepType
	SDT := article.SDT
	EDT := article.EDT

	if err == nil {

		//s := trackingID[:len(trackingID)-1]

		//chars := []rune(trackingID)
		//strArray := strings.Fields(MsgText)
		//var tx =
		//s := string(chars[0])
		//s2 := string(chars[1])
		boxes := []RepSummary{}
		if MerchantID != "" {
			ress2, err := db.Query("call SP_ShopReportSummary ('" + MerchantID + "','" + RepType + "','" + SDT + "','" + EDT + "')")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event RepSummary
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.UserID, &event.Status_Name, &event.CStatus_Name)

					if err != nil {
						panic(err)
					}
					// PMerchantID := event.PMerchantID
					// Ftypes := event.Ftypes
					// SumAmount := event.SumAmount

					boxes = append(boxes, RepSummary{UserID: event.UserID, Status_Name: event.Status_Name, CStatus_Name: event.CStatus_Name})

				}

			}
			defer ress2.Close()
			err = ress2.Close()

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//defer ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			return
		} else {

			respok := make(map[string]string)
			respok["Success"] = "false"
			jsonResp, err := json.Marshal(respok)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}
}
func OMSRepSummaryReportWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MerchantRep
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	// ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSRepSummaryReportWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	MerchantID := article.MerchantID
	RepType := article.RepType
	SDT := article.SDT
	EDT := article.EDT

	if err == nil {

		//s := trackingID[:len(trackingID)-1]

		//chars := []rune(trackingID)
		//strArray := strings.Fields(MsgText)
		//var tx =
		//s := string(chars[0])
		//s2 := string(chars[1])
		boxes := []RepSummary{}
		if MerchantID != "" {

			sp := "SP_ShopReportSummary"
			if RepType == "Daily" {
				sp = "SP_ShopReportSummary"
				// SDT = SDT + " 00:00:00"
				// EDT = EDT + " 23:59:00"
			}
			if RepType == "Monthly" {
				sp = "SP_ShopReportSummaryMonth"
			}
			if RepType == "Yearly" {
				sp = "SP_ShopReportSummaryYear"
			}
			//fmt.Println(sp)

			ress2, err := db.Query("call " + sp + " ('" + MerchantID + "','" + RepType + "','" + SDT + "','" + EDT + "')")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event RepSummary
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.UserID, &event.Status_Name, &event.CStatus_Name)

					if err != nil {
						panic(err)
					}
					// PMerchantID := event.PMerchantID
					// Ftypes := event.Ftypes
					// SumAmount := event.SumAmount

					boxes = append(boxes, RepSummary{UserID: event.UserID, Status_Name: event.Status_Name, CStatus_Name: event.CStatus_Name})

				}

			}
			defer ress2.Close()
			err = ress2.Close()

			b, _ := json.Marshal(boxes)

			// defer ress2.Close()
			// err = ress2.Close()
			//defer ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			return
		} else {

			respok := make(map[string]string)
			respok["Success"] = "false"
			jsonResp, err := json.Marshal(respok)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}
}
func OMSRepPayment(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	MerchantID := article.MerchantID

	if err == nil {

		//s := trackingID[:len(trackingID)-1]

		//chars := []rune(trackingID)
		//strArray := strings.Fields(MsgText)
		//var tx =
		//s := string(chars[0])
		//s2 := string(chars[1])
		boxes := []RepPayment{}
		if MerchantID != "" {
			ress2, err := db.Query("call sp_mobilereportpayment ('" + MerchantID + "')")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event RepPayment
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.PMerchantID, &event.Ftypes, &event.SumAmount)

					if err != nil {
						panic(err)
					}
					// PMerchantID := event.PMerchantID
					// Ftypes := event.Ftypes
					// SumAmount := event.SumAmount

					boxes = append(boxes, RepPayment{PMerchantID: event.PMerchantID, Ftypes: event.Ftypes, SumAmount: event.SumAmount})

				}

			}
			defer ress2.Close()
			err = ress2.Close()

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//defer ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			return
		} else {

			respok := make(map[string]string)
			respok["Success"] = "false"
			jsonResp, err := json.Marshal(respok)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}
}
func OMSRepPaymentWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSRepPaymentWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	MerchantID := article.MerchantID

	if err == nil {

		//s := trackingID[:len(trackingID)-1]

		//chars := []rune(trackingID)
		//strArray := strings.Fields(MsgText)
		//var tx =
		//s := string(chars[0])
		//s2 := string(chars[1])
		boxes := []RepPayment{}
		if MerchantID != "" {
			ress2, err := db.Query("call sp_mobilereportpayment ('" + MerchantID + "')")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event RepPayment
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.PMerchantID, &event.Ftypes, &event.SumAmount)

					if err != nil {
						panic(err)
					}
					// PMerchantID := event.PMerchantID
					// Ftypes := event.Ftypes
					// SumAmount := event.SumAmount

					boxes = append(boxes, RepPayment{PMerchantID: event.PMerchantID, Ftypes: event.Ftypes, SumAmount: event.SumAmount})

				}

			}
			defer ress2.Close()
			err = ress2.Close()

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//defer ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			return
		} else {

			respok := make(map[string]string)
			respok["Success"] = "false"
			jsonResp, err := json.Marshal(respok)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}
}
func OMSCheckPrice(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSCheckPrice','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	trackingID := article.TrackingID
	Jobid := ""
	TrackingID := ""
	THPDRatePrice := ""
	FTLFinalPrice := ""
	AmountCOD := ""

	if err == nil {

		//s := trackingID[:len(trackingID)-1]

		chars := []rune(trackingID)
		//strArray := strings.Fields(MsgText)
		//var tx =
		s := string(chars[0])
		s2 := string(chars[1])

		if s == "B" {
			ress2, err := db.Query("SELECT Jobid,'' as TrackingID , SUM(THPDRatePrice)THPDRatePrice, SUM(AmountCOD)AmountCOD, SUM(FTLFinalPrice)FTLFinalPrice  FROM THPDMPDB.vw_omschecksendprice Where jobid = '" + trackingID + "' GROUP BY Jobid,TrackingID")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event ChecksendPrice
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice)

					if err != nil {
						panic(err)
					}
					Jobid = event.Jobid
					TrackingID = event.TrackingID
					THPDRatePrice = event.THPDRatePrice
					AmountCOD = event.AmountCOD
					FTLFinalPrice = event.FTLFinalPrice

					// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

				}

			}
			defer ress2.Close()
			err = ress2.Close()
		} else if (s == "F" || s == "L" || s == "2") && s2 == "T" {
			ress2, err := db.Query("SELECT Jobid, TrackingID, THPDRatePrice, AmountCOD, FTLFinalPrice FROM THPDMPDB.vw_omschecksendprice Where TrackingID = '" + trackingID + "'")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event ChecksendPrice
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice)

					if err != nil {
						panic(err)
					}
					Jobid = event.Jobid
					TrackingID = event.TrackingID
					THPDRatePrice = event.THPDRatePrice
					AmountCOD = event.AmountCOD
					FTLFinalPrice = event.FTLFinalPrice

					// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

				}

			}
			defer ress2.Close()
			err = ress2.Close()
		} else if (s == "G" || s == "L") && s2 != "T" {
			/// Logispost
			ress2, err := db.Query("SELECT Jobid, TrackingID, THPDRatePrice, AmountCOD, FTLFinalPrice FROM THPDMPDB.vw_omschecksendprice  ")
			if err == nil {

				// boxes := []Merchant{}
				// resp := make(map[string]string)
				// resp["MerchantID"] = article.MerchantID

				for ress2.Next() {
					var event ChecksendPrice
					//JobID := ress2.Scan(&event.JobID)
					err = ress2.Scan(&event.Jobid, &event.TrackingID, &event.THPDRatePrice, &event.AmountCOD, &event.FTLFinalPrice)

					if err != nil {
						panic(err)
					}
					Jobid = event.Jobid
					TrackingID = event.TrackingID
					THPDRatePrice = event.THPDRatePrice
					AmountCOD = event.AmountCOD
					FTLFinalPrice = event.FTLFinalPrice

					// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

				}

			}
			defer ress2.Close()
			err = ress2.Close()
		}

		respok := make(map[string]string)
		respok["Success"] = "true"
		respok["Jobid"] = Jobid
		respok["TrackingID"] = TrackingID
		respok["THPDRatePrice"] = THPDRatePrice
		respok["AmountCOD"] = AmountCOD
		respok["FTLFinalPrice"] = FTLFinalPrice
		jsonResp, err := json.Marshal(respok)

		/// update driver get job
		//GetAPIDriverFromLoadboard(article.TrackingID)
		/// Update jobdetail

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	} else {

		respok := make(map[string]string)
		respok["Success"] = "false"
		jsonResp, err := json.Marshal(respok)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

}

func ChkAppVersion(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article MobileVersion
	json.Unmarshal(reqBody, &article)

	//fmt.Println("Body", article.VersionNow)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	MobileType := article.MobileType
	VersionNow := article.VersionNow

	ress3, err2 := db.Query("SELECT AppType, Version FROM THPDMPDB.tblmobileappversion WHERE AppType = '" + MobileType + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []Appversion{}

		for ress3.Next() {
			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.AppType, &event.Version)

			if err != nil {
				panic(err)
			}

			if event.Version == VersionNow {
				respok := make(map[string]string)
				respok["AppType"] = event.AppType
				respok["Version"] = "lastest" //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return

			}

			boxes = append(boxes, Appversion{AppType: event.AppType, Version: event.Version})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}
}
func GetWeight(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetWeight',requestcount+1,'" + channel[0] + "','" + article.WHName + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	//////

	WHName := article.WHName
	//Weight := article.Weight
	CTime := article.CTime

	if CTime == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte("PWASMART0805401073")
	t = jwt.New(jwt.SigningMethodHS256)
	s, err = t.SignedString(key)
	//fmt.Println(s)

	ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblcurrentweight WHERE Wid in (SELECT max(Wid) as Appid FROM raotdb.tblcurrentweight WHERE WHName = '" + WHName + "' and CShaftNum = '00') ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}

		for ress3.Next() {
			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT)
			if err != nil {
				panic(err)
			}
			if event.Version == WHName {
				respok := make(map[string]string)
				respok["WHName"] = WHName
				respok["Token"] = s //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return
			}

			boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Version, LastestDT: event.LastestDT, Token: s})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Content-Type", "text/plain")
		// w.WriteHeader(http.StatusOK)

		// if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
		// 	w.Header().Set("Access-Control-Allow-Headers", acrh[0])
		// }
		// w.Header().Set("Access-Control-Allow-Credentials", "True")
		// if acao, ok := r.Header["Access-Control-Allow-Origin"]; ok {
		// 	w.Header().Set("Access-Control-Allow-Origin", acao[0])
		// } else {
		// 	if _, oko := r.Header["Origin"]; oko {
		// 		w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
		// 	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}

func UpdateTransDetail(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASaveWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("TranSubID", article.TranSubID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	Token := article.Token

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('UpdateTransDetail',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	//CTime := article.CTime
	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}
	TransportID := ""
	ress3, err2 := db.Query("SELECT TransportSPID, TransportID, CustomName, CusAddress, ifnull(CusPhone,'')CusPhone, ContainerType, ContainerNo, ContainerSideNo,TruckTypeID, TruckLicenseNo, TruckLicenseCountry, TruckLicenseTrail, TruckLicenseTrailCountry,ToCustomDetail, CustomID, TransferDT, NetWeight,GrossWeight,ShippingName, RubberType, ifnull(GrossWeightOnSite,'')GrossWeightOnSite,ifnull(EstWeightOnSite,'')EstWeightOnSite, ifnull(Calctype,'')Calctype, ifnull(TimeOnSite,'')TimeOnSite, ifnull(Transubstatus,'')Transubstatus, ifnull(weightID,'')weightID  FROM raotdb.tbltransportsubdetail WHERE TransportSPID  = '" + article.TranSubID + "'  ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppCheckPointTranDetailversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransportSPID, &event.TransportID, &event.CustomName, &event.CusAddress, &event.CusPhone, &event.ContainerType, &event.ContainerNo, &event.ContainerSideNo, &event.TruckTypeID, &event.TruckLicenseNo, &event.TruckLicenseCountry, &event.TruckLicenseTrail, &event.TruckLicenseTrailCountry, &event.ToCustomDetail, &event.CustomID, &event.TransferDT, &event.NetWeight, &event.GrossWeight, &event.ShippingName, &event.RubberType, &event.GrossWeightOnSite, &event.EstWeightOnSite, &event.Calctype, &event.TimeOnSite, &event.Transubstatus, &event.WeightID)
			if err != nil {
				panic(err)
			}
			TransportID = string(event.TransportID)
		}
	}

	ress, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET grossWeightOnSite = '" + article.Weight + "', estWeightOnSite = '" + article.WeightCalcEst + "', weightID = '" + article.WeightID + "', transubstatus = '" + article.EType + "', timeOnSite = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') WHERE transportSPID = '" + article.TranSubID + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	} else {

		ress99, err := db.Query("INSERT INTO raotdb.tblimageupload (transportSPID, imgBase64, uploadDT) Values ('" + article.TranSubID + "','" + article.IMGBase64 + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE imgBase64 ='" + article.IMGBase64 + "'   ")
		defer ress99.Close()
		if err != nil {
			panic(err)
		}

		ress999, err := db.Query("call  raotdb.sp_updatestatustransport ('" + TransportID + "')  ")
		defer ress999.Close()
		if err != nil {
			panic(err)
		}

		boxes := []Appversion{}
		boxes = append(boxes, Appversion{Appid: article.TranSubID, AppType: article.Weight})
		b, _ := json.Marshal(boxes)

		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)
	}

	////fmt.Println("UserID", articleimg.ParsedText)

	////fmt.Println("UserID", article.WHName)

	// jsonResp, err := json.Marshal(string(body))
	// //fmt.Println(jsonResp)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")

}

func GetCalculate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASaveWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.TranSubID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	//CTime := article.TranSubID

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCalculate',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}
	/// ตรวจสอบจาก ข้อมูล trader ฝั่งพี่วัต ว่ามีรถใหม่หรือไม่ ถ้ามีใหม่ให้ insert

	//GetTruck

	// ress99, err := db.Query("call raotdb.sp_inserttradertrucktocalc")
	// defer ress99.Close()
	// if err != nil {
	// 	//panic(err)
	// }
	////

	strsql := "SELECT calcID, rubberType, containerType,truckWeight,trailWeight,tankweight,containerWeight,boxWeight,totalcontainer,totalbox,ifnull(containerTypeID,'1')containerTypeID  FROM raotdb.tblcalc "

	if article.SearchTxt != "" {
		strsql = strsql + " WHERE CONCAT(rubberType, containerType,truckWeight,trailWeight,tankweight,containerWeight,boxWeight,totalcontainer,totalbox ) like '%" + article.SearchTxt + "%'"
	}

	ress3, err2 := db.Query(strsql)
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppReturnCalc{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppGetCalc
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.calcID, &event.rubberType, &event.containerType, &event.truckWeight, &event.trailWeight, &event.tankweight, &event.containerWeight, &event.boxWeight, &event.totalcontainer, &event.totalbox, &event.containerTypeID)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnCalc{Id: cnt, CalcID: event.calcID, RubberType: event.rubberType, ContainerType: event.containerType, TruckWeight: event.truckWeight, TrailWeight: event.trailWeight, Tankweight: event.tankweight, ContainerWeight: event.containerWeight, BoxWeight: event.boxWeight, Totalcontainer: event.totalcontainer, Totalbox: event.totalbox, ContainerTypeID: event.containerTypeID})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func GetUpdateOverLicense(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASaveLicense
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.TranSubID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	// CTime := article.TranSubID
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetUpdateOverLicense',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET truckLPRlicenseOverwrite = '" + article.TruckLicenseID + "', truckLPRlicenseCountryOverwrite = '" + article.TruckDistinct + "', truckLPRTypeOverwrite = '" + article.TruckType + "', truckLPRTypeTruckOverwrite = '" + article.TruckTypeDetail + "', modifyDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') WHERE transportSPID = '" + article.TranSubID + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	}

	ress3, err2 := db.Query("SELECT calcID, rubberType, containerType,truckWeight,trailWeight,tankweight,containerWeight,boxWeight,totalcontainer,totalbox  FROM raotdb.tblcalc ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppReturnCalc{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppGetCalc
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.calcID, &event.rubberType, &event.containerType, &event.truckWeight, &event.trailWeight, &event.tankweight, &event.containerWeight, &event.boxWeight, &event.totalcontainer, &event.totalbox)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnCalc{Id: cnt, CalcID: event.calcID, RubberType: event.rubberType, ContainerType: event.containerType, TruckWeight: event.truckWeight, TrailWeight: event.trailWeight, Tankweight: event.tankweight, ContainerWeight: event.containerWeight, BoxWeight: event.boxWeight, Totalcontainer: event.totalcontainer, Totalbox: event.totalbox})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}

func GetIMGBase64(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASaveWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.TranSubID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetIMGBase64',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	CTime := article.TranSubID
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT TransportSPID, imgBase64, uploadDT  FROM raotdb.tblimageupload WHERE  transportSPID = '" + CTime + "'  ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppReturnIMG{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppGetIMG
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransportSPID, &event.imgBase64, &event.uploadDT)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnIMG{Id: cnt, TransportSPID: event.TransportSPID, IMGBase64: event.imgBase64, UploadDT: event.uploadDT})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}

func GetLicenseIMGBase64LPR(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	CTime := article.CTime
	img64 := strings.Split(CTime, ",")[1]
	Token := article.Token

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetLicenseIMGBase64LPR',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	// //fmt.Println(CTime)
	//fmt.Println(img64)

	if CTime == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	url := "https://cess-api.olive.co.th/api/v1/public/scan-license-plate/"
	method := "POST"

	payload := strings.NewReader(`{"image": "` + img64 + `"}`)

	// payload := &bytes.Buffer{}
	// writer := multipart.NewWriter(payload)
	// // _ = writer.WriteField("language", "eng")
	// // _ = writer.WriteField("isOverlayRequired", "true")
	// _ = writer.WriteField("image", img64)
	// // _ = writer.WriteField("OCREngine", "2")
	// err = writer.Close()

	// if err != nil {
	// 	//fmt.Println(err)
	// 	return
	// }

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	// req.Header.Add("apikey", "helloworld")

	// req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "f7b8c2512a00f0ca3eac801e3bfc03d95de4e39aa3bef17e3ba67c8709f18f01")

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body))

	var articleimg IMGGeneratedLPR
	json.Unmarshal(body, &articleimg)

	boxes := []AppLicense{}

	if articleimg.RecognizedData.LicenseNumber != "" {
		////fmt.Println(articleimg.ParsedResults[0].ParsedText)
		//var lines []string = regexp.MustCompile("\r?\n").Split(articleimg.ParsedResults[0].ParsedText, -1)

		boxes = append(boxes, AppLicense{Id: 1, License: articleimg.RecognizedData.LicenseNumber})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", fmt.Sprint(err))
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	} else {
		boxes = append(boxes, AppLicense{Id: 1, License: fmt.Sprint(err)})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	}

}
func resizeAnImage(imageFile multipart.File, width uint, imageType string, fileName string) {
	out, err := os.Create("resources/images/" + fileName)
	if err != nil {
		//fmt.Println(err)
	}
	switch imageType {
	case "image/jpeg":

		img, err := jpeg.Decode(imageFile)
		if err != nil {
			//fmt.Println(err)
		}

		jpegImg := resize.Resize(width, 0, img, resize.Lanczos2)

		jpeg.Encode(out, jpegImg, nil)
		break

	case "image/png":
		img, err := png.Decode(imageFile)
		if err != nil {
			//fmt.Println(err)
		}

		pngImg := resize.Resize(width, 0, img, resize.Lanczos2)

		png.Encode(out, pngImg)
		break

	case "image/gif":
		newGifImg := gif.GIF{}
		gifImg, err := gif.DecodeAll(imageFile)
		if err != nil {
			//fmt.Println(err)
		}

		for _, img := range gifImg.Image {
			resizedGifImg := resize.Resize(width, 0, img, resize.Lanczos2)
			palettedImg := image.NewPaletted(resizedGifImg.Bounds(), img.Palette)
			draw.FloydSteinberg.Draw(palettedImg, resizedGifImg.Bounds(), resizedGifImg, image.ZP)

			newGifImg.Image = append(newGifImg.Image, palettedImg)
			newGifImg.Delay = append(newGifImg.Delay, 25)
		}

		gif.EncodeAll(out, &newGifImg)

		break
	}

}

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		//fmt.Println(str)
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func GetImage(w http.ResponseWriter, r *http.Request) {

	// m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	// blue := color.RGBA{255, 0, 0, 255}
	// draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	// var img image.Image = m
	// writeImageWithTemplate(w, &img)

	// file, header, err := r.FormFile("file")
	// if err != nil {
	// 	fmt.Fprintln(w, err)
	// 	return
	// }
	// defer file.Close()

	now := time.Now()
	fmt.Println("Today : ", now.Format(time.ANSIC))

	reqBody, _ := ioutil.ReadAll(r.Body)

	//log.Println(reqBody)
	//tUnix := now.Unix()
	//strtUnix := strconv.FormatInt(tUnix, 10)

	// ioutil.WriteFile("images/IMG_"+strtUnix+".png", reqBody, 0666)

	str := base64.StdEncoding.EncodeToString(reqBody)
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
		log.Println(tmpl)
	}
	//else {
	// 	data := map[string]interface{}{"Image": str}
	// 	//fmt.Println(str)
	// 	if err = tmpl.Execute(w, data); err != nil {
	// 		log.Println("unable to execute template.")
	// 	}
	// }

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// //CUser := article.CUser
	// //EType := article.EType
	// ////fmt.Println(CUser)

	// ress99, err := db.Query("UPDATE raotdb.tblconfig SET getlastdt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), ref4='" + res41 + "' WHERE configdetail = 'truckentry' ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	ress3, err2 := db.Query("SELECT configdetail,ref1,ref2,ref3,ref4 FROM raotdb.tblconfig WHERE configdetail = 'truckentry'  Order by id  ")

	//defer ress3.Close()
	boxes := []AppReturnTruckEntry{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnTruckEntry
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Configdetail, &event.Ref1, &event.Ref2, &event.Ref3, &event.Ref4)
			if err != nil {
				panic(err)
			}
			//i, err := strconv.Atoi(event.Ref4)
			feetFloat, _ := strconv.ParseFloat(event.Ref4, 64)
			var x float64 = feetFloat
			var y int = int(x)
			//i := int(Round(event.Ref4))
			if y < 100 {
				//ioutil.WriteFile("images/IMG_"+strtUnix+".png", reqBody, 0666) /// write file

				ress99, err := db.Query("UPDATE raotdb.tblconfig SET getlastdt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), ref4='999' WHERE configdetail = 'truckentry' ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
			}

			ress99999, err := db.Query("DELETE FROM raotdb.tblimageupload WHERE transportSPID ='A6701A0167070'   ")
			defer ress99999.Close()
			if err != nil {
				panic(err)
			}

			ress999, err := db.Query("INSERT INTO raotdb.tblimageupload (transportSPID, imgBase64, uploadDT) Values ('A6701A0167070','" + str + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE imgBase64 ='" + str + "'   ")
			defer ress999.Close()
			if err != nil {
				panic(err)
			}

			//boxes = append(boxes, AppReturnTruckEntry{Configdetail: event.Configdetail, Ref1: event.Ref1, Ref2: event.Ref2, Ref3: event.Ref3, Ref4: event.Ref4})

		}

	}

	boxes = append(boxes, AppReturnTruckEntry{Configdetail: "", Ref1: "", Ref2: "", Ref3: "", Ref4: ""})

	b, _ := json.Marshal(boxes)

	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(b)
}

func GetLicenseIMGBase64(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetLicenseIMGBase64',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	CTime := article.CTime
	Token := article.Token

	if CTime == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	url := "https://api.ocr.space/parse/image"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("language", "eng")
	_ = writer.WriteField("isOverlayRequired", "true")
	_ = writer.WriteField("Base64Image", CTime)
	_ = writer.WriteField("OCREngine", "2")
	err = writer.Close()
	if err != nil {
		//fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	req.Header.Add("apikey", "helloworld")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body))

	var articleimg IMGGenerated
	json.Unmarshal(body, &articleimg)

	boxes := []AppLicense{}

	if articleimg.ParsedResults[0].ParsedText != "" {
		//fmt.Println(articleimg.ParsedResults[0].ParsedText)

		var lines []string = regexp.MustCompile("\r?\n").Split(articleimg.ParsedResults[0].ParsedText, -1)

		boxes = append(boxes, AppLicense{Id: 1, License: lines[len(lines)-1]})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	} else {
		boxes = append(boxes, AppLicense{Id: 1, License: "fail!"})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	}

}

func GetCheckPointTranDetail(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCheckPointTranDetail',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	CTime := article.CTime
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	//fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT TransportSPID, TransportID, CustomName, CusAddress, ifnull(CusPhone,'')CusPhone, ContainerType, ContainerNo, ContainerSideNo,TruckTypeID, TruckLicenseNo, TruckLicenseCountry, TruckLicenseTrail, TruckLicenseTrailCountry,ToCustomDetail, CustomID, TransferDT, NetWeight,GrossWeight,ShippingName, RubberType, ifnull(GrossWeightOnSite,'')GrossWeightOnSite,ifnull(EstWeightOnSite,'')EstWeightOnSite, ifnull(Calctype,'')Calctype, ifnull(TimeOnSite,'')TimeOnSite, ifnull(Transubstatus,'')Transubstatus, ifnull(weightID,'')weightID  FROM raotdb.tbltransportsubdetail WHERE TransportID in (SELECT TransportID FROM raotdb.tbltransport WHERE cessID = '" + CTime + "' ) ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppCheckPointTranDetailConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppCheckPointTranDetailversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransportSPID, &event.TransportID, &event.CustomName, &event.CusAddress, &event.CusPhone, &event.ContainerType, &event.ContainerNo, &event.ContainerSideNo, &event.TruckTypeID, &event.TruckLicenseNo, &event.TruckLicenseCountry, &event.TruckLicenseTrail, &event.TruckLicenseTrailCountry, &event.ToCustomDetail, &event.CustomID, &event.TransferDT, &event.NetWeight, &event.GrossWeight, &event.ShippingName, &event.RubberType, &event.GrossWeightOnSite, &event.EstWeightOnSite, &event.Calctype, &event.TimeOnSite, &event.Transubstatus, &event.WeightID)
			if err != nil {
				panic(err)
			}

			// /// SEND TO GPS

			// url := "https://api-smart-gps.aimer-psc.tech/v1/gps-data/find-license"
			// method := "POST"
			// // ทดสอบได้แค่ 65-6175

			// // 	payload := strings.NewReader(`{
			// // "license_num": "` + article.WHName + `"
			// // }`)
			// payload := strings.NewReader(`{
			// 	"license_num": "65-6175"
			// 	}`)

			// client := &http.Client{}
			// req, err := http.NewRequest(method, url, payload)

			// if err != nil {
			// 	//fmt.Println(err)
			// 	return
			// }
			// req.Header.Add("Content-Type", "application/json")

			// res, err := client.Do(req)
			// if err != nil {
			// 	//fmt.Println(err)
			// 	return
			// }
			// defer res.Body.Close()

			// body, err := ioutil.ReadAll(res.Body)
			// if err != nil {
			// 	//fmt.Println(err)
			// 	return
			// }
			// //fmt.Println(string(body))

			// ress9999, err := db.Query("INSERT INTO raotdb.tblapilog (apiname,apimessage, createdt) Values ('" + url + "','" + string(body) + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
			// defer ress9999.Close()
			// if err != nil {

			// 	// boxes = append(boxes, AppCheckPointTranDetailConfig{Id: cnt, TransportSPID: event.TransportSPID, CustomName: event.CustomName, CusAddress: event.CusAddress, CusPhone: event.CusPhone, ContainerType: event.ContainerType, ContainerNo: event.ContainerNo, ContainerSideNo: event.ContainerSideNo, TruckTypeID: event.TruckTypeID, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLicenseTrail: event.TruckLicenseTrail, TruckLicenseTrailCountry: event.TruckLicenseTrailCountry, ToCustomDetail: event.ToCustomDetail, CustomID: event.CustomID, TransferDT: event.TransferDT, NetWeight: event.NetWeight, GrossWeight: event.GrossWeight, ShippingName: event.ShippingName, RubberType: event.RubberType, GrossWeightOnSite: event.GrossWeightOnSite, EstWeightOnSite: event.EstWeightOnSite, Calctype: event.Calctype, TimeOnSite: event.TimeOnSite, Transubstatus: returnstatus, WeightID: event.WeightID, Statusflag: string(Transubstatus), GPSLAT: "", GPSLNG: "", GPSLastUpdateDT: ""})

			// 	panic(err)
			// }

			// var articlegps AutoGPSGenerated
			// json.Unmarshal(body, &articlegps)

			// GPSLat := ""
			// GPSLng := ""
			// ModifyDate := ""
			// if articlegps.Data.Gps.Lat != "" {
			// 	GPSLat = articlegps.Data.Gps.Lat
			// 	GPSLng = articlegps.Data.Gps.Lng
			// 	ModifyDate = articlegps.Data.Gps.ModifyDate.String()

			// }

			Transubstatus := string(event.Transubstatus)
			returnstatus := "/images/status0.png"
			if string(Transubstatus) == "0" {
				returnstatus = "/images/status0.png"
			} else if string(Transubstatus) == "1" {
				returnstatus = "/images/status1.png"
			} else if string(Transubstatus) == "2" {
				returnstatus = "/images/status2.png"
			} else if string(Transubstatus) == "3" {
				returnstatus = "/images/status3.png"
			} else if string(Transubstatus) == "4" {
				returnstatus = "/images/status4.png"
			} else if string(Transubstatus) == "S" {
				returnstatus = "/images/statussuccess.png"
			} else if string(Transubstatus) == "N" {
				returnstatus = "/images/statusnotsuccess.png"
			} else if string(Transubstatus) == "O" {
				returnstatus = "/images/sendstatussuccess.png"
			} else if string(Transubstatus) == "I" {
				returnstatus = "/images/statuscheckin.png"
			} else if string(Transubstatus) == "E" {
				returnstatus = "/images/statuscheckout.png"
			} else if string(Transubstatus) == "C" {
				returnstatus = "/images/statustocheckpoint.png"
			} else if string(Transubstatus) == "A" {
				returnstatus = "/images/statustocheckpoint.png"

			} else if string(Transubstatus) == "Q" {
				returnstatus = "/images/statusoutroute.png"

			} else if string(Transubstatus) == "X" {
				returnstatus = "/images/statusnotweight.png"

			} else if string(Transubstatus) == "Z" {
				returnstatus = "/images/statustocustom.png"

			} else {
				returnstatus = "/images/status0.png"
			}

			boxes = append(boxes, AppCheckPointTranDetailConfig{Id: cnt, TransportSPID: event.TransportSPID, CustomName: event.CustomName, CusAddress: event.CusAddress, CusPhone: event.CusPhone, ContainerType: event.ContainerType, ContainerNo: event.ContainerNo, ContainerSideNo: event.ContainerSideNo, TruckTypeID: event.TruckTypeID, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLicenseTrail: event.TruckLicenseTrail, TruckLicenseTrailCountry: event.TruckLicenseTrailCountry, ToCustomDetail: event.ToCustomDetail, CustomID: event.CustomID, TransferDT: event.TransferDT, NetWeight: event.NetWeight, GrossWeight: event.GrossWeight, ShippingName: event.ShippingName, RubberType: event.RubberType, GrossWeightOnSite: event.GrossWeightOnSite, EstWeightOnSite: event.EstWeightOnSite, Calctype: event.Calctype, TimeOnSite: event.TimeOnSite, Transubstatus: returnstatus, WeightID: event.WeightID, Statusflag: string(Transubstatus), GPSLAT: "", GPSLNG: "", GPSLastUpdateDT: ""})
		}

		if cnt == 0 {
			boxes = append(boxes, AppCheckPointTranDetailConfig{Id: 0, TransportSPID: "", CustomName: "", CusAddress: "", CusPhone: "", ContainerType: "", ContainerNo: "", ContainerSideNo: "", TruckTypeID: "", TruckLicenseNo: "", TruckLicenseCountry: "", TruckLicenseTrail: "", TruckLicenseTrailCountry: "", ToCustomDetail: "", CustomID: "", TransferDT: "-", NetWeight: "", GrossWeight: "", ShippingName: "", RubberType: "", GrossWeightOnSite: "", EstWeightOnSite: "", Calctype: "", TimeOnSite: "", Transubstatus: "", WeightID: "", Statusflag: "", GPSLAT: "", GPSLNG: "", GPSLastUpdateDT: ""})

		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func GetCheckPointTranDetailGPS(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCheckPointTranDetail',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	CTime := article.CTime
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT TransportSPID, TransportID, CustomName, CusAddress, ifnull(CusPhone,'')CusPhone, ContainerType, ContainerNo, ContainerSideNo,TruckTypeID, TruckLicenseNo, TruckLicenseCountry, TruckLicenseTrail, TruckLicenseTrailCountry,ToCustomDetail, CustomID, TransferDT, NetWeight,GrossWeight,ShippingName, RubberType, ifnull(GrossWeightOnSite,'')GrossWeightOnSite,ifnull(EstWeightOnSite,'')EstWeightOnSite, ifnull(Calctype,'')Calctype, ifnull(TimeOnSite,'')TimeOnSite, ifnull(Transubstatus,'')Transubstatus, ifnull(weightID,'')weightID  FROM raotdb.tbltransportsubdetail WHERE TransportID in (SELECT TransportID FROM raotdb.tbltransport WHERE cessID = '" + CTime + "' ) ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppCheckPointTranDetailConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppCheckPointTranDetailversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransportSPID, &event.TransportID, &event.CustomName, &event.CusAddress, &event.CusPhone, &event.ContainerType, &event.ContainerNo, &event.ContainerSideNo, &event.TruckTypeID, &event.TruckLicenseNo, &event.TruckLicenseCountry, &event.TruckLicenseTrail, &event.TruckLicenseTrailCountry, &event.ToCustomDetail, &event.CustomID, &event.TransferDT, &event.NetWeight, &event.GrossWeight, &event.ShippingName, &event.RubberType, &event.GrossWeightOnSite, &event.EstWeightOnSite, &event.Calctype, &event.TimeOnSite, &event.Transubstatus, &event.WeightID)
			if err != nil {
				panic(err)
			}

			Transubstatus := string(event.Transubstatus)
			returnstatus := "/images/status0.png"
			if string(Transubstatus) == "0" {
				returnstatus = "/images/status0.png"
			} else if string(Transubstatus) == "1" {
				returnstatus = "/images/status1.png"
			} else if string(Transubstatus) == "2" {
				returnstatus = "/images/status2.png"
			} else if string(Transubstatus) == "3" {
				returnstatus = "/images/status3.png"
			} else if string(Transubstatus) == "4" {
				returnstatus = "/images/status4.png"
			} else if string(Transubstatus) == "S" {
				returnstatus = "/images/statussuccess.png"
			} else if string(Transubstatus) == "N" {
				returnstatus = "/images/statusnotsuccess.png"
			} else if string(Transubstatus) == "O" {
				returnstatus = "/images/sendstatussuccess.png"
			} else if string(Transubstatus) == "I" {
				returnstatus = "/images/statuscheckin.png"
			} else if string(Transubstatus) == "E" {
				returnstatus = "/images/statuscheckout.png"
			} else if string(Transubstatus) == "C" {
				returnstatus = "/images/statustocheckpoint.png"

			} else if string(Transubstatus) == "A" {
				returnstatus = "/images/statustocheckpoint.png"

			} else if string(Transubstatus) == "Q" {
				returnstatus = "/images/statusoutroute.png"

			} else if string(Transubstatus) == "X" {
				returnstatus = "/images/statusnotweight.png"

			} else if string(Transubstatus) == "Z" {
				returnstatus = "/images/statustocustom.png"

			} else {
				returnstatus = "/images/status0.png"
			}

			/// SEND TO GPS

			url := "https://api-smart-gps.aimer-psc.tech/v1/gps-data/find-license"
			method := "POST"
			// ทดสอบได้แค่ 65-6175

			// 	payload := strings.NewReader(`{
			// "license_num": "` + article.WHName + `"
			// }`)
			payload := strings.NewReader(`{
				"license_num": "ษ-9763"
				}`)

			client := &http.Client{}
			req, err := http.NewRequest(method, url, payload)

			if err != nil {
				//fmt.Println(err)
				return
			}
			req.Header.Add("Content-Type", "application/json")

			res, err := client.Do(req)
			if err != nil {
				//fmt.Println(err)
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				//fmt.Println(err)
				return
			}
			//fmt.Println(string(body))

			ress9999, err := db.Query("INSERT INTO raotdb.tblapilog (apiname,apimessage, createdt) Values ('" + url + "','" + string(body) + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
			defer ress9999.Close()
			if err != nil {

				//panic(err)
			}

			if string(body) == "error code: 504" {
				boxes = append(boxes, AppCheckPointTranDetailConfig{Id: cnt, TransportSPID: event.TransportSPID, CustomName: event.CustomName, CusAddress: event.CusAddress, CusPhone: event.CusPhone, ContainerType: event.ContainerType, ContainerNo: event.ContainerNo, ContainerSideNo: event.ContainerSideNo, TruckTypeID: event.TruckTypeID, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLicenseTrail: event.TruckLicenseTrail, TruckLicenseTrailCountry: event.TruckLicenseTrailCountry, ToCustomDetail: event.ToCustomDetail, CustomID: event.CustomID, TransferDT: event.TransferDT, NetWeight: event.NetWeight, GrossWeight: event.GrossWeight, ShippingName: event.ShippingName, RubberType: event.RubberType, GrossWeightOnSite: event.GrossWeightOnSite, EstWeightOnSite: event.EstWeightOnSite, Calctype: event.Calctype, TimeOnSite: event.TimeOnSite, Transubstatus: returnstatus, WeightID: event.WeightID, Statusflag: string(Transubstatus), GPSLAT: "", GPSLNG: "", GPSLastUpdateDT: ""})

			} else {

				var articlegps AutoGPSGenerated
				json.Unmarshal(body, &articlegps)

				GPSLat := ""
				GPSLng := ""
				ModifyDate := ""
				if articlegps.Data.Gps.Lat != "" {
					GPSLat = articlegps.Data.Gps.Lat
					GPSLng = articlegps.Data.Gps.Lng
					ModifyDate = articlegps.Data.Gps.ModifyDate.String()

				}

				boxes = append(boxes, AppCheckPointTranDetailConfig{Id: cnt, TransportSPID: event.TransportSPID, CustomName: event.CustomName, CusAddress: event.CusAddress, CusPhone: event.CusPhone, ContainerType: event.ContainerType, ContainerNo: event.ContainerNo, ContainerSideNo: event.ContainerSideNo, TruckTypeID: event.TruckTypeID, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLicenseTrail: event.TruckLicenseTrail, TruckLicenseTrailCountry: event.TruckLicenseTrailCountry, ToCustomDetail: event.ToCustomDetail, CustomID: event.CustomID, TransferDT: event.TransferDT, NetWeight: event.NetWeight, GrossWeight: event.GrossWeight, ShippingName: event.ShippingName, RubberType: event.RubberType, GrossWeightOnSite: event.GrossWeightOnSite, EstWeightOnSite: event.EstWeightOnSite, Calctype: event.Calctype, TimeOnSite: event.TimeOnSite, Transubstatus: returnstatus, WeightID: event.WeightID, Statusflag: string(Transubstatus), GPSLAT: GPSLat, GPSLNG: GPSLng, GPSLastUpdateDT: ModifyDate})

			}

		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func GetGPSLicense(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	//CTime := article.CTime
	Token := article.Token

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetGPSLicense',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	/// SEND TO GPS

	url := "https://api-smart-gps.aimer-psc.tech/v1/gps-data/find-license"
	method := "POST"
	// ทดสอบได้แค่ 65-6175

	// 	payload := strings.NewReader(`{
	// "license_num": "` + article.WHName + `"
	// }`)
	payload := strings.NewReader(`{
	"license_num": "ษ-9763"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body))

	ress9998, err := db.Query("INSERT INTO raotdb.tblapilog (apiname,apimessage, createdt) Values ('" + url + "','" + string(body) + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
	defer ress9998.Close()
	if err != nil {
		panic(err)
	}

	var articlegps AutoGPSGenerated
	json.Unmarshal(body, &articlegps)

	boxes := []AppLicense{}

	if articlegps.Data.Gps.Lat != "" {

		boxes = append(boxes, AppLicense{Id: 1, License: article.WHName, LAT: articlegps.Data.Gps.Lat, LNG: articlegps.Data.Gps.Lng})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	} else {
		boxes = append(boxes, AppLicense{Id: 1, License: "fail!"})

		b, _ := json.Marshal(boxes)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

	}

	//AutoGPSGenerated

}

func GetCheckPoint(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)
	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	//CTime := article.CTime
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCheckPoint',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	//fmt.Println(s)

	SearchTxt := strings.Replace(article.SearchTxt, " ", "", -1)
	SearchStatus := article.SearchStatus

	strsql := "SELECT transportID as Appid, ifnull(cessID,'') as AppType, rubberType as Version ,customName as LastestDT, customName as CustomDetail ,ifnull(transubstatus,'') as TranStatus ,ifnull(transferDT,'') as shipStartDT  FROM raotdb.vw_repweightchecked WHERE transportSPID =  '" + SearchTxt + "' group by Appid, AppType, Version, LastestDT, CustomDetail, TranStatus, shipStartDT "

	strsql += "UNION SELECT transportID as Appid, cessID as AppType, rubberDesc as Version ,customerDetail as LastestDT, customName as CustomDetail ,ifnull(tranStatus,'') as TranStatus ,ifnull(shipStartDT,'') as shipStartDT  FROM raotdb.tbltransport "

	if SearchStatus == "()" && SearchTxt == "" {
		SearchStatus = "('S','N','1','2','3','O','I','E','C')"
		strsql = strsql + " WHERE tranStatus in " + SearchStatus + "   "

	} else if SearchTxt != "" || SearchStatus != "()" {
		if SearchTxt == "" {
			SearchTxt = " "
		}
		if SearchStatus == "" || SearchStatus == "('','','','','')" {
			SearchStatus = "('S','N','1','2','3','O','I','E','C')"
		}
		//strsql = strsql + " WHERE tranStatus in " + SearchStatus + " || transportID like '%" + SearchTxt + "%' || cessID like '%" + SearchTxt + "%'  "
		strsql = strsql + " WHERE   CONCAT(invoiceNo ,customerDetail , customName ,  cessID,transportID) like '%" + SearchTxt + "%' && tranStatus in " + SearchStatus + "   "
	}
	strsql = strsql + " ORDER BY shipStartDT DESC"

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppCheckPointConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppCheckPointversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT, &event.CustomDetail, &event.TranStatus, &event.ShipStartDT)
			if err != nil {
				panic(err)
			}
			returnstatus := "/images/status0.png"
			if string(event.TranStatus) == "0" {
				returnstatus = "/images/status0.png"
			} else if string(event.TranStatus) == "2" {
				returnstatus = "/images/status1.png"
			} else if string(event.TranStatus) == "3" {
				returnstatus = "/images/status2.png"
			} else if string(event.TranStatus) == "4" {
				returnstatus = "/images/status3.png"
			} else if string(event.TranStatus) == "5" {
				returnstatus = "/images/status4.png"
			} else if string(event.TranStatus) == "S" {
				returnstatus = "/images/statussuccess.png"
			} else if string(event.TranStatus) == "N" {
				returnstatus = "/images/statusnotsuccess.png"
			} else if string(event.TranStatus) == "O" {
				returnstatus = "/images/sendstatussuccess.png"
			} else if string(event.TranStatus) == "I" {
				returnstatus = "/images/statuscheckin.png"
			} else if string(event.TranStatus) == "E" {
				returnstatus = "/images/statuscheckout.png"
			} else if string(event.TranStatus) == "C" {
				returnstatus = "/images/statustocheckpoint.png"

			} else {
				returnstatus = "/images/status0.png"
			}

			boxes = append(boxes, AppCheckPointConfig{Id: cnt, AppType: event.Appid, LastestDT: event.Version, Token: event.LastestDT, TkenExpiredt: event.AppType, CustomDetail: event.CustomDetail, TranStatus: returnstatus, StatusCD: string(event.TranStatus), ShipStartDT: event.ShipStartDT})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func GetCessTran(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCessTran',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	TranId := article.CTime
	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	//fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT transportID as Appid, cessID as AppType, rubberDesc as Version ,customerDetail as LastestDT, customName as CustomDetail ,ifnull(tranStatus,'') as TranStatus, shipStartDT,netWeight,notifier  FROM raotdb.tbltransport WHERE cessID  = '" + TranId + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppCheckPointConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppCheckPointversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT, &event.CustomDetail, &event.TranStatus, &event.ShipStartDT, &event.NetWeight, &event.Notifier)
			if err != nil {
				panic(err)
			}
			returnstatus := "/images/status0.png"
			if string(event.TranStatus) == "0" {
				returnstatus = "/images/status0.png"
			} else if string(event.TranStatus) == "1" {
				returnstatus = "/images/status1.png"
			} else if string(event.TranStatus) == "2" {
				returnstatus = "/images/status2.png"
			} else if string(event.TranStatus) == "3" {
				returnstatus = "/images/status3.png"
			} else if string(event.TranStatus) == "4" {
				returnstatus = "/images/status4.png"
			} else if string(event.TranStatus) == "S" {
				returnstatus = "/images/statussuccess.png"
			} else if string(event.TranStatus) == "N" {
				returnstatus = "/images/statusnotsuccess.png"
			} else if string(event.TranStatus) == "O" {
				returnstatus = "/images/sendstatussuccess.png"
			} else if string(event.TranStatus) == "I" {
				returnstatus = "/images/statuscheckin.png"
			} else if string(event.TranStatus) == "E" {
				returnstatus = "/images/statuscheckout.png"
			} else if string(event.TranStatus) == "C" {
				returnstatus = "/images/statustocheckpoint.png"

			} else {
				returnstatus = "/images/status0.png"
			}

			boxes = append(boxes, AppCheckPointConfig{Id: cnt, AppType: event.Appid, LastestDT: event.Version, Token: event.LastestDT, TkenExpiredt: event.AppType, CustomDetail: event.CustomDetail, TranStatus: returnstatus, ShipStartDT: event.ShipStartDT, NetWeight: event.NetWeight, StatusCD: event.TranStatus, Notifier: event.Notifier})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}

func GetCheckPointLocation(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	//CTime := article.CTime
	Token := article.Token

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetCheckPointLocation',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	//fmt.Println(s)

	strsql := "SELECT idcheckpointlocation as Appid, locationname as AppType, ifnull( locationresponse,'') as Version ,locationdetail as LastestDT,ifnull(locationGPS,'') as locationGPS FROM raotdb.tblcheckpointlocation  "
	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	if article.SearchTxt != "" {
		strsql = strsql + " WHERE CONCAT(locationname,ifnull( locationresponse,''),locationdetail) like '%" + article.SearchTxt + "%'"
	}

	ress3, err2 := db.Query(strsql)
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppWeightConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT, &event.LocationGPS)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppWeightConfig{Id: cnt, AppType: event.Appid, LastestDT: event.Version, Token: event.LastestDT, TkenExpiredt: event.AppType, LocationGPS: event.LocationGPS})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func GetConfigWeight(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	// WHName := article.WHName
	// Weight := article.Weight
	//CTime := article.CTime

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetConfigWeight',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	//////

	Token := article.Token

	// if CTime == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }
	if Token == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	//fmt.Println(s)

	strsql := "SELECT weightname as Appid, weightcustomname as AppType, weightresponsename as Version ,createdt as LastestDT FROM raotdb.tblweight"

	if article.SearchCustomname != "" && article.SearchTxt != "" {
		strsql = strsql + " WHERE weightcustomname like '%" + article.SearchCustomname + "%' || CONCAT(weightname , weightresponsename,weightcustomname , weightlocation)  like '%" + article.SearchTxt + "%'"
	} else {

		if article.SearchCustomname != "" {
			strsql = strsql + " WHERE weightcustomname like '%" + article.SearchCustomname + "%'"
		}
		if article.SearchTxt != "" {
			strsql = strsql + " WHERE CONCAT(weightname , weightresponsename,weightcustomname , weightlocation)  like '%" + article.SearchTxt + "%'"
		}

	}

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query(strsql)
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppWeightConfig{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppWeightConfig{Id: cnt, AppType: event.Appid, LastestDT: event.Version, Token: event.LastestDT, TkenExpiredt: event.AppType})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Write(b)

		//counter = 0

	}
}
func ConfigWeight(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	WHName := article.WHName
	Weight := article.Weight
	CTime := article.CTime
	EType := article.EType
	CUser := article.CUser
	CLocal := article.CLocal

	fmt.Println(CTime)

	if Weight == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetWeight',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	//////

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte("RAOTSMART")
	t = jwt.New(jwt.SigningMethodHS256)
	s, err = t.SignedString(key)
	//fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT weightname as Appid, weightcustomname as AppType, weightresponsename as Version ,createdt as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + Weight + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT)
			if err != nil {
				panic(err)
			}
			// if event.Version == WHName {
			// 	respok := make(map[string]string)
			// 	respok["WHName"] = "Duplicate"
			// 	respok["Token"] = s //QR

			// 	jsonResp, err := json.Marshal(respok)

			// 	if err != nil {
			// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// 		return
			// 	}
			// 	w.Write(jsonResp)
			// 	return
			// }

			//boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Version, LastestDT: event.LastestDT, Token: s})
		}

		if cnt == 0 {
			ress99, err := db.Query("INSERT INTO raotdb.tblweight (weightname, weightresponsename,weightcustomname,weightlocation, createdt) Values ('" + Weight + "','" + CUser + "','" + WHName + "','" + CLocal + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightname ='" + Weight + "'   ")
			defer ress99.Close()
			if err != nil {
				panic(err)
			}
			var event Appversion

			boxes = append(boxes, AppToken{Appid: Weight, AppType: WHName, LastestDT: event.LastestDT, Token: s})

		} else {

			if EType == "edit" {
				ress99, err := db.Query("UPDATE  raotdb.tblweight SET weightresponsename = '" + CUser + "',weightcustomname = '" + WHName + "',weightlocation = '" + CLocal + "',modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')   WHERE weightname ='" + Weight + "' ")
				//ress99, err := db.Query("INSERT raotdb.tblweight (weightname, weightresponsename,weightcustomname, createdt) Values ('" + Weight + "','" + WHName + "','" + CTime + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightname ='" + Weight + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: Weight, AppType: WHName, LastestDT: event.LastestDT, Token: s})
			} else if EType == "delete" {
				ress99, err := db.Query("DELETE FROM  raotdb.tblweight  WHERE weightname ='" + Weight + "' and  weightname not in (SELECT WHName FROM raotdb.tblcurrentweight ) ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: Weight, AppType: WHName, LastestDT: event.LastestDT, Token: s})

			} else {

				respok := make(map[string]string)
				respok["WHName"] = "Duplicate"
				respok["Token"] = s //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return

			}

		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}

func SetJobDetail(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStartJob
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.JobID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//CheckPointResponse := article.CessID
	//EType := article.EType

	//fmt.Println(CUser)

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetJobDetail',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	ress99, err := db.Query("INSERT INTO raotdb.tbljobdetail (jobID, cessID, transportID,transportsubID,userid,createDT) Values ('" + article.JobID + "','" + article.CessID + "','" + article.TranID + "','','" + article.CUser + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightOnScale =''   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	ress999, err := db.Query("UPDATE raotdb.tbljobmaster SET cpstartDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') WHERE jobID = '" + article.JobID + "'   ")
	defer ress999.Close()
	if err != nil {
		panic(err)
	}

	strsql := "SELECT jobid , weightid  FROM raotdb.tbljobmaster WHERE jobID = '" + article.JobID + "' ORDER BY cpid DESC "

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnJobID{}
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnJobID
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.JobID, &event.WeightId)
			if err != nil {
				panic(err)
			}
			boxes = append(boxes, AppReturnJobID{JobID: event.JobID, WeightId: event.WeightId})
		}
	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}
func SetConsent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendConsent
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//CheckPointResponse := article.CessID
	//EType := article.EType

	//fmt.Println(CUser)

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetConsent',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	if article.EType == "Add" {
		ress99, err := db.Query("INSERT INTO raotdb.tblconsentmsg (msgid,consentmsg, modifydt ) Values ('1','" + article.ConsentMsg + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE ConsentMsg ='" + article.ConsentMsg + "'   ")
		defer ress99.Close()
		if err != nil {
			panic(err)
		}
	}

	strsql := "SELECT msgid , consentmsg  FROM raotdb.tblconsentmsg WHERE msgid = '1' ORDER BY msgid DESC "

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnConsent{}
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnConsent
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.MsgID, &event.ConsentMsg)
			if err != nil {
				panic(err)
			}
			boxes = append(boxes, AppReturnConsent{MsgID: event.MsgID, ConsentMsg: event.ConsentMsg})
		}
	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetStartCheckPoint(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	//WHName := article.WHName
	// StartDT := article.StartDT
	// StopDT := article.StopDT
	// CheckPointID := article.CheckPointID

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetStartCheckPoint',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	//CUser := article.CUser
	CheckPointResponse := article.CheckPointResponse
	//EType := article.EType

	//fmt.Println(CUser)

	// if WHName == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)
	cntjob := 0

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	strsql := "SELECT cpid , checkpointcode , checkpointname  ,DATE_FORMAT(StartDate, '%d-%m-%Y %T')StartDate  ,DATE_FORMAT(Stopdate, '%d-%m-%Y %T')Stopdate ,weightid,ifnull(locationresponse,'')locationresponse,createdt as LastestDT,ifnull(jobid,'')JobID, '' as HaveJob FROM raotdb.tblcheckpoint WHERE statususe is null  "

	if article.SearchCustomname != "" {
		strsql = strsql + "and (  checkpointname like '%" + article.SearchCustomname + "%'  ) "
	}
	if article.SearchDT != "" {
		strsql = strsql + "and (  StartDate between '" + strings.Split(article.SearchDT, "|")[0] + "' and '" + strings.Split(article.SearchDT, "|")[1] + "' ) "
	}
	if article.SearchTxt != "" {
		strsql = strsql + "and (  locationresponse like '%" + article.SearchTxt + "%' ) "
	}
	if article.SearchCustomname != "" && article.SearchDT != "" && article.SearchTxt != "" {
		strsql = strsql + "and (  locationresponse = '%" + article.SearchCustomname + "%' || StartDate between '" + strings.Split(article.SearchDT, "|")[0] + "' and '" + strings.Split(article.SearchDT, "|")[1] + "' || checkpointname like '%" + article.SearchTxt + " %' ) "
	}

	strsql = strsql + "ORDER BY cpid DESC"

	if CheckPointResponse != "" {

		ress33, err2 := db.Query("SELECT jobID,cessID,transportID  FROM raotdb.tbljobdetail WHERE cessID  = '" + article.WHName + "' ")
		if err2 == nil {
			for ress33.Next() {
				cntjob++
				var event AppGetJobDetail
				//JobID := ress2.Scan(&event.JobID)
				err := ress33.Scan(&event.CessID, &event.JobID, &event.TransportID)
				if err != nil {
					panic(err)
				}

				//boxes = append(boxes, AppGetJobDetail{CPID: event.CPID, CheckpointCode: event.CheckpointCode, CheckpointName: event.CheckpointName, StartDate: event.StartDate, StopDate: event.StopDate, WeightId: event.WeightId, LocationResponse: event.LocationResponse, CreateDate: event.LastestDT, JobID: event.JobID})

			}

		}
		defer ress33.Close()

		if cntjob >= 1 {
			strsql = "SELECT cpid , checkpointcode , checkpointname  ,DATE_FORMAT(StartDate, '%d-%m-%Y %T')StartDate  ,DATE_FORMAT(Stopdate, '%d-%m-%Y %T')Stopdate ,weightid,ifnull(locationresponse,'')locationresponse,createdt as LastestDT,ifnull(jobid,'')JobID , '1' as HaveJob FROM raotdb.tblcheckpoint WHERE statususe is null and locationresponse = '" + CheckPointResponse + "'  ORDER BY cpid DESC "

		} else {
			strsql = "SELECT cpid , checkpointcode , checkpointname  ,DATE_FORMAT(StartDate, '%d-%m-%Y %T')StartDate  ,DATE_FORMAT(Stopdate, '%d-%m-%Y %T')Stopdate ,weightid,ifnull(locationresponse,'')locationresponse,createdt as LastestDT,ifnull(jobid,'')JobID , '0' as HaveJob FROM raotdb.tblcheckpoint WHERE statususe is null and locationresponse = '" + CheckPointResponse + "' and StartDate >= DATE_SUB(NOW(), INTERVAL 2 DAY)  ORDER BY cpid DESC "

		}

	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnCheckpointStart{}
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppConfigCheckpointStart
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.CPID, &event.CheckpointCode, &event.CheckpointName, &event.WeightId, &event.StartDate, &event.StopDate, &event.LastestDT, &event.LocationResponse, &event.JobID, &event.HaveJob)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnCheckpointStart{CPID: event.CPID, CheckpointCode: event.CheckpointCode, CheckpointName: event.CheckpointName, StartDate: event.StartDate, StopDate: event.StopDate, WeightId: event.WeightId, LocationResponse: event.LocationResponse, CreateDate: event.LastestDT, JobID: event.JobID, HaveJob: event.HaveJob})

		}

	}

	if len(boxes) == 0 {
		boxes = append(boxes, AppReturnCheckpointStart{CPID: "", CheckpointCode: "", CheckpointName: "", StartDate: "", StopDate: "", WeightId: "", LocationResponse: "", CreateDate: "", JobID: "", HaveJob: ""})

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRepTrader(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//EType := article.EType

	//fmt.Println(CUser)

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRepTrader',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	strsql := "SELECT distinct(name_th)customName FROM raotdb.vw_tradertruck"

	// if article.SearchTxt != "" {
	// 	strsql = strsql + " WHERE CONCAT( transportSPID, transportID, customName,truckLicenseNo) like '%" + article.SearchTxt + "%'"
	// }

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnRepTrader{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnRepTrader
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.CustomName)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnRepTrader{CustomName: event.CustomName})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRepWeightChecked(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRepWeightChecked',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///
	//EType := article.EType

	//fmt.Println(CUser)
	//strsql := "SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
	strsql := "SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, CASE WHEN truckLPRlicenseCountryOverwrite like '%มาเล%' THEN CONCAT(truckLPRlicenseOverwrite , '-มาเลเซีย') ELSE ''  END as licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"

	if article.SearchTxt != "" && article.SearchCustomname != "" {
		strsql = strsql + " WHERE CONCAT( customName, transportID,cessID, transportSPID,licenseMY, truckLicenseNo,checkpointname,locationresponse) like '%" + article.SearchTxt + "%'  || customName like '%" + article.SearchCustomname + "%' and  arrivalDT is not null"

	} else {
		if article.SearchTxt != "" {
			strsql = strsql + " WHERE CONCAT(customName, transportID,cessID, transportSPID,licenseMY, truckLicenseNo,checkpointname,locationresponse) like '%" + article.SearchTxt + "%' and  arrivalDT is not null "
		} else {
			//strsql = strsql + " WHERE  arrivalDT is not null"
		}

		if article.SearchCustomname != "" {
			//strsql = strsql + " UNION SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
			strsql = strsql + " WHERE  customName like '%" + article.SearchCustomname + "%'   "
		}

		if article.SearchDT != "" {
			//strsql = strsql + " UNION SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
			strsql = strsql + " WHERE  arrivalDT between '" + strings.Split(article.SearchDT, "|")[0] + "' and '" + strings.Split(article.SearchDT, "|")[1] + "'  "
		}

		if article.SearchRubberType != "" {
			//strsql = strsql + "UNION SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
			strsql = strsql + " WHERE  rubberType like '%" + article.SearchRubberType + "%' "
		}

		if article.SearchResponse != "" {
			//strsql = strsql + "UNION SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
			strsql = strsql + " WHERE  locationresponse like '%" + article.SearchResponse + "%' "
		}

		if article.SearchLocationName != "" {
			//strsql = strsql + "UNION SELECT transferDT, rubberType, customName, ifnull(cessID,'')cessID, transportID, transportSPID, licenseMY, truckLicenseNo, truckLicenseCountry, ifnull(truckLPRlicenseOverwrite,'')truckLPRlicenseOverwrite, ifnull(truckLPRlicenseCountryOverwrite,'')truckLPRlicenseCountryOverwrite, truckTypeID,ifnull(checkpointname,'')checkpointname, ifnull(checkpointcode,'')checkpointcode,ifnull( name,'')name, ifnull(locationresponse,'')locationresponse, ifnull(grouporganize,'')grouporganize, ifnull(transubstatus,'')transubstatus, ifnull(arrivalDT,'')arrivalDT, ifnull(netWeight,'')netWeight, ifnull(estWeightOnSite,'')estWeightOnSite, ifnull(netallweight,'')netallweight, ifnull(SumEstWeightOnSite,'')SumEstWeightOnSite FROM raotdb.vw_repweightchecked"
			strsql = strsql + " WHERE  checkpointname like '%" + strings.Split(article.SearchLocationName, " ")[1] + "%' "
		}

		// if article.SearchCustomname != "" {
		// 	strsql = strsql + " WHERE name_th like '%" + article.SearchCustomname + "%'"
		// }
	}

	strsql = strsql + " ORDER BY customName, transferDT "
	// if article.SearchTxt != "" {
	// 	strsql = strsql + " WHERE CONCAT( transportSPID, transportID, customName,truckLicenseNo) like '%" + article.SearchTxt + "%'"
	// }

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnRepWeightChecked{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnRepWeightChecked
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransferDT, &event.RubberType, &event.CustomName, &event.CessID, &event.TransportID, &event.TransportSPID, &event.LicenseMY, &event.TruckLicenseNo, &event.TruckLicenseCountry, &event.TruckLPRlicenseOverwrite, &event.TruckLPRlicenseCountryOverwrite, &event.TruckTypeID, &event.Checkpointname, &event.Checkpointcode, &event.Name, &event.Locationresponse, &event.Grouporganize, &event.Transubstatus, &event.ArrivalDT, &event.NetWeight, &event.EstWeightOnSite, &event.Netallweight, &event.SumEstWeightOnSite)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnRepWeightChecked{TransferDT: event.TransferDT, RubberType: event.RubberType, CustomName: event.CustomName, CessID: event.CessID, TransportID: event.TransportID, TransportSPID: event.TransportSPID, LicenseMY: event.LicenseMY, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLPRlicenseOverwrite: event.TruckLPRlicenseOverwrite, TruckLPRlicenseCountryOverwrite: event.TruckLPRlicenseCountryOverwrite, TruckTypeID: event.TruckTypeID, Checkpointname: event.Checkpointname, Checkpointcode: event.Checkpointcode, Name: event.Name, Locationresponse: event.Locationresponse, Grouporganize: event.Grouporganize, Transubstatus: event.Transubstatus, ArrivalDT: event.ArrivalDT, NetWeight: event.NetWeight, EstWeightOnSite: event.EstWeightOnSite, Netallweight: event.Netallweight, SumEstWeightOnSite: event.SumEstWeightOnSite})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRepRubberType(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRepRubberType',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	//EType := article.EType

	//fmt.Println(CUser)
	strsql := "SELECT distinct(rubberType)rubberType FROM raotdb.vw_reptradertruck"

	// if article.SearchTxt != "" {
	// 	strsql = strsql + " WHERE CONCAT( transportSPID, transportID, customName,truckLicenseNo) like '%" + article.SearchTxt + "%'"
	// }

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnRubberType{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnRubberType
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.RubberType)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnRubberType{RubberType: event.RubberType})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}
func GetRepTraderTruck(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRepTraderTruck',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	//EType := article.EType

	//fmt.Println(CUser)
	strsql := "SELECT '' as transportSPID, '' as  transportID,name_th as customName, ifnull(tax_id,'')cusCitizenID, '' as  cusAddress, ifnull(phone,'')cusPhone,'' as containerType,'' as containerNo,'' as containerSideNo, ifnull(truck_type,'') as truckTypeID, ifnull(license_number_th,'') as truckLicenseNo,ifnull(license_province_th,'') as truckLicenseCountry,'' as truckLicenseTrail,'' as truckLicenseTrailCountry, '' as toCustomDetail, '' as  customID, '' as  transferDT,ifnull(weight_kg,'') as  netWeight,ifnull(weight_kg,'') as grossWeight,'' as  shippingName,'' as  totalBox,ifnull(registration_type,'') as rubberType,  '' as createDT,  '' as modifyDT, '' as transubstatus, '' as grossWeightOnSite, '' as estWeightOnSite,'' as calctype,  '' as timeOnSite,  '' as weightID, '' as userID, '' as truckLPRlicenseOverwrite, '' as truckLPRlicenseCountryOverwrite, '' as truckLPRTypeOverwrite, '' as truckLPRTypeTruckOverwrite, ifnull(gps_provider_name,'') as GPSProvider, ifnull(gps_box_id,'') as GPSBoxID,ifnull(license_number_my,'') as TruckLicenseNoMY FROM raotdb.vw_tradertruck  "

	if article.SearchTxt != "" && article.SearchCustomname != "" {
		strsql = strsql + " WHERE CONCAT( name_en, name_th, license_number_th,license_number_my) like '%" + article.SearchTxt + "%'  || customName like '%" + article.SearchCustomname + "%'"

	} else {
		if article.SearchTxt != "" {
			strsql = strsql + " WHERE CONCAT( name_en, name_th, license_number_th,license_number_my) like '%" + article.SearchTxt + "%'"
		}

		if article.SearchCustomname != "" {
			strsql = strsql + " WHERE name_th like '%" + article.SearchCustomname + "%'"
		}
	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnRepTraderTruck{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnRepTraderTruck
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TransportSPID, &event.TransportID, &event.CustomName, &event.CusCitizenID, &event.CusAddress, &event.CusPhone, &event.ContainerType, &event.ContainerNo, &event.ContainerSideNo, &event.TruckTypeID, &event.TruckLicenseNo, &event.TruckLicenseCountry, &event.TruckLicenseTrail, &event.TruckLicenseTrailCountry, &event.ToCustomDetail, &event.CustomID, &event.TransferDT, &event.NetWeight, &event.GrossWeight, &event.ShippingName, &event.TotalBox, &event.RubberType, &event.CreateDT, &event.ModifyDT, &event.Transubstatus, &event.GrossWeightOnSite, &event.EstWeightOnSite, &event.Calctype, &event.TimeOnSite, &event.WeightID, &event.UserID, &event.TruckLPRlicenseOverwrite, &event.TruckLPRlicenseCountryOverwrite, &event.TruckLPRTypeOverwrite, &event.TruckLPRTypeTruckOverwrite, &event.GPSProvider, &event.GPSBoxID, &event.TruckLicenseNoMY)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnRepTraderTruck{TransportSPID: event.TransportSPID, TransportID: event.TransportID, CustomName: event.CustomName, CusCitizenID: event.CusCitizenID, CusAddress: event.CusAddress, CusPhone: event.CusPhone, ContainerType: event.ContainerType, ContainerNo: event.ContainerNo, ContainerSideNo: event.ContainerSideNo, TruckTypeID: event.TruckTypeID, TruckLicenseNo: event.TruckLicenseNo, TruckLicenseCountry: event.TruckLicenseCountry, TruckLicenseTrail: event.TruckLicenseTrail, TruckLicenseTrailCountry: event.TruckLicenseTrailCountry, ToCustomDetail: event.ToCustomDetail, CustomID: event.CustomID, TransferDT: event.TransferDT, NetWeight: event.NetWeight, GrossWeight: event.GrossWeight, ShippingName: event.ShippingName, TotalBox: event.TotalBox, RubberType: event.RubberType, CreateDT: event.CreateDT, ModifyDT: event.ModifyDT, Transubstatus: event.Transubstatus, GrossWeightOnSite: event.GrossWeightOnSite, EstWeightOnSite: event.EstWeightOnSite, Calctype: event.Calctype, TimeOnSite: event.TimeOnSite, WeightID: event.WeightID, UserID: event.UserID, TruckLPRlicenseOverwrite: event.TruckLPRlicenseOverwrite, TruckLPRlicenseCountryOverwrite: event.TruckLPRlicenseCountryOverwrite, TruckLPRTypeOverwrite: event.TruckLPRTypeOverwrite, TruckLPRTypeTruckOverwrite: event.TruckLPRTypeTruckOverwrite, GPSProvider: event.GPSProvider, GPSBoxID: event.GPSBoxID, TruckLicenseNoMY: event.TruckLicenseNoMY})

		}

	}

	if len(boxes) == 0 {
		boxes = append(boxes, AppReturnRepTraderTruck{TransportSPID: "", TransportID: "", CustomName: "", CusCitizenID: "", CusAddress: "", CusPhone: "", ContainerType: "", ContainerNo: "", ContainerSideNo: "", TruckTypeID: "", TruckLicenseNo: "", TruckLicenseCountry: "", TruckLicenseTrail: "", TruckLicenseTrailCountry: "", ToCustomDetail: "", CustomID: "", TransferDT: "", NetWeight: "", GrossWeight: "", ShippingName: "", TotalBox: "", RubberType: "", CreateDT: "", ModifyDT: "", Transubstatus: "", GrossWeightOnSite: "", EstWeightOnSite: "", Calctype: "", TimeOnSite: "", WeightID: "", UserID: "", TruckLPRlicenseOverwrite: "", TruckLPRlicenseCountryOverwrite: "", TruckLPRTypeOverwrite: "", TruckLPRTypeTruckOverwrite: "", GPSProvider: "", GPSBoxID: "", TruckLicenseNoMY: ""})

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func APITruckCheckInOut(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckInOut
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.CustomerName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)

	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	// var TruckRAOT TruckRAOTDetail
	// jsonData2 := []byte(string(body2))

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('APITruckCheckInOut',requestcount+1,'" + channel[0] + "','GPS' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	if article.Status == "check-in" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'I' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}

		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ( 'ถึงสถานที่โหลดสินค้า','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}

	}
	if article.Status == "check-out" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'E' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}
		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ('ออกจากสถานที่โหลดสินค้า','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}
	}
	if article.Status == "arrivecheckpoint" || article.Status == "ArriveCheckpoint" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'A' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}
		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ('ถึงจุดตรวจ','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}
	}
	if article.Status == "outroute" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'Q' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}
		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ('ออกนอกเส้นทาง','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}
	}
	if article.Status == "notentrytoweightscale" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'X' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}
		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ('ไม่ขึ้นชั่งน้ำหนัก','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}
	}
	if article.Status == "tocustom" {

		ress9998, err := db.Query("UPDATE raotdb.tbltransportsubdetail SET transubstatus = 'Z' , GPSTruckLatLong = '" + article.LAT + "," + article.LNG + "' WHERE transportSPID = '" + article.TransportSubID + "' ")
		defer ress9998.Close()
		if err != nil {
			panic(err)
		}
		ress99988, err := db.Query("INSERT INTO raotdb.tblnotifier ( msgtxt, msgstatus, msgref1,msgref2,msgref3, msgdt) Values ('ถึงด่านศุลกากร','unread','" + article.License + "','" + article.TransportSubID + "','" + article.Status + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
		defer ress99988.Close()
		if err != nil {
			panic(err)
		}
	}
	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	// ress3, err2 := db.Query("SELECT id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, company, status FROM raotdb.tbltrucks   ")

	// //defer ress3.Close()
	// boxes := []TruckTrader{}
	// if err2 == nil {

	// 	cnt := 0

	// 	for ress3.Next() {

	// 		cnt++

	// 		var event TruckTrader
	// 		//JobID := ress2.Scan(&event.JobID)
	// 		err := ress3.Scan(&event.ID, &event.License_number_th, &event.License_province_th, &event.License_number_my, &event.Registration_type, &event.Ttruck_type, &event.Weight_kg, &event.Gps_box_id, &event.Gps_provider_name, &event.Company, &event.Status)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		boxes = append(boxes, TruckTrader{ID: event.ID, License_number_th: event.License_number_th, License_province_th: event.License_province_th, License_number_my: event.License_number_my, Registration_type: event.Registration_type, Ttruck_type: event.Ttruck_type, Weight_kg: event.Weight_kg, Gps_box_id: event.Gps_box_id, Gps_provider_name: event.Gps_provider_name, Company: event.Company, Status: event.Status})
	// 		break
	// 	}

	// }

	// b, _ := json.Marshal(boxes)

	// defer ress3.Close()
	// err = ress3.Close()
	// //defer ress2.Close()
	// //jsonResp, err := json.Marshal(b)
	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	return
	// }

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// // 	}
	// // }
	// // w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// // w.Header().Set("Connection", "Close")
	// w.Write(b)

	respok := make(map[string]string)
	respok["Success"] = "true"
	respok["response"] = "Update : license " + article.License + " Status : " + article.Status //QR

	jsonResp, err := json.Marshal(respok)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonResp)

	//counter = 0

}
func SndDatatoCESS1(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	// var TruckRAOT TruckRAOTDetail
	// jsonData2 := []byte(string(body2))

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SndDatatoCESS1',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	ress9998, err := db.Query("UPDATE raotdb.tbltransport SET tranStatus = 'O' WHERE cessID = '" + article.WHName + "' ")
	defer ress9998.Close()
	if err != nil {
		panic(err)
	}

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, company, status FROM raotdb.tbltrucks   ")

	//defer ress3.Close()
	boxes := []TruckTrader{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event TruckTrader
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.ID, &event.License_number_th, &event.License_province_th, &event.License_number_my, &event.Registration_type, &event.Ttruck_type, &event.Weight_kg, &event.Gps_box_id, &event.Gps_provider_name, &event.Company, &event.Status)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, TruckTrader{ID: event.ID, License_number_th: event.License_number_th, License_province_th: event.License_province_th, License_number_my: event.License_number_my, Registration_type: event.Registration_type, Ttruck_type: event.Ttruck_type, Weight_kg: event.Weight_kg, Gps_box_id: event.Gps_box_id, Gps_provider_name: event.Gps_provider_name, Company: event.Company, Status: event.Status})
			break
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SetDatatoCumtomer(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	tokenkey := raottokenrefresh

	if raottokenrefresh == "" {
		url := "https://cess-export.olive.co.th/api/token/"
		method := "POST"

		payload := strings.NewReader(`{
			"username": "ong_olive",
			"password": "olive1234"
		}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			//fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			//fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			//fmt.Println(err)
			return
		}
		//fmt.Println(string(body))

		var LoadTruck TruckRAOT
		jsonData := []byte(string(body))

		js := TruckRAOT{}

		json.Unmarshal(jsonData, &js)
		if err != nil {
			panic(err)
		}
		//fmt.Println(js)

		json.Unmarshal([]byte(body), &LoadTruck)

		////fmt.Println(string(logistic.Data.Booking[0].BookingID))
		//fmt.Println(LoadTruck.Access)
		//fmt.Println(LoadTruck.Refresh)

		tokenkey = LoadTruck.Access
		raottokenrefresh = LoadTruck.Refresh
	}

	// if LoadTruck.Refresh != "" {
	// 	tokenkey = LoadTruck.Refresh
	// }
	t := time.Now()
	tUnix := t.Unix()
	fmt.Printf("timeUnix: %d\n", tUnix)
	str := strconv.FormatInt(tUnix, 10)

	url := "https://cess-export.olive.co.th/api/declaration_invoices/DELARE_1/"
	method := "PUT"
	payload := strings.NewReader(`{
		"declaration_id": "` + article.WHName + `",
		"net_weight": "` + article.CUser + `",
		"status": "` + article.PositionName + `", 
		"status_updated_at": "` + str + `"
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	//req.Header.Add("Authorization", "f7b8c2512a00f0ca3eac801e3bfc03d95de4e39aa3bef17e3ba67c8709f18f01")
	req.Header.Add("Authorization", "Bearer "+tokenkey)
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println(req.Header)
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body2, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body2))

	// var TruckRAOT TruckRAOTDetail
	// jsonData2 := []byte(string(body2))

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetDatatoCumtomer',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	ress9998, err := db.Query("INSERT INTO raotdb.tblapilog (apiname,apimessage, createdt) Values ('" + url + "','" + string(body2) + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
	defer ress9998.Close()
	if err != nil {
		panic(err)
	}

	// if WHName == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }]
	//i := 0
	// for i < len(TruckRAOT) {
	// 	//i++
	// 	ress99, err := db.Query("INSERT INTO raotdb.tbltrucks (id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, created_at, updated_at, company, status) Values ('" + strconv.Itoa(TruckRAOT[i].ID) + "','" + TruckRAOT[i].LicenseNumberTh + "','" + TruckRAOT[i].LicenseProvinceTh + "','" + TruckRAOT[i].LicenseNumberMy + "','" + TruckRAOT[i].RegistrationType + "','" + TruckRAOT[i].TruckType + "','" + TruckRAOT[i].WeightKg + "','" + TruckRAOT[i].GpsBoxID + "','" + TruckRAOT[i].GpsProviderName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + TruckRAOT[i].Company + "','" + TruckRAOT[i].Status + "' ) ON DUPLICATE KEY UPDATE updated_at = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')   ")
	// 	defer ress99.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	i++
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, company, status FROM raotdb.tbltrucks  LIMIT 10 ")

	//defer ress3.Close()
	boxes := []TruckTrader{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event TruckTrader
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.ID, &event.License_number_th, &event.License_province_th, &event.License_number_my, &event.Registration_type, &event.Ttruck_type, &event.Weight_kg, &event.Gps_box_id, &event.Gps_provider_name, &event.Company, &event.Status)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, TruckTrader{ID: event.ID, License_number_th: event.License_number_th, License_province_th: event.License_province_th, License_number_my: event.License_number_my, Registration_type: event.Registration_type, Ttruck_type: event.Ttruck_type, Weight_kg: event.Weight_kg, Gps_box_id: event.Gps_box_id, Gps_provider_name: event.Gps_provider_name, Company: event.Company, Status: event.Status})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SetWidget(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWAWidget
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.UserID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if article.IMG != "img" {
		ress99, err := db.Query("UPDATE raotdb.tbluserraot SET imgprofile = '" + article.IMG + "' WHERE username = '" + article.UserID + "'  ")
		defer ress99.Close()
		if err != nil {
			panic(err)
		}
	}

	ress3, err2 := db.Query("SELECT imgprofile FROM raotdb.tbluserraot   WHERE username = '" + article.UserID + "'")

	//defer ress3.Close()
	boxes := []Widget{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event Widget
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.IMG)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Widget{IMG: event.IMG})
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SndNotiTracking(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//ress3, err2 := db.Query("call raotdb.sp_dashboard ('" + article.StartDT + "','" + article.StopDT + "');  ")

	ress99, err := db.Query("UPDATE  raotdb.tbltransport SET notifier = '" + article.PositionName + "'   WHERE transportID ='" + article.WHName + "' ")
	//ress99, err := db.Query("INSERT raotdb.tblweight (weightname, weightresponsename,weightcustomname, createdt) Values ('" + Weight + "','" + WHName + "','" + CTime + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightname ='" + Weight + "'   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	ress3, err2 := db.Query("call raotdb.sp_dashboard ('" + article.StartDT + "','" + article.StopDT + "');  ")

	defer ress3.Close()
	boxes := []Dashboard{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event Dashboard
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DMY, &event.Cntall, &event.Cntsuccess, &event.Cntincheckpoint, &event.Cntnotincheckpoint)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Dashboard{DMY: event.DMY, Cntall: event.Cntall, Cntsuccess: event.Cntsuccess, Cntincheckpoint: event.Cntincheckpoint, Cntnotincheckpoint: event.Cntnotincheckpoint})
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ress3, err2 := db.Query("call raotdb.sp_dashboard ('" + article.StartDT + "','" + article.StopDT + "');  ")
	defer ress3.Close()
	//defer ress3.Close()
	boxes := []Dashboard{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event Dashboard
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DMY, &event.Cntall, &event.Cntsuccess, &event.Cntincheckpoint, &event.Cntnotincheckpoint)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Dashboard{DMY: event.DMY, Cntall: event.Cntall, Cntsuccess: event.Cntsuccess, Cntincheckpoint: event.Cntincheckpoint, Cntnotincheckpoint: event.Cntnotincheckpoint})
		}

	}

	b, _ := json.Marshal(boxes)

	//err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetNoti(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ress3, err2 := db.Query("SELECT msgid, msgtxt, msgstatus, msgref1, msgdt,ifnull(transportID,'')transportID,ifnull(transportSPID,'')transportSPID, customName  FROM raotdb.vw_notitruck  ")

	//defer ress3.Close()
	boxes := []TruckNoti{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event TruckNoti
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.MSGID, &event.MSGTxt, &event.MSGStatus, &event.MSGRef1, &event.MSGDT, &event.TransportID, &event.TransportSPID, &event.CustomName)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, TruckNoti{MSGID: event.MSGID, MSGTxt: event.MSGTxt, MSGStatus: event.MSGStatus, MSGRef1: event.MSGRef1, MSGDT: event.MSGDT, TransportID: event.TransportID, TransportSPID: event.TransportSPID, CustomName: event.CustomName})
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetTruck(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	tokenkey := raottokenrefresh

	if raottokenrefresh == "" {
		url := "https://cess-export.olive.co.th/api/token/"
		method := "POST"

		payload := strings.NewReader(`{
			"username": "wiwat",
			"password": "guskeenneo1971"
		}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			//fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			//fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			//fmt.Println(err)
			return
		}
		//fmt.Println(string(body))

		var LoadTruck TruckRAOT
		jsonData := []byte(string(body))

		js := TruckRAOT{}

		json.Unmarshal(jsonData, &js)
		if err != nil {
			panic(err)
		}
		//fmt.Println(js)

		json.Unmarshal([]byte(body), &LoadTruck)

		////fmt.Println(string(logistic.Data.Booking[0].BookingID))
		//fmt.Println(LoadTruck.Access)
		//fmt.Println(LoadTruck.Refresh)

		tokenkey = LoadTruck.Access
		raottokenrefresh = LoadTruck.Refresh
	}

	// if LoadTruck.Refresh != "" {
	// 	tokenkey = LoadTruck.Refresh
	// }

	url := "https://cess-export.olive.co.th/api/trucks/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		//fmt.Println(err)
		return
	}
	//req.Header.Add("Authorization", "f7b8c2512a00f0ca3eac801e3bfc03d95de4e39aa3bef17e3ba67c8709f18f01")
	req.Header.Add("Authorization", "Bearer "+tokenkey)
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println(req.Header)
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body2, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body2))

	var TruckRAOT TruckRAOTDetailTrader
	jsonData2 := []byte(string(body2))

	js2 := TruckRAOTDetailTrader{}

	json.Unmarshal(jsonData2, &js2)
	if err != nil {
		panic(err)
	}
	//fmt.Println(js2)

	json.Unmarshal([]byte(body2), &TruckRAOT)

	////fmt.Println(string(logistic.Data.Booking[0].BookingID))
	//fmt.Println(TruckRAOT[0].ID)
	//fmt.Println(TruckRAOT[0].LicenseNumberTh)

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//EType := article.EType

	//fmt.Println(CUser)

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetTruck',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	// if WHName == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }]
	i := 0
	for i < len(TruckRAOT) {
		//i++
		ress99, err := db.Query("INSERT INTO raotdb.tbltrucks (id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, created_at, updated_at, company, status) Values ('" + strconv.Itoa(TruckRAOT[i].ID) + "','" + TruckRAOT[i].LicenseNumberTh + "','" + TruckRAOT[i].LicenseProvinceTh + "','" + TruckRAOT[i].LicenseNumberMy + "','" + TruckRAOT[i].RegistrationType + "','1','" + TruckRAOT[i].WeightKg + "','" + TruckRAOT[i].GpsBoxID + "','" + TruckRAOT[i].GpsProviderName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + TruckRAOT[i].Company + "','" + TruckRAOT[i].Status + "' ) ON DUPLICATE KEY UPDATE updated_at = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')   ")
		defer ress99.Close()
		if err != nil {
			panic(err)
		}
		i++
	}

	if err == nil {
		url := "https://cess-export.olive.co.th/api/trailers/"
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			//fmt.Println(err)
			return
		}
		//req.Header.Add("Authorization", "f7b8c2512a00f0ca3eac801e3bfc03d95de4e39aa3bef17e3ba67c8709f18f01")
		req.Header.Add("Authorization", "Bearer "+tokenkey)
		req.Header.Add("Content-Type", "application/json")
		//fmt.Println(req.Header)
		res, err := client.Do(req)
		if err != nil {
			//fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body2, err := ioutil.ReadAll(res.Body)
		if err != nil {
			//fmt.Println(err)
			return
		}
		//fmt.Println(string(body2))

		var trailers TrailersRAOTDetailTrader
		jsonData2 := []byte(string(body2))

		js2 := TrailersRAOTDetailTrader{}

		json.Unmarshal(jsonData2, &js2)
		if err != nil {
			panic(err)
		}
		//fmt.Println(js2)

		json.Unmarshal([]byte(body2), &trailers)

		////fmt.Println(string(logistic.Data.Booking[0].BookingID))
		//fmt.Println(trailers[0].ID)
		//fmt.Println(trailers[0].LicenseNumberTh)

		// dns := getDNSString(dblogin, userlogin, passlogin, conn)
		// db, err := sql.Open("mysql", dns)

		// if err != nil {
		// 	panic(err)
		// }
		// err = db.Ping()
		// if err != nil {
		// 	panic(err)
		// }
		// defer db.Close()

		//CUser := article.CUser
		//EType := article.EType

		//fmt.Println(CUser)

		// if WHName == "" {
		// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		// 	//return
		// }]
		i := 0
		for i < len(trailers) {
			//i++
			ress99, err := db.Query("INSERT INTO raotdb.tbltrucks (id, license_number_th, license_province_th, license_number_my, registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, created_at, updated_at, company, status) Values ('" + strconv.Itoa(trailers[i].ID) + "','" + trailers[i].LicenseNumberTh + "','" + trailers[i].LicenseProvinceTh + "','" + trailers[i].LicenseNumberMy + "','" + trailers[i].RegistrationType + "','2','" + trailers[i].WeightKg + "','" + trailers[i].GpsBoxID + "','" + trailers[i].GpsProviderName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + trailers[i].Company + "','" + trailers[i].Status + "' ) ON DUPLICATE KEY UPDATE updated_at = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')   ")
			defer ress99.Close()
			if err != nil {
				panic(err)
			}
			i++
		}
	}

	ress99, err := db.Query("call raotdb.sp_inserttradertrucktocalc")
	defer ress99.Close()
	if err != nil {
		//panic(err)
	}
	//

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT id, license_number_th, license_province_th, ifnull(license_number_my,'')license_number_my, ifnull(registration_type,'')registration_type, truck_type, weight_kg, gps_box_id, gps_provider_name, company, status FROM raotdb.tbltrucks   ")

	//defer ress3.Close()
	boxes := []TruckTrader{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event TruckTrader
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.ID, &event.License_number_th, &event.License_province_th, &event.License_number_my, &event.Registration_type, &event.Ttruck_type, &event.Weight_kg, &event.Gps_box_id, &event.Gps_provider_name, &event.Company, &event.Status)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, TruckTrader{ID: event.ID, License_number_th: event.License_number_th, License_province_th: event.License_province_th, License_number_my: event.License_number_my, Registration_type: event.Registration_type, Ttruck_type: event.Ttruck_type, Weight_kg: event.Weight_kg, Gps_box_id: event.Gps_box_id, Gps_provider_name: event.Gps_provider_name, Company: event.Company, Status: event.Status})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetTrader(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	//WHName := article.WHName
	// StartDT := article.StartDT
	// StopDT := article.StopDT
	// CheckPointID := article.CheckPointID
	//CUser := article.CUser
	//EType := article.EType

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetTrader',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	//fmt.Println(CUser)

	// if WHName == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	strsql := "SELECT dbid, trader_tax_id, trader_branch_no, trader_name, trader_name_en, trader_type, address, subdistrict, district, province, postal_code, phone_number, fax_number, email FROM raotdb.tbltrader   "

	if article.SearchTxt != "" {
		strsql = strsql + " WHERE CONCAT( trader_name, trader_name_en, trader_type, address, subdistrict, district, province, postal_code, phone_number, fax_number, email) like '%" + article.SearchTxt + "%'"
	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnTrader{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppConfigTrader
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Trader_tax_id, &event.Trader_branch_no, &event.Trader_name, &event.Trader_name_en, &event.Trader_type, &event.Address, &event.Subdistrict, &event.District, &event.Province, &event.Postal_code, &event.Phone_number, &event.Fax_number, &event.Email)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnTrader{DBID: event.DBID, Trader_tax_id: event.Trader_tax_id, Trader_branch_no: event.Trader_branch_no, Trader_name: event.Trader_name, Trader_name_en: event.Trader_name_en, Trader_type: event.Trader_type, Address: event.Address, Subdistrict: event.Subdistrict, District: event.District, Province: event.Province, Postal_code: event.Postal_code, Phone_number: event.Phone_number, Email: event.Email})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SetCalculate(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASetCalculate
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.Authorization)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetCalculate',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	TruckID := article.TruckID
	TruckRubberType := article.TruckRubberType
	ContainerType := article.ContainerType
	ContainerTypeID := article.ContainerTypeID
	TruckWeight := article.TruckWeight
	TrailWeight := article.TrailWeight
	TrailCount := article.TrailCount
	ContainerWeight := article.ContainerWeight
	ContainerCount := article.ContainerCount
	BoxWeight := article.BoxWeight
	BoxCount := article.BoxCount

	TruckWeight = strings.ReplaceAll(TruckWeight, ",", "")
	TrailWeight = strings.ReplaceAll(TrailWeight, ",", "")
	ContainerWeight = strings.ReplaceAll(ContainerWeight, ",", "")

	if BoxWeight == "" {
		BoxWeight = "0"
	}
	if BoxCount == "" {
		BoxCount = "0"
	}
	if TruckWeight == "" {
		TruckWeight = "0"
	}
	if TrailWeight == "" {
		TrailWeight = "0"
	}
	if TrailCount == "" {
		TrailCount = "0"
	}
	if ContainerWeight == "" {
		ContainerWeight = "0"
	}
	if ContainerCount == "" {
		ContainerCount = "0"
	}
	CUser := article.CUser
	EType := article.EType

	//fmt.Println(CUser)

	if CUser == "" {
		//log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	strsql := ""
	if EType == "add" {
		//log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
		strsql = "SELECT CalcID, RubberType, ContainerType ,TruckWeight  FROM raotdb.tblcalc WHERE rubberType  = '" + TruckRubberType + "'  and permanent is null "
	} else {
		strsql = "SELECT CalcID, RubberType, ContainerType ,TruckWeight FROM raotdb.tblcalc WHERE calcID  = '" + TruckID + "'  "
	}

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT cpid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT,stopdate  ,stopdate ,weightid FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppConfigCheckpointCalculate
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.CalcID, &event.RubberType, &event.ContainerType, &event.TruckWeight)
			if err != nil {
				panic(err)
			}
		}

		if cnt == 0 {
			ress99, err := db.Query("INSERT INTO raotdb.tblcalc (  rubberType, containerType,ContainerTypeID, truckWeight, trailWeight, trailCount, tankweight, containerWeight, totalcontainer, boxWeight, totalbox, createdt, userid) Values ('" + TruckRubberType + "','" + ContainerType + "','" + ContainerTypeID + "','" + TruckWeight + "','" + TrailWeight + "','" + TrailCount + "','" + ContainerWeight + "','" + ContainerWeight + "','" + ContainerCount + "','" + BoxWeight + "','" + BoxCount + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + CUser + "' ) ON DUPLICATE KEY UPDATE rubberType ='" + TruckRubberType + "'   ")
			defer ress99.Close()
			if err != nil {
				panic(err)
			}
			var event Appversion

			boxes = append(boxes, AppToken{Appid: TruckRubberType, AppType: ContainerType, LastestDT: event.LastestDT, Token: CUser})

		} else {

			if EType == "edit" {
				ress99, err := db.Query("UPDATE  raotdb.tblcalc SET rubberType = '" + TruckRubberType + "',containerType = '" + ContainerType + "',containerTypeID = '" + ContainerTypeID + "',truckWeight = '" + TruckWeight + "',trailWeight = '" + TrailWeight + "',trailCount = '" + TrailCount + "',tankweight = '" + ContainerWeight + "',containerWeight = '" + ContainerWeight + "',totalcontainer = '" + ContainerCount + "',boxWeight = '" + BoxWeight + "',totalbox = '" + BoxCount + "',modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),userid = '" + CUser + "'   WHERE calcID ='" + TruckID + "' ")
				//ress99, err := db.Query("INSERT raotdb.tblweight (weightname, weightresponsename,weightcustomname, createdt) Values ('" + Weight + "','" + WHName + "','" + CTime + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightname ='" + Weight + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: TruckRubberType, AppType: ContainerType, LastestDT: event.LastestDT, Token: CUser})
			} else if EType == "delete" {
				ress99, err := db.Query("DELETE FROM raotdb.tblcalc  WHERE calcID = '" + TruckID + "' and permanent is null  ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: TruckRubberType, AppType: ContainerType, LastestDT: event.LastestDT, Token: CUser})

			} else {

				respok := make(map[string]string)
				respok["WHName"] = "Duplicate"
				respok["Token"] = TruckRubberType //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return

			}

		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}

func ConfigStartCheckPoint(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('ConfigStartCheckPoint',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	WHName := article.WHName
	StartDT := article.StartDT
	StopDT := article.StopDT
	CheckPointID := article.CheckPointID
	CheckPointLocationName := article.CheckPointLocationName
	CheckPointResponse := article.CheckPointResponse
	CUser := article.CUser
	EType := article.EType

	//fmt.Println(CUser)

	if WHName == "" {
		//log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	strexpiredAt := strconv.FormatInt(expiredAt, 10)

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)
	strsql := ""
	if EType == "add" {
		//log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
		strsql = "SELECT WeightId, CheckpointCode, CheckpointName ,StartDate, StopDate FROM raotdb.tblcheckpoint WHERE checkpointcode = '" + CheckPointID + "' and weightid  = '" + WHName + "' and startdate = '" + StartDT + "'  and statususe is null "
	} else {
		strsql = "SELECT WeightId, CheckpointCode, CheckpointName ,StartDate, StopDate FROM raotdb.tblcheckpoint WHERE cpid  = '" + WHName + "'  "
	}

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT cpid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT,stopdate  ,stopdate ,weightid FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppConfigCheckpointStart
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.CheckpointCode, &event.CheckpointName, &event.WeightId, &event.StartDate, &event.StopDate)
			if err != nil {
				panic(err)
			}
			// if event.Version == WHName {
			// 	respok := make(map[string]string)
			// 	respok["WHName"] = "Duplicate"
			// 	respok["Token"] = s //QR

			// 	jsonResp, err := json.Marshal(respok)

			// 	if err != nil {
			// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// 		return
			// 	}
			// 	w.Write(jsonResp)
			// 	return
			// }

			//boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Version, LastestDT: event.LastestDT, Token: s})
		}

		if cnt == 0 {
			ress99, err := db.Query("INSERT INTO raotdb.tblcheckpoint (weightid, checkpointcode,checkpointname,locationresponse,startdate,stopdate, createdt,usercreateid,jobid) Values ('" + WHName + "','" + CheckPointID + "','" + CheckPointLocationName + "','" + CheckPointResponse + "',CONVERT_TZ('" + StartDT + "','+00:00','+0:00'),CONVERT_TZ('" + StopDT + "','+00:00','+0:00'),CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + CUser + "','J" + strexpiredAt + "' ) ON DUPLICATE KEY UPDATE weightid ='" + WHName + "'   ")
			defer ress99.Close()
			if err != nil {
				panic(err)
			}

			ress999, err := db.Query("INSERT INTO raotdb.tbljobmaster (jobID, userID, totalaudit,weightid, createDT) Values ('J" + strexpiredAt + "','" + CUser + "','0','" + WHName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
			defer ress999.Close()
			if err != nil {
				panic(err)
			}

			ress9999, err := db.Query("UPDATE raotdb.tbljobmaster tableB INNER JOIN raotdb.tblcheckpoint tableA ON tableB.jobID = tableA.jobID SET tableB.cpid = tableA.cpid  WHERE tableA.jobID = 'J" + strexpiredAt + "'  ")
			defer ress9999.Close()
			if err != nil {
				panic(err)
			}

			var event Appversion

			boxes = append(boxes, AppToken{Appid: CheckPointID, AppType: WHName, LastestDT: event.LastestDT, Token: WHName})

		} else {

			strsql = "SELECT jobID FROM raotdb.tblcheckpoint  WHERE cpid ='" + WHName + "' and jobID in (SELECT jobID FROM raotdb.tbljobdetail ) "
			ress333, err2 := db.Query(strsql)

			cnt333 := 0
			for ress333.Next() {

				cnt333++
			}
			if err2 != nil {
				panic(err)
			}
			defer ress333.Close()
			err = ress333.Close()

			if cnt333 > 0 {
				respok := make(map[string]string)
				respok["WHName"] = "not edit"
				respok["Token"] = WHName //QR

				//jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				boxes = append(boxes, AppToken{Appid: CheckPointID, AppType: "not edit", LastestDT: "", Token: WHName})

				// w.Write(jsonResp)
				// return
			} else {

				if EType == "edit" {
					ress99, err := db.Query("UPDATE  raotdb.tblcheckpoint SET checkpointcode = '" + CheckPointID + "',checkpointname = '" + WHName + "',startdate = '" + StartDT + "',stopdate = '" + StopDT + "',modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),usermodify = '" + CUser + "'   WHERE cpid ='" + WHName + "' ")
					//ress99, err := db.Query("INSERT raotdb.tblweight (weightname, weightresponsename,weightcustomname, createdt) Values ('" + Weight + "','" + WHName + "','" + CTime + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE weightname ='" + Weight + "'   ")
					defer ress99.Close()
					if err != nil {
						panic(err)
					}
					var event Appversion
					boxes = append(boxes, AppToken{Appid: CheckPointID, AppType: WHName, LastestDT: event.LastestDT, Token: WHName})
				} else if EType == "delete" {
					ress99, err := db.Query("UPDATE  raotdb.tblcheckpoint SET statususe = 'F',modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),usermodify = '" + CUser + "' WHERE cpid ='" + WHName + "'   ")
					defer ress99.Close()
					if err != nil {
						panic(err)
					}
					var event Appversion
					boxes = append(boxes, AppToken{Appid: CheckPointID, AppType: WHName, LastestDT: event.LastestDT, Token: WHName})

				} else {

					respok := make(map[string]string)
					respok["WHName"] = "Duplicate"
					respok["Token"] = WHName //QR

					jsonResp, err := json.Marshal(respok)

					if err != nil {
						log.Fatalf("Error happened in JSON marshal. Err: %s", err)
						return
					}
					w.Write(jsonResp)
					return

				}
			}

		}

		/// ส่งข้อมูลให้กับ GPS เพื่อกำหนดจุดตรวจ
		//WHName := article.WHName
		//StartDT := article.StartDT
		//StopDT := article.StopDT
		//CheckPointID := article.CheckPointID
		//CheckPointLocationName := article.CheckPointLocationName
		//CheckPointResponse := article.CheckPointResponse

		ress3, err2 := db.Query("SELECT locationgps,locationdetail FROM raotdb.tblcheckpointlocation WHERE locationresponse  = '" + CheckPointResponse + "' ")
		//defer ress3.Close()
		if err2 == nil {

			//boxes := []AppToken{}
			cnt := 0

			for ress3.Next() {

				cnt++

				var event Appversion
				//JobID := ress2.Scan(&event.JobID)
				err := ress3.Scan(&event.LocationGPS, &event.Locationdetail)
				if err != nil {
					panic(err)
				}

				url := "https://api-smart-gps.aimer-psc.tech/v1/check-point/create"
				method := "POST"

				payload := strings.NewReader(`{
					"CheckpointID": "` + article.CheckPointID + `",
					"LocationName": "` + CheckPointLocationName + `",
					"LocationResponse": "` + CheckPointResponse + `",
					"LocationDetail": "` + event.Locationdetail + `",
					"lat": "` + strings.Split(event.LocationGPS, ",")[0] + `",
					"lng":"` + strings.Split(event.LocationGPS, ",")[1] + `",
					"WeightID": "` + WHName + `",
					"StartDate": "` + article.StDategps + `",
					"StopDate": "` + article.EtDategps + `"
				}`)

				client := &http.Client{}
				req, err := http.NewRequest(method, url, payload)

				if err != nil {
					//fmt.Println(err)
					return
				}
				req.Header.Add("Content-Type", "application/json")

				res, err := client.Do(req)
				if err != nil {
					//fmt.Println(err)
					return
				}
				defer res.Body.Close()

				//body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					//fmt.Println(err)
					return
				}
				//fmt.Println(string(body))

				defer ress3.Close()
				err = ress3.Close()
				//defer ress2.Close()
				//jsonResp, err := json.Marshal(b)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
			}
		}
		b, _ := json.Marshal(boxes)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}
func ConfigCheckPointLocation(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendLocation
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.LocationID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('ConfigCheckPointLocation',requestcount+1,'" + channel[0] + "','" + article.LocationID + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	LocationID := article.LocationID
	LocationResponse := article.LocationResponse
	LocationName := article.LocationName
	LocationDetail := article.LocationDetail
	LocationGPS := article.LocationGPS
	EType := article.EType

	//fmt.Println(LocationID)

	/// SEND TO GPS

	url := "https://api-smart-gps.aimer-psc.tech/v1/check-point/create"
	method := "POST"

	payload := strings.NewReader(`{
	"CheckpointID": "` + article.LocationID + `",
	"LocationName": "` + article.LocationName + `",
	"LocationResponse": "` + article.LocationResponse + `",
  	"LocationDetail": "` + article.LocationDetail + `",
	"LocationGPS": "` + article.LocationGPS + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		//fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//APItTruckCheckInOut
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(body))

	ress9998, err := db.Query("INSERT INTO raotdb.tblapilog (apiname,apimessage, createdt) Values ('" + url + "','" + string(body) + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )   ")
	defer ress9998.Close()
	if err != nil {
		panic(err)
	}

	//////

	// if Weight == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte("RAOTSMART")
	t = jwt.New(jwt.SigningMethodHS256)
	s, err = t.SignedString(key)
	//fmt.Println(s)

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT idcheckpointlocation as Appid, locationname as AppType, ifnull(locationresponse,'') as Version ,ifnull(createdt,'') as LastestDT FROM raotdb.tblcheckpointlocation WHERE idcheckpointlocation  = '" + LocationID + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}
		cnt := 0

		for ress3.Next() {

			cnt++

			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.Appid, &event.AppType, &event.Version, &event.LastestDT)
			if err != nil {
				panic(err)
			}
			// if event.Version == WHName {
			// 	respok := make(map[string]string)
			// 	respok["WHName"] = "Duplicate"
			// 	respok["Token"] = s //QR

			// 	jsonResp, err := json.Marshal(respok)

			// 	if err != nil {
			// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// 		return
			// 	}
			// 	w.Write(jsonResp)
			// 	return
			// }

			//boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Version, LastestDT: event.LastestDT, Token: s})
		}

		if cnt == 0 {
			ress99, err := db.Query("INSERT INTO raotdb.tblcheckpointlocation (idcheckpointlocation, locationname, locationresponse, locationdetail,locationgps, createdt) Values ('" + LocationID + "','" + LocationName + "','" + LocationResponse + "','" + LocationDetail + "','" + LocationGPS + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE idcheckpointlocation ='" + LocationID + "'   ")
			defer ress99.Close()
			if err != nil {
				panic(err)
			}
			var event Appversion

			boxes = append(boxes, AppToken{Appid: LocationID, AppType: LocationName, LastestDT: event.LastestDT, Token: s})

		} else {

			if EType == "edit" {
				ress99, err := db.Query("UPDATE  raotdb.tblcheckpointlocation SET locationname = '" + LocationName + "',locationresponse = '" + LocationResponse + "',locationdetail = '" + LocationDetail + "',locationgps = '" + LocationGPS + "',modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')   WHERE idcheckpointlocation ='" + LocationID + "' ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: LocationID, AppType: LocationName, LastestDT: event.LastestDT, Token: s})
			} else if EType == "delete" {
				ress99, err := db.Query("DELETE FROM  raotdb.tblcheckpointlocation  WHERE idcheckpointlocation ='" + LocationID + "' and  idcheckpointlocation not in (SELECT WHName FROM raotdb.tblcurrentweight ) ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}
				var event Appversion
				boxes = append(boxes, AppToken{Appid: LocationID, AppType: LocationName, LastestDT: event.LastestDT, Token: s})

			} else {

				respok := make(map[string]string)
				respok["WHName"] = "Duplicate"
				respok["Token"] = s //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return

			}

		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}
func SendWeight(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendWeight
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn) //(dblogin, userlogin, passlogin, conn)

	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SendWeight',requestcount+1,'" + channel[0] + "','" + article.WHName + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	//////

	WHName := article.WHName
	Weight := article.Weight
	CTime := article.CTime

	if CTime == "" {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		//return
	}

	ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight,CShaftNum, CDateTime) Values ('" + WHName + "','" + Weight + "','" + CTime + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte("PWASMART0805401073")
	t = jwt.New(jwt.SigningMethodHS256)
	s, err = t.SignedString(key)
	//fmt.Println(s)

	ress3, err2 := db.Query("SELECT userid as AppType, weightid as Version FROM raotdb.tbluser WHERE username = '0805401073' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AppToken{}

		for ress3.Next() {
			var event Appversion
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.AppType, &event.Version)
			if err != nil {
				panic(err)
			}
			if event.Version == WHName {
				respok := make(map[string]string)
				respok["WHName"] = WHName
				respok["Token"] = s //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return
			}

			boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Version, Token: s})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Content-Type", "text/plain")
		// w.WriteHeader(http.StatusOK)

		// if acrh, ok := r.Header["Access-Control-Request-Headers"]; ok {
		// 	w.Header().Set("Access-Control-Allow-Headers", acrh[0])
		// }
		// w.Header().Set("Access-Control-Allow-Credentials", "True")
		// if acao, ok := r.Header["Access-Control-Allow-Origin"]; ok {
		// 	w.Header().Set("Access-Control-Allow-Origin", acao[0])
		// } else {
		// 	if _, oko := r.Header["Origin"]; oko {
		// 		w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
		// 	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// 	}
		// }
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		// w.Header().Set("Connection", "Close")
		w.Write(b)

		//counter = 0

	}
}

func SetRAOTUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article GetRaotUser
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.UserName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	CUser := article.UserName
	//EType := article.EType

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetRAOTUser',requestcount+1,'" + channel[0] + "','" + CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	//fmt.Println(CUser)

	expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	strexpiredAt := strconv.FormatInt(expiredAt, 10)

	ress3, err2 := db.Query("SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize , email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role  WHERE username = '" + article.Email + "'  ")

	//defer ress3.Close()
	boxes := []AppReturnUserRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Username, &event.Name, &event.Surename, &event.Position, &event.Department, &event.Locationresponse, &event.Grouporganize, &event.Email, &event.Permission, &event.Lastlogindt, &event.Logincounter, &event.KeycloakUserRole)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserRAOT{DBID: event.DBID, Username: event.Username, Name: event.Name, Surename: event.Surename, Position: event.Position, Department: event.Department, Locationresponse: event.Locationresponse, Grouporganize: event.Grouporganize, Email: event.Email, Permission: event.Permission, Lastlogindt: event.Lastlogindt, Logincounter: event.Logincounter, KeycloakUserRole: event.KeycloakUserRole})

		}

		if cnt >= 0 {

			if article.EType != "delete" {
				ress99, err := db.Query("INSERT INTO raotdb.tbluserraot (id, username, name, surename, position, department, locationresponse, email, permission,createdt) Values ( '" + strexpiredAt + "', '" + article.Email + "','" + article.UserName + "','','" + article.Position + "','" + article.Department + "','" + article.ResponseGroup + "','" + article.Email + "','" + article.Permission + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE name ='" + article.UserName + "' , surename = '', position='" + article.Position + "', department='" + article.Department + "', locationresponse='" + article.ResponseGroup + "', email='" + article.Email + "', permission='" + article.Permission + "',createdt=CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserRAOT{DBID: "1", Username: article.Email, Name: article.UserName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			} else {

				ress998, err := db.Query("INSERT INTO raotdb.tbluserraotbackup SELECT * FROM  raotdb.tbluserraot  WHERE  username  ='" + article.Email + "'  ")
				defer ress998.Close()
				if err != nil {
					panic(err)
				}

				ress99, err := db.Query("DELETE FROM raotdb.tbluserraot WHERE  username ='" + article.Email + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				ress999, err := db.Query("DELETE FROM raotdb.tbluser WHERE  username ='" + article.Email + "'   ")
				defer ress999.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserRAOT{DBID: "1", Username: article.Email, Name: article.UserName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			}

			//return
			return
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}
func SetRAOTUserRollback(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article GetRaotUser
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.UserName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	CUser := article.UserName
	//EType := article.EType

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetRAOTUser',requestcount+1,'" + channel[0] + "','" + CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	//fmt.Println(CUser)

	expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	strexpiredAt := strconv.FormatInt(expiredAt, 10)

	ress3, err2 := db.Query("SELECT id as DBID, username, name, surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize , email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role_backup  WHERE username = '" + article.Email + "'  ")

	//defer ress3.Close()
	boxes := []AppReturnUserRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Username, &event.Name, &event.Surename, &event.Position, &event.Department, &event.Locationresponse, &event.Grouporganize, &event.Email, &event.Permission, &event.Lastlogindt, &event.Logincounter, &event.KeycloakUserRole)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserRAOT{DBID: event.DBID, Username: event.Username, Name: event.Name, Surename: event.Surename, Position: event.Position, Department: event.Department, Locationresponse: event.Locationresponse, Grouporganize: event.Grouporganize, Email: event.Email, Permission: event.Permission, Lastlogindt: event.Lastlogindt, Logincounter: event.Logincounter, KeycloakUserRole: event.KeycloakUserRole})

		}

		if cnt >= 0 {

			if article.EType != "delete" {
				ress99, err := db.Query("INSERT INTO raotdb.tbluserraot (id, username, name, surename, position, department, locationresponse, email, permission,createdt) Values ( '" + strexpiredAt + "', '" + article.Email + "','" + article.UserName + "','','" + article.Position + "','" + article.Department + "','" + article.ResponseGroup + "','" + article.Email + "','" + article.Permission + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE name ='" + article.UserName + "' , surename = '', position='" + article.Position + "', department='" + article.Department + "', locationresponse='" + article.ResponseGroup + "', email='" + article.Email + "', permission='" + article.Permission + "',createdt=CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserRAOT{DBID: "1", Username: article.Email, Name: article.UserName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			} else {

				ress998, err := db.Query("INSERT INTO raotdb.tbluserraot SELECT * FROM  raotdb.tbluserraotbackup  WHERE  username  ='" + article.Email + "'  ")
				defer ress998.Close()
				if err != nil {
					panic(err)
				}

				ress99, err := db.Query("DELETE FROM raotdb.tbluserraotbackup WHERE  username ='" + article.Email + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserRAOT{DBID: "1", Username: article.Email, Name: article.UserName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			}

			//return
			return
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SetRAOTUserPosition(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article AppSetUserPositionRAOT
	json.Unmarshal(reqBody, &article)

	////fmt.Println("Positionname", article.Positionname)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	////fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.UserID
	//EType := article.EType

	////fmt.Println(CUser)

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetRAOTUserPosition',requestcount+1,'" + channel[0] + "','" + article.UserID + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///

	// if WHName == "" {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	//return
	// }

	// ress99, err := db.Query("INSERT INTO raotdb.tblcurrentweight (WHName, OnlineWeight, CDateTime) Values ('" + WHName + "','" + Weight + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE OnlineWeight ='" + Weight + "'   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("RAOTSMART")
	// t = jwt.New(jwt.SigningMethodHS256)
	// s, err = t.SignedString(key)
	// //fmt.Println(s)
	expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	strexpiredAt := strconv.FormatInt(expiredAt, 10)
	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query("SELECT dbid, positionname, ifnull(isdashboard,'')isdashboard, ifnull(ismanageweight,'')ismanageweight, ifnull(ismanagecheckpoint,'')ismanagecheckpoint, ifnull(ischeckpoint,'')ischeckpoint, ifnull(ischeckpointlist,'')ischeckpointlist, ifnull(Ismanagecalculate,'')Ismanagecalculate, ifnull(iscumtomerdata,'')iscumtomerdata, ifnull(istracking,'')istracking, ifnull(isreportlicense,'')isreportlicense, ifnull(isreporttruckchecked,'')isreporttruckchecked, ifnull(isreportweightchecked,'')isreportweightchecked, ifnull(ismanageemployee,'')ismanageemployee, ifnull(ismanageposition,'')ismanageposition, ifnull(ismanagelocationresponse,'')ismanagelocationresponse, ifnull(ismanagegroup,'')ismanagegroup, ifnull(iscondition,'')iscondition, ifnull(isreporttrader,'')isreporttrader FROM raotdb.tbluserposition  WHERE positionname = '" + article.Positionname + "'  ")

	//defer ress3.Close()
	boxes := []AppReturnUserPositionRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserPositionRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Positionname, &event.Isdashboard, &event.Ismanageweight, &event.Ismanagecheckpoint, &event.Ischeckpoint, &event.Ischeckpointlist, &event.Ismanagecalculate, &event.Iscumtomerdata, &event.Istracking, &event.Isreportlicense, &event.Isreporttruckchecked, &event.Isreportweightchecked, &event.Ismanageemployee, &event.Ismanageposition, &event.Ismanagelocationresponse, &event.Ismanagegroup, &event.Iscondition, &event.Isreporttrader)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserPositionRAOT{DBID: event.DBID, Positionname: event.Positionname, Isdashboard: event.Isdashboard, Ismanageweight: event.Ismanageweight, Ismanagecheckpoint: event.Ismanagecheckpoint, Ischeckpoint: event.Ischeckpoint, Ischeckpointlist: event.Ischeckpointlist, Ismanagecalculate: event.Ismanagecalculate, Iscumtomerdata: event.Iscumtomerdata, Istracking: event.Istracking, Isreportlicense: event.Isreportlicense, Isreporttruckchecked: event.Isreporttruckchecked, Isreportweightchecked: event.Isreportweightchecked, Ismanageemployee: event.Ismanageemployee, Ismanageposition: event.Ismanageposition, Ismanagelocationresponse: event.Ismanagelocationresponse, Ismanagegroup: event.Ismanagegroup, Iscondition: event.Iscondition, Isreporttrader: event.Isreporttrader})

		}

		if cnt >= 0 {

			if article.EType != "delete" {
				ress99, err := db.Query("INSERT INTO raotdb.tbluserposition (dbid,  positionname, isdashboard, ismanageweight, ismanagecheckpoint, ischeckpoint, ischeckpointlist,ismanagecalculate, iscumtomerdata, istracking, isreportlicense, isreporttruckchecked, isreportweightchecked, ismanageemployee, ismanageposition, ismanagelocationresponse, ismanagegroup, iscondition, isreporttrader, createdt, modifydt, editbyuser) Values ('" + strexpiredAt + "','" + article.Positionname + "','" + article.Isdashboard + "','" + article.Ismanageweight + "','" + article.Ismanagecheckpoint + "','" + article.Ischeckpoint + "','" + article.Ischeckpointlist + "','" + article.Ismanagecalculate + "','" + article.Iscumtomerdata + "','" + article.Istracking + "','" + article.Isreportlicense + "','" + article.Isreporttruckchecked + "','" + article.Isreportweightchecked + "','" + article.Ismanageemployee + "','" + article.Ismanageposition + "' ,'" + article.Ismanagelocationresponse + "','" + article.Ismanagegroup + "','" + article.Iscondition + "', isreporttrader='" + article.Isreporttrader + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + article.UserID + "' ) ON DUPLICATE KEY UPDATE isdashboard ='" + article.Isdashboard + "' ,ismanageweight='" + article.Ismanageweight + "', ismanagecheckpoint='" + article.Ismanagecheckpoint + "', ischeckpoint='" + article.Ischeckpoint + "', ischeckpointlist='" + article.Ischeckpointlist + "',ismanagecalculate='" + article.Ismanagecalculate + "', iscumtomerdata='" + article.Iscumtomerdata + "', istracking='" + article.Istracking + "', isreportlicense='" + article.Isreportlicense + "', isreporttruckchecked='" + article.Isreporttruckchecked + "', isreportweightchecked='" + article.Isreportweightchecked + "', ismanageemployee='" + article.Ismanageemployee + "', ismanageposition='" + article.Ismanageposition + "', ismanagelocationresponse='" + article.Ismanagelocationresponse + "', ismanagegroup='" + article.Ismanagegroup + "', iscondition='" + article.Iscondition + "', isreporttrader='" + article.Isreporttrader + "', modifydt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), editbyuser ='" + article.UserID + "' ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserPositionRAOT{DBID: "1", Positionname: article.Positionname})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			} else {
				ress99, err := db.Query("DELETE FROM raotdb.tbluserposition WHERE positionname ='" + article.Positionname + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnUserPositionRAOT{DBID: "1", Positionname: article.Positionname})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			}

			return
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRAOTUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	////fmt.Println("User", article.User)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	////fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//EType := article.EType

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRAOTUser',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	////fmt.Println(CUser)

	strsql := "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role  Order by id  "

	if article.User != "all" {
		strsql = "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role WHERE email = '" + article.User + "'  Order by id  "
	}

	if article.SearchTxt != "" {
		strsql = "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role WHERE CONCAT(username, name, position, department, locationresponse) like '%" + article.SearchTxt + "%'  Order by id  "
	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnUserRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Username, &event.Name, &event.Surename, &event.Position, &event.Department, &event.Locationresponse, &event.Grouporganize, &event.Email, &event.Permission, &event.Lastlogindt, &event.Logincounter, &event.KeycloakUserRole)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserRAOT{DBID: event.DBID, Username: event.Username, Name: event.Name, Surename: event.Surename, Position: event.Position, Department: event.Department, Locationresponse: event.Locationresponse, Grouporganize: event.Grouporganize, Email: event.Email, Permission: event.Permission, Lastlogindt: event.Lastlogindt, Logincounter: event.Logincounter, KeycloakUserRole: event.KeycloakUserRole})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRAOTUserBackup(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	////fmt.Println("User", article.User)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	////fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser
	//EType := article.EType

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRAOTUserBackup',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	////fmt.Println(CUser)

	strsql := "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role_backup  Order by id  "

	if article.User != "all" {
		strsql = "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role_backup WHERE email = '" + article.User + "'  Order by id  "
	}

	if article.SearchTxt != "" {
		strsql = "SELECT id as DBID, username, name, ifnull(surename,'')surename, position, department, locationresponse, ifnull(grouporganize,'')grouporganize, email, permission, ifnull(lastlogindt,'')lastlogindt, ifnull(logincounter,'')logincounter, ifnull(keycloakUserRole,'')keycloakUserRole FROM raotdb.vw_user_role_backup WHERE CONCAT(username, name, position, department, locationresponse) like '%" + article.SearchTxt + "%'  Order by id  "
	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnUserRAOT{}
	cnt := 0
	if err2 == nil {

		//cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Username, &event.Name, &event.Surename, &event.Position, &event.Department, &event.Locationresponse, &event.Grouporganize, &event.Email, &event.Permission, &event.Lastlogindt, &event.Logincounter, &event.KeycloakUserRole)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserRAOT{DBID: event.DBID, Username: event.Username, Name: event.Name, Surename: event.Surename, Position: event.Position, Department: event.Department, Locationresponse: event.Locationresponse, Grouporganize: event.Grouporganize, Email: event.Email, Permission: event.Permission, Lastlogindt: event.Lastlogindt, Logincounter: event.Logincounter, KeycloakUserRole: event.KeycloakUserRole})

		}

	}
	if cnt == 0 {
		boxes = append(boxes, AppReturnUserRAOT{DBID: "0", Username: "ว่าง", Name: "", Surename: "", Position: "", Department: "", Locationresponse: "", Grouporganize: "", Email: "", Permission: "", Lastlogindt: "", Logincounter: "", KeycloakUserRole: ""})

	}
	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}
func SetRAOTResponselocation(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWAResponseLocation
	json.Unmarshal(reqBody, &article)

	////fmt.Println("UserID", article.CUser)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	////fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetRAOTResponselocation',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///
	//EType := article.EType
	//expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	//strexpiredAt := strconv.FormatInt(expiredAt, 10)

	////fmt.Println(CUser)
	ress3, err2 := db.Query("SELECT dbid, locationresponsename FROM raotdb.tbllocationresponse WHERE dbid = '" + article.CUser + "' ORDER BY dbid  ")

	//defer ress3.Close()
	boxes := []AppReturnResponseLocationRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnResponseLocationRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.LocationResponsename)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: event.DBID, LocationResponsename: event.LocationResponsename})

		}

		if cnt >= 0 {

			if article.EType != "delete" {

				if article.EType == "edit" {
					ress99, err := db.Query("INSERT INTO raotdb.tbllocationresponse (dbid,  locationResponsename, createdt) Values ('" + article.CUser + "','" + article.ResponseLocationName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE locationResponsename = '" + article.ResponseLocationName + "'   ")
					defer ress99.Close()
					if err != nil {
						panic(err)
					}
				}

				if article.EType == "add" {
					ress99, err := db.Query("INSERT INTO raotdb.tbllocationresponse (dbid,  locationResponsename, createdt) SELECT  MAX(dbid)+1 ,'" + article.ResponseLocationName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') FROM raotdb.tbllocationresponse   ")
					defer ress99.Close()
					if err != nil {
						panic(err)
					}
				}

				boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: "1", LocationResponsename: article.ResponseLocationName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			} else {
				ress99, err := db.Query("DELETE FROM raotdb.tbllocationresponse WHERE locationResponsename ='" + article.ResponseLocationName + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: "1", LocationResponsename: article.ResponseLocationName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			}

			//return
			return
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRAOTResponselocation(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	////fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	////fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRAOTResponselocation',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///
	//EType := article.EType

	////fmt.Println(CUser)

	strsql := "SELECT dbid, locationresponsename FROM raotdb.tbllocationresponse ORDER BY dbid  "

	if article.SearchTxt != "" {
		strsql = "SELECT dbid, locationresponsename FROM raotdb.tbllocationresponse WHERE locationresponsename like '%" + article.SearchTxt + "%' ORDER BY dbid  "
	}

	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnResponseLocationRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnResponseLocationRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.LocationResponsename)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: event.DBID, LocationResponsename: event.LocationResponsename})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetRAOTGroup(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	////fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRAOTGroup',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///
	//EType := article.EType

	strsql := "SELECT dbid, groupname as locationresponsename FROM raotdb.tblraotgroup ORDER BY dbid "
	if article.SearchTxt != "" {
		strsql = "SELECT dbid, groupname as locationresponsename FROM raotdb.tblraotgroup WHERE groupname like '%" + article.SearchTxt + "%' ORDER BY dbid "
	}
	//fmt.Println(CUser)
	ress3, err2 := db.Query(strsql)

	//defer ress3.Close()
	boxes := []AppReturnResponseLocationRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnResponseLocationRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.LocationResponsename)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: event.DBID, LocationResponsename: event.LocationResponsename})

		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func SetRAOTGroup(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWAResponseLocation
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.ResponseLocationName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//CUser := article.CUser

	// keep log request
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('SetRAOTGroup',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	///
	//EType := article.EType

	expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
	strexpiredAt := strconv.FormatInt(expiredAt, 10)

	//fmt.Println(CUser)
	ress3, err2 := db.Query("SELECT dbid, groupname as locationresponsename FROM raotdb.tblraotgroup WHERE groupname = '" + article.ResponseLocationName + "' ORDER BY dbid  ")

	//defer ress3.Close()
	boxes := []AppReturnResponseLocationRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnResponseLocationRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.LocationResponsename)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: event.DBID, LocationResponsename: event.LocationResponsename})

		}

		if cnt >= 0 {

			if article.EType != "delete" {
				ress99, err := db.Query("INSERT INTO raotdb.tblraotgroup (dbid,  groupname, createdt) Values ('" + strexpiredAt + "','" + article.ResponseLocationName + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: "1", LocationResponsename: article.ResponseLocationName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			} else {
				ress99, err := db.Query("DELETE FROM raotdb.tblraotgroup WHERE groupname ='" + article.ResponseLocationName + "'   ")
				defer ress99.Close()
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnResponseLocationRAOT{DBID: "1", LocationResponsename: article.ResponseLocationName})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)
			}

			//return
			return
		}

	}

	b, _ := json.Marshal(boxes)

	defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}
func GetRAOTUserPosition(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWASendCheckpointStart
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.WHName)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	CUser := article.PositionName
	//EType := article.EType

	//fmt.Println(CUser)

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetRAOTUserPosition',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	strsql := "SELECT dbid, positionname, ifnull(isdashboard,'')isdashboard, ifnull(ismanageweight,'')ismanageweight, ifnull(ismanagecheckpoint,'')ismanagecheckpoint, ifnull(ischeckpoint,'')ischeckpoint, ifnull(ischeckpointlist,'')ischeckpointlist, ifnull(iscumtomerdata,'')iscumtomerdata, ifnull(istracking,'')istracking, ifnull(isreportlicense,'')isreportlicense, ifnull(isreporttruckchecked,'')isreporttruckchecked, ifnull(isreportweightchecked,'')isreportweightchecked, ifnull(ismanageemployee,'')ismanageemployee, ifnull(ismanageposition,'')ismanageposition, ifnull(ismanagelocationresponse,'')ismanagelocationresponse, ifnull(ismanagegroup,'')ismanagegroup, ifnull(iscondition,'')iscondition, ifnull(ismanagecalculate,'')ismanagecalculate, ifnull(isreporttrader,'')isreporttrader  FROM raotdb.tbluserposition ORDER BY dbid "
	if CUser != "" {
		strsql = "SELECT dbid, positionname, ifnull(isdashboard,'')isdashboard, ifnull(ismanageweight,'')ismanageweight, ifnull(ismanagecheckpoint,'')ismanagecheckpoint, ifnull(ischeckpoint,'')ischeckpoint, ifnull(ischeckpointlist,'')ischeckpointlist, ifnull(iscumtomerdata,'')iscumtomerdata, ifnull(istracking,'')istracking, ifnull(isreportlicense,'')isreportlicense, ifnull(isreporttruckchecked,'')isreporttruckchecked, ifnull(isreportweightchecked,'')isreportweightchecked, ifnull(ismanageemployee,'')ismanageemployee, ifnull(ismanageposition,'')ismanageposition, ifnull(ismanagelocationresponse,'')ismanagelocationresponse, ifnull(ismanagegroup,'')ismanagegroup, ifnull(iscondition,'')iscondition, ifnull(ismanagecalculate,'')ismanagecalculate, ifnull(isreporttrader,'')isreporttrader FROM raotdb.tbluserposition WHERE positionname = '" + CUser + "' ORDER BY dbid "

	}
	if article.SearchTxt != "" {
		strsql = "SELECT dbid, positionname, ifnull(isdashboard,'')isdashboard, ifnull(ismanageweight,'')ismanageweight, ifnull(ismanagecheckpoint,'')ismanagecheckpoint, ifnull(ischeckpoint,'')ischeckpoint, ifnull(ischeckpointlist,'')ischeckpointlist, ifnull(iscumtomerdata,'')iscumtomerdata, ifnull(istracking,'')istracking, ifnull(isreportlicense,'')isreportlicense, ifnull(isreporttruckchecked,'')isreporttruckchecked, ifnull(isreportweightchecked,'')isreportweightchecked, ifnull(ismanageemployee,'')ismanageemployee, ifnull(ismanageposition,'')ismanageposition, ifnull(ismanagelocationresponse,'')ismanagelocationresponse, ifnull(ismanagegroup,'')ismanagegroup, ifnull(iscondition,'')iscondition, ifnull(ismanagecalculate,'')ismanagecalculate, ifnull(isreporttrader,'')isreporttrader FROM raotdb.tbluserposition WHERE positionname like '%" + article.SearchTxt + "%' ORDER BY dbid "

	}

	//ress3, err2 := db.Query("SELECT Wid as Appid, WHName as AppType, OnlineWeight as Version ,CDateTime as LastestDT FROM raotdb.tblweight WHERE weightname  = '" + WHName + "' ")
	//ress3, err2 := db.Query("SELECT weightid as Appid, checkpointcode as AppType, checkpointname as Version ,createdt as LastestDT FROM raotdb.tblcheckpoint WHERE weightid  = '" + WHName + "' ")
	ress3, err2 := db.Query(strsql)

	defer ress3.Close()
	boxes := []AppReturnUserPositionRAOT{}
	if err2 == nil {

		cnt := 0

		for ress3.Next() {

			cnt++

			var event AppReturnUserPositionRAOT
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.DBID, &event.Positionname, &event.Isdashboard, &event.Ismanageweight, &event.Ismanagecheckpoint, &event.Ischeckpoint, &event.Ischeckpointlist, &event.Iscumtomerdata, &event.Istracking, &event.Isreportlicense, &event.Isreporttruckchecked, &event.Isreportweightchecked, &event.Ismanageemployee, &event.Ismanageposition, &event.Ismanagelocationresponse, &event.Ismanagegroup, &event.Iscondition, &event.Ismanagecalculate, &event.Isreporttrader)
			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AppReturnUserPositionRAOT{DBID: event.DBID, Positionname: event.Positionname, Isdashboard: event.Isdashboard, Ismanageweight: event.Ismanageweight, Ismanagecheckpoint: event.Ismanagecheckpoint, Ischeckpoint: event.Ischeckpoint, Ischeckpointlist: event.Ischeckpointlist, Iscumtomerdata: event.Iscumtomerdata, Istracking: event.Istracking, Isreportlicense: event.Isreportlicense, Isreporttruckchecked: event.Isreporttruckchecked, Isreportweightchecked: event.Isreportweightchecked, Ismanageemployee: event.Ismanageemployee, Ismanageposition: event.Ismanageposition, Ismanagelocationresponse: event.Ismanagelocationresponse, Ismanagegroup: event.Ismanagegroup, Iscondition: event.Iscondition, Ismanagecalculate: event.Ismanagecalculate, Isreporttrader: event.Isreporttrader})

		}

	}

	b, _ := json.Marshal(boxes)

	//defer ress3.Close()
	err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)

	//counter = 0

}

func GetSSOAccesstoken(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWALogin
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.UserID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	fmt.Println("ChkAuth", aa)

	if aa < 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	dns := getDNSString(dblogin, userlogin, passlogin, conn) //
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct

	// keep log request

	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('GetSSOAccesstoken',requestcount+1,'" + channel[0] + "','" + article.CUser + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	/////

	MobileType := article.UserID
	VersionNow := article.PassWord
	toKen := article.Token
	//fmt.Println(toKen)

	// var (
	// 	key []byte
	// 	t   *jwt.Token
	// 	s   string
	// )

	// key = []byte("PWASMART")
	// t = jwt.New(jwt.SigningMethodHS256)

	userData := map[string]interface{}{"USERID": MobileType, "VERSION": VersionNow}

	// accessToken, err := Sign(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
	// //fmt.Println(accessToken)
	accessToken := ""

	//s, err = t.SignedString(key)
	//fmt.Println(s)

	/////  SSO Check ////
	url := "http://61.19.236.5/realms/SmartCESS-CheckPoint/protocol/openid-connect/token"
	//url := "https://cess-sso.olive.co.th/realms/SmartCESS-CheckPoint/protocol/openid-connect/token"
	method := "POST"
	//toKen = "7383b48a-df2d-4561-93a4-c42e46f09f6b.5ea3beb9-225c-4a66-930d-a0feed1e865d.379cc971-ac0a-4369-8a25-fc4624f5ea15"

	payload := strings.NewReader("grant_type=authorization_code&redirect_uri=" + MobileType + "Login&code=" + toKen + "&client_id=authen-checkpoint&client_secret=DPCfvNgYsOuS5aBwBinVwnO33VnXpO1W")

	//payload := strings.NewReader("grant_type=authorization_code&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2FKeycloakVerify&code=74589a4b-7f61-40c4-943a-c71b669078db.e302d296-063a-4277-8873-d56b70a6fa6c.379cc971-ac0a-4369-8a25-fc4624f5ea15&client_id=authen-checkpoint&client_secret=DPCfvNgYsOuS5aBwBinVwnO33VnXpO1W")
	//payload := strings.NewReader("grant_type=client_credentials&redirect_uri=https%253A%252F%252Foauth.tools%252Fcallback%252Fcode&code= " + toKen + "&client_id=authen-checkpoint&client_secret=DPCfvNgYsOuS5aBwBinVwnO33VnXpO1W")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	var article2 SSOAccessToken
	json.Unmarshal(body, &article2)

	//fmt.Println("body", article2.AccessToken)
	//fmt.Println("ExpireToken", article2.ExpiresIn)

	tokenString := article2.AccessToken
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("DPCfvNgYsOuS5aBwBinVwnO33VnXpO1W"), nil
	})

	// token2, err := jwt.Parse(tokenString, nil)
	// if token == nil {
	// 	//return
	// }

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		aa := claims["paymentToken"].(string)
		fmt.Println("JWTToken", aa)
	}

	// roles := claims["roles"]
	//fmt.Println("Email", claims["email"])
	//fmt.Println("Role", claims["realm_access"])
	//fmt.Println("Name", claims["name"])

	str_realm_access := fmt.Sprintf("%v", claims["realm_access"])
	fmt.Println("realm_access", str_realm_access)

	//realm_access := claims["realm_access"]
	realm_access := claims["realm_access"]

	b, _ := json.Marshal(realm_access)

	fmt.Println(string(b))

	var rAOTRole RAOTRole
	json.Unmarshal(b, &rAOTRole)

	fmt.Println("Roles", rAOTRole.Roles[2])

	userRole := rAOTRole.Roles[0]
	roleLength := len(rAOTRole.Roles)
	//adminRold := ""

	// if roleLength > 4 {
	// 	adminRold = rAOTRole.Roles[4]
	// }

	for i := 1; i < roleLength; i++ {
		userRole = userRole + "," + rAOTRole.Roles[i]
	}

	//userRole = userRole + "," + adminRold

	exp := fmt.Sprintf("%f", claims["exp"])
	email := claims["email"].(string)

	// role := realm_access["role"]
	// respok["Token"] = accessToken //QR
	// respok["TokenExpireDT"] = strexpiredAt
	// //fmt.Println("JWTToken", role)
	//fmt.Println("realm_access", realm_access)
	fmt.Println("email", email)
	ExpiresIn := strconv.Itoa(article2.ExpiresIn)

	// ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSGetBankAccount','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	// defer ress99.Close()
	// if err != nil {
	// 	panic(err)
	// }

	ressSSO, err := db.Query("INSERT INTO raotdb.tbluser (userid,username,keycloakToken, keycloakExpire,tokenExpiredt,keycloakUserRole,lastlogindt) Values ('" + exp + "','" + email + "','" + article2.AccessToken + "','" + ExpiresIn + "','" + exp + "','" + userRole + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE keycloakToken = '" + article2.AccessToken + "', keycloakExpire = '" + ExpiresIn + "' , tokenExpiredt = '" + exp + "' , keycloakUserRole = '" + userRole + "'   ")
	defer ressSSO.Close()
	if err != nil {
		panic(err)
	}

	// ressSSO, err := db.Query("UPDATE raotdb.tbluser SET keycloakToken = '" + article2.AccessToken + "', keycloakExpire = '" + ExpiresIn + "' , tokenExpiredt = '" + exp + "' , keycloakUserRole = '" + userRole + "'  WHERE username = '" + email + "'  ")
	// defer ressSSO.Close()
	// if err != nil {
	// 	panic(err)
	// }

	//ress3, err2 := db.Query("SELECT userid as AppType, weightid as Version ,tokenExpiredt FROM raotdb.vw_user WHERE username = '" + MobileType + "' ")
	ress3, err2 := db.Query("SELECT AppType, ifnull(Version,'')Version, tokenExpiredt, ifnull(name,'')name, ifnull(surename,'')surename, ifnull(department,'')department, ifnull(locationresponse,'')locationresponse, ifnull(permission,'')permission,ifnull(position,'')position,ifnull(keycloakToken,'')keycloakToken,ifnull(keycloakUserRole,'')Role FROM raotdb.vw_user WHERE username = '" + email + "' ")

	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
		strexpiredAt := strconv.FormatInt(expiredAt, 10)

		boxes := []AppToken{}

		count := 0

		for ress3.Next() {

			count += 1
			var event AppToken
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.AppType, &event.Appid, &event.TkenExpiredt, &event.Name, &event.Surename, &event.Department, &event.Locationresponse, &event.Permission, &event.Position, &event.KeycloakToken, &event.KeyclockRole)
			if err != nil {
				panic(err)
			}

			// ress, err := db.Query("UPDATE raotdb.tbluser SET logincounter = logincounter + 1  WHERE username = '" + MobileType + "'")
			// defer ress.Close()
			// if err != nil {
			// 	panic(err)
			// }

			if toKen != "" {

				f64, err := strconv.ParseFloat(event.TkenExpiredt, 64)
				mf64 := math.Round(f64)
				i64 := int64(mf64)
				// if err != nil {
				// 	accessToken, err := SignRAOT(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
				// 	//fmt.Println(accessToken)

				// 	if err != nil {
				// 		panic(err)
				// 	}

				// 	boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken})
				// 	b, _ := json.Marshal(boxes)

				// 	if err != nil {
				// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				// 		return
				// 	}

				// 	w.Header().Set("Access-Control-Allow-Origin", "*")
				// 	w.Write(b)

				// 	return
				// 	//panic(err)
				// }

				if expiredAt > i64 {
					// accessToken, err := Sign(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
					// //fmt.Println(accessToken)

					// if err != nil {
					// 	panic(err)
					// }

					accessToken = "token_expire"

					boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, User: email, Role: event.Role})
					b, _ := json.Marshal(boxes)

					if err != nil {
						log.Fatalf("Error happened in JSON marshal. Err: %s", err)
						return
					}

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Write(b)

					//return
					return
				}

			} else {
				accessToken, err := SignRAOT(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
				//fmt.Println(accessToken)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, User: email, Role: event.Role})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)

				return
			}

			if event.Appid == VersionNow {
				respok := make(map[string]string)
				respok["AppType"] = event.AppType
				respok["Token"] = accessToken //QR
				respok["TokenExpireDT"] = strexpiredAt
				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return
			}

			//realm_access := claims["realm_access"]
			boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, User: email, Role: event.KeyclockRole})
		}

		if count == 0 {
			accessToken = "invalid_user"

			boxes = append(boxes, AppToken{Appid: "", AppType: "", Token: accessToken})
			b, _ := json.Marshal(boxes)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(b)

			//return
			return
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)

	}
}

func ChkApp(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var article PWALogin
	json.Unmarshal(reqBody, &article)

	//fmt.Println("UserID", article.UserID)

	reqHeader := r.Header["Authorization"]

	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	aa := 0
	auth := []string{article.Authorization}
	channel := []string{article.Channel}
	////fmt.Println("Header", reqHeaderChannel)
	if len(reqHeader) == 0 {

		aa = ChkAuth(auth, channel)
	} else {
		aa = ChkAuth(reqHeader, reqHeaderChannel)
	}

	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	MobileType := article.UserID
	VersionNow := article.PassWord
	toKen := article.Token
	//fmt.Println(toKen)

	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte("PWASMART")
	t = jwt.New(jwt.SigningMethodHS256)

	userData := map[string]interface{}{"USERID": MobileType, "VERSION": VersionNow}

	// accessToken, err := Sign(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
	// //fmt.Println(accessToken)
	accessToken := ""

	s, err = t.SignedString(key)
	fmt.Println(s)

	/////  Keep log count ////
	ress9999, err := db.Query("INSERT INTO raotdb.tblapirepuestlog (apiname, requestcount,channel,lastuser, lastrequestDT) Values ('ChkApp',requestcount+1,'" + channel[0] + "','" + article.UserID + "' ,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )  ON DUPLICATE KEY UPDATE lastrequestDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), requestcount = requestcount + 1 , channel = '" + channel[0] + "' ")
	defer ress9999.Close()
	if err != nil {
		panic(err)
	}
	//////

	/////  SSO Check ////

	//ress3, err2 := db.Query("SELECT userid as AppType, weightid as Version ,tokenExpiredt FROM raotdb.vw_user WHERE username = '" + MobileType + "' ")
	ress3, err2 := db.Query("SELECT AppType, ifnull(Version,'')Version, ifnull(tokenExpiredt,'')tokenExpiredt, ifnull(name,'')name, ifnull(surename,'')surename, ifnull(department,'')department, ifnull(locationresponse,'')locationresponse, ifnull(permission,'')permission,ifnull(position,'')position,ifnull(keycloakToken,'')keycloakToken,ifnull(keycloakUserRole,'')KeyclockRole, discrepancy,apikey FROM raotdb.vw_user WHERE username = '" + MobileType + "' ")

	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		expiredAt := time.Now().Add(time.Duration(time.Second) * 1).Unix()
		strexpiredAt := strconv.FormatInt(expiredAt, 10)

		boxes := []AppToken{}

		count := 0

		for ress3.Next() {

			count += 1
			var event AppToken
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.AppType, &event.Appid, &event.TkenExpiredt, &event.Name, &event.Surename, &event.Department, &event.Locationresponse, &event.Permission, &event.Position, &event.KeycloakToken, &event.KeyclockRole, &event.Discrepancy, &event.APIKey)
			if err != nil {
				panic(err)
			}

			// ress, err := db.Query("UPDATE raotdb.tbluser SET logincounter = logincounter + 1  WHERE username = '" + MobileType + "'")
			// defer ress.Close()
			// if err != nil {
			// 	panic(err)
			// }

			if toKen != "" {

				i64, err := strconv.ParseInt(event.TkenExpiredt, 10, 64)
				if err != nil {
					accessToken, err := SignRAOT(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
					//fmt.Println(accessToken)

					if err != nil {
						panic(err)
					}

					boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, KeyclockRole: event.KeyclockRole, Discrepancy: event.Discrepancy, APIKey: event.APIKey})
					b, _ := json.Marshal(boxes)

					if err != nil {
						log.Fatalf("Error happened in JSON marshal. Err: %s", err)
						return
					}

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Write(b)

					return
					//panic(err)
				}

				if expiredAt > i64 {
					// accessToken, err := Sign(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
					// //fmt.Println(accessToken)

					// if err != nil {
					// 	panic(err)
					// }

					accessToken = "token_expire"

					boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, KeyclockRole: event.KeyclockRole, Discrepancy: event.Discrepancy, APIKey: ""})
					b, _ := json.Marshal(boxes)

					if err != nil {
						log.Fatalf("Error happened in JSON marshal. Err: %s", err)
						return
					}

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Write(b)

					//return
					return
				}

			} else {
				accessToken, err := SignRAOT(userData, "9C0459C63EFB8231A9E23063BF0D413ADCAEFED8D2745B9E1CD5CB070E616BA8", 86400) // data -> secretkey env name -> expiredAt // Production
				//fmt.Println(accessToken)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, KeyclockRole: event.KeyclockRole, Discrepancy: event.Discrepancy, APIKey: event.APIKey})
				b, _ := json.Marshal(boxes)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Write(b)

				return
			}

			if event.Appid == VersionNow {
				respok := make(map[string]string)
				respok["AppType"] = event.AppType
				respok["Token"] = accessToken //QR
				respok["TokenExpireDT"] = strexpiredAt
				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				return
			}

			boxes = append(boxes, AppToken{Appid: event.AppType, AppType: event.Appid, Token: accessToken, Name: event.Name, Surename: event.Surename, Department: event.Department, Locationresponse: event.Locationresponse, Permission: event.Permission, Position: event.Position, KeycloakToken: event.KeycloakToken, KeyclockRole: event.KeyclockRole, Discrepancy: event.Discrepancy, APIKey: ""})
		}

		if count == 0 {
			accessToken = "invalid_user"

			boxes = append(boxes, AppToken{Appid: "", AppType: "", Token: accessToken})
			b, _ := json.Marshal(boxes)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(b)

			//return
			return
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)

	}
}

func OMSGetOrder(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article AccountBank
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	TrackingID := article.TrackingID

	ress3, err2 := db.Query("SELECT JobID, JobName, MerchantID, ProjectName, MerchantWHID, PostCodeWH, PostCodeTHPD, JobDmsID, JobDMSNet, JobVolumeWeightNet, JobQty, JobWeightID, JobWeightNet, JobGrossWeightNet, JobProductType, JobType, JobStatus, JobStatusName, JobDriverID, JobDriverName, JobTruckID, JobTruckPlate, JobOpenDT, JobCallReceiveDT, JobPickupDT, JobAssignDT, JobAppointmentDT, JobBranchReceiveDT, JobBranchAssignDT, JobPendingDT, JobSuspendDT, JobCancelDT, JobCompleteDT, JobRateNet, LTLDiscount, ContactName, ContactPhone, Remark, IsActual, IsRequestPrint, MsgToDriver, IsSurvey, TruckTypeID, FTLStartPrice, FTLAddServices, FTLAddSeason, FTLStandardPrice, FTLSpecialPrice, FTLExpectPrice, FTLFinalPrice, FTLRateOilKM, FTLDistanceKM, FTLOriginGPS, FTLDestinationGPS, FTLAddOnName, FTLDiscount, FTLBidType, FTLInsuranceFee, FTLInsurancePrice, DropPointCount, TMSSendStatus, TMSSendDT, JobCustomerRead, JobOpenRead, JobAssignRead, IsBookmark, JobUpdateDT, JobCreateDT, Customer_Po_1, Customer_Po, Order_No, TrackingID, UserID, WarehouseID, Customer_No, Customer_Name, Receive_name, Receive_Address, Receive_Tel, Receive_Country, Receive_City, Receive_State, Receive_Postal_Code, Receive_GPS, Item, Item_Description, Product_Category, Original_Product_Category, Qty, UnitType, Weight, Width, Length, Height, Ref_1, Ref_2, Ref_3, MessageToSeller, THPDRatePrice, Saved, PickupStatus, PickupTruckID, PickupTruckPlate, PickupDriverID, PickupDriverName, PickupDT, IsCOD, AmountCOD, IsInsurance, AmountPrice, ItemNo, DropPointNo, ItemStatus, PreDeliveryID, PreDeliveryDT, RemarkID, RemarkDesc ,IMG1 FROM THPDMPDB.vw_jobordermaster WHERE JobID = '" + TrackingID + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []ReOrder{}

		for ress3.Next() {
			var event ReOrder
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.JobID, &event.JobName, &event.MerchantID, &event.ProjectName, &event.MerchantWHID, &event.PostCodeWH, &event.PostCodeTHPD, &event.JobDmsID, &event.JobDMSNet, &event.JobVolumeWeightNet, &event.JobQty, &event.JobWeightID, &event.JobWeightNet, &event.JobGrossWeightNet, &event.JobProductType, &event.JobType, &event.JobStatus, &event.JobStatusName, &event.JobDriverID, &event.JobDriverName, &event.JobTruckID, &event.JobTruckPlate, &event.JobOpenDT, &event.JobCallReceiveDT, &event.JobPickupDT, &event.JobAssignDT, &event.JobAppointmentDT, &event.JobBranchReceiveDT, &event.JobBranchAssignDT, &event.JobPendingDT, &event.JobSuspendDT, &event.JobCancelDT, &event.JobCompleteDT, &event.JobRateNet, &event.LTLDiscount, &event.ContactName, &event.ContactPhone, &event.Remark, &event.IsActual, &event.IsRequestPrint, &event.MsgToDriver, &event.IsSurvey, &event.TruckTypeID, &event.FTLStartPrice, &event.FTLAddServices, &event.FTLAddSeason, &event.FTLStandardPrice, &event.FTLSpecialPrice, &event.FTLExpectPrice, &event.FTLFinalPrice, &event.FTLRateOilKM, &event.FTLDistanceKM, &event.FTLOriginGPS, &event.FTLDestinationGPS, &event.FTLAddOnName, &event.FTLDiscount, &event.FTLBidType, &event.FTLInsuranceFee, &event.FTLInsurancePrice, &event.DropPointCount, &event.TMSSendStatus, &event.TMSSendDT, &event.JobCustomerRead, &event.JobOpenRead, &event.JobAssignRead, &event.IsBookmark, &event.JobUpdateDT, &event.JobCreateDT, &event.Customer_Po_1, &event.Customer_Po, &event.Order_No, &event.TrackingID, &event.UserID, &event.WarehouseID, &event.Customer_No, &event.Customer_Name, &event.Receive_name, &event.Receive_Address, &event.Receive_Tel, &event.Receive_Country, &event.Receive_City, &event.Receive_State, &event.Receive_Postal_Code, &event.Receive_GPS, &event.Item, &event.Item_Description, &event.Product_Category, &event.Original_Product_Category, &event.Qty, &event.UnitType, &event.Weight, &event.Width, &event.Length, &event.Height, &event.Ref_1, &event.Ref_2, &event.Ref_3, &event.MessageToSeller, &event.THPDRatePrice, &event.Saved, &event.PickupStatus, &event.PickupTruckID, &event.PickupTruckPlate, &event.PickupDriverID, &event.PickupDriverName, &event.PickupDT, &event.IsCOD, &event.AmountCOD, &event.IsInsurance, &event.AmountPrice, &event.ItemNo, &event.DropPointNo, &event.ItemStatus, &event.PreDeliveryID, &event.PreDeliveryDT, &event.RemarkID, &event.RemarkDesc, &event.IMG1)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, ReOrder{JobID: event.JobID, JobName: event.JobName, MerchantID: event.MerchantID, ProjectName: event.ProjectName, MerchantWHID: event.MerchantWHID, PostCodeWH: event.PostCodeWH, PostCodeTHPD: event.PostCodeTHPD, JobDmsID: event.JobDmsID, JobDMSNet: event.JobDMSNet, JobVolumeWeightNet: event.JobVolumeWeightNet, JobQty: event.JobQty, JobWeightID: event.JobWeightID, JobWeightNet: event.JobWeightNet, JobGrossWeightNet: event.JobGrossWeightNet, JobProductType: event.JobProductType, JobType: event.JobType, JobStatus: event.JobStatus, JobStatusName: event.JobStatusName, JobDriverID: event.JobDriverID, JobDriverName: event.JobDriverName, JobTruckID: event.JobTruckID, JobTruckPlate: event.JobTruckPlate, JobOpenDT: event.JobOpenDT, JobCallReceiveDT: event.JobCallReceiveDT, JobPickupDT: event.JobPickupDT, JobAssignDT: event.JobAssignDT, JobAppointmentDT: event.JobAppointmentDT, JobBranchReceiveDT: event.JobBranchReceiveDT, JobBranchAssignDT: event.JobBranchAssignDT, JobPendingDT: event.JobPendingDT, JobSuspendDT: event.JobSuspendDT, JobCancelDT: event.JobCancelDT, JobCompleteDT: event.JobCompleteDT, JobRateNet: event.JobRateNet, LTLDiscount: event.LTLDiscount, ContactName: event.ContactName, ContactPhone: event.ContactPhone, Remark: event.Remark, IsActual: event.IsActual, IsRequestPrint: event.IsRequestPrint, MsgToDriver: event.MsgToDriver, IsSurvey: event.IsSurvey, TruckTypeID: event.TruckTypeID, FTLStartPrice: event.FTLStartPrice, FTLAddServices: event.FTLAddServices, FTLAddSeason: event.FTLAddSeason, FTLStandardPrice: event.FTLStandardPrice, FTLSpecialPrice: event.FTLSpecialPrice, FTLExpectPrice: event.FTLExpectPrice, FTLFinalPrice: event.FTLFinalPrice, FTLRateOilKM: event.FTLRateOilKM, FTLDistanceKM: event.FTLDistanceKM, FTLOriginGPS: event.FTLOriginGPS, FTLDestinationGPS: event.FTLDestinationGPS, FTLAddOnName: event.FTLAddOnName, FTLDiscount: event.FTLDiscount, FTLBidType: event.FTLBidType, FTLInsuranceFee: event.FTLInsuranceFee, FTLInsurancePrice: event.FTLInsurancePrice, DropPointCount: event.DropPointCount, TMSSendStatus: event.TMSSendStatus, TMSSendDT: event.TMSSendDT, JobCustomerRead: event.JobCustomerRead, JobOpenRead: event.JobOpenRead, JobAssignRead: event.JobAssignRead, IsBookmark: event.IsBookmark, JobUpdateDT: event.JobUpdateDT, JobCreateDT: event.JobCreateDT, Customer_Po_1: event.Customer_Po_1, Customer_Po: event.Customer_Po, Order_No: event.Order_No, TrackingID: event.TrackingID, UserID: event.UserID, WarehouseID: event.WarehouseID, Customer_No: event.Customer_No, Customer_Name: event.Customer_Name, Receive_name: event.Receive_name, Receive_Address: event.Receive_Address, Receive_Tel: event.Receive_Tel, Receive_Country: event.Receive_Country, Receive_City: event.Receive_City, Receive_State: event.Receive_State, Receive_Postal_Code: event.Receive_Postal_Code, Receive_GPS: event.Receive_GPS, Item: event.Item, Item_Description: event.Item_Description, Product_Category: event.Product_Category, Original_Product_Category: event.Original_Product_Category, Qty: event.Qty, UnitType: event.UnitType, Weight: event.Weight, Width: event.Width, Length: event.Length, Height: event.Height, Ref_1: event.Ref_1, Ref_2: event.Ref_2, Ref_3: event.Ref_3, MessageToSeller: event.MessageToSeller, THPDRatePrice: event.THPDRatePrice, Saved: event.Saved, PickupStatus: event.PickupStatus, PickupTruckID: event.PickupTruckID, PickupTruckPlate: event.PickupTruckPlate, PickupDriverID: event.PickupDriverID, PickupDriverName: event.PickupDriverName, PickupDT: event.PickupDT, IsCOD: event.IsCOD, AmountCOD: event.AmountCOD, IsInsurance: event.IsInsurance, AmountPrice: event.AmountPrice, ItemNo: event.ItemNo, DropPointNo: event.DropPointNo, ItemStatus: event.ItemStatus, PreDeliveryID: event.PreDeliveryID, PreDeliveryDT: event.PreDeliveryDT, RemarkID: event.RemarkID, RemarkDesc: event.RemarkDesc, IMG1: event.IMG1})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}
}
func OMSUpdateBookMark(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article AccountBank
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	//TrackingID := article.TrackingID

	ress, err := db.Query("UPDATE THPDMPDB.tblJobMaster SET IsBookmark = '" + article.UserType + "' WHERE JobID = '" + article.TrackingID + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	} else {
		boxes := []AmountCOD{}
		boxes = append(boxes, AmountCOD{TrackingID: article.TrackingID, AmountCOD: article.UserType})
		b, _ := json.Marshal(boxes)

		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
	}

}
func OMSGetCOD(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article AccountBank
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	TrackingID := article.TrackingID

	ress3, err2 := db.Query("SELECT TrackingID, AmountCOD FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + TrackingID + "' ")
	//defer ress3.Close()
	if err2 == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AmountCOD{}

		for ress3.Next() {
			var event AmountCOD
			//JobID := ress2.Scan(&event.JobID)
			err := ress3.Scan(&event.TrackingID, &event.AmountCOD)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AmountCOD{TrackingID: event.TrackingID, AmountCOD: event.AmountCOD})
		}

		b, _ := json.Marshal(boxes)

		defer ress3.Close()
		err = ress3.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}
}
func OMSGetBankAccount(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article AccountBank
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	if article.TaxID == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""               //QR
		respok["ErrMsg"] = "Invalid TaxID" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.BankCode == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                  //QR
		respok["ErrMsg"] = "Invalid BankCode" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.BankAccNo == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                   //QR
		respok["ErrMsg"] = "Invalid BankAccNo" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.ReceiptAddress == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                        //QR
		respok["ErrMsg"] = "Invalid ReceiptAddress" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.ReceiptDistrict == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                         //QR
		respok["ErrMsg"] = "Invalid ReceiptDistrict" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.ReceiptProvince == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                         //QR
		respok["ErrMsg"] = "Invalid ReceiptProvince" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if article.ReceiptZipcode == "" {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                        //QR
		respok["ErrMsg"] = "Invalid ReceiptZipcode" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSGetBankAccount','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	TrackingID := article.TrackingID
	AccountID := article.AccountID
	//AccountName := article.AccountName
	//PhoneNo := article.PhoneNo
	BankName := article.BankName
	BankAccNo := article.BankAccNo
	BankAccName := article.BankAccName
	//BankType := article.BankType
	PromptPayID := article.PromptPayID
	//	UserType := article.UserType
	PrimaryTransferAccount := article.PrimaryTransferAccount
	AmountPrice := article.AmountPrice
	Comcode := article.ComCode

	ReceiptAddress := article.ReceiptAddress
	ReceiptSubDistrict := article.ReceiptSubDistrict
	ReceiptDistrict := article.ReceiptDistrict
	ReceiptProvince := article.ReceiptProvince
	ReceiptZipcode := article.ReceiptZipcode

	receiptaddress := ReceiptAddress + "|" + ReceiptSubDistrict + "|" + ReceiptDistrict + "|" + ReceiptProvince + "|" + ReceiptZipcode

	//Corporation := article.Corporation
	//AccTransferPrimary := article.AccTransferPrimary
	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date

	////fmt.Println(article.Id)

	// firstEvent := Event5{}
	// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
	// if err != nil {
	// 	panic(err)
	// }

	// 	AccName           string
	// AccBankName        string
	// AccNo              string
	// AccPromptPay       string
	// AccTransferPrimary string
	// AccType            string

	ress2, err := db.Query("SELECT TrackingID,ComCode FROM THPDMPDB.tblpaymentaccount WHERE TrackingID = '" + TrackingID + "' and Comcode = '" + Comcode + "'  ")

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []AccountBank{}

		for ress2.Next() {
			var event AccountBank
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.TrackingID, &event.ComCode)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, AccountBank{TrackingID: event.TrackingID, ComCode: event.ComCode})

			resperr := make(map[string]string)
			resperr["errmsg"] = "Transaction is Duplicate in " + event.TrackingID + " ComCode : " + event.ComCode

			jsonResp, err := json.Marshal(resperr)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
	}

	//counter := 0
	// sqlstr := "INSERT INTO THPDMPDB.tblpaymentaccount (TrackingID, AccountID, AccountName, PhoneNo, Email, BankName, BankAccNo, BankType, PromptPayID, CreateDT, UserType, PrimaryTransferAccount,AmountPrice,ComCode ) Values ('" + article.TrackingID + "','" + article.AccountID + "','" + article.AccountName + "','" + article.PhoneNo + "','" + article.Email + "','" + article.BankName + "','" + article.BankAccNo + "','" + article.BankType + "','" + article.PromptPayID + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + article.UserType + "','" + article.PrimaryTransferAccount + "', '" + article.AmountPrice + "', '" + article.ComCode + "')"
	// //fmt.Println(sqlstr)
	ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblpaymentaccount (TrackingID, AccountID, AccountName, PhoneNo, Email, BankName, BankAccNo, BankAccName, BankType, PromptPayID, CreateDT, UserType, PrimaryTransferAccount,AmountPrice,ComCode,TaxID,BankCode,Ref3,Ref1 ) Values ('" + article.TrackingID + "','" + article.AccountID + "','" + article.AccountName + "','" + article.PhoneNo + "','" + article.Email + "','" + article.BankName + "','" + article.BankAccNo + "','" + article.BankAccName + "','" + article.BankType + "','" + article.PromptPayID + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'" + article.UserType + "','" + article.PrimaryTransferAccount + "', '" + article.AmountPrice + "', '" + article.ComCode + "', '" + article.TaxID + "', '" + article.BankCode + "','" + article.Corporation + "' ,'" + receiptaddress + "' )")
	//defer ress3.Close()
	if err2 != nil {

		resperr := make(map[string]string)
		resperr["errmsg"] = err2.Error()

		jsonResp, err := json.Marshal(resperr)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
		//panic(err)
	}
	defer ress3.Close()
	err = ress3.Close()

	if err2 != nil {
		//panic(err)
		resperr := make(map[string]string)
		resperr["errmsg"] = err.Error()

		jsonResp, err := json.Marshal(resperr)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	} else {
		// respok := make(map[string]string)
		// respok["Success"] = "Insert Success"
		// respok["TrackingID"] = TrackingID
		// respok["AmountPrice"] = AmountPrice
		// respok["AccountID"] = AccountID
		// respok["BankAccNo"] = BankAccNo
		// respok["BankName"] = BankName
		// respok["PromptPayID"] = PromptPayID

		respok := map[string]interface{}{
			"Success": true,
			"Data": struct {
				TrackingID             string
				AmountPrice            string
				AccountID              string
				BankAccNo              string
				BankAccName            string
				BankName               string
				PromptPayID            string
				PrimaryTransferAccount string
				ComCode                string
			}{TrackingID, AmountPrice, AccountID, BankAccNo, BankAccName, BankName, PromptPayID, PrimaryTransferAccount, Comcode},
		}

		jsonResp, err := json.Marshal(respok)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return

	}

	//ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber,MemberDiscount FROM THPDMPDB.tblMobileOMSMemberConfig WHERE PhoneNumber = '" + article.MobileID + "'  ")
	// ress2, err := db.Query("SELECT TrackingID, AccountID, AccountName, PhoneNo, Email, BankName, BankAccNo, BankType, PromptPayID, CreateDT, UserType, PrimaryTransferAccount,AmountPrice FROM THPDMPDB.tblpaymentaccount WHERE PhoneNo = '" + article.PhoneNo + "'  ")

	// if err == nil {

	// 	// boxes := []BoxData{}
	// 	// boxes = append(boxes, BoxData{Width: 10, Height: 20})
	// 	// boxes = append(boxes, BoxData{Width: 5, Height: 30})

	// 	boxes := []Mobile{}

	// 	for ress2.Next() {
	// 		var event Mobile
	// 		//JobID := ress2.Scan(&event.JobID)
	// 		err := ress2.Scan(&event.MobileID, &event.Lang, &event.BG, &event.AccName, &event.AccBankName, &event.AccNo, &event.AccPromptPay, &event.AccTransferPrimary, &event.AccType)

	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		boxes = append(boxes, Mobile{MobileID: event.MobileID, Lang: event.Lang, BG: event.BG, AccName: event.AccName, AccBankName: event.AccBankName, AccNo: event.AccNo, AccPromptPay: event.AccPromptPay, AccTransferPrimary: event.AccTransferPrimary, AccType: event.AccType})
	// 	}

	// 	b, _ := json.Marshal(boxes)

	// 	defer ress2.Close()
	// 	err = ress2.Close()
	// 	//defer ress2.Close()
	// 	//jsonResp, err := json.Marshal(b)
	// 	if err != nil {
	// 		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 		return
	// 	}
	// 	w.Write(b)
	//counter = 0

}
func OMSMobileUpdateLoginMember(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// rootCertPool := x509.NewCertPool()
	// pem, err := ioutil.ReadFile("BaltimoreCyberTrustRoot.crt.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
	// 	log.Fatal("Failed to append PEM.")
	// }
	// mysql.RegisterTLSConfig("custom", &tls.Config{
	// 	ServerName: conn,
	// 	RootCAs:    rootCertPool,
	// })
	// db, err := sql.Open("mysql", "admin:!10<>Oms!@tcp(prd-olivedb.mysql.database.azure.com:3306)/THPDMPDB?tls=custom")

	// rootCertPool := x509.NewCertPool()
	// pem, _ := ioutil.ReadFile("BaltimoreCyberTrustRoot.crt.pem")
	// if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
	// 	log.Fatal("Failed to append PEM.")
	// }
	// mysql.RegisterTLSConfig("TLS_AES_256_GCM_SHA384", &tls.Config{RootCAs: rootCertPool})
	// var connectionString string
	// connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=TLS_AES_256_GCM_SHA384", "admin", "!10<>Oms!", conn, "THPDMPDB")
	// db, err := sql.Open("mysql", connectionString)

	//db, err := sql.Open("mysql", "user=admin password=!10<>Oms! dbname=THPDMPDB sslmode=disable")
	//db, err := sql.Open("mysql", "mysql://admin:!10<>Oms!@prd-olivedb.mysql.database.azure.com/THPDMPDB?sslmode=require")

	dns := getDNSString("mysql", "admin", "!10<>Oms!", conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MobileID
	typeEvent := article.TypeEvent
	typelang := article.Lang
	bg := article.BG
	LoginType := article.LoginType
	AccName := article.AccName
	AccBankName := article.AccBankName
	AccNo := article.AccNo
	AccType := article.AccType
	AccPromptPay := article.AccPromptPay
	AccTransferPrimary := article.AccTransferPrimary
	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		// 	AccName            string
		// AccBankName        string
		// AccNo              string
		// AccPromptPay       string
		// AccTransferPrimary string
		// AccType            string

		//counter := 0

		if typeEvent == "Login" {
			ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblMobileOMSMemberConfig (MemberPhone, CreateDT, LoginCount, LastLoginDT, Lang, BG, MerchantID,LoginType,AccName,AccBankName,AccNo,AccType,AccPromptPay,AccTransferPrimary) Values ('" + article.MobileID + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),1,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'th','White','c','" + LoginType + "','','','','','','') ON DUPLICATE KEY UPDATE LastLoginDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') , LoginCount = LoginCount + 1 ,LoginType = '" + LoginType + "'")
			defer ress3.Close()
			if err2 != nil {
				panic(err)
			}
			defer ress3.Close()
			err = ress3.Close()

		} else if typeEvent == "GetLang" {

		} else if typeEvent == "SetLang" && AccName != "" {

			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSMemberConfig SET Lang = '" + typelang + "' , BG = '" + bg + "', AccName = '" + AccName + "' ,AccBankName = '" + AccBankName + "',AccNo = '" + AccNo + "',AccPromptPay = '" + AccPromptPay + "' ,AccType = '" + AccType + "',AccTransferPrimary = '" + AccTransferPrimary + "'  WHERE MemberPhone = '" + article.MobileID + "'  ")
			if err != nil {
				panic(err)
			}
			defer ress2.Close()
			err = ress2.Close()
		} else if typeEvent == "SetLang" && AccName == "" {

			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSMemberConfig SET Lang = '" + typelang + "' , BG = '" + bg + "'  WHERE MemberPhone = '" + article.MobileID + "'  ")
			if err != nil {
				panic(err)
			}
			defer ress2.Close()
			err = ress2.Close()
		} else {
			return
		}
	}

	//ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber,MemberDiscount FROM THPDMPDB.tblMobileOMSMemberConfig WHERE PhoneNumber = '" + article.MobileID + "'  ")
	ress2, err := db.Query("SELECT MemberPhone, Lang, BG,IFNULL(AccName,'')AccName, IFNULL(AccBankName,'')AccBankName,IFNULL(AccNo,'')AccNo,IFNULL(AccPromptPay,'')AccPromptPay,IFNULL(AccTransferPrimary,'')AccTransferPrimary,IFNULL(AccType,'')AccType FROM THPDMPDB.tblMobileOMSMemberConfig WHERE MemberPhone = '" + article.MobileID + "'  ")

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []Mobile{}

		for ress2.Next() {
			var event Mobile
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.MobileID, &event.Lang, &event.BG, &event.AccName, &event.AccBankName, &event.AccNo, &event.AccPromptPay, &event.AccTransferPrimary, &event.AccType)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Mobile{MobileID: event.MobileID, Lang: event.Lang, BG: event.BG, AccName: event.AccName, AccBankName: event.AccBankName, AccNo: event.AccNo, AccPromptPay: event.AccPromptPay, AccTransferPrimary: event.AccTransferPrimary, AccType: event.AccType})
		}

		b, _ := json.Marshal(boxes)

		defer ress2.Close()
		err = ress2.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}

}
func OMSMobileUpdateLoginMemberWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	// rootCertPool := x509.NewCertPool()
	// pem, err := ioutil.ReadFile("BaltimoreCyberTrustRoot.crt.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
	// 	log.Fatal("Failed to append PEM.")
	// }
	// mysql.RegisterTLSConfig("custom", &tls.Config{
	// 	ServerName: conn,
	// 	RootCAs:    rootCertPool,
	// })
	// db, err := sql.Open("mysql", "admin:!10<>Oms!@tcp(prd-olivedb.mysql.database.azure.com:3306)/THPDMPDB?tls=custom")

	// rootCertPool := x509.NewCertPool()
	// pem, _ := ioutil.ReadFile("BaltimoreCyberTrustRoot.crt.pem")
	// if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
	// 	log.Fatal("Failed to append PEM.")
	// }
	// mysql.RegisterTLSConfig("TLS_AES_256_GCM_SHA384", &tls.Config{RootCAs: rootCertPool})
	// var connectionString string
	// connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=TLS_AES_256_GCM_SHA384", "admin", "!10<>Oms!", conn, "THPDMPDB")
	// db, err := sql.Open("mysql", connectionString)

	//db, err := sql.Open("mysql", "user=admin password=!10<>Oms! dbname=THPDMPDB sslmode=disable")
	//db, err := sql.Open("mysql", "mysql://admin:!10<>Oms!@prd-olivedb.mysql.database.azure.com/THPDMPDB?sslmode=require")

	dns := getDNSString("mysql", "admin", "!10<>Oms!", conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MobileID
	typeEvent := article.TypeEvent
	typelang := article.Lang
	bg := article.BG
	LoginType := article.LoginType
	AccName := article.AccName
	AccBankName := article.AccBankName
	AccNo := article.AccNo
	AccType := article.AccType
	AccPromptPay := article.AccPromptPay
	AccTransferPrimary := article.AccTransferPrimary
	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		// 	AccName            string
		// AccBankName        string
		// AccNo              string
		// AccPromptPay       string
		// AccTransferPrimary string
		// AccType            string

		//counter := 0

		if typeEvent == "Login" {
			ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblMobileOMSMemberConfig (MemberPhone, CreateDT, LoginCount, LastLoginDT, Lang, BG, MerchantID,LoginType,AccName,AccBankName,AccNo,AccType,AccPromptPay,AccTransferPrimary) Values ('" + article.MobileID + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),1,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'th','White','c','" + LoginType + "','','','','','','') ON DUPLICATE KEY UPDATE LastLoginDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') , LoginCount = LoginCount + 1 ,LoginType = '" + LoginType + "'")
			defer ress3.Close()
			if err2 != nil {
				panic(err)
			}
			defer ress3.Close()
			err = ress3.Close()

		} else if typeEvent == "GetLang" {

		} else if typeEvent == "SetLang" && AccName != "" {

			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSMemberConfig SET Lang = '" + typelang + "' , BG = '" + bg + "', AccName = '" + AccName + "' ,AccBankName = '" + AccBankName + "',AccNo = '" + AccNo + "',AccPromptPay = '" + AccPromptPay + "' ,AccType = '" + AccType + "',AccTransferPrimary = '" + AccTransferPrimary + "'  WHERE MemberPhone = '" + article.MobileID + "'  ")
			if err != nil {
				panic(err)
			}
			defer ress2.Close()
			err = ress2.Close()
		} else if typeEvent == "SetLang" && AccName == "" {

			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSMemberConfig SET Lang = '" + typelang + "' , BG = '" + bg + "'  WHERE MemberPhone = '" + article.MobileID + "'  ")
			if err != nil {
				panic(err)
			}
			defer ress2.Close()
			err = ress2.Close()
		} else {
			return
		}
	}

	//ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber,MemberDiscount FROM THPDMPDB.tblMobileOMSMemberConfig WHERE PhoneNumber = '" + article.MobileID + "'  ")
	ress2, err := db.Query("SELECT MemberPhone, Lang, BG,IFNULL(AccName,'')AccName, IFNULL(AccBankName,'')AccBankName,IFNULL(AccNo,'')AccNo,IFNULL(AccPromptPay,'')AccPromptPay,IFNULL(AccTransferPrimary,'')AccTransferPrimary,IFNULL(AccType,'')AccType FROM THPDMPDB.tblMobileOMSMemberConfig WHERE MemberPhone = '" + article.MobileID + "'  ")

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []Mobile{}

		for ress2.Next() {
			var event Mobile
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.MobileID, &event.Lang, &event.BG, &event.AccName, &event.AccBankName, &event.AccNo, &event.AccPromptPay, &event.AccTransferPrimary, &event.AccType)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Mobile{MobileID: event.MobileID, Lang: event.Lang, BG: event.BG, AccName: event.AccName, AccBankName: event.AccBankName, AccNo: event.AccNo, AccPromptPay: event.AccPromptPay, AccTransferPrimary: event.AccTransferPrimary, AccType: event.AccType})
		}

		b, _ := json.Marshal(boxes)

		defer ress2.Close()
		err = ress2.Close()
		//defer ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}

}
func OMSMobileGetMember(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)
	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	// crc32q := crc32.MakeTable(0xFFFF)
	// fmt.Printf("%08x\n", crc32.Checksum([]byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"), crc32q))

	// CRC Checksum test
	data2 := []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304")
	target2 := uint16(0x0808)
	hex_value := strconv.FormatInt(int64(target2), 16)

	actual2 := ChecksumCCITTFalse(data2)
	if actual2 != target2 {
		fmt.Printf("%08x\n", actual2)
	}
	if hex_value != "" {
		//t.Fatalf("CCITT checksum did not return the correct value, expected %x, received %x", target, actual)
	}
	// test crc
	// tab := crc32.MakeTable(0x1021)

	// var crc uint32 = 0xffffffff
	// crc = crc32.Update(crc, tab, []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"))
	// fmt.Printf("%08x\n", crc)

	// var crc2 uint32 = 0xffffffff
	// crc2 = crc32c_le(crc2, tab, []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"))
	// fmt.Printf("%08x\n", crc2)

	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber,MemberDiscount FROM THPDMPDB.tblMerchant WHERE PhoneNumber = '" + article.MobileID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []EventMobile{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event EventMobile
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.MerchantName, &event.Email, &event.ContactName, &event.PhoneNumber, &event.MemberDiscount)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventMobile{MerchantID: event.MerchantID, ContactName: event.ContactName, Email: event.Email, MerchantName: event.MerchantName, PhoneNumber: event.PhoneNumber, MemberDiscount: event.MemberDiscount})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}

}
func CheckUser(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article KTBApproveJson
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	} else {
		// dns := getDNSString(dblogin, userlogin, passlogin, conn)
		// db, err := sql.Open("mysql", dns)
		// if err != nil {
		// 	panic(err)
		// }
		// err = db.Ping()
		// if err != nil {
		// 	panic(err)
		// }
		// defer db.Close()
		if DB.Ping() != nil {
			connectDb()
		}
		if DB.Stats().OpenConnections != 0 {
			//fmt.Println(DB.Stats().OpenConnections)
		} else {
			connectDb()
		}
		db := DB
		defer db.Close()
		// err := db.Ping()
		// if err != nil {
		// 	connectDb()

		// }

		User := article.User
		Pass := article.Password
		data := User + Pass
		sEnc := b64.StdEncoding.EncodeToString([]byte(data))

		ress2, err := db.Query("SELECT Cid, ChannelName, ChannelSecretKey FROM THPDMPDB.tblauthchannel WHERE ChannelName = '" + User + "' and ChannelSecretKey ='" + sEnc + "'")

		if err != nil {
			return
		}

		havedata := 0

		for ress2.Next() {
			var event AuthDevice
			//JobID := ress2.Scan(&event.JobID)
			err = ress2.Scan(&event.Cid, &event.ChannelName, &event.ChannelSecretKey)

			if err != nil {
				panic(err)
			}
			//mobileid = event.MobileID
			//SendSMSFirstJob = event.SendSMSFirstJob
			ress, err := db.Query("UPDATE  THPDMPDB.tblAuthChannel Set UseCount = UseCount + 1, LastUseDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  WHERE ChannelName = '" + User + "'  ")
			defer ress.Close()
			if err != nil {
				panic(err)
			} else {

			}

			havedata++
		}

		defer ress2.Close()
		err = ress2.Close()

		respok := make(map[string]string)

		if havedata == 1 {
			respok["QRcodetxt"] = ""
			respok["ExpDT"] = ""         //QR
			respok["ErrMsg"] = "Success" //QR
		} else {
			respok["QRcodetxt"] = ""
			respok["ExpDT"] = ""            //QR
			respok["ErrMsg"] = "Login Fail" //QR
		}

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return

	}
}
func OMSMobileGetMemberWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	// crc32q := crc32.MakeTable(0xFFFF)
	// fmt.Printf("%08x\n", crc32.Checksum([]byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"), crc32q))

	// CRC Checksum test
	data2 := []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304")
	target2 := uint16(0x0808)
	hex_value := strconv.FormatInt(int64(target2), 16)

	actual2 := ChecksumCCITTFalse(data2)
	if actual2 != target2 {
		fmt.Printf("%08x\n", actual2)
	}
	if hex_value != "" {
		//t.Fatalf("CCITT checksum did not return the correct value, expected %x, received %x", target, actual)
	}
	// test crc
	// tab := crc32.MakeTable(0x1021)

	// var crc uint32 = 0xffffffff
	// crc = crc32.Update(crc, tab, []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"))
	// fmt.Printf("%08x\n", crc)

	// var crc2 uint32 = 0xffffffff
	// crc2 = crc32c_le(crc2, tab, []byte("00020101021230770016A00000067701011201150105557004580010220100900000370300000710310B000003703530376454076000.005802TH6304"))
	// fmt.Printf("%08x\n", crc2)

	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber,MemberDiscount FROM THPDMPDB.tblMerchant WHERE PhoneNumber = '" + article.MobileID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []EventMobile{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event EventMobile
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.MerchantID, &event.MerchantName, &event.Email, &event.ContactName, &event.PhoneNumber, &event.MemberDiscount)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventMobile{MerchantID: event.MerchantID, ContactName: event.ContactName, Email: event.Email, MerchantName: event.MerchantName, PhoneNumber: event.PhoneNumber, MemberDiscount: event.MemberDiscount})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}

}
func OMSMobileGetPostcodeWithAuth(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Zipcode
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.PostId

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		strsql := "SELECT PostID,Tambon_TH as Tambon_Name,Tambon_EN,Amphoe_TH as Amphoe_Name,Amphoe_EN,Province_TH as Province_Name, GPS FROM THPDDB.VW_List_Location_SubDistrict WHERE Postid = '" + article.PostId + "'"

		if article.PostId == "00000" {

			strsql = "SELECT PostID,Tambon_TH as Tambon_Name,Tambon_EN,Amphoe_TH as Amphoe_Name,Amphoe_EN,Province_TH as Province_Name, GPS FROM THPDDB.VW_List_Location_SubDistrict ORDER BY Province_TH,Amphoe_TH,Tambon_TH "

		}

		ress2, err := db.Query(strsql)

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Zipcode{}
			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["PostId"] = article.PostId

			for ress2.Next() {
				var event Zipcode
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.PostId, &event.Tumbon, &event.TumbonEN, &event.District, &event.DistrictEN, &event.Province, &event.GPS)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Zipcode{PostId: event.PostId, Tumbon: event.Tumbon, TumbonEN: event.TumbonEN, District: event.District, DistrictEN: event.DistrictEN, Province: event.Province, GPS: event.GPS})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}

}
func OMSMobileGetPostcode(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Zipcode
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.PostId

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		strsql := "SELECT PostID,Tambon_TH as Tambon_Name,Tambon_EN,Amphoe_TH as Amphoe_Name,Amphoe_EN,Province_TH as Province_Name, GPS FROM THPDDB.VW_List_Location_SubDistrict WHERE Postid = '" + article.PostId + "'"

		if article.PostId == "00000" {

			strsql = "SELECT PostID,Tambon_TH as Tambon_Name,Tambon_EN,Amphoe_TH as Amphoe_Name,Amphoe_EN,Province_TH as Province_Name, GPS FROM THPDDB.VW_List_Location_SubDistrict "

		}

		ress2, err := db.Query(strsql)

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Zipcode{}
			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["PostId"] = article.PostId

			for ress2.Next() {
				var event Zipcode
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.PostId, &event.Tumbon, &event.TumbonEN, &event.District, &event.DistrictEN, &event.Province, &event.GPS)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Zipcode{PostId: event.PostId, Tumbon: event.Tumbon, TumbonEN: event.TumbonEN, District: event.District, DistrictEN: event.DistrictEN, Province: event.Province, GPS: event.GPS})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}

}
func OMSMobileSelectOrder(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT * FROM THPDMPDB.tblMobileOMSOrder WHERE MobileID = '" + article.MobileID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event Event5
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.IMG1, &event.IMG2, &event.IMG3, &event.IMG4, &event.CreateDt, &event.JobDesc, &event.WarehouseID, &event.MerchantID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail, IMG1: event.IMG1, IMG2: event.IMG2, IMG3: event.IMG3, IMG4: event.IMG4,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}

func OMSMobileExtraRate(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MobileVersion
	json.Unmarshal(reqBody, &article)

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()

	strsql := "SELECT RateID,RateName,RatePrice FROM thpddb.tblrateextra "

	ress2, err := db.Query(strsql)
	//fmt.Println(err)

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []RateExtra{}
		//boxes = append(boxes, BoxData{Width: 5, Height: 30})

		resp := make(map[string]string)
		resp["citi"] = article.VersionNow

		for ress2.Next() {
			var event RateExtra
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.RateID, &event.RateName, &event.RatePrice)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, RateExtra{RateID: event.RateID, RateName: event.RateName, RatePrice: event.RatePrice})
		}

		b, _ := json.Marshal(boxes)

		defer ress2.Close()
		err = ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}
	//var m MyStruct
	//typefind := article.MobileType

	// respok := make(map[string]string)
	// respok["mobile"] = article.MobileType
	// respok["response"] = "keep" //QR

	// jsonResp, err := json.Marshal(respok)

	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	return
	// }
	// w.Write(jsonResp)
	// return

	//return

}

func OMSMobileChkMerchant(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MobileVersion
	json.Unmarshal(reqBody, &article)

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()

	strsql := "SELECT PhoneNumber,CitizenID FROM thpdmpdb.tblmerchant WHERE PhoneNumber = '" + article.MobileType + "'"

	ress2, err := db.Query(strsql)
	//fmt.Println(err)

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []Citicen{}
		//boxes = append(boxes, BoxData{Width: 5, Height: 30})

		resp := make(map[string]string)
		resp["citi"] = article.VersionNow

		for ress2.Next() {
			var event Citicen
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.MobileId, &event.CiticenId)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Citicen{MobileId: event.MobileId, CiticenId: substr(event.CiticenId, 9, 4)})
		}

		b, _ := json.Marshal(boxes)

		defer ress2.Close()
		err = ress2.Close()
		//jsonResp, err := json.Marshal(b)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(b)
		//counter = 0

	}
	//var m MyStruct
	//typefind := article.MobileType

	// respok := make(map[string]string)
	// respok["mobile"] = article.MobileType
	// respok["response"] = "keep" //QR

	// jsonResp, err := json.Marshal(respok)

	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	return
	// }
	// w.Write(jsonResp)
	// return

	//return

}

func OMSUpdateTMS(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MobileVersion
	json.Unmarshal(reqBody, &article)

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// if err != nil {
	// 	connectDb()
	// }

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()

	trackingID := ""
	ress2, err := db.Query("SELECT TrackingID,ComCode FROM THPDMPDB.tblpaymentaccount WHERE IDs = '" + article.MobileType + "'  ")

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})
		// boxes := []AccountBank{}

		for ress2.Next() {
			var event AccountBank
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.TrackingID, &event.ComCode)
			if err != nil {
				panic(err)
			}
			trackingID = event.TrackingID
		}

	}

	if trackingID == "" {
		return
	}

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = article.MobileType

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	t := time.Now() //It will return time.Time object with current timestamp
	//t = t.Add(time.Hour * 7) เอาออกเนื่องจาก dx บวก 7 แล้ว

	tUnix := t.Unix()
	fmt.Printf("timeUnix: %d\n", tUnix)

	var n int64 = tUnix

	s := strconv.FormatInt(n, 10)

	fmt.Printf("%d\n", tUnix)
	//strOTP := strconv.FormatInt(tUnix, 6)

	//payload := strings.NewReader(`{"con_no": "` + trackID + `" ,"carrier_id": "` + carrier_id + `"}`)
	// status 86 อยู่ระหว่างตรวจสอบ
	// status 87 อนุมัติการการขอเบิกค่าขนส่ง
	// status 88 ไม่อนุมัติ
	// status 89 โอนเรียบร้อย
	status := article.VersionNow
	if article.VersionNow == "0" {
		status = "86" // กำลังตรวจสอบ
	} else if article.VersionNow == "1" {
		status = "86" // กำลังตรวจสอบ
	} else if article.VersionNow == "2" {
		status = "87" // อนุมิติ
	} else if article.VersionNow == "-2" {
		status = "88" // ไม่อนุมัติ
	} else if article.VersionNow == "3" {
		status = "89" //โอนเงินสำเร็จแล้ว
	}

	payload := strings.NewReader(`{
		"shipment_no": "` + trackingID + `",
		"status": "` + status + `",
		"actual_time": ` + s + `,
		"action_time": ` + s + `,
	
		"doc_link": ""
	  }`)

	//  check LTL ให้ส่ง trackid เป็น LT...DB

	spayload := `{
		"shipment_no": "` + trackingID + `",
		"status": "` + status + `",
		"actual_time": ` + s + `,
		"action_time": ` + s + `,
	
		"doc_link": ""
	  }`
	//  "jobs": [{ "con_no": "` + trackID + `" }],
	method := "POST"
	url := tmsapi + "API/BrokerGateways/updateStatus" // UAT
	//url := "https://tms-api.promptsong.co/API/BrokerGateways/updateStatus" //prodution

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)
	req.Header.Add("Content-Type", "application/json")

	resp2, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		//return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		//fmt.Println(err)
		//return string(err.Error())
	}
	//fmt.Println(string(body))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackingID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + ", " + spayload + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	//var m MyStruct
	//typefind := article.MobileType

	respok := make(map[string]string)
	respok["trackingID"] = trackingID
	respok["response"] = "Update TMS Payment Success : Status " + status //QR

	jsonResp, err := json.Marshal(respok)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.Write(jsonResp)
	return

	//return

}
func OMSMobileKeepLog(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article MobileVersion
	json.Unmarshal(reqBody, &article)

	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// if err != nil {
	// 	connectDb()

	// }

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	err := db.Ping()

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + article.MobileType + "' , 'Mobile Api Log', '" + strings.ReplaceAll(string(article.VersionNow), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}
	//var m MyStruct
	//typefind := article.MobileType

	respok := make(map[string]string)
	respok["mobile"] = article.MobileType
	respok["response"] = "keep" //QR

	jsonResp, err := json.Marshal(respok)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.Write(jsonResp)
	return

	//return

}
func OMSMobileDiscount(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT * FROM THPDMPDB.tblMobileOMSOrder WHERE MobileID = '" + article.MobileID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.TrackingID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, TrackingID: event.TrackingID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}

func OMSMobileGetJobDriverBookingMatchAlready(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)
	//", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, Receivename, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, Status, PickupStartDt, TrackingID, CarID, DriverName, DriverPhoneNo, Price,CustomerSelect,WarehouseName,WHAddressTH FROM THPDMPDB.VW_OMSMobileDriverSelectJobAlready WHERE TrackingID = '" + article.TrackingID + "'  ")

		if err == nil {

			boxes := []EventBookingMatchAlready{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event EventBookingMatchAlready
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.Receivename, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.Status, &event.PickupStartDt, &event.TrackingID, &event.CarID, &event.DriverName, &event.DriverPhoneNo, &event.Price, &event.CustomerSelect, &event.WarehouseName, &event.WHAddressTH)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventBookingMatchAlready{JobID: event.JobID, MobileID: event.MobileID, Receivename: event.Receivename, ReceiveAddress: event.ReceiveAddress, ReceiveTumbon: event.ReceiveTumbon, ReceiveDistrict: event.ReceiveDistrict, ReceiveProvince: event.ReceiveZipcode, ReceiveZipcode: event.ReceiveProvince, ReceivePhoneNo: event.ReceivePhoneNo, Status: event.Status, PickupStartDt: event.PickupStartDt, TrackingID: event.TrackingID, CarID: event.CarID, DriverName: event.DriverName, DriverPhoneNo: event.DriverPhoneNo, Price: event.Price, CustomerSelect: event.CustomerSelect, WarehouseName: event.WarehouseName, WHAddressTH: event.WHAddressTH})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileGetJobDriverBookingMatchAlreadyWithAuth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileGetJobDriverBookingMatchAlreadyWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, Receivename, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, Status, PickupStartDt, TrackingID, CarID, DriverName, DriverPhoneNo, Price,CustomerSelect,WarehouseName,WHAddressTH FROM THPDMPDB.VW_OMSMobileDriverSelectJobAlready WHERE TrackingID = '" + article.TrackingID + "'  ")

		if err == nil {

			boxes := []EventBookingMatchAlready{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event EventBookingMatchAlready
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.Receivename, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.Status, &event.PickupStartDt, &event.TrackingID, &event.CarID, &event.DriverName, &event.DriverPhoneNo, &event.Price, &event.CustomerSelect, &event.WarehouseName, &event.WHAddressTH)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventBookingMatchAlready{JobID: event.JobID, MobileID: event.MobileID, Receivename: event.Receivename, ReceiveAddress: event.ReceiveAddress, ReceiveTumbon: event.ReceiveTumbon, ReceiveDistrict: event.ReceiveDistrict, ReceiveProvince: event.ReceiveZipcode, ReceiveZipcode: event.ReceiveProvince, ReceivePhoneNo: event.ReceivePhoneNo, Status: event.Status, PickupStartDt: event.PickupStartDt, TrackingID: event.TrackingID, CarID: event.CarID, DriverName: event.DriverName, DriverPhoneNo: event.DriverPhoneNo, Price: event.Price, CustomerSelect: event.CustomerSelect, WarehouseName: event.WarehouseName, WHAddressTH: event.WHAddressTH})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileGetJobDriverBooking(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT BookingID, Oid, DriverName ,Company,DriverStar as Rating,CarID,Price FROM THPDMPDB.tblMobileOMSJobDriverBooking WHERE TrackingID = '" + article.TrackingID + "' and BookingID is not null ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []EventBooking{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event EventBooking
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.BookingID, &event.Oid, &event.DriverName, &event.Company, &event.Rating, &event.TruckPlate, &event.Price)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventBooking{BookingID: event.BookingID, Oid: event.Oid, DriverName: event.DriverName, Company: event.Company, Rating: event.Rating, TruckPlate: event.TruckPlate, Price: event.Price})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileGetJobDriverBookingWithAuth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT BookingID, Oid, DriverName ,Company,DriverStar as Rating,CarID,Price FROM THPDMPDB.tblMobileOMSJobDriverBooking WHERE TrackingID = '" + article.TrackingID + "' and BookingID is not null ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []EventBooking{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event EventBooking
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.BookingID, &event.Oid, &event.DriverName, &event.Company, &event.Rating, &event.TruckPlate, &event.Price)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, EventBooking{BookingID: event.BookingID, Oid: event.Oid, DriverName: event.DriverName, Company: event.Company, Rating: event.Rating, TruckPlate: event.TruckPlate, Price: event.Price})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderWithIMG(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, WarehouseID, MerchantID,IMG1 FROM THPDMPDB.VW_OMSMobileOrder WHERE TrackingID = '" + article.TrackingID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []OMSOrderStruct{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event OMSOrderStruct
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.WarehouseID, &event.MerchantID, &event.IMG1)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, OMSOrderStruct{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, IMG1: event.IMG1})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderWithIMGWithAuth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)
	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, WarehouseID, MerchantID,IMG1 FROM THPDMPDB.VW_OMSMobileOrder WHERE TrackingID = '" + article.TrackingID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []OMSOrderStruct{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event OMSOrderStruct
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.WarehouseID, &event.MerchantID, &event.IMG1)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, OMSOrderStruct{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, IMG1: event.IMG1})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderFromBookingWithAuth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, JobType, WarehouseID, MerchantID, Status,TrackingID FROM THPDMPDB.VW_OMSMobileOrderNoIMG2 WHERE jobtype like '%" + article.TrackingID + "%' ORDER BY JobID DESC")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.Status, &event.TrackingID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, Status: event.Status, TrackingID: event.TrackingID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderFromBooking(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.TrackingID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, JobType, WarehouseID, MerchantID, Status,TrackingID FROM THPDMPDB.VW_OMSMobileOrderNoIMG2 WHERE jobtype like '%" + article.TrackingID + "%' ORDER BY JobID DESC")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.Status, &event.TrackingID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, Status: event.Status, TrackingID: event.TrackingID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}

func OMSCheckCarrierGetJob(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()

	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSCheckCarrierGetJob','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, JobType, WarehouseID, MerchantID, Status,TrackingID FROM THPDMPDB.VW_OMSMobileOrderNoIMG2 WHERE MobileID = '" + article.MobileID + "' and Status like '%ผู้ขนส่งขอเสนอราคา%' ORDER BY JobID DESC")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.Status, &event.TrackingID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, Status: event.Status, TrackingID: event.TrackingID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderNoIMGWithAuth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)

	reqHeader := r.Header["Authorization"]
	////fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	////fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	//fmt.Println("aa", aa)

	if aa <= 0 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                                  //QR
		respok["ErrMsg"] = "Invalid Authorization Or Channel" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// if DB.Ping() != nil {
	// 	connectDb()
	// }
	// if DB.Stats().OpenConnections != 0 {
	// 	//fmt.Println(DB.Stats().OpenConnections)
	// } else {
	// 	connectDb()
	// }
	// db := DB
	// defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }

	ress99, err := db.Query("INSERT INTO THPDMPDB.tblmobileapiloghit (ApiName, HitCount, LastHitDT) Values ('OMSMobileSelectOrderNoIMGWithAuth','1',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') ) ON DUPLICATE KEY UPDATE LastHitDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), HitCount = HitCount + 1   ")
	defer ress99.Close()
	if err != nil {
		panic(err)
	}
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, JobType, WarehouseID, MerchantID, Status,TrackingID,IMG1,IMG2,IFNULL(IMG3,'')IMG3,IFNULL(IMG4,'')IMG4 FROM THPDMPDB.VW_OMSMobileOrderNoIMG2 WHERE MobileID = '" + article.MobileID + "' ORDER BY IMG4*1 desc, JOBID desc")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.Status, &event.TrackingID, &event.IMG1, &event.IMG2, &event.IMG3, &event.IMG4)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, Status: event.Status, TrackingID: event.TrackingID, IMG1: event.IMG1, IMG2: event.IMG2, IMG3: event.IMG3, IMG4: event.IMG4})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}
func OMSMobileSelectOrderNoIMG(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	// dns := getDNSString(dblogin, userlogin, passlogin, conn)
	// db, err := sql.Open("mysql", dns)

	// if err != nil {
	// 	panic(err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	if DB.Ping() != nil {
		connectDb()
	}
	if DB.Stats().OpenConnections != 0 {
		//fmt.Println(DB.Stats().OpenConnections)
	} else {
		connectDb()
	}
	db := DB
	defer db.Close()
	// err := db.Ping()
	// if err != nil {
	// 	connectDb()

	// }
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)
	//firstEvent2 := Event5{}
	//datett := article.Date
	if typefind != "" {
		////fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		ress2, err := db.Query("SELECT JobID, MobileID, ReceiveName, ReceiveAddress, ReceiveTumbon, ReceiveDistrict, ReceiveProvince, ReceiveZipcode, ReceivePhoneNo, SenderZipcode, SendPrice, PickupStartDt, DeliveryEndDt, PaymentFlag, PaymentDetail, CreateDt, JobDesc, JobType, WarehouseID, MerchantID, Status,TrackingID FROM THPDMPDB.VW_OMSMobileOrderNoIMG2 WHERE MobileID = '" + article.MobileID + "' ORDER BY JobID DESC")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []Event5noIMG{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			for ress2.Next() {
				var event Event5noIMG
				//JobID := ress2.Scan(&event.JobID)
				err := ress2.Scan(&event.JobID, &event.MobileID, &event.ReceiveName, &event.ReceiveAddress, &event.ReceiveTumbon, &event.ReceiveDistrict, &event.ReceiveProvince, &event.ReceiveZipcode, &event.ReceivePhoneNo, &event.SenderZipcode, &event.SendPrice, &event.PickupStartDt, &event.DeliveryEndDt, &event.PaymentFlag, &event.PaymentDetail, &event.CreateDt, &event.JobDesc, &event.JobType, &event.WarehouseID, &event.MerchantID, &event.Status, &event.TrackingID)

				if err != nil {
					panic(err)
				}

				boxes = append(boxes, Event5noIMG{JobID: event.JobID, MobileID: event.MobileID, ReceiveName: event.ReceiveName, ReceiveAddress: event.ReceiveAddress,
					ReceiveTumbon: event.ReceiveTumbon, ReceiveProvince: event.ReceiveProvince, ReceiveDistrict: event.ReceiveDistrict,
					ReceivePhoneNo: event.ReceivePhoneNo, ReceiveZipcode: event.ReceiveZipcode,
					SenderZipcode: event.SenderZipcode, SendPrice: event.SendPrice, PickupStartDt: event.PickupStartDt, DeliveryEndDt: event.DeliveryEndDt,
					PaymentFlag: event.PaymentFlag, PaymentDetail: event.PaymentDetail,
					CreateDt: event.CreateDt, JobDesc: event.JobDesc, JobType: event.JobType, WarehouseID: event.WarehouseID, MerchantID: event.MerchantID, Status: event.Status, TrackingID: event.TrackingID})
			}

			b, _ := json.Marshal(boxes)

			defer ress2.Close()
			err = ress2.Close()
			//jsonResp, err := json.Marshal(b)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(b)
			//counter = 0

		}

	}
	//return

}

func API(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.

	// reqBody, _ := ioutil.ReadAll(r.Body)
	// var article APIIot
	// json.Unmarshal(reqBody, &article)

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	bs := string(body)
	//fmt.Println(bs)

	res1 := strings.Split(bs, "&")
	//fmt.Println(strings.Split(res1[4], "=")[1])

	res41 := strings.Split(res1[4], "=")[1]

	//ioutil.WriteFile("dump", body, 0600)
	// update our global Articles array to include
	// our new Article
	//Articles = append(Articles, article)

	//json.NewEncoder(w).Encode(article)

	//w.Write(reqBody)

	dns := getDNSString(dblogin, userlogin, passlogin, conn)
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// //CUser := article.CUser
	// //EType := article.EType
	// ////fmt.Println(CUser)

	boxes := []AppReturnTruckEntry{}

	if res41 == "chk" {

		ress3, err2 := db.Query("SELECT configdetail,ref1,ref2,ref3,ref4 FROM raotdb.tblconfig WHERE configdetail = 'truckentry'  Order by id  ")

		//defer ress3.Close()

		if err2 == nil {

			cnt := 0

			for ress3.Next() {

				cnt++

				var event AppReturnTruckEntry
				//JobID := ress2.Scan(&event.JobID)
				err := ress3.Scan(&event.Configdetail, &event.Ref1, &event.Ref2, &event.Ref3, &event.Ref4)
				if err != nil {
					panic(err)
				}

				boxes = append(boxes, AppReturnTruckEntry{Configdetail: event.Configdetail, Ref1: event.Ref1, Ref2: event.Ref2, Ref3: event.Ref3, Ref4: event.Ref4})

			}

		}

	} else {

		ress99, err := db.Query("UPDATE raotdb.tblconfig SET  ref1=(ref1*1)+1, getlastdt = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'), ref4='" + res41 + "' WHERE configdetail = 'truckentry' ")
		defer ress99.Close()
		if err != nil {
			panic(err)
		}

		boxes = append(boxes, AppReturnTruckEntry{Configdetail: "", Ref1: "", Ref2: "", Ref3: "", Ref4: ""})

	}

	b, _ := json.Marshal(boxes)

	// defer ress3.Close()
	// err = ress3.Close()
	//defer ress2.Close()
	//jsonResp, err := json.Marshal(b)
	// if err != nil {
	// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	// 	return
	// }

	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	}
	// }
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	// w.Header().Set("Connection", "Close")
	w.Write(b)
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
	orderid = 0
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)

	/// RAOT API GATEWAY///

	myRouter.HandleFunc("/ChkApp", ChkApp).Methods("POST")                                         //test PWA Checkin
	myRouter.HandleFunc("/SendWeight", SendWeight).Methods("POST")                                 //test PWA SendWeight
	myRouter.HandleFunc("/GetWeight", GetWeight).Methods("POST")                                   //test PWA GetWeight
	myRouter.HandleFunc("/ConfigWeight", ConfigWeight).Methods("POST")                             //test PWA ConfigWeight
	myRouter.HandleFunc("/GetConfigWeight", GetConfigWeight).Methods("POST")                       //test PWA ConfigWeight
	myRouter.HandleFunc("/GetCheckPointLocation", GetCheckPointLocation).Methods("POST")           //test PWA ConfigWeight
	myRouter.HandleFunc("/ConfigCheckPointLocation", ConfigCheckPointLocation).Methods("POST")     //test PWA ConfigCheckPointLocation
	myRouter.HandleFunc("/GetCheckPoint", GetCheckPoint).Methods("POST")                           //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetCheckPointTranDetail", GetCheckPointTranDetail).Methods("POST")       //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetCheckPointTranDetailGPS", GetCheckPointTranDetailGPS).Methods("POST") //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetCessTran", GetCessTran).Methods("POST")                               //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetLicenseIMGBase64", GetLicenseIMGBase64).Methods("POST")               //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetLicenseIMGBase64LPR", GetLicenseIMGBase64LPR).Methods("POST")         //test PWA GetCheckPoint
	myRouter.HandleFunc("/UpdateTransDetail", UpdateTransDetail).Methods("POST")                   //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetIMGBase64", GetIMGBase64).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetCalculate", GetCalculate).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetUpdateOverLicense", GetUpdateOverLicense).Methods("POST")             //test PWA GetCheckPoint
	myRouter.HandleFunc("/ConfigStartCheckPoint", ConfigStartCheckPoint).Methods("POST")           //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetStartCheckPoint", GetStartCheckPoint).Methods("POST")                 //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetCalculate", SetCalculate).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetTrader", GetTrader).Methods("POST")                                   //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetSSOAccesstoken", GetSSOAccesstoken).Methods("POST")                   //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRAOTUser", GetRAOTUser).Methods("POST")                               //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRAOTUserBackup", GetRAOTUserBackup).Methods("POST")                   //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetRAOTUser", SetRAOTUser).Methods("POST")                               //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetRAOTUserRollback", SetRAOTUserRollback).Methods("POST")               //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRAOTUserPosition", GetRAOTUserPosition).Methods("POST")               //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetRAOTUserPosition", SetRAOTUserPosition).Methods("POST")               //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRAOTResponselocation", GetRAOTResponselocation).Methods("POST")       //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetRAOTResponselocation", SetRAOTResponselocation).Methods("POST")       //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRAOTGroup", GetRAOTGroup).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetRAOTGroup", SetRAOTGroup).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetJobDetail", SetJobDetail).Methods("POST")                             //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetConsent", SetConsent).Methods("POST")                                 //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetTruck", GetTruck).Methods("POST")                                     //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRepTraderTruck", GetRepTraderTruck).Methods("POST")                   //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRepTrader", GetRepTrader).Methods("POST")                             //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRepRubberType", GetRepRubberType).Methods("POST")                     //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetRepWeightChecked", GetRepWeightChecked).Methods("POST")               //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetGPSLicense", GetGPSLicense).Methods("POST")                           //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetDatatoCumtomer", SetDatatoCumtomer).Methods("POST")                   //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/SndDatatoCESS1", SndDatatoCESS1).Methods("POST")                         //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/APITruckCheckInOut", APITruckCheckInOut).Methods("POST")                 //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetNoti", GetNoti).Methods("POST")                                       //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/SetWidget", SetWidget).Methods("POST")                                   //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/GetDashboard", GetDashboard).Methods("POST")                             //from API p'wat  //test PWA GetCheckPoint
	myRouter.HandleFunc("/SndNotiTracking", SndNotiTracking).Methods("POST")                       //from API p'wat  //test PWA GetCheckPoint

	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)

	//myRouter.HandleFunc("/webhook", handleWebhook)
	//myRouter.HandleFunc("/GetImage", GetImage).Methods("POST")
	//myRouter.HandleFunc("/API", API)

	//http.Handle("/", v3.NewHandler("My API", "/swagger.json", "/"))
	//<<<<<<< HEAD
	//log.Fatal(http.ListenAndServe(":8080", myRouter)) // production RAOT
	//log.Fatal(http.ListenAndServe(":80", myRouter)) // production
	log.Fatal(http.ListenAndServe(":8081", myRouter)) // test
	//log.Fatal(http.ListenAndServe(":9080", myRouter)) // test
	//=======
	//log.Fatal(http.ListenAndServe(":80", myRouter)) // production
	//log.Fatal(http.ListenAndServe(":8081", myRouter)) // test
	//log.Fatal(http.ListenAndServe(":5678", myRouter)) // test
	//>>>>>>> c475cd4d25c7033671cff6ae7260c045e921dcef

	// fs := http.FileServer(http.Dir("dist"))
	// http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
}

func main() {

	connectDb()

	//http.ListenAndServe(":8080", nil)

	Articles = []Article{
		// Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		// Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Articcle Content", Token: "xxx"},
		Article{Id: "1", User: "User", Token: "xxx", Type: "InvoiceNo", Date: "2022-02-22T14:20:00"},
		Article{Id: "2", User: "User", Token: "xxx", Type: "CustomerNo", Date: "2022-02-22T14:20:00"},
		Article{Id: "3", User: "User", Token: "xxx", Type: "InvoiceNo", Date: "2022-02-22T14:20:00"},
	}
	OMSOrderStructs = []OMSOrderStruct{
		OMSOrderStruct{JobID: "1", MobileID: "User", Token: "xxx", ReceiveName: "ReceiveName", ReceiveAddress: "ReceiveAddress",
			ReceiveTumbon: "ReceiveTumbon", ReceiveDistrict: "ReceiveDistrict", ReceiveProvince: "ReceiveProvince", ReceivePhoneNo: "ReceivePhoneNo",
			ReceiveZipcode: "ReceiveZipcode", SenderZipcode: "SenderZipcode", SendPrice: "SendPrice", PickupStartDt: "PickupStartDt",
			DeliveryEndDt: "DeliveryEndDt", PaymentFlag: "PaymentFlag", PaymentDetail: "PaymentDetail", IMG1: "IMG1",
			IMG2: "IMG2", IMG3: "IMG3", IMG4: "IMG4", CreateDt: "CreateDt", MerchantID: "MerchantID", WarehouseID: "WarehouseID"},
	}
	Mobiles = []Mobile{
		// Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		// Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Articcle Content", Token: "xxx"},
		Mobile{MobileID: "1", Password: "xxxx", TypeEvent: "Login", Lang: "en", BG: "#ffffff", LoginType: "Finglescan"},
	}
	Merchants = []Merchant{
		// Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		// Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Articcle Content", Token: "xxx"},
		Merchant{MerchantID: "1"},
	}
	MerchantReps = []MerchantRep{
		// Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content", Token: "xxx"},
		// Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Articcle Content", Token: "xxx"},
		MerchantRep{MerchantID: "1", RepType: "Daily", SDT: "2022-09-01", EDT: "2022-09-22"},
	}
	DriverTrackings = []DriverTracking{
		DriverTracking{TrackingID: "1", DriverID: "xxx", DriverName: "xx", Plate: "xx", FinalPrice: "xx"},
	}
	DriverJobMasterBookings = []DriverJobMasterBooking{
		DriverJobMasterBooking{TrackingID: "1", JobDriverID: "xxx", JobDriverName: "xx", JobTruckID: "xx", JobTruckPlate: "xx"},
	}

	Get2c2ppaymentTokens = []Get2c2ppaymentToken{
		Get2c2ppaymentToken{TrackingID: "1", Amount: "100", PhoneNumber: "0809999999", PaymentType: "QR", ComCode: "1001", TransportationPrice: "0", ProductPrice: "0", NumberOfPieces: "0"},
	}
	// GetPostPrintOnlines = []GetPostPrintOnline{
	// GetPostPrintOnline{PostCode: "10010", CountItems: 1},
	// }
	GetCancels = []GetCancel{
		GetCancel{Shipment_no: "B000003xxx", Remark: "xxx"},
	}

	GetAccounts = []AccountBank{
		AccountBank{TrackingID: "B000003xxx", AccountID: "xxx", AccountName: "xxx", PhoneNo: "xxx", Email: "xxx", BankName: "xxx", BankAccNo: "xxx", BankAccName: "xxx", BankType: "xxx", PromptPayID: "xxx", UserType: "xxx", PrimaryTransferAccount: "xxx", AmountPrice: "xxx", ComCode: "1001"},
	}

	GetKTBApproveJson = []KTBApproveJson{
		KTBApproveJson{User: "B000003xxx", Password: "xxx", ComCode: "", ProdCode: "", Command: "", BankRef: "", DateTime: "", EffDate: "", Amount: 0.00, Channel: "", CusName: "", Ref1: "", Ref2: "", Ref3: "", Ref4: ""},
	}

	GetMobileVersion = []MobileVersion{
		MobileVersion{MobileType: "B000003xxx", VersionNow: "xxx"},
	}

	log.Println("server started")
	//http.HandleFunc("/webhook", handleWebhook)

	//log.Fatal(http.ListenAndServe(":8081", nil))

	handleRequests()

}
