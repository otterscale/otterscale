<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import CephPicker from '../utils/ceph-picker.svelte';
	import VolumePicker from '../utils/volume-picker.svelte';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable()
	}: { selectedScope: string; selectedFacility: string; selectedVolume: string } = $props();

	$effect(() => {
		selectedScope;
		selectedVolume = '';
	});
</script>

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Ceph</Picker.Label>
		<CephPicker bind:selectedScope bind:selectedFacility />
	</Picker.Wrapper>

	{#key selectedScope}
		<Picker.Wrapper class="*:h-8">
			<Picker.Label>Volume</Picker.Label>
			<VolumePicker {selectedScope} {selectedFacility} bind:selectedVolume />
		</Picker.Wrapper>
	{/key}
</Picker.Root>
