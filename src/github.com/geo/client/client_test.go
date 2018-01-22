package client

import (
	"errors"
	geoMock "github.com/DroiTaipei/droipkg/geo/mock"
	geopb "github.com/DroiTaipei/droipkg/geo/protobuf"
	"github.com/DroiTaipei/droipkg/grpc/connpool"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"os"
	"testing"
)

type mockSimplifiedPool func(...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error)

func (msp mockSimplifiedPool) Get(addr ...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error) {
	return msp(addr...)
}

type matcher func(interface{}) bool

func (m matcher) Matches(msg interface{}) bool {
	return m(msg)
}

func (m matcher) String() string {
	return "Matcher"
}

type mockPools func(...string) (*connpool.Pool, error)

func (mp mockPools) Get(addr string) (*connpool.Pool, error) {
	return mp(addr)
}

func (mp mockPools) GetRoundRobin() (*connpool.Pool, error) {
	return mp()
}

func (mp mockPools) Connect(addr string, maxConn int) error {
	if addr == notExistAddr {
		return errors.New("Could not connect")
	}
	return nil
}

const (
	existAddr    = "Exist.Server:8888"
	notExistAddr = "Not.Exist.Server:0000"
)

var (
	p         mockPools
	sp        mockSimplifiedPool
	spErr     mockSimplifiedPool
	cf        ClientFactory
	cMatcher  matcher
	ipMatcher matcher
)

func BeforeTest() {
	p = mockPools(func(addrs ...string) (*connpool.Pool, error) {
		if len(addrs) < 1 {
			return nil, errors.New("no servers")
		}
		switch addrs[0] {
		case existAddr:
			return connpool.NewPool(addrs[0], 1)
		case notExistAddr:
			return nil, errors.New("addr not in map")
		}
		return &connpool.Pool{}, nil
	})
	sp = mockSimplifiedPool(func(addr ...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error) {
		return nil, func(in *grpc.ClientConn) error { return nil }, nil
	})
	spErr = mockSimplifiedPool(func(addr ...string) (*grpc.ClientConn, func(*grpc.ClientConn) error, error) {
		return nil, func(in *grpc.ClientConn) error { return nil }, errors.New("Get Conn Failed")
	})
	cMatcher = matcher(func(msg interface{}) bool {
		_, ok := msg.(*geopb.Content)
		return ok
	})
	ipMatcher = matcher(func(msg interface{}) bool {
		_, ok := msg.(*geopb.IP)
		return ok
	})
}

func TestPoolError(t *testing.T) {
	assert := assert.New(t)
	sp := &sPool{p}
	_, _, err := sp.Get()
	assert.Error(err)
	_, _, err = sp.Get(existAddr)
	assert.NoError(err)
	_, _, err = sp.Get(notExistAddr)
	assert.Error(err)
	_, _, err = sp.Get("XX")
	assert.Error(err)

}

func TestSimplifiedPoolError(t *testing.T) {
	assert := assert.New(t)
	sp := &sPool{}
	_, _, err := sp.Get("a")
	assert.Error(err)
	err = nil
	pools := connpool.NewPools()
	sp = &sPool{pools}
	_, _, err = sp.Get("a", "b", "c")
	assert.Error(err)
}

func TestInitialize(t *testing.T) {
	assert := assert.New(t)
	poolsFactory = func() connpool.Pools {
		return p
	}
	agent, err := Initialize(existAddr, 1)
	assert.NoError(err)
	assert.NotNil(agent)
	agent, err = Initialize(notExistAddr, 1)
	assert.Error(err)
	assert.Nil(agent)
}

func TestEcho(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.EXPECT().Echo(gomock.Any(), cMatcher).Return(&geopb.Content{Headers: map[string]string{"Hello": "grpc"}}, nil)
		return mockObj
	}
	a := NewAgent(sp, cf)
	in := map[string]string{"Hello": "grpc"}
	ret, err := a.Echo(in)
	assert.NoError(err)
	assert.Equal(in, ret)
}

func TestEchoError(t *testing.T) {
	assert := assert.New(t)
	a := NewAgent(sp, nil)
	ret, err := a.Echo(nil)
	assert.Error(err)
	assert.Nil(ret)
}

func TestGetMaxmindCity(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	info := &geopb.MaxmindCityInfo{
		Message:     "",
		City:        "济宁市",
		Subdivision: "山东省",
		Country:     "CN",
		Zone:        "Asia/Shanghai",
		Latitude:    35.405,
		Longitude:   116.5814,
	}
	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.
			EXPECT().
			GetMaxmindCity(gomock.Any(), ipMatcher).
			Return(info, nil)
		return mockObj
	}
	a := NewAgent(sp, cf)
	ret, err := a.GetMaxmindCity("60.211.182.76", "")
	assert.NoError(err)
	assert.EqualValues(info, ret)
}

func TestGetIpipCity(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	info := &geopb.IpipCityInfo{
		Message:   "",
		Country:   "中国",
		Province:  "广东",
		City:      "深圳",
		Org:       "",
		ISP:       "电信",
		Latitude:  "22.547",
		Longitude: "114.085947",
		TimeZone:  "Asia/Shanghai",
		UTC:       "UTC+8",
		ChinaCode: "440300",
		PhoneCode: "86",
		ISO2:      "CN",
		Continent: "AP",
	}
	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.
			EXPECT().
			GetIpipCity(gomock.Any(), ipMatcher).
			Return(info, nil)
		return mockObj
	}
	a := NewAgent(sp, cf)
	ret, err := a.GetIpipCity("60.211.181.76", "")
	assert.NoError(err)
	assert.EqualValues(info, ret)
}

func TestSpError(t *testing.T) {
	assert := assert.New(t)
	a := NewAgent(spErr, nil)
	ret, err := a.Echo(map[string]string{"Hello": "grpc"})
	assert.Error(err)
	assert.Nil(ret)
	info1, err := a.GetMaxmindCity("60.211.182.76", "")
	assert.Error(err)
	assert.Nil(info1)
	info2, err := a.GetIpipCity("60.211.181.76", "")
	assert.Error(err)
	assert.Nil(info2)
}

func TestGrpcError(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.
			EXPECT().
			Echo(gomock.Any(), cMatcher).
			Return(nil, errors.New("Unknown"))
		return mockObj
	}
	a := NewAgent(sp, cf)
	ret, err := a.Echo(map[string]string{"Hello": "grpc"})
	assert.Error(err)
	assert.Nil(ret)
	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.
			EXPECT().
			GetMaxmindCity(gomock.Any(), ipMatcher).
			Return(nil, errors.New("invalid ip"))
		return mockObj
	}
	a = NewAgent(sp, cf)
	info1, err := a.GetMaxmindCity("60.211.182.76", "")
	assert.Error(err)
	assert.Nil(info1)
	cf = func(in *grpc.ClientConn) geopb.GeoClient {
		mockObj := geoMock.NewMockGeoClient(mockCtrl)
		mockObj.
			EXPECT().
			GetIpipCity(gomock.Any(), ipMatcher).
			Return(nil, errors.New("Unknown"))

		return mockObj
	}
	a = NewAgent(sp, cf)
	info2, err := a.GetIpipCity("60.211.181.76", "")
	assert.Error(err)
	assert.Nil(info2)
}

// Do somethings after all test cases
func AfterTest() {

}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}
