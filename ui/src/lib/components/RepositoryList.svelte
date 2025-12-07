<script lang="ts">
    import { onMount } from "svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";

    interface Repository {
        id: number;
        name: string;
        full_name: string;
        html_url: string;
        description: string;
        private: boolean;
    }

    let repositories = $state<Repository[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        try {
            const res = await fetch("/api/user/repos");
            if (!res.ok) {
                throw new Error("Failed to fetch repositories");
            }
            repositories = await res.json();
        } catch (e) {
            error = (e as Error).message;
        } finally {
            loading = false;
        }
    });

    function navigateToRepo(repoName: string) {
        window.location.hash = `#/repo/${repoName}`;
    }
</script>

<Card.Root class="w-full mt-4">
    <Card.Header>
        <Card.Title>Repositories</Card.Title>
        <Card.Description
            >Select a repository to manage secrets.</Card.Description
        >
    </Card.Header>
    <Card.Content>
        {#if loading}
            <div class="text-center py-4">Loading repositories...</div>
        {:else if error}
            <div class="text-destructive text-center py-4">{error}</div>
        {:else if repositories.length === 0}
            <div class="text-center py-4 text-muted-foreground">
                No repositories found.
            </div>
        {:else}
            <div class="flex flex-col gap-2">
                {#each repositories as repo}
                    <!-- svelte-ignore a11y_click_events_have_key_events -->
                    <!-- svelte-ignore a11y_no_static_element_interactions -->
                    <div
                        class="flex items-center justify-between p-3 border rounded-md hover:bg-muted/50 transition-colors cursor-pointer group"
                        onclick={() => navigateToRepo(repo.full_name)}
                    >
                        <div class="flex flex-col gap-1 overflow-hidden">
                            <span
                                class="font-medium truncate flex items-center gap-2"
                            >
                                {repo.name}
                                {#if repo.private}
                                    <span
                                        class="text-[10px] uppercase border px-1 rounded text-muted-foreground"
                                        >Private</span
                                    >
                                {/if}
                            </span>
                            <span
                                class="text-xs text-muted-foreground break-words"
                                >{repo.description || "No description"}</span
                            >
                        </div>
                        <Button variant="accent" size="sm" class=""
                            >Manage</Button
                        >
                    </div>
                {/each}
            </div>
        {/if}
    </Card.Content>
</Card.Root>
