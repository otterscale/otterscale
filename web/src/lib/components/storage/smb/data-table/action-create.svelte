<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import type {
		CreateSMBShareRequest,
		SMBShare_CommonConfig
	} from '$lib/api/storage/v1/storage_pb';
	import {
		SMBShare_CommonConfig_MapToGuest,
		SMBShare_SecurityConfig_Mode,
		StorageService
	} from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import * as Loading from '$lib/components/custom/loading';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';

	import CreateUser from './utils/create-user.svelte';
	import CreateUsers from './utils/create-users.svelte';
</script>

<script lang="ts">
	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);
	const applicationClient = createClient(ApplicationService, transport);

	let isNamespaceOptionsLoaded = $state(false);
	const namespaceOptions: Writable<SingleSelect.OptionType[]> = writable([]);
	async function fetchNamespaces() {
		applicationClient
			.listNamespaces({ scope })
			.then((response) => {
				namespaceOptions.set(
					response.namespaces.map((namespace) => ({
						value: namespace.name,
						label: namespace.name,
						icon: 'ph:cube'
					}))
				);
			})
			.catch((error) => {
				console.debug('Failed to fetch namespaces:', error);
			});
	}

	const defaults = {
		scope: scope,
		browsable: true,
		guestOk: false,
		readOnly: true,
		commonConfig: { mapToGuest: SMBShare_CommonConfig_MapToGuest.NEVER } as SMBShare_CommonConfig,
		securityConfig: { mode: SMBShare_SecurityConfig_Mode.USER }
	} as CreateSMBShareRequest;

	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let isNameInvalid = $state(false);
	let isNamespaceInvalid = $state(false);
	let isSizeInvalid = $state(false);
	const invalid = $derived(isNameInvalid || isNamespaceInvalid || isSizeInvalid);

	let open = $state(false);
	function close() {
		open = false;
	}

	let securityModeOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable([
			{
				value: SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY,
				label: 'Active Directory',
				icon: 'ph:database'
			},
			{
				value: SMBShare_SecurityConfig_Mode.USER,
				label: 'User',
				icon: 'ph:textbox'
			}
		])
	);
	let mapToGuestOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable([
			{
				value: SMBShare_CommonConfig_MapToGuest.NEVER,
				label: 'Never',
				icon: 'ph:shield-slash',
				information: m.map_to_guest_never_info()
			},
			{
				value: SMBShare_CommonConfig_MapToGuest.BAD_USER,
				label: 'Bad User',
				icon: 'ph:shield-warning',
				information: m.map_to_guest_bad_user_info()
			},
			{
				value: SMBShare_CommonConfig_MapToGuest.BAD_PASSWORD,
				label: 'Bad Password',
				icon: 'ph:shield-warning',
				information: m.map_to_guest_bad_password_info()
			}
		])
	);

	onMount(async () => {
		await fetchNamespaces();
		isNamespaceOptionsLoaded = true;
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_group()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General
						required
						type="text"
						bind:value={request.name}
						bind:invalid={isNameInvalid}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.namespace()}</Form.Label>
					{#if isNamespaceOptionsLoaded}
						<SingleSelect.Root
							required
							options={namespaceOptions}
							bind:value={request.namespace}
							bind:invalid={isNamespaceInvalid}
						>
							<SingleSelect.Trigger />
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
														class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
													/>
													<span class="flex flex-col">
														<p>{option.label}</p>
														<p class="text-xs text-muted-foreground">{option.information}</p>
													</span>
													<SingleSelect.Check {option} />
												</SingleSelect.Item>
											{/each}
										</SingleSelect.Group>
									</SingleSelect.List>
								</SingleSelect.Options>
							</SingleSelect.Content>
						</SingleSelect.Root>
					{:else}
						<Loading.Selection />
					{/if}
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						required
						type="number"
						transformer={(value) => String(value)}
						bind:value={request.sizeBytes}
						bind:invalid={isSizeInvalid}
						units={[
							{ value: 1024 ** 2, label: 'MB' } as SingleInput.UnitType,
							{ value: 1024 ** 3, label: 'GB' } as SingleInput.UnitType,
							{ value: 1024 ** 4, label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				{#if request.commonConfig}
					<Form.Field>
						<Form.Label>{m.map_to_guest()}</Form.Label>
						<Form.Help>
							{m.map_to_guest_help()}
						</Form.Help>
						<SingleSelect.Root
							bind:options={mapToGuestOptions}
							bind:value={request.commonConfig.mapToGuest}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $mapToGuestOptions as option (option.value)}
												<SingleSelect.Item {option}>
													<Icon
														icon={option.icon ? option.icon : 'ph:empty'}
														class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
													/>
													<span class="flex flex-col">
														<p>{option.label}</p>
														<p class="text-xs text-muted-foreground">{option.information}</p>
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
				{/if}

				{#if request.securityConfig}
					<Form.Field>
						<Form.Label>{m.mode()}</Form.Label>
						<SingleSelect.Root
							bind:options={securityModeOptions}
							bind:value={request.securityConfig.mode}
						>
							<SingleSelect.Trigger />
							<SingleSelect.Content>
								<SingleSelect.Options>
									<SingleSelect.Input class="border-0 ring-0" />
									<SingleSelect.List>
										<SingleSelect.Empty>{m.no_result()}</SingleSelect.Empty>
										<SingleSelect.Group>
											{#each $securityModeOptions as option (option.value)}
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
						<Form.Label>{m.local_user()}</Form.Label>
						{#if request.securityConfig.localUser}
							<div class="rounded-lg border p-2">
								<div class="flex items-center gap-2 rounded-lg p-2">
									<div class={cn('flex size-8 items-center justify-center rounded-full border-2')}>
										<Icon icon="ph:user" class="size-5" />
									</div>

									<div class="flex flex-col gap-1">
										<p class="text-xs text-muted-foreground">{m.user()}</p>
										<p class="text-sm">{request.securityConfig.localUser?.username}</p>
									</div>
								</div>
							</div>
						{/if}
						<CreateUser bind:user={request.securityConfig.localUser} />
					</Form.Field>

					<Form.Field>
						<Form.Label>{m.realm()}</Form.Label>
						<SingleInput.General type="text" bind:value={request.securityConfig.realm} />
					</Form.Field>

					<Form.Field>
						<Form.Label>{m.join_sources()}</Form.Label>
						{#if request.securityConfig.joinSources?.length > 0}
							<div class="max-h-40 overflow-y-auto rounded-lg border p-2">
								{#each request.securityConfig.joinSources as user, index (index)}
									<div class="flex items-center gap-2 rounded-lg p-2">
										<div
											class={cn('flex size-8 items-center justify-center rounded-full border-2')}
										>
											<Icon icon="ph:user" class="size-5" />
										</div>

										<div class="flex flex-col gap-1">
											<p class="text-xs text-muted-foreground">{m.user()}</p>
											<p class="text-sm">{user.username}</p>
										</div>
									</div>
								{/each}
							</div>
						{/if}
						<CreateUsers bind:users={request.securityConfig.joinSources} />
					</Form.Field>
				{/if}

				<Form.Field>
					<Form.Label>{m.valid_users()}</Form.Label>
					<MultipleInput.Root type="text" icon="ph:user" bind:values={request.validUsers}>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.browsable()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.browsable}
						descriptor={() => m.browsable_description()}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.read_only()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.readOnly}
						descriptor={() => m.read_only_description()}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.guest_accessible()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.guestOk}
						descriptor={() => m.guest_accessible_description()}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
					close();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				disabled={invalid}
				onclick={() => {
					toast.promise(() => storageClient.createSMBShare(request), {
						loading: `Creating ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Create ${request.name} successfully`;
						},
						error: (error) => {
							let message = `Fail to create ${request.name}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
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
