package converter

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Expand(i any, s *schema.Schema) (any, error) {
	switch {
	case s.Type == schema.TypeString || s.Type == schema.TypeBool || s.Type == schema.TypeInt || s.Type == schema.TypeFloat:
		return i, nil

	case s.Type == schema.TypeList && s.MaxItems == 1 && isResource(s.Elem):
		if i == nil {
			return nil, nil
		}

		ii, ok := i.(map[string]any)
		if !ok {
			return nil, errors.New("not a map")
		}

		res := s.Elem.(*schema.Resource)
		m := make(map[string]any, len(res.Schema))
		for k, s := range res.Schema {
			v, err := Expand(ii[k], s)
			if err != nil {
				return nil, err
			}
			m[k] = v
		}

		return []any{m}, nil

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
				v, err := Expand(ii[i], sch)
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
					v, err := Expand(iii[k], s)
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
		return nil, fmt.Errorf("expand not implemented: %v+; %v+", i, s)
	}
}
