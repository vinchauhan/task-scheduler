package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func TestObjectIDArrayFromHex(t *testing.T) {
	type args struct {
		objectIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []primitive.ObjectID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ObjectIDArrayFromHex(tt.args.objectIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ObjectIDArrayFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectIDArrayFromHex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
