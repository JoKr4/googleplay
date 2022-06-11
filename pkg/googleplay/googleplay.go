package googleplay

import (
	"fmt"
	"time"

	gp "github.com/JoKr4/googleplay"
)

func GetDetails(head *gp.Header, app string, parse bool) (*gp.Details, error) {
	details, err := head.Details(app)
	if err != nil {
		return nil, err
	}
	if parse {
		date, err := details.Time()
		if err != nil {
			return nil, err
		}
		details.UploadDate = date.String()
	}
	return details, nil
}

func DoDevice(dir, platform string, screenDensity int) error {
	device, err := gp.Phone.Checkin(platform, screenDensity)
	if err != nil {
		return err
	}
	fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
	time.Sleep(gp.Sleep)
	return device.Create(dir, platform+".json")
}

func GetHeader(dir, platform string, single bool) (*gp.Header, error) {
	token, err := gp.OpenToken(dir, "token.json")
	if err != nil {
		return nil, err
	}
	device, err := gp.OpenDevice(dir, platform+".json")
	if err != nil {
		return nil, err
	}
	return token.Header(device.AndroidID, single)
}
