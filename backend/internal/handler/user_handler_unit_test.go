package handler

import (
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/service"
)

func TestMaskIDCard(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{name: "empty", input: "", output: ""},
		{name: "short_2", input: "12", output: "**"},
		{name: "short_8", input: "12345678", output: "1******8"},
		{name: "normal_18", input: "110101199001011234", output: "1101**********1234"},
	}

	for _, tc := range testCases {
		got := maskIDCard(tc.input)
		if got != tc.output {
			t.Fatalf("%s: expected %s, got %s", tc.name, tc.output, got)
		}
	}
}

func TestCalculateAge(t *testing.T) {
	now := time.Date(2026, 2, 10, 12, 0, 0, 0, time.UTC)

	if got := calculateAge(time.Date(1990, 2, 10, 0, 0, 0, 0, time.UTC), now); got != 36 {
		t.Fatalf("expected 36, got %d", got)
	}
	if got := calculateAge(time.Date(1990, 2, 11, 0, 0, 0, 0, time.UTC), now); got != 35 {
		t.Fatalf("expected 35 before birthday, got %d", got)
	}
	if got := calculateAge(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), now); got != 0 {
		t.Fatalf("expected 0 for future birth date, got %d", got)
	}
}

func TestIDCardTokenGenerateAndVerify(t *testing.T) {
	h := NewUserHandler(nil, "unit-test-secret")

	idCardHash := h.idCardDigest("110101199001011234")
	token, err := h.generateIDCardToken(1001, idCardHash)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	if !h.verifyIDCardToken(1001, idCardHash, token) {
		t.Fatal("expected token verification success")
	}
	if h.verifyIDCardToken(1002, idCardHash, token) {
		t.Fatal("expected token verification fail for different user")
	}
	if h.verifyIDCardToken(1001, "another-hash", token) {
		t.Fatal("expected token verification fail for different hash")
	}

	tampered := token[:len(token)-1] + "x"
	if h.verifyIDCardToken(1001, idCardHash, tampered) {
		t.Fatal("expected token verification fail for tampered token")
	}
}

func TestToUserResponseIncludesDerivedFields(t *testing.T) {
	h := NewUserHandler(nil, "unit-test-secret")
	birthDate := time.Date(1990, 2, 10, 0, 0, 0, 0, time.UTC)

	user := &service.UserWithIdentities{
		User: &model.User{
			ID:        1,
			Phone:     "13800000001",
			Name:      "张三",
			Email:     "zhangsan@example.com",
			Avatar:    "https://example.com/avatar.jpg",
			Gender:    "male",
			BirthDate: birthDate,
			IDCard:    "110101199002101234",
			StationID: 1,
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PrimaryIdentity: &model.UserIdentity{IdentityType: "staff"},
		BEndIdentities:  []string{"staff"},
		CEndIdentities:  []string{},
	}

	resp := h.toUserResponse(user)

	if resp["id_card_masked"] != "1101**********1234" {
		t.Fatalf("expected masked id card, got %v", resp["id_card_masked"])
	}
	idCardHash, ok := resp["id_card_hash"].(string)
	if !ok || idCardHash == "" {
		t.Fatalf("expected non-empty id_card_hash, got %v", resp["id_card_hash"])
	}
	token, ok := resp["id_card_token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected non-empty id_card_token, got %v", resp["id_card_token"])
	}
	if !h.verifyIDCardToken(1, idCardHash, token) {
		t.Fatal("expected id_card_token to be verifiable")
	}

	if resp["birth_date"] != "1990-02-10" {
		t.Fatalf("expected birth_date 1990-02-10, got %v", resp["birth_date"])
	}
	if _, ok := resp["age"].(*int); !ok {
		t.Fatalf("expected age as *int, got %T", resp["age"])
	}
}
