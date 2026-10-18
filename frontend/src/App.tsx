import { Routes, Route, NavLink } from 'react-router-dom'
import './App.css'

import HomePage from './pages/HomePage'
import RsaPage from './pages/RsaPage'
import ElGamalPage from './pages/ElGamalPage'
import EcPage from './pages/EcPage'
import PolynomialPage from './pages/PolynomialPage'

const topics = [
  { path: '/rsa', label: 'RSA', num: 1 },
  { path: '/elgamal', label: 'Эль-Гамаль', num: 2 },
  { path: '/ec', label: 'Эллиптические кривые', num: 3 },
  { path: '/polynomial', label: 'Многочлены', num: 4 },
]

function App() {
  return (
    <div className="app">
      <header className="header">
        <NavLink to="/" className="logo">
          <span>Crypto</span> Lab
        </NavLink>

        <nav className="nav">
          {topics.map((t) => (
            <NavLink
              key={t.path}
              to={t.path}
              className={({ isActive }) =>
                `nav-btn${isActive ? ' active' : ''}`
              }
            >
              {t.num}. {t.label}
            </NavLink>
          ))}
        </nav>
      </header>

      <main className="main">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/rsa" element={<RsaPage />} />
          <Route path="/elgamal" element={<ElGamalPage />} />
          <Route path="/ec" element={<EcPage />} />
          <Route path="/polynomial" element={<PolynomialPage />} />
        </Routes>
      </main>
    </div>
  )
}

export default App
