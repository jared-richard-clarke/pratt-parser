import {
    ADD,
    CARRIAGE_RETURN,
    DIVIDE,
    EIGHT,
    EXPONENT,
    FIVE,
    FORM_FEED,
    FOUR,
    LINEFEED,
    MULTIPLY,
    NINE,
    ONE,
    SEVEN,
    SIX,
    SUBTRACT,
    TAB,
    THREE,
    TWO,
    VERTICAL_TAB,
    WHITE_SPACE,
    ZERO,
} from "./constants.js";

const mod = Object.create(null);

mod.fix_object = function (x) {
    return Object.freeze(
        Object.entries(x).reduce((accum, [key, value]) => {
            accum[key] = value;
            return accum;
        }, Object.create(null)),
    );
};

mod.is_space = (function () {
    const set = new Set([
        WHITE_SPACE,
        TAB,
        LINEFEED,
        CARRIAGE_RETURN,
        VERTICAL_TAB,
        FORM_FEED,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

mod.is_decimal = function (x) {
    return x === DECIMAL_POINT;
};

mod.is_digit = (function () {
    const set = new Set([
        ZERO,
        ONE,
        TWO,
        THREE,
        FOUR,
        FIVE,
        SIX,
        SEVEN,
        EIGHT,
        NINE,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

mod.is_operator = (function () {
    const set = new Set([ADD, SUBTRACT, MULTIPLY, DIVIDE, EXPONENT]);
    return function (x) {
        return set.has(x);
    };
})();

export default Object.freeze(mod);
