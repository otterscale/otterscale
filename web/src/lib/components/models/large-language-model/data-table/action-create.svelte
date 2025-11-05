<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { ModelService, type CreateModelRequest } from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);

	const reloadManager: ReloadManager = getContext('reloadManager');

	let isNameInvalid = $state(false);
	let isNamespaceInvalid = $state(false);
	let isModelNameInvalid = $state(false);

	let artifactOptions = $state(writable<SingleSelect.OptionType[]>([]));
	let isArtifactOptionsLoading = $state(true);

	const defaults = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		namespace: 'default',
	} as CreateModelRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	async function fetchArtifactOptions() {
		try {
			const response = await modelClient.listModelArtifacts({
				scope: $currentKubernetes?.scope,
				facility: $currentKubernetes?.name,
				namespace: 'default',
			});
			artifactOptions.set(
				response.modelArtifacts.map(
					(artifact) =>
						({
							value: artifact.name,
							label: artifact.name,
							icon: 'ph:robot',
						}) as SingleSelect.OptionType,
				),
			);

			isArtifactOptionsLoading = false;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}
	onMount(async () => {
		try {
			await fetchArtifactOptions();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.model_name()}</Form.Label>
					{#if isArtifactOptionsLoading}
						<Loading.Selection />
					{:else}
						<SingleSelect.Root
							required
							bind:options={artifactOptions}
							bind:value={request.modelName}
							bind:invalid={isModelNameInvalid}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $artifactOptions as option}
												<SingleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
													/>
													{option.label}
													<SingleSelect.Check {option} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					{/if}
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General required type="text" bind:value={request.name} bind:invalid={isNameInvalid} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={isNameInvalid || isNamespaceInvalid || isModelNameInvalid}
					onclick={() => {
						toast.promise(() => modelClient.createModel(request), {
							loading: `Creating ${request.modelName}...`,
							success: () => {
								reloadManager.force();
								return `Create ${request.modelName}`;
							},
							error: (error) => {
								let message = `Fail to create ${request.modelName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
								});
								return message;
							},
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
