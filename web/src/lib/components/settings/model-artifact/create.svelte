<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import { type CreateModelArtifactRequest, ModelService } from '$lib/api/model/v1/model_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import ModelsStore from './util-models-store.svelte';
	import { fetchHuggingFaceModels } from './utils.svelte';

	interface OptionType extends SingleSelect.OptionType {
		downloads: number;
		likes: number;
		createdAt: string;
	}
</script>

<script lang="ts">
	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);
	const applicationClient = createClient(ApplicationService, transport);

	const namespaceOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	async function fetchNamespaceOptions() {
		const response = await applicationClient.listNamespaces({ scope });
		namespaceOptions.set(
			response.namespaces.map((namespace) => ({
				value: namespace.name,
				label: namespace.name,
				icon: 'ph:cube'
			}))
		);
	}

	const huggingFaceModelOptions: Writable<OptionType[]> = writable([]);

	const defaults = {
		scope: scope,
		namespace: 'llm-d',
		size: BigInt(100 * 1024 ** 3)
	} as CreateModelArtifactRequest;
	let request = $state({ ...defaults });
	function reset() {
		request = { ...defaults };
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	async function fetchModelOptions() {
		try {
			const huggingfaceModels = await fetchHuggingFaceModels('RedHatAI', [], 'downloads', 10);
			huggingFaceModelOptions.set(
				huggingfaceModels.map((model) => ({
					value: model.id,
					label: model.id,
					icon: 'ph:robot',
					downloads: model.downloads,
					likes: model.likes,
					createdAt: model.createdAt
				}))
			);
		} catch (error) {
			console.error('Failed to fetch Hugging Face models:', error);
		}
	}

	let invalidity = $state({} as Booleanified<CreateModelArtifactRequest>);
	const invalid = $derived(
		invalidity.name || invalidity.namespace || invalidity.modelName || invalidity.size
	);

	onMount(async () => {
		try {
			await Promise.all([fetchNamespaceOptions(), fetchModelOptions()]);
		} catch (error) {
			console.error('Error during initialization:', error);
		}
	});
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			reset();
		}
	}}
>
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
					<SingleInput.General
						type="text"
						required
						bind:value={request.name}
						bind:invalid={invalidity.name}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>

					<SingleSelect.Root
						bind:value={request.namespace}
						options={namespaceOptions}
						required
						bind:invalid={invalidity.namespace}
					>
						<SingleSelect.Trigger>
							<Icon icon="ph:cube" />
							{request.namespace}
							<Icon icon="ph:caret-down" class="ml-auto size-5" />
						</SingleSelect.Trigger>
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $namespaceOptions as option (option.value)}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visible' : 'invisible')}
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
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.model_name()}</Form.Label>
					<ModelsStore
						bind:value={request.modelName}
						required
						bind:invalid={invalidity.modelName}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						bind:value={request.size}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
						required
						bind:invalid={invalidity.size}
						transformer={(value) => String(value)}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>

		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				disabled={invalid}
				onclick={() => {
					toast.promise(() => modelClient.createModelArtifact(request), {
						loading: `Creating model artifact ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Successfully created model artifact ${request.name}`;
						},
						error: (error) => {
							let message = `Failed to create model artifact ${request.name}`;
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
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
