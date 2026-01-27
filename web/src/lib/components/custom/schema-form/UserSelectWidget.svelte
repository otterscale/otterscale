<script lang="ts">
	import { Check, ChevronsUpDown, User } from '@lucide/svelte';
	import type { ComponentProps } from '@sjsf/form';

	import { Button } from '$lib/components/ui/button';
	import * as Command from '$lib/components/ui/command';
	import * as Popover from '$lib/components/ui/popover';
	import * as Select from '$lib/components/ui/select';
	import { searchUsers, type User as UserType, users } from '$lib/stores/users';
	import { cn } from '$lib/utils';

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

	// User list from store
	let userList = $state<UserType[]>([]);
	users.subscribe((u) => (userList = u));

	// User Popover state
	let userOpen = $state(false);
	let searchQuery = $state('');
	let triggerRef = $state<HTMLButtonElement | null>(null);

	// Filtered users based on search
	const filteredUsers = $derived(searchUsers(searchQuery, userList));

	// Find selected user
	const selectedUser = $derived(userList.find((u) => u.subject === currentSubject));

	// Display text for the user button
	const userDisplayText = $derived(
		selectedUser
			? `${selectedUser.name} (${selectedUser.email || selectedUser.subject})`
			: 'Select user...'
	);

	// Find selected role label
	const selectedRoleLabel = $derived(
		roleOptions.find((r) => r.value === currentRole)?.label ?? 'Select role...'
	);

	function handleUserSelect(user: UserType) {
		// Update subject and name in the value object
		value = {
			...(value as Record<string, unknown>),
			subject: user.subject,
			name: user.name
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
					<Command.Input placeholder="Search users..." bind:value={searchQuery} />
					<Command.List>
						<Command.Empty>No users found.</Command.Empty>
						<Command.Group>
							{#each filteredUsers as user (user.subject)}
								<Command.Item value={user.subject} onSelect={() => handleUserSelect(user)}>
									<Check
										class={cn(
											'mr-2 h-4 w-4',
											currentSubject !== user.subject && 'text-transparent'
										)}
									/>
									<div class="flex flex-col">
										<span class="font-medium">{user.name}</span>
										<span class="text-xs text-muted-foreground">{user.email || user.subject}</span>
									</div>
								</Command.Item>
							{/each}
						</Command.Group>
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
		<input type="hidden" name="subject" value={selectedUser.subject} />
		<input type="hidden" name="name" value={selectedUser.name} />
	{/if}
	{#if currentRole}
		<input type="hidden" name="role" value={currentRole} />
	{/if}
</div>
