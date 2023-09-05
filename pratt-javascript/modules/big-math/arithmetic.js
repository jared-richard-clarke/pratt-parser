import constants from "./constants.js";
import encoders from "./encoders.js";
import utils from "./utils.js";

// Negates big-float "x". Positive if negative. Negative if positive.
function neg(x) {
    return constants.make_bigfloat(-x.coefficient, x.exponent);
}

// Truncates the fractional component of big-float "x".
function floor(x) {
    const { coefficient, exponent } = x;
    if (exponent === 0) {
        return x;
    }
    if (exponent > 0) {
        return constants.make_bigfloat(
            coefficient * constants.BIGINT_TEN ** BigInt(exponent),
            0,
        );
    }
    return constants.make_bigfloat(
        coefficient / constants.BIGINT_TEN ** BigInt(-exponent),
        0,
    );
}

// Adjusts the components of big-float "x" and big-float "y"
// so that they can be added or subtracted.
function match_terms(op) {
    return function (x, y) {
        const differential = x.exponent - y.exponent;
        if (differential === 0) {
            return constants.make_bigfloat(
                op(x.coefficient, y.coefficient),
                x.exponent,
            );
        }
        if (differential > 0) {
            return constants.make_bigfloat(
                op(
                    x.coefficient *
                        constants.BIGINT_TEN ** BigInt(differential),
                    y.coefficient,
                ),
                y.exponent,
            );
        }
        return constants.make_bigfloat(
            op(
                x.coefficient,
                y.coefficient * constants.BIGINT_TEN ** BigInt(-differential),
            ),
            x.exponent,
        );
    };
}

// big-float "x" + big-float "y".
const add = match_terms((x, y) => x + y);
// big-float "x" - big-float "y".
const sub = match_terms((x, y) => x - y);

// big-float "x" = big-float "y".
function eq(x, y) {
    x = encoders.normalize(x);
    y = encoders.normalize(y);
    return (x.coefficient === y.coefficient) && (x.exponent === y.exponent);
}
// big-float "x" >= big-float "y"
function ge(x, y) {
    const difference = sub(x, y);
    return difference.coefficient >= constants.BIGINT_ZERO;
}
// big-float "x" > big-float "y"
function gt(x, y) {
    const difference = sub(x, y);
    return difference.coefficient > constants.BIGINT_ZERO;
}
// big-float "x" <= big-float "y"
function le(x, y) {
    const difference = sub(x, y);
    return difference.coefficient <= constants.BIGINT_ZERO;
}
// big-float "x" < big-float "y"
function lt(x, y) {
    const difference = sub(x, y);
    return difference.coefficient < constants.BIGINT_ZERO;
}

// big-float "x" × big-float "y"
function mul(x, y) {
    return constants.make_bigfloat(
        x.coefficient * y.coefficient,
        x.exponent + y.exponent,
    );
}

// big-float "x" ÷ big-float "y"
// Division by zero is undefined.
function div(x, y) {
    if (utils.is_zero(x)) {
        return constants.BIGFLOAT_ZERO;
    }
    if (utils.is_zero(y)) {
        return constants.DIVIDE_ZERO;
    }
    const precision = constants.PRECISION;
    let { coefficient, exponent } = x;
    exponent -= y.exponent;
    if (exponent > precision) {
        coefficient = coefficient *
            constants.BIGINT_TEN ** BigInt(exponent - precision);
        exponent = precision;
    }
    coefficient = coefficient / y.coefficient;
    const remainder = coefficient % y.coefficient;
    // round
    if (
        utils.bigint_abs(remainder + remainder) >=
            utils.bigint_abs(y.coefficient)
    ) {
        coefficient += utils.bigint_signum(x.coefficient);
    }
    return constants.make_bigfloat(coefficient, exponent);
}

// big-float "x" ^ big-float "y"
// Performs exponentiation by squaring.
function pow(x, y) {
    // Current implementation does not support non-integer exponents.
    if (encoders.normalize(y).exponent !== 0) {
        return constants.NON_INTEGER_EXPONENT;
    }
    if (utils.is_zero(x)) {
        return constants.BIGFLOAT_ZERO;
    }
    if (utils.is_zero(y)) {
        return constants.BIGFLOAT_ONE;
    }
    if (utils.is_negative(y)) {
        x = div(constants.BIGFLOAT_ONE, x);
        y = neg(y);
    }
    let z = constants.BIGFLOAT_ONE;
    while (gt(y, constants.BIGFLOAT_ONE)) {
        if (utils.is_odd(y)) {
            z = mul(x, z);
        }
        x = mul(x, x);
        y = floor(div(y, constants.BIGFLOAT_TWO));
    }
    return mul(x, z);
}

export default Object.freeze({
    eq,
    ge,
    gt,
    le,
    lt,
    neg,
    add,
    sub,
    mul,
    div,
    pow,
});
