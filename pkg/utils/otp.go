package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

func SendOtp(to string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	var TWILIO_ACCOUNT_SID string = cfg.ACCOUNTSID
	var TWILIO_AUTH_TOKEN string = cfg.AUTHTOKEN
	var VERIFY_SERVICE_SID string = cfg.SERVICESSID
	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})
	to = "+91" + to
	fmt.Println(to)
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Printf("error ocurred while generating otp")
		fmt.Println(err)
		return err
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		return nil
	}
}

func CheckOtp(to string, code string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	var TWILIO_ACCOUNT_SID string = cfg.ACCOUNTSID
	var TWILIO_AUTH_TOKEN string = cfg.AUTHTOKEN
	var VERIFY_SERVICE_SID string = cfg.SERVICESSID
	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})
	to = "+91" + to
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	fmt.Println(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("hiii")
	} else if *resp.Status == "approved" {
		return nil
	} else {
		return errors.New("invalid otp entered")
	}
	return errors.New("invalid otp entered")
}
