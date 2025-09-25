<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages.js';
	import { currentKubernetes } from '$lib/stores';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	const reloadManager: ReloadManager = getContext('reloadManager');

	let modelOptions = $state(
		writable<SingleSelect.OptionType[]>([
			{
				value: 'ibm-granite/granite-docling-258M',
				label: 'ibm-granite/granite-docling-258M',
				icon: 'ph:robot',
				information: 'Image-Text-to-Text · 0.3B',
			},
			{
				value: 'mistralai/Magistral-Small-2509',
				label: 'mistralai/Magistral-Small-2509',
				icon: 'ph:robot',
				information: '24B',
			},
			{
				value: 'microsoft/VibeVoice-1.5B',
				label: 'microsoft/VibeVoice-1.5B',
				icon: 'ph:robot',
				information: 'Text-to-Speech · 2.7B',
			},
		]),
	);
	let secretOptions = $state(
		writable<SingleSelect.OptionType[]>([
			{
				value: 'secret-llm-access',
				label: 'LLM Access Secret',
				icon: 'ph:key',
			},
			{
				value: 'secret-storage-credentials',
				label: 'Storage Credentials',
				icon: 'ph:key',
			},
			{
				value: 'secret-db-connection',
				label: 'Database Connection',
				icon: 'ph:key',
			},
		]),
	);
	let dataTypeOptions = $state(
		writable<SingleSelect.OptionType[]>([
			{
				value: 'float16',
				label: 'float16',
				icon: 'ph:binary',
			},
			{
				value: 'float32',
				label: 'float32',
				icon: 'ph:binary',
			},
			{
				value: 'bfloat16',
				label: 'bfloat16',
				icon: 'ph:binary',
			},
		]),
	);

	let isModelInvalid = $state(false);
	let isSecretInvalid = $state(false);

	const defaults = {
		scopeUuid: $currentKubernetes?.scopeUuid,
		facilityName: $currentKubernetes?.name,
		modelName: undefined,
		secret: undefined,
		dtype: undefined,
		temperature: undefined,
	};
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	onMount(() => {
		reloadManager.force();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Model</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.model()}</Form.Label>
					<SingleSelect.Root
						bind:options={modelOptions}
						bind:value={request.modelName}
						bind:invalid={isModelInvalid}
						required
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $modelOptions as option}
											<SingleSelect.Item {option}>
												<Icon
													icon={option.icon ? option.icon : 'ph:empty'}
													class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
												/>
												<span>
													<p>
														{option.label}
													</p>
													<p class="text-muted-foreground text-xs">
														{option.information}
													</p>
												</span>
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
					<Form.Label>Secret</Form.Label>
					<SingleSelect.Root
						bind:options={secretOptions}
						bind:value={request.secret}
						bind:invalid={isSecretInvalid}
						required
					>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $secretOptions as option}
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
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Data Type</Form.Label>
					<SingleSelect.Root bind:options={dataTypeOptions} bind:value={request.dtype}>
						<SingleSelect.Trigger />
						<SingleSelect.Content>
							<SingleSelect.Options>
								<SingleSelect.Input />
								<SingleSelect.List>
									<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
									<SingleSelect.Group>
										{#each $dataTypeOptions as option}
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
				</Form.Field>
				<Form.Field>
					<Form.Label>temperature</Form.Label>
					<SingleInput.General type="number" bind:value={request.temperature} min={0} max={1} step={0.1} />
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
					disabled={isModelInvalid || isSecretInvalid}
					onclick={() => {
						// toast.promise(() => storageClient.createBucket(request), {
						// 	loading: `Creating ${request.bucketName}...`,
						// 	success: () => {
						// 		reloadManager.force();
						// 		return `Create ${request.bucketName}`;
						// 	},
						// 	error: (error) => {
						// 		let message = `Fail to create ${request.bucketName}`;
						// 		toast.error(message, {
						// 			description: (error as ConnectError).message.toString(),
						// 			duration: Number.POSITIVE_INFINITY,
						// 		});
						// 		return message;
						// 	},
						// });
						// reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>
