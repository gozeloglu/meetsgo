package main

import (
	"testing"
)

type Email struct {
	email string
	result bool
}

func TestIsValidEmail(t *testing.T)  {
	var emails [10]Email

	emails[0] = Email{"a@mail.co", true}
	emails[1] = Email{"@mail.co", false}
	emails[2] = Email{"12@mail.com", true}
	emails[3] = Email{"a@mail.", false}
	emails[4] = Email{"a@mail.c", false}
	emails[5] = Email{"name@gmail.com", true}
	emails[6] = Email{".name@gmail.com", true}
	emails[7] = Email{"name.surname", false}
	emails[8] = Email{"name.surname@mail", false}
	emails[9] = Email{"name.surname@gmail.com", true}

	for i := 0; i < len(emails); i++ {
		res := IsValidEMail(emails[i].email)
		if res != emails[i].result {
			t.Errorf("Email is checked, got %t, want: %t", res, emails[i].result)
		}
	}
}