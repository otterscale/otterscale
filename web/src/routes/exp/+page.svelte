<!-- <script lang="ts" module>
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	let selectedScope = $state(0);
	let selectedFacility = $state(0);
	let selectedCeph = $state({});

	const cephOptions = writable<SingleSelect.OptionType[]>([]);
	cephOptions.set([
		{
			value: { scopeUuid: 1, facilityName: 1 },
			label: '1',
			icon: 'ph:cube'
		} as SingleSelect.OptionType,
		{
			value: { scopeUuid: 2, facilityName: 2 },
			label: '2',
			icon: 'ph:cube'
		} as SingleSelect.OptionType,
		{
			value: { scopeUuid: 3, facilityName: 3 },
			label: '3',
			icon: 'ph:cube'
		} as SingleSelect.OptionType
	]);
</script>

<SingleSelect.Root options={cephOptions} bind:value={selectedCeph}>
	<SingleSelect.Trigger />
	<SingleSelect.Content>
		<SingleSelect.Options>
			<SingleSelect.Input />
			<SingleSelect.List>
				<SingleSelect.Empty>No results found.</SingleSelect.Empty>
				<SingleSelect.Group>
					{#each $cephOptions as option}
						<SingleSelect.Item
							{option}
							onclick={() => {
								selectedScope = option.value.scopeUuid;
								selectedFacility = option.value.facilityName;
							}}
						>
							<Icon
								icon={option.icon ? option.icon : 'ph:empty'}
								class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
							/>
							{option.label}
							<SingleSelect.Check {option} />
						</SingleSelect.Item>
					{/each}
				</SingleSelect.Group>
			</SingleSelect.List>
		</SingleSelect.Options>
	</SingleSelect.Content>
</SingleSelect.Root> -->

<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { Tabs } from '$lib/components/custom/tabs';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';

	import * as MultipleStepModal from './mutiple-step-modal';

	let step = $state(1);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<MultipleStepModal.Root bind:open steps={3}>
	<MultipleStepModal.Trigger>Trigger</MultipleStepModal.Trigger>
	<MultipleStepModal.Content>
		<MultipleStepModal.Stepper>
			<MultipleStepModal.Steps>
				<MultipleStepModal.Step text="STEP 1" icon="ph:sun" />
				<MultipleStepModal.Step icon="ph:moon" />
				<MultipleStepModal.Step />
			</MultipleStepModal.Steps>
			<MultipleStepModal.Models>
				<MultipleStepModal.Model>Content for Tab 1</MultipleStepModal.Model>
				<MultipleStepModal.Model>Content for Tab 2</MultipleStepModal.Model>
				<MultipleStepModal.Model>Content for Tab 3</MultipleStepModal.Model>
			</MultipleStepModal.Models>
		</MultipleStepModal.Stepper>
		<MultipleStepModal.Header>Header</MultipleStepModal.Header>
		<MultipleStepModal.Footer>
			<MultipleStepModal.Cancel>Cancel</MultipleStepModal.Cancel>
			<MultipleStepModal.Confirm>Confirm</MultipleStepModal.Confirm>
			<MultipleStepModal.Controllers>
				<MultipleStepModal.Back>Back</MultipleStepModal.Back>
				<MultipleStepModal.Next>Next</MultipleStepModal.Next>
			</MultipleStepModal.Controllers>
		</MultipleStepModal.Footer>
	</MultipleStepModal.Content>
</MultipleStepModal.Root>
