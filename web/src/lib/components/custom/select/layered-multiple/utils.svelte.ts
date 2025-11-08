import { getAllAncestralOptions, getAncestralOptionsMap } from '../layered-single/index';

import type { AccessorType, AncestralOptionType, OptionType } from './types';

class OptionManager {
	visibility = 1;

	options: OptionType[];
	ancestralOptionsMap: Record<string, OptionType[]>;
	accessor: AccessorType;

	constructor(options: OptionType[], accessor: AccessorType) {
		this.ancestralOptionsMap = getAncestralOptionsMap(options);

		this.options = options;
		this.accessor = accessor;
	}

	get isSomeAncestralOptionsSelected(): boolean {
		return this.accessor.value.length > 0;
	}

	get selectedAncestralOptions(): AncestralOptionType[] {
		return this.accessor.value.map((value) => this.ancestralOptionsMap[JSON.stringify(value)]);
	}

	hasSubOptions(option: OptionType): boolean {
		if (!option.subOptions) {
			return false;
		}
		if (option.subOptions.length == 0) {
			return false;
		}
		return true;
	}

	isOptionSelected(option: OptionType, parents: OptionType[]): boolean {
		return this.accessor.value.some((value) => {
			return (
				JSON.stringify(value) ===
				JSON.stringify([...parents.map((parent) => parent.value), option.value])
			);
		});
	}

	handleSelect(option: OptionType, parents: OptionType[]) {
		const values = this.accessor.value;
		const candidateValue: any[] = [...parents.map((parent) => parent.value), option.value];

		if (this.isOptionSelected(option, parents)) {
			this.accessor.value = values.filter((v) => {
				return JSON.stringify(v) !== JSON.stringify(candidateValue);
			});
		} else {
			this.accessor.value = [...values, candidateValue];
		}
	}

	all() {
		const all = getAllAncestralOptions(this.options).map((ancestralOption) =>
			ancestralOption.map((option) => option.value)
		);
		this.accessor.value = all;
	}

	clear() {
		this.accessor.value = [];
	}
}

function validate(required: boolean | undefined, optionManager: OptionManager) {
	return required && !optionManager.isSomeAncestralOptionsSelected;
}

export { OptionManager, validate };
