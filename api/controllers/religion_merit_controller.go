package controllers

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/phapsu86/tamlinh/api/auth"
	"github.com/phapsu86/tamlinh/api/models"
	"github.com/phapsu86/tamlinh/api/responses"

	//"github.com/phapsu86/tamlinh/api/utils/dictionary"
	"github.com/phapsu86/tamlinh/api/utils/formaterror"
	"github.com/phapsu86/tamlinh/api/utils/formatresult"
)

type ParamInfo struct {
	MeritType int `json:"merit_type"`
	Status    int `json:"status"`
	Page      int `json:"page"`
}

type MeritOrder struct {
	ObjectID    int              `json:"object_id"`
	Items       []MeritOrderItem `json:"items"`
	ObjectType  int8             `json:"object_type"`
	IsAnonymous int8             `json:"is_anonymous"`
	Notes       string           `json:"notes"`
	TextToStore string           `json:"text_to_store"`
}

type MeritOrderItem struct {
	OfferingType int8   `json:"offering_type"`
	OfferingID   uint64 `json:"offering_id"`
	Amount       uint64 `json:"amount"`
	Price        uint64 `json:"price"`
}

type MeritItemSearch struct {
	ObjectType int8   `json:"object_type"`
	ObjectID   uint64 `json:"object_id"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
}

func (server *Server) CreateMerit(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	fmt.Println("tham so truyen len " + string(body))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return

	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	item := MeritOrder{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Check tồn tại chùa hoặc sự kiện có hợp lệ không
	var objectName string

	if item.ObjectType == 0 {
		mdReligionItem := models.ReligionItem{}
		objectData, err := mdReligionItem.CheckReligionItem(server.DB, item.ObjectID)
		if err != nil {
			err := formaterror.ErrMsg{}
			err.Msg = "RELIGION_ITEM_NOT_FOUND"
			err.Status = 1
			responses.JSON(w, http.StatusOK, err)
			return

		}
		objectName = objectData.Name

	} else {

		mdReligionEvent := models.ReligionEvent{}
		objectData, err := mdReligionEvent.CheckEventForMerit(server.DB, item.ObjectID)
		if err != nil {
			err := formaterror.ErrMsg{}
			err.Msg = "RELIGION_EVENT_NOT_FOUND"
			err.Status = 1
			responses.JSON(w, http.StatusOK, err)
			return

		}
		objectName = objectData.Name

	}

	var orderItems []MeritOrderItem
	orderItems = item.Items
	output := make(map[int][]MeritOrderItem)
	//Tao ma md5 de phat sinh cho giao dich
	data := []byte("hello")
	fmt.Printf("%x", md5.Sum(data))

	hashOrder := xid.New().String()
	if err != nil {
		err := formaterror.ErrMsg{}
		err.Msg = "CAN_NOT_CREATE_HASH "
		err.Status = 1
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(orderItems) > 0 {
		var totalAmount uint64
		var amountPoint uint64
		var err error
		totalAmount = 0

		//	arrItemInfo := [][]MeritOrderItem{}
		for i, _ := range orderItems {
			fmt.Printf("xxxxxx %i\n", orderItems[i].OfferingID)
			fmt.Printf("amount %i\n", orderItems[i].Amount)
			fmt.Printf("type %i\n", orderItems[i].OfferingType)
			//kiêm tra nếu là point
			if orderItems[i].OfferingType == 1 {
				totalAmount += orderItems[i].Amount
				amountPoint = orderItems[i].Amount
				fmt.Printf("tong cong so tièn %d\n", totalAmount)

			} else if orderItems[i].OfferingType == 2 { //truong hop la vat pham

				offeringID := orderItems[i].OfferingID
				amount := orderItems[i].Amount
				//check ton tai vat pham ban cho su kien /chua
				offItemSell := models.OfferingItemSell{}
				ItemData, errEx := offItemSell.CheckItemSellExist(server.DB, offeringID, item.ObjectID, item.ObjectType)
				if errEx != nil {
					err := formaterror.ReturnErr(errEx)
					responses.JSON(w, http.StatusOK, err)
					return

				}

				price := ItemData.Price
				orderItems[i].Price = price
				fmt.Printf("mang thong tin cac phan tu %v\n", price)
				output[ItemData.StoreID] = append(output[ItemData.StoreID], orderItems[i])

				fmt.Printf("mang out put day %v\n", output)

				totalAmount += amount * price

			} else {
				//Tra ve bi loi
				err := formaterror.ErrMsg{}
				err.Msg = "FORMAT_PARAM_ERROR"
				err.Status = 1

				responses.JSON(w, http.StatusOK, err)
				return

			}

		}

		//Truong hop error  = nil thi check point cua khach hang
		fmt.Printf("tong cong total day %v\n", totalAmount)
		// Neu khong du point thi tra ve error
		mPoint := models.Point{}
		numP, errP := mPoint.FindPointByID(server.DB, uid)

		if errP != nil {
			err := formaterror.ReturnErr(errP)
			responses.JSON(w, http.StatusOK, err)
			return

		}
		if uint64(numP) < totalAmount {
			err := formaterror.ErrMsg{}
			err.Msg = "POINT_NOT_ENOUGH"
			err.Status = 2

			responses.JSON(w, http.StatusOK, err)
			return
		}
		//neu du point thi thuc hien tao order va treo point==========================
		//++ Thuc hien tru point cua khach hang
		mPoint.Amount = -int64(totalAmount)
		mPoint.CreatedAt = time.Now()
		mPoint.Note = "Charge for transaction order"
		mPoint.Type = "CHARGE"
		mPoint.RefID = hashOrder
		mPoint.UserID = uid
		pointCharged, err := mPoint.SavePoint(server.DB)

		if err != nil {

			err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
		}
		// Thuc hien tạo order

		if pointCharged != nil {
			resutlt := server.createOrder(amountPoint, output, hashOrder, item.ObjectID, item.ObjectType, uid, item.IsAnonymous, item.Notes, item.TextToStore, objectName)
			if resutlt == 1 {
				result := formaterror.ReturnSuccess()

				responses.JSON(w, http.StatusOK, result)
				return

			} else {
				err := errors.New("PROCESS_ERROR_ORDER")
				result := formaterror.ReturnErr(err)

				responses.JSON(w, http.StatusOK, result)
				return
			}

		}

	}

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// if uid != item.AuthorID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	//itemCreated, err := item.SaveReligionList(server.DB)
	err = nil
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, 1))
	responses.JSON(w, http.StatusCreated, 1)
}

type ResultBank struct {
	Status      int    `json:"status"`
	Info        string `json:"info"`
	BankAccount string `json:"bank_account"`
	Content     string `json:"content"`
}

func (server *Server) GetInfoBank(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
	mtype, err := strconv.Atoi(vars["type"])
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	result := ResultBank{}
	if mtype == 0 { //Công đức địa chỉ
		result.Info = "- Tài khoản kho bạc:\n+ Chủ tài khoản: Nguyễn Duy Hùng\n+ Kho bạc Nhà nước - Chi nhánh Hai Bà Trưng, Hà Nội\n+ Số tài khoản: 000005657788,888\n+ Nội dung chuyển tiền: Tiền ủng hộ chùa MC" + string(mid)
		result.Content = "Công đức MC" + string(mid)
	} else if mtype == 1 {
		result.Info = "- Tài khoản kho bạc:\n+ Chủ tài khoản: Nguyễn Duy Hùng\n+ Kho bạc Nhà nước - Chi nhánh Hai Bà Trưng, Hà Nội\n+ Số tài khoản: 000.005.657.788.888\n+ Nội dung chuyển tiền: Tiền ủng hộ SK" + string(mid)
		result.Content = "Công đức SK" + string(mid)
	} else { //công đức đến
		result.Info = "- Tài khoản kho bạc:\n+ Chủ tài khoản: Nguyễn Duy Hùng\n+ Kho bạc Nhà nước - Chi nhánh Hai Bà Trưng, Hà Nội\n+ Số tài khoản: 3713.0.1057713\n+ Nội dung chuyển tiền: Tiền ủng hộ cho anh Hùng làm app"
		result.Content = "Công đức SK" + string(mid)
	}

	result.Status = 200
	result.BankAccount = "000005657788888"
	responses.JSON(w, http.StatusOK, result)
	return
}

func (server *Server) GetInfoBankAddress(w http.ResponseWriter, r *http.Request) {

	result := ResultBank{}
	result.Info = "- Tài khoản kho bạc:\n+ Chủ tài khoản: Nguyễn Duy Hùng\n+ Kho bạc Nhà nước - Chi nhánh Hai Bà Trưng, Hà Nội\n+ Số tài khoản: 3713.0.1057713\n+ Nội dung chuyển tiền: Tiền ủng hộ cho anh Hùng làm app"
	result.Status = 200
	result.BankAccount = "000005657788888"
	result.Content = "Tiền ủng hộ cho địa chỉ"
	responses.JSON(w, http.StatusOK, result)
	return
}

func (server *Server) createOrder(totalPoint uint64, items map[int][]MeritOrderItem, hashOrder string, ObjectID int, ObjectType int8, UserID uint64, IsAnonymous int8, Notes string, TextToStore string, ObjectName string) int {
	//Tao order voi cong duc bang point
	if totalPoint > 0 {
		modelMeritList := models.MeritList{}
		modelMeritList.UserID = UserID
		modelMeritList.IsAnonymous = IsAnonymous
		modelMeritList.Notes = Notes
		modelMeritList.TextToStore = TextToStore
		modelMeritList.ObjectID = ObjectID
		modelMeritList.ObjectType = ObjectType
		modelMeritList.MeritType = 0
		modelMeritList.Amount = totalPoint
		modelMeritList.TransactionCode = hashOrder
		modelMeritList.ObjectName = ObjectName
		_, err := modelMeritList.SaveMeritList(server.DB)
		if err != nil {

			return 0
		}

	}

	var total uint64

	if len(items) > 0 {

		for k, v := range items {
			// check item cua store

			total = 0
			modeltDetails := []models.MeritDetail{}
			for j, _ := range v {
				itemDetail := v[j]
				itemMerit := models.MeritDetail{}
				itemMerit.Amount = itemDetail.Amount
				itemMerit.OfferingID = itemDetail.OfferingID
				itemMerit.Price = uint64(itemDetail.Price)
				itemMerit.StoreID = k

				modeltDetails = append(modeltDetails, itemMerit)
				total += uint64(itemDetail.Amount) * itemDetail.Price

			}
			modelMeritList := models.MeritList{}
			modelMeritList.UserID = UserID
			modelMeritList.IsAnonymous = IsAnonymous
			modelMeritList.Notes = Notes
			modelMeritList.ObjectID = ObjectID
			modelMeritList.ObjectType = ObjectType
			modelMeritList.MeritType = 1
			modelMeritList.Amount = total
			modelMeritList.TransactionCode = hashOrder
			modelMeritList.StoreID = k
			modelMeritList.ObjectName = ObjectName
			meritItem, err := modelMeritList.SaveMeritList(server.DB)
			if err != nil {

				return 0
			} else {
				for _, itemdt := range modeltDetails {
					itemdt.MeritID = meritItem.ID
					_, err := itemdt.SaveMeritDetail(server.DB)
					if err != nil {
						return 0
					}

				}
				//return 1
			}

			//Insert vao bang itemdetails

		}

	}

	return 1

}

// Lấy danh sách người cung tiên vật phẩm
func (server *Server) GetMeritList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
	mtype, err := strconv.Atoi(vars["type"])
	mp, err := strconv.Atoi(vars["page"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	model := models.MeritList{}

	itemReceived, err := model.FindMeritListByObjectID(server.DB, mid, mtype, mp)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	var rs interface{}
	rs = itemReceived

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
}

// Lấy danh sách người cung tiên vật phẩm
func (server *Server) GetOfferingItemSells(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
	mtype, err := strconv.Atoi(vars["type"])
	mp, err := strconv.Atoi(vars["page"])
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	model := models.OfferingItemSell{}

	items, err := model.GetOfferingItemSellForObject(server.DB, mid, mtype, mp)
	// if err != nil {

	// 	responses.JSON(w, http.StatusOK, [])
	// 	return
	// }

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.ItemDetail.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].ItemDetail.Image = linkImg
			fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//	result := formatresult.ReturnOfferingItemSell(itemReceived)
	//	responses.JSON(w, http.StatusOK, result)
}

func (server *Server) GetMeritUserInfo(w http.ResponseWriter, r *http.Request) {
	//resultFail := ResultFail{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return

	}
	item := ParamInfo{}
	err = json.Unmarshal(body, &item)
	// vars := mux.Vars(r)
	// status, err := strconv.Atoi(vars["status"])
	// page, err := strconv.Atoi(vars["page"])
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		err := formaterror.ReturnErr(errors.New("Unauthorized"))
		responses.JSON(w, http.StatusOK, err)
		return
	}
	merit := models.MeritList{}

	items, err := merit.GetMeritInfoByStatus(server.DB, uid, item.Status, item.MeritType, item.Page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			for j, x := range v.MeritDetail {
				linkImg := server.getLink(x.OfferingDetail.Image, "tamlinh")
				//fmt.Printf("dfdfdfd %v", linkImg)
				items[i].MeritDetail[j].OfferingDetail.Image = linkImg
				//fmt.Printf("dfdfdfd %v", v)
			}

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)
	// result := formatresult.ReturnResultMerit(items)
	// responses.JSON(w, http.StatusOK, result)
}

// Công đức cho đia chi
func (server *Server) CreateMeritAddress(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return

	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	item := models.MeritAddressList{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	hashOrder := xid.New().String()
	//Truong hop error  = nil thi check point cua khach hang
	fmt.Printf("tong cong total day %v\n", item.Amount)
	// Neu khong du point thi tra ve error
	mPoint := models.Point{}
	numP, errP := mPoint.FindPointByID(server.DB, uid)

	if errP != nil {
		err := formaterror.ReturnErr(errP)
		responses.JSON(w, http.StatusOK, err)
		return

	}
	if uint64(numP) < item.Amount {
		err := formaterror.ErrMsg{}
		err.Msg = "POINT_NOT_ENOUGH"
		err.Status = 1

		responses.JSON(w, http.StatusOK, err)
		return
	}
	//neu du point thi thuc hien tao order va treo point==========================
	//++ Thuc hien tru point cua khach hang
	mPoint.Amount = -int64(item.Amount)
	mPoint.CreatedAt = time.Now()
	mPoint.Note = "Charge for transaction merit adress"
	mPoint.Type = "CHARGE_MERIT_ADDRESS"
	mPoint.RefID = hashOrder
	mPoint.UserID = uid
	_, err = mPoint.SavePoint(server.DB)

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}
	// Tạo công đức cho dịa  chỉ
	item.TransactionCode = hashOrder
	_, err = item.SaveMeritAddressList(server.DB)

	if err != nil {
		formattedError := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusInternalServerError, formattedError)
		return
	}

	result := formaterror.ReturnSuccess()
	responses.JSON(w, http.StatusOK, result)

}

//Search item by keyword

func (server *Server) SearchOfferItemByNames(w http.ResponseWriter, r *http.Request) {

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modelSearch := MeritItemSearch{}
	err = json.Unmarshal(body, &modelSearch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//locID := modelSearch.LocID
	Keyword := modelSearch.Keyword
	ObjectType := modelSearch.ObjectType
	ObjectID := modelSearch.ObjectID
	page := modelSearch.Page

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.OfferingItemSell{}

	items, err := model.FindOfferingItemSellByKeyword(server.DB, ObjectID, ObjectType, Keyword, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
			//	fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	//result := formatresult.ReturnGlobalArray(itemReceived)
	//responses.JSON(w, http.StatusOK, result)
}

// Tìm kiêm sự kiện cho công đức

type EventSearchMerit struct {
	ReligionID int    `json:"religion_id"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
}

func (server *Server) SearchEventForMerit(w http.ResponseWriter, r *http.Request) {

	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	modelSearch := EventSearchMerit{}
	err = json.Unmarshal(body, &modelSearch)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)

		return
	}

	Keyword := modelSearch.Keyword

	ReligionID := modelSearch.ReligionID

	page := modelSearch.Page

	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}
	model := models.ReligionEvent{}

	items, err := model.SearchEventForMeritByName(server.DB, ReligionID, Keyword, page)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return
	}

	if len(items) > 0 {
		for i, v := range items {
			linkImg := server.getLink(v.Image, "tamlinh")
			//fmt.Printf("dfdfdfd %v", linkImg)
			items[i].Image = linkImg
			//	fmt.Printf("dfdfdfd %v", v)

		}
	}

	var rs interface{}
	rs = items

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

	// var rs interface{}
	// rs = itemReceived

	// result := formatresult.ReturnGlobalArray(rs)
	// responses.JSON(w, http.StatusOK, result)
}
