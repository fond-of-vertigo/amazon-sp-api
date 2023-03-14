package apis

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/fond-of-vertigo/amazon-sp-api/constants"
)

type dummyHTTPClient struct {
	endpoint constants.Endpoint
	req      *http.Request
	resp     *http.Response
	errResp  error
}
type dummyBody struct {
	Message string
	Number  float64
}

func (r *dummyHTTPClient) Do(req *http.Request) (*http.Response, error) {
	r.req = req
	r.resp.Request = req
	return r.resp, r.errResp
}
func (r *dummyHTTPClient) GetEndpoint() constants.Endpoint {
	return r.endpoint
}
func (r *dummyHTTPClient) Close() {
}

func Test_call_Execute(t *testing.T) {
	type args struct {
		method              string
		endpoint            constants.Endpoint
		restrictedDataToken string
		url                 string
		queryParams         url.Values
		reqBodyObject       any
	}
	type want struct {
		url  string
		resp CallResponse[dummyBody]
	}
	type testCase struct {
		name    string
		args    args
		want    want
		wantErr bool
	}
	tests := []testCase{
		{
			name: "Simple get",
			args: args{
				endpoint: constants.NorthAmerica,
				url:      "/message",
				method:   http.MethodGet,
			},
			want: want{
				url: "https://sellingpartnerapi-na.amazon.com/message",
				resp: CallResponse[dummyBody]{
					ResponseBody: &dummyBody{
						Message: "All ok",
						Number:  4711.0815,
					},
				},
			},
		},
		{
			name: "Get with restricted data token",
			args: args{
				endpoint:            constants.NorthAmerica,
				url:                 "/message",
				method:              http.MethodGet,
				restrictedDataToken: "ABCDED",
			},
			want: want{
				url: "https://sellingpartnerapi-na.amazon.com/message",
				resp: CallResponse[dummyBody]{
					ResponseBody: &dummyBody{
						Message: "All ok",
						Number:  4711.0815,
					},
				},
			},
		},
		{
			name: "Post body with queryParam",
			args: args{
				endpoint: constants.NorthAmerica,
				url:      "/message",
				method:   http.MethodPost,
				queryParams: map[string][]string{
					"final":     {"true"},
					"messageID": {"1234"},
				},
				reqBodyObject: &dummyBody{
					Message: "Hello there...",
					Number:  47.0,
				},
			},
			want: want{
				url: "https://sellingpartnerapi-na.amazon.com/message?final=true&messageID=1234",
				resp: CallResponse[dummyBody]{
					ResponseBody: &dummyBody{
						Message: "All ok",
						Number:  4711.0815,
					},
				},
			},
		},
		{
			name: "Delete without body",
			args: args{
				endpoint: constants.NorthAmerica,
				url:      "/message/4711",
				method:   http.MethodDelete,
			},
			want: want{
				url: "https://sellingpartnerapi-na.amazon.com/message/4711",
			},
		},
		{
			name: "Error case",
			args: args{
				endpoint: constants.NorthAmerica,
				url:      "/message/4711",
				method:   http.MethodDelete,
			},
			want: want{
				url: "https://sellingpartnerapi-na.amazon.com/message/4711",
				resp: CallResponse[dummyBody]{
					ErrorList: &ErrorList{
						Errors: []Error{
							{
								Message: "Oooops",
								Code:    "4711",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given:
			reqBodyBytes, err := getJSONBytes(tt.args.reqBodyObject)
			if err != nil {
				t.Fatal(err)
			}
			mockResp, err := mockResponse(&tt.want.resp)
			if err != nil {
				t.Fatal(err)
			}
			client := &dummyHTTPClient{
				endpoint: tt.args.endpoint,
				resp:     mockResp,
			}

			// when:
			call := NewCall[dummyBody](tt.args.method, tt.args.url).
				WithQueryParams(tt.args.queryParams).
				WithBody(reqBodyBytes)

			if tt.want.resp.ErrorList != nil {
				call = call.WithParseErrorListOnError()
			}

			if tt.args.restrictedDataToken != "" {
				call = call.WithRestrictedDataToken(&tt.args.restrictedDataToken)
			}

			got, err := call.Execute(client)

			// then:
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() unexpected error = '%v'", err)
				return
			}
			if client.req.URL.String() != tt.want.url {
				t.Errorf("Execute(): url different. got = '%v', want '%v'", client.req.URL.String(), tt.want.url)
				return
			}
			if client.req.Method != tt.args.method {
				t.Errorf("Execute(): method different. got = '%v', want '%v'", client.req.Method, tt.want.url)
				return
			}
			if tt.args.restrictedDataToken != client.req.Header.Get(constants.AccessTokenHeader) {
				t.Errorf("Execute(): method AccessTokenHeader. got = '%v', want '%v'", client.req.Header.Get(constants.AccessTokenHeader), tt.args.restrictedDataToken)
				return
			}

			gotReqBodyBytes := make([]byte, len(reqBodyBytes))

			if len(gotReqBodyBytes) > 0 {
				if _, err = client.req.Body.Read(gotReqBodyBytes); err != nil {
					t.Fatal(err)
				}
			}
			if !bytes.Equal(gotReqBodyBytes, reqBodyBytes) {
				t.Errorf("Execute(): request body different.")
				return
			}

			if tt.want.resp.ResponseBody == nil && !reflect.ValueOf(got.ResponseBody).IsNil() {
				t.Errorf("Execute(): response different. got = '%v', want nil", got)
			} else {
				if tt.wantErr {
					if tt.want.resp.ResponseBody != nil && !reflect.DeepEqual(got.ResponseBody, tt.want.resp.ResponseBody) {
						t.Errorf("Execute(): error response different. got = '%v', want '%v'", got.ResponseBody, tt.want.resp.ResponseBody)
					}
				} else {
					if tt.want.resp.ResponseBody != nil && !reflect.DeepEqual(got.ResponseBody, tt.want.resp.ResponseBody) {
						t.Errorf("Execute(): response different. got = '%v', want '%v'", got.ResponseBody, tt.want.resp.ResponseBody)
					}
				}
			}
			if diff(tt.want.resp.ResponseBody, got.ResponseBody) {
				t.Errorf("Execute(): responseBody different. want = '%v', got '%v'", tt.want.resp.ErrorList, got.ErrorList)
			}
			if diff(tt.want.resp.ErrorList, got.ErrorList) {
				t.Errorf("Execute(): errorList different. want = '%v', got '%v'", tt.want.resp.ErrorList, got.ErrorList)
			}
		})
	}
}

func diff(want any, got any) bool {
	if want == nil && !reflect.ValueOf(want).IsNil() {
		return true
	} else {
		if want != nil && !reflect.DeepEqual(got, want) {
			return true
		}
	}
	return false
}

func mockResponse(callResp *CallResponse[dummyBody]) (*http.Response, error) {
	if callResp.ErrorList != nil {
		bodyBytes, err := getJSONBytes(callResp.ErrorList)
		if err != nil {
			return nil, err
		}
		return &http.Response{
			Status:        "500 Internal Server Error",
			StatusCode:    http.StatusInternalServerError,
			Body:          io.NopCloser(bytes.NewReader(bodyBytes)),
			ContentLength: int64(len(bodyBytes)),
		}, nil
	}
	bodyBytes, err := getJSONBytes(callResp.ResponseBody)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewReader(bodyBytes)),
		ContentLength: int64(len(bodyBytes)),
	}, nil
}
func getJSONBytes(obj any) ([]byte, error) {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		return []byte{}, nil
	}
	return json.Marshal(obj)
}

func Test_calcWaitTimeByRateLimit(t *testing.T) {
	type args struct {
		callsPer float32
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "0.5 req per sec, wait 2 seconds",
			args: args{0.5, time.Second},
			want: 2 * time.Second,
		}, {
			name: "2 req per sec, wait 500ms",
			args: args{2, time.Second},
			want: 500 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcWaitTimeByRateLimit(tt.args.callsPer, tt.args.duration); got != tt.want {
				t.Errorf("calcWaitTimeByRateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
