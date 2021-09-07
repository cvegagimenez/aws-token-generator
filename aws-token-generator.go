package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pborman/getopt/v2"
)

type userInfo struct {
	accountId string
	username string
	mfaArn string
}

func getIdentity(profile string) *sts.GetCallerIdentityOutput {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
	  Profile: profile,
  }))
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	identity, err := svc.GetCallerIdentity(input)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        fmt.Println(err.Error())
    }
	}

	return identity
}

func getUserData(profile string) userInfo {

	userRawInfo := getIdentity(profile)

	account := userRawInfo.Account
	username := (*userRawInfo.Arn)[strings.LastIndex(*userRawInfo.Arn, "/")+1:]
	mfaArn := strings.Replace(*userRawInfo.Arn, ":user/", ":mfa/", -1)
	
	return userInfo {accountId: *account, 
									username: username, 
									mfaArn: mfaArn,
								}
}

func getMfaCredentials(profile string, token string, mfaArn string, expirTime int64) *sts.GetSessionTokenOutput {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
	  Profile: profile,
  }))
	svc := sts.New(sess)
	input := &sts.GetSessionTokenInput{
    DurationSeconds: aws.Int64(expirTime),
    SerialNumber:    aws.String(mfaArn),
    TokenCode:       aws.String(token),
	}

	credentials, err := svc.GetSessionToken(input)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case sts.ErrCodeRegionDisabledException:
            fmt.Println(sts.ErrCodeRegionDisabledException, aerr.Error())
        default:
            fmt.Println(aerr.Error())
        }
    } else {
        fmt.Println(err.Error())
    }

	}

	return credentials
}

func setEnvironmentVariables(credentials *sts.GetSessionTokenOutput, region string) {

	fmt.Println("Please, COPY/PASTE this commands in your terminal to configure the credentials: ")
	fmt.Println("")
	fmt.Println("export AWS_ACCESS_KEY_ID=" + *credentials.Credentials.AccessKeyId)
	fmt.Println("export AWS_SECRET_ACCESS_KEY=" + *credentials.Credentials.SecretAccessKey)
	fmt.Println("export AWS_SESSION_TOKEN=" + *credentials.Credentials.SessionToken)
	fmt.Println("export AWS_DEFAULT_REGION=" + region)

}

func main() {
	var profile = getopt.StringLong("profile", 'p', "", "AWS profile to use.")
	var mfaToken = getopt.StringLong("token", 't', "", "MFA token provided by device.")
	var region = getopt.StringLong("region", 'r', "eu-west-1", "Region to configure in AWS.")
	var expirTime = getopt.Int64Long("expiration", 'e', 3600, "Expiration time for the token.")

	getopt.Parse()

	if *mfaToken == "" || len(*mfaToken) != 6 {
		fmt.Println("The token have 6 digits at least.")
		os.Exit(1)
	}

	if *profile == "" {
		fmt.Println("The profile must not be empty.")
		os.Exit(1)
	}

	if *expirTime < 900 || *expirTime > 129600 {
		fmt.Println("The expiration time must be between 900 (15m) and 129600 (12h).")
		os.Exit(1)
	}

	credentials := getMfaCredentials(*profile, *mfaToken, getUserData(*profile).mfaArn, *expirTime)

	setEnvironmentVariables(credentials, *region)

}