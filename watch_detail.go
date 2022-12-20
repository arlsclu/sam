package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	queryDetailUrl = `https://api-sams.walmartmobile.cn/api/v1/sams/goods-portal/spu/queryDetail`
)

const spuDetailReqTemp = `
{
    "source": "iOS",
    "channel": "1",
    "uid": "273522748",
    "storeInfoVOList": [
        {
            "storeId": "6558",
            "storeType": "16",
            "storeDeliveryAttr": [
                2,
                3,
                4,
                5,
                6,
                9,
                12,
                13
            ]
        },
        {
            "storeId": "9991",
            "storeType": "32",
            "storeDeliveryAttr": [
                10
            ]
        },
        {
            "storeId": "6696",
            "storeType": "4",
            "storeDeliveryAttr": [
                3,
                4
            ]
        },
        {
            "storeId": "9996",
            "storeType": "8",
            "storeDeliveryAttr": [
                1
            ]
        }
    ],
    "addressVO": {
        "detailAddress": "63号楼403室",
        "cityName": "上海市",
        "countryName": "中国",
        "districtName": "青浦区",
        "provinceName": "上海"
    },
    "storeId": "6696",
    "storeDeliveryTemplateId": "1254057665164892694",
    "spuId": %s,
    "areaBlockId": "312014188493013526"
}
`

type spuDetail struct {
	Data struct {
		SpuID           string        `json:"spuId"`
		HostItem        string        `json:"hostItem"`
		StoreID         string        `json:"storeId"`
		Title           string        `json:"title"`
		MasterBizType   int           `json:"masterBizType"`
		ViceBizType     int           `json:"viceBizType"`
		CategoryIDList  []string      `json:"categoryIdList"`
		Images          []string      `json:"images"`
		Videos          []string      `json:"videos"`
		DescVideo       []interface{} `json:"descVideo"`
		IsAvailable     bool          `json:"isAvailable"`
		IsPutOnSale     bool          `json:"isPutOnSale"`
		SevenDaysReturn bool          `json:"sevenDaysReturn"`
		Intro           string        `json:"intro"`
		SubTitle        string        `json:"subTitle"`
		BrandID         string        `json:"brandId"`
		Weight          float64       `json:"weight"`
		Desc            string        `json:"desc"`
		PriceInfo       []struct {
			PriceType     int    `json:"priceType"`
			Price         string `json:"price"`
			PriceTypeName string `json:"priceTypeName"`
		} `json:"priceInfo"`
		StockInfo struct {
			StockQuantity     int `json:"stockQuantity"`
			SafeStockQuantity int `json:"safeStockQuantity"`
			SoldQuantity      int `json:"soldQuantity"`
		} `json:"stockInfo"`
		LimitInfo []struct {
			LimitType int    `json:"limitType"`
			LimitNum  int    `json:"limitNum"`
			Text      string `json:"text"`
			StoreID   string `json:"storeId,omitempty"`
			CycleDays int    `json:"cycleDays"`
		} `json:"limitInfo"`
		PurchaseLimitText   string `json:"purchaseLimitText"`
		PurchaseLimitMinNum int    `json:"purchaseLimitMinNum"`
		TagInfo             []struct {
			Title    string `json:"title"`
			TagPlace int    `json:"tagPlace"`
			TagMark  string `json:"tagMark"`
		} `json:"tagInfo"`
		DeliveryAttr int  `json:"deliveryAttr"`
		Favorite     bool `json:"favorite"`
		Giveaway     bool `json:"giveaway"`
		SpuExtDTO    struct {
			SubTitle        string        `json:"subTitle"`
			Intro           string        `json:"intro"`
			HostUpc         []string      `json:"hostUpc"`
			DepartmentID    string        `json:"departmentId"`
			Valuable        bool          `json:"valuable"`
			DetailVideos    []interface{} `json:"detailVideos"`
			Temperature     float64       `json:"temperature"`
			Weight          float64       `json:"weight"`
			IsImport        bool          `json:"isImport"`
			DeliveryAttr    int           `json:"deliveryAttr"`
			SevenDaysReturn bool          `json:"sevenDaysReturn"`
			Giveaway        bool          `json:"giveaway"`
			IsRoutine       bool          `json:"isRoutine"`
			ThumbnailImage  string        `json:"thumbnailImage"`
			Status          int           `json:"status"`
		} `json:"spuExtDTO"`
		BeltInfo     []interface{} `json:"beltInfo"`
		Valuable     bool          `json:"valuable"`
		DetailVideos []interface{} `json:"detailVideos"`
		Temperature  float64       `json:"temperature"`
		IsImport     bool          `json:"isImport"`
		SpuSpecInfo  []interface{} `json:"spuSpecInfo"`
		SpecList     struct {
		} `json:"specList"`
		SpecInfo      []interface{} `json:"specInfo"`
		AttrGroupInfo []interface{} `json:"attrGroupInfo"`
		AttrInfo      []struct {
			AttrID        string `json:"attrId"`
			Title         string `json:"title"`
			AttrValueList []struct {
				AttrValueID string `json:"attrValueId"`
				Value       string `json:"value"`
			} `json:"attrValueList"`
		} `json:"attrInfo"`
		ExtendedWarrantyList      []interface{} `json:"extendedWarrantyList"`
		CouponContentList         []interface{} `json:"couponContentList"`
		CouponList                []interface{} `json:"couponList"`
		PromotionList             []interface{} `json:"promotionList"`
		PromotionDetailList       []interface{} `json:"promotionDetailList"`
		DeliveryCapacityCountList []struct {
			StrDate string `json:"strDate"`
			List    []struct {
				StartTime  string `json:"startTime"`
				EndTime    string `json:"endTime"`
				CloseDate  string `json:"closeDate"`
				CloseTime  string `json:"closeTime"`
				TimeISFull bool   `json:"timeISFull"`
				Disabled   bool   `json:"disabled"`
			} `json:"list"`
		} `json:"deliveryCapacityCountList"`
		IsCollectOrder int `json:"isCollectOrder"`
		ComplianceInfo struct {
			ID    string `json:"id"`
			Value string `json:"value"`
		} `json:"complianceInfo"`
		PreSellList   []interface{} `json:"preSellList"`
		OnlyStoreSale bool          `json:"onlyStoreSale"`
		ServiceInfo   []interface{} `json:"serviceInfo"`
		IsStoreExtent bool          `json:"isStoreExtent"`
		IsTicket      bool          `json:"isTicket"`
	} `json:"data"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	ErrorMsg  string `json:"errorMsg"`
	TraceID   string `json:"traceId"`
	RequestID string `json:"requestId"`
	Rt        int    `json:"rt"`
	Success   bool   `json:"success"`
}

// func  to check if fullfilled  by  detail
func checkerByDetail(spuID string) (bool, error) {

	resp, err := http.Post(queryDetailUrl, "application/json", strings.NewReader(fmt.Sprintf(spuDetailReqTemp, spuID)))
	if err != nil {
		return false, err
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}
	var c = spuDetail{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return false, err
	}
	if !c.Success {
		return false, errors.New("返回状态success字段异常")
	}
	if c.Data.StockInfo.StockQuantity > 0 {
		return true, nil
	}
	return false, nil
}

// others
// unimplemented
func checkerByOther(name string) (bool, error) {
	return false, nil
}
