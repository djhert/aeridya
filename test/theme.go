package main

import (
	"fmt"
	"github.com/hlfstr/aeridya/router"
	"github.com/hlfstr/aeridya/router/staticrouter"
	"net/http"
)

type Router struct {
	staticrouter.Route
}
