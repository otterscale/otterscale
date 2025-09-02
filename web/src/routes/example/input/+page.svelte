<script lang="ts" module>
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import {
		LayeredMultiple as LayeredMultipleSelect,
		LayeredSingle as LayeredSingleSelect,
		Multiple as MultipleSelect,
		Single as SingleSelect,
	} from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { writable, type Writable } from 'svelte/store';

	type Values = {
		value1: any;
		value2: any;
		value3: any;
		value4: any;
		value5: any;
		value6: any;
		value7: string;
		value8: any;
		value9: any;
		value10: any;
		values1: any[];
		values2: any[];
		values3: any[];
		values4: any[];
		values5: any[];
		values6: any[];
		values7: any[];
		values8: any[];
		values9: any[];
		values0: any[];
	};
</script>

<script lang="ts">
	let options1: Writable<SingleSelect.OptionType[]> = $state(
		writable([
			{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
			{
				value: 'star',
				label: 'Star',
				icon: 'ph:star',
			},
			{
				value: 'sun',
				label: 'Sun',
				icon: 'ph:sun',
			},
		]),
	);
	let options2: Writable<MultipleSelect.OptionType[]> = $state(
		writable([
			{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
			{
				value: 'star',
				label: 'Star',
				icon: 'ph:star',
			},
			{
				value: 'sun',
				label: 'Sun',
				icon: 'ph:sun',
			},
		]),
	);
	let options3: LayeredSingleSelect.OptionType[] = [
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star',
		},
		{
			value: 'sun',
			label: 'Sun',
			icon: 'ph:sun',
			subOptions: [
				{ value: 'sunrise', label: 'Sunrise', icon: 'ph:arrow-up' },
				{
					value: 'sunset',
					label: 'Sunset',
					icon: 'ph:arrow-down',
					subOptions: [
						{
							value: 'golden-hour',
							label: 'Golden Hour',
							icon: 'ph:clock',
						},
					],
				},
			],
		},
	];
	const options4: LayeredMultipleSelect.OptionType[] = [
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star',
		},
		{
			value: 'sun',
			label: 'Sun',
			icon: 'ph:sun',
			subOptions: [
				{ value: 'sunrise', label: 'Sunrise', icon: 'ph:arrow-up' },
				{
					value: 'sunset',
					label: 'Sunset',
					icon: 'ph:arrow-down',
					subOptions: [
						{
							value: 'golden-hour',
							label: 'Golden Hour',
							icon: 'ph:clock',
						},
					],
				},
			],
		},
	];
	let values: Values = $state({} as Values);
</script>

<main class="flex flex-col gap-4 p-4">
	<SingleInput.General id="1" required type="text" bind:value={values.value1} />

	<SingleInput.General id="2" required type="number" bind:value={values.value2} />

	<SingleInput.Boolean id="3" required bind:value={values.value3} />

	<SingleInput.Boolean id="4" required format="switch" bind:value={values.value4} />

	<SingleInput.Password id="5" required bind:value={values.value5} />

	<SingleInput.Confirm id="6" required target="value6" bind:value={values.value6} />

	<SingleInput.Measurement
		id="7"
		required
		units={[
			{ value: 1, label: 'I' } as SingleInput.UnitType,
			{ value: 2, label: 'II' } as SingleInput.UnitType,
			{ value: 3, label: 'III' } as SingleInput.UnitType,
		]}
		bind:value={values.value7}
	/>

	<SingleInput.Structure id="14" language="json" required bind:value={values.value7} />

	<MultipleInput.Root type="text" required bind:values={values.values1} id="8">
		<MultipleInput.Viewer />
		<MultipleInput.Controller>
			<MultipleInput.Input />
			<MultipleInput.Add />
			<MultipleInput.Clear />
		</MultipleInput.Controller>
	</MultipleInput.Root>

	<MultipleInput.Root type="number" required bind:values={values.values2} id="9">
		<MultipleInput.Viewer />
		<MultipleInput.Controller>
			<MultipleInput.Input />
			<MultipleInput.Add />
			<MultipleInput.Clear />
		</MultipleInput.Controller>
	</MultipleInput.Root>

	<SingleSelect.Root id="10" bind:options={options1} bind:value={values.value8} required>
		<SingleSelect.Trigger />
		<SingleSelect.Content>
			<SingleSelect.Options>
				<SingleSelect.Input />
				<SingleSelect.List>
					<SingleSelect.Empty>No results found.</SingleSelect.Empty>
					<SingleSelect.Group>
						{#each $options1 as option}
							<SingleSelect.Item {option}>
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
	</SingleSelect.Root>

	<MultipleSelect.Root id="11" bind:options={options2} bind:value={values.values2} required>
		<MultipleSelect.Viewer />
		<MultipleSelect.Controller>
			<MultipleSelect.Trigger />
			<MultipleSelect.Content>
				<MultipleSelect.Options>
					<MultipleSelect.Input />
					<MultipleSelect.List>
						<MultipleSelect.Empty>No results found.</MultipleSelect.Empty>
						<MultipleSelect.Group>
							{#each $options2 as option}
								<MultipleSelect.Item {option}>
									<Icon
										icon={option.icon ? option.icon : 'ph:empty'}
										class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
									/>
									{option.label}
									<MultipleSelect.Check {option} />
								</MultipleSelect.Item>
							{/each}
						</MultipleSelect.Group>
					</MultipleSelect.List>
					<MultipleSelect.Actions>
						<MultipleSelect.ActionAll>All</MultipleSelect.ActionAll>
						<MultipleSelect.ActionClear>Clear</MultipleSelect.ActionClear>
					</MultipleSelect.Actions>
				</MultipleSelect.Options>
			</MultipleSelect.Content>
		</MultipleSelect.Controller>
	</MultipleSelect.Root>
	<LayeredSingleSelect.Root id="12" bind:value={values.value9} options={options3} required>
		<LayeredSingleSelect.Trigger />
		<LayeredSingleSelect.Content>
			<LayeredSingleSelect.Group>
				{#each options3 as option}
					{#if option.subOptions && option.subOptions.length > 0}
						{#snippet Branch(
							options: LayeredSingleSelect.OptionType[],
							option: LayeredSingleSelect.OptionType,
							parents: LayeredSingleSelect.OptionType[],
							level: number = 0,
						)}
							<LayeredSingleSelect.Sub>
								<LayeredSingleSelect.SubTrigger>
									<Icon
										icon={option.icon && option.icon ? option.icon : 'ph:empty'}
										class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
									/>
									{option.label}
								</LayeredSingleSelect.SubTrigger>
								<LayeredSingleSelect.SubContent>
									{#each options as option}
										{#if option.subOptions && option.subOptions.length > 0}
											{@render Branch(option.subOptions, option, [...parents, option], level + 1)}
										{:else}
											<LayeredSingleSelect.Item {option} {parents}>
												<Icon
													icon={option.icon && option.icon ? option.icon : 'ph:empty'}
													class={cn(
														'size-5',
														option.icon && option.icon ? 'visibale' : 'invisible',
													)}
												/>
												{option.label}
												<LayeredSingleSelect.Check {option} {parents} />
											</LayeredSingleSelect.Item>
										{/if}
									{/each}
								</LayeredSingleSelect.SubContent>
							</LayeredSingleSelect.Sub>
						{/snippet}
						{@render Branch(option.subOptions, option, [option])}
					{:else}
						<LayeredSingleSelect.Item {option}>
							<Icon
								icon={option.icon && option.icon ? option.icon : 'ph:empty'}
								class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
							/>
							{option.label}
							<LayeredSingleSelect.Check {option} />
						</LayeredSingleSelect.Item>
					{/if}
				{/each}
			</LayeredSingleSelect.Group>
		</LayeredSingleSelect.Content>
	</LayeredSingleSelect.Root>

	<LayeredMultipleSelect.Root id="13" bind:value={values.values3} options={options4} required>
		<LayeredMultipleSelect.Viewer />
		<LayeredMultipleSelect.Controller>
			<LayeredMultipleSelect.Trigger />
			<LayeredMultipleSelect.Content>
				<LayeredMultipleSelect.Group>
					{#each options4 as option}
						{#if option.subOptions && option.subOptions.length > 0}
							{#snippet Branch(
								options: LayeredMultipleSelect.OptionType[],
								option: LayeredMultipleSelect.OptionType,
								parents: LayeredMultipleSelect.OptionType[],
								level: number = 0,
							)}
								<LayeredMultipleSelect.Sub>
									<LayeredMultipleSelect.SubTrigger>
										<Icon
											icon={option.icon && option.icon ? option.icon : 'ph:empty'}
											class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
										/>
										{option.label}
									</LayeredMultipleSelect.SubTrigger>
									<LayeredMultipleSelect.SubContent>
										{#each options as option}
											{#if option.subOptions && option.subOptions.length > 0}
												{@render Branch(
													option.subOptions,
													option,
													[...parents, option],
													level + 1,
												)}
											{:else}
												<LayeredMultipleSelect.Item {option} {parents}>
													<Icon
														icon={option.icon && option.icon ? option.icon : 'ph:empty'}
														class={cn(
															'size-5',
															option.icon && option.icon ? 'visibale' : 'invisible',
														)}
													/>
													{option.label}
													<LayeredMultipleSelect.Check {option} {parents} />
												</LayeredMultipleSelect.Item>
											{/if}
										{/each}
									</LayeredMultipleSelect.SubContent>
								</LayeredMultipleSelect.Sub>
							{/snippet}

							{@render Branch(option.subOptions, option, [option])}
						{:else}
							<LayeredMultipleSelect.Item {option}>
								<Icon
									icon={option.icon && option.icon ? option.icon : 'ph:empty'}
									class={cn('size-5', option.icon && option.icon ? 'visibale' : 'invisible')}
								/>
								{option.label}
								<LayeredMultipleSelect.Check {option} />
							</LayeredMultipleSelect.Item>
						{/if}
					{/each}
				</LayeredMultipleSelect.Group>
				<LayeredMultipleSelect.Actions>
					<LayeredMultipleSelect.ActionAll>All</LayeredMultipleSelect.ActionAll>
					<LayeredMultipleSelect.ActionClear>Clear</LayeredMultipleSelect.ActionClear>
				</LayeredMultipleSelect.Actions>
			</LayeredMultipleSelect.Content>
		</LayeredMultipleSelect.Controller>
	</LayeredMultipleSelect.Root>
</main>
