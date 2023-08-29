import constants from "./constants.js";
import arithmetic from "./big-math/arithmetic.js";

function identity(x) {
    return x;
}
function neg(x) {
    return arithmetic.neg(x);
}
function add(x, y) {
    return arithmetic.add(x, y);
}
function sub(x, y) {
    return arithmetic.sub(x, y);
}
function mul(x, y) {
    return arithmetic.mul(x, y);
}

function div(x, y) {
    return arithmetic.div(x, y);
}

function pow(x, y) {
    return arithmetic.pow(x, y);
}

const unary_operation = Object.freeze({
    [constants.ADD]: identity,
    [constants.SUBTRACT]: neg,
    [constants.SUBTRACT_ALT]: neg,
});

const binary_operation = Object.freeze({
    [constants.ADD]: add,
    [constants.SUBTRACT]: sub,
    [constants.SUBTRACT_ALT]: sub,
    [constants.IMPLIED_MULTIPLY]: mul,
    [constants.MULTIPLY]: mul,
    [constants.MULTIPLY_ALT]: mul,
    [constants.DIVIDE]: div,
    [constants.DIVIDE_ALT]: div,
    [constants.EXPONENT]: pow,
});

function is_exceeding(x) {
    const EXPONENT_LIMIT = 21;
    const COEFFICIENT_LIMIT = 1000000000000000000000n;
    return Math.abs(x.exponent) >= EXPONENT_LIMIT ||
        x.coefficient >= COEFFICIENT_LIMIT;
}

const is_space = (function () {
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

function is_decimal(x) {
    return x === constants.DECIMAL_POINT;
}

function is_open_paren(x) {
    return x === constants.OPEN_PAREN;
}

function is_close_paren(x) {
    return x === constants.CLOSE_PAREN;
}

const is_paren = (function () {
    const set = new Set([constants.OPEN_PAREN, constants.CLOSE_PAREN]);
    return function (x) {
        return set.has(x);
    };
})();

function is_exponent(x) {
    return x === constants.UPPER_E || x === constants.LOWER_E;
}

const is_plus_minus = (function () {
    const set = new Set([
        constants.ADD,
        constants.SUBTRACT,
        constants.SUBTRACT_ALT,
    ]);
    return function (x) {
        return set.has(x);
    };
})();

function is_zero(x) {
    return x === constants.ZERO;
}

const is_digit = (function () {
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

const is_operator = (function () {
    const set = new Set([
        constants.ADD,
        constants.SUBTRACT,
        constants.SUBTRACT_ALT,
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

export default Object.freeze({
    is_exceeding,
    unary_operation,
    binary_operation,
    is_space,
    is_decimal,
    is_open_paren,
    is_close_paren,
    is_paren,
    is_exponent,
    is_plus_minus,
    is_zero,
    is_digit,
    is_operator,
});
