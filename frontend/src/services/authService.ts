// API service for authentication
// Use environment variables for API URL with fallback to localhost for development
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

interface RegisterRequest {
  username: string;
  password: string;
}

interface LoginRequest {
  username: string;
  password: string;
  threshold?: number;
}

// Standard API response structure
interface ApiResponse<T> {
  status: string;
  success: boolean;
  message?: string;
  data?: T;
}

// Login response data structure
interface LoginResponseData {
  username: string;
  similarity: number;
  threshold: number;
}

// Report data structure
interface ReportItem {
  input: string;
  similarity: number;
  timestamp: string;
  passed: boolean;
}

export const authService = {
  // Register a new user
  async register(username: string, password: string): Promise<ApiResponse<{username: string}>> {
    try {
      const registerData: RegisterRequest = { username, password };
      const response = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(registerData),
      });
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return {
          status: 'error',
          success: false,
          message: errorData.message || `Server error: ${response.status}`,
        };
      }

      const data = await response.json();
      
      if (!data.success) {
        return {
          status: 'error',
          success: false,
          message: data.message || 'Registration failed',
        };
      }
      
      return data;
    } catch (error) {
      return {
        status: 'error',
        success: false,
        message: error instanceof Error ? error.message : 'Network error during registration',
      };
    }
  },

  // Login a user
  async login(username: string, password: string, threshold?: number): Promise<ApiResponse<LoginResponseData>> {
    try {
      const loginData: LoginRequest = { username, password };
      if (threshold) {
        loginData.threshold = threshold;
      }

      const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginData),
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return {
          status: 'error',
          success: false,
          message: errorData.message || `Server error: ${response.status}`,
        };
      }

      const data = await response.json();
      
      if (!data.success) {
        return {
          status: 'error',
          success: false,
          message: data.message || 'Login failed',
        };
      }

      return data;
    } catch (error) {
      return {
        status: 'error',
        success: false,
        message: error instanceof Error ? error.message : 'Network error during login',
      };
    }
  },

  // Get report data
  async getReport(threshold?: string): Promise<ApiResponse<ReportItem[]>> {
    try {
      // Build query parameters
      const params = new URLSearchParams();
      if (threshold !== undefined) {
        params.append('threshold', threshold);
      }
      
      const queryString = params.toString() ? `?${params.toString()}` : '';
      const response = await fetch(`${API_URL}/report${queryString}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        return {
          status: 'error',
          success: false,
          message: errorData.message || `Server error: ${response.status}`,
        };
      }

      const data = await response.json();
      
      if (!data.success) {
        return {
          status: 'error',
          success: false,
          message: data.message || 'Failed to fetch report data',
        };
      }

      return data;
    } catch (error) {
      return {
        status: 'error',
        success: false,
        message: error instanceof Error ? error.message : 'Network error while fetching report data',
        data: [],
      };
    }
  }
};

export default authService;
