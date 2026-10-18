package rsa

import (
	"errors"
	"fmt"
	"math/big"
)

// AttackFactor — пробное деление n (для малых ключей, bits ≈ 32–48).
func AttackFactor(nStr, eStr, ciphertext string) (*AttackResult, error) {
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	e, ok := new(big.Int).SetString(eStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга e")
	}

	steps := []Step{
		{Title: "Метод", Detail: "Пробное деление (trial division)"},
		{Title: "Цель", Detail: "Найти простой делитель p ≤ √n"},
	}

	sqrtN := new(big.Int).Sqrt(n)
	p := big.NewInt(2)
	one := big.NewInt(1)
	two := big.NewInt(2)

	// чётность
	if new(big.Int).Mod(n, two).Cmp(big.NewInt(0)) == 0 {
		return finishFactor(n, e, two, ciphertext, steps)
	}
	p.SetInt64(3)

	for p.Cmp(sqrtN) <= 0 {
		if new(big.Int).Mod(n, p).Cmp(big.NewInt(0)) == 0 {
			steps = append(steps, Step{
				Title:  "Найден делитель",
				Detail: fmt.Sprintf("p = %s", p.String()),
			})
			return finishFactor(n, e, p, ciphertext, steps)
		}
		p.Add(p, two) // только нечётные
	}

	_ = one
	return nil, errors.New("делители не найдены (n слишком большое для пробного деления)")
}

// AttackFermat — метод Ферма (эффективен при p ≈ q).
func AttackFermat(nStr, eStr, ciphertext string) (*AttackResult, error) {
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	e, ok := new(big.Int).SetString(eStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга e")
	}

	steps := []Step{
		{Title: "Метод", Detail: "Факторизация Ферма"},
		{Title: "Идея", Detail: "n = a² − b² = (a−b)(a+b), ищем a ≥ ceil(√n)"},
	}

	a := new(big.Int).Sqrt(n)
	// если n не полный квадрат — a = floor(√n)+1
	if new(big.Int).Mul(a, a).Cmp(n) < 0 {
		a.Add(a, big.NewInt(1))
	}

	maxIter := 1_000_000 // защита от бесконечного цикла
	for i := 0; i < maxIter; i++ {
		b2 := new(big.Int).Sub(new(big.Int).Mul(a, a), n) // a² − n
		b := new(big.Int).Sqrt(b2)
		if new(big.Int).Mul(b, b).Cmp(b2) == 0 {
			// нашли: p = a−b, q = a+b
			p := new(big.Int).Sub(a, b)
			q := new(big.Int).Add(a, b)
			steps = append(steps,
				Step{Title: "Найдено a, b", Detail: fmt.Sprintf("a = %s, b = %s (итераций: %d)", a.String(), b.String(), i+1)},
				Step{Title: "Делители", Detail: fmt.Sprintf("p = a−b = %s, q = a+b = %s", p.String(), q.String())},
			)
			return finishFactor(n, e, p, ciphertext, steps)
		}
		a.Add(a, big.NewInt(1))
	}

	return nil, errors.New("метод Ферма: не найдено за отведённое число итераций (p и q слишком далеки)")
}

// AttackSmallE — атака при малом e, когда m^e < n (целочисленный корень).
func AttackSmallE(nStr, eStr, ciphertext string) (*AttackResult, error) {
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	e, ok := new(big.Int).SetString(eStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга e")
	}
	c, ok := new(big.Int).SetString(ciphertext, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга ciphertext")
	}

	steps := []Step{
		{Title: "Метод", Detail: "Атака на малый e (когда m^e < n)"},
		{Title: "Идея", Detail: "c = m^e ⇒ m = ∛[e]{c} (целочисленный корень)"},
	}

	if !e.IsInt64() || e.Int64() < 2 {
		return nil, errors.New("e слишком большое или < 2 для этой атаки")
	}
	exp := e.Int64()

	// бинарный поиск целочисленного корня степени exp
	m := intRoot(c, exp)
	if m == nil {
		return nil, errors.New("не удалось извлечь корень")
	}

	if m.Cmp(n) >= 0 {
		return nil, errors.New("m ≥ n — атака неприменима")
	}

	// проверка: m^e == c
	check := new(big.Int).Exp(m, e, nil)
	if check.Cmp(c) != 0 {
		steps = append(steps, Step{
			Title:  "Проверка",
			Detail: "m^e ≠ c — условие m^e < n не выполнено, атака неприменима",
		})
		return nil, errors.New("атака неприменима: m^e ≥ n (нужен малый e и короткое сообщение)")
	}

	plaintext := string(m.Bytes())
	steps = append(steps,
		Step{Title: "Корень", Detail: fmt.Sprintf("m = c^(1/%s) = %s", e.String(), m.String())},
		Step{Title: "Текст", Detail: "plaintext = \"" + plaintext + "\""},
	)

	return &AttackResult{
		Method:    "small_e",
		Plaintext: plaintext,
		Details:   "восстановлено без факторизации",
		Steps:     steps,
	}, nil
}

// AttackWiener — атака Винера (малый d, d < n^{1/4}/3).
func AttackWiener(nStr, eStr, ciphertext string) (*AttackResult, error) {
	n, ok := new(big.Int).SetString(nStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга n")
	}
	e, ok := new(big.Int).SetString(eStr, 10)
	if !ok {
		return nil, errors.New("ошибка парсинга e")
	}

	steps := []Step{
		{Title: "Метод", Detail: "Атака Винера (цепные дроби)"},
		{Title: "Условие", Detail: "d < n^{1/4}/3"},
		{Title: "Идея", Detail: "подходящие дроби e/n → кандидаты k/d"},
	}

	// цепная дробь для e/n
	convergents := continuedFraction(e, n)
	steps = append(steps, Step{
		Title:  "Подходящие дроби",
		Detail: fmt.Sprintf("найдено %d подходящих дробей k/d", len(convergents)),
	})

	one := big.NewInt(1)
	for _, frac := range convergents {
		k, d := frac[0], frac[1]
		if d.Cmp(one) == 0 {
			continue
		}
		// φ = (e*d − 1) / k  должно быть целым
		edMinus1 := new(big.Int).Sub(new(big.Int).Mul(e, d), one)
		if k.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		if new(big.Int).Mod(edMinus1, k).Cmp(big.NewInt(0)) != 0 {
			continue
		}
		phi := new(big.Int).Div(edMinus1, k)

		// уравнение x² − (n−φ+1)x + n = 0
		// sum = p+q = n − phi + 1
		sum := new(big.Int).Add(new(big.Int).Sub(n, phi), one)
		// disc = sum² − 4n
		disc := new(big.Int).Sub(new(big.Int).Mul(sum, sum), new(big.Int).Lsh(n, 2))
		if disc.Sign() < 0 {
			continue
		}
		sqrtDisc := new(big.Int).Sqrt(disc)
		if new(big.Int).Mul(sqrtDisc, sqrtDisc).Cmp(disc) != 0 {
			continue
		}

		p := new(big.Int).Rsh(new(big.Int).Add(sum, sqrtDisc), 1)
		q := new(big.Int).Rsh(new(big.Int).Sub(sum, sqrtDisc), 1)

		if new(big.Int).Mul(p, q).Cmp(n) != 0 {
			continue
		}

		steps = append(steps,
			Step{Title: "Кандидат d", Detail: fmt.Sprintf("d = %s, k = %s", d.String(), k.String())},
			Step{Title: "φ(n)", Detail: phi.String()},
			Step{Title: "Факторы", Detail: fmt.Sprintf("p = %s, q = %s", p.String(), q.String())},
		)

		return finishFactorWithD(n, e, p, q, d, ciphertext, steps)
	}

	return nil, errors.New("атака Винера не удалась (d недостаточно мал)")
}

// ——— helpers ———

func finishFactor(n, e, p *big.Int, ciphertext string, steps []Step) (*AttackResult, error) {
	q := new(big.Int).Div(n, p)
	one := big.NewInt(1)
	phi := new(big.Int).Mul(new(big.Int).Sub(p, one), new(big.Int).Sub(q, one))
	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		return nil, errors.New("не удалось восстановить d")
	}

	steps = append(steps,
		Step{Title: "q = n/p", Detail: q.String()},
		Step{Title: "φ(n)", Detail: phi.String()},
		Step{Title: "d = e⁻¹ mod φ", Detail: d.String()},
	)

	result := &AttackResult{
		Method:  "factor",
		P:       p.String(),
		Q:       q.String(),
		D:       d.String(),
		Details: "ключ восстановлен через факторизацию",
		Steps:   steps,
	}

	if ciphertext != "" {
		c, ok := new(big.Int).SetString(ciphertext, 10)
		if ok {
			m := new(big.Int).Exp(c, d, n)
			result.Plaintext = string(m.Bytes())
			steps = append(steps, Step{Title: "Расшифровка", Detail: "plaintext = \"" + result.Plaintext + "\""})
			result.Steps = steps
		}
	}

	return result, nil
}

func finishFactorWithD(n, e, p, q, d *big.Int, ciphertext string, steps []Step) (*AttackResult, error) {
	result := &AttackResult{
		Method:  "wiener",
		P:       p.String(),
		Q:       q.String(),
		D:       d.String(),
		Details: "ключ восстановлен атакой Винера",
		Steps:   steps,
	}

	if ciphertext != "" {
		c, ok := new(big.Int).SetString(ciphertext, 10)
		if ok {
			m := new(big.Int).Exp(c, d, n)
			result.Plaintext = string(m.Bytes())
			result.Steps = append(result.Steps, Step{
				Title:  "Расшифровка",
				Detail: "plaintext = \"" + result.Plaintext + "\"",
			})
		}
	}
	_ = e
	return result, nil
}

// intRoot — целочисленный корень степени exp из x (бинарный поиск).
func intRoot(x *big.Int, exp int64) *big.Int {
	if x.Sign() < 0 {
		return nil
	}
	if x.Sign() == 0 {
		return big.NewInt(0)
	}

	low := big.NewInt(1)
	high := new(big.Int).Set(x)
	e := big.NewInt(exp)

	for low.Cmp(high) <= 0 {
		mid := new(big.Int).Rsh(new(big.Int).Add(low, high), 1)
		pow := new(big.Int).Exp(mid, e, nil)
		cmp := pow.Cmp(x)
		if cmp == 0 {
			return mid
		}
		if cmp < 0 {
			low.Add(mid, big.NewInt(1))
		} else {
			high.Sub(mid, big.NewInt(1))
		}
	}
	// ближайший снизу
	return high
}

// continuedFraction — подходящие дроби для a/b.
func continuedFraction(a, b *big.Int) [][2]*big.Int {
	var cf []*big.Int
	aa, bb := new(big.Int).Set(a), new(big.Int).Set(b)

	for bb.Sign() != 0 {
		q := new(big.Int).Div(aa, bb)
		cf = append(cf, new(big.Int).Set(q))
		aa, bb = bb, new(big.Int).Mod(aa, bb)
	}

	// подходящие дроби hₙ/kₙ
	var result [][2]*big.Int
	h1, h2 := big.NewInt(0), big.NewInt(1) // h₋₂, h₋₁
	k1, k2 := big.NewInt(1), big.NewInt(0) // k₋₂, k₋₁

	for _, q := range cf {
		h := new(big.Int).Add(new(big.Int).Mul(q, h2), h1)
		k := new(big.Int).Add(new(big.Int).Mul(q, k2), k1)
		result = append(result, [2]*big.Int{new(big.Int).Set(h), new(big.Int).Set(k)})
		h1, h2 = h2, h
		k1, k2 = k2, k
	}
	return result
}
