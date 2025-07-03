<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import CephPicker from './ceph-picker.svelte';
	import VolumePicker from './volume-picker.svelte';
	import SubvolumeGroupPicker from './subvolume-group-picker.svelte';
</script>

<script lang="ts">
	let {
		selectedScope = $bindable(),
		selectedFacility = $bindable(),
		selectedVolume = $bindable(),
		selectedSubvolumeGroup = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroup: string;
	} = $props();

	$effect(() => {
		selectedScope;
		selectedVolume = '';
		selectedSubvolumeGroup = '';
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

		{#key selectedVolume}
			<Picker.Wrapper class="*:h-8">
				<Picker.Label>Group</Picker.Label>
				<SubvolumeGroupPicker
					{selectedScope}
					{selectedFacility}
					{selectedVolume}
					bind:selectedSubvolumeGroup
				/>
			</Picker.Wrapper>
		{/key}
	{/key}
</Picker.Root>
