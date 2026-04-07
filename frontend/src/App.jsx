import { useState } from 'react'
import axios from 'axios'
import './App.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8081/api'

function App() {
  const [pollId, setPollId] = useState('')
  const [pollList, setPollList] = useState([])
  const [poll, setPoll] = useState(null)
  const [selectedOption, setSelectedOption] = useState('')
  const [voteMessage, setVoteMessage] = useState('')
  const [voteStatus, setVoteStatus] = useState('notice')
  const [results, setResults] = useState(null)

  const createPoll = async (e) => {
    e.preventDefault()
    const formData = new FormData(e.target)
    const question = formData.get('question')
    const options = formData
      .get('options')
      .split(',')
      .map((opt) => opt.trim())
      .filter(Boolean)

    if (options.length < 2) {
      alert('Please provide at least 2 non-empty options separated by commas.')
      return
    }

    try {
      const response = await axios.post(`${API_URL}/polls`, { question, options })
      alert(`Poll created! ID: ${response.data.poll_id}`)
      setPollId(response.data.poll_id)
      await loadRecentPolls()
      e.target.reset()
    } catch (error) {
      const backendMessage = error.response?.data?.error
      alert(backendMessage || 'Error creating poll')
    }
  }

  const loadPoll = async (id = pollId) => {
    if (!id) return
    try {
      const response = await axios.get(`${API_URL}/polls/${id}`)
      setPollId(id)
      setPoll(response.data)
      setSelectedOption('')
      setResults(null)
      setVoteMessage('')
      setVoteStatus('notice')
    } catch {
      alert('Poll not found')
    }
  }

  const loadRecentPolls = async () => {
    try {
      const response = await axios.get(`${API_URL}/polls?limit=20`)
      setPollList(response.data.polls || [])
    } catch {
      alert('Could not load recent polls')
    }
  }

  const vote = async () => {
    if (!selectedOption) return
    try {
      await axios.post(`${API_URL}/vote`, {
        poll_id: pollId,
        option_id: selectedOption,
      })
      setVoteMessage('Success: vote recorded.')
      setVoteStatus('success')
      loadResults()
    } catch (error) {
      if (error.response?.status === 409) {
        setVoteMessage('Notice: you already voted in this poll.')
        setVoteStatus('notice')
      } else {
        setVoteMessage('Error: could not record vote.')
        setVoteStatus('error')
      }
    }
  }

  const loadResults = async () => {
    try {
      const response = await axios.get(`${API_URL}/results/${pollId}`)
      setResults(response.data.results)
    } catch {
      console.error('Error loading results')
    }
  }

  return (
    <main className="container">
      <header className="hero">
        <p className="hero-kicker">Community Decisions</p>
        <h1>Online Voting System</h1>
        <p className="hero-subtitle">Create a poll, cast a vote, and watch live results update instantly.</p>
      </header>

      <section className="card">
        <div className="section-head">
          <h2>Create New Poll</h2>
          <p>Use comma-separated options, for example: Go, Python, JavaScript.</p>
        </div>

        <form onSubmit={createPoll} className="stack-form">
          <input name="question" placeholder="Poll question" required />
          <input name="options" placeholder="Option A, Option B, Option C" required />
          <button type="submit" className="btn btn-primary">Create Poll</button>
        </form>
      </section>

      <section className="card">
        <div className="section-head">
          <h2>Vote in Poll</h2>
          <p>Paste a poll ID and load options before submitting your vote.</p>
        </div>

        <div className="poll-input">
          <input
            type="text"
            placeholder="Enter poll ID"
            value={pollId}
            onChange={(e) => setPollId(e.target.value)}
          />
          <button onClick={loadPoll} className="btn btn-primary">Load Poll</button>
        </div>

        <div className="poll-input">
          <button onClick={loadRecentPolls} className="btn btn-secondary">Load Recent Polls</button>
          <select value={pollId} onChange={(e) => setPollId(e.target.value)}>
            <option value="">Select a poll</option>
            {pollList.map((p) => (
              <option key={p.id} value={p.id}>
                {p.question} ({p.id.slice(0, 8)})
              </option>
            ))}
          </select>
          <button onClick={() => loadPoll(pollId)} className="btn btn-primary">Open Selected</button>
        </div>

        {poll && (
          <div className="poll-details">
            <h3>{poll.poll.question}</h3>
            <div className="options">
              {poll.options.map((opt) => (
                <label key={opt.id} className="option">
                  <input
                    type="radio"
                    name="option"
                    value={opt.id}
                    onChange={(e) => setSelectedOption(e.target.value)}
                  />
                  <span>{opt.value}</span>
                  <strong>{opt.votes_count} votes</strong>
                </label>
              ))}
            </div>
            <button onClick={vote} className="btn btn-danger">Submit Vote</button>
            {voteMessage && <p className={`message message-${voteStatus}`}>{voteMessage}</p>}
          </div>
        )}
      </section>

      <section className="card">
        <div className="section-head">
          <h2>View Results</h2>
          <p>Refresh to pull latest percentages for the selected poll.</p>
        </div>

        <button onClick={loadResults} className="btn btn-secondary">Refresh Results</button>
        {results && (
          <div className="results">
            {results.map((r, idx) => (
              <div key={idx} className="result-bar">
                <span className="result-label">{r.option}</span>
                <div className="bar">
                  <div style={{ width: `${r.percentage}%` }} />
                </div>
                <span className="result-meta">{r.votes} votes ({r.percentage}%)</span>
              </div>
            ))}
          </div>
        )}
      </section>
    </main>
  )
}

export default App

