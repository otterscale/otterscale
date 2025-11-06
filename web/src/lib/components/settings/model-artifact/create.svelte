<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { HuggingFaceModel } from './types';
	import { fetchModels } from './utils.svelte';

	import { ModelService, type CreateModelArtifactRequest } from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { formatBigNumber, formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';

	interface OptionType extends SingleSelect.OptionType {
		downloads: number;
		likes: number;
		createdAt: string;
	}
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');
	const modelClient = createClient(ModelService, transport);

	let huggingFaceModels = $state([] as HuggingFaceModel[]);
	let huggingFaceModelOptions: Writable<OptionType[]> = $derived(
		writable(
			huggingFaceModels.map((model) => ({
				value: model.id,
				label: model.id,
				icon: 'ph:robot',
				downloads: model.downloads,
				likes: model.likes,
				createdAt: model.createdAt,
			})),
		),
	);

	let isLoading = $state(true);

	let isNameInvalid = $state(false);
	let isModelNameInvalid = $state(false);
	let isSizeInvalid = $state(false);

	const defaults = {
		scope: $currentKubernetes?.scope,
		facility: $currentKubernetes?.name,
		namespace: 'default',
	} as CreateModelArtifactRequest;
	let request = $state({ ...defaults });
	function reset() {
		request = { ...defaults };
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(async () => {
		await fetchModels('RedHatAI', undefined, undefined, 13)
			.then((response) => {
				huggingFaceModels = response;

				isLoading = false;
			})
			.catch((error) => {
				console.error('Error fetching models on create modal mount:', error);
			});
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>

	<Modal.Content>
		<Modal.Header>{m.create()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" required bind:value={request.name} bind:invalid={isNameInvalid} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.model_name()}</Form.Label>
					{#if isLoading}
						<Loading.Selection />
					{:else}
						<SingleSelect.Root
							required
							bind:options={huggingFaceModelOptions}
							bind:value={request.modelName}
							bind:invalid={isModelNameInvalid}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input />
									<SingleSelect.List>
										<SingleSelect.Empty>No results found.</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $huggingFaceModelOptions as option}
												<SingleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-6', option.icon ? 'visibale' : 'invisible')}
													/>
													<div class="font-mono">
														<h6 class="text-sm">{option.label}</h6>
														<div class="text-muted-foreground flex items-center text-xs">
															<span class="flex items-center gap-1">
																<Icon icon="ph:clock" />
																<p>{formatTimeAgo(new Date(option.createdAt))}</p>
															</span>
															<Icon icon="ph:dot-bold" />
															<span class="flex items-center gap-1">
																<Icon icon="ph:download-simple" />
																<p>{formatBigNumber(option.downloads)}</p>
															</span>
															<Icon icon="ph:dot-bold" />
															<span class="flex items-center gap-1">
																<Icon icon="ph:heart" />
																<p>{formatBigNumber(option.likes)}</p>
															</span>
														</div>
													</div>
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
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.General
						type="number"
						required
						bind:value={request.size}
						bind:invalid={isSizeInvalid}
					/>
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

			<Modal.Action
				disabled={isNameInvalid || isModelNameInvalid || isSizeInvalid}
				onclick={() => {
					toast.promise(() => modelClient.createModelArtifact(request), {
						loading: `Creating model ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Successfully created model ${request.name}`;
						},
						error: (error) => {
							let message = `Failed to create model ${request.name}`;
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
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
