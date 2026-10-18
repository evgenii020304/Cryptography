export interface ApiResponse<T> {
  data?: T
  error?: string
}

export interface Step {
  title: string
  detail: string
}

export interface RsaKeyPair {
  n: string
  e: string
  d: string
  p?: string
  q?: string
  phi?: string
  steps?: Step[]
}

export interface RsaEncryptRequest {
  message: string
  n: string
  e: string
}

export interface RsaEncryptResponse {
  ciphertext: string
  steps?: Step[]
}

export interface RsaDecryptRequest {
  ciphertext: string
  n: string
  d: string
}

export interface RsaDecryptResponse {
  plaintext: string
  steps?: Step[]
}

export interface RsaAttackRequest {
  n: string
  e: string
  ciphertext?: string
  method?: string
}

export interface RsaAttackResponse {
  method: string
  p?: string
  q?: string
  d?: string
  factors?: { p: string; q: string }
  privateKey?: string
  plaintext?: string
  details?: string
  steps?: Step[]
}