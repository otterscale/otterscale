<script lang="ts">
	import Icon from '@iconify/svelte';

	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';

	import type { ModeSource } from '../types';
	import SelectCloudModel from './util-select-cloud-model.svelte';
	import SelectLocalModel from './util-select-local-model.svelte';

	let {
		value = $bindable(),
		modelSource = $bindable(),
		scope,
		namespace,
		required,
		invalid = $bindable(),
		fromLocal = $bindable()
	}: {
		value: string;
		modelSource?: ModeSource;
		scope: string;
		namespace: string;
		required: boolean;
		invalid: boolean;
		fromLocal: boolean;
	} = $props();

	$effect(() => {
		invalid = required && !value;
	});
</script>

<ButtonGroup.Root class="w-full" aria-invalid={invalid}>
	<InputGroup.Root>
		<InputGroup.Input
			placeholder="Select from Artifacts or HuggingFace"
			readonly
			class="cursor-default"
			{value}
		/>
		<InputGroup.Addon>
			<Icon icon="ph:robot" />
		</InputGroup.Addon>
	</InputGroup.Root>
	<SelectLocalModel bind:value {scope} {namespace} bind:fromLocal />
	<SelectCloudModel bind:value />
</ButtonGroup.Root>
