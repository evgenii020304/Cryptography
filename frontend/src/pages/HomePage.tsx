import { Link } from 'react-router-dom'

const topics = [
  {
    path: '/rsa',
    num: 1,
    title: 'RSA',
    desc: 'Криптосистема RSA и алгоритмы атаки на систему',
    status: 'ready' as const,
  },
  {
    path: '/elgamal',
    num: 2,
    title: 'Эль-Гамаль',
    desc: 'Криптосистема Эль-Гамаля',
    status: 'soon' as const,
  },
  {
    path: '/ec',
    num: 3,
    title: 'Эллиптические кривые',
    desc: 'Криптосистема на эллиптических кривых',
    status: 'soon' as const,
  },
  {
    path: '/polynomial',
    num: 4,
    title: 'Многочлены',
    desc: 'Криптосистема с использованием многочленов',
    status: 'soon' as const,
  },
]

const statusLabel = {
  ready: 'В работе',
  wip: 'В разработке',
  soon: 'Скоро',
}

export default function HomePage() {
  return (
    <div className="content" style={{ maxWidth: 720, margin: '0 auto' }}>
      <div className="card">
        <h2>Математические основы криптологии</h2>
        <p>
          Портфолио заданий по криптографическим методам. Backend — Go,
          Frontend — React + TypeScript. Начните с RSA.
        </p>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
        {topics.map((t) => (
          <Link
            key={t.path}
            to={t.path}
            className="card"
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 16,
              textDecoration: 'none',
              color: 'inherit',
              transition: 'border-color 0.15s',
            }}
          >
            <div className="topic-num" style={{ width: 36, height: 36, fontSize: 14 }}>
              {t.num}
            </div>
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 600, marginBottom: 2 }}>{t.title}</div>
              <div style={{ fontSize: 13, color: 'var(--text-muted)' }}>{t.desc}</div>
            </div>
            <span className={`badge badge-${t.status}`}>{statusLabel[t.status]}</span>
          </Link>
        ))}
      </div>
    </div>
  )
}
