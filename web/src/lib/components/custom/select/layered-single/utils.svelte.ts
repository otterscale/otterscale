import type { OptionType, AncestralOptionType, valueSetterType } from './types';


class OptionManager {
    selectedAncestralOption: AncestralOptionType;
    valueSetter: valueSetterType;

    constructor(selectedAncestralOption: AncestralOptionType, valueSetter: valueSetterType) {
        this.selectedAncestralOption = $state(selectedAncestralOption)
        this.valueSetter = valueSetter;
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
        const selectedAncestralOption: AncestralOptionType = [...parents, option];
        return JSON.stringify(this.selectedAncestralOption) === JSON.stringify(selectedAncestralOption);
    }

    handleSelect(option: OptionType, parents: OptionType[]) {
        const selectedAncestralOption: AncestralOptionType = [...parents, option];
        this.selectedAncestralOption = selectedAncestralOption
        this.valueSetter(this.selectedAncestralOption);
    }
}

export {
    OptionManager
}