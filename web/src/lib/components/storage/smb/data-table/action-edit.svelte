<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { SMBShare, UpdateSMBShareRequest } from '$lib/api/storage/v1/storage_pb';
	import type {
		CreateSMBShareRequest,
		SMBShare_SecurityConfig
	} from '$lib/api/storage/v1/storage_pb';
	import {
		SMBShare_CommonConfig_MapToGuest,
		SMBShare_SecurityConfig_Mode,
		StorageService
	} from '$lib/api/storage/v1/storage_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';

	import ManipulateUser from './util-manipulate-user.svelte';
	import ManipulateUsers from './util-manipulate-users.svelte';
	import ManipulateValidUsers from './util-manipulate-valid-users.svelte';
	import ManipulateVerifiedValidUsers from './util-manipulate-verified-valid-users.svelte';
</script>

<script lang="ts">
	let {
		smbShare,
		scope,
		reloadManager
	}: {
		smbShare: SMBShare;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	const defaults = {
		scope: scope,
		name: smbShare.name,
		sizeBytes: smbShare.sizeBytes,
		browsable: smbShare.browsable,
		guestOk: smbShare.guestOk,
		readOnly: smbShare.readOnly,
		validUsers: smbShare.validUsers,
		commonConfig: { ...smbShare.commonConfig },
		securityConfig: { ...smbShare.securityConfig }
	} as UpdateSMBShareRequest;

	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let invaliditySMBShare = $state({} as Booleanified<CreateSMBShareRequest>);
	let invaliditySecurityConfig = $state({} as Booleanified<SMBShare_SecurityConfig>);
	const invalid = $derived(
		invaliditySMBShare.sizeBytes ||
			invaliditySecurityConfig.mode ||
			(request.securityConfig?.mode === SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY &&
				(invaliditySecurityConfig.realm || invaliditySecurityConfig.joinSource)) ||
			(request.securityConfig?.mode === SMBShare_SecurityConfig_Mode.USER &&
				invaliditySecurityConfig.localUsers)
	);

	let open = $state(false);
	function close() {
		open = false;
	}

	const securityModeOptions: Writable<SingleSelect.OptionType[]> = writable([
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
	]);
	const mapToGuestOptions: Writable<SingleSelect.OptionType[]> = writable([
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
	]);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:pencil" />
		{m.update()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.update()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General disabled type="text" bind:value={request.name} />
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.size()}</Form.Label>
					<SingleInput.Measurement
						required
						type="number"
						transformer={(value) => String(value)}
						bind:value={request.sizeBytes}
						bind:invalid={invaliditySMBShare.sizeBytes}
						units={[
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
							options={mapToGuestOptions}
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
														class={cn('size-5', option.icon ? 'visible' : 'invisible')}
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
							required
							options={securityModeOptions}
							bind:value={request.securityConfig.mode}
							bind:invalid={invaliditySecurityConfig.mode}
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

					{#if request.securityConfig.mode === SMBShare_SecurityConfig_Mode.USER}
						<Form.Field>
							<Form.Label>{m.local_users()}</Form.Label>
							{#if request.securityConfig.localUsers?.length > 0}
								<div class="group max-h-40 overflow-y-auto rounded-lg border p-2">
									{#each request.securityConfig.localUsers as user, index (index)}
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
											<CopyButton
												class="invisible ml-auto size-4 group-hover:visible"
												text={user.username}
											/>
										</div>
									{/each}
								</div>
							{/if}
							<ManipulateUsers
								type="update"
								required={request.securityConfig.mode === SMBShare_SecurityConfig_Mode.USER}
								bind:users={request.securityConfig.localUsers}
								bind:invalid={invaliditySecurityConfig.localUsers}
							/>
						</Form.Field>
					{/if}

					{#if request.securityConfig.mode === SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY}
						<Form.Field>
							<Form.Label>{m.realm()}</Form.Label>
							<SingleInput.General
								type="text"
								bind:value={request.securityConfig.realm}
								bind:invalid={invaliditySecurityConfig.realm}
								required={request.securityConfig.mode ===
									SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY}
							/>
						</Form.Field>

						<Form.Field>
							<Form.Label>{m.join_source()}</Form.Label>
							{#if request.securityConfig.joinSource}
								<div class="rounded-lg border p-2">
									<div class="flex items-center gap-2 rounded-lg p-2">
										<div
											class={cn('flex size-8 items-center justify-center rounded-full border-2')}
										>
											<Icon icon="ph:user" class="size-5" />
										</div>

										<div class="flex flex-col gap-1">
											<p class="text-xs text-muted-foreground">{m.user()}</p>
											<p class="text-sm">{request.securityConfig.joinSource.username}</p>
										</div>
									</div>
								</div>
							{/if}
							<ManipulateUser
								type="update"
								bind:user={request.securityConfig.joinSource}
								bind:invalid={invaliditySecurityConfig.joinSource}
								required={request.securityConfig.mode ===
									SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY}
							/>
						</Form.Field>
					{/if}

					<Form.Field>
						<Form.Label>{m.valid_users()}</Form.Label>
						{#if request.securityConfig.mode === SMBShare_SecurityConfig_Mode.USER}
							<ManipulateValidUsers bind:validUsers={request.validUsers} type="update" />
						{:else if request.securityConfig.mode === SMBShare_SecurityConfig_Mode.ACTIVE_DIRECTORY}
							<ManipulateVerifiedValidUsers
								bind:validUsers={request.validUsers}
								type="update"
								realm={request.securityConfig.realm}
								joinSource={request.securityConfig.joinSource}
							/>
						{/if}
					</Form.Field>
				{/if}
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
					toast.promise(() => storageClient.updateSMBShare(request), {
						loading: `Updating ${request.name}...`,
						success: () => {
							reloadManager.force();
							return `Update ${request.name} successfully`;
						},
						error: (error) => {
							let message = `Fail to update ${request.name}`;
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
