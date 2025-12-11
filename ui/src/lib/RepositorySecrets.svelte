<script lang="ts">
    import { onMount } from "svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button, buttonVariants } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import AddSecretDialog from "$lib/components/AddSecretDialog.svelte";

    let { owner, repo } = $props<{ owner: string; repo: string }>();

    let secrets = $state<string[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);

    let csrfToken = $state<string | null>(null);

    async function fetchSecrets() {
        loading = true;
        error = null;
        try {
            const res = await fetch(`/api/repo/${owner}/${repo}/secrets`);
            if (!res.ok) {
                throw new Error("Failed to fetch secrets");
            }
            secrets = (await res.json()) || [];
        } catch (e) {
            error = (e as Error).message;
        } finally {
            loading = false;
        }
    }

    async function fetchCsrfToken() {
        try {
            const res = await fetch("/api/csrf-token");
            if (res.ok) {
                const data = await res.json();
                csrfToken = data.token;
            }
        } catch (e) {
            console.error("Failed to fetch CSRF token", e);
        }
    }

    async function deleteSecret(name: string) {
        if (!csrfToken) {
            error = "CSRF token missing";
            return;
        }

        try {
            const res = await fetch(
                `/api/repo/${owner}/${repo}/secrets/${name}`,
                {
                    method: "DELETE",
                    headers: {
                        "X-CSRF-Token": csrfToken,
                    },
                },
            );

            if (!res.ok) {
                throw new Error("Failed to delete secret");
            }

            // Remove from list immediately
            secrets = secrets.filter((s) => s !== name);
        } catch (e) {
            error = (e as Error).message;
        }
    }

    onMount(() => {
        fetchSecrets();
        fetchCsrfToken();
    });
</script>

<div class="container mx-auto py-8 max-w-4xl">
    <div class="mb-4">
        <a
            href="#/dashboard"
            class="text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
        >
            ← Back to Dashboard
        </a>
    </div>

    <Card.Root>
        <Card.Header>
            <div class="flex justify-between items-center">
                <Card.Title class="text-2xl"
                    >Secrets for {owner}/{repo}</Card.Title
                >
                <AddSecretDialog
                    {owner}
                    {repo}
                    onSecretAdded={(name: string) => {
                        if (!secrets.includes(name)) {
                            secrets = [...secrets, name];
                        }
                    }}
                />
            </div>
            <Card.Description
                >Manage secrets for this repository.</Card.Description
            >
        </Card.Header>
        <Card.Content>
            {#if loading}
                <div class="text-center py-8 text-muted-foreground">
                    Loading secrets...
                </div>
            {:else if error}
                <div class="text-destructive text-center py-8">{error}</div>
            {:else if secrets.length === 0}
                <div class="text-center py-8 text-muted-foreground">
                    No secrets found for this repository.
                </div>
            {:else}
                <div class="grid gap-2">
                    {#each secrets as secret}
                        <div
                            class="flex items-center justify-between p-3 border rounded-md"
                        >
                            <span class="font-mono">{secret}</span>
                            <div class="flex gap-2">
                                <AlertDialog.Root>
                                    <AlertDialog.Trigger
                                        class={buttonVariants({
                                            variant: "destructive",
                                            size: "sm",
                                        })}
                                        disabled={!csrfToken}
                                    >
                                        Delete
                                    </AlertDialog.Trigger>
                                    <AlertDialog.Content>
                                        <AlertDialog.Header>
                                            <AlertDialog.Title
                                                >Are you absolutely sure?</AlertDialog.Title
                                            >
                                            <AlertDialog.Description>
                                                This action cannot be undone.
                                                This will permanently delete the
                                                secret
                                                <span
                                                    class="font-mono font-bold"
                                                    >{secret}</span
                                                >
                                                from the repository.
                                            </AlertDialog.Description>
                                        </AlertDialog.Header>
                                        <AlertDialog.Footer>
                                            <AlertDialog.Cancel
                                                >Cancel</AlertDialog.Cancel
                                            >
                                            <AlertDialog.Action
                                                class={buttonVariants({
                                                    variant: "destructive",
                                                })}
                                                onclick={() =>
                                                    deleteSecret(secret)}
                                                >Continue</AlertDialog.Action
                                            >
                                        </AlertDialog.Footer>
                                    </AlertDialog.Content>
                                </AlertDialog.Root>
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </Card.Content>
    </Card.Root>
</div>
