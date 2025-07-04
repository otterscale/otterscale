<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { DataTable } from './data-table';
	import type { Image, User } from '$gen/api/storage/v1/storage_pb';
	import { writable, type Writable } from 'svelte/store';

	let {
		selectedScope,
		selectedFacility,
		user,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		user: User;
		data: Writable<User[]>;
	} = $props();
</script>

<div class="flex items-center justify-end gap-1">
	{user.keys.length}
	<Sheet.Root>
		<Sheet.Trigger>
			<Icon icon="ph:arrow-square-out" />
		</Sheet.Trigger>
		<Sheet.Content class="min-w-[38vw]">
			<DataTable {selectedScope} {selectedFacility} {user} bind:data />
		</Sheet.Content>
	</Sheet.Root>
</div>
