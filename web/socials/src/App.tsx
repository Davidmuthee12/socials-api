import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import './App.css'

const API_URL = import.meta.env.VITE_API_URL

function App() {
  return (
    <main className="page-shell">
      <section className="card hero-card">
        <p className="eyebrow">Socials</p>
        <h1>Home page</h1>
        <p className="support-copy">
          Your account activation link lands on a dedicated confirmation screen.
        </p>
      </section>
    </main>
  )
}

export function ActivateAccountPage() {
  const { token = ' ' } = useParams()
  const navigate = useNavigate()
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle')
  const [message, setMessage] = useState('Press activate to confirm your email address.')

  useEffect(() => {
    if (status !== 'success') {
      return undefined
    }

    const timeoutId = window.setTimeout(() => {
      navigate('/', { replace: true })
    }, 1200)

    return () => window.clearTimeout(timeoutId)
  }, [navigate, status])

  async function handleActivate() {
    const trimmedToken = token.trim()

    if (!trimmedToken) {
      setStatus('error')
      setMessage('Activation token is missing.')
      return
    }

    setStatus('loading')
    setMessage('Activating your account...')

    try {
      const response = await fetch(`${API_URL}/users/activate/${trimmedToken}`, {
        method: 'PUT',
      })

      if (!response.ok) {
        throw new Error('Activation failed')
      }

      setStatus('success')
      setMessage('Email confirmed. Redirecting to the home page...')
    } catch {
      setStatus('error')
      setMessage('Unable to activate your account. Please try the link again.')
    }
  }

  return (
    <main className="page-shell">
      <section className="card confirmation-card">
        <p className="eyebrow">Account activation</p>
        <div>
          <h1>Confirm email</h1>
          <p className="support-copy">Use the button below to activate your account.</p>
        </div>
        <button
          type="button"
          className="primary-button"
          onClick={handleActivate}
          disabled={status === 'loading' || status === 'success'}
        >
          {status === 'loading' ? 'Activating...' : 'Activate'}
        </button>
        <p className={`status-message ${status}`}>{message}</p>
      </section>
    </main>
  )
}

export default App
