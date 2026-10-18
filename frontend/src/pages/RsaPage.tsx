import { useState } from 'react'
import { apiPost } from '../utils/api'
import type { RsaKeyPair, RsaEncryptResponse, RsaDecryptResponse, RsaAttackResponse } from '../types'

type Step = { title: string; detail: string }

type KeyPairWithSteps = RsaKeyPair & { steps?: Step[]; phi?: string }
type EncryptWithSteps = RsaEncryptResponse & { steps?: Step[] }
type DecryptWithSteps = RsaDecryptResponse & { steps?: Step[] }
type AttackWithSteps = RsaAttackResponse & {
    steps?: Step[]
    p?: string
    q?: string
    d?: string
    details?: string
}

function StepsList({ steps }: { steps?: Step[] }) {
    if (!steps?.length) return null
    return (
        <div style={{ marginTop: 12 }}>
            <div className="result-label">Ход решения</div>
            {steps.map((s, i) => (
                <div key={i} className="result" style={{ marginTop: 6 }}>
                    <div className="result-label">
                        {i + 1}. {s.title}
                    </div>
                    <div className="mono">{s.detail}</div>
                </div>
            ))}
        </div>
    )
}

function ErrorBox({ message }: { message: string | null }) {
    if (!message) return null
    return (
        <div
            style={{
                marginBottom: 12,
                padding: '10px 14px',
                borderRadius: 6,
                border: '1px solid rgba(240, 113, 120, 0.4)',
                background: 'rgba(240, 113, 120, 0.1)',
                color: 'var(--danger)',
                fontSize: 14,
            }}
        >
            {message}
        </div>
    )
}

export default function RsaPage() {
    const [bits, setBits] = useState(64)
    const [keys, setKeys] = useState<KeyPairWithSteps | null>(null)
    const [message, setMessage] = useState('Hello RSA')
    const [ciphertext, setCiphertext] = useState('')
    const [plaintext, setPlaintext] = useState('')
    const [encryptSteps, setEncryptSteps] = useState<Step[] | undefined>()
    const [decryptSteps, setDecryptSteps] = useState<Step[] | undefined>()
    const [attackMethod, setAttackMethod] = useState('factor')
    const [attackResult, setAttackResult] = useState<AttackWithSteps | null>(null)

    const [genError, setGenError] = useState<string | null>(null)
    const [cryptoError, setCryptoError] = useState<string | null>(null)
    const [attackError, setAttackError] = useState<string | null>(null)
    const [loadingAction, setLoadingAction] = useState<string | null>(null)

    async function run<T>(
        action: string,
        setErr: (msg: string | null) => void,
        fn: () => Promise<T>,
    ): Promise<T | null> {
        setErr(null)
        setLoadingAction(action)
        try {
            return await fn()
        } catch (e) {
            setErr(e instanceof Error ? e.message : 'Ошибка')
            return null
        } finally {
            setLoadingAction(null)
        }
    }

    const onGenerate = () =>
        run('generate', setGenError, async () => {
            const data = await apiPost<KeyPairWithSteps>('/api/rsa/generate', { bits })
            setKeys(data)
            setCiphertext('')
            setPlaintext('')
            setEncryptSteps(undefined)
            setDecryptSteps(undefined)
            setAttackResult(null)
            setAttackError(null)
            setCryptoError(null)
            return data
        })

    const onEncrypt = () =>
        run('encrypt', setCryptoError, async () => {
            if (!keys) throw new Error('Сначала сгенерируйте ключи')
            const data = await apiPost<EncryptWithSteps>('/api/rsa/encrypt', {
                message,
                n: keys.n,
                e: keys.e,
            })
            setCiphertext(data.ciphertext)
            setEncryptSteps(data.steps)
            setPlaintext('')
            setDecryptSteps(undefined)
            return data
        })

    const onDecrypt = () =>
        run('decrypt', setCryptoError, async () => {
            if (!keys || !ciphertext) throw new Error('Нужны ключи и шифротекст')
            const data = await apiPost<DecryptWithSteps>('/api/rsa/decrypt', {
                ciphertext,
                n: keys.n,
                d: keys.d,
            })
            setPlaintext(data.plaintext)
            setDecryptSteps(data.steps)
            return data
        })

    const onAttack = () =>
        run('attack', setAttackError, async () => {
            if (!keys) throw new Error('Сначала сгенерируйте ключи')
            setAttackResult(null)
            const data = await apiPost<AttackWithSteps>('/api/rsa/attack', {
                n: keys.n,
                e: keys.e,
                ciphertext,
                method: attackMethod,
            })
            setAttackResult(data)
            return data
        })

    return (
        <div className="content">
            <div className="card">
                <h2>1. Криптосистема RSA</h2>
                <p>
                    Генерация ключей, шифрование / расшифрование и атаки. Backend — Go,
                    ход решения приходит в поле <code>steps</code>.
                </p>
            </div>

            {/* ——— Генерация ——— */}
            <div className="card">
                <h3>Генерация ключей</h3>
                <ErrorBox message={genError} />
                <div className="form-row">
                    <div className="form-group">
                        <label>Размер n (бит)</label>
                        <input
                            type="number"
                            value={bits}
                            min={20}
                            max={512}
                            onChange={(e) => setBits(Number(e.target.value))}
                        />
                    </div>
                </div>
                <p style={{ fontSize: 13, color: 'var(--text-muted)', marginBottom: 12 }}>
                    Для атак бери 32–48 бит. Для обычного шифрования — 128+.
                </p>
                <button className="btn btn-primary" disabled={!!loadingAction} onClick={onGenerate}>
                    {loadingAction === 'generate' ? 'Генерация…' : 'Сгенерировать'}
                </button>

                {keys && (
                    <>
                        <div className="form-row" style={{ marginTop: 16 }}>
                            <div className="form-group">
                                <label>n (модуль)</label>
                                <input readOnly value={keys.n} className="mono" />
                            </div>
                            <div className="form-group">
                                <label>e (открытая)</label>
                                <input readOnly value={keys.e} className="mono" />
                            </div>
                        </div>
                        <div className="form-row">
                            <div className="form-group">
                                <label>d (секретная)</label>
                                <input readOnly value={keys.d} className="mono" />
                            </div>
                            <div className="form-group">
                                <label>φ(n)</label>
                                <input readOnly value={keys.phi ?? ''} className="mono" />
                            </div>
                        </div>
                        <div className="form-row">
                            <div className="form-group">
                                <label>p</label>
                                <input readOnly value={keys.p ?? ''} className="mono" />
                            </div>
                            <div className="form-group">
                                <label>q</label>
                                <input readOnly value={keys.q ?? ''} className="mono" />
                            </div>
                        </div>
                        <StepsList steps={keys.steps} />
                    </>
                )}
            </div>

            {/* ——— Шифрование / расшифрование ——— */}
            <div className="card">
                <h3>Шифрование / расшифрование</h3>
                <ErrorBox message={cryptoError} />
                <div className="form-group">
                    <label>Сообщение</label>
                    <input value={message} onChange={(e) => setMessage(e.target.value)} />
                </div>
                <div className="btn-group">
                    <button className="btn btn-primary" disabled={!!loadingAction || !keys} onClick={onEncrypt}>
                        {loadingAction === 'encrypt' ? 'Шифрование…' : 'Зашифровать'}
                    </button>
                    <button
                        className="btn btn-secondary"
                        disabled={!!loadingAction || !keys || !ciphertext}
                        onClick={onDecrypt}
                    >
                        {loadingAction === 'decrypt' ? 'Расшифрование…' : 'Расшифровать'}
                    </button>
                </div>

                {ciphertext && (
                    <div className="form-group" style={{ marginTop: 14 }}>
                        <label>Шифротекст c</label>
                        <textarea readOnly value={ciphertext} className="mono" rows={3} />
                    </div>
                )}
                <StepsList steps={encryptSteps} />

                {plaintext !== '' && (
                    <div className="result" style={{ marginTop: 12 }}>
                        <div className="result-label">Расшифрованный текст</div>
                        <div className="mono">{plaintext}</div>
                    </div>
                )}
                <StepsList steps={decryptSteps} />
            </div>

            {/* ——— Атаки ——— */}
            <div className="card">
                <h3>Атаки</h3>
                <ErrorBox message={attackError} />
                <div className="form-group">
                    <label>Метод</label>
                    <select value={attackMethod} onChange={(e) => setAttackMethod(e.target.value)}>
                        <option value="factor">Пробное деление (factor)</option>
                        <option value="fermat">Метод Ферма (fermat)</option>
                        <option value="small_e">Малый e (small_e)</option>
                        <option value="wiener">Атака Винера (wiener)</option>
                    </select>
                </div>
                <button className="btn btn-danger" disabled={!!loadingAction || !keys} onClick={onAttack}>
                    {loadingAction === 'attack' ? 'Атака…' : 'Запустить атаку'}
                </button>

                {attackResult && (
                    <>
                        <div className="result" style={{ marginTop: 14 }}>
                            <div className="result-label">Метод: {attackResult.method}</div>
                            {attackResult.details && <div>{attackResult.details}</div>}
                            {attackResult.p && (
                                <div className="mono" style={{ marginTop: 6 }}>
                                    p = {attackResult.p}
                                    <br />
                                    q = {attackResult.q}
                                    <br />
                                    d = {attackResult.d}
                                </div>
                            )}
                            {attackResult.plaintext && (
                                <div className="mono" style={{ marginTop: 6 }}>
                                    plaintext = &quot;{attackResult.plaintext}&quot;
                                </div>
                            )}
                        </div>
                        <StepsList steps={attackResult.steps} />
                    </>
                )}
            </div>
        </div>
    )
}
