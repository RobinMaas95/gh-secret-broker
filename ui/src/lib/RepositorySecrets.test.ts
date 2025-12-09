import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import userEvent from '@testing-library/user-event';
import RepositorySecrets from './RepositorySecrets.svelte';

describe('RepositorySecrets', () => {
    const owner = 'test-owner';
    const repo = 'test-repo';

    beforeEach(() => {
        window.fetch = vi.fn();
    });

    it('shows loading state initially', () => {
        (window.fetch as any).mockReturnValue(new Promise(() => { })); // Never resolves
        render(RepositorySecrets, { owner, repo });
        expect(screen.getByText('Loading secrets...')).toBeInTheDocument();
    });

    it('displays secrets list on success', async () => {
        (window.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => ['SECRET_1', 'SECRET_2'],
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('SECRET_1')).toBeInTheDocument();
            expect(screen.getByText('SECRET_2')).toBeInTheDocument();
        });
    });

    it('shows empty state when no secrets found', async () => {
        (window.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => [],
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('No secrets found for this repository.')).toBeInTheDocument();
        });
    });

    it('handles null response from backend gracefully', async () => {
        (window.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => null,
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('No secrets found for this repository.')).toBeInTheDocument();
        });
    });

    it('displays error message on fetch failure', async () => {
        (window.fetch as any).mockResolvedValue({
            ok: false,
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('Failed to fetch secrets')).toBeInTheDocument();
        });
    });

    it('deletes a secret successfully', async () => {
        // Mock initial fetch and CSRF token
        (window.fetch as any).mockImplementation((url: string, options: any) => {
            if (url.endsWith('/secrets')) {
                return Promise.resolve({
                    ok: true,
                    json: async () => ['SECRET_TO_DELETE'],
                });
            }
            if (url.endsWith('/csrf-token')) {
                return Promise.resolve({
                    ok: true,
                    json: async () => ({ token: 'mock-csrf-token' }),
                });
            }
            if (url.includes('/secrets/SECRET_TO_DELETE') && options?.method === 'DELETE') {
                return Promise.resolve({
                    ok: true,
                });
            }
            return Promise.reject(new Error(`Unknown URL: ${url}`));
        });

        const user = userEvent.setup();
        render(RepositorySecrets, { owner, repo });

        // Wait for secret to appear
        await waitFor(() => {
            expect(screen.getByText('SECRET_TO_DELETE')).toBeInTheDocument();
        });

        // Click Delete button
        const deleteBtn = screen.getByRole('button', { name: 'Delete' });
        await user.click(deleteBtn);

        // Check if dialog appears
        expect(await screen.findByText('Are you absolutely sure?')).toBeInTheDocument();

        // Click Continue
        const continueBtn = screen.getByRole('button', { name: 'Continue' });
        await user.click(continueBtn);

        // Check fetch was called with DELETE
        await waitFor(() => {
            expect(window.fetch).toHaveBeenCalledWith(
                expect.stringContaining(`/api/repo/${owner}/${repo}/secrets/SECRET_TO_DELETE`),
                expect.objectContaining({
                    method: 'DELETE',
                    headers: expect.objectContaining({
                        'X-CSRF-Token': 'mock-csrf-token',
                    }),
                })
            );
        });

        // Check secret is removed from list
        await waitFor(() => {
            expect(screen.queryByText('SECRET_TO_DELETE')).not.toBeInTheDocument();
        });
    });
});
