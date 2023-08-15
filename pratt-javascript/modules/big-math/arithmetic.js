import constants from "./constants.js";
import encoders from "./encoders.js";
import utils from "./utils.js";

function neg(x) {
    return encoders.make_bigfloat(-x.coefficient, x.exponent);
}
function abs(x) {
    return x.coefficient < constants.BIGINT_ZERO ? neg(x) : x;
}
function conform_operation(op) {
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

const add = conform_operation((x, y) => x + y);
const sub = conform_operation((x, y) => x - y);

function eq(x, y) {
    return utils.is_zero(sub(x, y));
}
function lt(x, y) {
    return utils.is_negative(sub(x, y));
}
function gt(x, y) {
    return utils.is_positive(sub(x, y));
}

function mul(x, y) {
    return encoders.make_bigfloat(
        x.coefficient * y.coefficient,
        x.exponent + y.exponent,
    );
}

function div(x, y) {
    if (utils.is_zero(x)) {
        return constants.ZERO;
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
