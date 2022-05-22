package googleplay

import (
	"os"
	"time"

	gp "github.com/89z/googleplay"
)

func Header(platform string, agent int64) (*gp.Header, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	token, err := gp.OpenToken(cache, "googleplay/token.json")
	if err != nil {
		return nil, err
	}
	device, err := gp.OpenDevice(cache, "googleplay", platform+".json")
	if err != nil {
		return nil, err
	}
	return token.Header(device.AndroidID, gp.Agents[agent])
}

func Details(head *gp.Header, app string) (*gp.Details, error) {
	details, err := head.Details(app)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(gp.DateInput, string(details.UploadDate))
	if err == nil {
		details.UploadDate = gp.String(date.Format(gp.DateOutput))
	}
	return details, nil
}
