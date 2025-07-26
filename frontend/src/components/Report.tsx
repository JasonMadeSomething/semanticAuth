import { useState, useEffect } from 'react';
import authService from '../services/authService';

interface LoginAttempt {
  input: string;
  similarity: number;
  timestamp: string;
  passed: boolean;
}

const Report = () => {
  const [reportData, setReportData] = useState<LoginAttempt[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchReport();
  }, []);

  const [username, setUsername] = useState('');
  const [threshold, setThreshold] = useState(0.88);

  const fetchReport = async () => {
    if (!username) {
      setError('Please enter a username');
      return;
    }
    
    setIsLoading(true);
    setError('');
    
    try {
      const response = await authService.getReport(username, threshold);
      setReportData(response.data || []);
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Failed to fetch report');
      setReportData([]);
    } finally {
      setIsLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleString();
    } catch (e) {
      return dateString;
    }
  };

  return (
    <div className="report-container">
      <h2>Login Attempts Report</h2>
      
      <div className="report-form">
        <div className="form-group">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Enter username to view report"
            required
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="threshold">Threshold ({threshold})</label>
          <input
            type="range"
            id="threshold"
            min="0.5"
            max="1"
            step="0.01"
            value={threshold}
            onChange={(e) => setThreshold(parseFloat(e.target.value))}
          />
        </div>
        
        <button 
          onClick={fetchReport} 
          disabled={isLoading}
          className="refresh-button"
        >
          {isLoading ? 'Loading...' : 'Get Report'}
        </button>
      </div>
      
      {error && <div className="error-message">{error}</div>}
      
      {isLoading ? (
        <div className="loading">Loading report data...</div>
      ) : reportData.length > 0 ? (
        <div className="report-table-container">
          <table className="report-table">
            <thead>
              <tr>
                <th>Username</th>
                <th>Password Attempt</th>
                <th>Similarity</th>
                <th>Timestamp</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {reportData.map((attempt, index) => (
                <tr key={index} className={attempt.passed ? 'success-row' : 'failure-row'}>
                  <td>{username}</td>
                  <td>{attempt.input}</td>
                  <td>{attempt.similarity.toFixed(4)}</td>
                  <td>{formatDate(attempt.timestamp)}</td>
                  <td>{attempt.passed ? 'Success' : 'Failure'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="no-data">No login attempts found</div>
      )}
    </div>
  );
};

export default Report;
