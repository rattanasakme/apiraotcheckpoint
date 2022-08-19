// main.go
package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
type DriverLoadboard struct {
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
type Mobile struct {
	MobileID  string
	Password  string
	TypeEvent string
	Lang      string
	BG        string
	LoginType string
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
type Merchant struct {
	MerchantID   string
	Address1     string
	Address2     string
	LocationGPS  string
	WarehouseID  string
	SubDistrict  string
	District     string
	ProvinceName string
	PostCode     string
}
type Coupon struct {
	CouponName     string
	CouponCode     string
	CouponDiscount string
	CouponExpireDT string
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
type UrlAdvise struct {
	UrlText string
}
type MobileID struct {
	MobileID        string
	SendSMSFirstJob int
	TrackingID      string
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
	TrackingID  string
	Amount      string
	PhoneNumber string
	PaymentType string
}
type GetPostPrintOnline struct {
	PostCode   string
	CountItems int
}
type TblTMSItemEvents struct {
	DBID        string `json:"DBID"`
	TrackingID  string `json:"TrackingID"`
	Status_Name string `json:"Status_Name"`
	CreateDT    string `json:"CreateDT"`
	PaymentFlag string `json:"PaymentFlag"`
	OfferPrice  string `json:"OfferPrice"`
}

// type response1 struct {
// 	MobileID       string
// 	OMSOrderStruct []string
// }
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
	flagKey     = flag.String("key", "CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2", "path to key file or '-' to read from stdin")
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

var (
	APIName = ""
	//selectEventByIdQueryTOAT    = `SELECT * FROM THPDTOATDB.VW_TOAT_Track WHERE InvoiceNo = ? and InvoiceNo <> '' and Truck_type <> ''  limit 1`
	//selectEventByCusIdQueryTOAT = `SELECT * FROM THPDTOATDB.VW_TOAT_Track WHERE CustomerNo = ? and InvoiceDate = ? and InvoiceNo <> '' and Truck_type <> ''  limit 1`
	selectEventByToken  = `SELECT Email,MerchantName FROM THPDMPDB.tblMerchant WHERE MerchantID = ? and APIAuthenKey = ? limit 1`
	insertOrderItem     = `INSERT INTO THPDMPDB.tblMobileOMSOrder(JobID,MobileID,ReceiveName,ReceiveAddress,ReceiveTumbon,ReceiveDistrict,ReceiveProvince,ReceiveZipcode,ReceivePhoneNo,SenderZipcode,SendPrice,PickupStartDt,DeliveryEndDt,PaymentFlag,PaymentDetail,CreateDt,JobDesc,JobType,WarehouseID,MerchantID,Status) values (  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),?,?,?,?,? )`
	insertOrderItemDesc = `INSERT INTO THPDMPDB.tblMobileOMSOrderDesc(JobID,MobileID,IMG1,IMG2,IMG3,IMG4,CreateDt) values ( ?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )`
	insertGetDriver     = `INSERT INTO THPDMPDB.tblMobileOMSJobDriverBooking(JobID, TrackingID, CarID, DriverName, DriverPhoneNo, Company, Price, CustomerSelect, CustomConfirmDT, GetBeforeDT,CreateDT) values ( ?, ?, ?, ?,?, ?, ?, ?, ?, ?, CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') )`
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
var DriverTrackings []DriverTracking
var DriverJobMasterBookings []DriverJobMasterBooking
var Get2c2ppaymentTokens []Get2c2ppaymentToken
var GetPostPrintOnlines []GetPostPrintOnline

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

	fmt.Printf("ReplyToken: \n", ReplyToken)
	fmt.Printf("MsgText: \n", MsgText)

	matched, err := regexp.MatchString(`AX([0-9]+)TX`, MsgText)
	fmt.Println(matched)

	//var bearerToken = "YBAnQHA776xhZVp5IiPGWG556tuo5hU3h0Ke0afe2LBi2pCAymiz0fgLcjLAEeWQZ+9+oJkjH1RFYyA676vOmk78/CCb7Bgns1eSm7GRrGf7GKlwGhp944byiEQZbhV5X1QXpAQMXSM0nN/zusI6yAdB04t89/1O/w1cDnyilFU=" line oongang
	var bearerToken = "cpQfHAyvCAB1gAnYn7elBKCo8DQJhxsa+tWVSWz0EbkkTbFQbkm47EZgoeWwMZQrzSXzvsbN7Ay9MK/ET7EZQMcob5ZuL3FmW6zO5WTMpv++kVZmI+/VV6Z82IbOkLQkATft7vFh8Hca48OdMsMhUAdB04t89/1O/w1cDnyilFU="
	var myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"ระบบบริหารแพลตฟอร์มขนส่ง THPD DE ยินดีให้บริการ "},{"type":"text","text":"ท่านสามารถค้นหาสิ่งของโดยพิมพ์หมายเลข TrackID 13 หลัก"}]}`)

	if matched {

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

		chkdigit := 0

		if tall == 0 {
			chkdigit = 5
		} else if tall == 1 {
			chkdigit = 0
		} else {
			chkdigit = 11 - int(tall)
		}
		chkdigit2 := int(t9)
		if int(chkdigit) != chkdigit2 {

			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" หมายเลข TrackID ของท่านไม่ถูกต้อง! "}]}`)

			//return
		} else {

			//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
			dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

			db, err := sql.Open("mysql", dns)

			if err != nil {
				panic(err)
			}
			err = db.Ping()
			if err != nil {
				panic(err)
			}
			count := 0
			ress2, err := db.Query("SELECT  DBID,TrackingID,Status_Name,CreateDT,PaymentFlag,OfferPrice FROM THPDMPDB.VW_OMSMobileJobPayment  where trackingid = '" + MsgText + "' order by DBID DESC ")

			if err == nil {

				for ress2.Next() {
					var event TblTMSItemEvents

					//JobID := ress2.Scan(&event.JobID)
					err := ress2.Scan(&event.DBID, &event.TrackingID, &event.Status_Name, &event.CreateDT, &event.PaymentFlag, &event.OfferPrice)

					if err != nil {
						panic(err)
					}
					paymentMassage := "ยังไม่ชำระเงินค่าขนส่ง"
					if event.PaymentFlag == "Y" {
						paymentMassage = "ชำระเงินค่าขนส่งแล้วจำนวน " + event.OfferPrice
					}

					dt := strings.Split(event.CreateDT, "T")[0]
					dt = strings.Split(dt, "-")[2] + "-" + strings.Split(dt, "-")[1] + "-" + strings.Split(dt, "-")[0]

					tt := strings.Split(event.CreateDT, "T")[1]
					tt = strings.Replace(tt, "Z", "", -1)

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
								"label": "เข้าสู่เว็ปไซต์มาเก็ตเพรส",
								"uri": "http://example.com/page/123"
							},
							"actions": [
								{
									"type": "uri",
									"label": "ชำระผ่าน Credit Card",
									"uri": "https://report.thpdlogistics.com?Credit=` + MsgText + `"  
									
								},
								{
								  "type": "uri",
								  "label": "ชำระผ่าน QRCode",
								  "uri": "https://report.thpdlogistics.com?QRCode=` + MsgText + `" 
								}
							]
						}
					  }	 ]}`)

					if event.PaymentFlag == "Y" {

						myJsonString = []byte(`{"replyToken":"` + ReplyToken + `","messages":[{"type":"text","text":"สถานะการค้นหา "},{"type":"text","text":" TrackID : ` + MsgText + `"},{"type":"text","text":" สถานะจัดส่งล่าสุด : ` + event.Status_Name + ` วันที่ ` + dt + ` เวลา ` + tt + ` "},{"type":"text","text":" การชำระเงิน : ` + paymentMassage + `"}	 ]}`)
					}

					count = 1
					break
				}

				defer ress2.Close()
				err = ress2.Close()

				if count == 0 {
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
			fmt.Println(req.Header)

			//req.Body.Read(json_data)
			//Send req using http Client
			client := &http.Client{}
			resp2, err := client.Do(req)

			if resp2 != nil {
				//fmt.Println(resp2)
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
							"thumbnailImageUrl": "https://www.thpd.co.th/images/c-home/aw_infographic-5-steps-service-chain_2.png",
							"imageAspectRatio": "rectangle",
							"imageSize": "cover",
							"imageBackgroundColor": "#FFFFFF",
							"title": "เลือกช่องทางที่ท่านสนใจ",
							"text": "Please select",
							"defaultAction": {
								"type": "uri",
								"label": "เข้าสู่เว็ปไซต์มาเก็ตเพรส",
								"uri": "http://example.com/page/123"
							},
							"actions": [
								{
									"type": "uri",
									"label": "เข้าสู่เว็ปไซต์มาเก็ตเพรส",
									"uri": "http://example.com/page/123"
								},
								{
								  "type": "uri",
								  "label": "เข้าสู่ระบบ OMS",
								  "uri": "http://example.com/page/123"
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
							  "thumbnailImageUrl": "https://img.freepik.com/free-photo/durian-fruit-with-slices-leaves-isolated-white-surface_252965-910.jpg?w=2000",
							  "imageBackgroundColor": "#FFFFFF",
							  "title": "ทุเรียนหมอนทองระยอง  ",
							  "text": "โลละ 80 บาท (มีแค่ 200 กก.เท่านั้น) รีบจองเลย",
							  "defaultAction": {
								  "type": "uri",
								  "label": "View detail",
								  "uri": "http://example.com/page/123"
							  },"actions": [
								{
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/111"
								}
							]
							 
							},
							{
								"thumbnailImageUrl": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQeJagqaIctMtB-So04pAmAPr_MrUaXdNMyry08N16ghCsL7fMrUX9MUI7Hsg_OJcjXZzg&usqp=CAU",
								"imageBackgroundColor": "#FFFFFF",
								"title": "ทุเรียนภูเขาไฟ ศรีษะเกษแท้ ",
								"text": "โลละ 80 บาท",
								"defaultAction": {
									"type": "uri",
									"label": "View detail",
									"uri": "http://example.com/page/123"
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
		fmt.Println(req.Header)

		//req.Body.Read(json_data)
		//Send req using http Client
		client := &http.Client{}
		resp2, err := client.Do(req)

		if resp2 != nil {
			//fmt.Println(resp2)
			//panic(err)
			myJsonString = []byte(`{"replyToken":"` + ReplyToken + `"}`)
		}
		if err == nil {

		}

	}

}

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
	fmt.Printf("timeUnix: \n", typefind)
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
			return
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
func getDNSString(dbName, dbUser, dbPassword, conn string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&timeout=60s&readTimeout=60s",
		dbUser,
		dbPassword,
		conn,
		dbName)
}

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
	res, err := db.Exec(insertGetDriver, event.JobID, event.TrackingID, event.CarID, event.DriverName, event.DriverPhoneNo, event.Company, event.Price, event.CustomerSelect,
		event.CustomConfirmDT, event.GetBeforeDT)
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
// 		//fmt.Println(article.Id)

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
// 		//fmt.Println(article.Id)
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
// 	// 	fmt.Println(article.Id)
// 	// }

// 	//fmt.Println(article.Id)
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

func OMSMobileCreateOrder(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var m MyStruct
	typefind := article.MobileID

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	if typefind != "" {
		//var unixTime int64 = int64(m.Date)

		event5 := Event5{
			JobID:           string(article.JobID),
			MobileID:        string(article.MobileID),
			ReceiveName:     string(article.ReceiveName),
			ReceiveAddress:  string(article.ReceiveAddress),
			ReceiveTumbon:   string(article.ReceiveTumbon),
			ReceiveDistrict: string(article.ReceiveDistrict),
			ReceiveProvince: string(article.ReceiveProvince),
			ReceivePhoneNo:  string(article.ReceivePhoneNo),
			ReceiveZipcode:  string(article.ReceiveZipcode),
			SenderZipcode:   string(article.SenderZipcode),
			SendPrice:       string(article.SendPrice),
			PickupStartDt:   string(article.PickupStartDt),
			DeliveryEndDt:   string(article.DeliveryEndDt),
			PaymentFlag:     string(article.PaymentFlag),
			PaymentDetail:   string(article.PaymentDetail),
			IMG1:            string(article.IMG1),
			IMG2:            string(article.IMG2),
			IMG3:            string(article.IMG3),
			IMG4:            string(article.IMG4),
			JobDesc:         string(article.JobDesc),
			JobType:         string(article.JobType),
			MerchantID:      string(article.MerchantID),
			WarehouseID:     string(article.WarehouseID),
			Status:          "รอดำเนินการ",
			//CreateDt:        string(article.CreateDt),
		}
		insertedId, err := insertOrderItems(db, event5)
		fmt.Println(insertedId)
		if err != nil {
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
			respok := make(map[string]string)
			respok["Success"] = "Insert Success"
			jsonResp, err := json.Marshal(respok)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}
		//fmt.Println(insertedId)
	}

	//return

}

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
			/// ยิง api เลือกคนขับ
			UpdateAPIDriverToLoadboard(article.TrackingID, article.DriverID, drivercompany)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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

			ress2, err2 := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET CustomerSelect = '1', CustomConfirmDT = NOW() WHERE  TrackingID = '" + article.TrackingID + "' and DriverName = '" + drivername + "' ")
			defer ress2.Close()
			if err2 != nil {
				panic(err)
			}
			ress3, err3 := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'Booking' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + article.TrackingID + "')")
			defer ress3.Close()
			if err3 != nil {
				panic(err)
			}
			/// ยิง api เลือกคนขับ
			UpdateAPIDriverToLoadboard(article.TrackingID, article.DriverID, drivercompany)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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

func OMSMobileGetCoupon(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
func OMSMobileGetRateAddOn(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
func OMSMobileGetWareHouse(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Merchant
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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

		ress2, err := db.Query("SELECT MerchantID, Address1, Address2, IFNULL(LocationGPS,'0,0')LocationGPS, WarehouseID, IFNULL(SubDistrict,'')SubDistrict, IFNULL(District,'0,0')District,ProvinceName,PostCode FROM THPDMPDB.tblMerchantWH WHERE MerchantID = '" + article.MerchantID + "'  ")

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
func UpdateAPIDriverToLoadboard(trackID string, carrier_id string, drivername string) string {

	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	payload := strings.NewReader(`{"con_no": "` + trackID + `" ,"carrier_id": "` + carrier_id + `"}`)

	method := "POST"
	url := "https://broker.dxplace.com/API/Gateways/selectCarrier"

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
		fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}
	fmt.Println(string(body))

	//sqlstr := "INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  "
	//fmt.Println(string(sqlstr))

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

		SendMessageToDriver(trackID, "0850270971", carrier_id, drivername)

	}
	if err == nil {
	}

	return string("ok")
}
func UpdateAPIPaymentToLoadboard(trackID string, paymentDetail string, drivername string) string {

	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	t := time.Now() //It will return time.Time object with current timestamp
	tUnix := t.Unix()
	fmt.Printf("timeUnix: %d\n", tUnix)

	var n int64 = tUnix

	s := strconv.FormatInt(n, 10)

	fmt.Printf("%d\n", tUnix)
	//strOTP := strconv.FormatInt(tUnix, 6)

	//payload := strings.NewReader(`{"con_no": "` + trackID + `" ,"carrier_id": "` + carrier_id + `"}`)
	payload := strings.NewReader(`{
		"status": 51,
		"actual_time": ` + s + `,
		"action_time": ` + s + `,
		"jobs": [{ "con_no": "` + trackID + `" }],
		"doc_link": ""
	  }`)

	method := "POST"
	url := "https://broker.dxplace.com/API/Gateways/updateStatus"

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
		fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}
	fmt.Println(string(body))

	ress3, err2 := db.Query("INSERT INTO  THPDMPDB.tblMobileAPILog ( TrackingID, APIName, LogResponse, CreateDT) Values  ( '" + trackID + "' , '" + strings.ReplaceAll(url, "'", `\'`) + "', '" + strings.ReplaceAll(string(body), "'", `\'`) + "', CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'))  ")
	defer ress3.Close()
	if err2 != nil {
		panic(err)
	}

	msgbody := paymentDetail + " , DXResponse: " + strings.ReplaceAll(string(body), "'", `\'`)
	// price unit = ลูกค้าชำระเงินแล้ว
	ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder a  INNER JOIN THPDMPDB.VW_OMSMobileOrderNoIMG2 b INNER JOIN THPDMPDB.tblPayment c ON a.JobID = b.JobID and b.TrackingID = c.TrackingID SET a.PaymentFlag =  'Y' , a.PaymentDetail = '" + msgbody + "'   WHERE b.TrackingID = '" + trackID + "'  ")
	defer ress.Close()
	if err != nil {
		panic(err)
	} else {
		/// SMS ให้คนขับ ///

		//SendMessageToDriver(trackID, "0850270971", carrier_id, drivername)

	}
	if err == nil {
	}

	return string("Success")
}
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
func ChkAuth(Authorization []string, Channel []string) int {

	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ress2, err := db.Query("SELECT Cid, ChannelName, ChannelSecretKey FROM THPDMPDB.tblAuthChannel  WHERE ChannelName = '" + Channel[0] + "'  ")

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
		ress, err := db.Query("UPDATE  THPDMPDB.tblAuthChannel Set UseCount = UseCount + 1, LastUseDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00')  WHERE ChannelName = '" + Channel[0] + "'  ")
		defer ress.Close()
		if err != nil {
			panic(err)
		} else {

		}

		return event.Cid
		// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

	}

	defer ress2.Close()
	err = ress2.Close()

	return returnval
}
func FntGetpaymentToken(w http.ResponseWriter, r *http.Request) {
	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//PreparePayload("JT04")

	reqHeader := r.Header["Authorization"]
	//fmt.Println("Header", reqHeader)

	reqHeaderChannel := r.Header["Channel"]
	//fmt.Println("Header", reqHeaderChannel)

	aa := ChkAuth(reqHeader, reqHeaderChannel)
	fmt.Println("aa", aa)

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

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Get2c2ppaymentToken
	json.Unmarshal(reqBody, &article)

	TrackingID := article.TrackingID
	Amount, err := strconv.Atoi(article.Amount)
	PhoneNumber := article.PhoneNumber
	TypePayment := article.PaymentType

	chars := []rune(TrackingID)
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

	chkdigit := 0

	if tall == 0 {
		chkdigit = 5
	} else if tall == 1 {
		chkdigit = 0
	} else {
		chkdigit = 11 - int(tall)
	}
	chkdigit2 := int(t9)
	if int(chkdigit) != chkdigit2 {
		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                        //QR
		respok["ErrMsg"] = "Invalid TrackID Format" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		respok := make(map[string]string)
		respok["err"] = err.Error()

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)

		os.Exit(1)
	}
	//QRCODE
	//userData := map[string]interface{}{"merchantID": "JT04", "invoiceNo": TrackingID, "description": PhoneNumber, "amount": Amount, "currencyCode": "THB"}
	//CREDITCARD
	var emp PaymentChannel
	emp.Amount = Amount
	emp.CurrencyCode = "THB"
	emp.Description = PhoneNumber
	emp.InvoiceNo = TrackingID
	emp.MerchantID = "JT04"
	emp.PaymentChannel = []string{"CC"}
	emp.Tokenize = true

	var data2 []byte
	// Convert struct to json
	data2, _ = json.MarshalIndent(emp, "", "    ")

	//userData := map[string]interface{}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(data2), &jsonMap)

	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	userData := map[string]interface{}{"merchantID": "JT04", "invoiceNo": TrackingID, "description": PhoneNumber, "amount": Amount, "currencyCode": "THB"}

	if TypePayment == "CreditCard" {
		userData = jsonMap
	}

	//userData := map[string]interface{}{"merchantID": "JT04", "invoiceNo": TrackingID, "description": PhoneNumber, "amount": Amount, "currencyCode": "THB", "tokenize": true, "paymentChannel": ["CC"]}
	//userData := b

	accessToken, err := Sign(userData, "CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2", 1) // data -> secretkey env name -> expiredAt

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		respok := make(map[string]string)
		respok["err"] = err.Error()

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)

		os.Exit(1)
	}

	// respok := make(map[string]string)
	// respok["Step1"] = accessToken
	// jsonResp, err := json.Marshal(respok)
	// w.Write(jsonResp)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  One of the following flags is required: sign, verify\n")
		flag.PrintDefaults()
	}

	// Parse command line options
	flag.Parse()

	// create payment token

	tokenString := CreatePaymentToken(accessToken, "CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2", "JT04", TypePayment)

	if tokenString == "Existing Invoice Number" {

		// update payment already
		UpdateAPIPaymentToLoadboard(TrackingID, tokenString, "")
		//

		respok := make(map[string]string)
		//respok["err"] = "Existing Invoice Number"

		respok["QRcodetxt"] = ""
		respok["ExpDT"] = ""                         //QR
		respok["ErrMsg"] = "Existing Invoice Number" //QR

		jsonResp, err := json.Marshal(respok)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}
		w.Write(jsonResp)

		//os.Exit(1)
		//return
	} else {

		///  "paymentToken": "kSAops9Zwhos8
		//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3ZWJQYXltZW50VXJsIjoiaHR0cHM6Ly9zYW5kYm94LXBndy11aS4yYzJwLmNvbS9wYXltZW50LzQuMS8jL3Rva2VuL2tTQW9wczlad2hvczhoU1RTZUxUVWJMTEhodndCRzBlaENkUE4wSDN0QmJjVHhpeXpkdnl4VVM4WWx3bHRRUHZqSllDZzdTMXpkTEVMOXhLOXZDdHkyY2JBTEYyRUVxYXdZcVgwdmREQUZBJTNkIiwicGF5bWVudFRva2VuIjoia1NBb3BzOVp3aG9zOGhTVFNlTFRVYkxMSGh2d0JHMGVoQ2RQTjBIM3RCYmNUeGl5emR2eXhVUzhZbHdsdFFQdmpKWUNnN1MxemRMRUw5eEs5dkN0eTJjYkFMRjJFRXFhd1lxWDB2ZERBRkE9IiwicmVzcENvZGUiOiIwMDAwIiwicmVzcERlc2MiOiJTdWNjZXNzIn0.QeXmFdOTj7krJB22yU-q8sss4MQKT5t3ne12jMUZpOc"
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("CD229682D3297390B9F66FF4020B758F4A5E625AF4992E5D75D311D6458B38E2"), nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)

			respok := make(map[string]string)
			respok["err"] = err.Error()

			jsonResp, err := json.Marshal(respok)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(jsonResp)

			os.Exit(1)
		}

		if TypePayment == "CreditCard" {

			var result User2
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				result.paymentToken = claims["paymentToken"].(string)
				result.webPaymentUrl = claims["webPaymentUrl"].(string)
				//result.ExpDT = claims["expiryDescription"].(string)
				//json.NewEncoder(w).Encode(result)
				//return
			}
			fmt.Println("my paymentToken here", result.paymentToken)
			fmt.Println("webPaymentUrl", result.webPaymentUrl)

			respok := make(map[string]string)
			respok["paymentToken"] = result.paymentToken
			respok["webPaymentUrl"] = result.webPaymentUrl

			jsonResp, err := json.Marshal(respok)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(jsonResp)

		} else if TypePayment == "QRCode" {

			var result User
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				result.paymentToken = claims["paymentToken"].(string)
				//result.webPaymentUrl = claims["webPaymentUrl"].(string)
				//result.ExpDT = claims["expiryDescription"].(string)
				//json.NewEncoder(w).Encode(result)
				//return
			}
			fmt.Println("my paymentToken here", result.paymentToken)
			//fmt.Println("webPaymentUrl", result.webPaymentUrl)

			BarcodeQR := CreatePaymentQRCode(result.paymentToken)
			//fmt.Println("my Barcode here", BarcodeQR)

			if BarcodeQR != "" {

				respok := make(map[string]string)
				respok["QRcodetxt"] = BarcodeQR
				respok["ExpDT"] = result.ExpDT //QR
				respok["ErrMsg"] = ""          //QR

				jsonResp, err := json.Marshal(respok)

				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
					return
				}
				w.Write(jsonResp)
				//counter = 0

			}
		} else {

			t := time.Now() //It will return time.Time object with current timestamp
			fmt.Printf("time.Time %s\n", t)

			tUnix := t.Unix()
			fmt.Printf("timeUnix: %d\n", tUnix)

			var result RespCode
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

				result.respDesc = claims["respDesc"].(string)
				result.refCode = claims["referenceNo"].(string)
				result.amount = claims["amount"].(float64)
				result.invoiceNo = claims["invoiceNo"].(string)
				result.agentCode = claims["agentCode"].(string)
				result.cardType = claims["cardType"].(string)
				result.cardNo = claims["cardNo"].(string)
				result.transactionDatetime = claims["transactionDateTime"].(string)
				result.refID = tUnix
				result.phoneNo = PhoneNumber

			}

			//fmt.Println("my respDesc here", result.respDesc)
			//fmt.Println("webPaymentUrl", result.webPaymentUrl)
			//BarcodeQR := CreatePaymentQRCode(result.respDesc)
			//fmt.Println("my Barcode here", BarcodeQR)
			//if BarcodeQR != "" {

			if result.respDesc == "Success" {

				event5 := EventPayment{
					respDesc:            result.respDesc,
					amount:              fmt.Sprintf("%f", result.amount),
					refCode:             result.refCode,
					invoiceNo:           result.invoiceNo,
					agentCode:           result.agentCode,
					cardType:            result.cardType,
					transactionDatetime: result.transactionDatetime,
					phoneNo:             result.phoneNo,
					cardNo:              result.cardNo,
					refID:               strconv.FormatInt(result.refID, 10),
					//CreateDt:        string(article.CreateDt),
				}

				insertedId, err := insertPayment(db, event5)
				fmt.Println(insertedId)

				// update payment already
				if err == nil {
					UpdateAPIPaymentToLoadboard(result.invoiceNo, result.respDesc, "")
				}

			}
			//

			respok := make(map[string]string)
			respok["respDesc"] = result.respDesc
			respok["amount"] = fmt.Sprintf("%f", result.amount) //QR
			respok["refCode"] = result.refCode                  //QR
			respok["invoiceNo"] = result.invoiceNo
			respok["agentCode"] = result.agentCode
			respok["cardType"] = result.cardType
			respok["transactionDatetime"] = result.transactionDatetime

			jsonResp, err := json.Marshal(respok)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				return
			}
			w.Write(jsonResp)
		}
	}

	////return string("ok")
}
func CreatePaymentToken(accessToken string, SECRETKEY string, MERCHANTID string, TypePayment string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)

	//	 var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	//var emp THSMS
	//emp.payload = accessToken

	//var data []byte
	// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")
	payload := strings.NewReader(`{"payload": "` + accessToken + `" }`)

	url := "https://sandbox-pgw.2c2p.com/payment/4.1/paymentToken"

	if TypePayment == "Check" {
		url = "https://sandbox-pgw.2c2p.com/payment/4.1/paymentInquiry"
	}

	//url := "https://sandbox-pgw.2c2p.com/payment/4.1/paymentToken"

	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("SECRETKEY", SECRETKEY)
	req.Header.Add("MERCHANTID", MERCHANTID)
	req.Header.Add("typ", "JWT")
	req.Header.Add("alg", "HS256")
	req.Header.Add("Content-Type", "application/*+json")
	//fmt.Println(req.Header)
	//req.Body.Read(json_data)
	//Send req using http Client
	//client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := io.ReadAll(resp2.Body)
	bodyString := string(body)
	//log.Info(bodyString)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return string(err.Error())
	// }
	//fmt.Println(string(body))

	if body != nil {

		fmt.Print(bodyString)

		var article Payload2

		json.Unmarshal(body, &article)

		if article.RespDesc == "Existing Invoice Number" {
			article.Payload = "Existing Invoice Number"
		}

		//b, _ := json.Marshal(body)

		//fmt.Printf(article.Payload)
		//fmt.Println(b)
		return article.Payload
		//panic(err)
	}
	if err != nil {

		return "null"
	}

	return string("OK")
}
func CreatePaymentQRCode(paymentToken string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)

	//	 var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	//var emp THSMS
	//emp.payload = paymentToken

	//var data []byte
	// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")
	payload := strings.NewReader(`{"responseReturnUrl": "https://sandbox-pgw-ui.2c2p.com/payment/4.1/#/info/",
    "payment": {
        "code": {
            "channelCode": "PPQR"
        },
        "data": {
            "name": "Terrance Tay",
            "email": "terrance.tay@2c2p.com"
        }
    },
    "clientIP": "",
    "paymentToken": "` + paymentToken + `",
    "locale": "th",
    "clientID": "30c7cf51-75c4-4265-a70a-effddfbbb0ff" }`)

	url := "https://sandbox-pgw.2c2p.com/payment/4.1/Payment"
	//req, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))

	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	// req.Header.Add("SECRETKEY", SECRETKEY)
	// req.Header.Add("MERCHANTID", MERCHANTID)
	// req.Header.Add("typ", "JWT")
	// req.Header.Add("alg", "HS256")
	req.Header.Add("Content-Type", "application/*+json")
	//fmt.Println(req.Header)
	//req.Body.Read(json_data)
	//Send req using http Client
	//client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := io.ReadAll(resp2.Body)
	//bodyString := string(body)
	//log.Info(bodyString)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return string(err.Error())
	// }
	//fmt.Println(string(body))

	if body != nil {

		//fmt.Print(bodyString)

		var article PaymentQRStruct

		json.Unmarshal(body, &article)
		//b, _ := json.Marshal(body)

		//fmt.Printf(article.Data)
		//fmt.Println(b)
		if article.Data != "" {
			return article.Data
		} else {
			return article.RespDesc
		}

		//panic(err)
	}
	if err != nil {

		return "null"
	}

	return string("")
}
func Sign(Data map[string]interface{}, SecrePublicKeyEnvName string, ExpiredAt time.Duration) (string, error) {

	expiredAt := time.Now().Add(time.Duration(time.Second) * ExpiredAt).Unix()

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

	return accessToken, nil
}
func GetAPIDriverFromLoadboard(trackID string) string {

	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var app_id = "de_oms"
	var app_key = "dx1234"

	var emp Con_no
	emp.con_no = trackID

	//var data []byte
	//// Convert struct to json
	//data, _ = json.MarshalIndent(emp, "", "    ")

	payload := strings.NewReader(`{"con_no": "` + trackID + `"}`)

	method := "POST"
	url := "https://broker.dxplace.com/API/Gateways/viewOrder"

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
		fmt.Println(err)
		return err.Error()
	}
	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}
	//fmt.Println(string(body))

	var Loadboard DriverLoadboard
	jsonData := []byte(string(body))

	js := DriverLoadboard{}

	json.Unmarshal(jsonData, &js)
	if err != nil {
		panic(err)
	}
	fmt.Println(js)

	json.Unmarshal([]byte(body), &Loadboard)

	//fmt.Println(string(logistic.Data.Booking[0].BookingID))
	fmt.Println(Loadboard.Data.Shipment.BookingCount)
	//for i := 0; i < Loadboard.Data.Shipment.BookingCount; i++ {
	for i := 0; i < len(Loadboard.Data.Booking); i++ {
		fmt.Println(string(Loadboard.Data.Booking[i].BookingID))

		eventBooking := EventBooking{
			BookingID:         string(Loadboard.Data.Booking[i].BookingID),
			BookingStatus:     string(Loadboard.Data.Booking[i].BookingStatus),
			Price:             string(Loadboard.Data.Booking[i].Price),
			PriceUnit:         string(Loadboard.Data.Booking[i].PriceUnit),
			TruckPlate:        string(Loadboard.Data.Booking[i].TruckPlate),
			DriverName:        string(Loadboard.Data.Booking[i].DriverName),
			Oid:               string(Loadboard.Data.Booking[i].Oid),
			Company:           string(Loadboard.Data.Booking[i].Company),
			Address:           string(Loadboard.Data.Booking[i].Address),
			Contact:           string(Loadboard.Data.Booking[i].Contact),
			Phone:             string(Loadboard.Data.Booking[i].Phone),
			Email:             string(Loadboard.Data.Booking[i].Email),
			Rating:            string(Loadboard.Data.Booking[i].Rating),
			BookingScore:      string(Loadboard.Data.Booking[i].BookingScore),
			Time:              string(Loadboard.Data.Booking[i].Time),
			PositionLatitude:  string(Loadboard.Data.Booking[i].PositionLatitude),
			PositionLongitude: string(Loadboard.Data.Booking[i].PositionLongitude),
			PositionTime:      string(Loadboard.Data.Booking[i].PositionTime),
		}

		ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET BookingID = '" + eventBooking.BookingID + "',DriverPhoneNo = '" + eventBooking.Phone + "', Company = '" + eventBooking.Company + "' , DriverStar = '" + eventBooking.Rating + "' , Price = '" + eventBooking.Price + "', Oid = '" + eventBooking.Oid + "', Contact = '" + eventBooking.Contact + "', BookingScore = '" + eventBooking.BookingScore + "', Time = '" + eventBooking.Time + "', DriverEmail = '" + eventBooking.Email + "', PriceUnit = '" + eventBooking.PriceUnit + "'  WHERE TrackingID = '" + trackID + "' and  CarID = '" + eventBooking.TruckPlate + "' and DriverName = '" + eventBooking.DriverName + "' ")
		defer ress.Close()
		if err != nil {
			panic(err)
		} else {

		}
		//return "OK"

	}

	if err == nil {
	}

	return string("ok")
}

func SendMessageToDriver(trackingID string, CustomerPhone string, CarID string, DriverName string) string {

	//dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	//db, err := sql.Open("mysql", dns)

	var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	var emp THSMS
	emp.Sender = "Privileged"
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
	fmt.Println(req.Header)

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
func SendMessage(trackingID string, CustomerPhone string, CarID string, DriverName string) string {

	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var bearerToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC93d3cudGhzbXMuY29tXC9hcGkta2V5IiwiaWF0IjoxNjQ4MDkwNjUyLCJuYmYiOjE2NDgxMDEzMjUsImp0aSI6Ik00RkdVQjN5OFd6NnZuYzciLCJzdWIiOjEwNDE2MCwicHJ2IjoiMjNiZDVjODk0OWY2MDBhZGIzOWU3MDFjNDAwODcyZGI3YTU5NzZmNyJ9.PwmdMYwIdXIWRftvcrnqTDiulwTfcFVsLDpBj4REyI4"

	var emp THSMS
	emp.Sender = "Privileged"
	emp.Msisdn = []string{CustomerPhone}
	emp.Message = "OMS DE : หมายเลข Track: " + trackingID + "\n  ผู้ขนส่ง: " + DriverName + "\n  ทะเบียนรถ: " + CarID + "\n รับงานของท่านแล้ว!"

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
	fmt.Println(req.Header)

	//req.Body.Read(json_data)
	//Send req using http Client
	client := &http.Client{}
	resp2, err := client.Do(req)

	if resp2 != nil {
		//panic(err)
	}
	if err == nil {
		/// update tracking send message already
		ress, err := db.Query("UPDATE THPDMPDB.tblMobileOMSJobDriverBooking SET SendSMSFirstJob = 1 WHERE TrackingID = '" + trackingID + "'")
		defer ress.Close()
		if err != nil {

			panic(err)
		} else {
			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSOrder SET Status = 'ผู้ขนส่งรับงานแล้ว' WHERE JobID in ( SELECT Customer_Po FROM THPDMPDB.tblOrderMaster WHERE TrackingID = '" + trackingID + "')")
			defer ress2.Close()
			if err != nil {
				panic(err)
			} else {

			}

		}
		//panic(err)
	}
	return string("OK")
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func OMSMobileGetAdvertise(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var m MyStruct
	//typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//strOTP := strconv.FormatInt(tUnix, 6)
	//strOTP := EncodeToString(6) // substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}
	//datett := article.Date

	// event5 := EventDriverTrack{
	// 	JobID:           "",
	// 	TrackingID:      string(article.TrackingID),
	// 	CarID:           string(article.Plate),
	// 	DriverName:      string(article.DriverName),
	// 	DriverPhoneNo:   "",
	// 	Company:         "",
	// 	Price:           "",
	// 	CustomerSelect:  "",
	// 	CustomConfirmDT: "",
	// 	GetBeforeDT:     "",

	// 	//CreateDt:        string(article.CreateDt),
	// }
	// insertedId, err := insertGetDriverTrack(db, event5)
	// fmt.Println(insertedId)

	//boxes := []DriverTracking{}
	//boxes = append(boxes, BoxData{Width: 5, Height: 30})

	// resp := make(map[string]string)
	// resp["TrackingID"] = article.TrackingID

	//boxes = append(boxes, DriverTracking{TrackingID: article.TrackingID, DriverID: article.DriverID, DriverName: article.DriverName, Plate: article.Plate})

	//b, err := json.Marshal(boxes)
	UrlText := ""

	if err == nil {

		ress2, err := db.Query("SELECT UrlText FROM THPDMPDB.tblMobileAdvertise  ")

		if err == nil {

			// boxes := []Merchant{}
			// resp := make(map[string]string)
			// resp["MerchantID"] = article.MerchantID

			for ress2.Next() {
				var event UrlAdvise
				//JobID := ress2.Scan(&event.JobID)
				err = ress2.Scan(&event.UrlText)

				if err != nil {
					panic(err)
				}
				UrlText = event.UrlText

				// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

			}

		}
		defer ress2.Close()
		err = ress2.Close()

		respok := make(map[string]string)
		respok["Success"] = "true"
		respok["UrlText"] = UrlText
		jsonResp, err := json.Marshal(respok)

		/// update driver get job
		GetAPIDriverFromLoadboard(article.TrackingID)
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

func OMSMobileGetDriver(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article DriverTracking
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var m MyStruct
	typefind := article.TrackingID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//strOTP := strconv.FormatInt(tUnix, 6)
	//strOTP := EncodeToString(6) // substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}
	//datett := article.Date

	if typefind != "" {

		if err == nil {

			event5 := EventDriverTrack{
				JobID:           "",
				TrackingID:      string(article.TrackingID),
				CarID:           string(article.Plate),
				DriverName:      string(article.DriverName),
				DriverPhoneNo:   "",
				Company:         "",
				Price:           "",
				CustomerSelect:  "",
				CustomConfirmDT: "",
				GetBeforeDT:     "",

				//CreateDt:        string(article.CreateDt),
			}
			insertedId, err := insertGetDriverTrack(db, event5)
			fmt.Println(insertedId)

			//boxes := []DriverTracking{}
			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["TrackingID"] = article.TrackingID

			//boxes = append(boxes, DriverTracking{TrackingID: article.TrackingID, DriverID: article.DriverID, DriverName: article.DriverName, Plate: article.Plate})

			//b, err := json.Marshal(boxes)
			mobileid := ""
			SendSMSFirstJob := 0
			if err == nil {

				ress2, err := db.Query("SELECT * FROM THPDMPDB.VW_NewDriverSendMassage  WHERE TrackingID = '" + article.TrackingID + "'  ")

				if err == nil {

					// boxes := []Merchant{}
					// resp := make(map[string]string)
					// resp["MerchantID"] = article.MerchantID

					for ress2.Next() {
						var event MobileID
						//JobID := ress2.Scan(&event.JobID)
						err = ress2.Scan(&event.MobileID, &event.SendSMSFirstJob, &event.TrackingID)

						if err != nil {
							panic(err)
						}
						mobileid = event.MobileID
						SendSMSFirstJob = event.SendSMSFirstJob
						// boxes = append(boxes, Merchant{MerchantID: event.MerchantID, Address1: event.Address1, Address2: event.Address2, LocationGPS: event.LocationGPS, WarehouseID: event.WarehouseID, SubDistrict: event.SubDistrict, District: event.District, ProvinceName: event.ProvinceName, PostCode: event.PostCode})

					}

				}
				defer ress2.Close()
				err = ress2.Close()

				if SendSMSFirstJob == 0 && mobileid != "" {
					SendMessage(article.TrackingID, mobileid, article.Plate, article.DriverName)
				}

				respok := make(map[string]string)
				respok["Success"] = "true"
				respok["TrackingID"] = article.TrackingID
				respok["Plate"] = article.Plate
				respok["DriverName"] = article.DriverName
				jsonResp, err := json.Marshal(respok)

				/// update driver get job
				GetAPIDriverFromLoadboard(article.TrackingID)
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

			// if err != nil {
			// 	log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			// 	return
			// }
			//w.Write(b)
			//counter = 0

		}

	}

}
func OMSMobileGetOTP(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var m MyStruct
	typefind := article.MobileID

	t := time.Now() //It will return time.Time object with current timestamp
	fmt.Printf("time.Time %s\n", t)

	// tUnix := t.Unix()
	// fmt.Printf("timeUnix: %d\n", tUnix)
	// strOTP := strconv.FormatInt(tUnix, 6)
	strOTP := EncodeToString(6) // substr(strOTP, 5, 6)
	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		//fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		//ress2, err := db.Query("SELECT MerchantID,MerchantName,Email,ContactName,PhoneNumber FROM THPDMPDB.tblMerchant WHERE PhoneNumber = '" + article.MobileID + "'  ")

		if err == nil {

			// boxes := []BoxData{}
			// boxes = append(boxes, BoxData{Width: 10, Height: 20})
			// boxes = append(boxes, BoxData{Width: 5, Height: 30})

			boxes := []MobileOTP{}

			//boxes = append(boxes, BoxData{Width: 5, Height: 30})

			resp := make(map[string]string)
			resp["MobileID"] = article.MobileID

			boxes = append(boxes, MobileOTP{MobileID: article.MobileID, OTPPassword: strOTP})

			b, _ := json.Marshal(boxes)

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
func OMSMobileUpdateLoginMember(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Mobile
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//var m MyStruct
	typefind := article.MobileID
	typeEvent := article.TypeEvent
	typelang := article.Lang
	bg := article.BG
	LoginType := article.LoginType

	//t := time.Now() //It will return time.Time object with current timestamp
	//fmt.Printf("time.Time %s\n", t)

	//tUnix := t.Unix()
	//fmt.Printf("timeUnix: %d\n", tUnix)
	//str := strconv.FormatInt(tUnix, 10)

	//firstEvent2 := Event5{}

	//datett := article.Date
	if typefind != "" {
		//fmt.Println(article.Id)

		// firstEvent := Event5{}
		// err = selectOMSOrder(db, article.MobileID, article.Token, &firstEvent)
		// if err != nil {
		// 	panic(err)
		// }

		//counter := 0

		if typeEvent == "Login" {
			ress3, err2 := db.Query("INSERT INTO THPDMPDB.tblMobileOMSMemberConfig (MemberPhone, CreateDT, LoginCount, LastLoginDT, Lang, BG, MerchantID,LoginType) Values ('" + article.MobileID + "',CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),1,CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00'),'th','White','c','" + LoginType + "') ON DUPLICATE KEY UPDATE LastLoginDT = CONVERT_TZ(CURRENT_TIME(),'+00:00','+7:00') , LoginCount = LoginCount + 1 ,LoginType = '" + LoginType + "'")
			defer ress3.Close()
			if err2 != nil {
				panic(err)
			}
			defer ress3.Close()
			err = ress3.Close()

		} else if typeEvent == "GetLang" {

		} else if typeEvent == "SetLang" {
			ress2, err := db.Query("UPDATE THPDMPDB.tblMobileOMSMemberConfig SET Lang = '" + typelang + "' , BG = '" + bg + "' WHERE MemberPhone = '" + article.MobileID + "'  ")
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
	ress2, err := db.Query("SELECT MemberPhone, Lang, BG FROM THPDMPDB.tblMobileOMSMemberConfig WHERE MemberPhone = '" + article.MobileID + "'  ")

	if err == nil {

		// boxes := []BoxData{}
		// boxes = append(boxes, BoxData{Width: 10, Height: 20})
		// boxes = append(boxes, BoxData{Width: 5, Height: 30})

		boxes := []Mobile{}

		for ress2.Next() {
			var event Mobile
			//JobID := ress2.Scan(&event.JobID)
			err := ress2.Scan(&event.MobileID, &event.Lang, &event.BG)

			if err != nil {
				panic(err)
			}

			boxes = append(boxes, Mobile{MobileID: event.MobileID, Lang: event.Lang, BG: event.BG})
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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
func OMSMobileGetPostcode(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Zipcode
	json.Unmarshal(reqBody, &article)

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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

func OMSMobileDiscount(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
func OMSMobileSelectOrderNoIMG(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var article OMSOrderStruct
	json.Unmarshal(reqBody, &article)
	dns := getDNSString("THPDMPDB", "admin", "dedb<>10!", "thpd-dedb.cluster-ro-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")

	//dns := getDNSString("THPDMPDB", "admin", "vX0r7qBIEGk9eTBWBu7S", "thpddb.cluster-crcn7eiyated.ap-southeast-1.rds.amazonaws.com")
	db, err := sql.Open("mysql", dns)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
		//fmt.Println(article.Id)

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
	//myRouter.HandleFunc("/TOATGetByInvoiceID", returnSingleArticle2).Methods("POST")
	myRouter.HandleFunc("/OMSMobileCreateOrder", OMSMobileCreateOrder).Methods("POST")
	myRouter.HandleFunc("/OMSMobileSelectOrder", OMSMobileSelectOrder).Methods("POST")
	myRouter.HandleFunc("/OMSMobileSelectOrderNoIMG", OMSMobileSelectOrderNoIMG).Methods("POST")
	myRouter.HandleFunc("/OMSMobileSelectOrderWithIMG", OMSMobileSelectOrderWithIMG).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetPostcode", OMSMobileGetPostcode).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetMember", OMSMobileGetMember).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetJobDriverBooking", OMSMobileGetJobDriverBooking).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetJobDriverBookingMatchAlready", OMSMobileGetJobDriverBookingMatchAlready).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetOTP", OMSMobileGetOTP).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetDriver", OMSMobileGetDriver).Methods("POST")
	myRouter.HandleFunc("/OMSMobileUpdateLoginMember", OMSMobileUpdateLoginMember).Methods("POST")

	myRouter.HandleFunc("/OMSMobileGetAdvertise", OMSMobileGetAdvertise).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetCoupon", OMSMobileGetCoupon).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetRateAddOn", OMSMobileGetRateAddOn).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetWareHouse", OMSMobileGetWareHouse).Methods("POST")
	myRouter.HandleFunc("/OMSMobileGetTrcukSize", OMSMobileGetTrcukSize).Methods("POST")
	myRouter.HandleFunc("/OMSMobileUpdateDriverJobMaster", OMSMobileUpdateDriverJobMaster).Methods("POST")

	myRouter.HandleFunc("/OMSMobileConnect", OMSMobileConnect).Methods("POST")
	myRouter.HandleFunc("/GenToken", tokenGenerator).Methods("POST")
	myRouter.HandleFunc("/GenlblPrint", GenlblPrint).Methods("POST")
	myRouter.HandleFunc("/Callback", Callback).Methods("POST")
	myRouter.HandleFunc("/FntGetpaymentToken", FntGetpaymentToken).Methods("POST")
	myRouter.HandleFunc("/OMSMobileDiscount", OMSMobileDiscount).Methods("POST")
	//myRouter.HandleFunc("/TOATGetByCustomerID", returnSingleArticle2).Methods("POST")
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)

	myRouter.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":80", myRouter)) // production
	//log.Fatal(http.ListenAndServe(":8081", myRouter)) // test
	//log.Fatal(http.ListenAndServe(":5678", myRouter)) // test
}

func main() {

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
	DriverTrackings = []DriverTracking{
		DriverTracking{TrackingID: "1", DriverID: "xxx", DriverName: "xx", Plate: "xx", FinalPrice: "xx"},
	}
	DriverJobMasterBookings = []DriverJobMasterBooking{
		DriverJobMasterBooking{TrackingID: "1", JobDriverID: "xxx", JobDriverName: "xx", JobTruckID: "xx", JobTruckPlate: "xx"},
	}
	Get2c2ppaymentTokens = []Get2c2ppaymentToken{
		Get2c2ppaymentToken{TrackingID: "1", Amount: "100", PhoneNumber: "0809999999", PaymentType: "QR"},
	}
	GetPostPrintOnlines = []GetPostPrintOnline{
		GetPostPrintOnline{PostCode: "10010", CountItems: 1},
	}

	log.Println("server started")
	//http.HandleFunc("/webhook", handleWebhook)

	//log.Fatal(http.ListenAndServe(":8081", nil))

	handleRequests()

}
