import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import RepositorySecrets from './RepositorySecrets.svelte';

describe('RepositorySecrets', () => {
    const owner = 'test-owner';
    const repo = 'test-repo';

    beforeEach(() => {
        global.fetch = vi.fn();
    });

    it('shows loading state initially', () => {
        (global.fetch as any).mockReturnValue(new Promise(() => { })); // Never resolves
        render(RepositorySecrets, { owner, repo });
        expect(screen.getByText('Loading secrets...')).toBeInTheDocument();
    });

    it('displays secrets list on success', async () => {
        (global.fetch as any).mockResolvedValue({
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
        (global.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => [],
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('No secrets found for this repository.')).toBeInTheDocument();
        });
    });

    it('handles null response from backend gracefully', async () => {
        (global.fetch as any).mockResolvedValue({
            ok: true,
            json: async () => null,
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('No secrets found for this repository.')).toBeInTheDocument();
        });
    });

    it('displays error message on fetch failure', async () => {
        (global.fetch as any).mockResolvedValue({
            ok: false,
        });

        render(RepositorySecrets, { owner, repo });

        await waitFor(() => {
            expect(screen.getByText('Failed to fetch secrets')).toBeInTheDocument();
        });
    });
});
