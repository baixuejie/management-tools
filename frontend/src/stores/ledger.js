import { defineStore } from 'pinia';
import { ledgerAPI } from '../api';

export const useLedgerStore = defineStore('ledger', {
  state: () => ({
    statistics: {
      total_cost: 0,
      total_revenue: 0,
      total_commission: 0,
      net_profit: 0,
      new_customers: 0,
      renewal_customers: 0,
    },
    transactions: [],
    transactionTotal: 0,
    customers: [],
  }),

  actions: {
    async fetchStatistics() {
      const response = await ledgerAPI.getStatistics();
      this.statistics = response.data;
      return this.statistics;
    },

    async fetchTransactions(params = {}) {
      const response = await ledgerAPI.getTransactions(params);
      this.transactions = response.data.items || [];
      this.transactionTotal = response.data.total || 0;
      return response.data;
    },

    async fetchCustomers(search = '') {
      const response = await ledgerAPI.getCustomers(search);
      this.customers = response.data.items || [];
      return this.customers;
    },

    async createTransaction(data) {
      const response = await ledgerAPI.createTransaction(data);
      return response.data;
    },

    async createCost(data) {
      const response = await ledgerAPI.createCost(data);
      return response.data;
    },

    async fetchCosts(params = {}) {
      const response = await ledgerAPI.getCosts(params);
      return response.data;
    },

    async deleteCost(id) {
      const response = await ledgerAPI.deleteCost(id);
      return response.data;
    },

    async updateCustomer(id, data) {
      const response = await ledgerAPI.updateCustomer(id, data);
      return response.data;
    },
  },
});
