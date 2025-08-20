<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';

	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<Alert.Root>
	<Icon icon="ph:spinner" class="size-4 animate-spin" />
	<Alert.Title>{$machine.status}</Alert.Title>
	<Alert.Description>{$machine.statusMessage}</Alert.Description>
</Alert.Root>
