<script lang="ts">
    import * as Form from "$lib/components/ui/form";
    import { Input } from "$lib/components/ui/input";
    import { Textarea } from "$lib/components/ui/textarea";
    import * as Dialog from "$lib/components/ui/dialog";
    import { Button } from "$lib/components/ui/button";
    import { z } from "zod";
    import { defaults, superForm } from "sveltekit-superforms";
    import { zod } from "sveltekit-superforms/adapters";
    import { toast } from "svelte-sonner";
    import * as InputGroup from "$lib/components/ui/input-group";
    import * as Tooltip from "$lib/components/ui/tooltip";
    import Info from "lucide-svelte/icons/info";

    let { owner, repo, onSecretAdded, csrfToken } = $props<{
        owner: string;
        repo: string;
        onSecretAdded: (name: string) => void;
        csrfToken: string | null;
    }>();

    let open = $state(false);

    const formSchema = z.object({
        name: z
            .string()
            .regex(
                /^[a-zA-Z_][a-zA-Z0-9_]*$/,
                "Name must start with a letter or underscore and contain only alphanumeric characters and underscores. (Regex '^[a-zA-Z_][a-zA-Z0-9_]*$')",
            ),
        value: z.string().min(1, "Secret value is required"),
    });

    type FormSchema = z.infer<typeof formSchema>;

    const form = superForm<FormSchema>(
        defaults(zod(formSchema as any) as any),
        {
            SPA: true,
            validators: zod(formSchema as any),
            onUpdate: async ({ form: f }) => {
                if (f.valid) {
                    try {
                        if (!csrfToken)
                            throw new Error("CSRF token is missing");

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

                        onSecretAdded(f.data.name.toUpperCase());
                        toast.success(
                            `Secret ${f.data.name} created successfully`,
                        );
                        open = false;
                    } catch (e) {
                        toast.error((e as Error).message);
                    }
                }
            },
        },
    );

    const { form: formData, enhance, submitting, allErrors } = form;
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
            <Form.Field {form} name="name">
                <Form.Control>
                    {#snippet children({ props })}
                        <Form.Label>Name</Form.Label>
                        <InputGroup.Root>
                            <InputGroup.Input
                                {...props}
                                bind:value={$formData.name}
                                placeholder="SECRET_NAME"
                                disabled={$submitting}
                            />
                            <InputGroup.Addon align="inline-end">
                                <Tooltip.Root>
                                    <Tooltip.Trigger>
                                        {#snippet child({
                                            props: tooltipProps,
                                        })}
                                            <InputGroup.Button
                                                {...tooltipProps}
                                                aria-label="Info"
                                            >
                                                <Info class="size-4" />
                                            </InputGroup.Button>
                                        {/snippet}
                                    </Tooltip.Trigger>
                                    <Tooltip.Content>
                                        <p>
                                            Secret name must start with a letter
                                            or underscore and contain only
                                            alphanumeric characters and
                                            underscores (GitHub will convert it
                                            to uppercase)
                                        </p>
                                    </Tooltip.Content>
                                </Tooltip.Root>
                            </InputGroup.Addon>
                        </InputGroup.Root>
                    {/snippet}
                </Form.Control>
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
                            disabled={$submitting}
                        />
                    {/snippet}
                </Form.Control>
                <Form.FieldErrors />
            </Form.Field>

            <Dialog.Footer>
                <Button
                    type="submit"
                    disabled={$submitting ||
                        !$formData.name ||
                        !$formData.value ||
                        $allErrors.length > 0}
                    variant="default"
                >
                    {#if $submitting}Saving...{:else}Save Secret{/if}
                </Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>
