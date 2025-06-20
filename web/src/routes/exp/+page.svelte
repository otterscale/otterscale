<script lang="ts">
	import { Single as Input } from '$lib/components/custom/input/index';
	let value = $state('');
</script>

<div class="p-4">
	<Input.Measurement
		required
		bind:value
		units={[
			{ label: 'KB', value: 1024 },
			{ label: 'MB', value: 1024 * 1024 }
		]}
	/>
	{value}
</div>
