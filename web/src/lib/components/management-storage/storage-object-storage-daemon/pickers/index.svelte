<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import CephPicker from './ceph-picker.svelte';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable()
	}: { selectedScope: string; selectedFacility: string } = $props();
</script>

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Ceph</Picker.Label>
		<CephPicker bind:selectedScope bind:selectedFacility />
	</Picker.Wrapper>
</Picker.Root>
