package rsa

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Generate создаёт пару ключей RSA.
// bits — размер модуля n в битах (для демо атак бери 32–64, для «нормальных» — 256+).
//
// Алгоритм:
//  1. Сгенерировать два простых p, q
//  2. n = p * q
//  3. φ(n) = (p-1)*(q-1)
//  4. Выбрать e: gcd(e, φ) = 1 (часто 65537)
//  5. d = e⁻¹ mod φ
//
// Верни KeyPair со steps — каждый шаг алгоритма.

func Generate(bits int) (*KeyPair, error) {
	if bits < 32 {
		bits = 32
	}

	p, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, err
	}

	q, err := rand.Prime(rand.Reader, bits/2)
	if err != nil {
		return nil, err
	}

	for p.Cmp(q) == 0 {
		q, err = rand.Prime(rand.Reader, bits/2)
		if err != nil {
			return nil, err
		}
	}

	n := new(big.Int)
	n.Mul(p, q)

	one := big.NewInt(1)
	pMinus1 := new(big.Int).Sub(p, one)
	qMinus1 := new(big.Int).Sub(q, one)
	phi := new(big.Int).Mul(pMinus1, qMinus1)

	e := big.NewInt(65537)
	res := new(big.Int).GCD(nil, nil, e, phi)
	if res.Cmp(one) != 0 {
		e = big.NewInt(3)
		for {
			res.GCD(nil, nil, e, phi)
			if res.Cmp(one) == 0 {
				break
			}
			e.Add(e, big.NewInt(2))
		}
	}

	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		return nil, errors.New("не удалось вычислить d: нет обратного элемента")
	}

	return &KeyPair{
		N:   n.String(),
		E:   e.String(),
		D:   d.String(),
		P:   p.String(),
		Q:   q.String(),
		Phi: phi.String(),
		Steps: []Step{
			{Title: "Простые p и q", Detail: "p = " + p.String() + ", q = " + q.String()},
			{Title: "Модуль n", Detail: "n = p·q = " + n.String()},
			{Title: "Функция Эйлера", Detail: "φ(n) = (p-1)(q-1) = " + phi.String()},
			{Title: "Открытая экспонента e", Detail: "e = " + e.String()},
			{Title: "Секретная экспонента d", Detail: "d ≡ e⁻¹ (mod φ) = " + d.String()},
		},
	}, nil
}
