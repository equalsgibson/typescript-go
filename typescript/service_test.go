package typescript_test

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aaronellington/typescript-go/typescript"
)

func TestPrimary(t *testing.T) {
	type UserID uint64

	type Group struct {
		Name      string `json:"groupName"`
		UpdatedAt time.Time
		DeletedAt *time.Time
		Data      any
		MoreData  interface{}
	}

	type TypeNotGivenToTheRegistry string

	type GroupMap map[string]Group

	type User struct {
		Reports        map[UserID]bool
		UserID         UserID `json:"userID"`
		PrimaryGroup   Group  `json:"primaryGroup"`
		UnknownType    TypeNotGivenToTheRegistry
		SecondaryGroup *Group   `json:"secondaryGroup,omitempty"`
		Tags           []string `json:"user_tags"`
		Private        any      `json:"-"`
		unexported     any
	}

	type BaseResponse[T any] struct {
		UpdatedAt time.Time `json:"updated_at"`
		GroupMap  GroupMap  `json:"group_map"`
		Data      T         `json:"data"`
	}

	_ = User{}.unexported

	type UsersResponse BaseResponse[[]User]

	service := typescript.New(
		map[string]any{
			"TestUserID":    UserID(0),
			"GroupResponse": BaseResponse[Group]{},
			"UserResponse":  UsersResponse{},
			"group":         Group{},
			"SystemUser":    User{},
			"GroupMapA":     GroupMap{},
			"GroupMapB":     map[string]Group{},
		},
	)

	testThePackage(t, service)
}

func testThePackage(t *testing.T, service *typescript.Service) {
	actualFileName := "test_files/" + t.Name() + "_actual.ts"
	actualFile, err := os.Create(actualFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer actualFile.Close()

	actualFileBuffer := bytes.NewBuffer([]byte{})

	writer := io.MultiWriter(actualFile, actualFileBuffer)

	if err := service.Generate(writer); err != nil {
		t.Fatal(err)
	}

	expectedContents, err := os.ReadFile("test_files/" + t.Name() + "_expected.ts")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actualFileBuffer.Bytes(), expectedContents) {
		wd, _ := os.Getwd()
		t.Fatal("contents don't match: " + wd + "/" + actualFileName)
	}
}
