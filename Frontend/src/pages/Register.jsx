import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import api from '../api/client'
import './Auth.css'

export default function Register() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ name: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  const handle = e => setForm(f => ({ ...f, [e.target.name]: e.target.value }))

  const submit = async e => {
    e.preventDefault()
    setError(''); setSuccess(''); setLoading(true)
    try {
      await api.post('/register', form)
      setSuccess('> user created. redirecting...')
      setTimeout(() => navigate('/login'), 1200)
    } catch (err) {
      setError(err.response?.data?.error || 'Registration failed')
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
          <div className="auth-terminal-line">
            <span style={{ color: 'var(--neon2)' }}>&gt; creating new user instance...</span>
            <span className="cursor">█</span>
          </div>
        </div>

        <div className="auth-card card">
          <div className="auth-card-header">
            <span className="auth-card-title">// REGISTER</span>
            <div className="auth-dots">
              <span className="dot dot-red" />
              <span className="dot dot-yellow" />
              <span className="dot dot-green" />
            </div>
          </div>

          <form className="auth-form" onSubmit={submit}>
            <div className="field">
              <label className="field-label"><span className="label-prefix">$</span> NAME</label>
              <input name="name" value={form.name} onChange={handle} placeholder="your_handle" required />
            </div>
            <div className="field">
              <label className="field-label"><span className="label-prefix">$</span> EMAIL</label>
              <input name="email" type="email" value={form.email} onChange={handle} placeholder="user@domain.com" required />
            </div>
            <div className="field">
              <label className="field-label"><span className="label-prefix">$</span> PASSWORD <span style={{ color: 'var(--text-dim)', fontSize: '0.65rem' }}>(min 7 chars)</span></label>
              <input name="password" type="password" value={form.password} onChange={handle} placeholder="••••••••" required />
            </div>

            {error   && <div className="err">⚠ {error}</div>}
            {success && <div className="success">{success}</div>}

            <button className="btn auth-btn" type="submit" disabled={loading}>
              {loading ? <><span className="spinner" /> REGISTERING...</> : '> CREATE ACCOUNT'}
            </button>
          </form>

          <div className="auth-footer">
            <span className="text-dim">have account? </span>
            <Link to="/login" className="auth-link">login --existing</Link>
          </div>
        </div>
      </div>
    </div>
  )
}
