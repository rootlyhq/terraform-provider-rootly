package converter

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Flatten(i any, s *schema.Schema) (any, error) {
	switch {
	case s.Type == schema.TypeString:
		v, ok := i.(string)
		if !ok {
			return nil, errors.New("error converting to string")
		}
		return v, nil

	case s.Type == schema.TypeBool:
		v, ok := i.(bool)
		if !ok {
			return nil, errors.New("error converting to bool")
		}
		return v, nil

	case s.Type == schema.TypeInt:
		if v, ok := i.(int); ok {
			return v, nil
		} else if v, ok := i.(int16); ok {
			return v, nil
		} else if v, ok := i.(int32); ok {
			return v, nil
		} else if v, ok := i.(int64); ok {
			return v, nil
		}
		return nil, errors.New("error converting to int")

	case s.Type == schema.TypeFloat:
		if v, ok := i.(float32); ok {
			return v, nil
		} else if v, ok := i.(float64); ok {
			return v, nil
		}
		return nil, errors.New("error converting to float")

	case s.Type == schema.TypeList && s.MaxItems == 1 && isResource(s.Elem):
		if i == nil {
			return nil, nil
		}
		ii, ok := i.([]any)
		if !ok {
			return nil, errors.New("not a slice")
		} else if len(ii) == 0 {
			return map[string]any{}, nil
		} else if len(ii) > 1 {
			return nil, fmt.Errorf("got %d elements, want 1", len(ii))
		}

		iii, ok := ii[0].(map[string]any)
		if !ok {
			return nil, errors.New("not a slice with one map")
		}

		res := s.Elem.(*schema.Resource)
		m := make(map[string]any, len(res.Schema))
		for k, s := range res.Schema {
			v, err := Flatten(iii[k], s)
			if err != nil {
				return nil, err
			}
			m[k] = v
		}
		return m, nil

	case s.Type == schema.TypeList:
		if i == nil {
			return nil, nil
		}

		ii, ok := i.([]any)
		if !ok {
			return nil, errors.New("not a slice")
		}

		l := make([]any, len(ii))

		if sch, ok := s.Elem.(*schema.Schema); ok {
			for i := range len(ii) {
				v, err := Flatten(ii[i], sch)
				if err != nil {
					return nil, err
				}
				l[i] = v
			}
		} else if res, ok := s.Elem.(*schema.Resource); ok {
			for i := range len(ii) {
				iii, ok := ii[i].(map[string]any)
				if !ok {
					return nil, errors.New("not a map")
				}

				m := make(map[string]any, len(res.Schema))
				for k, s := range res.Schema {
					v, err := Flatten(iii[k], s)
					if err != nil {
						return nil, err
					}
					m[k] = v
				}

				l[i] = m
			}
		} else {
			return nil, errors.New("invalid list element")
		}

		return l, nil

	default:
		return nil, fmt.Errorf("flatten not implemented: %v+; %v+", i, s)
	}
}
