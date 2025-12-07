import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createUserState } from './user_state.svelte';
import { tick } from 'svelte';

// Mock fetch
const fetchMock = vi.fn();
global.fetch = fetchMock;

describe('UserState', () => {
    beforeEach(() => {
        fetchMock.mockReset();
    });

    it('should initialize with null user', () => {
        const userState = createUserState();
        expect(userState.current).toBeNull();
        expect(userState.loading).toBe(false); // Assuming it starts loading or we have a loading state
    });

    it('should load user successfully', async () => {
        const mockUser = {
            Name: 'Test User',
            AvatarURL: 'http://example.com/avatar.png',
            UserID: '123',
            Provider: 'github'
        };

        fetchMock.mockResolvedValueOnce({
            ok: true,
            json: async () => mockUser
        });

        const userState = createUserState();
        await userState.load();

        expect(userState.current).toEqual(mockUser);
        expect(userState.error).toBeNull();
        expect(userState.loading).toBe(false);
    });

    it('should handle unauthorized error', async () => {
        fetchMock.mockResolvedValueOnce({
            ok: false,
            status: 401
        });

        const userState = createUserState();
        await userState.load();

        expect(userState.current).toBeNull();
        expect(userState.error).toBe('Unauthorized');
        expect(userState.loading).toBe(false);
    });
});
