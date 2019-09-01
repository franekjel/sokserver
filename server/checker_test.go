package server

import (
	"testing"

	log "github.com/franekjel/sokserver/logger"

	"github.com/franekjel/sokserver/tasks"
)

func TestCompilationError(t *testing.T) {
	sub := tasks.Submission{Code: "int main(){return 0}"} //missing semicolon
	ok, err := compileCode(&sub)
	log.Info(err)
	if ok {
		t.Error("there is undetected missing semicolon")
	}
}

func TestCompilationOk(t *testing.T) {
	sub := tasks.Submission{Code: "int main(){return 0;}"} //good code
	ok, err := compileCode(&sub)
	log.Info(err)
	if !ok {
		t.Error("this is good code")
	}
}
