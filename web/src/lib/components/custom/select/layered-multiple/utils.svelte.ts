import { getAllAncestralOptions, getAncestralOptionsMap } from '../layered-single/index';
import type { AncestralOptionType, OptionType, valuesGetterType, valuesSetterType } from './types';

class OptionManager {
    visibility = 1

    options: OptionType[];
    ancestralOptionsMap: Record<string, OptionType[]>

    valuesSetter: valuesSetterType;
    valuesGetter: valuesGetterType;

    constructor(options: OptionType[], valuesSetter: valuesSetterType, valuesGetter: valuesGetterType) {
        this.ancestralOptionsMap = getAncestralOptionsMap(options);

        this.options = options;
        this.valuesSetter = valuesSetter;
        this.valuesGetter = valuesGetter;
    }

    get isSomeAncestralOptionsSelected(): boolean {
        return this.valuesGetter().length > 0;
    }

    get selectedAncestralOptions(): AncestralOptionType[] {
        return this.valuesGetter().map((value) => (this.ancestralOptionsMap[JSON.stringify(value)]))
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
        return this.valuesGetter().some(
            (value) => {
                return JSON.stringify(value) === JSON.stringify([...parents.map((parent) => (parent.value)), option.value])
            }
        );
    }

    handleSelect(option: OptionType, parents: OptionType[]) {
        const values = this.valuesGetter()
        const candidateValue: any[] = [...parents.map((parent) => (parent.value)), option.value]

        if (this.isOptionSelected(option, parents)) {
            this.valuesSetter(values.filter((v) => {
                return JSON.stringify(v) !== JSON.stringify(candidateValue)
            }));
        } else {
            this.valuesSetter([...values, candidateValue])
        }
    }

    all() {
        const all = getAllAncestralOptions(this.options).map((ancestralOption) => (ancestralOption.map((option) => (option.value))))
        this.valuesSetter(all);
    }

    clear() {
        this.valuesSetter([]);
    }
}

export {
    OptionManager
};
