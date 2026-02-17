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

  it('renders landing message and purpose sections', () => {
    const wrapper = mount(Home, {
      global: {
        stubs: {
          RouterLink: true,
          RouterView: true
        }
      }
    });

    expect(wrapper.text()).toContain('enterprise-ready foundation');
    expect(wrapper.text()).toContain('Project Purpose');
    expect(wrapper.text()).toContain('Unified visibility');
    expect(wrapper.text()).toContain('Backend in one command');
  });
});
