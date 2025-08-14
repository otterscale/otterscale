<script lang="ts" module>
	import * as Picker from '$lib/components/custom/picker';
	import SubvolumeGroupPicker from '../../utils/subvolume-group-picker.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		selectedVolume,
		selectedSubvolumeGroupName = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		selectedVolume: string;
		selectedSubvolumeGroupName: string;
	} = $props();
</script>

<Picker.Root align="left">
	<Picker.Wrapper class="*:h-8">
		<Picker.Label>Group</Picker.Label>
		<SubvolumeGroupPicker
			{selectedScopeUuid}
			{selectedFacility}
			{selectedVolume}
			bind:selectedSubvolumeGroupName
		/>
	</Picker.Wrapper>
</Picker.Root>
