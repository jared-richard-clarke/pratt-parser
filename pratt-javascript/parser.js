import constants from "./modules/constants.js";
import scan from "./modules/lexer.js";
import utils from "./modules/utils.js";
import encoders from "./modules/big-math/encoders.js";

const parser = (function () {
    // === parser: private methods ===
    function parse_expression(rbp) {
        const token = next();
        const [prefix, ok] = table.get_parser("prefix", token.type);
        if (!ok) {
            token.message += constants.NO_PREFIX;
            return [null, token];
        }
        let [x, error] = prefix(token);
        if (error !== null) {
            return [null, error];
        }
        while (rbp < table.get_binding("bind", peek())) {
            const token = next();
            const [infix, ok] = table.get_parser("infix", token.type);
            if (!ok) {
                token.message += constants.NO_INFIX;
                return [null, token];
            }
            [x, error] = infix(x, token);
            if (error !== null) {
                return [null, error];
            }
        }
        return [x, null];
    }

    function parse_eof(token) {
        if (state.length === 1) {
            // If the expression is empty, then the error spans it.
            token.length = token.column;
            token.column = 0;
            token.message += constants.EMPTY_EXPRESSION;
            return [null, token];
        }
        token.message += constants.INCOMPLETE_EXPRESSION;
        return [null, token];
    }

    function parse_unary_error(token) {
        return [null, token];
    }
    function parse_binary_error(x, token) {
        x = null;
        return [x, token];
    }

    function parse_number(token) {
        const number = encoders.decode(token.value);
        return [number, null];
    }

    function parse_unary(token) {
        const bind = table.get_binding("prebind", token.type);
        const [x, error] = parse_expression(bind);
        if (error !== null) {
            return [null, error];
        }
        const operation = utils.unary_operation[token.type];
        return [operation(x), null];
    }

    function parse_binary(left) {
        return function (x, token) {
            const bind = table.get_binding("bind", token.type);
            const [y, error] = parse_expression(left ? bind : bind - 1);
            if (error !== null) {
                return [null, error];
            }
            const operation = utils.binary_operation[token.type];
            const value = operation(x, y);
            if (typeof value === "string") {
                token.message += value;
                return [null, token];
            }
            return [value, null];
        };
    }
    const parse_left = parse_binary(true);
    const parse_right = parse_binary(false);

    function parse_grouping(token) {
        const [x, error] = parse_expression(0);
        if (error !== null) {
            return [null, error];
        }
        if (!match(constants.CLOSE_PAREN)) {
            token.message += constants.MISMATCHED_PAREN;
            return [null, token];
        }
        next();
        return [x, null];
    }

    function next() {
        if (state.index >= state.end) {
            return state.source[state.end];
        }
        const token = state.source[state.index];
        state.index += 1;
        return token;
    }

    function peek() {
        return state.source[state.index].type;
    }

    function match(expect) {
        return peek() === expect;
    }

    function flush_errors(error) {
        const errors = [error];
        while (state.index < state.end) {
            const token = next();
            if (token.type === "error") {
                errors.push(token);
            }
        }
        return errors;
    }

    // === parser: lookup table ===
    const table = (function () {
        function register(bind, type, parser) {
            registry.prefix[type] = parser;
            registry.bind[type] = bind;
        }

        function register_unary(bind, types, parser) {
            types.forEach((type) => {
                registry.prefix[type] = parser;
                registry.prebind[type] = bind;
            });
        }

        function register_binary(bind, types, parser) {
            types.forEach((type) => {
                registry.infix[type] = parser;
                registry.bind[type] = bind;
            });
        }

        // === table: internal state ===
        const registry = {
            prefix: {},
            infix: {},
            bind: {},
            prebind: {},
        };

        register(0, constants.EOF, parse_eof);
        register(0, constants.NUMBER, parse_number);
        register(0, constants.OPEN_PAREN, parse_grouping);
        register_binary(10, [
            constants.ADD,
            constants.SUBTRACT,
            constants.SUBTRACT_ALT,
        ], parse_left);
        register_binary(
            20,
            [
                constants.MULTIPLY,
                constants.MULTIPLY_ALT,
                constants.DIVIDE,
                constants.DIVIDE_ALT,
            ],
            parse_left,
        );
        register_binary(30, [constants.EXPONENT], parse_right);
        register_binary(40, [constants.IMPLIED_MULTIPLY], parse_left);
        register_unary(50, [
            constants.ADD,
            constants.SUBTRACT,
            constants.SUBTRACT_ALT,
        ], parse_unary);
        register(60, constants.ERROR, parse_unary_error);
        register_binary(60, [constants.ERROR], parse_binary_error);

        // === table: public methods ===
        const methods = Object.create(null);

        methods.get_parser = function (category, type) {
            const parser = registry[category][type];
            return parser === undefined ? [null, false] : [parser, true];
        };

        methods.get_binding = function (category, type) {
            return registry[category][type];
        };

        return Object.freeze(methods);
    })();

    // parser: internal state
    const state = {
        source: [],
        length: 0,
        index: 0,
        end: 0,
    };

    // === parser: public methods ===
    const methods = Object.create(null);

    methods.input = function (text) {
        const tokens = scan(text);
        state.source = tokens;
        state.length = tokens.length;
        state.index = 0;
        state.end = tokens.length - 1;
        return methods;
    };

    methods.run = function () {
        const [x, error] = parse_expression(0);
        if (error !== null) {
            const errors = flush_errors(error);
            return [null, errors];
        }
        // Check for unused tokens.
        if (state.index < state.end) {
            const token = next();
            token.message += constants.INCOMPLETE_EXPRESSION;
            const errors = flush_errors(token);
            return [null, errors];
        }
        // Use scientific notation for exceedingly large or small numbers.
        if (utils.is_exceeding(x)) {
            return [encoders.encode_scientific(x), null];
        }
        return [encoders.encode(x), null];
    };

    return Object.freeze(methods);
})();

export default function parse(text) {
    return parser.input(text).run();
}
