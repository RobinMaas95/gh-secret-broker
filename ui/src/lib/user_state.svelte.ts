interface User {
    AvatarURL: string;
    Name: string;
    NickName: string;
    Email: string;
    Location: string;
    Description: string;
    UserID: string;
    Provider: string;
}

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
    };
}

export const userState = createUserState();
