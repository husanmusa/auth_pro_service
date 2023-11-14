package helper

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/saidamir98/udevs_pkg/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
	)

	for k, v := range params {
		if k != "" && strings.Contains(namedQuery, ":"+k) {
			namedQuery = strings.ReplaceAll(namedQuery, ":"+k, "$"+strconv.Itoa(i))
			args = append(args, v)
			i++
		}
	}

	return namedQuery, args
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func HandleError(log logger.LoggerI, err error, message string, req interface{}, code codes.Code) error {
	if code != 0 {
		log.Error(message, logger.Error(err), logger.Any("req", req))
		return status.Error(code, message)
	} else if err == sql.ErrNoRows {
		log.Error(message+", Not Found", logger.Error(err), logger.Any("req", req))
		return status.Error(codes.NotFound, "Not Found")
	} else if err != nil {
		log.Error(message, logger.Error(err), logger.Any("req", req))
		return status.Error(codes.Internal, message+err.Error())
	}
	return nil
}

func ToNullString(s *wrappers.StringValue) (res sql.NullString) {
	if s.GetValue() != "" {
		res.String = s.Value
		res.Valid = true
	}
	return res
}

func ToStringValue(s sql.NullString) *wrappers.StringValue {
	if s.Valid {
		return &wrappers.StringValue{Value: s.String}
	}
	return nil
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}
