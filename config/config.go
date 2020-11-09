package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func GetConfig() Configuration {
	conf := Configuration{}
	gonfig.GetConf("config/config.json", &conf)

	// fmt.Println(conf)
	return conf
}

type ObjectStorageConfiguration struct {
	OSS_Endpoint        string
	OSS_AccessKeyId     string
	OSS_AccessKeySecret string
	OSS_Bucket          string
}

func ObjectStorageConfig() ObjectStorageConfiguration {
	// conf := ObjectStorageConfiguration{}

	val := ObjectStorageConfiguration{
		"http://oss-ap-southeast-5.aliyuncs.com",
		"LTAI4FiLiiCr9HoP23gxj477",
		"sxW9RrMkJ4dMsa6ZBOH3DkVY59rlfE",
		"akhirbucket",
	}

	return val
}
