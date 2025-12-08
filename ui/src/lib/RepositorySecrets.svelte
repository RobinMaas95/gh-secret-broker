<script lang="ts">
    import { onMount } from "svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";

    let { owner, repo } = $props<{ owner: string; repo: string }>();

    let secrets = $state<string[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);

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

    onMount(() => {
        fetchSecrets();
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
            <Card.Title class="text-2xl">Secrets for {owner}/{repo}</Card.Title>
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
                                <Button variant="destructive" size="sm" disabled
                                    >Delete</Button
                                >
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </Card.Content>
    </Card.Root>
</div>
