<script lang="ts">
	import { shortcut } from '$lib/actions/shortcut.svelte';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import * as ContextMenu from '$lib/components/ui/context-menu/index.js';
	import * as Item from '$lib/components/ui/item';
	import * as Kbd from '$lib/components/ui/kbd/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import {
		Content,
		createForm,
		type FailureValidationResult,
		Form,
		type FormState,
		type FormValue,
		type FormValueValidator,
		getPseudoIdByPath,
		getValueSnapshot,
		type Schema,
		setFormContext,
		setValue,
		SubmitButton,
		type UiSchemaRoot,
		type ValidationResult,
		type ValidatorFactoryOptions
	} from '@sjsf/form';
	import { chain, fromFactories, fromRecord, overrideByRecord } from '@sjsf/form/lib/resolver';
	import lodash from 'lodash';
	import { mode as themeMode } from 'mode-watcher';
	import { onMount, tick } from 'svelte';
	import Monaco from 'svelte-monaco';
	import { schema as data } from './schema';

	import * as defaults from '$lib/components/dynamic-form/defaults';

	import ObjectPropertyField from '$lib/components/dynamic-form/fields/object-property.svelte';
	import ArrayItemTemplate from '$lib/components/dynamic-form/templates/array-item.svelte';
	import ArrayTemplate from '$lib/components/dynamic-form/templates/array.svelte';
	import MultiFieldTemplate from '$lib/components/dynamic-form/templates/multi-field.svelte';
	import ObjectPropertyTemplate from '$lib/components/dynamic-form/templates/object-property.svelte';
	import ObjectTemplate from '$lib/components/dynamic-form/templates/object.svelte';
	import {
		openAPISchemaToJSONSchema,
		toVersionedJSONSchema
	} from '$lib/components/dynamic-form/utils.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { FileCodeCornerIcon, FormIcon, LocateFixedIcon } from '@lucide/svelte';
	import type { SchemaObjectValue, SchemaValue } from '@sjsf/form/core';
	import { createFocusOnFirstError } from '@sjsf/form/focus-on-first-error';
	import { toast } from 'svelte-sonner';
	import { parse, stringify } from 'yaml';

	// Clean schema from unnecessary keywords and convert OpenAPI schema to JSON Schema Draft-07, which is compatible with AJV and most JSON Schema validators. This step is crucial to ensure that the form can be generated correctly without running into issues caused by unsupported keywords or schema versions.
	const apiResourceSchema = toVersionedJSONSchema(openAPISchemaToJSONSchema(data), 'draft-07');

	// Form
	const schema: Schema = {
		// ...(lodash.omit(apiResourceSchema, ['$schema', 'properties', 'description']) as any),
		properties: {
			metadata: {
				...(lodash.omit(lodash.get(apiResourceSchema, 'properties.metadata') as any, [
					'properties'
				]) as any),
				title: 'Metadata',
				properties: {
					name: {
						...(lodash.get(apiResourceSchema, 'properties.metadata.properties.name') as any),
						title: 'Name'
					}
				}
			},
			spec: {
				...(lodash.omit(lodash.get(apiResourceSchema, 'properties.spec') as any, [
					'description',
					'properties'
				]) as any),
				title: 'Specification',
				properties: {
					namespace: {
						...(lodash.get(apiResourceSchema, 'properties.spec.properties.namespace') as any),
						title: 'Namespace'
					},
					resourceQuota: {
						...(lodash.omit(
							lodash.get(apiResourceSchema, 'properties.spec.properties.resourceQuota'),
							'properties'
						) as any),
						title: 'Resource Quota',
						properties: {
							hard: {
								...(lodash.get(
									apiResourceSchema,
									'properties.spec.properties.resourceQuota.properties.hard'
								) as any),
								additionalProperties: {
									...(lodash.get(
										apiResourceSchema,
										'properties.spec.properties.resourceQuota.properties.hard.additionalProperties'
									) as any)
								}
							}
						}
					},
					limitRange: {
						...(lodash.get(apiResourceSchema, 'properties.spec.properties.limitRange') as any),
						title: 'Limit Range'
					},
					users: {
						...lodash.omit(
							lodash.get(apiResourceSchema, 'properties.spec.properties.users') as any,
							'items'
						),
						title: 'Users',
						items: {
							...lodash.omit(
								lodash.get(apiResourceSchema, 'properties.spec.properties.users.items') as any,
								'properties'
							),
							properties: {
								name: {
									...(lodash.get(
										apiResourceSchema,
										'properties.spec.properties.users.items.properties.name'
									) as any),
									enum: ['En-Yao', 'Chang', 'Elliott']
								},
								role: {
									...(lodash.get(
										apiResourceSchema,
										'properties.spec.properties.users.items.properties.role'
									) as any)
								},
								subject: {
									...(lodash.get(
										apiResourceSchema,
										'properties.spec.properties.users.items.properties.subject'
									) as any)
								}
							}
						}
					},
					networkIsolation: {
						...(lodash.get(
							apiResourceSchema,
							'properties.spec.properties.networkIsolation'
						) as any),
						title: 'Network Isolation'
					}
				}
			}
		}
	} as const;

	const uiSchema = {
		spec: {
			resourceQuota: {
				hard: {
					additionalProperties: {
						'ui:options': {
							translations: {
								'key-input-title': 'limit to'
							},
							hideTitle: true
						}
					}
				}
			},
			limitRange: {
				limits: {
					'ui:options': {
						itemTitle: () => 'limit'
					},
					items: {
						default: {
							additionalProperties: {
								'ui:options': {
									translations: {
										'key-input-title': 'default for'
									},
									hideTitle: true
								}
							}
						},
						defaultRequest: {
							additionalProperties: {
								'ui:options': {
									translations: {
										'key-input-title': 'default request for'
									},
									hideTitle: true
								}
							}
						},
						max: {
							additionalProperties: {
								'ui:options': {
									translations: {
										'key-input-title': 'maximum for'
									},
									hideTitle: true
								}
							}
						},
						maxLimitRequestRatio: {
							additionalProperties: {
								'ui:options': {
									translations: {
										'key-input-title': 'maximum limit request ratio for'
									},
									hideTitle: true
								}
							}
						},
						min: {
							additionalProperties: {
								'ui:options': {
									translations: {
										'key-input-title': 'minimum for'
									},
									hideTitle: true
								}
							}
						}
					}
				}
			},
			users: {
				'ui:options': {
					itemTitle: () => 'User'
				},
				items: {
					name: {
						'ui:components': {
							stringField: 'enumField',
							selectWidget: 'comboboxWidget'
						}
					},
					role: {
						'ui:components': {
							stringField: 'enumField',
							selectWidget: 'comboboxWidget'
						}
					}
				}
			},
			networkIsolation: {
				allowedNamespaces: {
					'ui:options': {
						itemTitle: () => 'Namespace'
					}
				}
			}
		}
	} satisfies UiSchemaRoot;

	const initialValue = {
		apiVersion: 'tenant.otterscale.io/v1alpha1',
		kind: 'Workspace',
		metadata: {
			name: ''
		},
		spec: {
			namespace: '',
			resourceQuota: {
				hard: {
					cpu: '',
					memory: '',
					vgpu: ''
				}
			},
			limitRange: {
				limits: [
					{
						default: {
							cpu: '',
							memory: '',
							vgpu: ''
						},
						defaultRequest: {
							cpu: '',
							memory: '',
							vgpu: ''
						},
						max: {
							cpu: '',
							memory: '',
							vgpu: ''
						},
						maxLimitRequestRatio: {
							cpu: '',
							memory: '',
							vgpu: ''
						},
						min: {
							cpu: '',
							memory: '',
							vgpu: ''
						},
						type: ''
					}
				]
			},
			users: [
				{
					name: '',
					role: ''
				}
			],
			networkIsolation: {
				allowedNamespaces: ['']
			}
		}
	};

	const sections = [
		{
			title: 'Name',
			path: ['metadata', 'name']
		},
		{
			title: 'Namespace',
			path: ['spec', 'namespace']
		},
		{
			title: 'Resource Quota',
			path: ['spec', 'resourceQuota']
		},
		{
			title: 'Limit Range',
			path: ['spec', 'limitRange']
		},
		{
			title: 'Users',
			path: ['spec', 'users']
		},
		{
			title: 'Network Isolation',
			path: ['spec', 'networkIsolation']
		}
	];

	function transfer(value: FormValue): FormValue {
		const temporaryValue = value as SchemaObjectValue;

		let users: SchemaObjectValue[] = lodash.get(value, 'spec.users', []) as SchemaObjectValue[];
		users = users.map((user) => ({
			...user,
			subject: user?.name
		})) as SchemaObjectValue[];
		lodash.set(temporaryValue, 'spec.users', users as SchemaValue[]);

		setValue(form, temporaryValue);
		return getValueSnapshot(form);
	}

	const extraUiOptions = chain(
		fromRecord({
			useLabel: false
		}),
		fromFactories({})
	);

	const theme = overrideByRecord(defaults.theme, {
		// Fields
		objectPropertyField: ObjectPropertyField,
		// Templates
		arrayItemTemplate: ArrayItemTemplate,
		arrayTemplate: ArrayTemplate,
		multiFieldTemplate: MultiFieldTemplate,
		objectPropertyTemplate: ObjectPropertyTemplate,
		objectTemplate: ObjectTemplate
	});

	let validationResult: ValidationResult<FormValue> | null = $state(null);
	function validator(options: ValidatorFactoryOptions) {
		const validator = defaults.validator<FormValue>(options);
		return {
			...validator,
			validateFormValue(schema: Schema, formValue: FormValue) {
				let value = {} as FormValue;
				switch (mode) {
					case 'yaml':
						value = parse(yamlValue);
						break;
					case 'form':
						value = formValue;
						break;
				}
				const transferredValue = transfer(value);
				validationResult = validator.validateFormValue(schema, transferredValue);
				if (validationResult && validationResult.errors && validationResult.errors.length > 0) {
					validationResult.errors.forEach((error) => {
						toast.error(error.message, {
							description: `[${error.path.join('.')}]`,
							duration: Number.POSITIVE_INFINITY,
							closeButton: true
						});
					});
				}
				return validationResult;
			}
		} satisfies FormValueValidator<FormValue>;
	}

	function onSubmit() {
		console.log(getValueSnapshot(form));
	}

	function onSubmitError(
		result: FailureValidationResult,
		event: SubmitEvent,
		form: FormState<FormValue>
	) {
		if (result.errors.length > 0) createFocusOnFirstError()(result, event, form);
	}

	const form = createForm<FormValue>({
		...defaults,
		theme,
		schema,
		uiSchema,
		extraUiOptions,
		initialValue,
		validator,
		onSubmit,
		onSubmitError
	});

	// Sections
	let activeSectionIndex = $state(0);
	function scrollTo(identifier: string, options?: ScrollIntoViewOptions) {
		const element = document.getElementById(identifier);
		if (element) {
			element.scrollIntoView({ behavior: 'smooth', block: 'start', ...options });
		}
	}
	function handleSectionNavigation(index: number) {
		const section = sections[index];
		if (section) {
			activeSectionIndex = index;
			scrollTo(getPseudoIdByPath(form, section.path, 'title'));
		}
	}
	function handleNextSection() {
		activeSectionIndex = Math.min(sections.length - 1, activeSectionIndex + 1);
		scrollTo(getPseudoIdByPath(form, sections[activeSectionIndex].path, 'title'));
	}
	function handlePreviousSection() {
		activeSectionIndex = Math.max(0, activeSectionIndex - 1);
		scrollTo(getPseudoIdByPath(form, sections[activeSectionIndex].path, 'title'));
	}

	onMount(() => {
		const sectionElements = sections
			.map((section) => {
				const id = getPseudoIdByPath(form, section.path, 'title');
				return document.getElementById(id);
			})
			.filter(Boolean) as HTMLElement[];

		const observer = new IntersectionObserver(
			(entries) => {
				const visibleEntries = entries.filter((entry) => entry.isIntersecting);
				if (visibleEntries.length > 0) {
					const topVisibleEntry = visibleEntries.reduce((previous, next) => {
						return next.boundingClientRect.top < previous.boundingClientRect.top ? next : previous;
					});
					const index = sectionElements.findIndex((element) => element === topVisibleEntry.target);
					if (index !== -1) {
						activeSectionIndex = index;
					}
				}
			},
			{ threshold: 0.5 }
		);

		sectionElements.forEach((element) => observer.observe(element));

		return () => observer.disconnect();
	});

	// YAML
	// Reorder attributes in YAML editor to match the form schema, making it more intuitive for users to find and edit values.
	// This is achieved by creating a new object based on the form schema and populating it with values from the current form state, ensuring that the order of attributes in the YAML editor reflects the structure defined in the form schema.
	setValue(form, parse(stringify(getValueSnapshot(form))));
	let yamlValue = $state(stringify(getValueSnapshot(form)));

	function onReady(event: CustomEvent) {
		const editor = event.detail;
	}

	// Tab
	async function synchronizeToYAML() {
		yamlValue = stringify(getValueSnapshot(form));
	}
	async function synchronizeToForm() {
		setValue(form, parse(yamlValue));
		await tick();
		setValue(form, parse(yamlValue));
	}

	let mode = $state('form');
	async function changeMode(targetMode: string) {
		try {
			switch (targetMode) {
				case 'yaml':
					await synchronizeToYAML();
					break;
				case 'form':
					await synchronizeToForm();
					break;
			}
			mode = targetMode;
			toast.success(`Switched to ${mode.toUpperCase()} mode`);
		} catch (error) {
			toast.error(
				`Failed to switch to ${targetMode.toUpperCase()} mode: ${(error as Error).message}`,
				{
					duration: 5000,
					closeButton: true
				}
			);
			return;
		}
	}

	setFormContext(form);
</script>

<svelte:window
	use:shortcut={{
		key: 'f',
		ctrl: true,
		callback: async () => {
			await changeMode('form');
		}
	}}
	use:shortcut={{
		key: 'y',
		ctrl: true,
		callback: async () => {
			await changeMode('yaml');
		}
	}}
	use:shortcut={{
		key: 'p',
		ctrl: true,
		callback: handlePreviousSection
	}}
	use:shortcut={{
		key: 'n',
		ctrl: true,
		callback: handleNextSection
	}}
/>
<Tabs.Root bind:value={mode} class="mx-auto max-w-3xl p-4">
	<Item.Root class="h-20 w-full p-0">
		<Item.Content class="text-left">
			<!-- Header -->
			<Item.Title class="text-lg font-bold">Workspace</Item.Title>
			<Item.Description class="text-sm">
				{apiResourceSchema.description}
			</Item.Description>
		</Item.Content>
		<Item.Actions>
			<!-- Mode Switcher -->
			<Tooltip.Provider>
				<ButtonGroup.Root>
					<Button
						size="icon-sm"
						onclick={async () => {
							await changeMode('form');
						}}
					>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<FormIcon />
							</Tooltip.Trigger>
							<Tooltip.Content class="flex items-center gap-1">
								Form
								<Kbd.Group>
									<Kbd.Root>ctrl</Kbd.Root>
									<Kbd.Root>F</Kbd.Root>
								</Kbd.Group>
							</Tooltip.Content>
						</Tooltip.Root>
					</Button>
					<Button
						size="icon-sm"
						onclick={async () => {
							await changeMode('yaml');
						}}
					>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<FileCodeCornerIcon />
							</Tooltip.Trigger>
							<Tooltip.Content class="flex items-center gap-1">
								YAML
								<Kbd.Group>
									<Kbd.Root>ctrl</Kbd.Root>
									<Kbd.Root>Y</Kbd.Root>
								</Kbd.Group>
							</Tooltip.Content>
						</Tooltip.Root>
					</Button>
				</ButtonGroup.Root>
			</Tooltip.Provider>
		</Item.Actions>
	</Item.Root>
	<Tabs.Content value="form">
		<ContextMenu.Root>
			<ContextMenu.Trigger>
				<!-- Form -->
				<Form attributes={{ novalidate: true }}>
					<Content />
					<SubmitButton />
				</Form>
			</ContextMenu.Trigger>
			<ContextMenu.Content>
				<ContextMenu.Group>
					<ContextMenu.GroupHeading>Navigation</ContextMenu.GroupHeading>
					<ContextMenu.Separator />
					{#each sections as section, index}
						<ContextMenu.Item
							disabled={activeSectionIndex === index}
							onclick={() => {
								handleSectionNavigation(index);
							}}
						>
							<LocateFixedIcon />{section.title}
							<ContextMenu.Shortcut>
								{#if activeSectionIndex - 1 === index}
									<Kbd.Group>
										<Kbd.Root>ctrl</Kbd.Root>
										<Kbd.Root>P</Kbd.Root>
									</Kbd.Group>
								{/if}
								{#if activeSectionIndex + 1 === index}
									<Kbd.Group>
										<Kbd.Root>ctrl</Kbd.Root>
										<Kbd.Root>N</Kbd.Root>
									</Kbd.Group>
								{/if}
							</ContextMenu.Shortcut>
						</ContextMenu.Item>
					{/each}
				</ContextMenu.Group>
			</ContextMenu.Content>
		</ContextMenu.Root>
	</Tabs.Content>
	<Tabs.Content value="yaml" class="h-[calc(100vh-7.5rem)]">
		<!-- YAML -->
		<Monaco
			bind:value={yamlValue}
			options={{
				automaticLayout: true,
				language: 'yaml',
				extraEditorClassName: 'h-full',
				folding: true,
				padding: { top: 24 },
				renderLineHighlight: 'all',
				theme: themeMode.current === 'dark' ? 'vs-dark' : 'vs-light'
			}}
			on:ready={onReady}
		/>
	</Tabs.Content>
</Tabs.Root>
