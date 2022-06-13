package googleplay

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/89z/format"
)

func GetDetails(head *Header, app string, parse bool) (*Details, error) {
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

func DoDelivery(head *Header, app string, ver uint64, dir string) error {
	download := func(addr, name string) error {
		fmt.Println("GET", addr)
		res, err := http.Get(addr)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
		pro := format.ProgressBytes(file, res.ContentLength)
		if _, err := io.Copy(pro, res.Body); err != nil {
			return err
		}
		return nil
	}
	del, err := head.Delivery(app, ver)
	if err != nil {
		return err
	}
	folder := del.PackageName + "@" + strconv.FormatInt(int64(del.VersionCode), 10)
	folderpath := filepath.Join(dir, folder)
	_ = os.Mkdir(folderpath, 0755)
	for _, split := range del.SplitDeliveryData {
		fp := filepath.Join(folderpath, split.ID+".apk")
		err := download(split.DownloadURL, fp)
		if err != nil {
			return err
		}
	}
	for _, file := range del.AdditionalFile {
		fp := filepath.Join(folderpath, del.Additional(file.FileType))
		err := download(file.DownloadURL, fp)
		if err != nil {
			return err
		}
	}
	fp := filepath.Join(folderpath, "full.apk")
	if len(del.SplitDeliveryData) > 0 {
		fp = filepath.Join(folderpath, "base.apk")
	}
	return download(del.DownloadURL, fp)
}

func DoToken(dir, email, password string) error {
	token, err := NewToken(email, password)
	if err != nil {
		return err
	}
	return token.Create(dir, "token.json")
}

func DoDevice(dir, platform string, screenDensity int) error {
	device, err := Phone.Checkin(platform, screenDensity)
	if err != nil {
		return err
	}
	fmt.Printf("Sleeping %v for server to process\n", Sleep)
	time.Sleep(Sleep)
	return device.Create(dir, platform+".json")
}

func GetHeader(dir, platform string, single bool) (*Header, error) {
	token, err := OpenToken(dir, "token.json")
	if err != nil {
		return nil, err
	}
	device, err := OpenDevice(dir, platform+".json")
	if err != nil {
		return nil, err
	}
	return token.Header(device.AndroidID, single)
}
