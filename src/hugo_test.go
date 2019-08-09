package hugo

import (
	"strings"
	"testing"
)

func TestRepository_GetBranches(t *testing.T) {
	tests := []struct {
		name    string
		r       *Repository
		wantErr bool
	}{
		{
			"aa",
			&Repository{
				URL:  "https://github.com/gohugoio/hugo.git",
				Name: "hugo",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bList, err := tt.r.GetBranches()
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetBranches() error = %v, wantErr %v", err, tt.wantErr)
			}
			// t.Logf("%s", bList)
			if len(bList) > 2 {
				t.Errorf("Repository.GetBranches() bList = %v, wantErr %v", bList, tt.wantErr)
			}
		})
	}
}

func TestBranch_checkout(t *testing.T) {
	type fields struct {
		Name       string
		Hash       string
		Repository *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "aa",
			fields: fields{
				Hash: "209898998",
				Name: "master",
				Repository: &Repository{
					URL:  "https://github.com/gohugoio/hugo.git",
					Name: "hugo",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Branch{
				Name:       tt.fields.Name,
				Hash:       tt.fields.Hash,
				Repository: tt.fields.Repository,
			}
			if _, err := b.checkout(); (err != nil) != tt.wantErr {
				t.Errorf("Branch.checkout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBranch_build(t *testing.T) {
	type fields struct {
		Name       string
		Hash       string
		Repository *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "aa",
			fields: fields{
				Hash: "209898998",
				Name: "master",
				Repository: &Repository{
					URL:  "https://github.com/gohugoio/hugo.git",
					Name: "hugo",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Branch{
				Name:       tt.fields.Name,
				Hash:       tt.fields.Hash,
				Repository: tt.fields.Repository,
			}
			output, err := b.build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Branch.build() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !strings.Contains(output, "Total in") {
				t.Errorf("Output does not contain 'Total in', output is: %s", output)
			}
		})
	}
}

func TestBranch_buildAndUpload(t *testing.T) {
	type fields struct {
		Name       string
		Hash       string
		Repository *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "aa",
			fields: fields{
				Hash: "209898998",
				Name: "master",
				Repository: &Repository{
					URL:  "https://github.com/gohugoio/hugo.git",
					Name: "hugo",
				},
			},
			want:    "hugo-operator-hugo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Branch{
				Name:       tt.fields.Name,
				Hash:       tt.fields.Hash,
				Repository: tt.fields.Repository,
			}
			got, err := b.buildAndUpload()
			if (err != nil) != tt.wantErr {
				t.Errorf("Branch.buildAndUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Branch.buildAndUpload() = %v, want %v", got, tt.want)
			}
		})
	}
}
