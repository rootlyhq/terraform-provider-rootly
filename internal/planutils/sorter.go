package planutils

import (
	"context"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/diagutils"
)

// KeyFunc computes a stable identity key for an item of type T, and
// reports whether one was available.
type KeyFunc[T any] func(item T) (key string, ok bool)

// KeyRegistry maps a Go struct type to the KeyFunc that identifies
// "the same logical item" for that type across API reads.
type KeyRegistry struct {
	fns map[reflect.Type]func(reflect.Value) (string, bool)
}

// NewKeyRegistry creates an empty registry. Populate it with RegisterKey.
func NewKeyRegistry() *KeyRegistry {
	return &KeyRegistry{fns: map[reflect.Type]func(reflect.Value) (string, bool){}}
}

// RegisterKey configures the stable key function for items of type T.
func RegisterKey[T any](r *KeyRegistry, fn KeyFunc[T]) {
	var zero T
	r.fns[reflect.TypeOf(zero)] = func(v reflect.Value) (string, bool) {
		item, ok := v.Interface().(T)
		if !ok {
			return "", false
		}
		return fn(item)
	}
}

func (r *KeyRegistry) keyFor(v reflect.Value) (string, bool) {
	if r == nil {
		return "", false
	}
	fn, ok := r.fns[v.Type()]
	if !ok {
		return "", false
	}
	return fn(v)
}

func SortNested[T any](
	ctx context.Context,
	keys *KeyRegistry,
	target interface {
		IsKnown() bool
		Get(context.Context) ([]*T, diag.Diagnostics)
		Set(context.Context, []*T) diag.Diagnostics
	},
	existing interface {
		IsKnown() bool
		Get(context.Context) ([]*T, diag.Diagnostics)
	},
) ([]MatchedItem[T], diag.Diagnostics) {
	var diags diag.Diagnostics
	if !target.IsKnown() {
		return nil, diags
	}
	newItems := diagutils.MergeDiagnostics(target.Get(ctx))(&diags)
	if diags.HasError() {
		return nil, diags
	}

	var existingItems []*T
	hasExisting := existing.IsKnown()
	if hasExisting {
		existingItems = diagutils.MergeDiagnostics(existing.Get(ctx))(&diags)
		if diags.HasError() {
			return nil, diags
		}
	}

	rMatches := matchByKeyReflect(reflect.ValueOf(newItems), reflect.ValueOf(existingItems), keys)

	for _, m := range rMatches {
		if !m.matchedPriorValid {
			continue
		}
		diags.Append(walkAndSortNestedFields(ctx, keys, m.item.Elem(), m.matchedPrior.Elem())...)
	}

	result := make([]MatchedItem[T], len(rMatches))
	items := make([]*T, len(rMatches))
	for i, m := range rMatches {
		item, _ := m.item.Interface().(*T)
		result[i] = MatchedItem[T]{Item: *item}
		if m.matchedPriorValid {
			p, _ := m.matchedPrior.Interface().(*T)
			result[i].MatchedPrior = p
		}
		items[i] = item
	}

	diags.Append(target.Set(ctx, items)...)
	return result, diags
}

func walkAndSortNestedFields(ctx context.Context, keys *KeyRegistry, newElem, existingElem reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics
	if newElem.Kind() != reflect.Struct || existingElem.Kind() != reflect.Struct {
		return diags
	}

	t := newElem.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if !isNestedListWrapperType(f.Type) {
			continue
		}
		diags.Append(sortWrapperFieldReflect(ctx, keys, newElem.Field(i), existingElem.Field(i))...)
	}
	return diags
}

func sortWrapperFieldReflect(ctx context.Context, keys *KeyRegistry, newField, existingField reflect.Value) diag.Diagnostics {
	var diags diag.Diagnostics
	if !callIsKnown(newField) {
		return diags
	}

	getRes := newField.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(ctx)})
	newSlice := getRes[0]
	appendDiagResult(&diags, getRes[1])

	hasExisting := callIsKnown(existingField)
	var existingSlice reflect.Value
	if hasExisting {
		er := existingField.MethodByName("Get").Call([]reflect.Value{reflect.ValueOf(ctx)})
		existingSlice = er[0]
		appendDiagResult(&diags, er[1])
	} else {
		existingSlice = reflect.Zero(newSlice.Type())
	}

	matches := matchByKeyReflect(newSlice, existingSlice, keys)

	for _, m := range matches {
		if !m.matchedPriorValid {
			continue
		}
		diags.Append(walkAndSortNestedFields(ctx, keys, m.item.Elem(), m.matchedPrior.Elem())...)
	}

	resultSlice := reflect.MakeSlice(newSlice.Type(), len(matches), len(matches))
	for i, m := range matches {
		resultSlice.Index(i).Set(m.item)
	}

	if !newField.CanAddr() {
		diags.AddError(
			"Unable to sort nested block",
			"internal error: nested field was not addressable when writing back sorted elements",
		)
		return diags
	}
	setRes := newField.Addr().MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(ctx), resultSlice})
	appendDiagResult(&diags, setRes[0])
	return diags
}

func matchByKeyReflect(newSlice, existingSlice reflect.Value, keys *KeyRegistry) []reflectMatch {
	n := newSlice.Len()
	if n == 0 {
		return nil
	}

	type existingEntry struct {
		idx int
		val reflect.Value
	}
	existingByKey := map[string]existingEntry{}
	seenKeys := map[string]bool{}
	if existingSlice.IsValid() {
		for j := 0; j < existingSlice.Len(); j++ {
			ev := existingSlice.Index(j)
			k, ok := keyFor(ev, keys)
			if !ok {
				continue
			}
			if seenKeys[k] {
				delete(existingByKey, k) // ambiguous duplicate key; don't match on it
				continue
			}
			seenKeys[k] = true
			existingByKey[k] = existingEntry{idx: j, val: ev}
		}
	}

	type ranked struct {
		idx       int
		priorRank int // -1 if unmatched
	}
	matchedVal := make([]reflect.Value, n)
	hasMatch := make([]bool, n)
	rankedItems := make([]ranked, n)

	for i := 0; i < n; i++ {
		nv := newSlice.Index(i)
		priorRank := -1
		if k, ok := keyFor(nv, keys); ok {
			if e, found := existingByKey[k]; found {
				priorRank = e.idx
				matchedVal[i] = e.val
				hasMatch[i] = true
			}
		}
		rankedItems[i] = ranked{idx: i, priorRank: priorRank}
	}

	sort.SliceStable(rankedItems, func(x, y int) bool {
		a, b := rankedItems[x], rankedItems[y]
		switch {
		case a.priorRank >= 0 && b.priorRank >= 0:
			return a.priorRank < b.priorRank
		case a.priorRank >= 0:
			return true
		case b.priorRank >= 0:
			return false
		default:
			return a.idx < b.idx
		}
	})

	result := make([]reflectMatch, n)
	for outIdx, r := range rankedItems {
		rm := reflectMatch{item: newSlice.Index(r.idx)}
		if hasMatch[r.idx] {
			rm.matchedPrior = matchedVal[r.idx]
			rm.matchedPriorValid = true
		}
		result[outIdx] = rm
	}
	return result
}

func keyFor(v reflect.Value, keys *KeyRegistry) (string, bool) {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", false
		}
		v = v.Elem()
	}
	return keys.keyFor(v)
}

func callIsKnown(v reflect.Value) bool {
	m := v.MethodByName("IsKnown")
	if !m.IsValid() {
		return false
	}
	res := m.Call(nil)
	if len(res) != 1 {
		return false
	}
	return res[0].Bool()
}

func appendDiagResult(diags *diag.Diagnostics, v reflect.Value) {
	if d, ok := v.Interface().(diag.Diagnostics); ok {
		diags.Append(d...)
	}
}

var (
	contextType     = reflect.TypeFor[*context.Context]()
	diagnosticsType = reflect.TypeFor[diag.Diagnostics]()
)

type MatchedItem[T any] struct {
	Item         T
	MatchedPrior *T
}

func isNestedListWrapperType(t reflect.Type) bool {
	getMethod, ok := t.MethodByName("Get")
	if !ok || getMethod.Type.NumIn() != 2 || getMethod.Type.NumOut() != 2 {
		return false
	}
	if getMethod.Type.In(1) != contextType {
		return false
	}
	out0 := getMethod.Type.Out(0)
	if out0.Kind() != reflect.Slice || out0.Elem().Kind() != reflect.Pointer {
		return false
	}
	if getMethod.Type.Out(1) != diagnosticsType {
		return false
	}

	ptrType := reflect.PointerTo(t)
	setMethod, ok := ptrType.MethodByName("Set")
	if !ok || setMethod.Type.NumIn() != 3 || setMethod.Type.NumOut() != 1 {
		return false
	}
	if setMethod.Type.In(2) != out0 {
		return false // Set's slice param must match Get's returned slice type
	}
	if setMethod.Type.Out(0) != diagnosticsType {
		return false
	}

	if _, ok := t.MethodByName("IsKnown"); !ok {
		return false
	}

	return true
}

type reflectMatch struct {
	item              reflect.Value
	matchedPrior      reflect.Value
	matchedPriorValid bool
}
