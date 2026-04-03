import { useState } from 'react'
import axios from 'axios'
import './App.css'

const API_URL = 'http://localhost:8080/api'

function App() {
  const [pollId, setPollId] = useState('')
  const [poll, setPoll] = useState(null)
  const [selectedOption, setSelectedOption] = useState('')
  const [voteMessage, setVoteMessage] = useState('')
  const [results, setResults] = useState(null)

  const createPoll = async (e) => {
    e.preventDefault()
    const formData = new FormData(e.target)
    const question = formData.get('question')
    const options = formData.get('options').split(',').map((opt) => opt.trim())

    try {
      const response = await axios.post(`${API_URL}/polls`, { question, options })
      alert(`Poll created! ID: ${response.data.poll_id}`)
      e.target.reset()
    } catch {
      alert('Error creating poll')
    }
  }

  const loadPoll = async () => {
    if (!pollId) return
    try {
      const response = await axios.get(`${API_URL}/polls/${pollId}`)
      setPoll(response.data)
      setResults(null)
      setVoteMessage('')
    } catch {
      alert('Poll not found')
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
      loadResults()
    } catch (error) {
      if (error.response?.status === 409) {
        setVoteMessage('Notice: you already voted in this poll.')
      } else {
        setVoteMessage('Error: could not record vote.')
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
          <button type="submit">Create Poll</button>
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
          <button onClick={loadPoll}>Load Poll</button>
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
            <button onClick={vote}>Submit Vote</button>
            {voteMessage && <p className="message">{voteMessage}</p>}
          </div>
        )}
      </section>

      <section className="card">
        <div className="section-head">
          <h2>View Results</h2>
          <p>Refresh to pull latest percentages for the selected poll.</p>
        </div>

        <button onClick={loadResults}>Refresh Results</button>
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
