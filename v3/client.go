package aliyun

import (
	openapiV2 "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	dysmsapi20180501 "github.com/alibabacloud-go/dysmsapi-20180501/v2/client"
	"github.com/ghinknet/smsutils/v3/errors"
	"github.com/ghinknet/smsutils/v3/model"
)

type Client struct {
	Globe         *dysmsapi20180501.Client
	CN            *dysmsapi20170525.Client
	GlobeTemplate map[string]string
	// JSON
	Marshal   func(any) ([]byte, error)
	Unmarshal func([]byte, any) error
}

type Driver struct{}

func (d Driver) NewClient(params model.DriverClientParam) (model.Client, error) {
	// Check credential
	keyID, keySecret := params.Credential[AccessKeyID], params.Credential[AccessKeySecret]
	if keyID == "" || keySecret == "" {
		return Client{}, errors.ErrDriverCredentialInvalid
	}

	// Try to decode globe template
	globeTemplate := make(map[string]string)
	if params.Credential[GlobeTemplate] != "" {
		if err := params.Unmarshal([]byte(params.Credential[GlobeTemplate]), &globeTemplate); err != nil {
			return nil, err
		}
	}

	// Struct aliyun client config
	clientConfigV2Globe := &openapiV2.Config{
		AccessKeyId:     &keyID,
		AccessKeySecret: &keySecret,
	}
	clientConfigV2CN := &openapiV2.Config{
		AccessKeyId:     &keyID,
		AccessKeySecret: &keySecret,
	}

	// Set aliyun endpoint
	clientConfigV2Globe.Endpoint = new(EndpointGlobe)
	clientConfigV2CN.Endpoint = new(EndpointCN)

	// Create aliyun client
	globeResult := new(dysmsapi20180501.Client)
	globeResult, err := dysmsapi20180501.NewClient(clientConfigV2Globe)
	if err != nil {
		return Client{}, err
	}
	cnResult := new(dysmsapi20170525.Client)
	cnResult, err = dysmsapi20170525.NewClient(clientConfigV2CN)
	if err != nil {
		return Client{}, err
	}

	return Client{
		Globe:         globeResult,
		CN:            cnResult,
		GlobeTemplate: globeTemplate,
		Marshal:       params.Marshal,
		Unmarshal:     params.Unmarshal,
	}, nil
}
