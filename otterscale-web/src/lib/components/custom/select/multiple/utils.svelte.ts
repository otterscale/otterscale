import type { AccessorType, OptionType } from './types';

class OptionManager {
    visibility = 1;
    options = $state([] as OptionType[]);
    accessor: AccessorType

    constructor(
        options: OptionType[],
        accessor: AccessorType,
    ) {
        this.options = options;
        this.accessor = accessor;
    }

    updateOptions(newOptions: OptionType[]) {
        this.options = newOptions
    }

    get selectedOptions(): OptionType[] {
        return this.options.filter((option) => (this.accessor.value.includes(option.value)))
    }

    get isSomeOptionsSelected(): boolean {
        return this.accessor.value.length > 0;
    }

    isOptionSelected(option: OptionType): boolean {
        return this.accessor.value.includes(option.value);
    }

    handleSelect(option: OptionType) {
        if (this.isOptionSelected(option)) {
            this.accessor.value = this.accessor.value.filter((o) => o !== option.value);
        } else {
            this.accessor.value = [...this.accessor.value, option.value];
        }
    }

    all() {
        this.accessor.value = this.options.map((option) => (option.value));
    }

    clear() {
        this.accessor.value = [];
    }
}

function validate(required: boolean | undefined, optionManager: OptionManager) {
    return required && !optionManager.isSomeOptionsSelected
}

export {
    OptionManager, validate
};
