import constants from "./constants.js";
import encoders from "./encoders.js";

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

function mul(x, y) {
    return encoders.make_bigfloat(
        x.coefficient * y.coefficient,
        x.exponent + y.exponent,
    );
}
