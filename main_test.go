package main

import (
	"testing"
)

type Email struct {
	email  string
	result bool
}

type Password struct {
	password string
	result   bool
}

type TestUser struct {
	user    User
	isValid bool
	reason  InvalidReason
}

func TestIsValidEmail(t *testing.T) {
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

func TestIsValidPassword(t *testing.T) {
	var passwords [10]Password

	passwords[0] = Password{"", false}
	passwords[1] = Password{"1", false}
	passwords[2] = Password{"12345678", false}
	passwords[3] = Password{"asdfghj", false}
	passwords[4] = Password{"1234as.!", true}
	passwords[5] = Password{"a", false}
	passwords[6] = Password{".", false}
	passwords[7] = Password{".12", false}
	passwords[8] = Password{"asdf1234", false}
	passwords[9] = Password{"aaaa1234.-", true}

	for i := 0; i < len(passwords); i++ {
		res, _ := IsValidPassword(passwords[i].password)
		if res != passwords[i].result {
			t.Errorf("%s is checked, got %t, want: %t", passwords[i].password, res, passwords[i].result)
		}
	}
}

func TestIsValidUser(t *testing.T) {
	var testUsers [10]TestUser
	var emptyMeetupArr []*Meetup

	testUsers[0] = TestUser{User{1, "", "john", "jack", "asdf1234.-", "a@mail.com", 10, false, emptyMeetupArr}, false, UsernameShort}
	testUsers[1] = TestUser{User{1, "j", "john", "jack", "asdf1234.-", "a@mail.com", 10, false, emptyMeetupArr}, false, UsernameShort}
	testUsers[2] = TestUser{User{1, "john", "", "jack", "asdf1234.-", "a@mail.com", 10, false, emptyMeetupArr}, false, NameEmpty}
	testUsers[3] = TestUser{User{1, "john", "john", "", "asdf1234.-", "a@mail.com", 10, false, emptyMeetupArr}, false, SurnameEmpty}
	testUsers[4] = TestUser{User{1, "john", "john", "jack", "", "a@mail.com", 10, false, emptyMeetupArr}, false, PasswordShort}
	testUsers[5] = TestUser{User{1, "john", "john", "jack", "asdf1234", "a@mail.com", 10, false, emptyMeetupArr}, false, PasswordWeak}
	testUsers[6] = TestUser{User{1, "john", "john", "jack", "asdf", "a@mail.com", 10, false, emptyMeetupArr}, false, PasswordShort}
	testUsers[7] = TestUser{User{1, "john", "john", "jack", "asdf1234.-", "a@mail.", 10, false, emptyMeetupArr}, false, EmailINotValid}
	testUsers[8] = TestUser{User{1, "john", "john", "jack", "asdf1234.-", "a@mail.com", -1, false, emptyMeetupArr}, false, AgeNotValid}
	testUsers[9] = TestUser{User{1, "john", "john", "jack", "asdf1234.-", "a@mail.com", 10, false, emptyMeetupArr}, true, IsValid}

	for i := 0; i < len(testUsers); i++ {
		isValid, reason := IsValidUser(testUsers[i].user)

		if isValid != testUsers[i].isValid || reason != testUsers[i].reason {
			t.Errorf("%v", testUsers[i])
			t.Errorf("User is checked, got %t, want: %t", isValid, testUsers[i].isValid)
			t.Errorf("User's reason, got %d, want: %d", reason, testUsers[i].reason)
		}
	}

}
