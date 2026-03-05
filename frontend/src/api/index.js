import axios from 'axios';

const DEFAULT_API_BASE_URL = '/api';
const runtimeConfig = typeof window !== 'undefined' ? window.__APP_CONFIG__ : undefined;
const configuredBaseURL =
  typeof runtimeConfig?.apiBaseUrl === 'string' ? runtimeConfig.apiBaseUrl.trim() : '';

const apiClient = axios.create({
  baseURL: configuredBaseURL || DEFAULT_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (username, password, rememberMe = false) =>
    apiClient.post('/auth/login', { username, password, remember_me: rememberMe }),
  me: () => apiClient.get('/auth/me'),
};

export const keySpecAPI = {
  list: () => apiClient.get('/key-specs'),
  get: (id) => apiClient.get(`/key-specs/${id}`),
  create: (data) => apiClient.post('/key-specs', data),
  update: (id, data) => apiClient.put(`/key-specs/${id}`, data),
  delete: (id) => apiClient.delete(`/key-specs/${id}`),
  reorder: (ids) => apiClient.put('/key-specs/reorder', { ids }),
};

export const keyAPI = {
  list: (specId, { onlyUnused = false, limit = 50, offset = 0 } = {}) =>
    apiClient.get('/keys', {
      params: { spec_id: specId, only_unused: onlyUnused, limit, offset },
    }),
  batchUpload: (specId, keys) => apiClient.post('/keys/batch', { spec_id: specId, keys }),
  getAvailable: (specId) => apiClient.get(`/keys/available/${specId}`),
  markUsed: (id) => apiClient.put(`/keys/${id}/use`),
  delete: (id) => apiClient.delete(`/keys/${id}`),
};

export const configAPI = {
  getTemplate: () => apiClient.get('/config/copy-template'),
  updateTemplate: (template) => apiClient.put('/config/copy-template', { template }),
};

export const ledgerAPI = {
  getCosts: ({ page = 1, limit = 20 } = {}) =>
    apiClient.get('/ledger/costs', { params: { page, limit } }),
  createCost: (data) => apiClient.post('/ledger/costs', data),
  deleteCost: (id) => apiClient.delete(`/ledger/costs/${id}`),

  getCustomers: (search = '') =>
    apiClient.get('/ledger/customers', { params: search ? { search } : {} }),
  createCustomer: (data) => apiClient.post('/ledger/customers', data),
  updateCustomer: (id, data) => apiClient.put(`/ledger/customers/${id}`, data),

  getTransactions: ({ page = 1, limit = 20, customerId, isNewCustomer } = {}) => {
    const params = { page, limit };
    if (customerId) params.customer_id = customerId;
    if (typeof isNewCustomer === 'boolean') {
      params.is_new_customer = isNewCustomer ? 1 : 0;
    }
    return apiClient.get('/ledger/transactions', { params });
  },
  createTransaction: (data) => apiClient.post('/ledger/transactions', data),

  getStatistics: () => apiClient.get('/ledger/statistics'),
};

export default apiClient;
