import type { AncestralOptionType, OptionType, valueGetterType, valueSetterType } from './types';

const getAllAncestralOptions = (
    options: OptionType[],
    parents: OptionType[] = []
): AncestralOptionType[] => {
    return options.flatMap(option => {
        const ancestorOption = [...parents, option];

        if (option.subOptions) {
            return getAllAncestralOptions(option.subOptions!, [...parents, option]);
        } else {
            return [ancestorOption];
        }
    });
};

function getAncestralOptionsMap(
    options: OptionType[]
): Record<string, OptionType[]> {
    let allAncestralOptionsMap: Record<string, OptionType[]> = {};

    const allAncestralOptions = getAllAncestralOptions(options)
    allAncestralOptions.forEach((line) => {
        allAncestralOptionsMap[JSON.stringify(line.map((node) => node.value))] = line;
    });

    return allAncestralOptionsMap;
}

class OptionManager {
    options: OptionType[];
    ancestralOptionsMap: Record<string, OptionType[]>

    valueSetter: valueSetterType;
    valueGetter: valueGetterType;

    constructor(options: OptionType[], valueSetter: valueSetterType, valueGetter: valueGetterType) {
        this.ancestralOptionsMap = getAncestralOptionsMap(options);

        this.options = options
        this.valueSetter = valueSetter;
        this.valueGetter = valueGetter
    }

    get selectedAncestralOption(): AncestralOptionType {
        return this.ancestralOptionsMap[JSON.stringify(this.valueGetter())]
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
        return JSON.stringify(this.valueGetter()) === JSON.stringify([...parents.map((parent) => (parent.value)), option.value]);
    }

    handleSelect(option: OptionType, parents: OptionType[]) {
        this.valueSetter([...parents.map((parent) => (parent.value)), option.value]);
    }
}

export {
    getAllAncestralOptions, getAncestralOptionsMap, OptionManager
};
