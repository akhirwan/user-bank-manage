package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"user-bank-manage/config"

	"user-bank-manage/db"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-playground/validator/v10"
)

type BankManage struct {
	BankId         int    `json:"bank_id"`
	BankIdentifier string `json:"bank_identifier" validate:"required"`
	BankActive     string `json:"bank_active" validate:"required"`
	BankAddedOn    string `json:"bank_added_on" validate:"required"`
	BankDeleted    string `json:"bank_deleted" validate:"required"`
}

func FetchAllBanks() (AllResponse, error) {
	var obj BankManage
	var arrobj []BankManage
	var res AllResponse

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM tbl_user_bank_manage"

	rows, err := con.Query(sqlStatement)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.BankId, &obj.BankIdentifier, &obj.BankActive, &obj.BankAddedOn, &obj.BankDeleted)
		// log.Println(err)
		if err != nil {
			return res, err
		}

		arrobj = append(arrobj, obj)
		// fmt.Println(arrobj)
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Data bank"
	res.Length = len(arrobj)
	res.Result = arrobj

	// fmt.Println(res)

	return res, nil
}

func DetailBanks(bank_id int) (DetailResponse, error) {
	var obj BankManage
	var arrobj []BankManage
	var res DetailResponse

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM tbl_user_bank_manage WHERE bank_id = ? "

	rows, err := con.Query(sqlStatement, bank_id)
	defer rows.Close()

	// log.Println(rows, err)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.BankId, &obj.BankIdentifier, &obj.BankActive, &obj.BankAddedOn, &obj.BankDeleted)
		if err != nil {
			return res, err
		}

		// fmt.Println(arrobj)
	}
	arrobj = append(arrobj, obj)

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Detail bank"
	res.Id = obj.BankId
	res.Result = map[string]string{
		"bank_identifier": obj.BankIdentifier,
		"bank_active":     obj.BankActive,
		"bank_added_on":   obj.BankAddedOn,
		"bank_deleted":    obj.BankDeleted,
	}

	return res, nil
}

func StoreBanks(bank_identifier string, bank_active string, bank_added_on string, bank_deleted string) (Response, error) {
	var res Response

	v := validator.New()

	banks := BankManage{
		BankIdentifier: bank_identifier,
		BankActive:     bank_active,
		BankAddedOn:    bank_added_on,
		BankDeleted:    bank_deleted,
	}

	err := v.Struct(banks)
	if err != nil {
		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT tbl_user_bank_manage (bank_identifier, bank_active, bank_added_on, bank_deleted) VALUES (?, ?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(bank_identifier, bank_active, bank_added_on, bank_deleted)
	if err != nil {
		return res, err
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Data berhasil disimpan"
	res.Result = map[string]int64{
		"last_inserted_id": lastInsertedId,
	}

	return res, nil
}

// func UpdateBanks(bank_id int, bank_identifier string, bank_active string, bank_added_on string, bank_deleted string) (DetailResponse, error) {
func UpdateBanks(bank_id int, bank_identifier string, bank_active int) (DetailResponse, error) {
	var res DetailResponse

	con := db.CreateCon()

	// sqlStatement := "UPDATE tbl_user_bank_manage SET bank_identifier = ?, bank_active = ?, bank_added_on = ?, bank_deleted = ? WHERE bank_id = ?"
	sqlStatement := "UPDATE tbl_user_bank_manage SET bank_identifier = ?, bank_active = ? WHERE bank_id = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	// result, err := stmt.Exec(bank_identifier, bank_active, bank_added_on, bank_deleted, bank_id)
	result, err := stmt.Exec(bank_identifier, bank_active, bank_id)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Data berhasil diperbaharui"
	res.Id = bank_id
	res.Result = map[string]int64{
		"rows_affected": rowsAffected,
	}

	return res, nil
}

func DeleteBanks(bank_id int) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM tbl_user_bank_manage WHERE bank_id = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(bank_id)
	if err != nil {
		return res, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Data berhasil dihapus"
	res.Result = map[string]int64{
		"rows_affected": rowsAffected,
	}

	return res, nil
}

type AuthAdminToken struct {
	AdminId         int    `json:"admauth_admin_id"`
	AdminToken      string `json:"admauth_token" validate:"required"`
	AdminWebToken   string `json:"admauth_web_token" validate:"required"`
	AdminExpiry     string `json:"admauth_expiry" validate:"required"`
	AdminBrowser    string `json:"admauth_browser" validate:"required"`
	AdminLastAccess string `json:"admauth_last_access" validate:"required"`
	AdminLastIp     string `json:"admauth_last_ip" validate:"required"`
}

func getAdminToken(token string) (AllResponse, error) {
	var obj AuthAdminToken
	var arrobj []AuthAdminToken
	var res AllResponse

	con := db.CreateCon()

	sqlStatement := "SELECT * FROM tbl_admin_auth_token WHERE admauth_web_token = ?"

	rows, err := con.Query(sqlStatement, token)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.AdminId, &obj.AdminToken, &obj.AdminWebToken, &obj.AdminExpiry, &obj.AdminBrowser, &obj.AdminLastAccess, &obj.AdminLastIp)
		// log.Println(err)
		if err != nil {
			return res, err
		}

		arrobj = append(arrobj, obj)
		// fmt.Println(arrobj)
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Admin Authentication"
	res.Length = len(arrobj)
	res.Result = arrobj

	fmt.Println(res)

	return res, nil
}

type ObjectStorage struct {
	File string `json:"file_object" validate:"required"`
}

func ListObjectStorage() (AllResponse, error) {
	// var obj ObjectStorage
	var arrobj []ObjectStorage
	var res AllResponse

	conf := config.ObjectStorageConfig()

	// client, err := oss.New("Endpoint", "AccessKeyId", "AccessKeySecret")
	client, err := oss.New(conf.OSS_Endpoint, conf.OSS_AccessKeyId, conf.OSS_AccessKeySecret)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	bucket, err := client.Bucket(conf.OSS_Bucket)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	lsRes, err := bucket.ListObjects()
	if err != nil {
		// HandleError(err)
		return res, err
	}

	objStorage := lsRes.Objects
	for _, value := range objStorage {

		list := ObjectStorage{
			value.Key,
		}

		arrobj = append(arrobj, list)
	}

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "List Objects"
	res.Length = len(objStorage)
	res.Result = arrobj

	return res, nil
}

func UploadObjectStorage(fileName string) (AllResponse, error) {
	var res AllResponse
	conf := config.ObjectStorageConfig()

	folder := "test/"

	client, err := oss.New(conf.OSS_Endpoint, conf.OSS_AccessKeyId, conf.OSS_AccessKeySecret)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	bucket, err := client.Bucket(conf.OSS_Bucket)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	// tempFile, err := ioutil.TempFile("file-images", "prefix"+fileName)
	tempFile, err := ioutil.TempFile(os.TempDir(), fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// data, err := ioutil.ReadAll(tempFile)
	data, err := ioutil.ReadFile(tempFile.Name())
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(tmpFile.Name())
	if _, err := tempFile.Write(data); err != nil {
		fmt.Println(err)
	}

	// err = bucket.PutObjectFromFile(folder+file, "D:/hd-icon/"+file)
	err = bucket.PutObjectFromFile(folder+fileName, tempFile.Name())
	// err = bucket.PutObjectFromFile(folder+file, string(readFile))
	if err != nil {
		// HandleError(err)
		res.Status = http.StatusOK
		res.Success = true
		res.Message = "failed"
		return res, err
	}

	defer os.Remove(tempFile.Name())

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Object " + fileName
	// res.Result = data

	return res, nil
}

func MultipartObjectStorage(fileName string, tmpDir string) (AllResponse, error) {
	var res AllResponse
	conf := config.ObjectStorageConfig()

	folder := "test/"

	client, err := oss.New(conf.OSS_Endpoint, conf.OSS_AccessKeyId, conf.OSS_AccessKeySecret)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	bucket, err := client.Bucket(conf.OSS_Bucket)
	if err != nil {
		// HandleError(err)
		return res, err
	}

	// tempFile, err := ioutil.TempFile("file-images", "prefix"+fileName)
	// tempFile, err := ioutil.TempFile(os.TempDir(), fileName)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer tempFile.Close()

	// // data, err := ioutil.ReadAll(tempFile)
	// data, err := ioutil.ReadFile(tempFile.Name())
	// // if our program was unable to read the file
	// // print out the reason why it can't
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // fmt.Println(tmpFile.Name())
	// if _, err := tempFile.Write(data); err != nil {
	// 	fmt.Println(err)
	// }

	// err = bucket.PutObjectFromFile(folder+file, "D:/hd-icon/"+file)
	// err = bucket.PutObjectFromFile(folder+fileName, tempFile.Name())
	err = bucket.PutObjectFromFile(folder+fileName, tmpDir)
	if err != nil {
		// HandleError(err)
		res.Status = http.StatusOK
		res.Success = true
		res.Message = "failed"
		return res, err
	}

	// defer os.Remove(tempFile.Name())

	res.Status = http.StatusOK
	res.Success = true
	res.Message = "Object " + fileName
	// res.Result = data

	return res, nil
}
