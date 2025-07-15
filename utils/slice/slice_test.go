package slice

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestStructToSlice(t *testing.T) {
	type order struct {
		OrderNo     string
		UserId      int64
		ProductName string
	}
	o := order{
		OrderNo:     "order_no",
		UserId:      1,
		ProductName: "test",
	}
	s1 := StructToSlice(o)
	t.Logf("%#v", s1) // []interface {}{"order_no", 1, "test"}
}

func TestStructToSliceByTagValues(t *testing.T) {
	type order struct {
		OrderNo     string `colName:"a"`
		UserId      int64  `colName:"b"`
		ProductName string `colName:"c"`
	}
	o := order{
		OrderNo:     "order_no",
		UserId:      1,
		ProductName: "test",
	}
	s1 := StructToSliceByTagValues(o, "colName", []string{"a", "c"})
	t.Logf("%#v", s1) // []interface {}{"order_no", "test"}
}

func TestPluck(t *testing.T) {
	type order struct {
		OrderNo string
		UserId  int64
	}
	o := []order{
		{
			OrderNo: "order_no",
			UserId:  1,
		},
		{
			OrderNo: "order_no2",
			UserId:  2,
		},
	}
	t.Log(SlicePluckInt(o, "UserId"))     // [1 2]
	t.Log(SlicePluckString(o, "OrderNo")) // [order_no order_no2]
}

func TestColumnSumDecimal(t *testing.T) {
	type order struct {
		OrderNo string
		Price   decimal.Decimal
	}
	b := []order{
		{
			OrderNo: "100011",
			Price:   decimal.NewFromInt(1),
		},
		{
			OrderNo: "100012",
			Price:   decimal.NewFromInt(2),
		},
	}

	t.Log(SliceColumnSumDecimal(b, "Price")) // 3
}

func TestGroupBy(t *testing.T) {
	type order struct {
		OrderNo string
		UserId  int64
	}
	s := []order{
		{
			OrderNo: "order_no1",
			UserId:  1,
		},
		{
			OrderNo: "order_no2",
			UserId:  1,
		},
		{
			OrderNo: "order_no2",
			UserId:  2,
		},
	}
	g1 := make(map[int64][]order)
	g2 := make(map[string][]order)
	SliceGroupBy(s, "UserId", g1)
	SliceGroupBy(s, "OrderNo", g2)
	t.Logf("%+v", g1) // map[1:[{OrderNo:order_no1 UserId:1} {OrderNo:order_no2 UserId:1}] 2:[{OrderNo:order_no2 UserId:2}]]
	t.Logf("%+v", g2) // map[order_no1:[{OrderNo:order_no1 UserId:1}] order_no2:[{OrderNo:order_no2 UserId:1} {OrderNo:order_no2 UserId:2}]]
}

func TestSliceColumnUnique(t *testing.T) {
	type order struct {
		OrderNo string
		UserId  int64
	}
	s := []order{
		{
			OrderNo: "order_no1",
			UserId:  1,
		},
		{
			OrderNo: "order_no1",
			UserId:  1,
		},
	}

	r1 := SliceColumnUniqueString("OrderNo", s)
	r2 := SliceColumnUniqueInt64("UserId", s)
	t.Log(r1) // [order_no1]
	t.Log(r2) // [1]
}

func TestSliceElemPos(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []int{1, 3, 4}
	s1Pos := SliceElemPos("b", s1)
	s2Pos := SliceElemPos(4, s2)
	t.Log(s1Pos) // 1
	t.Log(s2Pos) // 2
}
