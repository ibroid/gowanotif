package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestGantiKodeNoHp(t *testing.T) {
	nomorHandphone := "6289633855678"

	nomorHandphone = strings.Replace(nomorHandphone, "62", "0", 1)

	fmt.Println(nomorHandphone)
}
