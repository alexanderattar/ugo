package api

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/brynbellomy/go-structomancer"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// Takes an *http.Request and decodes its contents into a struct containing special struct tags
// recognized by this function.  DecodeRequest can extract values from URL parameters, from the
// query string, and from the request body.  Additionally, it performs validation using the
// struct tags defined by gopkg.in/go-playground/validator.v9.
//
// For detailed usage and instructions, see README.md in this package.
func DecodeRequest(r *http.Request, dst interface{}) error {
	// A structomancer is a reflection tool that makes it possible to iterate over struct fields
	// as though they were map entries.  Here we initialize one and tell it to focus on fields that
	// are marked by the `api:` struct tag.
	z := structomancer.New(dst, "api")

	// Loop over the fields in the provided struct.
	for _, fname := range z.FieldNames() {
		field := z.Field(fname)

		// Decode the URL params and query string
		if field.IsFlagged("@query") || field.IsFlagged("@url_param") {
			var param string

			// Grab the value out of the query string or URL.
			if field.IsFlagged("@query") {
				param = r.URL.Query().Get(fname)
			} else if field.IsFlagged("@url_param") {
				param = chi.URLParam(r, fname)
			}

			if param == "" {
				continue
			}

			// All query string and URL params are returned as strings, so we do a bit of reflection
			// magic to convert them to the appropriate type.
			newVal, err := convertStringToFieldType(param, field.Type())
			if err != nil {
				return err
			}
			// Then we take the converted value and store it in the corresponding struct field.
			err = z.SetFieldValueV(reflect.ValueOf(dst), field.Nickname(), newVal)
			if err != nil {
				return err
			}

		} else if field.IsFlagged("@body") {
			// Decode the body

			// This reflection voodoo basically ensures that we're dealing with an "addressable"
			// pointer to the @body field on the struct.  Otherwise, when we attempt to store the
			// value in the @body field, the Go runtime will panic because it can't set unaddressable
			// values.
			rv := reflect.ValueOf(dst)
			if rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			rv = rv.FieldByName(field.Name())

			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				// If the field is a pointer, and is currently nil, fill it with a new empty struct
				rv.Set(reflect.New(rv.Type().Elem()))
			} else if rv.Kind() != reflect.Ptr {
				// If the field is a plain struct (not a pointer), get an addressable (mutable) pointer to it
				rv = rv.Addr()
			}
			// By this point, we should have an addressable pointer to this struct field, regardless of what type it is

			// If the @body field is a struct that supports the `render.Binder` interface, use that.
			// Otherwise, just use `render.Decode`.  This allows customization by the calling code
			// if desired, but a sensible fallback (and no extra code) if not.
			bodyDst := rv.Interface()
			if bodyBinder, isBinder := bodyDst.(render.Binder); isBinder {
				if err := render.Bind(r, bodyBinder); err != nil {
					return err
				}
			} else {
				if err := render.Decode(r, bodyDst); err != nil {
					return err
				}
			}
		}
	}

	// Run the validator
	return validate.Struct(dst)
}

func convertStringToFieldType(strValue string, fieldType reflect.Type) (reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.Ptr:
		innerVal, err := convertStringToFieldType(strValue, fieldType.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptrToVal := reflect.New(fieldType.Elem())
		ptrToVal.Elem().Set(innerVal)
		return ptrToVal, nil

	case reflect.String:
		return reflect.ValueOf(strValue), nil

	case reflect.Bool:
		b, err := strconv.ParseBool(strValue)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(b), nil

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(i), nil

	case reflect.Float32,
		reflect.Float64:
		f, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(f), nil

	default:
		return reflect.Value{}, errors.Errorf("unknown type %v", fieldType)
	}
}
