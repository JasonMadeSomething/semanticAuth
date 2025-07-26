import { useRef, useEffect } from 'react';
import { Scatter } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  LinearScale,
  TimeScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend
} from 'chart.js';
import 'chartjs-adapter-date-fns';

// Register Chart.js components
ChartJS.register(
  LinearScale,
  TimeScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend
);

// Define the LoginAttempt interface
interface LoginAttempt {
  username?: string;
  input: string;
  similarity: number;
  timestamp: string;
  passed: boolean;
}

// Props for the chart component
interface LoginAttemptsChartProps {
  reportData: LoginAttempt[];
  threshold: number;
}

const LoginAttemptsChart = ({ reportData, threshold }: LoginAttemptsChartProps) => {
  const chartRef = useRef<ChartJS | null>(null);
  
  // Clean up chart instance when component unmounts
  useEffect(() => {
    return () => {
      if (chartRef.current) {
        chartRef.current.destroy();
      }
    };
  }, []);

  // Prepare data for the scatter plot
  const prepareChartData = () => {
    if (!reportData || reportData.length === 0) {
      return { datasets: [] };
    }

    // Sort by timestamp and take the most recent 100 attempts
    const recentAttempts = [...reportData]
      .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime())
      .slice(-100);

    // Define chart point type
    interface ChartPoint {
      x: number;
      y: number;
      similarity: number;
      input: string;
      timestamp: string;
      passed: boolean;
    }
    
    // Separate successful and failed attempts
    const successData: ChartPoint[] = [];
    const failData: ChartPoint[] = [];
    
    // Create data points with sequential x values
    recentAttempts.forEach((attempt, index) => {
      // Parse and clamp similarity value between 0 and 1
      let similarity = parseFloat(attempt.similarity.toString());
      if (isNaN(similarity)) {
        console.warn("Invalid similarity value", attempt.similarity, attempt);
        similarity = 0; // Default to 0 for invalid values
      } else {
        // Clamp between 0 and 1 to handle floating-point precision issues
        similarity = Math.max(0, Math.min(1, similarity));
      }
      
      const point = {
        x: index + 1, // Use sequential integers for x-axis
        y: similarity,
        similarity,
        input: attempt.input,
        timestamp: attempt.timestamp,
        passed: attempt.passed
      };
      
      if (attempt.passed) {
        successData.push(point);
      } else {
        failData.push(point);
      }
    });
    
    // Create datasets
    const datasets = [
      {
        label: 'Successful Attempts',
        data: successData,
        backgroundColor: 'rgba(75, 192, 192, 0.6)',
        pointRadius: 4, // Smaller point radius to prevent Y-axis expansion
        pointStyle: 'circle',
        borderColor: 'rgba(75, 192, 192, 1)',
        borderWidth: 1
      },
      {
        label: 'Failed Attempts',
        data: failData,
        backgroundColor: 'rgba(255, 99, 132, 0.5)',
        pointRadius: 4, // Smaller point radius to prevent Y-axis expansion
        pointStyle: 'triangle',
        borderColor: 'rgba(255, 99, 132, 1)',
        borderWidth: 1
      }
    ];
    
    // Add threshold line if we have data
    if (recentAttempts.length > 0) {
      const thresholdValue = parseFloat(threshold.toString());
      
      // Use a different type for threshold line
      type ThresholdPoint = { x: number; y: number };
      
      datasets.push({
        label: `Threshold (${threshold})`,
        data: [
          { x: 0, y: thresholdValue } as ThresholdPoint,  // Start before first point
          { x: recentAttempts.length + 1, y: thresholdValue } as ThresholdPoint  // End after last point
        ],
        backgroundColor: 'rgba(0, 0, 0, 0.7)',
        borderColor: 'rgba(0, 0, 0, 0.7)',
        borderWidth: 2,
        borderDash: [5, 5],
        pointRadius: 0,
        type: 'line' as const
      } as any); // Use type assertion for the threshold dataset
    }
    
    // Add debugging to verify counts
    console.log("Success:", successData.length, "Fail:", failData.length, "Total:", successData.length + failData.length);
    
    return { datasets };
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'linear' as const,
        position: 'bottom' as const,
        title: {
          display: true,
          text: 'Login Attempts (Time Sequence)'
        },
        ticks: {
          stepSize: 1,
          precision: 0,  // Force whole numbers
          callback: function(value: any) {
            // Only show integer values
            if (Number.isInteger(value)) {
              return value;
            }
            return '';
          }
        },
        min: 1, // Start at first point (inclusive)
        max: reportData?.length || 10 // End at last point (inclusive)
      },
      y: {
        type: 'linear' as const,
        position: 'left' as const,
        title: {
          display: true,
          text: 'Similarity'
        },
        beginAtZero: true,
        min: 0,
        max: 1,
        ticks: {
          stepSize: 0.1,
          precision: 1,
          // Force specific ticks
          callback: function(value: any) {
            return value.toFixed(1);
          }
        },
        // Disable auto-scaling
        grace: 0
      }
    },
    plugins: {
      tooltip: {
        callbacks: {
          label: function(context: any) {
            const dataPoint = context.raw;
            return [
              `Password Input: ${dataPoint.input || 'N/A'}`,
              `Similarity: ${typeof dataPoint.similarity === 'number' ? dataPoint.similarity.toFixed(4) : 'N/A'}`,
              `Status: ${dataPoint.passed ? 'Success' : 'Failure'}`,
              `Time: ${dataPoint.timestamp}`
            ];
          }
        }
      },
      legend: {
        position: 'top' as const
      }
    }
  };

  return (
    <div className="chart-container" style={{ height: '400px', maxHeight: '400px', position: 'relative' }}>
      <Scatter
        data={prepareChartData()}
        options={chartOptions}
        ref={(reference) => {
          if (reference) {
            chartRef.current = reference;
          }
        }}
      />
      <div className="chart-explanation">
        <h3>Chart Explanation</h3>
        <ul>
          <li>Each point represents a login attempt</li>
          <li>X-axis: Time of the attempt</li>
          <li>Y-axis: Similarity score (higher is better)</li>
          <li>Circles: Successful logins</li>
          <li>Triangles: Failed logins</li>
          <li>The threshold line shows the current similarity threshold ({threshold})</li>
        </ul>
      </div>
    </div>
  );
};

export default LoginAttemptsChart;
