import type { PageLoad } from "./$types";
import { error, redirect } from "@sveltejs/kit";

export const load: PageLoad = async ({ fetch, params }) => {
    const { owner, repo } = params;

    const [secretsRes, csrfRes] = await Promise.all([
        fetch(`/api/repo/${owner}/${repo}/secrets`),
        fetch("/api/csrf-token"),
    ]);

    if (secretsRes.status === 401 || csrfRes.status === 401) {
        throw redirect(302, "/login?unauthorized=1");
    }

    if (!secretsRes.ok) {
        throw error(secretsRes.status, "Failed to fetch secrets");
    }

    const secrets = (await secretsRes.json()) || [];
    let csrfToken: string | null = null;
    if (csrfRes.ok) {
        const data = await csrfRes.json();
        csrfToken = data.token;
    }

    return { owner, repo, secrets, csrfToken };
};
