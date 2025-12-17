<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import { goto } from "$app/navigation";

    interface Repository {
        id: number;
        name: string;
        full_name: string;
        html_url: string;
        description: string;
        private: boolean;
    }

    let { repositories } = $props<{ repositories: Repository[] }>();

    function navigateToRepo(repoName: string) {
        goto(`/repo/${repoName}`);
    }
</script>

<Card.Root class="w-full mt-4">
    <Card.Header class="text-center">
        <Card.Title>Repositories</Card.Title>
        <Card.Description
            >Select a repository to manage secrets.</Card.Description
        >
    </Card.Header>
    <Card.Content>
        {#if repositories.length === 0}
            <div class="text-center py-4 text-muted-foreground">
                No repositories found.
            </div>
        {:else}
            <div class="flex flex-col gap-2">
                {#each repositories as repo}
                    <button
                        class="flex items-center justify-between p-3 border rounded-md hover:bg-muted/50 transition-colors cursor-pointer group text-left w-full"
                        type="button"
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
                        <Button variant="default" size="sm">Manage</Button>
                    </button>
                {/each}
            </div>
        {/if}
    </Card.Content>
</Card.Root>
