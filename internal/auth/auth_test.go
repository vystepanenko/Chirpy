package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	pass_1 := "SomePassword!"
	pass_2 := "AnotherPassword!"

	hash_1, _ := HashPassword(pass_1)
	hash_2, _ := HashPassword(pass_2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: pass_1,
			hash:     hash_1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash_1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: pass_1,
			hash:     hash_2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash_1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: pass_1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
