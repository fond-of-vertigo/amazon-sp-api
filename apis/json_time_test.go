package apis

import (
	"reflect"
	"testing"
	"time"
)

func TestJsonTimeISO8601_MarshalJSON(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				Time: time.Date(2022, 02, 26, 9, 12, 11, 0, time.UTC),
			},
			want:    []byte("\"2022-02-26T09:12:11Z\""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := JsonTimeISO8601{
				Time: tt.fields.Time,
			}
			got, err := t.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t1.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
