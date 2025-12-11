
import { readable } from 'svelte/store';

// $app/environment
export const browser = true;
export const dev = true;
export const building = false;
export const version = 'test';

// $app/stores
export const page = readable({
    url: new URL('http://localhost'),
    params: {},
    route: { id: 'test' },
    status: 200,
    error: null,
    data: {},
    form: null
});
export const navigating = readable(null);
export const updated = readable(false);

// $app/navigation
export const goto = () => Promise.resolve();
export const invalidate = () => Promise.resolve();
export const invalidateAll = () => Promise.resolve();
export const preloadData = () => Promise.resolve();
export const preloadCode = () => Promise.resolve();
export const beforeNavigate = () => { };
export const afterNavigate = () => { };

// $app/forms
export const applyAction = () => { };
export const deserialize = () => { };
export const enhance = () => { };
