package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/cecepsprd/foodstore-server/constans"
	"github.com/cecepsprd/foodstore-server/model"
	"github.com/cecepsprd/foodstore-server/utils/logger"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func ConvertPrimitiveID(id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	return objectID
}

func ConvertPrimitiveIDs(ids []string) []primitive.ObjectID {
	var primitiveIDs []primitive.ObjectID
	for _, id := range ids {
		primitiveIDs = append(primitiveIDs, ConvertPrimitiveID(id))
	}
	return primitiveIDs
}

func GetParamsValue(c echo.Context, from interface{}, to interface{}) {
	m := model.Paging{}
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		n := v.Type().Field(i).Name
		q := strings.ToLower(c.QueryParam(n))
		p := reflect.ValueOf(q)
		reflect.ValueOf(&m).Elem().FieldByName(n).Set(p)
	}
}

func HashPassword(password string) (hashedPassword string, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}

	return string(hashed), nil
}

func MappingRequest(request interface{}, model interface{}) error {
	// convert interface to json
	jsonRecords, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("Error encode records: %s", err.Error())
	}

	// bind json to struct
	if err := json.Unmarshal(jsonRecords, model); err != nil {
		return fmt.Errorf("Error decode json to struct: %s", err.Error())
	}

	return nil
}

func DecodeQueryParams(queryStr string, key string, req interface{}) error {
	params, err := url.ParseQuery(queryStr)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	err = json.Unmarshal([]byte(params.Get(key)), &req)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	return nil
}

func SetHTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case constans.ErrInternalServerError:
		return http.StatusInternalServerError
	case constans.ErrNotFound:
		return http.StatusNotFound
	case constans.ErrConflict:
		return http.StatusConflict
	case constans.ErrWrongEmailOrPassword:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
