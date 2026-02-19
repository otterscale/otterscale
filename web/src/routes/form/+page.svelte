<script lang="ts">
	import { mode as themeMode } from 'mode-watcher';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as Item from '$lib/components/ui/item';
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import {
		Content,
		createForm,
		type FailureValidationResult,
		Form,
		type FormState,
		type FormValue,
		type FormValueValidator,
		getValueSnapshot,
		type Schema,
		setFormContext,
		setValue,
		SubmitButton,
		type UiSchemaRoot,
		type ValidationResult,
		type ValidatorFactoryOptions
	} from '@sjsf/form';
	import { overrideByRecord } from '@sjsf/form/lib/resolver';
	import { JSONSchemaFaker } from 'json-schema-faker';
	import lodash from 'lodash';
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
	import { CircleAlertIcon, FileCodeCornerIcon, FormIcon, SaveIcon } from '@lucide/svelte';
	import type { SchemaObjectValue, SchemaValue } from '@sjsf/form/core';
	import { createFocusOnFirstError } from '@sjsf/form/focus-on-first-error';
	import { toast } from 'svelte-sonner';
	import { parse, stringify } from 'yaml';

	JSONSchemaFaker.option({
		alwaysFakeOptionals: false,
		fillProperties: false,
		maxItems: 1,
		renderTitle: true,
		renderDescription: true,
		renderComment: true,
		requiredOnly: true,
		useDefaultValue: true
	});

	const apiResourceSchema = toVersionedJSONSchema(openAPISchemaToJSONSchema(data), 'draft-07');

	const schema: Schema = {
		...(lodash.omit(apiResourceSchema, ['$schema', 'properties', 'description']) as any),
		// title: 'Workspace',
		properties: {
			spec: {
				...(lodash.omit(lodash.get(apiResourceSchema, 'properties.spec'), [
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
						itemTitle: () => 'Limit'
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
		spec: {
			resourceQuota: {
				hard: {
					cpu: ''
				}
			},
			limitRange: {
				limits: [
					{
						default: {
							cpu: ''
						},
						defaultRequest: {
							cpu: ''
						},
						max: {
							cpu: ''
						},
						maxLimitRequestRatio: {
							cpu: ''
						},
						min: {
							cpu: ''
						},
						type: 'Container'
					}
				]
			},
			users: [
				{
					name: '',
					role: 'admin'
				}
			]
		}
	};

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

	function transfer(value: SchemaObjectValue): FormValue {
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

	let validationResult: ValidationResult<FormValue> | null = $state(null);
	function validator(options: ValidatorFactoryOptions) {
		const validator = defaults.validator<FormValue>(options);
		return {
			...validator,
			validateFormValue(schema: Schema, formValue: FormValue) {
				const value = mode === 'yaml' ? parse(yamlValue) : formValue;
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
		initialValue,
		validator,
		onSubmit,
		onSubmitError
	});

	function scrollTo(identifier: string, options?: ScrollIntoViewOptions) {
		const element = document.getElementById(identifier);
		if (element) {
			element.scrollIntoView({ behavior: 'smooth', block: 'start', ...options });
		}
	}

	let yamlValue = $state(stringify(getValueSnapshot(form)));

	let mode = $state('form');

	function synchronizeToYAML() {
		yamlValue = stringify(getValueSnapshot(form));
	}
	function synchronizeToForm() {
		setValue(form, parse(yamlValue));
	}

	let isYAMLEditing = $state(false);
	function onReady(event: CustomEvent) {
		const yamlEditor = event.detail;
		yamlEditor.onDidChangeModelContent(() => {
			if (mode === 'yaml') {
				if (!isYAMLEditing) isYAMLEditing = true;
			}
		});
	}
	function handleYAMLSave() {
		if (mode === 'yaml') {
			try {
				synchronizeToForm();
				isYAMLEditing = false;
			} catch (error: any) {
				console.error(error);
				toast.error('Invalid YAML syntax.', {
					description: error.message.toString(),
					duration: Number.POSITIVE_INFINITY,
					closeButton: true
				});
			}
		}
	}

	setFormContext(form);
</script>

<!-- Avoid unnecessary nesting of elements to prevent layout issues and ensure better readability. -->
<main class="max-h-screen overflow-hidden">
	<Tabs.Root
		bind:value={mode}
		class="space-y-2 p-4"
		onValueChange={() => {
			if (mode === 'yaml') {
				synchronizeToYAML();
			} else if (mode === 'form') {
				synchronizeToForm();
			}
		}}
	>
		<Item.Root class="h-20 w-full p-0">
			<Item.Content class="text-left">
				<!-- Header -->
				<Item.Title class="text-lg font-bold">Workspace</Item.Title>
				<Item.Description class="text-sm">
					{apiResourceSchema.description}
				</Item.Description>
			</Item.Content>
			<Item.Actions>
				<Button size="icon" class={isYAMLEditing ? undefined : 'hidden'} onclick={handleYAMLSave}>
					<SaveIcon />
				</Button>
				<Tabs.List>
					<!-- Mode Switcher -->
					<Tabs.Trigger value="form" disabled={isYAMLEditing}>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<FormIcon />
								</Tooltip.Trigger>
								<Tooltip.Content>Form</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Tabs.Trigger>
					<Tabs.Trigger value="yaml">
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<FileCodeCornerIcon />
								</Tooltip.Trigger>
								<Tooltip.Content>YAML</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Tabs.Trigger>
				</Tabs.List>
			</Item.Actions>
		</Item.Root>
		<Form attributes={{ novalidate: true }}>
			<ScrollArea
				class="h-[calc(100vh-11rem)]"
				orientation={mode === 'yaml' ? 'horizontal' : 'both'}
			>
				<Tabs.Content value="form" class="h-screen">
					<!-- Form -->
					<!-- Force refresh: anyOfField may retain previous schema and clear value when schema changes -->
					<!-- {#key mode === 'form'} -->
					<Content />
					<!-- {/key} -->
				</Tabs.Content>
				<Tabs.Content value="yaml" class="flex h-screen gap-4">
					<!-- Code -->
					<Monaco
						bind:value={yamlValue}
						options={{
							automaticLayout: true,
							language: 'yaml',
							extraEditorClassName: 'h-full',
							folding: true,
							renderLineHighlight: 'all',
							theme: themeMode.current === 'dark' ? 'vs-dark' : 'vs-light'
						}}
						on:ready={onReady}
					/>
					<!-- Errors -->
					<!-- {#if validationResult && validationResult.errors && validationResult.errors.length > 0}
						{@const errors = validationResult.errors}
						<div class="fixed right-0 bottom-12 left-0 z-50 bg-card/90 p-4">
							<Alert.Root
								variant="destructive"
								class="max-h-[33vh] overflow-auto border-destructive bg-destructive/10"
							>
								<CircleAlertIcon />
								<Alert.Title>{errors?.length ?? 0} errors</Alert.Title>
								<Alert.Description>
									{#each errors as error}
										<p class="text-sm">[{error.path.join('.')}]: {error.message}</p>
									{/each}
								</Alert.Description>
							</Alert.Root>
						</div>
					{/if} -->
				</Tabs.Content>
			</ScrollArea>
			<!-- Submit -->
			<SubmitButton />
		</Form>
	</Tabs.Root>
</main>
