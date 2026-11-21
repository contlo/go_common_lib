package gomysql

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, run func()) string {
	t.Helper()
	stdout, stderr := os.Stdout, os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout, os.Stderr = w, w
	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()
	func() {
		defer func() {
			os.Stdout, os.Stderr = stdout, stderr
			w.Close()
		}()
		run()
	}()
	return <-done
}

func TestInitDoesNotWriteThePassword(t *testing.T) {
	db := &MysqlDB{MysqlConfig: &MysqlConfigData{
		Database: "segment_db",
		Username: "segment_user",
		Password: "pw-must-not-leak-7f3a",
		Host:     "127.0.0.1",
		Port:     "1",
	}}

	var recovered interface{}
	output := captureOutput(t, func() {
		defer func() { recovered = recover() }()
		db.Init()
	})

	if recovered == nil {
		t.Fatal("expected Init to panic when the database is unreachable")
	}
	if strings.Contains(output, "pw-must-not-leak-7f3a") {
		t.Fatalf("Init wrote the password to the output: %d bytes written", len(output))
	}
	if strings.Contains(output, "segment_user:") {
		t.Fatalf("Init wrote the connection string to the output: %d bytes written", len(output))
	}
}
