package libs

import "testing"

func TestToken(t *testing.T) {
	uid, deviceId, scope := "a", "b", "c"
	token, err := SuperJWTSign(uid, deviceId, scope)
	ok, claim := SuperJWTVerify(token)
	
	if err!=nil || !ok {
		t.Error("SuperJWTVerify Error")
	}

	if uid != claim.UserId || deviceId != claim.DeviceId || scope != claim.Scope {
		t.Error("SuperJWTVerify Error")
	}

}
