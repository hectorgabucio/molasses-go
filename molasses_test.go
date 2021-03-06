package molasses_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/molassesapp/molasses-go"
	"github.com/stretchr/testify/assert"
)

func TestInitWithValidFeatureAndStop(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/get-features", req.URL.String())

		rw.Write([]byte(`{"data":{"name":"Production","updatedAt":"2020-08-26T02:11:44Z","features":[{"id":"f603f621-83ba-46f0-adf5-70ed2d668646","key":"GOOGLE_SSO","description":"asdfasdf","active":true,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]},{"id":"f3fae17d-a8d2-446f-8e85-bfa408562b73","key":"MOBILE_CHECKOUT","description":"asdfasd","active":false,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]},{"id":"d1f276ce-80b7-45a7-a70d-dad190abcd6e","key":"NEW_CHECKOUT","description":"this is the new checkout screen","active":false,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]}]}}`))
	}))

	client, err := molasses.Init(molasses.ClientOptions{
		HTTPClient: server.Client(),
		APIKey:     "API_KEY",
		URL:        server.URL,
		SendEvents: molasses.Bool(false),
	})
	assert.True(t, client.IsInitiated())
	if err != nil {
		t.Error(err)
	}
	assert.True(t, client.IsActive("GOOGLE_SSO"))
	assert.False(t, client.IsActive("MOBILE_CHECKOUT", molasses.User{ID: "USERID1"}))
	client.Stop()
	assert.False(t, client.IsInitiated())
}

func TestDefaultsAreSet(t *testing.T) {
	client, err := molasses.Init(molasses.ClientOptions{
		APIKey:     "API_KEY",
		SendEvents: molasses.Bool(false),
	})
	assert.True(t, client.IsInitiated())
	if err != nil {
		t.Error(err)
	}
}

func TestErrorsWhenAPIKeyIsNotSet(t *testing.T) {
	_, err := molasses.Init(molasses.ClientOptions{
		APIKey:     "",
		SendEvents: molasses.Bool(false),
	})
	if err != nil {
		assert.Error(t, err)
	}
}

func TestInitWithValidFeatureWithUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/get-features" {

			assert.Equal(t, "/get-features", req.URL.String())
			rw.Write([]byte(`{"data":{"name":"Production","updatedAt":"2020-08-26T02:11:44Z","features":[{"id":"f603f621-83ba-46f0-adf5-70ed2d668646","key":"GOOGLE_SSO","description":"asdfasdf","active":true,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]},{"id":"f3fae17d-a8d2-446f-8e85-bfa408562b73","key":"MOBILE_CHECKOUT","description":"asdfasd","active":false,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]},{"id":"d1f276ce-80b7-45a7-a70d-dad190abcd6e","key":"NEW_CHECKOUT","description":"this is the new checkout screen","active":false,"segments":[{"segmentType":"everyoneElse","userConstraints":[{"operator":"all","values":"","userParam":"","userParamType":""}],"percentage":100}]}]}}`))
			return
		}
		assert.Equal(t, "/analytics", req.URL.String())
		rw.Write([]byte(`{}`))
		return
	}))

	client, err := molasses.Init(molasses.ClientOptions{
		HTTPClient: server.Client(),
		APIKey:     "API_KEY",
		URL:        server.URL,
	})
	assert.True(t, client.IsInitiated())
	if err != nil {
		t.Error(err)
	}
	assert.True(t, client.IsActive("GOOGLE_SSO", molasses.User{
		ID: "1234",
		Params: map[string]string{
			"foo": "bar",
		},
	}))
	assert.False(t, client.IsActive("MOBILE_CHECKOUT", molasses.User{
		ID: "1234",
		Params: map[string]string{
			"foo": "bar",
		},
	}))
	client.Stop()
	time.Sleep(1 * time.Second)
	assert.False(t, client.IsInitiated())
}

func TestOtherSegments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/get-features" {

			assert.Equal(t, "/get-features", req.URL.String())
			rw.Write([]byte(`{
				"data": {
					"name": "Production",
					"updatedAt": "2020-08-26T02:11:44Z",
					"features": [
						{
							"id": "f603f621-83ba-46f0-adf5-70ed2d668646",
							"key": "GOOGLE_SSO",
							"description": "asdfasdf",
							"active": true,
							"segments": [
								{
									"segmentType": "alwaysControl",
									"userConstraints": [
										{
											"operator": "equals",
											"values": "true",
											"userParam": "controlUser",
											"userParamType": ""
										},
										{
											"operator": "nin",
											"values": "yes,maybe,definitely",
											"userParam": "experimentUser",
											"userParamType": ""
										}
									],
									"percentage": 100
								},
								{
									"segmentType": "alwaysExperiment",
									"userConstraints": [
										{
											"operator": "contains",
											"values": "fals",
											"userParam": "controlUser",
											"userParamType": ""
										},
										{
											"operator": "in",
											"values": "yes,maybe,definitely",
											"userParam": "experimentUser",
											"userParamType": ""
										}
									],
									"percentage": 100
								},
								{
									"segmentType": "everyoneElse",
									"userConstraints": [
										{
											"operator": "all",
											"values": "",
											"userParam": "",
											"userParamType": ""
										}
									],
									"percentage": 50
								}
							]
						}
					]
				}
			}`))
			return
		}
		assert.Equal(t, "/analytics", req.URL.String())
		rw.Write([]byte(`{}`))
		return
	}))

	client, err := molasses.Init(molasses.ClientOptions{
		HTTPClient: server.Client(),
		APIKey:     "API_KEY",
		URL:        server.URL,
	})
	assert.True(t, client.IsInitiated())
	if err != nil {
		t.Error(err)
	}
	experimentUser := molasses.User{
		ID: "1235",
		Params: map[string]string{
			"controlUser":    "false",
			"experimentUser": "yes",
		},
	}
	controlUser := molasses.User{
		ID: "1234",
		Params: map[string]string{
			"controlUser":    "true",
			"experimentUser": "nope",
		},
	}
	assert.False(t, client.IsActive("GOOGLE_SSO", controlUser))
	client.ExperimentSuccess("GOOGLE_SSO", controlUser, map[string]string{})
	assert.True(t, client.IsActive("GOOGLE_SSO", experimentUser))
	client.ExperimentSuccess("GOOGLE_SSO", experimentUser, map[string]string{})
	assert.False(t, client.IsActive("GOOGLE_SSO", molasses.User{
		ID: "1",
		Params: map[string]string{
			"controlUser": "bar",
		},
	}))

	assert.True(t, client.IsActive("GOOGLE_SSO", molasses.User{
		ID: "2",
		Params: map[string]string{
			"controlUser": "bar",
		},
	}))
	client.Stop()
	time.Sleep(1 * time.Second)
	assert.False(t, client.IsInitiated())
}

func TestMoreSegments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/get-features" {

			assert.Equal(t, "/get-features", req.URL.String())
			rw.Write([]byte(`{
				"data": {
					"name": "Production",
					"updatedAt": "2020-08-26T02:11:44Z",
					"features": [
						{
							"id": "f603f621-83ba-46f0-adf5-70ed2d668646",
							"key": "GOOGLE_SSO",
							"description": "asdfasdf",
							"active": true,
							"segments": [
								{
									"segmentType": "alwaysControl",
									"userConstraints": [
										{
											"operator": "doesNotEqual",
											"values": "false",
											"userParam": "controlUser",
											"userParamType": ""
										},
										{
											"operator": "doesNotContain",
											"values": "yes",
											"userParam": "experimentUser",
											"userParamType": ""
										}
									],
									"percentage": 100
								},
								{
									"segmentType": "alwaysExperiment",
									"constraint":  "any",
									"userConstraints": [
										{
											"operator": "contains",
											"values": "fals",
											"userParam": "controlUser",
											"userParamType": ""
										},
										{
											"operator": "in",
											"values": "yes,maybe,definitely",
											"userParam": "experimentUser",
											"userParamType": ""
										},
										{
											"operator": "in",
											"values": "1235,123,1",
											"userParam": "id",
											"userParamType": ""
										}
									],
									"percentage": 100
								},
								{
									"segmentType": "everyoneElse",
									"userConstraints": [
										{
											"operator": "all",
											"values": "",
											"userParam": "",
											"userParamType": ""
										}
									],
									"percentage": 50
								}
							]
						}
					]
				}
			}`))
			return
		}
		assert.Equal(t, "/analytics", req.URL.String())
		rw.Write([]byte(`{}`))
		return
	}))

	client, err := molasses.Init(molasses.ClientOptions{
		HTTPClient: server.Client(),
		APIKey:     "API_KEY",
		URL:        server.URL,
	})
	assert.True(t, client.IsInitiated())
	if err != nil {
		t.Error(err)
	}
	experimentUser := molasses.User{
		ID: "1235",
		Params: map[string]string{
			"experimentUser": "yes",
		},
	}
	controlUser := molasses.User{
		ID: "1234",
		Params: map[string]string{
			"controlUser":    "true",
			"experimentUser": "nope",
		},
	}
	assert.False(t, client.IsActive("GOOGLE_SSO", controlUser))
	client.ExperimentSuccess("GOOGLE_SSO", controlUser, map[string]string{})
	assert.True(t, client.IsActive("GOOGLE_SSO", experimentUser))
	client.ExperimentSuccess("GOOGLE_SSO", experimentUser, map[string]string{
		"experiment_id": "hello",
		"button_color":  "green",
	})
	assert.False(t, client.IsActive("GOOGLE_SSO", molasses.User{
		ID: "5",
		Params: map[string]string{
			"controlUser": "bar",
		},
	}))

	assert.True(t, client.IsActive("GOOGLE_SSO", molasses.User{
		ID: "2",
		Params: map[string]string{
			"controlUser": "bar",
		},
	}))
	client.Stop()
	time.Sleep(1 * time.Second)
	assert.False(t, client.IsInitiated())
}
