import type { AccessorType, OptionType } from './types';

class OptionManager {
    options = $state([] as OptionType[]);
    accessor: AccessorType

    constructor(options: OptionType[], accessor: AccessorType) {
        this.options = options;
        this.accessor = accessor
    }

    updateOptions(newOptions: OptionType[]) {
        this.options = newOptions
    }

    get selectedOption(): OptionType {
        return this.options.find((option) => (JSON.stringify(option.value) === JSON.stringify(this.accessor.value))) ?? {} as OptionType
    }

    isOptionSelected(option: OptionType): boolean {
        return JSON.stringify(option.value) === JSON.stringify(this.accessor.value)
    }

    handleSelect(option: OptionType) {
        this.accessor.value = option.value;
    }
}

function validate(required: boolean | undefined, optionManager: OptionManager) {
    return required && !optionManager.selectedOption.value
}


export {
    OptionManager, validate
};
