import type { PageLoad } from "./$types";
import { error, redirect } from "@sveltejs/kit";

export const load: PageLoad = async ({ fetch }) => {
    const res = await fetch("/api/user/repos");

    if (res.status === 401 || res.status === 403) {
        throw redirect(302, "/login?unauthorized=1");
    }

    if (!res.ok) {
        throw error(res.status, "Failed to fetch repositories");
    }

    const repositories = await res.json();
    return { repositories };
};
