import type { OptionType, valuesSetterType } from './types';

class OptionManager {
    visibility = 1;
    options: OptionType[];
    selectedOptions: OptionType[] = $state([] as OptionType[]);
    valuesSetter: valuesSetterType;

    constructor(options: OptionType[], selectedOptions: OptionType[], valuesSetter: valuesSetterType) {
        this.options = options;
        this.selectedOptions = selectedOptions
        this.valuesSetter = valuesSetter;
    }

    get isSomeOptionsSelected() {
        return Array.isArray(this.selectedOptions) && this.selectedOptions.length > 0;
    }

    isOptionSelected(option: OptionType): boolean {
        return this.selectedOptions.map((o) => o.value).includes(option.value);
    }

    handleSelect(option: OptionType) {
        if (Array.isArray(this.selectedOptions) && this.isOptionSelected(option)) {
            this.selectedOptions = this.selectedOptions.filter((o) => o.value !== option.value);
        } else {
            this.selectedOptions = [...this.selectedOptions, option];
        }
        this.valuesSetter(this.selectedOptions);
    }

    all() {
        this.selectedOptions = [...this.options];
        this.valuesSetter(this.selectedOptions);
    }

    clear() {
        this.selectedOptions = [];
        this.valuesSetter(this.selectedOptions);
    }
}

export {
    OptionManager
}