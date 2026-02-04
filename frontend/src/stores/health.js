import { defineStore } from 'pinia';
import { fetchHealth } from '@/services/api';

const defaultState = () => ({
  user: {
    status: 'Unknown',
    message: 'Awaiting health probe.',
    checkedAt: 'Not checked',
    variant: 'neutral'
  },
  order: {
    status: 'Unknown',
    message: 'Awaiting health probe.',
    checkedAt: 'Not checked',
    variant: 'neutral'
  },
  auditLog: {
    status: 'Unknown',
    message: 'Awaiting health probe.',
    checkedAt: 'Not checked',
    variant: 'neutral'
  },
  gateway: {
    checkedAt: new Date().toLocaleTimeString()
  }
});

export const useHealthStore = defineStore('health', {
  state: defaultState,
  actions: {
    async refresh() {
      await Promise.all([
        this.fetchService('user'),
        this.fetchService('order'),
        this.fetchService('audit-log', 'auditLog')
      ]);
      this.gateway.checkedAt = new Date().toLocaleTimeString();
    },
    async fetchService(service, stateKey = service) {
      try {
        await fetchHealth(service);
        this[stateKey] = {
          status: 'Healthy',
          message: 'Service is responding within threshold.',
          checkedAt: new Date().toLocaleTimeString(),
          variant: 'positive'
        };
      } catch (error) {
        this[stateKey] = {
          status: 'Degraded',
          message: error?.message || 'Service did not respond to health check.',
          checkedAt: new Date().toLocaleTimeString(),
          variant: 'warning'
        };
      }
    }
  }
});
