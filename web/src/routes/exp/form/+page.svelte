<script lang="ts">
	import Label from '$lib/components/ui/label/label.svelte';
	import Icon from '@iconify/svelte';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { Multiple as MultipleInput } from '$lib/components/custom/input';
	import { Single as SingleSelect } from '$lib/components/custom/select';
	import { Multiple as MultipleSelect } from '$lib/components/custom/select';
	import { LayeredSingle as LayeredSingleSelect } from '$lib/components/custom/select';
	import { LayeredMultiple as LayeredMultipleSelect } from '$lib/components/custom/select';
	import { cn } from '$lib/utils';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Form from '$lib/components/custom/form';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { z } from 'zod';
	import { writable, type Writable } from 'svelte/store';

	let options1: Writable<SingleSelect.OptionType[]> = writable([
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star'
		},
		{
			value: 'sun',
			label: 'Sun',
			icon: 'ph:sun'
		}
	]);
	let options2: Writable<MultipleSelect.OptionType[]> = writable([
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star'
		},
		{
			value: 'sun',
			label: 'Sun',
			icon: 'ph:sun'
		},
		{
			value: 'cloud',
			label: 'Cloud',
			icon: 'ph:cloud'
		},
		{
			value: 'rainbow',
			label: 'Rainbow',
			icon: 'ph:rainbow'
		},
		{
			value: 'comet',
			label: 'Comet',
			icon: 'ph:comet'
		},
		{
			value: 'planet',
			label: 'Planet',
			icon: 'ph:planet'
		},
		{
			value: 'meteor',
			label: 'Meteor',
			icon: 'ph:shooting-star'
		},
		{
			value: 'sparkle',
			label: 'Sparkle',
			icon: 'ph:sparkle'
		},
		{
			value: 'nebula',
			label: 'Nebula',
			icon: 'ph:cloud-fog'
		}
	]);
	let options3: LayeredSingleSelect.OptionType[] = [
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star'
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
							icon: 'ph:clock'
						}
					]
				}
			]
		}
	];
	const options4: LayeredMultipleSelect.OptionType[] = [
		{ value: 'moon', label: 'Moon', icon: 'ph:moon' },
		{
			value: 'star',
			label: 'Star',
			icon: 'ph:star'
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
							icon: 'ph:clock'
						}
					]
				}
			]
		}
	];

	type Values = {
		value1: any;
		value2: any;
		value3: any;
		value4: any;
		value5: any;
		value6: any;
		value7: any;
		value8: any;
		value9: any;
		value10: any;
		values1: any[];
		values2: any[];
		values3: any[];
		values4: any[];
		values5: any[];
	};

	let values: Values = $state({} as Values);

	function reset() {
		values = {} as Values;
	}

	function ListenInputs(...values: any[]) {
		return values.some((value) => {
			if (value) {
				return true;
			}
		});
	}

	const stateController = new DialogStateController(true);

	import * as Enterprise from '$lib/components/custom/enterprise/index';
	let invalid = $state(false);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<Enterprise.Hint>
		<AlertDialog.Trigger class={cn('hover:cursor-pointer', buttonVariants({ variant: 'outline' }))}>
			Trigger
		</AlertDialog.Trigger>
	</Enterprise.Hint>
	<AlertDialog.Content>
		<AlertDialog.Header class="flex items-center justify-center text-xl font-bold"
			>Header</AlertDialog.Header
		>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Legend>Fieldset I</Form.Legend>
				<Form.Description>
					This is a long description for the form fieldset. It provides detailed information about
					the purpose of this section.
				</Form.Description>
				<Form.Field>
					<Form.Label>
						Field I
						{#snippet information()}
							This is a detailed description for the form field. It provides additional context
							about the purpose and usage of this field.
						{/snippet}
					</Form.Label>
					<SingleInput.General
						required
						schema={z.string().min(3)}
						type="text"
						id="1"
						bind:value={values.value1}
					/>
					<Form.Help>
						This is a help text for Field I. Please enter the required information.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Form.Label>
						Field II
						{#snippet information()}
							This is a new description for Field II. Please enter a number in this field as
							required.
						{/snippet}
					</Form.Label>
					<SingleInput.General required type="number" id="2" bind:value={values.value2} />
					<Form.Help>This is a help text for Field II. Please enter a valid number.</Form.Help>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Fieldset II</Form.Legend>
				<Form.Description>
					This is a detailed description for the form fieldset. It offers additional context about
					the purpose and usage of this section.
				</Form.Description>
				<Form.Field>
					<Label>Field III</Label>
					<SingleInput.Boolean required id="3" bind:value={values.value3} />
					<Form.Help>
						Enable this option if you want to activate Field III. This is a help text for the
						boolean input.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field IV</Label>
					<SingleInput.Password required id="4" bind:value={values.value4} />
					<Form.Help>Please enter your password. Make sure it is strong and secure.</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field V</Label>
					<SingleInput.Color id="5" required bind:value={values.value5} />
					<Form.Help>Select your favorite color using the color picker above.</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field VI</Label>
					<SingleInput.DeletionConfirm id="6" required target="value6" bind:value={values.value6} />
					<Form.Help>
						This field requires confirmation before deletion. Double check before proceeding.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field VII</Label>
					<SingleInput.Structure id="7" language="json" required bind:value={values.value7} />
					<Form.Help>
						Enter valid JSON structure. The input will validate your JSON format.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field VIII</Label>
					<SingleInput.Measurement
						id="8"
						bind:value={values.value8}
						required
						units={[
							{ value: 1, label: 'I' } as SingleInput.UnitType,
							{ value: 2, label: 'II' } as SingleInput.UnitType,
							{ value: 3, label: 'III' } as SingleInput.UnitType
						]}
					/>
					<Form.Help>
						Enter a value and select a unit of measurement from the dropdown options.
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Fieldset III</Form.Legend>
				<Form.Description>
					This section demonstrates a multiple input field. You can add, view, or clear multiple
					values as needed.
				</Form.Description>
				<Form.Field>
					<Label>Field IV</Label>
					<MultipleInput.Root type="number" required bind:values={values.values1} id="9">
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
					<Form.Help>
						This is a help text for Field VI. You can add multiple numbers, view them, or clear the
						list as needed.
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>

			<Form.Separator>Or</Form.Separator>

			<Form.Fieldset>
				<Form.Legend>Fieldset IV</Form.Legend>
				<Form.Description>
					This section demonstrates various select components, including single, multiple, and
					layered selects with nested options.
				</Form.Description>
				<Form.Field>
					<Label>Field V</Label>
					<SingleSelect.Root id="10" bind:options={options1} bind:value={values.value9} required>
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
					<Form.Help>
						Please select an option from the dropdown above. This is a help text for Field VII.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field XI</Label>
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
					<Form.Help>
						You can select multiple options from the dropdown above. This is a help text for Field
						VIII.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field XII</Label>
					<LayeredSingleSelect.Root id="12" bind:value={values.value10} options={options3} required>
						<LayeredSingleSelect.Trigger />
						<LayeredSingleSelect.Content>
							<LayeredSingleSelect.Group>
								{#each options3 as option}
									{#if option.subOptions && option.subOptions.length > 0}
										{#snippet Branch(
											options: LayeredSingleSelect.OptionType[],
											option: LayeredSingleSelect.OptionType,
											parents: LayeredSingleSelect.OptionType[],
											level: number = 0
										)}
											<LayeredSingleSelect.Sub>
												<LayeredSingleSelect.SubTrigger>
													<Icon
														icon={option.icon && option.icon ? option.icon : 'ph:empty'}
														class={cn(
															'size-5',
															option.icon && option.icon ? 'visibale' : 'invisible'
														)}
													/>
													{option.label}
												</LayeredSingleSelect.SubTrigger>
												<LayeredSingleSelect.SubContent>
													{#each options as option}
														{#if option.subOptions && option.subOptions.length > 0}
															{@render Branch(
																option.subOptions,
																option,
																[...parents, option],
																level + 1
															)}
														{:else}
															<LayeredSingleSelect.Item {option} {parents}>
																<Icon
																	icon={option.icon && option.icon ? option.icon : 'ph:empty'}
																	class={cn(
																		'size-5',
																		option.icon && option.icon ? 'visibale' : 'invisible'
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
					<Form.Help>
						Please select a nested option from the layered select above. This is a help text for
						Field IX.
					</Form.Help>
				</Form.Field>
				<Form.Field>
					<Label>Field XIII</Label>
					<LayeredMultipleSelect.Root
						id="13"
						bind:value={values.values3}
						options={options4}
						required
					>
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
												level: number = 0
											)}
												<LayeredMultipleSelect.Sub>
													<LayeredMultipleSelect.SubTrigger>
														<Icon
															icon={option.icon && option.icon ? option.icon : 'ph:empty'}
															class={cn(
																'size-5',
																option.icon && option.icon ? 'visibale' : 'invisible'
															)}
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
																	level + 1
																)}
															{:else}
																<LayeredMultipleSelect.Item {option} {parents}>
																	<Icon
																		icon={option.icon && option.icon ? option.icon : 'ph:empty'}
																		class={cn(
																			'size-5',
																			option.icon && option.icon ? 'visibale' : 'invisible'
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
													class={cn(
														'size-5',
														option.icon && option.icon ? 'visibale' : 'invisible'
													)}
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
					<Form.Help>
						You can select multiple nested options from the layered select above. This is a help
						text for Field X.
					</Form.Help>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Actions>
					<AlertDialog.ActionGroup>
						<AlertDialog.ActionGroupHeading>Group</AlertDialog.ActionGroupHeading>
						<AlertDialog.ActionItem
							onclick={() => {
								console.log(values);
							}}
						>
							Debug
							<AlertDialog.ActionShortcut>⌘S</AlertDialog.ActionShortcut>
						</AlertDialog.ActionItem>
					</AlertDialog.ActionGroup>
					<AlertDialog.ActionGroup>
						<AlertDialog.ActionGroupHeading>Group</AlertDialog.ActionGroupHeading>
						<AlertDialog.ActionItem
							onclick={() => {
								console.log('Action 2 clicked');
								stateController.close();
							}}
						>
							Action
							<AlertDialog.ActionShortcut>⌘S</AlertDialog.ActionShortcut>
						</AlertDialog.ActionItem>
						<AlertDialog.ActionItem
							onclick={() => {
								console.log('Action 3 clicked');
								stateController.close();
							}}
						>
							Action
							<AlertDialog.ActionShortcut>⌘S</AlertDialog.ActionShortcut>
						</AlertDialog.ActionItem>
					</AlertDialog.ActionGroup>
				</AlertDialog.Actions>
				<AlertDialog.Action disabled={invalid} onclick={reset}>Confirm</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
