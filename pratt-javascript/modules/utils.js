import constants from "./constants.js";

const mod = Object.create(null);

mod.is_space = (function () {
    const set = new Set([
        constants.WHITE_SPACE,
        constants.TAB,
        constants.LINEFEED,
        constants.CARRIAGE_RETURN,
        constants.VERTICAL_TAB,
        constants.FORM_FEED,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

mod.is_decimal = function (x) {
    return x === constants.DECIMAL_POINT;
};

mod.is_paren = (function () {
    const set = new Set([constants.OPEN_PAREN, constants.CLOSE_PAREN]);
    return function (x) {
        return set.has(x);
    };
})();

mod.is_zero = function (x) {
    return x === constants.ZERO;
};

mod.is_digit = (function () {
    const set = new Set([
        constants.ZERO,
        constants.ONE,
        constants.TWO,
        constants.THREE,
        constants.FOUR,
        constants.FIVE,
        constants.SIX,
        constants.SEVEN,
        constants.EIGHT,
        constants.NINE,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

mod.is_operator = (function () {
    const set = new Set([
        constants.ADD,
        constants.SUBTRACT,
        constants.MULTIPLY,
        constants.MULTIPLY_ALT,
        constants.DIVIDE,
        constants.DIVIDE_ALT,
        constants.EXPONENT,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

export default Object.freeze(mod);
