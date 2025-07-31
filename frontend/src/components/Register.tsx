import { useState } from 'react';
import authService from '../services/authService';

const Register = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage('');
    setIsError(false);

    try {
      const response = await authService.register(username, password);
      setMessage(response.message || 'Registration successful');
      setUsername('');
      setPassword('');
    } catch (error) {
      setIsError(true);
      setMessage(error instanceof Error ? error.message : 'Registration failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-form">
      <h2>Register for Semantic Auth</h2>
      <div className="info-tooltip">
        <span className="info-icon">ℹ️</span>
        <div className="tooltip-content">
          <p>Semantic Auth is a password system that uses similarity matching instead of exact matches.</p>
          <p>When you log in later, your password will be compared to what you register now, and if it's similar enough, you'll be granted access.</p>
        </div>
      </div>
      <div className="warning-banner">
        <span className="warning-icon">⚠️</span>
        <p>This is a demo system. Please do not use real passwords or sensitive information.</p>
      </div>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            placeholder="Enter username"
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            placeholder="Enter password"
          />
        </div>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Registering...' : 'Register'}
        </button>
      </form>
      {message && (
        <div className={`message ${isError ? 'error' : 'success'}`}>
          {message}
        </div>
      )}
    </div>
  );
};

export default Register;
