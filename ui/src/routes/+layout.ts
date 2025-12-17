import type { LayoutLoad } from "./$types";
import { error, redirect } from "@sveltejs/kit";

export const load: LayoutLoad = async ({ fetch, url }) => {
    // Skip user fetch on public login route to avoid redirect loops
    if (url.pathname.startsWith("/login")) {
        return { user: null };
    }

    const res = await fetch("/api/user");

    if (res.status === 401 || res.status === 403) {
        throw redirect(302, "/login?unauthorized=1");
    }

    if (!res.ok) {
        throw error(res.status, "Failed to load user");
    }

    const user = await res.json();
    return { user };
};
