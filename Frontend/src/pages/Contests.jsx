import { useState, useEffect } from 'react'
import api from '../api/client'
import ContestCard from '../components/ContestCard'
import './Contests.css'

const PLATFORMS = ['All', 'CodeChef', 'Codeforces', 'Leetcode']

const PLATFORM_COLORS = {
  CodeChef:   'var(--neon)',        // orange but we'll use neon for tab
  Codeforces: 'var(--neon2)',
  Leetcode:   '#ffd700',
}

export default function Contests() {
  const [contests, setContests]   = useState([])
  const [platform, setPlatform]   = useState('All')
  const [loading, setLoading]     = useState(true)
  const [error, setError]         = useState('')
  const [myReminders, setMyReminders] = useState([]) // contest IDs with reminders

  const fetchContests = async (p) => {
    setLoading(true); setError('')
    try {
      const params = p !== 'All' ? { platform: p } : {}
      const res = await api.get('/contests', { params })
      setContests(res.data.contests || [])
    } catch {
      setError('> failed to fetch contests. is the backend running?')
    } finally {
      setLoading(false)
    }
  }

  const fetchMyReminders = async () => {
    try {
      const res = await api.get('/reminders')
      const ids = (res.data.reminders || []).map(r => r.ContestID || r.contest_id)
      setMyReminders(ids)
    } catch {
      // not logged in or no reminders
    }
  }

  useEffect(() => { fetchContests(platform) }, [platform])
  useEffect(() => { fetchMyReminders() }, [])

  const onReminderSet = (contestId) => {
    setMyReminders(prev => [...prev, contestId])
  }
  const onReminderDelete = (contestId) => {
    setMyReminders(prev => prev.filter(id => id !== contestId))
  }

  return (
    <div className="contests-page">
      {/* Header */}
      <div className="contests-header">
        <div className="contests-title-row">
          <h1 className="contests-title">
            <span className="title-prompt">root@coders_hub:~$</span>
            <span className="title-cmd"> get contests</span>
          </h1>
          <div className="contests-count">
            {!loading && <span><span className="neon">{contests.length}</span> contests found</span>}
          </div>
        </div>

        {/* Platform tabs */}
        <div className="platform-tabs">
          {PLATFORMS.map(p => (
            <button
              key={p}
              className={`platform-tab ${platform === p ? 'active' : ''}`}
              onClick={() => setPlatform(p)}
              style={platform === p && p !== 'All' ? { '--tab-color': PLATFORM_COLORS[p] } : {}}
            >
              {p !== 'All' && <span className="tab-dot" style={{ background: PLATFORM_COLORS[p] }} />}
              {p}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="contests-content">
        {loading && (
          <div className="loading-screen">
            <div className="loading-spinner" />
            <span className="loading-text">&gt; fetching contests<span className="dots-anim" /></span>
          </div>
        )}

        {!loading && error && (
          <div className="error-screen">
            <div className="error-code">ERROR</div>
            <div className="error-msg">{error}</div>
          </div>
        )}

        {!loading && !error && contests.length === 0 && (
          <div className="empty-screen">
            <div className="empty-code">[]</div>
            <div className="empty-msg">&gt; no contests found for {platform}</div>
          </div>
        )}

        {!loading && !error && contests.length > 0 && (
          <div className="contests-grid">
            {contests.map((c, i) => (
              <ContestCard
                key={c.ID || i}
                contest={c}
                hasReminder={myReminders.includes(c.ID)}
                onReminderSet={onReminderSet}
                onReminderDelete={onReminderDelete}
                animDelay={i * 60}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
