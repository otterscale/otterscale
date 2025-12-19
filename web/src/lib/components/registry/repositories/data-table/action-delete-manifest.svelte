<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DeleteManifestRequest, Manifest } from '$lib/api/registry/v1/registry_pb';
	import { RegistryService } from '$lib/api/registry/v1/registry_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		manifest,
		scope,
		reloadManager,
		onSuccess
	}: {
		manifest: Manifest;
		scope: string;
		reloadManager: ReloadManager;
		onSuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	let request = $state({} as DeleteManifestRequest);
	let verification = $state({ tag: '' });
	let open = $state(false);

	function init() {
		request = {
			scope: scope,
			digest: manifest.digest,
			repositoryName: manifest.repositoryName
		} as DeleteManifestRequest;
		verification = { tag: '' };
	}

	function close() {
		open = false;
	}

	let invalidity = $state({} as Booleanified<{ repositoryName: string; tag: string }>);
	const invalid = $derived(invalidity.repositoryName || invalidity.tag);
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger
		variant="ghost"
		class="group flex h-7 w-7 items-center justify-center text-destructive"
	>
		<Icon icon="ph:trash" class="transition-transform duration-200 group-hover:scale-110" />
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_manifest()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.repository_name()}</Form.Label>
					<SingleInput.Confirm
						required
						disabled
						target={manifest.repositoryName}
						bind:value={request.repositoryName}
						bind:invalid={invalidity.repositoryName}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.tag()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: 'Tag' })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={manifest.tag}
						bind:value={verification.tag}
						bind:invalid={invalidity.tag}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						const repoName = request.repositoryName;
						toast.promise(() => registryClient.deleteManifest(request), {
							loading: `Deleting ${repoName}...`,
							success: () => {
								reloadManager.force();
								if (onSuccess) onSuccess();
								return `Delete ${repoName} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${repoName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
