import { useState } from 'react';
import authService from '../services/authService';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [threshold, setThreshold] = useState(0.88);
  const [message, setMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);
  const [showThreshold, setShowThreshold] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage('');
    setIsError(false);

    try {
      const response = await authService.login(username, password, threshold);
      const successMessage = `Login successful! Similarity: ${response.data?.similarity.toFixed(4)}`;
      setMessage(response.message || successMessage);
    } catch (error) {
      setIsError(true);
      setMessage(error instanceof Error ? error.message : 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="auth-form">
      <h2>Login</h2>
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
        <div className="advanced-options">
          <button 
            type="button" 
            className="toggle-button"
            onClick={() => setShowThreshold(!showThreshold)}
          >
            {showThreshold ? 'Hide Advanced Options' : 'Show Advanced Options'}
          </button>
          
          {showThreshold && (
            <div className="form-group">
              <label htmlFor="threshold">
                Similarity Threshold ({threshold})
              </label>
              <input
                type="range"
                id="threshold"
                min="0.5"
                max="1"
                step="0.01"
                value={threshold}
                onChange={(e) => setThreshold(parseFloat(e.target.value))}
              />
              <div className="threshold-description">
                <small>Lower values are more lenient, higher values are stricter</small>
              </div>
            </div>
          )}
        </div>
        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Logging in...' : 'Login'}
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

export default Login;
