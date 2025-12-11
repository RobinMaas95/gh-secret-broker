<script lang="ts">
    import * as Form from "$lib/components/ui/form";
    import { Input } from "$lib/components/ui/input";
    import { Textarea } from "$lib/components/ui/textarea";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { z } from "zod";
    import { defaults, superForm } from "sveltekit-superforms";
    import { zod } from "sveltekit-superforms/adapters";

    let { owner, repo, onSecretAdded } = $props<{
        owner: string;
        repo: string;
        onSecretAdded: (name: string) => void;
    }>();

    let open = $state(false);
    let loading = $state(false);
    let submitError = $state<string | null>(null);

    const formSchema = z.object({
        name: z
            .string()
            .regex(
                /^[a-zA-Z_][a-zA-Z0-9_]*$/,
                "Name must start with a letter or underscore and contain only alphanumeric characters and underscores.",
            ),
        value: z.string().min(1, "Secret value is required"),
    });

    type FormSchema = z.infer<typeof formSchema>;

    const form = superForm(defaults(zod(formSchema)), {
        SPA: true,
        validators: zod(formSchema),
        onUpdate: async ({ form: f }) => {
            if (f.valid) {
                loading = true;
                submitError = null;
                try {
                    // Get CSRF Token first
                    const csrfRes = await fetch("/api/csrf-token");
                    if (!csrfRes.ok)
                        throw new Error("Failed to get CSRF token");
                    const { token: csrfToken } = await csrfRes.json();

                    const res = await fetch(
                        `/api/repo/${owner}/${repo}/secrets/${f.data.name}`,
                        {
                            method: "PUT",
                            headers: {
                                "Content-Type": "application/json",
                                "X-CSRF-Token": csrfToken,
                            },
                            body: JSON.stringify({ value: f.data.value }),
                        },
                    );

                    if (!res.ok) {
                        throw new Error("Failed to create secret");
                    }

                    onSecretAdded(f.data.name);
                    open = false;

                    // Reset form manually? superForm usually handles reset but in SPA mode on success we might need to be explicit if dialog reopens
                    // But defaults() recreated it initially.
                    // Let's just rely on binding 'open' to reset if needed or just letting it be.
                } catch (e) {
                    submitError = (e as Error).message;
                } finally {
                    loading = false;
                }
            }
        },
    });

    const { form: formData, enhance } = form;
</script>

<Dialog.Root bind:open>
    <Dialog.Trigger>
        <Button>Add Secret</Button>
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-2xl">
        <Dialog.Header>
            <Dialog.Title>Add New Secret</Dialog.Title>
            <Dialog.Description>
                Add a new secret to this repository. The value will be
                encrypted.
            </Dialog.Description>
        </Dialog.Header>

        <form method="POST" use:enhance class="grid gap-4 py-4">
            {#if submitError}
                <div class="text-destructive text-sm">{submitError}</div>
            {/if}
            <Form.Field {form} name="name">
                <Form.Control>
                    {#snippet children({ props })}
                        <Form.Label>Name</Form.Label>
                        <Input
                            {...props}
                            bind:value={$formData.name}
                            placeholder="SECRET_NAME"
                            disabled={loading}
                        />
                    {/snippet}
                </Form.Control>
                <Form.Description
                    >Secret name must start with a letter or underscore and
                    contain only alphanumeric characters and underscores (GitHub
                    will convert to uppercase)</Form.Description
                >
                <Form.FieldErrors />
            </Form.Field>

            <Form.Field {form} name="value">
                <Form.Control>
                    {#snippet children({ props })}
                        <Form.Label>Value</Form.Label>
                        <Textarea
                            {...props}
                            bind:value={$formData.value}
                            placeholder="Secret value..."
                            class="font-mono"
                            disabled={loading}
                        />
                    {/snippet}
                </Form.Control>
                <Form.FieldErrors />
            </Form.Field>

            <Dialog.Footer>
                <Button type="submit" disabled={loading}>
                    {#if loading}Saving...{:else}Save Secret{/if}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
