package cloudpan

import (
	"encoding/xml"
	"fmt"
	"github.com/tickstep/cloudpan189-api/cloudpan/apierror"
	"github.com/tickstep/cloudpan189-api/cloudpan/apiutil"
	"github.com/tickstep/library-go/logger"
	"net/url"
	"strings"
)

func (p *PanClient) AppFamilyRenameFile(familyId int64, renameFileId, newName string) (*AppFileEntity, *apierror.ApiError) {
	fullUrl := &strings.Builder{}
	fmt.Fprintf(fullUrl, "%s/family/file/renameFile.action?familyId=%d&fileId=%s&destFileName=%s&%s",
		API_URL,
		familyId, renameFileId, url.QueryEscape(newName),
		apiutil.PcClientInfoSuffixParam())

	sessionKey := p.appToken.FamilySessionKey
	sessionSecret := p.appToken.FamilySessionSecret
	httpMethod := "GET"
	dateOfGmt := apiutil.DateOfGmtStr()
	headers := map[string]string {
		"Date": dateOfGmt,
		"SessionKey": sessionKey,
		"Signature": apiutil.SignatureOfHmac(sessionSecret, sessionKey, httpMethod, fullUrl.String(), dateOfGmt),
		"X-Request-ID": apiutil.XRequestId(),
	}

	logger.Verboseln("do request url: " + fullUrl.String())
	respBody, err1 := p.client.Fetch(httpMethod, fullUrl.String(), nil, headers)
	if err1 != nil {
		logger.Verboseln("AppFamilyRenameFile occurs error: ", err1.Error())
		return nil, apierror.NewApiErrorWithError(err1)
	}
	logger.Verboseln("response: " + string(respBody))

	er := &apierror.AppErrorXmlResp{}
	if err := xml.Unmarshal(respBody, er); err == nil {
		if er.Code != "" {
			if er.Code == "FileAlreadyExists" {
				return nil, apierror.NewApiError(apierror.ApiCodeFileAlreadyExisted, "文件已存在")
			}
		}
	}
	item := &AppFileEntity{}
	if err := xml.Unmarshal(respBody, item); err != nil {
		logger.Verboseln("AppFamilyRenameFile parse response failed")
		return nil, apierror.NewApiErrorWithError(err)
	}
	return item, nil
}
