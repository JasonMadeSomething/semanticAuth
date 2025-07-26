import { useState } from 'react'
import './App.css'
import Register from './components/Register'
import Login from './components/Login'
import Report from './components/Report'

function App() {
  const [activeTab, setActiveTab] = useState('login')

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Semantic Authentication</h1>
        <nav className="app-nav">
          <button 
            className={activeTab === 'login' ? 'active' : ''}
            onClick={() => setActiveTab('login')}
          >
            Login
          </button>
          <button 
            className={activeTab === 'register' ? 'active' : ''}
            onClick={() => setActiveTab('register')}
          >
            Register
          </button>
          <button 
            className={activeTab === 'report' ? 'active' : ''}
            onClick={() => setActiveTab('report')}
          >
            Report
          </button>
        </nav>
      </header>

      <main className="app-content">
        {activeTab === 'login' && <Login />}
        {activeTab === 'register' && <Register />}
        {activeTab === 'report' && <Report />}
      </main>

      <footer className="app-footer">
        <p>Semantic Authentication - Using embedding similarity for password verification</p>
      </footer>
    </div>
  )
}

export default App
