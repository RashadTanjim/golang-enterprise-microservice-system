import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fetchHealth, listUsers, ApiError } from '@/services/api';

const mockResponse = (overrides = {}) => ({
  ok: true,
  status: 200,
  headers: {
    get: () => 'application/json'
  },
  json: async () => ({ success: true, data: [] }),
  ...overrides
});

describe('api service', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  it('rejects invalid health service names', async () => {
    await expect(fetchHealth('invalid')).rejects.toBeInstanceOf(ApiError);
  });

  it('requires a token for listUsers', async () => {
    await expect(listUsers('')).rejects.toBeInstanceOf(ApiError);
  });

  it('propagates api error payloads', async () => {
    global.fetch.mockResolvedValue(
      mockResponse({
        ok: false,
        status: 400,
        json: async () => ({ error: { message: 'Bad request', code: 'BAD_REQUEST' } })
      })
    );

    await expect(listUsers('token')).rejects.toMatchObject({
      message: 'Bad request',
      code: 'BAD_REQUEST',
      status: 400
    });
  });
});
