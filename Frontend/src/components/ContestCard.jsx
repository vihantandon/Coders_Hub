import { useState, useEffect } from 'react'
import api from '../api/client'
import './ContestCard.css'

const PLATFORM_META = {
  CodeChef:   { color: '#ff8c00', label: 'CC', url: 'https://www.codechef.com/contests' },
  Codeforces: { color: '#00d4ff', label: 'CF', url: 'https://codeforces.com/contests' },
  Leetcode:   { color: '#ffd700', label: 'LC', url: 'https://leetcode.com/contest' },
}

function useCountdown(startStr) {
  const [diff, setDiff] = useState(0)

  useEffect(() => {
    const calc = () => {
      const start = new Date(startStr.replace(' ', 'T') + 'Z')
      setDiff(Math.max(0, Math.floor((start - Date.now()) / 1000)))
    }
    calc()
    const t = setInterval(calc, 1000)
    return () => clearInterval(t)
  }, [startStr])

  const d = Math.floor(diff / 86400)
  const h = Math.floor((diff % 86400) / 3600)
  const m = Math.floor((diff % 3600) / 60)
  const s = diff % 60
  const pad = n => String(n).padStart(2, '0')

  return { d, h, m, s: pad(s), urgent: diff < 3600, soon: diff < 86400, diff }
}

const HOURS_OPTIONS = [1, 6, 24, 48]

export default function ContestCard({ contest, hasReminder, onReminderSet, onReminderDelete, animDelay }) {
  const { d, h, m, s, urgent, soon, diff } = useCountdown(contest.Start)
  const meta = PLATFORM_META[contest.Platform] || { color: 'var(--neon)', label: '??', url: '#' }

  const [showPicker, setShowPicker] = useState(false)
  const [selectedHours, setSelectedHours] = useState(24)
  const [loading, setLoading] = useState(false)
  const [msg, setMsg] = useState('')

  const setReminder = async () => {
    setLoading(true); setMsg('')
    try {
      await api.post('/reminders', { contest_id: contest.ID, hours_before: selectedHours })
      setMsg(`✓ reminder set for ${selectedHours}h before`)
      onReminderSet(contest.ID)
      setTimeout(() => { setShowPicker(false); setMsg('') }, 1500)
    } catch (err) {
      setMsg(`✗ ${err.response?.data?.error || 'failed'}`)
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (str) => {
    const d = new Date(str.replace(' ', 'T') + 'Z')
    return d.toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short', timeZone: 'Asia/Kolkata' }) + ' IST'
  }

  return (
    <div
      className={`contest-card card ${urgent ? 'urgent' : soon ? 'soon' : ''}`}
      style={{ animationDelay: `${animDelay}ms` }}
    >
      {/* Top bar */}
      <div className="cc-top">
        <div className="cc-platform" style={{ '--pcolor': meta.color }}>
          <span className="cc-platform-badge">{meta.label}</span>
          <span className="cc-platform-name">{contest.Platform}</span>
        </div>
        {urgent && <span className="cc-live-badge">LIVE SOON</span>}
      </div>

      {/* Name */}
      <div className="cc-name">{contest.Name}</div>

      {/* Countdown */}
      <div className={`cc-countdown ${urgent ? 'countdown-urgent' : ''}`}>
        <div className="cc-countdown-label">&gt; starts in</div>
        <div className="cc-countdown-display">
          {d > 0 && <><span className="cc-num">{d}</span><span className="cc-unit">d</span></>}
          <span className="cc-num">{String(h).padStart(2,'0')}</span><span className="cc-unit">h</span>
          <span className="cc-num">{String(m).padStart(2,'0')}</span><span className="cc-unit">m</span>
          <span className={`cc-num ${urgent ? 'cc-num-blink' : ''}`}>{s}</span><span className="cc-unit">s</span>
        </div>
      </div>

      {/* Times */}
      <div className="cc-times">
        <div className="cc-time-row">
          <span className="cc-time-label">start</span>
          <span className="cc-time-val">{formatDate(contest.Start)}</span>
        </div>
        <div className="cc-time-row">
          <span className="cc-time-label">end</span>
          <span className="cc-time-val">{formatDate(contest.End)}</span>
        </div>
      </div>

      {/* Actions */}
      <div className="cc-actions">
        <a href={meta.url} target="_blank" rel="noopener noreferrer" className="btn btn-ghost cc-goto">
          open →
        </a>

        {hasReminder ? (
          <button className="btn btn-danger cc-remind" disabled>
            ✓ reminded
          </button>
        ) : (
          <button
            className="btn cc-remind"
            onClick={() => setShowPicker(p => !p)}
          >
            ⏰ remind me
          </button>
        )}
      </div>

      {/* Reminder picker */}
      {showPicker && !hasReminder && (
        <div className="cc-picker">
          <div className="cc-picker-label">notify me before:</div>
          <div className="cc-picker-opts">
            {HOURS_OPTIONS.map(h => (
              <button
                key={h}
                className={`cc-picker-opt ${selectedHours === h ? 'selected' : ''}`}
                onClick={() => setSelectedHours(h)}
              >
                {h}h
              </button>
            ))}
          </div>
          <button className="btn cc-picker-confirm" onClick={setReminder} disabled={loading}>
            {loading ? <span className="spinner" /> : '> confirm'}
          </button>
          {msg && <div className={msg.startsWith('✓') ? 'success' : 'err'} style={{ fontSize: '0.7rem' }}>{msg}</div>}
        </div>
      )}
    </div>
  )
}
