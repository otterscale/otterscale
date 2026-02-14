<script lang="ts">
	import { Check, ChevronsUpDown, User } from '@lucide/svelte';
	import type { ComponentProps } from '@sjsf/form';
	import { onMount } from 'svelte';

	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils';

	// User type from Keycloak
	interface KeycloakUser {
		id: string;
		username: string;
		email?: string;
		firstName?: string;
		lastName?: string;
	}

	// Role options based on the K8s schema
	const roleOptions = [
		{ value: 'admin', label: 'Admin' },
		{ value: 'edit', label: 'Edit' },
		{ value: 'view', label: 'View' }
	];

	// This widget handles the entire user object (subject + name + role)
	// It replaces the default objectField for the user item in the array
	let { value = $bindable(), config }: ComponentProps['objectField'] = $props();

	// Extract current values from the form value
	let currentSubject = $derived((value as Record<string, unknown>)?.subject as string | undefined);
	let currentRole = $derived((value as Record<string, unknown>)?.role as string | undefined);

	// User list from API
	let userList = $state<KeycloakUser[]>([]);
	let loading = $state(true);

	// Debounce timer
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	// Fetch users from API with search parameter
	async function fetchUsers(search: string = '') {
		loading = true;
		try {
			const queryParams = search ? `search=${encodeURIComponent(search)}&max=10` : 'max=10';
			const res = await fetch(`/rest/users?${queryParams}`);
			if (res.ok) {
				userList = await res.json();
			} else {
				console.error('Failed to fetch users:', res.statusText);
			}
		} catch (error) {
			console.error('Error fetching users:', error);
		} finally {
			loading = false;
		}
	}

	// Fetch initial users on mount
	onMount(() => {
		fetchUsers();
	});

	// User Popover state
	let userOpen = $state(false);
	let searchQuery = $state('');
	let triggerRef = $state<HTMLButtonElement | null>(null);

	// Debounced search - call API when input changes
	function handleSearchChange(query: string) {
		searchQuery = query;

		if (debounceTimer) {
			clearTimeout(debounceTimer);
		}

		debounceTimer = setTimeout(() => {
			fetchUsers(query);
		}, 300); // 300ms debounce
	}

	// Use userList directly since filtering is done server-side
	const filteredUsers = $derived(userList);

	// Find selected user
	const selectedUser = $derived(userList.find((u) => u.id === currentSubject));

	// Get display name for a user
	function getDisplayName(user: KeycloakUser): string {
		if (user.firstName || user.lastName) {
			return `${user.firstName || ''} ${user.lastName || ''}`.trim();
		}
		return user.username;
	}

	// Display text for the user button
	const userDisplayText = $derived(
		selectedUser
			? `${getDisplayName(selectedUser)} (${selectedUser.email || selectedUser.username})`
			: ((value as Record<string, unknown>)?.name as string) || 'Select user...'
	);

	// Find selected role label
	const selectedRoleLabel = $derived(
		roleOptions.find((r) => r.value === currentRole)?.label ?? 'Select role...'
	);

	function handleUserSelect(user: KeycloakUser) {
		// Update subject and name in the value object
		value = {
			...(value as Record<string, unknown>),
			subject: user.id,
			name: `${getDisplayName(user)} (${user.email || user.username})`
		};
		userOpen = false;
		searchQuery = '';
	}

	function handleRoleChange(newRole: string) {
		value = {
			...(value as Record<string, unknown>),
			role: newRole
		};
	}

	// Check if the form is disabled or read-only
	const isDisabled = $derived(config.schema.readOnly === true);
</script>

<div class="user-select-widget">
	<div class="flex items-center gap-2">
		<!-- User Selector (flex-1 to take remaining space) -->
		<Popover.Root bind:open={userOpen}>
			<Popover.Trigger bind:ref={triggerRef} class="flex-1">
				{#snippet child({ props })}
					<Button
						variant="outline"
						class={cn('w-full justify-between', !selectedUser && 'text-muted-foreground')}
						disabled={isDisabled}
						{...props}
					>
						<span class="flex items-center gap-2 truncate">
							<User class="h-4 w-4 shrink-0" />
							<span class="truncate">{userDisplayText}</span>
						</span>
						<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
					</Button>
				{/snippet}
			</Popover.Trigger>
			<Popover.Content class="w-75 p-0" align="start">
				<Command.Root shouldFilter={false}>
					<Command.Input
						placeholder="Search users..."
						value={searchQuery}
						oninput={(e) => handleSearchChange(e.currentTarget.value)}
					/>
					<Command.List>
						{#if loading}
							<Command.Loading>Loading users...</Command.Loading>
						{:else}
							<Command.Empty>No users found.</Command.Empty>
							<Command.Group>
								{#each filteredUsers as user (user.id)}
									<Command.Item value={user.id} onSelect={() => handleUserSelect(user)}>
										<Check
											class={cn('mr-2 h-4 w-4', currentSubject !== user.id && 'text-transparent')}
										/>
										<div class="flex flex-col">
											<span class="font-medium">{getDisplayName(user)}</span>
											<span class="text-xs text-muted-foreground"
												>{user.email || user.username}</span
											>
										</div>
									</Command.Item>
								{/each}
							</Command.Group>
						{/if}
					</Command.List>
				</Command.Root>
			</Popover.Content>
		</Popover.Root>

		<!-- Role Selector (fixed width) -->
		<Select.Root type="single" bind:value={currentRole} onValueChange={handleRoleChange}>
			<Select.Trigger class="w-28">
				{#if currentRole}
					<span>{selectedRoleLabel}</span>
				{:else}
					<span class="text-muted-foreground">Role</span>
				{/if}
			</Select.Trigger>
			<Select.Content>
				{#each roleOptions as role (role.value)}
					<Select.Item value={role.value}>
						{role.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	<!-- Hidden fields to satisfy form validation if needed -->
	{#if selectedUser}
		<input type="hidden" name="subject" value={selectedUser.id} />
		<input type="hidden" name="name" value={getDisplayName(selectedUser)} />
	{/if}
	{#if currentRole}
		<input type="hidden" name="role" value={currentRole} />
	{/if}
</div>
