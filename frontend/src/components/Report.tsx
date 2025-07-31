import { useState, useEffect, useCallback } from 'react';
import authService from '../services/authService';
import LoginAttemptsChart from './LoginAttemptsChart';

interface LoginAttempt {
  username?: string;
  input: string;
  similarity: number;
  timestamp: string;
  passed: boolean;
}

const Report = () => {
  const [threshold, setThreshold] = useState(0.88);
  const [reportData, setReportData] = useState<LoginAttempt[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [chartView, setChartView] = useState(true); // Default to chart view
  const [chartKey, setChartKey] = useState<number>(Date.now());

  // Define fetchReport with useCallback to prevent it from changing on every render
  const fetchReport = useCallback(async () => {
    setIsLoading(true);
    setError('');
    try {
      const response = await authService.getReport(threshold.toString());
      if (response.success) {
        setReportData(response.data || []);
      } else {
        setError(response.message || 'Failed to fetch report data');
        setReportData([]);
      }
    } catch (err) {
      setError('Error fetching report: ' + (err instanceof Error ? err.message : String(err)));
      setReportData([]);
    } finally {
      setIsLoading(false);
    }
  }, [threshold]);

  useEffect(() => {
    fetchReport();
  }, [fetchReport]);
  
  // Force chart refresh when report data changes
  useEffect(() => {
    if (reportData.length > 0) {
      setChartKey(Date.now());
    }
  }, [reportData]);

  // fetchReport is now defined above with useCallback

  return (
    <div className="report-container">
      <div className="report-card">
        <h2>Login Attempts Report</h2>
        
        <div className="report-controls">
          <div className="threshold-control">
            <label htmlFor="threshold">Similarity Threshold:</label>
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
            <div className="threshold-explanation">
              <p>* This value affects which attempts are shown as "Success" or "Failure" in the report (0.5-1.0)</p>
            </div>
          </div>
          
          <div className="view-toggle-switch">
            <span className={chartView ? 'active-label' : ''}>Chart</span>
            <label className="switch">
              <input 
                type="checkbox" 
                checked={!chartView}
                onChange={() => setChartView(!chartView)} 
              />
              <span className="slider round"></span>
            </label>
            <span className={!chartView ? 'active-label' : ''}>Table</span>
          </div>
          
          <button className="get-report-btn" onClick={fetchReport}>
            Get Report
          </button>
        </div>
        
        {error && <div className="error-message">{error}</div>}
        
        {isLoading ? (
          <p>Loading...</p>
        ) : reportData.length > 0 && chartView ? (
          // Use our new chart component
          <LoginAttemptsChart 
            key={chartKey}
            reportData={reportData}
            threshold={threshold}
          />
        ) : reportData.length > 0 ? (
          <div className="table-container">
            <table className="report-table">
              <thead>
                <tr>
                  <th>Password Attempt</th>
                  <th>Similarity</th>
                  <th>Timestamp</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {reportData.map((attempt, index) => (
                  <tr 
                    key={index} 
                    className={attempt.passed ? 'success-row' : 'failure-row'}
                  >
                    <td>{attempt.input}</td>
                    <td>{attempt.similarity.toFixed(4)}</td>
                    <td>{attempt.timestamp}</td>
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
    </div>
  );
};

export default Report;
