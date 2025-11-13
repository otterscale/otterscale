<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type {
		SMBShare,
		SMBShare_User,
		UpdateSMBShareRequest
	} from '$lib/api/storage/v1/storage_pb';
	import {
		type SMBShare_ActiveDirectory,
		type SMBShare_LocalUser,
		SMBShare_MapToGuest,
		SMBShare_SecurityMode,
		StorageService
	} from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		scope,
		facility,
		namespace,
		smbShare,
		reloadManager
	}: {
		scope: string;
		facility: string;
		namespace: string;
		smbShare: SMBShare;
		reloadManager: ReloadManager;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	const defaults = {
		scope: scope,
		facility: facility,
		namespace: namespace,
		name: smbShare.name,
		sizeBytes: smbShare.sizeBytes,
		securityMode: smbShare.securityMode,
		auth: smbShare.auth,
		mapToGuest: smbShare.mapToGuest,
		browsable: smbShare.browsable,
		guestOk: smbShare.guestOk,
		readOnly: smbShare.readOnly,
		validUsers: smbShare.validUsers
	} as UpdateSMBShareRequest;
	const defaults_activeDirectory =
		smbShare.auth.case === 'activeDirectory'
			? smbShare.auth.value
			: ({} as SMBShare_ActiveDirectory);
	const defaults_localUser = {
		users: smbShare.auth.case === 'localUser' ? smbShare.auth.value.users : ([] as SMBShare_User[])
	} as SMBShare_LocalUser;
	const defaults_validUsers = smbShare.validUsers;
	const defaults_user = {} as SMBShare_User;

	let request = $state(defaults);
	let request_activeDirectory = $state(defaults_activeDirectory);
	let request_localUser = $state(defaults_localUser);
	let request_validUsers = $state(defaults_validUsers);
	let request_user = $state(defaults_user);

	function reset_activeDirectory() {
		request_activeDirectory = defaults_activeDirectory;
	}
	function reset_localUser() {
		request_localUser = defaults_localUser;
	}
	function reset_validUsers() {
		request_validUsers = defaults_validUsers;
	}
	function reset_user() {
		request_user = defaults_user;
	}

	function reset() {
		request = defaults;
		reset_activeDirectory();
		reset_localUser();
		reset_validUsers();
		reset_user();
	}

	$effect(() => {
		if (request.securityMode === SMBShare_SecurityMode.ACTIVE_DIRECTORY) {
			request.auth = {
				value: request_activeDirectory,
				case: 'activeDirectory'
			};
		} else if (request.securityMode === SMBShare_SecurityMode.USER) {
			request.auth = {
				value: request_localUser,
				case: 'localUser'
			};
		}
		request.validUsers = request_validUsers;
	});

	let isNameInvalid = $state(false);
	let isSizeInvalid = $state(false);
	let isSecurityModeInvalid = $state(false);
	let isMapToGuestInvalid = $state(false);
	let isBrowseableInvalid = $state(false);
	let isReadOnlyInvalid = $state(false);
	let isGuestOKInvalid = $state(false);
	let isRealmInvalid = $state(false);
	let isJoinSourceInvalid = $state(false);
	const invalid = $derived(
		isNameInvalid ||
			isSizeInvalid ||
			isSecurityModeInvalid ||
			isMapToGuestInvalid ||
			isBrowseableInvalid ||
			isReadOnlyInvalid ||
			isGuestOKInvalid ||
			(request.securityMode === SMBShare_SecurityMode.ACTIVE_DIRECTORY &&
				(isRealmInvalid || isJoinSourceInvalid)) ||
			(request.securityMode === SMBShare_SecurityMode.USER &&
				request.auth.case === 'localUser' &&
				request.auth.value.users.length === 0)
	);

	let open = $state(false);
	function close() {
		open = false;
	}

	let securityModeOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable([
			{
				value: SMBShare_SecurityMode.ACTIVE_DIRECTORY,
				label: 'Active Directory',
				icon: 'ph:database'
			},
			{
				value: SMBShare_SecurityMode.USER,
				label: 'User',
				icon: 'ph:textbox'
			}
		])
	);
	let mapToGuestOptions: Writable<SingleSelect.OptionType[]> = $state(
		writable([
			{
				value: SMBShare_MapToGuest.NEVER,
				label: 'Never',
				icon: 'ph:shield-slash',
				information: 'Always refuse access'
			},
			{
				value: SMBShare_MapToGuest.BAD_USER,
				label: 'Bad User',
				icon: 'ph:shield-warning',
				information: 'Map to guest when username is incorrect'
			},
			{
				value: SMBShare_MapToGuest.BAD_PASSWORD,
				label: 'Bad Password',
				icon: 'ph:shield-warning',
				information: 'Map to guest when password is incorrect'
			}
		])
	);
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:penvil" />
		{m.update()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.update()}</Modal.Header>
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
				<Form.Field>
					<Form.Label>{m.security_mode()}</Form.Label>
					<SingleSelect.Root
						required
						bind:options={securityModeOptions}
						bind:value={request.securityMode}
						bind:invalid={isSecurityModeInvalid}
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

				{#if request.securityMode === SMBShare_SecurityMode.ACTIVE_DIRECTORY}
					<Form.Field>
						<Form.Label>{m.realm()}</Form.Label>
						<SingleInput.General
							type="text"
							required={request.securityMode === SMBShare_SecurityMode.ACTIVE_DIRECTORY}
							bind:value={request_activeDirectory.realm}
							bind:invalid={isRealmInvalid}
						/>
					</Form.Field>
					<Form.Field>
						<Form.Label>{m.join_source()}</Form.Label>
						<SingleInput.General
							type="text"
							bind:value={request_activeDirectory.joinSource}
							bind:invalid={isJoinSourceInvalid}
						/>
					</Form.Field>
				{:else if request.securityMode === SMBShare_SecurityMode.USER}
					<Form.Field>
						<Form.Label>{m.users()}</Form.Label>
						{#if request.auth.case === 'localUser' && request_localUser.users.length > 0}
							<div class="rounded-lg border p-2">
								{#each request_localUser.users as user, index (index)}
									<div class="flex items-center gap-2 rounded-lg p-2">
										<div
											class={cn(
												'flex size-8 items-center justify-center rounded-full border-2 transition-all duration-300',
												request_validUsers.includes(user.username)
													? 'border-green-500 bg-green-100 text-green-700'
													: 'border-red-500 bg-red-100 text-red-700'
											)}
										>
											<Tooltip.Provider>
												<Tooltip.Root>
													<Tooltip.Trigger>
														<Icon
															icon={request_validUsers.includes(user.username)
																? 'ph:user-check'
																: 'ph:user'}
															class="size-5"
														/>
													</Tooltip.Trigger>
													<Tooltip.Content>
														{#if request_validUsers.includes(user.username)}
															Accessible for SMB Authentication
														{:else}
															Not accessible for SMB Authentication
														{/if}
													</Tooltip.Content>
												</Tooltip.Root>
											</Tooltip.Provider>
										</div>

										<div class="flex flex-col gap-1">
											<p class="text-xs text-muted-foreground">{m.user()}</p>
											<p class="text-sm">{user.username}</p>
										</div>

										<div class="ml-auto">
											<Button
												variant="ghost"
												size="icon"
												disabled={smbShare.validUsers.includes(user.username)}
												onclick={() => {
													if (request.auth.case === 'localUser') {
														request_localUser.users.splice(index, 1);
													}
												}}
											>
												<Icon icon="ph:trash" class="size-4 text-destructive" />
											</Button>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<div
								class="rounded-lg border border-red-300 bg-destructive/10 p-4 text-center text-xs text-destructive"
							>
								There is no user. Please add users for this share.
							</div>
						{/if}

						<Form.Label>{m.name()}</Form.Label>
						<SingleInput.General
							type="text"
							bind:value={request_user.username}
							required={request_localUser.users.length === 0}
						/>
						<Form.Label>{m.password()}</Form.Label>
						<SingleInput.General
							type="password"
							bind:value={request_user.password}
							required={request_localUser.users.length === 0}
						/>

						<Button
							onclick={() => {
								if (request.auth.case === 'localUser') {
									request_localUser.users = [...request_localUser.users, request_user];
									request_validUsers = [...request_validUsers, request_user.username];
									reset_user();
								}
							}}
						>
							<Icon icon="ph:plus" class="size-4" />
						</Button>
					</Form.Field>
				{/if}
				<Form.Field>
					<Form.Label>{m.map_to_guest()}</Form.Label>
					<Form.Help>
						Configure authentication failure handling and guest access behavior.
					</Form.Help>
					<SingleSelect.Root
						bind:options={mapToGuestOptions}
						bind:value={request.mapToGuest}
						bind:invalid={isMapToGuestInvalid}
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
				<Form.Field>
					<MultipleInput.Root type="text" bind:values={request_validUsers}>
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
						bind:invalid={isBrowseableInvalid}
						descriptor={() =>
							'Allow this share to be visible in network browsing and share enumeration.'}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.read_only()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.readOnly}
						bind:invalid={isReadOnlyInvalid}
						descriptor={() =>
							'Prevent write operations to this share. Users can only read and download files.'}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.guest_accessible()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.guestOk}
						bind:invalid={isGuestOKInvalid}
						descriptor={() => 'Allow anonymous guest access without authentication credentials.'}
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
					console.log(request);
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
