import constants from "./constants.js";

function identity(x) {
    return x;
}
function neg(x) {
    return -x;
}
function add(x, y) {
    return x + y;
}
function sub(x, y) {
    return x - y;
}
function mul(x, y) {
    return x * y;
}

function div(x, y) {
    if (y === 0) {
        return constants.DIVIDE_ZERO;
    }
    return x / y;
}

function pow(x, y) {
    return Math.pow(x, y);
}

const unary_operation = Object.freeze({
    [constants.ADD]: identity,
    [constants.SUBTRACT]: neg,
});

const binary_operation = Object.freeze({
    [constants.ADD]: add,
    [constants.SUBTRACT]: sub,
    [constants.IMPLIED_MULTIPLY]: mul,
    [constants.MULTIPLY]: mul,
    [constants.MULTIPLY_ALT]: mul,
    [constants.DIVIDE]: div,
    [constants.DIVIDE_ALT]: div,
    [constants.EXPONENT]: pow,
});

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
    unary_operation,
    binary_operation,
    is_space,
    is_decimal,
    is_open_paren,
    is_close_paren,
    is_paren,
    is_zero,
    is_digit,
    is_operator,
});
