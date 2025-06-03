import type { OptionType, valueSetterType } from './types';

class OptionManager {
    selectedOption: OptionType;
    valueSetter: valueSetterType;

    constructor(selectedOption: OptionType, valueSetter: valueSetterType) {
        this.selectedOption = $state(selectedOption ?? {} as OptionType);
        this.valueSetter = valueSetter
    }

    handleSelect(option: OptionType) {
        this.selectedOption = option;
        this.valueSetter(this.selectedOption);
    }
}

export {
    OptionManager
}