// API service for authentication
const API_URL = 'http://localhost:8080';

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
    const registerData: RegisterRequest = { username, password };
    const response = await fetch(`${API_URL}/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(registerData),
    });

    const data = await response.json();
    
    if (!response.ok || !data.success) {
      throw new Error(data.message || 'Registration failed');
    }

    return data;
  },

  // Login a user
  async login(username: string, password: string, threshold?: number): Promise<ApiResponse<LoginResponseData>> {
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

    const data = await response.json();
    
    if (!response.ok || !data.success) {
      throw new Error(data.message || 'Login failed');
    }

    return data;
  },

  // Get report data
  async getReport(threshold?: string): Promise<ApiResponse<ReportItem[]>> {
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

    const data = await response.json();
    
    if (!response.ok || !data.success) {
      throw new Error(data.message || 'Failed to get report');
    }

    return data;
  }
};

export default authService;
