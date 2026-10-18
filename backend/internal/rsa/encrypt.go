package rsa

import (
	"errors"
	"math/big"
)

func Encrypt(message, nStr, eStr string) (*EncryptResult, error) {
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	e, ok := new(big.Int).SetString(eStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга e")
	}

	m := new(big.Int).SetBytes([]byte(message))
	if m.Cmp(n) >= 0 {
		return nil, errors.New("сообщение слишком длинное для n")
	}

	c := new(big.Int).Exp(m, e, n)

	return &EncryptResult{
		Ciphertext: c.String(),
		Steps: []Step{
			{Title: "Сообщение → число m", Detail: "message = \"" + message + "\", m = " + m.String()},
			{Title: "Проверка m < n", Detail: "m < n: OK"},
			{Title: "Шифрование", Detail: "c = m^e mod n = " + c.String()},
		},
	}, nil
}

func Decrypt(ciphertext, nStr, dStr string) (*DecryptResult, error) {
	c, ok := new(big.Int).SetString(ciphertext, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга c")
	}
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	d, ok := new(big.Int).SetString(dStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга d")
	}

	m := new(big.Int).Exp(c, d, n)
	plaintext := string(m.Bytes())

	return &DecryptResult{
		Plaintext: plaintext,
		Steps: []Step{
			{Title: "Шифротекст → число c", Detail: "c = " + c.String()},
			{Title: "Расшифрование", Detail: "m = c^d mod n = " + m.String()},
			{Title: "Число → текст", Detail: "plaintext = \"" + plaintext + "\""},
		},
	}, nil
}
