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
  const [showPassword, setShowPassword] = useState(false);

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
      <h2>Login to Semantic Auth</h2>
      <p className="auth-description">
        Semantic Auth uses fuzzy password matching based on similarity scores rather than exact matches.
        Your password will be compared to what you registered, and if it's similar enough, you'll be logged in.
      </p>
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
          <div className="password-input-container">
            <input
              type={showPassword ? "text" : "password"}
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Enter password"
            />
            <button 
              type="button" 
              className="password-toggle-icon"
              onClick={() => setShowPassword(!showPassword)}
              aria-label={showPassword ? "Hide password" : "Show password"}
              title={showPassword ? "Hide password" : "Show password"}
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="18" height="18" fill="currentColor">
                {showPassword ? (
                  <path d="M12 7c-2.76 0-5 2.24-5 5 0 .65.13 1.26.36 1.83l2.92-2.92c.29-.29.77-.29 1.06 0 .29.29.29.77 0 1.06l-2.92 2.92c.57.23 1.18.36 1.83.36 2.76 0 5-2.24 5-5s-2.24-5-5-5zm0-5C7 2 2.73 4.11 1 7.5 2.73 10.89 7 13 12 13s9.27-2.11 11-5.5C21.27 4.11 17 2 12 2zm0 13c-3.79 0-7.17-2.13-9-5.5 1.83-3.37 5.21-5.5 9-5.5s7.17 2.13 9 5.5c-1.83 3.37-5.21 5.5-9 5.5zm-3.5-5c0-1.93 1.57-3.5 3.5-3.5s3.5 1.57 3.5 3.5-1.57 3.5-3.5 3.5-3.5-1.57-3.5-3.5z"/>
                ) : (
                  <path d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46C3.08 8.3 1.78 10.02 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.31-.78l3.15 3.15.02-.16c0-1.66-1.34-3-3-3l-.17.01z"/>
                )}
              </svg>
            </button>
          </div>
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
                Similarity Threshold
              </label>
              <div className="numeric-input-container">
                <input
                  type="number"
                  id="threshold"
                  min="0.5"
                  max="1.0"
                  step="0.01"
                  value={threshold}
                  onChange={(e) => {
                    // Allow any input during typing
                    if (e.target.value === '' || !isNaN(parseFloat(e.target.value))) {
                      setThreshold(e.target.value === '' ? 0 : parseFloat(e.target.value));
                    }
                  }}
                  onBlur={(e) => {
                    // Validate on blur
                    let value = parseFloat(e.target.value);
                    if (isNaN(value) || value < 0.5) {
                      value = 0.5;
                    } else if (value > 1.0) {
                      value = 1.0;
                    }
                    setThreshold(value);
                  }}
                  className="threshold-input"
                />
              </div>
              <div className="threshold-description">
                <small>Lower values (0.5) are more lenient, higher values (1.0) are stricter</small>
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
