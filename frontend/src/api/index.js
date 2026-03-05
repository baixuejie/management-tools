import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (username, password) => apiClient.post('/auth/login', { username, password }),
  logout: () => apiClient.post('/auth/logout'),
};

export const keySpecAPI = {
  list: () => apiClient.get('/key-specs'),
  create: (data) => apiClient.post('/key-specs', data),
  update: (id, data) => apiClient.put(`/key-specs/${id}`, data),
  delete: (id) => apiClient.delete(`/key-specs/${id}`),
};

export const keyAPI = {
  list: (specId) => apiClient.get(`/keys?spec_id=${specId}`),
  create: (data) => apiClient.post('/keys', data),
  get: (id) => apiClient.get(`/keys/${id}`),
  delete: (id) => apiClient.delete(`/keys/${id}`),
};

export const configAPI = {
  get: () => apiClient.get('/config'),
  update: (data) => apiClient.put('/config', data),
};

export default apiClient;
