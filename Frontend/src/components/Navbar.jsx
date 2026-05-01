import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useState, useEffect } from 'react'
import './Navbar.css'

export default function Navbar() {
  const navigate = useNavigate()
  const location = useLocation()
  const [time, setTime] = useState(new Date())
  const user = JSON.parse(localStorage.getItem('user') || '{}')

  useEffect(() => {
    const t = setInterval(() => setTime(new Date()), 1000)
    return () => clearInterval(t)
  }, [])

  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/login')
  }

  const fmt = (n) => String(n).padStart(2, '0')
  const timeStr = `${fmt(time.getHours())}:${fmt(time.getMinutes())}:${fmt(time.getSeconds())}`

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <Link to="/" className="navbar-logo">
          <span className="logo-bracket">[</span>
          <span className="logo-text">CODERS_HUB</span>
          <span className="logo-bracket">]</span>
        </Link>
        <div className="navbar-links">
          <Link to="/" className={`nav-link ${location.pathname === '/' ? 'active' : ''}`}>
            <span className="nav-prefix">./</span>contests
          </Link>
          <Link to="/reminders" className={`nav-link ${location.pathname === '/reminders' ? 'active' : ''}`}>
            <span className="nav-prefix">./</span>reminders
          </Link>
        </div>
      </div>
      <div className="navbar-right">
        <span className="nav-clock">{timeStr}</span>
        <span className="nav-user">
          <span className="nav-prefix">@</span>{user.name || 'user'}
        </span>
        <button className="btn btn-ghost nav-logout" onClick={logout}>logout</button>
      </div>
    </nav>
  )
}
