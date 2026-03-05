import { defineStore } from 'pinia';
import { authAPI } from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: JSON.parse(localStorage.getItem('user') || 'null'),
  }),

  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    displayName: (state) => state.user?.display_name || state.user?.username || '用户',
  },

  actions: {
    setAuth(token, user) {
      this.token = token;
      this.user = user;
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
    },

    async login(username, password, rememberMe = false) {
      const response = await authAPI.login(username, password, rememberMe);
      this.setAuth(response.data.token, response.data.user);
      return response.data.user;
    },

    async fetchMe() {
      if (!this.token) return null;
      const response = await authAPI.me();
      this.user = response.data;
      localStorage.setItem('user', JSON.stringify(this.user));
      return this.user;
    },

    logout() {
      this.token = null;
      this.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
  },
});
