package govalidator

import (
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var rulesFuncMap = make(map[string]func(string, string, string, interface{}) error)

// AddCustomRule help to add custom rules for validator
// First argument it takes the rule name and second arg a func
// Second arg must have this signature below
// fn func(name string, fn func(field string, rule string, message string, value interface{}) error
// see example in readme: https://github.com/thedevsaddam/govalidator#add-custom-rules
func AddCustomRule(name string, fn func(field string, rule string, message string, value interface{}) error) {
	if isRuleExist(name) {
		panic(fmt.Errorf("govalidator: %s is already defined in rules", name))
	}
	rulesFuncMap[name] = fn
}

// validateCustomRules validate custom rules
func validateCustomRules(field string, rule string, message string, value interface{}, errsBag url.Values) {
	for k, v := range rulesFuncMap {
		if k == rule || strings.HasPrefix(rule, k+":") {
			err := v(field, rule, message, value)
			if err != nil {
				errsBag.Add(field, err.Error())
			}
			break
		}
	}
}

func init() {

	// Required check the Required fields
	AddCustomRule("required", func(field, rule, message string, value interface{}) error {
		err := fmt.Errorf("The %s field is required", field)
		if message != "" {
			err = errors.New(message)
		}
		if value == nil {
			return err
		}
		if _, ok := value.(multipart.File); ok {
			return nil
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if rv.Len() == 0 {
				return err
			}
		case reflect.Int:
			if isEmpty(value.(int)) {
				return err
			}
		case reflect.Int8:
			if isEmpty(value.(int8)) {
				return err
			}
		case reflect.Int16:
			if isEmpty(value.(int16)) {
				return err
			}
		case reflect.Int32:
			if isEmpty(value.(int32)) {
				return err
			}
		case reflect.Int64:
			if isEmpty(value.(int64)) {
				return err
			}
		case reflect.Float32:
			if isEmpty(value.(float32)) {
				return err
			}
		case reflect.Float64:
			if isEmpty(value.(float64)) {
				return err
			}
		case reflect.Uint:
			if isEmpty(value.(uint)) {
				return err
			}
		case reflect.Uint8:
			if isEmpty(value.(uint8)) {
				return err
			}
		case reflect.Uint16:
			if isEmpty(value.(uint16)) {
				return err
			}
		case reflect.Uint32:
			if isEmpty(value.(uint32)) {
				return err
			}
		case reflect.Uint64:
			if isEmpty(value.(uint64)) {
				return err
			}
		case reflect.Uintptr:
			if isEmpty(value.(uintptr)) {
				return err
			}
		case reflect.Struct:
			switch rv.Type().String() {
			case "govalidator.Int":
				if v, ok := value.(Int); ok {
					if !v.IsSet {
						return err
					}
				}
			case "govalidator.Int64":
				if v, ok := value.(Int64); ok {
					if !v.IsSet {
						return err
					}
				}
			case "govalidator.Float32":
				if v, ok := value.(Float32); ok {
					if !v.IsSet {
						return err
					}
				}
			case "govalidator.Float64":
				if v, ok := value.(Float64); ok {
					if !v.IsSet {
						return err
					}
				}
			case "govalidator.Bool":
				if v, ok := value.(Bool); ok {
					if !v.IsSet {
						return err
					}
				}
			default:
				panic("govalidator: invalid custom type for required rule")

			}

		default:
			panic("govalidator: invalid type for required rule")

		}
		return nil
	})

	// Regex check the custom Regex rules
	// Regex:^[a-zA-Z]+$ means this field can only contain alphabet (a-z and A-Z)
	AddCustomRule("regex", func(field, rule, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field format is invalid", field)
		if message != "" {
			err = errors.New(message)
		}
		rxStr := strings.TrimPrefix(rule, "regex:")
		if !isMatchedRegex(rxStr, str) {
			return err
		}
		return nil
	})

	// Alpha check if provided field contains valid letters
	AddCustomRule("alpha", func(field string, vlaue string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s may only contain letters", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isAlpha(str) {
			return err
		}
		return nil
	})

	// AlphaDash check if provided field contains valid letters, numbers, underscore and dash
	AddCustomRule("alpha_dash", func(field string, vlaue string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s may only contain letters, numbers, and dashes", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isAlphaDash(str) {
			return err
		}
		return nil
	})

	// AlphaDash check if provided field contains valid letters, numbers, underscore and dash
	AddCustomRule("alpha_space", func(field string, vlaue string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s may only contain letters, numbers, dashes, space", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isAlphaSpace(str) {
			return err
		}
		return nil
	})

	// AlphaNumeric check if provided field contains valid letters and numbers
	AddCustomRule("alpha_num", func(field string, vlaue string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s may only contain letters and numbers", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isAlphaNumeric(str) {
			return err
		}
		return nil
	})

	// Boolean check if provided field contains Boolean
	// in this case: "0", "1", 0, 1, "true", "false", true, false etc
	AddCustomRule("bool", func(field string, vlaue string, message string, value interface{}) error {
		err := fmt.Errorf("The %s may only contain boolean value, string or int 0, 1", field)
		if message != "" {
			err = errors.New(message)
		}
		switch t := value.(type) {
		case bool:
			//if value is boolean then pass
		case string:
			if !isBoolean(t) {
				return err
			}
		case int:
			if t != 0 && t != 1 {
				return err
			}
		case int8:
			if t != 0 && t != 1 {
				return err
			}
		case int16:
			if t != 0 && t != 1 {
				return err
			}
		case int32:
			if t != 0 && t != 1 {
				return err
			}
		case int64:
			if t != 0 && t != 1 {
				return err
			}
		case uint:
			if t != 0 && t != 1 {
				return err
			}
		case uint8:
			if t != 0 && t != 1 {
				return err
			}
		case uint16:
			if t != 0 && t != 1 {
				return err
			}
		case uint32:
			if t != 0 && t != 1 {
				return err
			}
		case uint64:
			if t != 0 && t != 1 {
				return err
			}
		case uintptr:
			if t != 0 && t != 1 {
				return err
			}
		}
		return nil
	})

	// Between check the fields character length range
	// if the field is array, map, slice then the valdiation rule will be the length of the data
	// if the value is int or float then the valdiation rule will be the value comparison
	AddCustomRule("between", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "between:"), ",")
		if len(rng) != 2 {
			panic(errInvalidArgument)
		}
		minFloat, err := strconv.ParseFloat(rng[0], 64)
		if err != nil {
			panic(errStringToInt)
		}
		maxFloat, err := strconv.ParseFloat(rng[1], 64)
		if err != nil {
			panic(errStringToInt)
		}
		min := int(minFloat)

		max := int(maxFloat)

		err = fmt.Errorf("The %s field must be between %d and %d", field, min, max)
		if message != "" {
			err = fmt.Errorf(message, min, max)
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
			inLen := rv.Len()
			if !(inLen >= min && inLen <= max) {
				return err
			}
		case reflect.Int:
			in := value.(int)
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Int8:
			in := int(value.(int8))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Int16:
			in := int(value.(int16))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Int32:
			in := int(value.(int32))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Int64:
			in := int(value.(int64))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uint:
			in := int(value.(uint))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uint8:
			in := int(value.(uint8))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uint16:
			in := int(value.(uint16))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uint32:
			in := int(value.(uint32))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uint64:
			in := int(value.(uint64))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Uintptr:
			in := int(value.(uintptr))
			if !(in >= min && in <= max) {
				return err
			}
		case reflect.Float32:
			in := float64(value.(float32))
			if !(in >= minFloat && in <= maxFloat) {
				return fmt.Errorf("The %s field must be between %f and %f", field, minFloat, maxFloat)
			}
		case reflect.Float64:
			in := value.(float64)
			if !(in >= minFloat && in <= maxFloat) {
				return fmt.Errorf("The %s field must be between %f and %f", field, minFloat, maxFloat)
			}

		}

		return nil
	})

	// CreditCard check if provided field contains valid credit card number
	// Accepted cards are Visa, MasterCard, American Express, Diners Club, Discover and JCB card
	AddCustomRule("credit_card", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid credit card number", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isCreditCard(str) {
			return err
		}
		return nil
	})

	// Coordinate check if provided field contains valid Coordinate
	AddCustomRule("coordinate", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid coordinate", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isCoordinate(str) {
			return err
		}
		return nil
	})

	// ValidateCSSColor check if provided field contains a valid CSS color code
	AddCustomRule("css_color", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid CSS color code", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isCSSColor(str) {
			return err
		}
		return nil
	})

	// Digits check the exact matching length of digit (0,9)
	// Digits:5 means the field must have 5 digit of length.
	// e.g: 12345 or 98997 etc
	AddCustomRule("digits", func(field string, rule string, message string, value interface{}) error {
		l, err := strconv.Atoi(strings.TrimPrefix(rule, "digits:"))
		if err != nil {
			panic(errStringToInt)
		}
		err = fmt.Errorf("The %s field must be %d digits", field, l)
		if l == 1 {
			err = fmt.Errorf("The %s field must be 1 digit", field)
		}
		if message != "" {
			err = fmt.Errorf(message, l)
		}
		var str string
		switch v := value.(type) {
		case string:
			str = v
		case float64:
			str = toString(int64(v))
		case float32:
			str = toString(int64(v))
		default:
			str = toString(v)
		}
		if len(str) != l || !regexDigits.MatchString(str) {
			return err
		}

		return nil
	})

	// DigitsBetween check if the field contains only digit and length between provided range
	// e.g: digits_between:4,5 means the field can have value like: 8887 or 12345 etc
	AddCustomRule("digits_between", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "digits_between:"), ",")
		if len(rng) != 2 {
			panic(errInvalidArgument)
		}
		min, err := strconv.Atoi(rng[0])
		if err != nil {
			panic(errStringToInt)
		}
		max, err := strconv.Atoi(rng[1])
		if err != nil {
			panic(errStringToInt)
		}
		err = fmt.Errorf("The %s field must be digits between %d and %d", field, min, max)
		if message != "" {
			err = fmt.Errorf(message, min, max)
		}
		str := toString(value)
		if !isNumeric(str) || !(len(str) >= min && len(str) <= max) {
			return err
		}

		return nil
	})

	// Date check the provided field is valid Date
	AddCustomRule("date", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)

		switch rule {
		case "date:dd-mm-yyyy":
			if !isDateDDMMYY(str) {
				if message != "" {
					return errors.New(message)
				}
				return fmt.Errorf("The %s field must be a valid date format. e.g: dd-mm-yyyy, dd/mm/yyyy etc", field)
			}
		default:
			if !isDate(str) {
				if message != "" {
					return errors.New(message)
				}
				return fmt.Errorf("The %s field must be a valid date format. e.g: yyyy-mm-dd, yyyy/mm/dd etc", field)
			}
		}

		return nil
	})

	// Email check the provided field is valid Email
	AddCustomRule("email", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid email address", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isEmail(str) {
			return err
		}
		return nil
	})

	// validFloat check the provided field is valid float number
	AddCustomRule("float", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a float number", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isFloat(str) {
			return err
		}
		return nil
	})

	// IP check if provided field is valid IP address
	AddCustomRule("ip", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid IP address", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isIP(str) {
			return err
		}
		return nil
	})

	// IP check if provided field is valid IP v4 address
	AddCustomRule("ip_v4", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid IPv4 address", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isIPV4(str) {
			return err
		}
		return nil
	})

	// IP check if provided field is valid IP v6 address
	AddCustomRule("ip_v6", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid IPv6 address", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isIPV6(str) {
			return err
		}
		return nil
	})

	// ValidateJSON check if provided field contains valid json string
	AddCustomRule("json", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid JSON string", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isJSON(str) {
			return err
		}
		return nil
	})

	/// Latitude check if provided field contains valid Latitude
	AddCustomRule("lat", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid latitude", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isLatitude(str) {
			return err
		}
		return nil
	})

	// Longitude check if provided field contains valid Longitude
	AddCustomRule("lon", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid longitude", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isLongitude(str) {
			return err
		}
		return nil
	})

	// Length check the field's character Length
	AddCustomRule("len", func(field string, rule string, message string, value interface{}) error {
		l, err := strconv.Atoi(strings.TrimPrefix(rule, "len:"))
		if err != nil {
			panic(errStringToInt)
		}
		err = fmt.Errorf("The %s field must be length of %d", field, l)
		if message != "" {
			err = fmt.Errorf(message, l)
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Map, reflect.Slice:
			vLen := rv.Len()
			if vLen != l {
				return err
			}
		default:
			str := toString(value) //force the value to be string
			if len(str) != l {
				return err
			}
		}

		return nil
	})

	// Min check the field's minimum character length for string, value for int, float and size for array, map, slice
	AddCustomRule("min", func(field string, rule string, message string, value interface{}) error {
		mustLen := strings.TrimPrefix(rule, "min:")
		lenInt, err := strconv.Atoi(mustLen)
		if err != nil {
			panic(errStringToInt)
		}
		lenFloat, err := strconv.ParseFloat(mustLen, 64)
		if err != nil {
			panic(errStringToFloat)
		}
		errMsg := fmt.Errorf("The %s field value can not be less than %d", field, lenInt)
		if message != "" {
			errMsg = fmt.Errorf(message, lenInt)
		}
		errMsgFloat := fmt.Errorf("The %s field value can not be less than %f", field, lenFloat)
		if message != "" {
			errMsgFloat = fmt.Errorf(message, lenFloat)
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String:
			inLen := rv.Len()
			if inLen < lenInt {
				if message != "" {
					return fmt.Errorf(message, lenInt)
				}
				return fmt.Errorf("The %s field must be minimum %d char", field, lenInt)
			}
		case reflect.Array, reflect.Map, reflect.Slice:
			inLen := rv.Len()
			if inLen < lenInt {
				if message != "" {
					return fmt.Errorf(message, lenInt)
				}
				return fmt.Errorf("The %s field must be minimum %d in size", field, lenInt)
			}
		case reflect.Int:
			in := value.(int)
			if in < lenInt {
				return errMsg
			}
		case reflect.Int8:
			in := int(value.(int8))
			if in < lenInt {
				return errMsg
			}
		case reflect.Int16:
			in := int(value.(int16))
			if in < lenInt {
				return errMsg
			}
		case reflect.Int32:
			in := int(value.(int32))
			if in < lenInt {
				return errMsg
			}
		case reflect.Int64:
			in := int(value.(int64))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uint:
			in := int(value.(uint))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uint8:
			in := int(value.(uint8))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uint16:
			in := int(value.(uint16))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uint32:
			in := int(value.(uint32))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uint64:
			in := int(value.(uint64))
			if in < lenInt {
				return errMsg
			}
		case reflect.Uintptr:
			in := int(value.(uintptr))
			if in < lenInt {
				return errMsg
			}
		case reflect.Float32:
			in := value.(float32)
			if in < float32(lenFloat) {
				return errMsgFloat
			}
		case reflect.Float64:
			in := value.(float64)
			if in < lenFloat {
				return errMsgFloat
			}

		}

		return nil
	})

	// Max check the field's maximum character length for string, value for int, float and size for array, map, slice
	AddCustomRule("max", func(field string, rule string, message string, value interface{}) error {
		mustLen := strings.TrimPrefix(rule, "max:")
		lenInt, err := strconv.Atoi(mustLen)
		if err != nil {
			panic(errStringToInt)
		}
		lenFloat, err := strconv.ParseFloat(mustLen, 64)
		if err != nil {
			panic(errStringToFloat)
		}
		errMsg := fmt.Errorf("The %s field value can not be greater than %d", field, lenInt)
		if message != "" {
			errMsg = fmt.Errorf(message, lenInt)
		}
		errMsgFloat := fmt.Errorf("The %s field value can not be greater than %f", field, lenFloat)
		if message != "" {
			errMsgFloat = fmt.Errorf(message, lenFloat)
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String:
			inLen := rv.Len()
			if inLen > lenInt {
				if message != "" {
					return fmt.Errorf(message, lenInt)
				}
				return fmt.Errorf("The %s field must be maximum %d char", field, lenInt)
			}
		case reflect.Array, reflect.Map, reflect.Slice:
			inLen := rv.Len()
			if inLen > lenInt {
				if message != "" {
					return fmt.Errorf(message, lenInt)
				}
				return fmt.Errorf("The %s field must be minimum %d in size", field, lenInt)
			}
		case reflect.Int:
			in := value.(int)
			if in > lenInt {
				return errMsg
			}
		case reflect.Int8:
			in := int(value.(int8))
			if in > lenInt {
				return errMsg
			}
		case reflect.Int16:
			in := int(value.(int16))
			if in > lenInt {
				return errMsg
			}
		case reflect.Int32:
			in := int(value.(int32))
			if in > lenInt {
				return errMsg
			}
		case reflect.Int64:
			in := int(value.(int64))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uint:
			in := int(value.(uint))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uint8:
			in := int(value.(uint8))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uint16:
			in := int(value.(uint16))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uint32:
			in := int(value.(uint32))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uint64:
			in := int(value.(uint64))
			if in > lenInt {
				return errMsg
			}
		case reflect.Uintptr:
			in := int(value.(uintptr))
			if in > lenInt {
				return errMsg
			}
		case reflect.Float32:
			in := value.(float32)
			if in > float32(lenFloat) {
				return errMsgFloat
			}
		case reflect.Float64:
			in := value.(float64)
			if in > lenFloat {
				return errMsgFloat
			}

		}

		return nil
	})

	// Numeric check if the value of the field is Numeric
	AddCustomRule("mac_address", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be a valid Mac Address", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isMacAddress(str) {
			return err
		}
		return nil
	})

	// Numeric check if the value of the field is Numeric
	AddCustomRule("numeric", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must be numeric", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isNumeric(str) {
			return err
		}
		return nil
	})

	// NumericBetween check if the value field numeric value range
	// e.g: numeric_between:18, 65 means number value must be in between a numeric value 18 & 65
	// Both of the bounds can be omited turning it into a min only (`10,`) or a max only (`,10`)
	AddCustomRule("numeric_between", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "numeric_between:"), ",")
		if len(rng) != 2 {
			panic(errInvalidArgument)
		}

		if rng[0] == "" && rng[1] == "" {
			panic(errInvalidArgument)
		}

		// check for integer value
		min := math.MinInt64
		if rng[0] != "" {
			_min, err := strconv.ParseFloat(rng[0], 64)
			if err != nil {
				panic(errStringToInt)
			}
			min = int(_min)
		}

		max := math.MaxInt64
		if rng[1] != "" {
			_max, err := strconv.ParseFloat(rng[1], 64)
			if err != nil {
				panic(errStringToInt)
			}
			max = int(_max)
		}

		var errMsg error
		switch {
		case rng[0] == "":
			errMsg = fmt.Errorf("The %s field value can not be greater than %d", field, max)
		case rng[1] == "":
			errMsg = fmt.Errorf("The %s field value can not be less than %d", field, min)
		default:
			errMsg = fmt.Errorf("The %s field must be numeric value between %d and %d", field, min, max)
		}

		if message != "" {
			switch {
			case rng[0] == "":
				errMsg = fmt.Errorf(message, max)
			case rng[1] == "":
				errMsg = fmt.Errorf(message, min)
			default:
				errMsg = fmt.Errorf(message, min, max)
			}
		}

		val := toString(value)

		if !strings.Contains(rng[0], ".") || !strings.Contains(rng[1], ".") {
			digit, errs := strconv.Atoi(val)
			if errs != nil {
				return errMsg
			}
			if !(digit >= min && digit <= max) {
				return errMsg
			}
		}
		// check for float value
		var err error
		minFloat := -math.MaxFloat64
		if rng[0] != "" {
			minFloat, err = strconv.ParseFloat(rng[0], 64)
			if err != nil {
				panic(errStringToFloat)
			}
		}

		maxFloat := math.MaxFloat64
		if rng[1] != "" {
			maxFloat, err = strconv.ParseFloat(rng[1], 64)
			if err != nil {
				panic(errStringToFloat)
			}
		}

		switch {
		case rng[0] == "":
			errMsg = fmt.Errorf("The %s field value can not be greater than %f", field, maxFloat)
		case rng[1] == "":
			errMsg = fmt.Errorf("The %s field value can not be less than %f", field, minFloat)
		default:
			errMsg = fmt.Errorf("The %s field must be numeric value between %f and %f", field, minFloat, maxFloat)
		}

		if message != "" {
			switch {
			case rng[0] == "":
				errMsg = fmt.Errorf(message, maxFloat)
			case rng[1] == "":
				errMsg = fmt.Errorf(message, minFloat)
			default:
				errMsg = fmt.Errorf(message, minFloat, maxFloat)
			}
		}

		digit, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errMsg
		}
		if !(digit >= minFloat && digit <= maxFloat) {
			return errMsg
		}
		return nil
	})

	// ValidateURL check if provided field is valid URL
	AddCustomRule("url", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field format is invalid", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isURL(str) {
			return err
		}
		return nil
	})

	// UUID check if provided field contains valid UUID
	AddCustomRule("uuid", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid UUID", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isUUID(str) {
			return err
		}
		return nil
	})

	// UUID3 check if provided field contains valid UUID of version 3
	AddCustomRule("uuid_v3", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid UUID V3", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isUUID3(str) {
			return err
		}
		return nil
	})

	// UUID4 check if provided field contains valid UUID of version 4
	AddCustomRule("uuid_v4", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid UUID V4", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isUUID4(str) {
			return err
		}
		return nil
	})

	// UUID5 check if provided field contains valid UUID of version 5
	AddCustomRule("uuid_v5", func(field string, rule string, message string, value interface{}) error {
		str := toString(value)
		err := fmt.Errorf("The %s field must contain valid UUID V5", field)
		if message != "" {
			err = errors.New(message)
		}
		if !isUUID5(str) {
			return err
		}
		return nil
	})

	// In check if provided field equals one of the values specified in the rule
	AddCustomRule("in", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "in:"), ",")
		if len(rng) == 0 {
			panic(errInvalidArgument)
		}
		str := toString(value)
		err := fmt.Errorf("The %s field must be one of %v", field, strings.Join(rng, ", "))
		if message != "" {
			err = fmt.Errorf(message, strings.Join(rng, ", "))
		}
		if !isIn(rng, str) {
			return err
		}
		return nil
	})

	// In check if provided field equals one of the values specified in the rule
	AddCustomRule("not_in", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_in:"), ",")
		if len(rng) == 0 {
			panic(errInvalidArgument)
		}
		str := toString(value)
		err := fmt.Errorf("The %s field must not be any of %v", field, strings.Join(rng, ", "))
		if message != "" {
			err = fmt.Errorf(message, strings.Join(rng, ", "))
		}
		if isIn(rng, str) {
			return err
		}
		return nil
	})
}
