package ciserver

import (
	"github.com/bndr/gojenkins"
)

var _STATUS = map[string]Status{
	"aborted":        ABORTED,
	"aborted_anime":  BUILDING_FROM_ABORTED,
	"blue":           SUCCESS,
	"blue_anime":     BUILDING_FROM_SUCCESS,
	"disabled":       DISABLED,
	"disabled_anime": BUILDING_FROM_DISABLED,
	"grey":           UNKNOWN,
	"grey_anime":     BUILDING_FROM_UNKNOWN,
	"notbuilt":       NOT_BUILT,
	"notbuilt_anime": BUILDING_FROM_NOT_BUILT,
	"red":            FAILURE,
	"red_anime":      BUILDING_FROM_FAILURE,
	"yellow":         UNSTABLE,
	"yellow_anime":   BUILDING_FROM_UNSTABLE,
}

//jenkinsClient *gojenkins.Jenkins

//InitClient initializes the Jenkins client - connects to jenkins instance
func InitClient(jenkinsURL string, username string, password string) (*gojenkins.Jenkins, error) {
	jenkins, err := gojenkins.CreateJenkins(jenkinsURL, username, password).Init()
	return jenkins, err
}
