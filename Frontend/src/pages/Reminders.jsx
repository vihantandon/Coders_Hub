import { useState, useEffect } from 'react'
import api from '../api/client'
import './Reminders.css'

export default function Reminders() {
  const [reminders, setReminders] = useState([])
  const [loading, setLoading]     = useState(true)
  const [deleting, setDeleting]   = useState(null)
  const [error, setError]         = useState('')

  const fetch = async () => {
    setLoading(true)
    try {
      const res = await api.get('/reminders')
      setReminders(res.data.reminders || [])
    } catch {
      setError('> failed to load reminders')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetch() }, [])

  const del = async (id) => {
    setDeleting(id)
    try {
      await api.delete(`/reminders/${id}`)
      setReminders(prev => prev.filter(r => r.ID !== id))
    } catch {
      // ignore
    } finally {
      setDeleting(null)
    }
  }

  const formatDate = (str) => new Date(str).toLocaleString('en-IN', {
    dateStyle: 'medium', timeStyle: 'short', timeZone: 'Asia/Kolkata'
  }) + ' IST'

  const timeUntil = (sendAt) => {
    const diff = new Date(sendAt) - Date.now()
    if (diff < 0) return 'queued'
    const h = Math.floor(diff / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    return h > 0 ? `${h}h ${m}m` : `${m}m`
  }

  const PLATFORM_COLORS = {
    CodeChef: '#ff8c00',
    Codeforces: '#00d4ff',
    Leetcode: '#ffd700',
  }

  return (
    <div className="reminders-page">
      <div className="reminders-header">
        <h1 className="reminders-title">
          <span className="title-prompt">root@coders_hub:~$</span>
          <span className="title-cmd"> ls reminders/</span>
        </h1>
        <div className="reminders-count">
          {!loading && <span><span className="neon">{reminders.length}</span> active reminders</span>}
        </div>
      </div>

      {loading && (
        <div className="rem-loading">
          <div className="loading-spinner" />
          <span style={{ fontFamily: 'var(--mono)', fontSize: '0.8rem', color: 'var(--text-dim)' }}>
            &gt; loading reminders...
          </span>
        </div>
      )}

      {!loading && error && (
        <div className="rem-error">{error}</div>
      )}

      {!loading && !error && reminders.length === 0 && (
        <div className="rem-empty">
          <div className="rem-empty-icon">⏰</div>
          <div className="rem-empty-title">no reminders set</div>
          <div className="rem-empty-sub">&gt; go to contests page and set some!</div>
        </div>
      )}

      {!loading && !error && reminders.length > 0 && (
        <div className="reminders-list">
          {reminders.map((r, i) => {
            const contest = r.Contest || {}
            const color = PLATFORM_COLORS[contest.Platform] || 'var(--neon)'
            const isPast = r.Sent || new Date(r.SendAt) < Date.now()

            return (
              <div
                key={r.ID}
                className={`rem-card card ${isPast ? 'rem-sent' : ''}`}
                style={{ animationDelay: `${i * 70}ms` }}
              >
                <div className="rem-left">
                  <div className="rem-index" style={{ color }}>
                    {String(i + 1).padStart(2, '0')}
                  </div>
                  <div className="rem-info">
                    <div className="rem-contest-name">{contest.Name || '—'}</div>
                    <div className="rem-meta">
                      <span className="rem-platform" style={{ color }}>{contest.Platform}</span>
                      <span className="rem-sep">·</span>
                      <span className="rem-sendat">
                        notify at {formatDate(r.SendAt)}
                      </span>
                    </div>
                    <div className="rem-contest-start">
                      contest starts: {contest.Start ? formatDate(contest.Start) : '—'}
                    </div>
                  </div>
                </div>

                <div className="rem-right">
                  {isPast ? (
                    <span className="rem-badge-sent">SENT ✓</span>
                  ) : (
                    <span className="rem-badge-pending">
                      in {timeUntil(r.SendAt)}
                    </span>
                  )}
                  <button
                    className="btn btn-danger rem-del"
                    onClick={() => del(r.ID)}
                    disabled={deleting === r.ID}
                  >
                    {deleting === r.ID ? <span className="spinner" /> : 'rm'}
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
