const API_BASE_URL = 'http://127.0.0.1:8080/api';

// Helper function to get auth token from localStorage
const getToken = () => {
  return localStorage.getItem('token');
};

// Helper function to set auth token
const setToken = (token) => {
  localStorage.setItem('token', token);
};

// Helper function to remove auth token
const removeToken = () => {
  localStorage.removeItem('token');
};

// Helper function to get auth headers
const getAuthHeaders = () => {
  const token = getToken();
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

// API request helper
const request = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    ...options,
    headers: {
      ...getAuthHeaders(),
      ...options.headers,
    },
  };

  try {
    const response = await fetch(url, config);
    const contentType = response.headers.get('content-type') || '';
    const rawBody = await response.text();

    const data =
      rawBody && contentType.includes('application/json')
        ? JSON.parse(rawBody)
        : rawBody;

    if (!response.ok) {
      const message =
        (data && typeof data === 'object' && data.error) ||
        rawBody ||
        response.statusText ||
        'Request failed';
      throw new Error(message);
    }

    return data;
  } catch (error) {
    throw error;
  }
};

// Auth API
export const authAPI = {
  register: async (email, password, name) => {
    const data = await request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, name }),
    });
    if (data.token) {
      setToken(data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    }
    return data;
  },

  login: async (email, password) => {
    const data = await request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    if (data.token) {
      setToken(data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    }
    return data;
  },

  logout: () => {
    removeToken();
    localStorage.removeItem('user');
  },

  getProfile: async () => {
    return request('/profile');
  },
};

// Hotel API
export const hotelAPI = {
  getCities: async () => {
    const data = await request('/cities');
    return data.cities || [];
  },

  getAllHotels: async (sort = 'city', order = 'asc') => {
    return request(`/hotels?sort=${sort}&order=${order}`);
  },

  getHotelsByCity: async (city, sort = 'price', order = 'asc') => {
    return request(`/hotels/city/${encodeURIComponent(city)}?sort=${sort}&order=${order}`);
  },

  getPriceHistory: async (hotelId) => {
    return request(`/hotels/${hotelId}/history`);
  },
};

// Admin API
export const adminAPI = {
  triggerScrape: async () => {
    return request('/scrape/trigger', {
      method: 'POST',
    });
  },

  getScrapingLogs: async () => {
    return request('/logs');
  },

  getAllUsers: async () => {
    return request('/admin/users');
  },
};

export { getToken, setToken, removeToken };

