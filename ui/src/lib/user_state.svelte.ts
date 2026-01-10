import type { User } from "$lib/types";

export function createUserState() {
    let current = $state<User | null>(null);
    let loading = $state(false);
    let error = $state<string | null>(null);

    async function load() {
        loading = true;
        error = null;
        try {
            const res = await fetch("/api/user");
            if (!res.ok) {
                throw new Error("Unauthorized");
            }
            current = await res.json();
        } catch (e) {
            error = (e as Error).message;
            current = null;
        } finally {
            loading = false;
        }
    }

    return {
        get current() { return current },
        get loading() { return loading },
        get error() { return error },
        load,
        setLoaded(user: User | null, err: string | null = null) {
            current = user;
            error = err;
            loading = false;
        },
    };
}

export const userState = createUserState();
