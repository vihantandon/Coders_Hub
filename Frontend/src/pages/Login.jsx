import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../api/client'
import './Auth.css'

export default function Login() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [typed, setTyped] = useState('')

  const intro = '> initializing secure connection...'
  useEffect(() => {
    let i = 0
    const t = setInterval(() => {
      setTyped(intro.slice(0, ++i))
      if (i >= intro.length) clearInterval(t)
    }, 35)
    return () => clearInterval(t)
  }, [])

  const handle = e => setForm(f => ({ ...f, [e.target.name]: e.target.value }))

  const submit = async e => {
    e.preventDefault()
    setError(''); setLoading(true)
    try {
      const res = await api.post('/login', form)
      localStorage.setItem('token', res.data.token)
      // decode name from token or store email
      localStorage.setItem('user', JSON.stringify({ name: form.email.split('@')[0] }))
      navigate('/')
    } catch (err) {
      setError(err.response?.data?.error || 'Connection failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-bg">
      <div className="auth-particles">
        {[...Array(20)].map((_, i) => (
          <div key={i} className="particle" style={{
            left: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 3}s`,
            animationDuration: `${2 + Math.random() * 3}s`
          }} />
        ))}
      </div>

      <div className="auth-container">
        <div className="auth-header">
          <div className="auth-logo">CODERS_HUB</div>
          <div className="auth-terminal-line">{typed}<span className="cursor">█</span></div>
        </div>

        <div className="auth-card card">
          <div className="auth-card-header">
            <span className="auth-card-title">// LOGIN</span>
            <div className="auth-dots">
              <span className="dot dot-red" />
              <span className="dot dot-yellow" />
              <span className="dot dot-green" />
            </div>
          </div>

          <form className="auth-form" onSubmit={submit}>
            <div className="field">
              <label className="field-label"><span className="label-prefix">$</span> EMAIL</label>
              <input
                name="email"
                type="email"
                value={form.email}
                onChange={handle}
                placeholder="user@domain.com"
                required
                autoComplete="off"
              />
            </div>
            <div className="field">
              <label className="field-label"><span className="label-prefix">$</span> PASSWORD</label>
              <input
                name="password"
                type="password"
                value={form.password}
                onChange={handle}
                placeholder="••••••••"
                required
              />
            </div>

            {error && <div className="err">⚠ {error}</div>}

            <button className="btn auth-btn" type="submit" disabled={loading}>
              {loading ? <><span className="spinner" /> AUTHENTICATING...</> : '> LOGIN'}
            </button>
          </form>

          <div className="auth-footer">
            <span className="text-dim">no account? </span>
            <Link to="/register" className="auth-link">register --new-user</Link>
          </div>
        </div>
      </div>
    </div>
  )
}
