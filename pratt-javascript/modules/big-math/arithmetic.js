import constants from "./constants.js";
import encoders from "./encoders.js";
import utils from "./utils.js";

function equal(x, y) {
    x = encoders.normalize(x);
    y = encoders.normalize(y);
    return (x.coefficient === y.coefficient) && (x.exponent === y.exponent);
}
function neg(x) {
    return encoders.make_bigfloat(-x.coefficient, x.exponent);
}
function floor(x) {
    const { coefficient, exponent } = x;
    if (exponent === 0) {
        return x;
    }
    if (exponent > 0) {
        return encoders.make_bigfloat(
            coefficient * constants.BIGINT_TEN ** BigInt(exponent),
            0,
        );
    }
    return encoders.make_bigfloat(
        coefficient / constants.BIGINT_TEN ** BigInt(-exponent),
        0,
    );
}
function adjust_terms(op) {
    return function (x, y) {
        const differential = x.exponent - y.exponent;
        if (differential === 0) {
            return encoders.make_bigfloat(
                op(x.coefficient, y.coefficient),
                x.exponent,
            );
        }
        if (differential > 0) {
            return encoders.make_bigfloat(
                op(
                    x.coefficient *
                        constants.BIGINT_TEN ** BigInt(differential),
                    y.coefficient,
                ),
                y.exponent,
            );
        }
        return encoders.make_bigfloat(
            op(
                x.coefficient,
                y.coefficient * constants.BIGINT_TEN ** BigInt(-differential),
            ),
            x.exponent,
        );
    };
}

const add = adjust_terms((x, y) => x + y);
const sub = adjust_terms((x, y) => x - y);

function gt(x, y) {
    const difference = sub(x, y);
    return difference.coefficient > constants.BIGINT_ZERO;
}

function mul(x, y) {
    return encoders.make_bigfloat(
        x.coefficient * y.coefficient,
        x.exponent + y.exponent,
    );
}

// - Division by zero is undefined.
function div(x, y) {
    if (utils.is_zero(x)) {
        return constants.BIGFLOAT_ZERO;
    }
    if (utils.is_zero(y)) {
        return undefined;
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
    return encoders.make_bigfloat(coefficient, exponent);
}

// Exponentiation by squaring: efficiently computes integer exponents.
function pow(x, y) {
    // Current implementation does not support non-integer exponents.
    if (encoders.normalize(y).exponent !== 0) {
        return undefined;
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
