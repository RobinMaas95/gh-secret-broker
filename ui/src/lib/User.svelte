<script lang="ts">
    import { userState } from "./user_state.svelte";
    import * as Card from "$lib/components/ui/card";
    import { Button } from "$lib/components/ui/button";
    import RepositoryList from "$lib/components/RepositoryList.svelte";
    import type { PageData } from "../routes/$types";

    let { data } = $props<{ data?: PageData }>();
    // userState is hydrated from layout load
</script>

<div class="flex flex-col min-h-[calc(100vh-60px)] items-center p-4 gap-4 mt-8">
    {#if userState.error}
        <Card.Root class="w-full max-w-lg">
            <Card.Content class="pt-6">
                <div class="text-destructive text-center mb-4">
                    {userState.error}
                </div>
                <div class="text-center">
                    <Button variant="link" href="/login"
                        >Go back to login</Button
                    >
                </div>
            </Card.Content>
        </Card.Root>
    {:else if userState.current}
        <div class="w-full max-w-5xl">
            <h2 class="text-3xl font-bold tracking-tight mb-6 text-center">
                Your Repositories
            </h2>
            {#if data?.repositories}
                <RepositoryList repositories={data.repositories} />
            {:else}
                <div class="text-muted-foreground">Loading repositories...</div>
            {/if}
        </div>
    {:else}
        <div class="w-full max-w-7xl space-y-6">
            <!-- Unauthorized Banner -->
            <div
                class="border-l-4 border-yellow-400 bg-yellow-50 dark:bg-yellow-900/20 p-4 rounded-lg"
            >
                <div class="flex items-center gap-3">
                    <div class="flex-shrink-0">
                        <svg
                            class="h-6 w-6 text-yellow-600 dark:text-yellow-400"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 9v2m0 4h.01m-6.938 4h.138M12 19v-6m0 4h.01M12 5v-6m0 4h.01M9 9h6m-6 4h6"
                            />
                        </svg>
                    </div>
                    <div>
                        <h3
                            class="text-lg font-semibold text-yellow-800 dark:text-yellow-200"
                        >
                            Authentication Required
                        </h3>
                        <p class="text-yellow-700 dark:text-yellow-300 mt-1">
                            Please log in to access your GitHub repository
                            secrets.
                        </p>
                        <div class="mt-3">
                            <Button
                                href="/login"
                                class="bg-yellow-600 hover:bg-yellow-700 text-white"
                            >
                                Go to Login
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Login Card -->
            <Card.Root class="w-full max-w-md">
                <Card.Header>
                    <Card.Title class="text-2xl font-bold text-center"
                        >GH Secret Broker</Card.Title
                    >
                    <Card.Description class="text-center">
                        Please log in to manage your GitHub secrets.
                    </Card.Description>
                </Card.Header>
                <Card.Content class="pt-6">
                    <div class="text-center">
                        <Button href="/login">Go to Login</Button>
                    </div>
                </Card.Content>
            </Card.Root>
        </div>
    {/if}
</div>
