<script lang="ts">
    import { untrack } from "svelte";
    import { userState } from "./user_state.svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import RepositoryList from "$lib/components/RepositoryList.svelte";

    // Trigger load if not loaded? Or assume App handles it.
    // Let's make sure it's loaded.
    $effect(() => {
        untrack(() => {
            if (!userState.current && !userState.loading && !userState.error) {
                userState.load();
            }
        });
    });
</script>

<div class="flex flex-col min-h-[calc(100vh-60px)] items-center p-4 gap-4 mt-8">
    {#if userState.error}
        <Card.Root class="w-full max-w-lg">
            <Card.Content class="pt-6">
                <div class="text-destructive text-center mb-4">
                    {userState.error}
                </div>
                <div class="text-center">
                    <Button variant="link" href="/">Go back to login</Button>
                </div>
            </Card.Content>
        </Card.Root>
    {:else if userState.current}
        <div class="w-full max-w-7xl">
            <h2 class="text-3xl font-bold tracking-tight mb-6">
                Your Repositories
            </h2>
            <RepositoryList />
        </div>
    {:else}
        <div class="flex items-center justify-center h-full">
            <p class="text-muted-foreground">Loading...</p>
        </div>
    {/if}
</div>
