<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable, writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { CreateSMBShareRequest } from '$lib/api/storage/v1/storage_pb';
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
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';

	import CreateUser from './utils/create-user.svelte';
	import CreateUsers from './utils/create-users.svelte';
</script>

<script lang="ts">
	let {
		scope,
		facility,
		namespace,
		reloadManager
	}: { scope: string; facility: string; namespace: string; reloadManager: ReloadManager } =
		$props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	const defaults = {
		scope: scope,
		facility: facility,
		namespace: namespace,
		browsable: true,
		guestOk: false,
		readOnly: true,
		mapToGuest: SMBShare_MapToGuest.NEVER
	} as CreateSMBShareRequest;
	const defaults_activeDirectory = {} as SMBShare_ActiveDirectory;
	const defaults_validUsers = [] as string[];
	const defaults_localUser = {} as SMBShare_LocalUser;

	let request = $state(defaults);
	let request_activeDirectory = $state(defaults_activeDirectory);
	let request_localUser = $state(defaults_localUser);
	let request_validUsers = $state(defaults_validUsers);

	function reset_activeDirectory() {
		request_activeDirectory = defaults_activeDirectory;
	}
	function reset_localUser() {
		request_localUser = defaults_localUser;
	}
	function reset_validUsers() {
		request_validUsers = defaults_validUsers;
	}
	function reset() {
		request = defaults;
		reset_activeDirectory();
		reset_localUser();
		reset_validUsers();
	}

	let isNameInvalid = $state(false);
	let isSizeInvalid = $state(false);
	let isSecurityModeInvalid = $state(false);
	let isMapToGuestInvalid = $state(false);
	let isBrowseableInvalid = $state(false);
	let isReadOnlyInvalid = $state(false);
	let isGuestOKInvalid = $state(false);
	let isRealmInvalid = $state(false);
	let isJoinSourceInvalid = $state(false);
	let isLocalUsersInvalid = $state(false);
	let isValidUsersInvalid = $state(false);
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
				isLocalUsersInvalid) ||
			isValidUsersInvalid
	);

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
				information: m.map_to_guest_never_info()
			},
			{
				value: SMBShare_MapToGuest.BAD_USER,
				label: 'Bad User',
				icon: 'ph:shield-warning',
				information: m.map_to_guest_bad_user_info()
			},
			{
				value: SMBShare_MapToGuest.BAD_PASSWORD,
				label: 'Bad Password',
				icon: 'ph:shield-warning',
				information: m.map_to_guest_bad_password_info()
			}
		])
	);
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
						<CreateUser
							bind:user={request_activeDirectory.joinSource}
							bind:invalid={isJoinSourceInvalid}
						/>
					</Form.Field>
				{:else if request.securityMode === SMBShare_SecurityMode.USER}
					<Form.Field>
						<Form.Label>{m.users()}</Form.Label>
						<CreateUsers bind:users={request_localUser.users} bind:invalid={isLocalUsersInvalid} />
					</Form.Field>
				{/if}
				<Form.Field>
					<Form.Label>{m.map_to_guest()}</Form.Label>
					<Form.Help>
						{m.map_to_guest_help()}
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
					<Form.Label>{m.valid_users()}</Form.Label>

					<MultipleInput.Root
						type="text"
						icon="ph:user"
						bind:values={request.validUsers}
						bind:invalid={isValidUsersInvalid}
						required
					>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							{#if request.securityMode === SMBShare_SecurityMode.USER}
								<MultipleInput.Import
									values={request_localUser.users.map((user) => user.username)}
								/>
							{/if}
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
						descriptor={() => m.browsable_description()}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.read_only()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.readOnly}
						bind:invalid={isReadOnlyInvalid}
						descriptor={() => m.read_only_description()}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.guest_accessible()}</Form.Label>
					<SingleInput.Boolean
						bind:value={request.guestOk}
						bind:invalid={isGuestOKInvalid}
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
