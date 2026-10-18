package rsa

// Step — один шаг хода решения (для фронта)
type Step struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// KeyPair — пара ключей RSA
type KeyPair struct {
	N     string `json:"n"`
	E     string `json:"e"`
	D     string `json:"d"`
	P     string `json:"p"`
	Q     string `json:"q"`
	Phi   string `json:"phi"`
	Steps []Step `json:"steps"`
}

// EncryptResult — результат шифрования
type EncryptResult struct {
	Ciphertext string `json:"ciphertext"`
	Steps      []Step `json:"steps"`
}

// DecryptResult — результат расшифрования
type DecryptResult struct {
	Plaintext string `json:"plaintext"`
	Steps     []Step `json:"steps"`
}

// AttackResult — результат атаки
type AttackResult struct {
	Method    string `json:"method"`
	P         string `json:"p,omitempty"`
	Q         string `json:"q,omitempty"`
	D         string `json:"d,omitempty"`
	Plaintext string `json:"plaintext,omitempty"`
	Details   string `json:"details,omitempty"`
	Steps     []Step `json:"steps"`
}
