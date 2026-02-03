import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Home from '@/views/Home.vue';

vi.mock('@/services/api', () => ({
  fetchHealth: vi.fn().mockResolvedValue({ status: 'healthy' })
}));

describe('Home view', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('renders headline and status cards', () => {
    const wrapper = mount(Home, {
      global: {
        stubs: {
          RouterLink: true,
          RouterView: true
        }
      }
    });

    expect(wrapper.text()).toContain('Enterprise Microservice Portal');
    expect(wrapper.text()).toContain('User Service');
    expect(wrapper.text()).toContain('Order Service');
  });
});
