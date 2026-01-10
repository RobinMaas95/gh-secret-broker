<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { buttonVariants } from "$lib/components/ui/button";
    import * as AlertDialog from "$lib/components/ui/alert-dialog";
    import AddSecretDialog from "$lib/components/AddSecretDialog.svelte";
    import { toast } from "svelte-sonner";
    import Trash2 from "lucide-svelte/icons/trash-2";

    let {
        owner,
        repo,
        secrets: initialSecrets,
        csrfToken: initialCsrfToken,
    } = $props<{
        owner: string;
        repo: string;
        secrets: string[];
        csrfToken: string | null;
    }>();

    let secrets = $state<string[]>([]);
    let error = $state<string | null>(null);
    let csrfToken = $state<string | null>(null);

    $effect(() => {
        secrets = initialSecrets || [];
        csrfToken = initialCsrfToken || null;
    });

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
            toast("Secret deleted successfully", {
                icon: Trash2 as any,
            });
        } catch (e) {
            error = (e as Error).message;
        }
    }
</script>

<div class="container mx-auto py-8 max-w-4xl">
    <div class="mb-4">
        <a
            href="/"
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
                    {csrfToken}
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
            {#if error}
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
