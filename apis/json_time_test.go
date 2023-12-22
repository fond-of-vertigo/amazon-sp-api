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
			name: "Seconds in UTC",
			fields: fields{
				Time: time.Date(2022, 02, 26, 9, 12, 11, 0, time.UTC),
			},
			want:    []byte("\"2022-02-26T09:12:11Z\""),
			wantErr: false,
		},
		{
			name: "Zero Time",
			fields: fields{
				Time: time.Time{},
			},
			want:    []byte("\"0001-01-01T00:00:00Z\""),
			wantErr: false,
		},
		{
			name: "Time in PST",
			fields: fields{
				Time: time.Date(2022, 02, 26, 1, 12, 11, 0, time.FixedZone("PST", -8*3600)),
			},
			want:    []byte("\"2022-02-26T09:12:11Z\""), // Converted to UTC
			wantErr: false,
		},
		{
			name: "Time in IST",
			fields: fields{
				Time: time.Date(2022, 02, 26, 14, 42, 11, 0, time.FixedZone("IST", 5*3600+1800)),
			},
			want:    []byte("\"2022-02-26T09:12:11Z\""), // Converted to UTCc
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
				t1.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
