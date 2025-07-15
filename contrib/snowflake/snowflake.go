package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var n *snowflake.Node

func NewNode(node int64) (*snowflake.Node, error) {
	var err error
	n, err = snowflake.NewNode(node)
	return n, err
}

func Generate() snowflake.ID {
	if n == nil {
		panic("snowflake node not init")
	}
	return n.Generate()
}
