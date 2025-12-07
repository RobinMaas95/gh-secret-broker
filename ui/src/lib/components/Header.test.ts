import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import Header from './Header.svelte';
import { userState } from '../user_state.svelte';

// Mock userState
vi.mock('../user_state.svelte', () => {
    return {
        userState: {
            current: null,
            loading: false
        }
    };
});

vi.mock('mode-watcher', () => {
    return {
        toggleMode: vi.fn(),
        setMode: vi.fn(),
        resetMode: vi.fn(),
        mode: { subscribe: (fn: any) => { fn('light'); return () => { }; } }
    };
});

// ModeToggle mock removed to use actual component


describe('Header', () => {
    it('should render title', () => {
        render(Header);
        expect(screen.getByText('GH Secret Broker')).toBeInTheDocument();
    });
});
