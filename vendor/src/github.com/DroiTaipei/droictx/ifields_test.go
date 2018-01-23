package droictx

import "testing"

// Mock Peeker
type mockPeeker func(key string) []byte

func (p mockPeeker) Peek(key string) []byte { return p(key) }

//Mock Getter
type mockGetter func(key string) string

func (g mockGetter) Get(key string) string { return g(key) }

type mockSetter map[string]string

func (s mockSetter) Set(key, value string) { s[key] = value }

//Mock GetHeaderer
type mockGetHer func(key string) string

func (g mockGetHer) GetHeader(key string) string { return g(key) }

func TestHeaderMap(t *testing.T) {
	c := &DoneContext{}
	c.Set("Aid", "ZXC123ASDQWE")
	c.Set("Rid", "1029384756")
	c.Set("Ak", "abcdefg123")
	c.Set("SAid", "4RRFFV3edc")
	c.Set("R", "16384")
	c.Set("St", "lkjnhbhjw-sdahu")
	c.Set("XUri", "/api/v2/xxx/yyy/zzz")
	c.Set("XMd", "PATCH")
	c.Set("XIp", "192.168.3.4")
	c.Set("XPort", "26841")
	c.Set("OpMode", "soft")

	m := c.HeaderMap()
	v, ok := m["X-Droi-AppID"]
	if !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	v, ok = m["X-Droi-ReqID"]
	if !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	v, ok = m["X-Droi-Api-Key"]
	if !ok || v != "abcdefg123" {
		t.Error("ApiKey not match Ak")
	}
	v, ok = m["X-Droi-Service-AppID"]
	if !ok || v != "4RRFFV3edc" {
		t.Error("Serivce AppID not match SAid")
	}
	v, ok = m["X-Droi-Role"]
	if !ok || v != "16384" {
		t.Error("Role not match R")
	}
	v, ok = m["X-Droi-Session-Token"]
	if !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Session Token not match St")
	}
	v, ok = m["X-Droi-URI"]
	if !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	v, ok = m["X-Droi-Method"]
	if !ok || v != "PATCH" {
		t.Error("Method not match XMd")
	}
	v, ok = m["X-Droi-Remote-IP"]
	if !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	v, ok = m["X-Droi-Remote-Port"]
	if !ok || v != "26841" {
		t.Error("Remote Port not match Xport")
	}
	v, ok = m["X-Droi-Op-Mode"]
	if !ok || v != "soft" {
		t.Error("Op Mode not match OpMode")
	}
}

func TestHeaderSet(t *testing.T) {
	c := &DoneContext{}
	c.HeaderSet("X-Droi-AppID", "ZXC123ASDQWE")
	c.HeaderSet("X-Droi-ReqID", "1029384756")
	c.HeaderSet("X-Droi-Api-Key", "abcdefg123")
	c.HeaderSet("X-Droi-Service-AppID", "4RRFFV3edc")
	c.HeaderSet("X-Droi-Role", "16384")
	c.HeaderSet("X-Droi-Session-Token", "lkjnhbhjw-sdahu")
	c.HeaderSet("X-Droi-URI", "/api/v2/xxx/yyy/zzz")
	c.HeaderSet("X-Droi-Method", "PATCH")
	c.HeaderSet("X-Droi-Remote-IP", "192.168.3.4")
	c.HeaderSet("X-Droi-Remote-Port", "26841")
	c.HeaderSet("X-Droi-Op-Mode", "soft")

	if v, ok := c.GetString("Aid"); !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	if v, ok := c.GetString("Rid"); !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	if v, ok := c.GetString("Ak"); !ok || v != "abcdefg123" {
		t.Error("APIKey not match Ak")
	}
	if v, ok := c.GetString("SAid"); !ok || v != "4RRFFV3edc" {
		t.Error("Service Aid not match SAid")
	}
	if v, ok := c.GetString("R"); !ok || v != "16384" {
		t.Error("Role not match R")
	}
	if v, ok := c.GetString("St"); !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Session Token not match St")
	}
	if v, ok := c.GetString("XUri"); !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	if v, ok := c.GetString("XMd"); !ok || v != "PATCH" {
		t.Error("Method Aid not match XMd")
	}
	if v, ok := c.GetString("XIp"); !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	if v, ok := c.GetString("XPort"); !ok || v != "26841" {
		t.Error("Remote Port not match XPort")
	}
	if v, ok := c.GetString("OpMode"); !ok || v != "soft" {
		t.Error("Op Mode not match OpMode")
	}
}

func TestPeeker(t *testing.T) {

	peeker := mockPeeker(func(key string) []byte {
		if key == "X-Droi-AppID" {
			return []byte("ZXC123ASDQWE")
		}
		if key == "X-Droi-ReqID" {
			return []byte("1029384756")
		}
		if key == "X-Droi-Api-Key" {
			return []byte("abcdefg123")
		}
		if key == "X-Droi-Service-AppID" {
			return []byte("4RRFFV3edc")
		}
		if key == "X-Droi-Role" {
			return []byte("16384")
		}
		if key == "X-Droi-Session-Token" {
			return []byte("lkjnhbhjw-sdahu")
		}
		if key == "X-Droi-URI" {
			return []byte("/api/v2/xxx/yyy/zzz")
		}
		if key == "X-Droi-Method" {
			return []byte("PATCH")
		}
		if key == "X-Droi-Remote-IP" {
			return []byte("192.168.3.4")
		}
		if key == "X-Droi-Remote-Port" {
			return []byte("26841")
		}
		return []byte{}
	})

	c := GetContextFromPeeker(peeker)
	if v, ok := c.GetString("Aid"); !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	if v, ok := c.GetString("Rid"); !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	if v, ok := c.GetString("Ak"); !ok || v != "abcdefg123" {
		t.Error("APIKey not match Ak")
	}
	if v, ok := c.GetString("SAid"); !ok || v != "4RRFFV3edc" {
		t.Error("Service AppID not match SAid")
	}
	if v, ok := c.GetString("R"); !ok || v != "16384" {
		t.Error("Role not match R")
	}
	if v, ok := c.GetString("St"); !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Session Token not match St")
	}
	if v, ok := c.GetString("XUri"); !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	if v, ok := c.GetString("XMd"); !ok || v != "PATCH" {
		t.Error("Method Aid not match XMd")
	}
	if v, ok := c.GetString("XIp"); !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	if v, ok := c.GetString("XPort"); !ok || v != "26841" {
		t.Error("Remote Port not match XPort")
	}
	// Test Empty Field
	m := c.Map()
	if _, ok := m["Aidm"]; ok {
		t.Error("Aidm should not exists")
	}
}

func TestGetter(t *testing.T) {

	getter := mockGetter(func(key string) string {
		if key == "X-Droi-AppID" {
			return "ZXC123ASDQWE"
		}
		if key == "X-Droi-ReqID" {
			return "1029384756"
		}
		if key == "X-Droi-Api-Key" {
			return "abcdefg123"
		}
		if key == "X-Droi-Service-AppID" {
			return "4RRFFV3edc"
		}
		if key == "X-Droi-Role" {
			return "16384"
		}
		if key == "X-Droi-Session-Token" {
			return "lkjnhbhjw-sdahu"
		}
		if key == "X-Droi-URI" {
			return "/api/v2/xxx/yyy/zzz"
		}
		if key == "X-Droi-Method" {
			return "PATCH"
		}
		if key == "X-Droi-Remote-IP" {
			return "192.168.3.4"
		}
		if key == "X-Droi-Remote-Port" {
			return "26841"
		}
		return ""
	})

	c := GetContextFromGetter(getter)
	if v, ok := c.GetString("Aid"); !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	if v, ok := c.GetString("Rid"); !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	if v, ok := c.GetString("Ak"); !ok || v != "abcdefg123" {
		t.Error("APIKey not match Ak")
	}
	if v, ok := c.GetString("SAid"); !ok || v != "4RRFFV3edc" {
		t.Error("Service AppID not match SAid")
	}
	if v, ok := c.GetString("R"); !ok || v != "16384" {
		t.Error("Role not match R")
	}
	if v, ok := c.GetString("St"); !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Session Token not match St")
	}
	if v, ok := c.GetString("XUri"); !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	if v, ok := c.GetString("XMd"); !ok || v != "PATCH" {
		t.Error("Method Aid not match XMd")
	}
	if v, ok := c.GetString("XIp"); !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	if v, ok := c.GetString("XPort"); !ok || v != "26841" {
		t.Error("Remote Port not match XPort")
	}
	// Test Empty Field
	m := c.Map()
	if _, ok := m["Aidm"]; ok {
		t.Error("Aidm should not exists")
	}
}

func TestSetter(t *testing.T) {

	var s = make(mockSetter)
	c := &DoneContext{}
	c.HeaderSet("X-Droi-AppID", "ZXC123ASDQWE")
	c.HeaderSet("X-Droi-ReqID", "1029384756")
	c.HeaderSet("X-Droi-Api-Key", "abcdefg123")
	c.HeaderSet("X-Droi-Service-AppID", "4RRFFV3edc")
	c.HeaderSet("X-Droi-Role", "16384")
	c.HeaderSet("X-Droi-Session-Token", "lkjnhbhjw-sdahu")
	c.HeaderSet("X-Droi-URI", "/api/v2/xxx/yyy/zzz")
	c.HeaderSet("X-Droi-Method", "PATCH")
	c.HeaderSet("X-Droi-Remote-IP", "192.168.3.4")
	c.HeaderSet("X-Droi-Remote-Port", "26841")

	c.SetHTTPHeaders(s)
	if v, ok := s["X-Droi-AppID"]; !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	if v, ok := s["X-Droi-ReqID"]; !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	if v, ok := s["X-Droi-Api-Key"]; !ok || v != "abcdefg123" {
		t.Error("APIKey not match Ak")
	}
	if v, ok := s["X-Droi-Service-AppID"]; !ok || v != "4RRFFV3edc" {
		t.Error("Service AppID not match SAid")
	}
	if v, ok := s["X-Droi-Role"]; !ok || v != "16384" {
		t.Error("Role not match R")
	}
	if v, ok := s["X-Droi-Session-Token"]; !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Role not match R")
	}
	if v, ok := s["X-Droi-URI"]; !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	if v, ok := s["X-Droi-Method"]; !ok || v != "PATCH" {
		t.Error("Method not match XMd")
	}
	if v, ok := s["X-Droi-Remote-IP"]; !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	if v, ok := s["X-Droi-Remote-Port"]; !ok || v != "26841" {
		t.Error("Remote Port not match Xport")
	}
	// Test Empty Field
	if _, ok := s["X-Droi-AidMode"]; ok {
		t.Error("Aidm should not exists")
	}
}

func TestGeHeaderer(t *testing.T) {

	getter := mockGetHer(func(key string) string {
		if key == "X-Droi-AppID" {
			return "ZXC123ASDQWE"
		}
		if key == "X-Droi-ReqID" {
			return "1029384756"
		}
		if key == "X-Droi-Api-Key" {
			return "abcdefg123"
		}
		if key == "X-Droi-Service-AppID" {
			return "4RRFFV3edc"
		}
		if key == "X-Droi-Role" {
			return "16384"
		}
		if key == "X-Droi-Session-Token" {
			return "lkjnhbhjw-sdahu"
		}
		if key == "X-Droi-URI" {
			return "/api/v2/xxx/yyy/zzz"
		}
		if key == "X-Droi-Method" {
			return "PATCH"
		}
		if key == "X-Droi-Remote-IP" {
			return "192.168.3.4"
		}
		if key == "X-Droi-Remote-Port" {
			return "26841"
		}
		return ""
	})

	c := GetContextFromGetHeader(getter)
	if v, ok := c.GetString("Aid"); !ok || v != "ZXC123ASDQWE" {
		t.Error("AppID not match Aid")
	}
	if v, ok := c.GetString("Rid"); !ok || v != "1029384756" {
		t.Error("ReqID not match Rid")
	}
	if v, ok := c.GetString("Ak"); !ok || v != "abcdefg123" {
		t.Error("APIKey not match Ak")
	}
	if v, ok := c.GetString("SAid"); !ok || v != "4RRFFV3edc" {
		t.Error("Service AppID not match SAid")
	}
	if v, ok := c.GetString("R"); !ok || v != "16384" {
		t.Error("Role not match R")
	}
	if v, ok := c.GetString("St"); !ok || v != "lkjnhbhjw-sdahu" {
		t.Error("Session Token not match St")
	}
	if v, ok := c.GetString("XUri"); !ok || v != "/api/v2/xxx/yyy/zzz" {
		t.Error("URI not match XUri")
	}
	if v, ok := c.GetString("XMd"); !ok || v != "PATCH" {
		t.Error("Method Aid not match XMd")
	}
	if v, ok := c.GetString("XIp"); !ok || v != "192.168.3.4" {
		t.Error("Remote IP not match XIp")
	}
	if v, ok := c.GetString("XPort"); !ok || v != "26841" {
		t.Error("Remote Port not match XPort")
	}
	// Test Empty Field
	m := c.Map()
	if _, ok := m["Aidm"]; ok {
		t.Error("Aidm should not exists")
	}
}
