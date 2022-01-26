package accounts

import (
	"errors"
	"fmt"
)

type Account struct {
	name    string
	balance int
	deposit int
}

func NewAccount(name string, balance int) *Account {
	return &Account{name: name, balance: balance}
}

func (account *Account) Deposit(deposit int) {
	account.balance += deposit
}
func (account *Account) Balance() {
	fmt.Print(account.balance)
}

func (account *Account) Withdraw(balance int) error {
	account.balance -= balance
	return errors.New("hi")
	return nil
}
